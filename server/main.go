package main

import (
	"log"
	"math/rand"
	"net/http"

	c "github.com/freddie-nelson/scuffed-chess/server/chess"
	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/cors"
)

var games map[string]*c.GameController

func main() {
	games = make(map[string]*c.GameController)
	server := socketio.NewServer(nil)

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("error: ", e)
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("connected: ", s.ID())

		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	server.OnEvent("/", "game:create", func(s socketio.Conn, username string) string {
		code := ""
		for i := 0; i < 6; i++ {
			var char rune = 'a' + rune(rand.Intn(25))
			code += string(char)
		}

		if _, exists := games[code]; exists {
			return ""
		}

		g := c.NewGame(code)
		p := c.NewPlayer(username, false, s)
		g.You = p

		games[code] = g
		log.Printf("create game (%s): %s \n", username, code)

		return code
	})

	server.OnEvent("/", "game:join", func(s socketio.Conn, username string, code string) string {
		if _, exists := games[code]; !exists {
			return ""
		}

		g := games[code]
		p := c.NewPlayer(username, true, s)
		g.Opponent = p

		g.BroadcastData()
		log.Printf("join game (%s): %s \n", username, code)

		return code
	})

	go server.Serve()
	defer server.Close()

	mux := http.NewServeMux()
	mux.Handle("/socket.io/", server)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://192.168.1.84:8080"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	handler := c.Handler(mux)

	log.Println("Serving at localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}
