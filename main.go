package main

import (
	"fmt"
	"github.com/devashar13/instagram-api/controllers"
	"log"
	"net/http"
)

func newUserHandler() *controllers.UserHandler {
	return &controllers.UserHandler{}

}

func newPostHandler() *controllers.PostHandler {
	return &controllers.PostHandler{}

}

func main() {

	uh := newUserHandler()
	ph := newPostHandler()

	port := ":5000"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	http.Handle("/users", uh)
	http.Handle("/users/", uh)
	http.Handle("/posts", ph)
	http.Handle("/posts/users/", ph)
	http.Handle("/posts/", ph)

	fmt.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
