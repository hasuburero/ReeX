package main

import (
	"fmt"
	"github.com/hasuburero/ReeX/lib/controller/api/session"
)

const (
	filename = "session.json"
)

func main() {
	Session, err := session.NewSession(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, ctx := range Session.Hosts {
		fmt.Print(ctx, " ")
		fmt.Printf("%p\n", ctx)
	}

	fmt.Println("")
	for key, ctx := range Session.Groups {
		fmt.Println(key, ctx)
	}

	return
}
