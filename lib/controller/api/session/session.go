package session

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
)

import (
	"github.com/hasuburero/ReeX/lib/controller/config/confsession"
)

// section1
type Session struct {
	Hosts        map[string]*Host        // key:nodename, value: pointer(host)
	Groups       map[string][]*Host      // key:groupname, value:slice(pointer(host))
	Transactions map[string]*Transaction // key:tid, value:pointer(transaction)
	Tid          int64
	Mux          sync.Mutex
}

type Host struct {
	NodeName     string
	HostName     string
	Transactions map[string]*Transaction
	Mux          sync.Mutex
	Client       *http.Client
}

type Transaction struct {
	NodeName string
	Tid      string
	Pid      string
	Chan     chan bool
	Command  string
}

// section2
const (
	TransactionIdSize = 16 //16 characters
	Base              = 10

	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusFinished   = "finished"
)

const (
	BasePath = "/api/v1"
	ExecPath = BasePath + "/exec"
	KillPath = BasePath + "/kill"

	Post   = "POST"
	Get    = "GET"
	Put    = "PUT"
	Delete = "DELETE"
)

var (
	NotExistsError = errors.New("Does not exists\n")
	UnknownError   = errors.New("Unknown Error has occured\n")
	TimeoutError   = errors.New("timeout\n")
	WaitingError   = errors.New("Waiting error occured\n")
)

var (
	TidMux sync.Mutex
)

// section3
func (self *Session) NewTid() string {
	TidMux.Lock()
	tid := strconv.FormatInt(self.Tid, Base)
	self.Tid++
	TidMux.Unlock()
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
	new_session.Tid = 1
	new_session.Hosts = make(map[string]*Host)
	new_session.Groups = make(map[string][]*Host)
	new_session.Transactions = make(map[string]*Transaction)
	for _, ctx := range config.Node {
		var new_host = new(Host)
		new_host.NodeName = ctx.NodeName
		new_host.HostName = ctx.IP + ":" + ctx.Port
		new_host.Transactions = make(map[string]*Transaction)
		new_host.Client = &http.Client{
			Transport: &http.Transport{
				DisableKeepAlives: false,
			},
		}
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
		for _, node := range ctx.NodeName {
			if _, exists := node_buf[node]; exists {
				return nil, errors.New("NodeName is already exists, inside the Group field\n")
			}
			node_buf[node] = new_session.Hosts[node]
			node_list = append(node_list, new_session.Hosts[node])
		}
		new_session.Groups[ctx.Name] = node_list
	}
	return new_session, nil
}
