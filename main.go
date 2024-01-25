package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
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
    err = connection.WriteMessage(mt, message)
    if err != nil {
      log.Println("write:", err)
      break
    }
  }
}

func home(w http.ResponseWriter, r *http.Request) {
  var buffer []byte
  r.Body.Read(buffer)
  fmt.Printf("%b", buffer)

  response := "Hello!"
  w.Write([]byte(response))
}

func main() {
  http.HandleFunc("/", home)
  http.HandleFunc("/echo", echo)

  log.Printf("Listening on port: 8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
