package storage

import (
	"fmt"
	"path/filepath"

	"github.com/euforia/thrap/pkg/pb"
	"github.com/hashicorp/consul/api"
)

// ConsulDeployStorage implements a consul backed DeployStorage interface
type ConsulDeployStorage struct {
	client *api.Client
	// thrap/deployment /{project}/instance/{}/{}
	prefix string
}

func NewConsulDeployStorageFromClient(client *api.Client, prefix string) *ConsulDeployStorage {
	return &ConsulDeployStorage{client: client, prefix: prefix}
}

// NewConsulDeployStorage returns a new instance of ConsulDeployStorage
func NewConsulDeployStorage(conf *api.Config, prefix string) (*ConsulDeployStorage, error) {
	client, err := newConsulClient(conf)
	if err != nil {
		return nil, err
	}

	return NewConsulDeployStorageFromClient(client, prefix), nil
}

// Create satisfies the DeployStorage interface
func (s *ConsulDeployStorage) Create(project, profile string, depl *pb.Deployment) error {
	key := s.keyPath(project, profile, depl.Name)
	val, err := depl.Marshal()
	if err != nil {
		return err
	}

	kv := s.client.KV()
	q := &api.QueryOptions{}
	kvp, _, err := kv.Get(key, q)
	if err != nil {
		return err
	}
	if kvp != nil {
		return fmt.Errorf("deployment exists: %s", depl.Name)
	}

	return putKV(kv, key, val)
}

// Get satisfies the DeployStorage interface
func (s *ConsulDeployStorage) Get(project, profile, id string) (*pb.Deployment, error) {
	kv := s.client.KV()
	q := &api.QueryOptions{}
	kvp, _, err := kv.Get(s.keyPath(project, profile, id), q)
	if err != nil {
		return nil, err
	}

	if kvp == nil {
		return nil, fmt.Errorf("deployment not found: %s", id)
	}

	var depl pb.Deployment
	err = depl.Unmarshal(kvp.Value)
	return &depl, err

}

// Update satisfies the DeployStorage interface
func (s *ConsulDeployStorage) Update(project, profile string, depl *pb.Deployment) error {
	key := s.keyPath(project, profile, depl.Name)
	val, err := depl.Marshal()
	if err != nil {
		return err
	}

	kv := s.client.KV()
	q := &api.QueryOptions{}
	kvp, _, err := kv.Get(key, q)
	if err != nil {
		return err
	}

	if kvp == nil {
		return fmt.Errorf("deployment not found: %s", depl.Name)
	}

	return putKV(kv, key, val)
}

// Delete satisfies the DeployStorage interface
func (s *ConsulDeployStorage) Delete(project, profile, id string) error {
	kv := s.client.KV()
	opt := &api.WriteOptions{}

	key := s.keyPath(project, profile, id)
	fmt.Println("DELETE", key)

	_, err := kv.Delete(s.keyPath(project, profile, id), opt)
	return err
}

// List satisfies the DeployStorage interface
func (s *ConsulDeployStorage) List(project, start string) ([]*pb.Deployment, error) {
	kv := s.client.KV()
	q := &api.QueryOptions{}

	prefix := filepath.Join(s.prefix, project, "instance")
	if start != "" {
		prefix = filepath.Join(prefix, start)
	}

	kvps, _, err := kv.List(prefix, q)
	if err != nil {
		return nil, err
	}

	out := make([]*pb.Deployment, 0, len(kvps))
	for _, kvp := range kvps {
		var depl pb.Deployment
		err = depl.Unmarshal(kvp.Value)
		if err == nil {
			out = append(out, &depl)
		}
	}
	return out, nil
}

func (s *ConsulDeployStorage) keyPath(project, profile, id string) string {
	return filepath.Join(s.prefix, project, "instance", profile, id)
}
