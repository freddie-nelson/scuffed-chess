package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/cors"
)

func main() {
	server := socketio.NewServer(nil)

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("error:", e)
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("connected: ", s.ID())

		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	server.OnEvent("/", "game:create", func(s socketio.Conn, username string) bool {
		code := ""
		for i := 0; i < 6; i++ {
			var char rune = 'a' + rune(rand.Intn(25))
			code += string(char)
		}

		log.Printf("create game (%s): %s \n", username, code)

		return true
	})

	go server.Serve()
	defer server.Close()

	mux := http.NewServeMux()
	mux.Handle("/socket.io/", server)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://192.168.1.84:8080"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	handler := c.Handler(mux)

	log.Println("Serving at localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}
