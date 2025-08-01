package exec

import (
	"os/exec"
	"strings"
)

// external package
import (
	cstrings "github.com/hasuburero/util/strings"
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

func (self *Session) Kill(tid string) error {
	self.Mux.Lock()
	transaction, exists := self.Transactions[tid]
	if !exists {
		return TransactionNotExistsError
	}
	self.Mux.Unlock()

	err := transaction.Execmd.Process.Kill()
	if err != nil {
		return KillError
	}
	return nil
}

func (self *Transaction) Wait() {
	go func() {
		err := self.Execmd.Wait()
		if err != nil {
			close(self.StatusFailed)
		} else {
			close(self.StatusFinished)
		}
	}()

	return
}

func (self *Session) Exec(tid, cmd string) (*Transaction, error) {
	var new_transaction = new(Transaction)
	new_transaction.Tid = tid
	new_transaction.Cmd = cmd
	new_transaction.StatusFailed = make(chan bool)
	new_transaction.StatusProcessing = make(chan bool)
	new_transaction.StatusFinished = make(chan bool)

	slice := strings.Split(cmd, " ")
	slice = cstrings.TrimSlice(slice, "")

	new_transaction.Execmd = exec.Command(slice[0], slice[1:]...)
	err := new_transaction.Execmd.Start()
	if err != nil {
		close(new_transaction.StatusFailed)
		return nil, ExecError
	}

	err = self.AddTransaction(new_transaction)
	if err != nil {
		close(new_transaction.StatusFailed)
		new_transaction.Execmd.Process.Kill()
		return nil, err
	}

	close(new_transaction.StatusProcessing)

	return new_transaction, err
}

func Exec(sessionid, tid, cmd string) (string, error) {
	session, exists := IsSession(sessionid)
	if !exists {
		return "", SessionIdEmptyError
	}

	transaction, err := session.Exec(tid, cmd)

	return transaction.Pid, err
}
