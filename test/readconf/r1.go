package main

import (
	"fmt"
	"github.com/hasuburero/ReeX/lib/config/global"
	"github.com/hasuburero/ReeX/lib/config/session"
)

func main() {
	func() {
		nodes, err := global.Read("test1.json")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(nodes)
		for _, ctx := range nodes {
			if _, exists := ctx.AuthType["pubkey"]; exists {
				fmt.Println("pubkey exists")
			}
			for key, value := range ctx.AuthType {
				fmt.Println(key, value)
			}
		}
	}()

	func() {
		nodes, err := session.Read("test2.json")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(nodes)
	}()

	return
}
