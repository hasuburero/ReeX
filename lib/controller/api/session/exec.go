package session

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
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

const (
	StatusOK    = http.StatusOK
	StatusError = 500
)

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

	json_buf, err := json.Marshal(ctx)
	if err != nil {
		return "", err
	}

	reqbody := bytes.NewBuffer(json_buf)
	request, err := http.NewRequest(Post, host.Hostname, reqbody)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: false,
		},
	}

	res, err := client.Do(request)
	if err != nil {
		return "", err
	}

	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case StatusOK:
		var ctx Post_Exec_Struct
		err = json.Unmarshal(res_body, &ctx)
		if err != nil {
			return "", err
		}
		return ctx.Pid, nil
	case StatusError:
		var ctx Error
		err = json.Unmarshal(res_body, &ctx)
		if err != nil {
			return "", err
		}
		return "", errors.New(ctx.Message)
	default:
		return "", UnknownError
	}
}
