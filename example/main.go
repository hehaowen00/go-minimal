package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	minimal "github.com/hehaowen00/go-minimal"
)

func main() {
	router := minimal.NewRouter()

	router.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next(w, r)
			end := time.Now()

			log.Printf("%s %s %s %dus", r.Method, r.URL.String(), r.RemoteAddr, end.Sub(start).Microseconds())
		}
	})

	router.Use(minimal.CorsMiddleware(minimal.CorsOptions{
		AllowedOrigins: []string{
			"localhost:8080",
		},
	}))

	router.GET("/",
		func(w http.ResponseWriter, r *http.Request) {
			err := minimal.MarshalJSON(w, http.StatusOK, minimal.JSON{
				"data": "hello, world!",
			})
			if err != nil {
				log.Println(err)
				minimal.MarshalJSON(w, http.StatusInternalServerError, err.Error())
				return
			}
		},
		minimal.GzipMiddleware,
	)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}

	go func() {
		log.Println("server started", server.Addr)

		if err := server.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error - %v", err)
	}
}
