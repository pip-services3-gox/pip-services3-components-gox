package auth

import (
	"context"

	"github.com/pip-services3-gox/pip-services3-commons-gox/config"
	"github.com/pip-services3-gox/pip-services3-commons-gox/errors"
)

// MemoryCredentialStore Credential store that keeps credentials in memory.
//		Configuration parameters:
//			- [credential key 1]:
//				- ... credential parameters for key 1
//			- [credential key 2]:
//				- ... credential parameters for key N
//				- ...
// see ICredentialStore
// see CredentialParams
//
//	Example:
//		config := NewConfigParamsFromTuples(
//			"key1.user", "jdoe",
//			"key1.pass", "pass123",
//			"key2.user", "bsmith",
//			"key2.pass", "mypass"
//		);
//		credentialStore := NewEmptyMemoryCredentialStore();
//		credentialStore.ReadCredentials(config);
//		res, err := credentialStore.Lookup(context.Backgroudn(), "123", "key1");
type MemoryCredentialStore struct {
	items map[string]*CredentialParams
}

// NewEmptyMemoryCredentialStore creates a new instance of the credential store.
//	Returns: *MemoryCredentialStore
func NewEmptyMemoryCredentialStore() *MemoryCredentialStore {
	return &MemoryCredentialStore{
		items: make(map[string]*CredentialParams),
	}
}

// NewMemoryCredentialStore creates a new instance of the credential store.
//	Parameters:
//		- ctx context.Context
//		- config *config.ConfigParams configuration with credential parameters.
//	Returns: *MemoryCredentialStore
func NewMemoryCredentialStore(ctx context.Context, config *config.ConfigParams) *MemoryCredentialStore {
	c := &MemoryCredentialStore{
		items: make(map[string]*CredentialParams),
	}

	if config != nil {
		c.Configure(ctx, config)
	}

	return c
}

// Configure configures component by passing configuration parameters.
//	Parameters:
//		- config *config.ConfigParams configuration parameters to be set.
func (c *MemoryCredentialStore) Configure(ctx context.Context, config *config.ConfigParams) {
	c.ReadCredentials(config)
}

// ReadCredentials reads credentials from configuration parameters.
// Each section represents an individual CredentialParams
//	Parameters:
//		- config *config.ConfigParams configuration parameters to be read
func (c *MemoryCredentialStore) ReadCredentials(config *config.ConfigParams) {
	c.items = make(map[string]*CredentialParams)

	sections := config.GetSectionNames()
	for _, section := range sections {
		value := config.GetSection(section)
		credential := NewCredentialParams(value.Value())
		c.items[section] = credential
	}
}

// Store credential parameters into the store.
//	Parameters:
//		- ctx context.Context.
//		- correlationId string transaction id to trace execution through call chain.
//		- key string a key to uniquely identify the credential parameters.
//		- credential *CredentialParams a credential parameters to be stored.
//	Returns: error
func (c *MemoryCredentialStore) Store(ctx context.Context, correlationId string, key string,
	credential *CredentialParams) error {

	if credential != nil {
		c.items[key] = credential
	} else {
		delete(c.items, key)
	}

	return nil
}

// Lookup credential parameters by its key.
//	Parameters:
//		- ctx context.Context.
//		- correlationId string transaction id to trace execution through call chain.
//		- key string a key to uniquely identify the credential parameters.
//	Returns: result *CredentialParams, err error result of lookup and error message
func (c *MemoryCredentialStore) Lookup(ctx context.Context, correlationId string,
	key string) (result *CredentialParams, err error) {

	if credential, ok := c.items[key]; ok && credential != nil {
		return credential, nil
	}

	return nil, errors.NewConfigError(
		correlationId, "MISSING_CREDENTIALS", "missing credential param: "+key)
}
