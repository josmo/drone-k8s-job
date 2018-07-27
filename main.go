package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	apiv1 "k8s.io/api/core/v1"
)

var build string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "drone k8s job"
	app.Usage = "drone k8s job"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.0+%s", build)

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
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Config: Config{
			URL:       c.String("url"),
			Token:     c.String("token"),
			Insecure:  c.Bool("insecure"),
			Namespace: c.String("namespace"),
			Template:  c.String("template"),
			Cleanup:   c.BoolT("cleanup"),
			Timeout:   c.Int64("timeout"),
			Debug:     c.Bool("debug"),
		},
	}
	return plugin.Exec()
}
