package main

//go:generate buf generate

import (
	"log"
	"net/http"
	"strings"

	"study4cash/impl"

	"connectrpc.com/grpcreflect"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	db := configDB()

	mux := http.NewServeMux()
	// users
	usersPath, usersHandler := impl.NewUsersServer(db)
	log.Println(usersPath)
	mux.Handle(usersPath, usersHandler)

	// Reflectors
	reflector := grpcreflect.NewStaticReflector(
		strings.ReplaceAll(usersPath, "/", ""),
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	p := new(http.Protocols)
	p.SetHTTP1(true)
	p.SetUnencryptedHTTP2(true)
	s := http.Server{
		Addr:      "localhost:8080",
		Handler:   mux,
		Protocols: p,
	}
	log.Println("Server is running on port :8080 with HTTP and gRPC support")
	s.ListenAndServe()
}
