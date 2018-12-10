package fargatedeploy

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetOutput(os.Stderr)
	log.SetPrefix(appName + ": ")
}

// Run ...
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

		envs := conf.NewEnvironments(c.String("env"))
		cmds := conf.NewCommands()
		cs := conf.NewContainers()

		// deploy selected env
		for im, m := range conf.Materials {
			if im == c.String("env") {
				Env = c.String("env")
				for _, container := range cs {
					for ie, e := range envs {
						// validate container name
						if container.Name == ie {
							container.Environments = e
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

// NewEnvironments ...  Initialize by parse Environments in config
// Output : map[app:map[NODE_ENV:production PORT:3000], db:map[PASSWORD:xxxx]]
func (c *Config) NewEnvironments(env string) map[string][]Environment {
	envs := make(map[string][]Environment)

	for i, e := range c.Environments {
		for _, ee := range e {
			containerName := strings.Split(i, ".")

			if len(containerName) == 1 || strings.HasSuffix(i, env) {
				envs[containerName[0]] = append(envs[containerName[0]], Environment{ee.Name, ee.Value})
			}
		}
	}
	return envs
}

// NewCommands ...  Initialize by parse Commands in config
// Output : map[app:["go", "run", "hoge.go"], db:["mysql", "help"]]
func (c *Config) NewCommands() map[string]Command {
	cmd := make(map[string]Command)

	for i, e := range c.Commands {
		for _, ee := range e {
			containerName := strings.Split(i, ".")
			cmd[containerName[0]] = ee
		}
	}
	return cmd
}

// NewContainers ... Initialize by parse Containers in config
func (c *Config) NewContainers() []Container {
	var containers []Container
	for i, cs := range c.Containers {
		v := Container{
			i,
			cs.Image,
			cs.Ecr,
			cs.CPU,
			cs.Memory,
			cs.Port,
			cs.EntryPoint,
			cs.WorkingDirectory,
			nil,
			Command{},
		}
		containers = append(containers, v)
	}
	return containers
}
