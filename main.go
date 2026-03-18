package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/blue-script/coal_mine/enterprise"
	"github.com/blue-script/coal_mine/rest"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	ent := enterprise.NewEnterprise(ctx, cancel)
	ent.RunPassiveIncome()

	handlers := rest.NewHTTPHandlers(ent)
	server := rest.NewHTTPServer(handlers)

	var wg sync.WaitGroup
	wg.Add(1)

	handlers.SetOnFinish(func() {
		defer wg.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v\n", err)
		}

		fmt.Println("Shutdown: game finished")
	})

	err := server.Start()
	wg.Wait()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	fmt.Println("Finished game")
}
