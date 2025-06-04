package types

type Dependency struct {
	InternalId int     `json:"internal_id"`
	Package    Package `json:"package"`
	DependsOn  Package `json:"depends_on"`
}
