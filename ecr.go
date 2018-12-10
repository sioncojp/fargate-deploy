package fargatedeploy

import (
	"fmt"

	"encoding/base64"

	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// ECR ... Store Session for ECR.
type ECR struct {
	Session *ecr.ECR
}

// EcrLogin ... Login ECR.
func (m *Material) EcrLogin() error {
	ecr := &ECR{}
	err := ecr.NewSession(m.AWSProfile)
	if err != nil {
		return err
	}

	token, err := ecr.GetAuthorizationToken()
	if err != nil {
		return err
	}

	auth, err := NewECRAuth(token)
	if err != nil {
		return err
	}

	if err := DockerLogin(auth); err != nil {
		return err
	}

	return nil
}

// GetAuthorizationToken ... Retrieves a token that is valid for a specified registry for 12 hours for ECR.
func (e *ECR) GetAuthorizationToken() (*ecr.GetAuthorizationTokenOutput, error) {
	input := &ecr.GetAuthorizationTokenInput{
		RegistryIds: GetRegistryIds(),
	}

	result, err := e.Session.GetAuthorizationToken(input)
	if err != nil {
		return nil, fmt.Errorf("authorizing: %s", err)
	}
	return result, nil
}

// GetRegistryIds ... get list of registries from env, leave empty for default
func GetRegistryIds() []*string {
	var registryIds []*string
	registryIds = append(registryIds, aws.String(EcrID))

	return registryIds
}

// NewECRAuth ... New creates a auth data of the ECR to login.
func NewECRAuth(token *ecr.GetAuthorizationTokenOutput) (ECRAuth, error) {
	auth := token.AuthorizationData[0]
	data, err := base64.StdEncoding.DecodeString(*auth.AuthorizationToken)
	if err != nil {
		return ECRAuth{}, fmt.Errorf("decode to base64: %s", err)
	}
	// extract username and password
	t := strings.SplitN(string(data), ":", 2)

	// object to pass to template
	a := ECRAuth{
		Token:         *auth.AuthorizationToken,
		User:          t[0],
		Pass:          t[1],
		ProxyEndpoint: *(auth.ProxyEndpoint),
		ExpiresAt:     *(auth.ExpiresAt),
	}

	return a, nil
}

// NewSession ... New creates a new instance of the ECR client with a session.
func (e *ECR) NewSession(awsProfile string) error {
	sess, err := NewSession(awsProfile)
	if err != nil {
		return err
	}
	e.Session = ecr.New(sess, aws.NewConfig().WithMaxRetries(10).WithRegion(AWSRegion))
	return nil
}
