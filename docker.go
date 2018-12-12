package fargatedeploy

import (
	"fmt"
	"os/exec"

	"log"
)

// DockerImageUpload ... Docker build and tagging and push image
func DockerImageUpload() error {
	images := NewDockerImage()

	if err := DockerBuild(images); err != nil {
		return err
	}

	if err := DockerSetTag(images); err != nil {
		return err
	}

	if err := DockerPush(images); err != nil {
		return err
	}

	return nil
}

// DockerBuild ... Build an image from a Dockerfile
func DockerBuild(images []string) error {
	if out, err := exec.Command(
		"docker",
		"build",
		"--pull",
		"-t",
		images[0],
		".",
	).CombinedOutput(); err != nil {
		return fmt.Errorf("docker build: %s: %s", err, string(out))
	}
	log.Printf("docker build done: %s\n", images[0])
	return nil
}

// DockerPush ...Push an image or a repository to a registry
func DockerPush(images []string) error {
	for _, i := range images {
		if out, err := exec.Command(
			"docker",
			"push",
			i,
		).CombinedOutput(); err != nil {
			return fmt.Errorf("docker push: %s: %s", err, string(out))
		}
		log.Printf("docker push done: %s\n", i)
	}

	return nil
}

// DockerSetTag ... Create a tag TARGET_IMAGE that refers to SOURCE_IMAGE
func DockerSetTag(images []string) error {
	if out, err := exec.Command(
		"docker",
		"tag",
		images[0],
		images[1],
	).CombinedOutput(); err != nil {
		return fmt.Errorf("docker tag: %s: %s", err, string(out))
	}

	// For prod, tagging with latest
	if Env == "prod" {
		if out, err := exec.Command(
			"docker",
			"tag",
			images[0],
			images[2],
		).CombinedOutput(); err != nil {
			return fmt.Errorf("docker tag: %s: %s", err, string(out))
		}
	}

	return nil
}

// DockerLogin ... login ECR repository
func DockerLogin(auth ECRAuth) error {
	out, err := exec.Command(
		"docker",
		"login",
		"-u",
		auth.User,
		"-p",
		auth.Pass,
		auth.ProxyEndpoint,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker login: %s: %s", err, string(out))
	}

	return nil
}

// NewDockerImage ... New creates docker image names.
func NewDockerImage() []string {
	var images []string
	if EcrID != "" {
		images = []string{
			fmt.Sprintf(ecrImageFormat, EcrID, AWSRegion, BuildImage, CommitHash),
			fmt.Sprintf(ecrImageFormat, EcrID, AWSRegion, BuildImage, Env),
		}
		if Env == "prod" {
			images = append(images, fmt.Sprintf(ecrImageFormat, EcrID, AWSRegion, BuildImage, string("latest")))
		}
	} else {
		images = []string{
			fmt.Sprintf("%s:%s", BuildImage, CommitHash),
			fmt.Sprintf("%s:%s", BuildImage, Env),
		}
		if Env == "prod" {
			images = append(images, fmt.Sprintf("%s:latest", BuildImage))
		}
	}

	return images
}
