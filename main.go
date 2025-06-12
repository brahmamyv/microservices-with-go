package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	helloHandler := handlers.NewHello(l)

	goodbyeHandler := handlers.NewGoodbye(l)

	sm := http.NewServeMux()

	sm.Handle("/", helloHandler)
	sm.Handle("/goodbye", goodbyeHandler)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal(err)
		}
	}()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, os.Kill)

	sig := <-signalChan
	l.Println("received a termination, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	server.Shutdown(tc)

}
