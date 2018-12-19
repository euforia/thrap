package storage

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/euforia/thrap/pkg/pb"
	"github.com/hashicorp/consul/api"
)

const (
	// DefaultSpecVersion defines the default spec version
	DefaultSpecVersion = "default"
)

var (
	errMimeTypeRequired = errors.New("descriptor mime type required")
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
	key := s.versionPath(projectID, DefaultSpecVersion)
	kv := s.client.KV()
	kvp, _, err := kv.Get(key, &api.QueryOptions{})
	if err != nil {
		return nil, err
	}
	if kvp == nil {
		return nil, errors.New("default deployment descriptor not found: " + projectID)
	}

	var desc pb.DeploymentDescriptor
	err = desc.Unmarshal(kvp.Value)
	return &desc, err
}

// GetVersion satisfies the DeployDescStorage interface
func (s *ConsulDeployDescStorage) GetVersion(projectID string, version string) (*pb.DeploymentDescriptor, error) {
	key := s.versionPath(projectID, version)

	kv := s.client.KV()
	kvp, _, err := kv.Get(key, &api.QueryOptions{})
	if err != nil {
		return nil, err
	}
	if kvp == nil {
		return nil, fmt.Errorf("deployment descriptor not found project=%s version=%s",
			projectID, version)
	}

	var desc pb.DeploymentDescriptor
	err = desc.Unmarshal(kvp.Value)
	return &desc, err
}

// Set satisfies the DeployDescStorage interface
func (s *ConsulDeployDescStorage) Set(projectID string, desc *pb.DeploymentDescriptor) error {
	if len(desc.Mime) == 0 {
		return errMimeTypeRequired
	}
	if len(desc.ID) == 0 {
		desc.ID = DefaultSpecVersion
	}

	key := s.versionPath(projectID, desc.ID)
	val, err := desc.Marshal()
	if err != nil {
		return err
	}

	kv := s.client.KV()
	return putKV(kv, key, val)
}

// Delete satisfies the DeployDescStorage interface
func (s *ConsulDeployDescStorage) Delete(projectID string) error {
	key := s.versionPath(projectID, DefaultSpecVersion)
	kv := s.client.KV()
	_, err := kv.Delete(key, &api.WriteOptions{})
	return err
}

// DeleteVersion satisfies the DeployDescStorage interface
func (s *ConsulDeployDescStorage) DeleteVersion(projectID, version string) error {
	key := s.versionPath(projectID, version)
	kv := s.client.KV()
	_, err := kv.Delete(key, &api.WriteOptions{})
	return err
}

// ListVersions satisfies the DeployDescStorage interface
func (s *ConsulDeployDescStorage) ListVersions(projectID string) ([]string, error) {
	key := s.keyPath(projectID)

	kv := s.client.KV()
	kvps, _, err := kv.List(key, &api.QueryOptions{})
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(kvps))
	if kvps == nil {
		// return nil, errors.New("no versions found: " + projectID)
		return versions, nil
	}

	for _, kvp := range kvps {
		if kvp.Key == s.keyPath(projectID) {
			continue
		}
		versions = append(versions, filepath.Base(kvp.Key))
	}

	return versions, nil
}

func (s *ConsulDeployDescStorage) versionPath(projectID, version string) string {
	return filepath.Join(s.keyPath(projectID), version)
}

func (s *ConsulDeployDescStorage) keyPath(projectID string) string {
	return filepath.Join(s.prefix, projectID)
}
