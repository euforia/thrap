package storage

import (
	"fmt"
	"path/filepath"

	"github.com/euforia/thrap/thrapb"
	"github.com/hashicorp/consul/api"
)

type ConsulDeployDescStorage struct {
	client *api.Client
	// This should be {thrap}/deployment
	prefix string
}

func NewConsulDeployDescStorage(conf *api.Config, prefix string) (*ConsulDeployDescStorage, error) {
	client, err := newConsulClient(conf)
	if err != nil {
		return nil, err
	}

	s := &ConsulDeployDescStorage{
		prefix: prefix,
		client: client,
	}

	return s, nil
}

func (s *ConsulDeployDescStorage) Get(projectID string) (*thrapb.DeploymentDescriptor, error) {
	key := s.keyPath(projectID)
	kv := s.client.KV()
	kvp, _, err := kv.Get(key, &api.QueryOptions{})
	if err != nil {
		return nil, err
	}
	if kvp == nil {
		return nil, fmt.Errorf("deployment descriptor not found: %s", projectID)
	}

	var desc thrapb.DeploymentDescriptor
	err = desc.Unmarshal(kvp.Value)
	return &desc, err
}
func (s *ConsulDeployDescStorage) Set(projectID string, desc *thrapb.DeploymentDescriptor) error {
	key := s.keyPath(projectID)
	val, err := desc.Marshal()
	if err != nil {
		return err
	}

	kv := s.client.KV()
	return putKV(kv, key, val)
}

func (s *ConsulDeployDescStorage) Delete(projectID string) error {
	key := s.keyPath(projectID)
	kv := s.client.KV()
	_, err := kv.Delete(key, &api.WriteOptions{})
	return err
}

func (s *ConsulDeployDescStorage) keyPath(id string) string {
	return filepath.Join(s.prefix, id, "descriptor")
}
