package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	apiv1 "k8s.io/api/core/v1"
)

var (
	version = "0.0.0"
	build   = "0"
)

func main() {
	app := cli.NewApp()
	app.Name = "drone k8s job"
	app.Usage = "drone k8s job"
	app.Action = run
	app.Version = fmt.Sprintf("%s+%s", version, build)

	app.Flags = []cli.Flag{

		cli.StringFlag{
			Name:   "url",
			Usage:  "url to the k8s api",
			EnvVar: "PLUGIN_URL, KUBERNETES_URL",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "kubernetes token",
			EnvVar: "PLUGIN_TOKEN, KUBERNETES_TOKEN",
		},
		cli.StringFlag{
			Name:   "ca",
			Usage:  "Certificate Authority file encoded into base64: e.g: run: `cat ca.pem | base64` to get this value",
			EnvVar: "PLUGIN_CA,KUBERNETES_CA",
		},
		cli.BoolFlag{
			Name:   "insecure",
			Usage:  "Insecure connection",
			EnvVar: "PLUGIN_INSECURE",
		},
		cli.StringFlag{
			Name:   "namespace",
			Usage:  "namespace for the job",
			Value:  apiv1.NamespaceDefault,
			EnvVar: "PLUGIN_NAMESPACE",
		},
		cli.StringFlag{
			Name:   "template",
			Usage:  "template file to use for deployment: mydeployment.yaml :-)",
			EnvVar: "JOB_TEMPLATE,PLUGIN_TEMPLATE",
		},
		cli.BoolTFlag{
			Name:   "cleanup",
			Usage:  "Will delete the job after running to success",
			EnvVar: "PLUGIN_CLEANUP",
		},
		cli.Int64Flag{
			Name:   "timeout",
			Usage:  "How long will the service listen till it times out",
			EnvVar: "PLUGIN_TIMEOUT",
			Value:  120,
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "Enable debugging messages",
			EnvVar: "PLUGIN_DEBUG",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.Int64Flag{
			Name:   "build.started",
			Usage:  "build started",
			EnvVar: "DRONE_BUILD_STARTED",
		},
		cli.Int64Flag{
			Name:   "build.created",
			Usage:  "build created",
			EnvVar: "DRONE_BUILD_CREATED",
		},
		cli.StringFlag{
			Name:   "build.tag",
			Usage:  "build tag",
			EnvVar: "DRONE_TAG",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Tag:     c.String("build.tag"),
			Number:  c.Int("build.number"),
			Event:   c.String("build.event"),
			Status:  c.String("build.status"),
			Commit:  c.String("commit.sha"),
			Ref:     c.String("commit.ref"),
			Branch:  c.String("commit.branch"),
			Author:  c.String("commit.author"),
			Link:    c.String("build.link"),
			Started: c.Int64("build.started"),
			Created: c.Int64("build.created"),
		},
		Job: Job{
			Started: c.Int64("job.started"),
		},
		Config: Config{
			URL:       c.String("url"),
			Token:     c.String("token"),
			Insecure:  c.Bool("insecure"),
			Ca:        c.String("ca"),
			Namespace: c.String("namespace"),
			Template:  c.String("template"),
			Cleanup:   c.BoolT("cleanup"),
			Timeout:   c.Int64("timeout"),
			Debug:     c.Bool("debug"),
		},
	}
	return plugin.Exec()
}
