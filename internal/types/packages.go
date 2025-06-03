package types

type Package struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Ecosystem string `json:"ecosystem"`
}
