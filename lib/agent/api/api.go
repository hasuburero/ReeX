package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
)

// internal package
import (
	"github.com/hasuburero/ReeX/lib/common"
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
	Get    = http.MethodGet
	Post   = http.MethodPost
	Put    = http.MethodPut
	Delete = http.MethodDelete

	ContentType     = common.ContentType
	TextPlain       = common.TextPalin
	ApplicationJson = common.ApplicationJson
	OctetStream     = common.OctetStream

	ApiPath  = "/api/v1"
	ExecPath = ApiPath + "/exec"
	KillPath = ApiPath + "/kill"

	StatusMethodError         = common.StatusMethodError
	StatusInternalServerError = common.StatusInternalServerError
	ServerError               = common.ServerError
)

var (
	IPError     = errors.New("Empty IP Error\n")
	PortError   = errors.New("Empty Port Error\n")
	MethodError = errors.New("Invalid Method Error\n")
)

// setting status code and message to http.ResponseWriter. If error has occured, setting internal server error and err
func MakeError(w http.ResponseWriter, code int, message string) error {
	var ctx Error = Error{
		Code:    code,
		Message: message,
	}
	json_buf, err := json.Marshal(ctx)
	if err != nil {
		w.Header().Add(ContentType, TextPlain)
		buf := bytes.NewBufferString(ServerError)
		io.Copy(w, buf)
		return err
	}

	w.WriteHeader(code)
	buf := bytes.NewBuffer(json_buf)
	io.Copy(w, buf)

	return nil
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

	agent.ServeMux.HandleFunc(ExecPath, Exec)
	agent.ServeMux.HandleFunc(KillPath, Kill)

	return agent, nil
}
