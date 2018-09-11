package registry

import (
	"encoding/base64"
	"errors"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/docker/docker/api/types"
	"github.com/euforia/thrap/pkg/provider"
)

var (
	errRepoTagNotFound = errors.New("repository tag not found")
)

const defaultECRRepoPolicy = `{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Sid": "AllowPull",
            "Effect": "Allow",
            "Principal": "*",
            "Action": [
                "ecr:GetDownloadUrlForLayer",
                "ecr:BatchGetImage",
                "ecr:BatchCheckLayerAvailability"
            ]
        }
    ]
}`

type awsContainerRegistry struct {
	sess *session.Session
	ecr  *ecr.ECR
	conf *provider.Config
}

// Envionment Variables:
// AWS_ACCESS_KEY_ID
// AWS_SECRET_ACCESS_KEY
// AWS_DEFAULT_REGION
func (ar *awsContainerRegistry) Init(rconf *provider.Config) error {
	ar.conf = rconf

	conf := aws.NewConfig()

	c := rconf.Config

	var awsCreds credentials.Value
	// Override region if supplied
	if val, ok := c["region"]; ok {
		region, ok := val.(string)
		if ok {
			conf = conf.WithRegion(region)
		}
	}

	if val, ok := c["key"]; ok {
		id, ok := val.(string)
		if ok {
			awsCreds.AccessKeyID = id
		}
	}

	if val, ok := c["secret"]; ok {
		id, ok := val.(string)
		if ok {
			awsCreds.SecretAccessKey = id
		}
	}

	creds := credentials.NewStaticCredentialsFromCreds(awsCreds)
	conf = conf.WithCredentials(creds)

	sess, err := session.NewSession(conf)
	if err == nil {
		ar.sess = sess
		ar.ecr = ecr.New(sess)
	}

	return err
}

func (ar *awsContainerRegistry) ID() string {
	return ar.conf.ID
}

func (ar *awsContainerRegistry) GetImageManifest(name, tag string) (interface{}, error) {
	imageID := &ecr.ImageIdentifier{}
	imageID.SetImageTag(tag)

	mediaType := "application/vnd.docker.distribution.manifest.v2+json"
	getImgReq := &ecr.BatchGetImageInput{
		AcceptedMediaTypes: []*string{&mediaType},
	}
	getImgReq.SetRepositoryName(name)
	getImgReq.SetImageIds([]*ecr.ImageIdentifier{imageID})

	resp, err := ar.ecr.BatchGetImage(getImgReq)
	if err != nil {
		return nil, err
	}

	if len(resp.Failures) > 0 {
		return nil, errors.New(*resp.Failures[0].FailureCode)
	}

	return resp.Images[0], nil
}

// Name returns the name used to prefix an image for ecr
func (ar *awsContainerRegistry) Name() string {
	return ar.conf.Addr
}

// ImageName returns the name prepended with the registry address delimited by /
func (ar *awsContainerRegistry) ImageName(name string) string {
	return filepath.Join(ar.conf.Addr, name)
}

func (ar *awsContainerRegistry) GetAuthConfig() (types.AuthConfig, error) {
	var auth types.AuthConfig

	i := strings.Index(ar.conf.Addr, ".")
	if i < 1 {
		return auth, errors.New("could not get registry ID from ECR address")
	}

	in := &ecr.GetAuthorizationTokenInput{
		RegistryIds: []*string{aws.String(ar.conf.Addr[:i])},
	}
	resp, err := ar.ecr.GetAuthorizationToken(in)
	if err != nil {
		return auth, err
	}

	authData := resp.AuthorizationData[0]

	data, err := base64.StdEncoding.DecodeString(*authData.AuthorizationToken)
	if err != nil {
		return auth, err
	}
	// extract username and password
	token := strings.SplitN(string(data), ":", 2)

	// object to pass to template
	auth = types.AuthConfig{
		Auth:          *authData.AuthorizationToken,
		Username:      token[0],
		Password:      token[1],
		ServerAddress: *authData.ProxyEndpoint,
	}

	return auth, nil
}

func (ar *awsContainerRegistry) GetRepo(name string) (interface{}, error) {
	input := &ecr.DescribeRepositoriesInput{
		RepositoryNames: []*string{&name},
	}

	resp, err := ar.ecr.DescribeRepositories(input)
	if err == nil {
		return resp.Repositories[0], nil
	}

	awsErr := err.(awserr.Error)
	return nil, errors.New(awsErr.Code())
}

func (ar *awsContainerRegistry) DeleteRepo(name string) (interface{}, error) {
	req := &ecr.DeleteRepositoryInput{}
	req.SetRepositoryName(name)

	resp, err := ar.ecr.DeleteRepository(req)
	if err == nil {
		return resp.Repository, nil
	}
	return nil, err
}

func (ar *awsContainerRegistry) CreateRepo(name string) (interface{}, error) {
	in := &ecr.CreateRepositoryInput{
		RepositoryName: aws.String(name),
	}
	var (
		out, err = ar.ecr.CreateRepository(in)
		repo     *ecr.Repository
	)

	if err != nil {
		return nil, err
	}

	repo = out.Repository

	// Set default read-only for everyone
	policy := &ecr.SetRepositoryPolicyInput{}
	policy.SetPolicyText(defaultECRRepoPolicy)
	policy.SetRepositoryName(name)
	_, err = ar.ecr.SetRepositoryPolicy(policy)

	return repo, err
}
