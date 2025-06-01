package types

type DependencyTree struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Dependencies []*DependencyTree `json:"dependencies,omitempty"`
}
