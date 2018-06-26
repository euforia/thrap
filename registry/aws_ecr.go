package registry

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

var (
	errRepoTagNotFound = errors.New("repository tag not found")
)

type awsContainerRegistry struct {
	sess *session.Session
	ecr  *ecr.ECR
}

// Envionment Variables:
// AWS_ACCESS_KEY_ID
// AWS_SECRET_ACCESS_KEY
// AWS_DEFAULT_REGION
func (ar *awsContainerRegistry) Init(config Config) error {

	conf := aws.NewConfig()

	c := config.Conf
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
			os.Setenv("AWS_ACCESS_KEY_ID", id)
		}
	}

	if val, ok := c["secret"]; ok {
		id, ok := val.(string)
		if ok {
			os.Setenv("AWS_SECRET_ACCESS_KEY", id)
		}
	}

	sess, err := session.NewSession(conf)
	if err == nil {
		ar.sess = sess
		ar.ecr = ecr.New(sess)
	}

	return err
}

func (ar *awsContainerRegistry) ID() string {
	return "ecr"
}

func (ar *awsContainerRegistry) Type() Type {
	return TypeContainer
}

func (ar *awsContainerRegistry) GetManifest(name, tag string) (interface{}, error) {
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
	return resp.Images[0], nil
}

func (ar *awsContainerRegistry) Get(name string) (interface{}, error) {
	input := &ecr.DescribeRepositoriesInput{
		RepositoryNames: []*string{&name},
	}

	resp, err := ar.ecr.DescribeRepositories(input)
	if err == nil {
		return resp.Repositories[0], nil
	}
	return nil, err
}

func (ar *awsContainerRegistry) Delete(name string) (interface{}, error) {
	req := &ecr.DeleteRepositoryInput{}
	req.SetRepositoryName(name)

	resp, err := ar.ecr.DeleteRepository(req)
	if err == nil {
		return resp.Repository, nil
	}
	return nil, err
}

func (ar *awsContainerRegistry) Create(name string) (interface{}, error) {
	in := new(ecr.CreateRepositoryInput)
	in.SetRepositoryName(name)

	var (
		out, err = ar.ecr.CreateRepository(in)
		repo     *ecr.Repository
	)

	if err == nil {
		repo = out.Repository
	}

	return repo, err
}
