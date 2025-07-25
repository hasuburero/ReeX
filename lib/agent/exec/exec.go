package exec

import (
	"fmt"
	"os"
	"sync"
)

func GetStatus(sessionid, tid string) (Status, error) {
	session, exists := IsSession(sessionid)
	if !exists {
		return Status{}, SessionIdEmptyError
	}

	transaction, exists := session.IsTransaction(tid)
	if !exists {
		return Status{}, TransactionIdEmptyError
	}

	return Status{Status: transaction.GetStatus(), Pid: transaction.Pid, Tid: transaction.Tid}, nil
}

func WaitFinish(sessionid, tid string, timeout int) (Status, error) {
	session, exists := IsSession(sessionid)
	if !exists {
		return Status{}, SessionIdEmptyError
	}

	transaction, exists := session.IsTransaction(tid)
	if !exists {
		return Status{}, TransactionIdEmptyError
	}

	status, err := transaction.WaitFinish(timeout)
	if err != nil {
		return Status{}, err
	}

	return Status{Status: status, Pid: transaction.Pid, Tid: transaction.Tid}, nil
}

func (self *Session) Exec(tid, cmd string) (string, error) {
	var new_transaction = new(Transaction)
	new_transaction.Tid = tid
	new_transaction.Cmd = cmd
	new_transaction.StatusFailed = make(chan bool)
	new_transaction.StatusProcessing = make(chan bool)
	new_transaction.StatusFinished = make(chan bool)

	return pid, nil
}

func Exec(sessionid, tid, cmd string) (string, error) {
	session, exists := IsSession(sessionid)
	if !exists {
		return "", SessionIdEmptyError
	}

	pid, err := session.Exec(tid, cmd)

	return pid, err
}
