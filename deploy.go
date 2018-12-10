package fargatedeploy

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

// Deploy ... (build image) and service update for Fargate.
func (m *Material) Deploy() error {
	if BuildImage != "" {
		if EcrID != "" {
			if err := m.EcrLogin(); err != nil {
				return err
			}
		}
		c, err := GitGetCommitHash()
		if err != nil {
			return err
		}
		CommitHash = c

		// TODO: Implementation
		if Lifecycle {
			if err := ImageLifecycle(); err != nil {
				return err
			}
		}

		if err := DockerImageUpload(); err != nil {
			return err
		}
	}

	if err := m.ECSServiceUpdate(); err != nil {
		return err
	}
	return nil
}

// NewSession ... Creates a new instance of the ECS client with a session.
func NewSession(awsProfile string) (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           awsProfile,
	})

	return sess, err
}
