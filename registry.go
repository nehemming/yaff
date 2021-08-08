/*
Copyright (c) 2020-2021 The yaff Authors (Neil Hemming)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package yaff

import (
	"sync"

	"github.com/nehemming/lpax"
	"github.com/nehemming/yaff/langpack"
)

// Registry holds a list of available formatters.
// It is possible to create multiple formatters, but by default the
// default shared yaff.Formatters registry can be used
// Registry implementors must be threadsafe across all calls.
type Registry interface {

	// Register adds or updates an available format.
	// The NewFormatter factory is used to create an
	// instance oif the registered formatter.
	Register(format Format, factory NewFormatter)

	// GetFormatter returns the formatter supporting the format or an error if no formatter can be found.
	GetFormatter(format Format) (Formatter, error)

	// Formats returns a slice of supported formats.
	Formats() []Format
}

type registry struct {
	factories map[Format]NewFormatter
	mu        sync.Mutex
}

// NewRegistry creates a new registry.
func NewRegistry() Registry {
	return &registry{
		factories: make(map[Format]NewFormatter),
	}
}

func (r *registry) Formats() []Format {
	// Lock to maintain thread safety.
	r.mu.Lock()
	defer r.mu.Unlock()

	c := make([]Format, len(r.factories))
	i := 0
	for k := range r.factories {
		c[i] = k
		i++
	}

	return c
}

func (r *registry) Register(format Format, factory NewFormatter) {
	// Lock to maintain thread safety
	r.mu.Lock()
	defer r.mu.Unlock()

	// Passing nil will "deregister" the factory
	r.factories[format] = factory
}

func (r *registry) GetFormatter(format Format) (Formatter, error) {
	// getFactory locks registry, but keeps call to factory outside lock
	factory := r.getFactory(format)

	if factory == nil {
		return nil, lpax.Errorf(langpack.ErrorUnknownFormatter, format)
	}

	return factory()
}

func (r *registry) getFactory(format Format) NewFormatter {
	// Lock to maintain thread safety
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.factories[format]
}

var sharedRegistry = NewRegistry()

// Formatters is the shared registry of formatters.
func Formatters() Registry {
	return sharedRegistry
}
