package main

import (
	"context"
	"fmt"
	"log"

	"github.com/blue-script/coal_mine/enterprise"
	"github.com/blue-script/coal_mine/rest"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	ent := enterprise.NewEnterprise(ctx, cancel)
	ent.RunPassiveIncome()

	handlers := rest.NewHTTPHandlers(ent)
	server := rest.NewHTTPServer(handlers)
	err := server.Start()
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Finished game")
}
