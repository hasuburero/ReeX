package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// internal package
import (
	"github.com/hasuburero/ReeX/lib/agent/exec"
)

// external package
import ()

type Get_Exec_Struct struct {
	Pid    string `json:"pid"`
	Tid    string `json:"tid"`
	Cmd    string `json:"cmd"`
	Status string `json:"status"`
}

type Post_Exec_Struct struct {
	Pid string `json:"pid"`
	Tid string `json:"tid"`
	Cmd string `json:"cmd"`
}

const (
	KeySessionID = "sessionid"
	KeyTid       = "tid"
	KeyTimeout   = "timeout"
)

const (
	defaultTimeout = -1
)

func Get_Exec(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	params := url.Query()

	var err error
	var timeout int = defaultTimeout
	var tid string
	var sessionid string
	for key, value := range params {
		switch key {
		case KeySessionID:
			sessionid = value[0]
		case KeyTid:
			tid = value[0]
		case KeyTimeout:
			length, err := strconv.Atoi(value[0])
			if err != nil {
				continue
			}
			timeout = length
		default:
			continue
		}
	}

	var status exec.Status
	switch timeout {
	case defaultTimeout:
		status, err = exec.GetStatus(sessionid, tid)
	default:
		status, err = exec.WaitFinish(sessionid, tid, timeout)
	}

	if err != nil {
		fmt.Println(err)
	}

	var ctx Get_Exec_Struct = Get_Exec_Struct{
		Pid:    status.Pid,
		Tid:    status.Tid,
		Status: status.Status,
	}

	json_buf, err := json.Marshal(ctx)
	if err != nil {
		fmt.Println(err)
		err = MakeError(w, http.StatusInternalServerError, StatusInternalServerError)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	res_buf := bytes.NewBuffer(json_buf)
	io.Copy(w, res_buf)

	return
}

func Post_Exec(w http.ResponseWriter, r *http.Request) {
	var ctx Post_Exec_Struct
	req_body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		err = MakeError(w, http.StatusInternalServerError, StatusInternalServerError)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(req_body, &ctx)
	if err != nil {
		fmt.Println(err)
		err = MakeError(w, http.StatusInternalServerError, StatusInternalServerError)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	process := exec.Exec(ctx.SessionID, ctx.Cmd, ctx.Tid)

	ctx.Pid = process.Pid
	res_json, err := json.Marshal(ctx)
	if err != nil {
		fmt.Println(err)
		err = MakeError(w, http.StatusInternalServerError, StatusInternalServerError)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	res_body := bytes.NewBuffer(res_json)
	io.Copy(w, res_body)

	return
}

func Exec(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case Get:
		Get_Exec(w, r)
	case Post:
		Post_Exec(w, r)
	default:
		var ctx Error
		ctx.Code = http.StatusMethodNotAllowed
		ctx.Message = StatusMethodError
		json_buf, err := json.Marshal(ctx)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
		buf := bytes.NewBuffer(json_buf)
		_, err = io.Copy(w, buf)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	}
}
