package main

import (
	"context"
	"golang-chat/handlers"
	"golang-chat/lib"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		cancel()
	}()

	if err := serve(ctx, "127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}
}

func serve(ctx context.Context, addr string) error {
	stream := lib.NewStream()
	store := lib.NewStore()

	r := mux.NewRouter()
	r.Handle("/api/messages", handlers.NewRestApiHandler(store, stream)).
		Methods("GET", "POST")
	r.Handle("/ws/messages", handlers.NewWebSocketHandler(stream))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))
	http.Handle("/", r)

	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.Printf("Server listening at %s", addr)
	<-ctx.Done()
	log.Print("Server stopping")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream.Shutdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		return err
	}

	log.Print("Server shutdown")
	return nil
}
