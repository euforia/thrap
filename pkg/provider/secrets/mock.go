package secrets

import "github.com/euforia/thrap/pkg/provider"

// Mock is secrets provider used in tests
type Mock struct {
}

// Init implements Secrects interface
func (s *Mock) Init(conf *provider.Config) error {
	return nil
}

// SecretsPath implements Secrects interface
func (s *Mock) SecretsPath(key string) string {
	return ""
}

// Create implements Secrects interface
func (s *Mock) Create(path string) error {
	return nil
}

// List implements Secrects interface
func (s *Mock) List(prefix string) ([]string, error) {
	return nil, nil
}

// Get implements Secrects interface
func (s *Mock) Get(path string) (map[string]interface{}, error) {
	return nil, nil
}

// RecursiveGet implements Secrects interface
func (s *Mock) RecursiveGet(path string) (map[string]map[string]interface{}, error) {
	return nil, nil
}

// Set implements Secrects interface
func (s *Mock) Set(path string, data map[string]interface{}) error {
	return nil
}
