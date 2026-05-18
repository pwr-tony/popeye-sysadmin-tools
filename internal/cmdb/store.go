package cmdb

import (
	"encoding/json"
	"os"
)

type Store struct {
	path string
	data CMDB
}

func NewStore(path string) *Store {
	return &Store{
		path: path,
		data: CMDB{
			Servers:  make([]Server, 0),
			Networks: make([]Network, 0),
		},
	}
}

func (s *Store) Load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &s.data)
}

func (s *Store) Save() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0644)
}

func (s *Store) GetServers() []Server {
	return s.data.Servers
}

func (s *Store) GetNetworks() []Network {
	return s.data.Networks
}

func (s *Store) AddServer(server Server) {
	s.data.Servers = append(s.data.Servers, server)
}

func (s *Store) AddNetwork(network Network) {
	s.data.Networks = append(s.data.Networks, network)
}

func (s *Store) DeleteServer(id string) bool {
	for i, srv := range s.data.Servers {
		if srv.ID == id {
			s.data.Servers = append(s.data.Servers[:i], s.data.Servers[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Store) DeleteNetwork(id string) bool {
	for i, net := range s.data.Networks {
		if net.ID == id {
			s.data.Networks = append(s.data.Networks[:i], s.data.Networks[i+1:]...)
			return true
		}
	}
	return false
}
