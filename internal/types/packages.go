package types

type Package struct {
	Name     string    `json:"name"`
	Versions []Version `json:"versions"`
}

type Version struct {
	Version      string          `json:"version"`
	Dependencies *DependencyTree `json:"dependencies,omitempty"`
}
