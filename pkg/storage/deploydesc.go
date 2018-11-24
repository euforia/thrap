package storage

import (
	"errors"
	"path/filepath"

	"github.com/euforia/thrap/pkg/pb"
	"github.com/hashicorp/consul/api"
)

// ConsulDeployDescStorage implements a consul backed DeployDescStorage
type ConsulDeployDescStorage struct {
	client *api.Client
	// This should be {thrap}/deployment
	prefix string
}

// NewConsulDeployDescStorageFromClient returns a new ConsulDeployDescStorage instance with the given
// client
func NewConsulDeployDescStorageFromClient(client *api.Client, prefix string) *ConsulDeployDescStorage {
	return &ConsulDeployDescStorage{
		client: client,
		prefix: prefix,
	}
}

// NewConsulDeployDescStorage returns a new ConsulDeployDescStorage instance or an error
func NewConsulDeployDescStorage(conf *api.Config, prefix string) (*ConsulDeployDescStorage, error) {
	client, err := newConsulClient(conf)
	if err != nil {
		return nil, err
	}

	return NewConsulDeployDescStorageFromClient(client, prefix), nil
}

// Get satisfies the DeployDescStorage interface
func (s *ConsulDeployDescStorage) Get(projectID string) (*pb.DeploymentDescriptor, error) {
	key := s.keyPath(projectID)
	kv := s.client.KV()
	kvp, _, err := kv.Get(key, &api.QueryOptions{})
	if err != nil {
		return nil, err
	}
	if kvp == nil {
		return nil, errors.New("deployment descriptor not found: " + projectID)
	}

	var desc pb.DeploymentDescriptor
	err = desc.Unmarshal(kvp.Value)
	return &desc, err
}

// Set satisfies the DeployDescStorage interface
func (s *ConsulDeployDescStorage) Set(projectID string, desc *pb.DeploymentDescriptor) error {
	key := s.keyPath(projectID)
	val, err := desc.Marshal()
	if err != nil {
		return err
	}

	kv := s.client.KV()
	return putKV(kv, key, val)
}

// Delete satisfies the DeployDescStorage interface
func (s *ConsulDeployDescStorage) Delete(projectID string) error {
	key := s.keyPath(projectID)
	kv := s.client.KV()
	_, err := kv.Delete(key, &api.WriteOptions{})
	return err
}

func (s *ConsulDeployDescStorage) keyPath(projectID string) string {
	return filepath.Join(s.prefix, projectID, "descriptor")
}
