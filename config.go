package fargatedeploy

import (
	"time"

	"github.com/BurntSushi/toml"
)

const (
	appName              = "fargate-deploy"
	ecrImageFormat       = "%s.dkr.ecr.%s.amazonaws.com/%s:%s"
	ecrSetImageTagFormat = "%s.dkr.ecr.%s.amazonaws.com/%s"
)

var (
	// AWSRegion ... Set aws region.
	AWSRegion string

	// Service ... Set ECS service name.
	Service string

	// BuildImage ... Set build image name from Dockerfile.
	BuildImage string

	// EcrID ... Set a ID for logging in to ECR.
	EcrID string

	// Env ... Set a Env for choosing env from config.
	Env string

	// Lifecycle ... Lifecycle docker image in target repository.
	Lifecycle bool

	// CommitHash ... Set a git commit hash in current directroy.
	CommitHash string
)

// AWS ... Reacts to the aws NewSession.
type AWS interface {
	NewSession() error
}

// Config ... Store data from loaded config
type Config struct {
	AWSRegion    string                   `toml:"aws_region"`
	Service      string                   `toml:"service"`
	BuildImage   string                   `toml:"build_image"`
	Materials    map[string]Material      `toml:"material"`
	Containers   map[string]Container     `toml:"container"`
	Specs        map[string][]Spec        `toml:"Spec"`
	Environments map[string][]Environment `toml:"environment"`
	Commands     map[string][]Command     `toml:"command"`
	Secrets      map[string][]Secret      `toml:"secret"`
}

// Material ... Register task and store the information necessary for updating the service.
type Material struct {
	AWSProfile       string `toml:"aws_profile"`
	ECSCluster       string `toml:"ecs_cluster"`
	RoleArn          string `toml:"role_arn"`
	ExecutionRoleArn string `toml:"execution_role_arn"`
	LogGroup         string `toml:"log_group"`
	LogPrefix        string `toml:"log_prefix"`
	Family           string `toml:"family"`
	Containers       []Container
}

// Container ... Store a container data.
type Container struct {
	Name             string
	Image            string
	Ecr              bool
	CPU              int64
	Memory           int64
	Port             int64
	EntryPoint       string
	WorkingDirectory string
	Specs            []Spec
	Environments     []Environment
	Secrets          []Secret
	Command          Command
}

// Spec ... Add cpu, memory variables for container.
type Spec struct {
	CPU    int64
	Memory int64
}

// Secret ... Add secret variables for container.
type Secret struct {
	Name  string
	Value string
}

// Environment ... Add environment variables for container.
type Environment struct {
	Name  string
	Value string
}

// Command ... Add command for container.
type Command struct {
	Value []string
}

// ECRAuth ... Store a auth data of the ECR to login.
type ECRAuth struct {
	Token         string
	User          string
	Pass          string
	ProxyEndpoint string
	ExpiresAt     time.Time
}

// LoadToml ... Load toml file
func LoadToml(c string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(c, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
