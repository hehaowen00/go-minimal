package main

import (
	minimal "go-minimal"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
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

	go router.Serve(":8080")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
