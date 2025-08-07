package session

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"sync"
)

// internal package
import (
	"github.com/hasuburero/ReeX/lib/common"
)

const (
	StatusOK      = http.StatusOK
	StatusError   = 500
	StatusTimeout = 501
)

func (self *Session) WaitTids(tids []string, timeout int) ([]string, []string, error) {
	var wg sync.WaitGroup

	var finished = make([]string, 0)
	var errored = make([]string, 0)
	for _, tid := range tids {
		wg.Add(1)
		go func(tid string) {
			return_tid, err := self.Wait(tid, timeout)
			if err != nil {
				err = errors.New(return_tid + ": " + err.Error())
				errored = append(errored, return_tid)
			} else {
				finished = append(finished, return_tid)
			}
			wg.Done()
		}(tid)
	}

	wg.Wait()

	if len(errored) != 0 {
		return finished, errored, WaitingError
	}
	return finished, errored, nil
}

func (self *Session) Wait(tid string, timeout int) (string, error) {
	self.Mux.Lock()
	transaction, exists := self.Transactions[tid]
	if !exists {
		return "", NotExistsError
	}

	host, exists := self.Hosts[transaction.NodeName]
	if !exists {
		return "", NotExistsError
	}

	params := "?sessionid=" + self.SessionID
	params += "&tid=" + tid
	if timeout > 0 {
		to := strconv.Itoa(timeout)
		params += "&timeout=" + to
	}

	request, err := http.NewRequest(Post, host.HostName+ExecPath+params, nil)
	if err != nil {
		return "", err
	}

	res, err := host.Client.Do(request)
	if err != nil {
		return "", err
	}

	resbody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case StatusOK:
		var ctx common.Get_Exec_Struct
		err = json.Unmarshal(resbody, &ctx)
		if err != nil {
			return "", err
		}

		return ctx.Tid, err
	case StatusError:
		var ctx common.Error
		err = json.Unmarshal(resbody, &ctx)
		if err != nil {
			return "", err
		}
		return "", errors.New(ctx.Message)
	case StatusTimeout:
		return "", TimeoutError
	default:
		return "", UnknownError
	}
}

// return []Transaction ID, error. if len(Transaction ID) != 0, error includes failed nodename.
func (self *Session) ExecGroup(group, cmd string) ([]string, error) {
	var wg sync.WaitGroup

	self.Mux.Lock()
	hosts, exists := self.Groups[group]
	self.Mux.Unlock()
	if !exists {
		return []string{}, NotExistsError
	}

	var tids = make([]string, 0)
	var errored = make([]string, 0)
	for _, host := range hosts {
		wg.Add(1)
		go func(host *Host) {
			// adding Host.ExecAsync
			tid, err := self.Exec(host.NodeName, cmd)
			if err != nil {
				errored = append(errored, host.NodeName)
			}
			tids = append(tids, tid)
			wg.Done()
		}(host)
	}
	wg.Wait()

	msg := ""
	for i := 0; i < len(errored); i++ {
		msg += errored[i]
		if i != len(errored)-1 {
			msg += " "
		}
	}

	return tids, errors.New(msg)
}

// return Transaction ID and error
func (self *Session) Exec(nodename, cmd string) (string, error) {
	self.Mux.Lock()
	host, exists := self.Hosts[nodename]
	self.Mux.Unlock()
	if !exists {
		return "", NotExistsError
	}

	tid := self.NewTid()
	var ctx common.Post_Exec_Struct
	ctx.SessionID = self.SessionID
	ctx.Cmd = cmd
	ctx.Tid = tid

	var transaction = new(Transaction)
	transaction.Chan = make(chan bool)
	transaction.Tid = tid
	transaction.Command = cmd
	transaction.NodeName = host.NodeName

	json_buf, err := json.Marshal(ctx)
	if err != nil {
		return "", err
	}

	reqbody := bytes.NewBuffer(json_buf)
	request, err := http.NewRequest(Post, host.HostName+ExecPath, reqbody)
	if err != nil {
		return "", err
	}

	res, err := host.Client.Do(request)
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
		var ctx common.Post_Exec_Struct
		err = json.Unmarshal(res_body, &ctx)
		if err != nil {
			return "", err
		}

		transaction.Pid = ctx.Pid

		host.Mux.Lock()
		host.Transactions[transaction.Tid] = transaction
		host.Mux.Unlock()
		self.Mux.Lock()
		self.Transactions[transaction.Tid] = transaction
		self.Mux.Unlock()
		return transaction.Tid, nil
	case StatusError:
		var ctx common.Error
		err = json.Unmarshal(res_body, &ctx)
		if err != nil {
			return "", err
		}
		return "", errors.New(ctx.Message)
	default:
		return "", UnknownError
	}
}
