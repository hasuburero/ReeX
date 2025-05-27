package session

import (
	"encoding/json"
	"io"
	"os"
)

type Node struct {
	NodeName string `json:"nodename"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Workdir  string `json:"workdir"`
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

	var nodes []Node
	err = json.Unmarshal(buf, &nodes)
	if err != nil {
		return []Node{}, err
	}

	return nodes, nil
}
