package filter

import (
	"fmt"

	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/pkg/errors"
)

// factory function for a filter
type factory func(interface{}) (Filter, error)

// Filter is a filter instance
type Filter interface {
	Apply(*runtime.Node) *runtime.Node
}

// Set is a list of configured filters
type Set []Filter

var filters = make(map[string]factory)

// Register registers a new filter
func Register(name string, f factory) {
	if _, ok := filters[name]; ok {
		panic("already registered: " + name)
	}
	filters[name] = f
}

// New returns initializes a set of filters
func New(configs map[string]interface{}) (set Set, errs []error) {
	for name := range configs {
		if _, ok := filters[name]; !ok {
			errs = append(errs, fmt.Errorf("unknown filter: %s", name))
		}
	}
	for name, f := range filters {
		config := configs[name]
		filter, err := f(config)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "unable to initialize filter %s", name))
			continue
		}
		if filter != nil {
			set = append(set, filter)
		}

	}
	return
}

// Apply applies the filter set to the given node list and returns a new node list
func (set Set) Apply(nodesOrigin *runtime.Nodes) *runtime.Nodes {
	nodes := runtime.NewNodes(&runtime.NodesConfig{})

	nodesOrigin.Lock()
	defer nodesOrigin.Unlock()

	for _, nodeOrigin := range nodesOrigin.List {
		//maybe cloning of this object is better?
		node := nodeOrigin
		for _, filter := range set {
			node = filter.Apply(node)
			if node == nil {
				break
			}
		}

		if node != nil {
			nodes.AddNode(node)
		}
	}
	return nodes
}
