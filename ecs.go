package fargatedeploy

import (
	"fmt"

	"log"

	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// ECS ... Store Session for ECS.
type ECS struct {
	Session *ecs.ECS
}

// ECSServiceUpdate ... Register a new task and update service to it.
func (m *Material) ECSServiceUpdate() error {
	ecs := &ECS{}
	err := ecs.NewSession(m.AWSProfile)
	if err != nil {
		return err
	}

	containerDefinition := ecs.NewContainerDefinition(m)
	taskRoleArn, err := ecs.ECSRegisterTaskDefinition(m, containerDefinition)
	if err != nil {
		return err
	}

	if err := ecs.UpdateService(m, taskRoleArn); err != nil {
		return err
	}
	return nil
}

// NewContainerDefinition ... Initialize Container definitions.
func (e *ECS) NewContainerDefinition(m *Material) []*ecs.ContainerDefinition {
	var result []*ecs.ContainerDefinition
	var workingdir *string
	var env []*ecs.KeyValuePair
	var entrypoint []*string
	var port []*ecs.PortMapping
	var image string
	var spec Spec

	for _, container := range m.Containers {

		if container.EntryPoint != "" {
			entrypoint = append(entrypoint, aws.String(container.EntryPoint))
		}

		if container.WorkingDirectory != "" {
			workingdir = aws.String(container.WorkingDirectory)
		}

		for _, v := range container.Environments {
			env = append(env, &ecs.KeyValuePair{
				Name:  aws.String(v.Name),
				Value: aws.String(v.Value),
			})
		}

		// Set Image Name
		image = fmt.Sprintf("%s:%s", container.Image, Env)
		if container.Ecr && EcrID != "" {
			image = fmt.Sprintf(ecrImageFormat, EcrID, AWSRegion, container.Image, Env)
		}

		// set CPU, Memory
		spec.CPU = container.CPU
		spec.Memory = container.Memory
		for _, v := range container.Specs {
			if spec.CPU == 0 {
				spec.CPU = v.CPU
			}

			if spec.Memory == 0 {
				spec.Memory = v.Memory
			}
		}

		cd := &ecs.ContainerDefinition{
			Name:      aws.String(container.Name),
			Cpu:       aws.Int64(spec.CPU),
			Memory:    aws.Int64(spec.Memory),
			Essential: aws.Bool(true),
			Image: aws.String(
				image,
			),
			EntryPoint:       entrypoint,
			WorkingDirectory: workingdir,
			Environment:      env,
			LogConfiguration: &ecs.LogConfiguration{
				LogDriver: aws.String("awslogs"),
				Options:   e.NewLogOption(m),
			},
			LinuxParameters: &ecs.LinuxParameters{InitProcessEnabled: aws.Bool(true)},
			Ulimits:         e.NewUlimit(m),
		}

		// set portMapping
		if container.Port != 0 {
			port = append(port, &ecs.PortMapping{ContainerPort: aws.Int64(container.Port)})
			cd.SetPortMappings(port)
		}

		// set command
		if len(container.Command.Value) != 0 {
			cd.SetCommand(aws.StringSlice(container.Command.Value))
		}

		result = append(result, cd)

	}

	return result
}

// ECSRegisterTaskDefinition ... Registers a new task definition from the supplied task(family) and containerDefinitions.
func (e *ECS) ECSRegisterTaskDefinition(m *Material, cd []*ecs.ContainerDefinition) (*string, error) {
	var compatibilities []*string
	var cpu int64
	var memory int64

	compatibilities = append(compatibilities, aws.String("FARGATE"))

	if m.ExecutionRoleArn == "" {
		m.ExecutionRoleArn = m.RoleArn
	}

	// Add the values ​​of all cpu, memory
	for _, v := range cd {
		cpu += *v.Cpu
		memory += *v.Memory
	}

	if m.Family == "" {
		m.Family = m.ECSCluster
	}

	input := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    cd,
		Family:                  aws.String(m.Family),
		TaskRoleArn:             aws.String(m.RoleArn),
		ExecutionRoleArn:        aws.String(m.ExecutionRoleArn),
		RequiresCompatibilities: compatibilities,
		NetworkMode:             aws.String("awsvpc"),
		Cpu:                     aws.String(strconv.FormatInt(cpu, 10)),
		Memory:                  aws.String(strconv.FormatInt(memory, 10)),
	}

	result, err := e.Session.RegisterTaskDefinition(input)
	if err != nil {
		return nil, fmt.Errorf("register task: %s", err)
	}
	log.Printf("registered task: %s\n", *result.TaskDefinition.TaskDefinitionArn)

	return result.TaskDefinition.TaskDefinitionArn, nil
}

// UpdateService ... Update ECS Service.
func (e *ECS) UpdateService(m *Material, taskRoleArn *string) error {
	input := &ecs.UpdateServiceInput{
		Cluster:        aws.String(m.ECSCluster),
		Service:        aws.String(Service),
		TaskDefinition: taskRoleArn,
	}
	result, err := e.Session.UpdateService(input)
	if err != nil {
		return fmt.Errorf("update service: %s", err)
	}
	log.Printf("updated service: %s\n", *result.Service.ServiceArn)

	return nil
}

// NewSession ... New creates a new instance of the ECS client with a session.
func (e *ECS) NewSession(awsProfile string) error {
	sess, err := NewSession(awsProfile)
	if err != nil {
		return err
	}
	e.Session = ecs.New(sess, aws.NewConfig().WithMaxRetries(10).WithRegion(AWSRegion))
	return nil
}

// NewLogOption ... New creates a log option.
func (e *ECS) NewLogOption(m *Material) map[string]*string {
	group := m.ECSCluster
	prefix := m.ECSCluster

	if m.LogGroup != "" {
		group = m.LogGroup
	}

	if m.LogPrefix != "" {
		prefix = m.LogPrefix
	}

	result := map[string]*string{
		"awslogs-group":         aws.String(group),
		"awslogs-region":        aws.String(AWSRegion),
		"awslogs-stream-prefix": aws.String(prefix),
	}
	return result
}

// NewUlimit ... New creates a ulimit option.
func (e *ECS) NewUlimit(m *Material) []*ecs.Ulimit {
	var ulimit []*ecs.Ulimit

	u := &ecs.Ulimit{
		HardLimit: aws.Int64(65536),
		SoftLimit: aws.Int64(65536),
		Name:      aws.String("nofile"),
	}

	ulimit = append(ulimit, u)
	return ulimit
}
