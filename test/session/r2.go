package main

import (
	"fmt"
	"github.com/hasuburero/ReeX/lib/controller/api/session"
)

func main() {
	new_session := session.Session{}
	new_session.Tid = 1
	fmt.Println(new_session.NewTid())
	fmt.Println(new_session.NewTid())

	return
}
