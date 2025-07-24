package exec

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Status struct {
	Status string
	Pid    string
	Tid    string
}

type Session struct {
	SessionID    string
	Transactions map[string]*Transaction
	Mux          sync.Mutex
}

type Transaction struct {
	Tid              string
	Pid              string
	Cmd              string
	StatusProcessing chan bool
	StatusFinished   chan string
	Mux              sync.Mutex
}

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusFinished   = "finished"
)

var (
	SessionIdConflictError  = errors.New("Session IDs are conflicted\n")
	SessionIdEmptyError     = errors.New("Session ID is Empty\n")
	TransactionIdEmptyError = errors.New("Transaction ID is Empty\n")
	InvalidStatusError      = errors.New("Invalid Status Error\n")
	TimeoutError            = errors.New("Timeout\n")
)

var (
	Sessions map[string]*Session
	Mux      sync.Mutex
)

func (self *Transaction) GetStatus() string {
	select {
	case _ = <-self.StatusProcessing:
	default:
		return StatusPending
	}
	select {
	case _ = <-self.StatusFinished:
		return StatusFinished
	default:
		return StatusProcessing
	}
}

func (self *Transaction) WaitFinish(timeout int) (string, error) {
	select {
	case _ = <-self.StatusFinished:
		return StatusFinished, nil
	case <-time.After(time.Duration(timeout) * time.Second):
		return "", TimeoutError
	}
}

func IsSession(sessionid string) (*Session, bool) {
	Mux.Lock()
	session, exists := Sessions[sessionid]
	Mux.Unlock()
	return session, exists
}

func (self *Session) IsTransaction(tid string) (*Transaction, bool) {
	self.Mux.Lock()
	transaction, exists := self.Transactions[tid]
	self.Mux.Unlock()
  :q

	return transaction, exists
}

func NewSession(uuid string) error {
	new_session := new(Session)
	_, exists := Sessions[uuid]
	if exists {
		return SessionIdConflictError
	}
	Sessions[uuid] = new_session

	return nil
}

func init() {
	Sessions = make(map[string]*Session)
}
