package fargatedeploy

import (
	"log"
	"os"
	"strings"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetPrefix(appName + ": ")
}

// NewSpecs ...  Initialize by parse Specs in config
// Output : map[app:map[CPU:1024 Memory:2048], db:map[CPU:256 Memory:512]]
func (c *Config) NewSpecs(env string) map[string][]Spec {
	specs := make(map[string][]Spec)

	for i, e := range c.Specs {
		for _, ee := range e {
			containerName := strings.Split(i, "_")

			if len(containerName) == 1 || strings.HasSuffix(i, env) {
				specs[containerName[0]] = append(specs[containerName[0]], Spec{ee.CPU, ee.Memory})
			}
		}
	}
	return specs
}

// NewSecrets ...  Initialize by parse Secrets in config
// Output : map[app:map[CPU:1024 Memory:2048], db:map[CPU:256 Memory:512]]
func (c *Config) NewSecrets(env string) map[string][]Secret {
	secrets := make(map[string][]Secret)

	for i, e := range c.Secrets {
		for _, ee := range e {
			containerName := strings.Split(i, "_")

			if len(containerName) == 1 || strings.HasSuffix(i, env) {
				secrets[containerName[0]] = append(secrets[containerName[0]], Secret{ee.Name, ee.Value})
			}
		}
	}
	return secrets
}

// NewEnvironments ...  Initialize by parse Environments in config
// Output : map[app:map[NODE_ENV:production PORT:3000], db:map[PASSWORD:xxxx]]
func (c *Config) NewEnvironments(env string) map[string][]Environment {
	envs := make(map[string][]Environment)

	for i, e := range c.Environments {
		for _, ee := range e {
			containerName := strings.Split(i, "_")

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
			containerName := strings.Split(i, "_")
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
			nil,
			nil,
			Command{},
		}
		containers = append(containers, v)
	}
	return containers
}
