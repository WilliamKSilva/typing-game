package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

type Game struct {
  ID string
  PlayerOne Player 
  PlayerTwo Player 
}

type Player struct {
  Name string
}

type JoinGameRequest struct {
  GameID string
  PlayerName string
}

type CreateGameRequest struct {
  PlayerName string
}

var upgrader = websocket.Upgrader{}

/* Will change this to an mutex */
var games []Game

func wsConnect(w http.ResponseWriter, r *http.Request) {
  connection, err := upgrader.Upgrade(w, r, nil)

  if err != nil {
    log.Print("upgrade:", err)
    return
  }

  defer connection.Close()
  for {
    mt, message, err := connection.ReadMessage()
    if err != nil {
      log.Println("read:", err)
      break
    }
    log.Printf("recv: %s", message)
    err = connection.WriteMessage(mt, []byte("recv: done"))
    if err != nil {
      log.Println("write:", err)
      break
    }
  }
}

func home(w http.ResponseWriter, r *http.Request) {
  response := "Listening!"
  w.Write([]byte(response))
}

func joinGame(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var gameJoinRequest JoinGameRequest
  err := decoder.Decode(&gameJoinRequest)

  if err != nil {
    response := "Unprocessable request"
    w.Write([]byte(response))
    log.Println(response)
    return
  }

  for _, game := range games {
    if game.ID == gameJoinRequest.GameID {
      game.PlayerTwo = Player{
        Name: gameJoinRequest.PlayerName,
      }
    }
  }
}

func createGame(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var gameCreateRequest CreateGameRequest 
  err := decoder.Decode(&gameCreateRequest)
  if err != nil {
    response := "Unprocessable request"
    w.Write([]byte(response))
    log.Println(response)
    return
  }

  playerOne := Player{
    Name: gameCreateRequest.PlayerName,
  }

  game := Game{
    ID: uuid.New().String(),
    PlayerOne: playerOne,
  }

  games = append(games, game)

  log.Println(len(games))
}

func main() {
  http.HandleFunc("/", home)
  http.HandleFunc("/create-game", createGame)
  http.HandleFunc("/join-game", joinGame)
  http.HandleFunc("/ws", wsConnect)

  log.Printf("Listening on port: 8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
