// websocket.go
package websocket

import (
    "GoMessageApp/internal/models"
    "encoding/json"
    "log"
    "net/http"
    "sync"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type Client struct {
    conn *websocket.Conn
    userID uint
}

type Message struct {
    Type string      `json:"type"`
    Data interface{} `json:"data"`
}

var (
    clients    = make(map[*Client]bool)
    clientsMux sync.Mutex
)

func HandleWebSocket(c *gin.Context) {
    userInterface, _ := c.Get("user")
    currentUser, _ := userInterface.(models.User)

    conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println(err)
        return
    }

    client := &Client{conn: conn, userID: currentUser.ID}

    clientsMux.Lock()
    clients[client] = true
    clientsMux.Unlock()

    defer func() {
        clientsMux.Lock()
        delete(clients, client)
        clientsMux.Unlock()
        conn.Close()
    }()

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            break
        }
    }
}

func BroadcastMessage(message models.Message) {
    data, err := json.Marshal(Message{Type: "new_message", Data: message})
    if err != nil {
        log.Println("Error marshalling message:", err)
        return
    }

    clientsMux.Lock()
    defer clientsMux.Unlock()

    for client := range clients {
        if client.userID == message.SenderID || client.userID == message.ReceiverID {
            err := client.conn.WriteMessage(websocket.TextMessage, data)
            if err != nil {
                log.Println("Error sending message to client:", err)
                client.conn.Close()
                delete(clients, client)
            }
        }
    }
}

func BroadcastUserUpdate(user models.User) {
    data, err := json.Marshal(Message{Type: "user_update", Data: user})
    if err != nil {
        log.Println("Error marshalling user update:", err)
        return
    }

    clientsMux.Lock()
    defer clientsMux.Unlock()

    for client := range clients {
        err := client.conn.WriteMessage(websocket.TextMessage, data)
        if err != nil {
            log.Println("Error sending user update to client:", err)
            client.conn.Close()
            delete(clients, client)
        }
    }
}

