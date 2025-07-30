package router

import (
	"os"
	"sync"
	"sync/atomic"
	"gopkg.in/yaml.v3"
)

type PolicyMap map[string]string

type Router struct {
	policies atomic.Value // stores PolicyMap
	mu       sync.Mutex
	file     string
}

func NewRouter(policyFile string) (*Router, error) {
	r := &Router{file: policyFile}
	if err := r.reload(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Router) Lookup(team string) (backend string, ok bool) {
	pm, _ := r.policies.Load().(PolicyMap)
	backend, ok = pm[team]
	return
}

func (r *Router) reload() error {
	   data, err := os.ReadFile(r.file)
	   if err != nil {
			   return err
	   }
	   var wrapper struct {
			   Policies PolicyMap `yaml:"policies"`
	   }
	   if err := yaml.Unmarshal(data, &wrapper); err != nil {
			   return err
	   }
	   r.policies.Store(wrapper.Policies)
	   return nil
}

func (r *Router) HotReload() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reload()
}
