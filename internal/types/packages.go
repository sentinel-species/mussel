package types

type Package struct {
	InternalId int       `json:"internal_id"`
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	Ecosystem  Ecosystem `json:"ecosystem"`
}
