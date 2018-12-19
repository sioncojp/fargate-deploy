package fargatedeploy

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli"
)

// Run ... Running fargate-deploy
func Run() int {
	app := FlagSet()

	app.Action = func(c *cli.Context) error {
		if c.String("config") == "" {
			return errors.New("required -c option")
		}

		if c.String("env") == "" {
			return errors.New("required -e option")
		}

		// load config
		conf, err := LoadToml(c.String("config"))
		if err != nil {
			return err
		}

		// initialize
		AWSRegion = conf.AWSRegion
		Service = conf.Service
		BuildImage = conf.BuildImage
		EcrID = c.String("id")
		Lifecycle = c.Bool("lifecycle")

		specs := conf.NewSpecs(c.String("env"))
		envs := conf.NewEnvironments(c.String("env"))
		secrets := conf.NewSecrets(c.String("env"))
		cmds := conf.NewCommands()
		cs := conf.NewContainers()

		// deploy selected env
		for im, m := range conf.Materials {
			if im == c.String("env") {
				Env = c.String("env")
				for _, container := range cs {
					for ie, e := range specs {
						// validate container name
						if container.Name == ie {
							container.Specs = e
						}
					}
					for ie, e := range envs {
						// validate container name
						if container.Name == ie {
							container.Environments = e
						}
					}
					for ie, e := range secrets {
						// validate container name
						if container.Name == ie {
							container.Secrets = e
						}
					}
					for ie, e := range cmds {
						// validate container name
						if container.Name == ie {
							container.Command = e
						}
					}
					m.Containers = append(m.Containers, container)
				}
				if err := m.Deploy(); err != nil {
					return err
				}
			}
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("failed to %s", err)
	}
	return 0
}
