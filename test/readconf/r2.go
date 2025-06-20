package main

import (
	"fmt"
	"github.com/hasuburero/ReeX/lib/controller/config/confsession"
)

const (
	conf = "./session.json"
)

func main() {
	config, err := confsession.Read(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(config)
}
