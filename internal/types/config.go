package types

type Config struct {
	Database DatabaseConfig `json:"database"`
	Pypi     PypiConfig     `json:"pypi"`
}
