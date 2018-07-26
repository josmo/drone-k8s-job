package main
//TODO: This needs to be simplified a ton!
//Just the initial hack
import (
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	apiv1 "k8s.io/api/core/v1"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	"github.com/pkg/errors"
	"flag"
	"k8s.io/apimachinery/pkg/labels"
	v12 "k8s.io/client-go/kubernetes/typed/batch/v1"
	"bufio"
)

type (
	Config struct {
	URL             string
	Token           string
	Insecure        bool
	Namespace      string
	Template     string
	Cleanup      bool
	Timeout      int64
	Tail         bool
	Clientset    *kubernetes.Clientset
}
	Build struct {
		Tag     string
		Event   string
		Number  int
		Commit  string
		Ref     string
		Branch  string
		Author  string
		Status  string
		Link    string
		Started int64
		Created int64
	}
	Job struct {
		Started int64
	}
	Repo struct {
		Owner string
		Name  string
	}
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
		Job    Job
	}
)

func (p Plugin) Exec() error {
	log.Info("Drone k8s deployment plugin")


	if p.Config.URL == "" {
		return errors.New("eek: Must have the server url")
	}
	if p.Config.Token == "" || len(p.Config.Template) <= 0  {
		return errors.New("eek: Must have a Token")
	}
	if p.Config.Template == ""  {
		return errors.New("eek: Must have a Template")
	}
	// parse the template file and do substitutions
	txt, err := openAndSub(p.Config.Template, p)
	if err != nil {
		return err
	}
	json, err := utilyaml.ToJSON([]byte(txt))
	if err != nil {
		return err
	}

	var dep v1.Job
	scheme := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(scheme)
	e := runtime.DecodeInto(codecs.UniversalDecoder(), json, &dep)
	if e != nil {
		log.Fatal("Error decoding yaml file to json", e)
	}

	config := &rest.Config{
		Host:            p.Config.URL,
		BearerToken:     p.Config.Token,
		TLSClientConfig: rest.TLSClientConfig{Insecure: p.Config.Insecure},
	}
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return err
	}
	p.Config.Clientset = clientset

	jobClient := p.Config.Clientset.BatchV1().Jobs(p.Config.Namespace)
	job, err := jobClient.Create(&dep)
	if err != nil {
		return err
	}

	watchTimeoute := int64(p.Config.Timeout)
	var field string
	flag.StringVar(&field, "metadata.name", "", "Field selector")
	flag.Parse()
	labelSelect := labels.Set(job.Spec.Selector.MatchLabels)
	listOptions := metav1.ListOptions{
		LabelSelector:  labelSelect.String(),
		TimeoutSeconds: &watchTimeoute,
	}

	watcher, watchErr := jobClient.Watch(listOptions)
	if watchErr != nil {
		return watchErr
	}

	ch := watcher.ResultChan()
	for {
		event, ok := <-ch
		log.Printf("Service Event %v: %+v", event.Type, event.Object)
		if event.Type == "ADDED" {
			p.logJobContainers(job.Name,jobClient, labelSelect)
		} else if event.Type == "MODIFIED" {
			pods, _ := p.getPods(p.Config.Namespace, labelSelect.AsSelector().String())
			for _, pod := range pods.Items {
			if pod.Status.Phase == "Failed" {
				delErr := deleteJob(job.Name, jobClient, p.Config.Cleanup)
				if delErr != nil {
					return delErr
				}
				return errors.New("This job failed!!!")
			}
			delErr := deleteJob(job.Name, jobClient, p.Config.Cleanup)
			if delErr != nil {
				return delErr
			}
			}
			break
		} else if event.Type == "DELETED" {
			return errors.New("Job was delete out of band")
			//break
		} else if event.Type == "" {
			return errors.New("Service watch timed out")
		} else {
			return errors.New("Service watch unhandled event")
		}
		if !ok {
			break
		}
	}
	log.Info("job created ", job)
	return err
}
func(p Plugin) logJobContainers(name string, jobClient v12.JobInterface, labelSelect labels.Set) error {
	pods, _ := p.getPods(p.Config.Namespace, labelSelect.AsSelector().String())

	watcher, watchErr := p.Config.Clientset.CoreV1().Pods(p.Config.Namespace).Watch(metav1.ListOptions{LabelSelector:labelSelect.String()})
	if watchErr != nil {
		return watchErr
	}
	ch := watcher.ResultChan()
	for {
		event, ok := <-ch
		log.Printf("Service Event %v: %+v", event.Type, event.Object)
		if event.Type == "ADDED" {
			//c.newService(event.Object)
		} else if event.Type == "MODIFIED" {
			for _, pod := range pods.Items {
				req := p.Config.Clientset.CoreV1().Pods(p.Config.Namespace).GetLogs(pod.Name, &apiv1.PodLogOptions{
					Follow: true,
				})
				readCloser, err := req.Stream()
				if err != nil {
					return err
				}
				defer func() {
					_ = readCloser.Close()
					return
				}()
				scanner := bufio.NewScanner(readCloser)
				for scanner.Scan() {
					log.Info("Job Log: ", scanner.Text())
				}
			}
			log.Info("pods", pods)
			break
		} else if event.Type == "DELETED" {
			return errors.New("Job was delete out of band")
			//break
		} else if event.Type == "" {
			return errors.New("Service watch timed out")
		} else {
			return errors.New("Service watch unhandled event")
		}
		if !ok {
			break
		}
	}
	return nil
}

func deleteJob(name string, jobClient v12.JobInterface, cleanup bool) error {
	if cleanup {
		delProp := metav1.DeletionPropagation(metav1.DeletePropagationForeground)
		deleteErr := jobClient.Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &delProp,
		})
		return deleteErr
	}
	return nil
}


// open up the template and then sub variables in. Handlebar stuff.
func openAndSub(templateFile string, p Plugin) (string, error) {
	t, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return "", err
	}
	//potty humor!  Render trim toilet paper!  Ha ha, so funny.
	return RenderTrim(string(t), p)
}


//we should really do a watch, but this is fine for now
func(p Plugin) getPods(namespace string, selector string) (*apiv1.PodList, error) {
	pods, err := p.Config.Clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to list pods")
	}
	return pods, nil
}
