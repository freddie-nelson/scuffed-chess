package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

	c "github.com/freddie-nelson/scuffed-chess/server/chess"
	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/cors"
)

const PRODUCTION = true

var games map[string]*c.GameController
var players map[string]string

func main() {
	// disable logging in production
	if PRODUCTION {
		log.SetOutput(ioutil.Discard)
	}

	games = make(map[string]*c.GameController)
	players = make(map[string]string)

	server := socketio.NewServer(nil)

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("error: ", e)
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		players[s.ID()] = ""

		log.Println("connected: ", s.ID())
		return nil
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {

		if players[s.ID()] != "" {
			g, exists := games[players[s.ID()]]
			if exists {
				isOpponent := false
				if g.Opponent != nil && g.Opponent.CompareID((s.ID())) {
					isOpponent = true
				}
				leaveGame(s, players[s.ID()], isOpponent)
			}
		}

		delete(players, s.ID())

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
		players[s.ID()] = code

		log.Printf("create game (%s): %s \n", username, code)

		return code
	})

	server.OnEvent("/", "game:join", func(s socketio.Conn, username string, code string) string {
		if _, exists := games[code]; !exists || games[code].Opponent != nil {
			return ""
		}

		g := games[code]
		p := c.NewPlayer(username, true, s)
		g.Opponent = p

		g.StartGame()
		g.BroadcastData()

		players[s.ID()] = code

		log.Printf("join game (%s): %s \n", username, code)

		return code
	})

	server.OnEvent("/", "game:move", func(s socketio.Conn, code string, file, rank, dFile, dRank, promotion int) bool {
		if _, exists := games[code]; !exists {
			return false
		}

		g := games[code]
		madeMove := false
		if (g.You.CompareID(s.ID()) && g.IsCurrentlyPlaying(g.You)) || (g.Opponent.CompareID(s.ID()) && g.IsCurrentlyPlaying(g.Opponent)) {
			madeMove = g.MakeMove(file, rank, dFile, dRank, promotion)
		}

		g.BroadcastData()
		log.Printf("made move (%s): (%d, %d) (%d, %d) \n", code, file, rank, dFile, dRank)

		return madeMove
	})

	server.OnEvent("/", "game:valid-moves", func(s socketio.Conn, code string, file, rank int) string {
		if _, exists := games[code]; !exists {
			return "[]"
		}

		g := games[code]
		moves := []c.Spot{}
		if g.You.CompareID(s.ID()) {
			moves = g.GetValidMoves(file, rank, c.Black)
		} else if g.Opponent.CompareID(s.ID()) {
			moves = g.GetValidMoves(file, rank, c.White)
		}

		if moves != nil {
			json := "["

			for i := 0; i < len(moves); i++ {
				if i != 0 {
					json += ","
				}

				json += fmt.Sprintf("{ \"file\": %d, \"rank\": %d }", moves[i].GetFile(), moves[i].GetRank())

			}

			return json + "]"
		} else {
			return "[]"
		}
	})

	server.OnEvent("/", "game:leave", func(s socketio.Conn, code string, isOpponent bool) bool {
		return leaveGame(s, code, isOpponent)
	})

	go server.Serve()
	defer server.Close()

	mux := http.NewServeMux()
	mux.Handle("/socket.io/", server)

	port := ":8000"
	allowedOrigins := []string{"http://localhost:8080", "http://192.168.1.84:8080"}

	if PRODUCTION {
		port = ":" + os.Getenv("PORT")
		allowedOrigins = []string{"https://scuffedchess.netlify.app"}
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	handler := c.Handler(mux)

	fmt.Println("Serving at localhost" + port)
	fmt.Println(http.ListenAndServe(port, handler))
}

func leaveGame(s socketio.Conn, code string, isOpponent bool) bool {
	if _, exists := games[code]; !exists {
		return false
	}

	g := games[code]
	delete(games, code)

	players[s.ID()] = ""

	left := false
	if isOpponent && g.Opponent != nil && g.Opponent.CompareID(s.ID()) {
		g.Opponent = nil

		os := g.You.GetSocket()
		os.Emit("game:end-state", "disconnect")
		players[os.ID()] = ""

		left = true
	} else if g.You != nil && g.You.CompareID((s.ID())) {
		g.You = nil

		os := g.You.GetSocket()
		os.Emit("game:end-state", "disconnect")
		players[os.ID()] = ""

		left = true
	}

	return left
}
