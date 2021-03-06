package cache

import (
	"github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	"github.com/pip-services3-gox/pip-services3-components-gox/build"
)

// Creates ICache components by their descriptors.

var NullCacheDescriptor = refer.NewDescriptor("pip-services", "cache", "null", "*", "1.0")
var MemoryCacheDescriptor = refer.NewDescriptor("pip-services", "cache", "memory", "*", "1.0")

// NewDefaultCacheFactory create a new instance of the factory.
//	Returns: *build.Factory
func NewDefaultCacheFactory() *build.Factory {
	factory := build.NewFactory()

	factory.RegisterType(NullCacheDescriptor, NewNullCache[any])
	factory.RegisterType(MemoryCacheDescriptor, NewMemoryCache[any])

	return factory
}
