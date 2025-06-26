package session

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

type Exec_Struct struct {
	Pid string `json:"pid"`
	Tid string `json:"tid"`
	Cmd string `json:"cmd"`
}

// return
func (self *Session) ExecGroup(group, cmd string) ([]string, error) {
	var wg sync.WaitGroup
	var mux sync.Mutex
	var pids = make([]string, 0)
	go func() {
		wg.Add(1)

		// adding Host.ExecAsync

		wg.Done()
	}()
	wg.Wait()

	return pids, nil
}

// return
func (self *Session) Exec(nodename, cmd string) (string, error) {
	host, exists := self.Hosts[nodename]
	if !exists {
		return "", NotExistsError
	}

	tid := self.NewTid()
	var ctx Exec_Struct
	ctx.Cmd = cmd
	ctx.Tid = tid

	json, err := json.Marshal(ctx)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(Post)

	return ctx.Pid, nil
}
