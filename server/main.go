package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

type Game struct {
  ID string `json:"id"`
  PlayerOne Player `json:"playerOne"` 
  PlayerTwo Player `json:"playerTwo"`
}

type Player struct {
  Name string
}

type JoinGameRequest struct {
  GameID string `json:"gameId"`
  PlayerName string `json:"playerName"`
}

var upgrader = websocket.Upgrader{}

/* Will change this to an mutex */
var games []Game

func setupHeaders(w http.ResponseWriter) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Headers", "*")
  w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
}

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
  setupHeaders(w)

  decoder := json.NewDecoder(r.Body)
  var gameJoinRequest JoinGameRequest
  err := decoder.Decode(&gameJoinRequest)

  if err != nil {
    response := "Unprocessable request"
    w.Write([]byte(response))
    log.Println(response)
    return
  }

  var gameFound Game
  for _, game := range games {
    if game.ID == gameJoinRequest.GameID {
      if game.PlayerOne.Name != "" {
        game.PlayerTwo = Player{
          Name: gameJoinRequest.PlayerName,
        }

        gameFound = game 

        break
      }

      game.PlayerOne = Player{
          Name: gameJoinRequest.PlayerName,
        }

      gameFound = game
    }
  }

  log.Println(gameFound)

  encoder := json.NewEncoder(w)
  err = encoder.Encode(gameFound)

  if err != nil {
    response := "Internal Server Error"
    log.Println(response)
  }
}

func createGame(w http.ResponseWriter, r *http.Request) {
  setupHeaders(w)

  game := Game{
    ID: uuid.New().String(),
  }

  games = append(games, game)

  encoder := json.NewEncoder(w)
  err := encoder.Encode(game)
  
  if err != nil {
    response := "Internal Server Error"
    w.Write([]byte(response))
    log.Println(response)
    return
  }
}

func main() {
  http.HandleFunc("/", home)
  http.HandleFunc("/create-game", createGame)
  http.HandleFunc("/join-game", joinGame)
  http.HandleFunc("/ws", wsConnect)

  log.Printf("Listening on port: 8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
