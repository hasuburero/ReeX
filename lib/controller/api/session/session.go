package session

import (
	"errors"
	"strconv"
)

import (
	"github.com/hasuburero/ReeX/lib/controller/config/confsession"
)

// section1
type Session struct {
	Hosts  map[string]*Host
	Groups map[string]Group
	Tid    int64
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
	TransactionIdSize = 16 //16 characters
	Base              = 10
)

func (self *Session) NewTid() string {
	tid := strconv.FormatInt(self.Tid, Base)
	zero := ""
	for i := 0; i < TransactionIdSize-len(tid); i++ {
		zero += "0"
	}

	return zero + tid
}

func NewSession(filename string) (*Session, error) {
	config, err := confsession.Read(filename)
	if err != nil {
		return nil, err
	}

	var new_session = new(Session)
	new_session.Hosts = make(map[string]*Host)
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
		var node_list = make([]*Host, 0)
		for _, node := range ctx.Nodename {
			if _, exists := node_buf[node]; exists {
				return nil, errors.New("NodeName is already exists, inside the Group field\n")
			}
			node_buf[node] = new_session.Hosts[node]
			node_list = append(node_list, new_session.Hosts[node])
		}
		new_session.Groups[ctx.Name] = Group{node_list}
	}
	return new_session, nil
}
