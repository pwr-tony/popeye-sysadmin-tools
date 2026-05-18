package cmdb

type Server struct {
	ID          string   `json:"id"`
	Hostname    string   `json:"hostname"`
	IP          string   `json:"ip"`
	OS          string   `json:"os"`
	Environment string   `json:"environment"`
	Tags        []string `json:"tags"`
	Notes       string   `json:"notes"`
}

type Network struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	CIDR    string `json:"cidr"`
	VLAN    int    `json:"vlan"`
	Gateway string `json:"gateway"`
}

type CMDB struct {
	Servers  []Server  `json:"servers"`
	Networks []Network `json:"networks"`
}
