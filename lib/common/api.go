package common

import ()

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Post_Kill_Struct struct {
	SessionID string `json:"sessionid"`
	Tid       string `json:"tid"`
	Pid       string `json:"pid"`
}

type Get_Exec_Struct struct {
	Pid    string `json:"pid"`
	Tid    string `json:"tid"`
	Cmd    string `json:"cmd"`
	Status string `json:"status"`
}

type Post_Exec_Struct struct {
	SessionID string `json:"sessionid"`
	Pid       string `json:"pid"`
	Tid       string `json:"tid"`
	Cmd       string `json:"cmd"`
}
