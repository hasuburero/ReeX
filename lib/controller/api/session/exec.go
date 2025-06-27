package session

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Post_Exec_Struct struct {
	Pid string `json:"pid"`
	Tid string `json:"tid"`
	Cmd string `json:"cmd"`
}

type Get_Exec_Struct struct {
	Pid    string `json:"pid"`
	Tid    string `json:"tid"`
	Cmd    string `json:"cmd"`
	Output string `json:"output"`
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
	var ctx Post_Exec_Struct
	ctx.Cmd = cmd
	ctx.Tid = tid

	json, err := json.Marshal(ctx)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(Post, host., , json)

	return ctx.Pid, nil
}
