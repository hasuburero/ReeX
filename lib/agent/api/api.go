package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

type Transaction struct {
	Tid     string
	Pid     string
	Chan    chan bool
	Command string
}

type Agent struct {
	HostName     string
	NodeName     string
	Port         string
	Transactions map[string]*Transaction
	Mux          sync.Mutex
	Server       http.Server
	ServeMux     *http.ServeMux
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	Get    = "GET"
	Post   = "POST"
	Put    = "PUT"
	Delete = "DELETE"

	ApiPath  = "/api/v1"
	ExecPath = ApiPath + "/exec"
	KillPath = ApiPath + "/kill"

	StatusMethodError = "Method Not Allowed\n"
)

var (
	IPError     = errors.New("Empty IP Error\n")
	PortError   = errors.New("Empty Port Error\n")
	MethodError = errors.New("Invalid Method Error\n")
)

func MakeError(code int, message string) (Error, error) {
	var ctx Error = Error{
		Code:    code,
		Message: message,
	}
	json_buf, err := json.Marshal(ctx)
	if err != nil {
		return
	}
	return
}

func Start(ip string, port string) (Agent, error) {
	if ip == "" {
		return Agent{}, IPError
	} else if port == "" {
		return Agent{}, PortError
	}

	var agent Agent
	agent.ServeMux = http.NewServeMux()
	agent.Server = http.Server{
		Addr:    ip + ":" + port,
		Handler: agent.ServeMux,
	}

	agent.ServeMux.handleFunc(ExecPath, Exec)
	agent.ServeMux.handleFunc(KillPath, Kill)

	return agent, nil
}
