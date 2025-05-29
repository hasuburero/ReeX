package confssh

import (
	"encoding/json"
	"io"
	"os"
)

type Node struct {
	Nodename string            `json:"nodename"`
	IP       string            `json:"ip"`
	User     string            `json:"user"`
	AuthType map[string]string `json:"authtype"`
}

func Read(filename string) ([]Node, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return []Node{}, err
	}
	defer fd.Close()

	buf, err := io.ReadAll(fd)
	if err != nil {
		return []Node{}, err
	}

	var Nodes []Node
	err = json.Unmarshal(buf, &Nodes)
	if err != nil {
		return []Node{}, err
	}

	return Nodes, nil
}
