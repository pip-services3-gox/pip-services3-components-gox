package config

import (
	"context"
	"fmt"
	"io/ioutil"

	cconfig "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	"github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	"github.com/pip-services3-gox/pip-services3-commons-gox/errors"
)

// JsonConfigReader is a config reader that reads configuration from JSON file.
// The reader supports parameterization using Handlebar template engine.
//	Configuration parameters:
//		- path: path to configuration file
//		- parameters: this entire section is used as template parameters
//		- ...
//	see IConfigReader
//	see FileConfigReader
//	Example:
//		======== config.json ======
//		{ "key1": "{{KEY1_VALUE}}", "key2": "{{KEY2_VALUE}}" }
//		===========================
//
//		configReader := NewJsonConfigReader("config.json")
//		parameters := NewConfigParamsFromTuples("KEY1_VALUE", 123, "KEY2_VALUE", "ABC")
//		res, err := configReader.ReadConfig(context.Background(), "123", parameters)
type JsonConfigReader struct {
	*FileConfigReader
}

// NewEmptyJsonConfigReader creates a new instance of the config reader.
//	Returns: *JsonConfigReader
func NewEmptyJsonConfigReader() *JsonConfigReader {
	return &JsonConfigReader{
		FileConfigReader: NewEmptyFileConfigReader(),
	}
}

// NewJsonConfigReader creates a new instance of the config reader.
//	Parameters: path string a path to configuration file.
//	Returns: *JsonConfigReader
func NewJsonConfigReader(path string) *JsonConfigReader {
	return &JsonConfigReader{
		FileConfigReader: NewFileConfigReader(path),
	}
}

// ReadObject reads configuration file, parameterizes its content and converts it into JSON object.
//	Parameters:
//		- ctx context.Context
//		- correlationId string transaction id to trace execution through call chain.
//		- parameters *cconfig.ConfigParams values to parameters the configuration.
//	Returns: any, error a JSON object with configuration adn error.
func (c *JsonConfigReader) ReadObject(ctx context.Context, correlationId string,
	parameters *cconfig.ConfigParams) (any, error) {

	if c.Path() == "" {
		return nil, errors.NewConfigError(correlationId, "NO_PATH", "Missing config file path")
	}

	b, err := ioutil.ReadFile(c.Path())
	if err != nil {
		err = errors.NewFileError(
			correlationId,
			"READ_FAILED",
			"Failed reading configuration "+c.Path()+": "+err.Error(),
		).
			WithDetails("path", c.Path()).WithCause(err)
		return nil, err
	}

	data := string(b)
	data, err = c.Parameterize(data, parameters)
	if err != nil {
		return nil, err
	}

	return convert.JsonConverter.ToMap(data), nil
}

// ReadConfig кeads configuration from a file, parameterize
// it with given values and returns a new ConfigParams object.
//	Parameters:
//		- ctx context.Context
//		- correlationId string transaction id to trace execution through call chain.
//		- parameters *cconfig.ConfigParams values to parameters the configuration.
//	Returns: *cconfig.ConfigParams, error
func (c *JsonConfigReader) ReadConfig(ctx context.Context, correlationId string,
	parameters *cconfig.ConfigParams) (result *cconfig.ConfigParams, err error) {

	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}
		}
	}()

	value, err := c.ReadObject(ctx, correlationId, parameters)
	if err != nil {
		return nil, err
	}

	config := cconfig.NewConfigParamsFromValue(value)
	return config, err
}

// ReadJsonObject reads configuration file, parameterizes its content and converts it into JSON object.
//	Parameters:
//		- ctx context.Context
//		- correlationId string transaction id to trace execution through call chain.
//		- path string
//		- parameters *cconfig.ConfigParams values to parameters the configuration.
//	Returns: any, error a JSON object with configuration.
func ReadJsonObject(ctx context.Context, correlationId string, path string,
	parameters *cconfig.ConfigParams) (any, error) {

	reader := NewJsonConfigReader(path)
	return reader.ReadObject(ctx, correlationId, parameters)
}

// ReadJsonConfig reads configuration from a file, parameterize it
// with given values and returns a new ConfigParams object.
//	Parameters:
//		- ctx context.Context
//		- correlationId string transaction id to trace execution through call chain.
//		- path string
//		- parameters *cconfig.ConfigParams values to parameters the configuration.
//	Returns: *cconfig.ConfigParams, error
func ReadJsonConfig(ctx context.Context, correlationId string, path string,
	parameters *cconfig.ConfigParams) (*cconfig.ConfigParams, error) {

	reader := NewJsonConfigReader(path)
	return reader.ReadConfig(ctx, correlationId, parameters)
}
