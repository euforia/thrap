package storage

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/euforia/thrap/thrapb"
	"github.com/hashicorp/consul/api"
)

type ConsulProjectStorage struct {
	client *api.Client
	prefix string
}

func NewConsulProjectStorage(conf *api.Config, prefix string) (*ConsulProjectStorage, error) {
	conf.HttpClient = &http.Client{Timeout: 2 * time.Second}
	client, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}

	return &ConsulProjectStorage{
		client: client,
		prefix: prefix,
	}, nil
}

func (s *ConsulProjectStorage) Iter(start string, cb func(*thrapb.Project) error) error {
	kv := s.client.KV()
	q := &api.QueryOptions{}
	prefix := s.keyPath(start)

	kvps, _, err := kv.List(prefix, q)
	if err != nil {
		return err
	}
	for _, kvp := range kvps {
		proj := &thrapb.Project{}
		err = proj.Unmarshal(kvp.Value)
		if err != nil {
			log.Println("ERR", err)
			continue
		}
		if err = cb(proj); err != nil {
			break
		}
	}
	return err
}

func (s *ConsulProjectStorage) Create(proj *thrapb.Project) error {
	key := s.keyPath(proj.ID)
	val, err := proj.Marshal()
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
		return fmt.Errorf("project exists: %s", proj.ID)
	}

	return putKV(kv, key, val)
}

func (s *ConsulProjectStorage) Update(proj *thrapb.Project) error {
	key := s.keyPath(proj.ID)
	val, err := proj.Marshal()
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
		return fmt.Errorf("project not found: %s", proj.ID)
	}

	return putKV(kv, key, val)
}

func (s *ConsulProjectStorage) Get(id string) (*thrapb.Project, error) {
	kv := s.client.KV()
	q := &api.QueryOptions{}
	kvp, _, err := kv.Get(s.keyPath(id), q)
	if err != nil {
		return nil, err
	}
	if kvp == nil {
		return nil, fmt.Errorf("project not found: %s", id)
	}
	var proj thrapb.Project
	err = proj.Unmarshal(kvp.Value)
	return &proj, err
}

func (s *ConsulProjectStorage) Delete(id string) error {
	kv := s.client.KV()
	opt := &api.WriteOptions{}

	key := s.keyPath(id)
	fmt.Println("DELETE", key)

	_, err := kv.Delete(s.keyPath(id), opt)
	return err
}

func (s *ConsulProjectStorage) keyPath(k string) string {
	return filepath.Join(s.prefix, k)
}

func putKV(kv *api.KV, key string, val []byte) error {
	p := &api.KVPair{
		Key:   key,
		Value: val,
	}

	_, err := kv.Put(p, &api.WriteOptions{})
	return err
}