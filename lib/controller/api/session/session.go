package session

import (
	"errors"
	"net/http"
)

import (
	"github.com/hasuburero/ReeX/lib/controller/config/confsession"
)

// section1
type Session struct {
	Hosts  map[string]*Host
	Groups map[string]Group
	Tid    int
}

type Host struct {
	NodeName     string
	IP           string
	Transactions map[string]Transaction
}

type Transaction struct {
	Tid     string
	Pid     string
	Chan    chan bool
	Command string
}

type Group struct {
	Hosts []*Host
}

// section2
const (
	TransactionIdSize = 32 //32 characters
)

func (self *Session) NewTid() (string, error) {

}

func NewSession(filename string) (*Session, error) {
	config, err := confsession.Read(filename)
	if err != nil {
		return nil, err
	}

	var new_session = new(Session)
	new_session.Groups = make(map[string]Group)
	for _, ctx := range config.Node {
		var new_host = new(Host)
		new_host.NodeName = ctx.NodeName
		new_host.IP = ctx.IP
		new_host.Transactions = make(map[string]Transaction)
		if _, exists := new_session.Hosts[new_host.NodeName]; exists {
			return nil, errors.New("NodeName is already exists\n")
		}
		new_session.Hosts[new_host.NodeName] = new_host
	}

	for _, ctx := range config.Group {
		if _, exists := new_session.Groups[ctx.Name]; exists {
			return nil, errors.New("GroupName is already exists\n")
		}

		var node_buf = make(map[string]*Host)
		for _, node := range ctx.Nodename {
			if _, exists := node_buf[node]; exists {
				return nil, errors.New("NodeName is already exists, inside the Group field\n")
			}
			node_buf[node] = new_session.Hosts[node]
		}
		new_session.Groups[ctx.Name]
	}
	return new_session, nil
}
