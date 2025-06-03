package types

type DependencyTree struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Dependencies []*DependencyTree `json:"dependencies,omitempty"`
}

type Dependency struct {
	Id        int     `json:"id"`
	Package   Package `json:"package"`
	DependsOn Package `json:"depends_on"`
}
