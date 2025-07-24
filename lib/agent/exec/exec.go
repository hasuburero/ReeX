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
