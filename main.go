package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

//go:embed client/dist/*
var distFS embed.FS

var addr = flag.String("addr", "0.0.0.0:8080", "HTTP serving address")

var upgrader = websocket.Upgrader{}

type Hub struct {
    clients map[*Client]bool
    broadcast chan []byte
    register chan *Client
    unregister chan *Client
    mutex sync.RWMutex
}
func createNewHub() *Hub {
    return &Hub{
        clients: make(map[*Client]bool),
        broadcast: make(chan []byte),
        register: make(chan *Client),
        unregister: make(chan *Client),
    }
}
func (h *Hub) run() {
    for {
        select {
        case client := <- h.register:
            h.mutex.Lock()
            h.clients[client] = true
            h.mutex.Unlock()
        case client := <- h.unregister:
            h.mutex.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mutex.Unlock()
        case message := <- h.broadcast:
            h.mutex.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mutex.RUnlock()
        }
    }
}

type Client struct {
    hub *Hub
    conn *websocket.Conn
    send chan []byte
}
func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway,websocket.CloseAbnormalClosure) {
                log.Printf("Unexpected Close: %v\n", err)
            }
            break
        }
        c.broadcastToOthers(message)
    }
}

func (c *Client) writePump() {
    defer c.conn.Close()

    for {
        select {
        case message, ok := <- c.send:
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
            }
            if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
                log.Printf("Failed to write message: %v\n", err)
                return
            }
        }
    }
}

func (c *Client) broadcastToOthers(message []byte) {
    c.hub.mutex.RLock()
    defer c.hub.mutex.RUnlock()

    for client := range c.hub.clients {
        if client != c {
            select {
            case client.send <- message:
            default:
                close(client.send)
                delete(c.hub.clients, client)
            }
        }
    }
}

func serveWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade: %v\n", err)
        return
    }
    client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
    client.hub.register <- client

    go client.writePump()
    go client.readPump()
}

func main() {
    flag.Parse()
    log.SetFlags(0)

    hub := createNewHub()
    go hub.run()

    staticFS, err := fs.Sub(distFS, "client/dist")
    if err != nil {
        panic(err)
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        filePath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
        if filePath == "" {
            filePath = "index.html"
        }
        if _, err := staticFS.Open(filePath); err != nil {
            filePath = "index.html"
        }
        http.ServeFileFS(w, r, staticFS, filePath)
    })

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWebSocket(hub, w, r)
    })

    log.Printf("Server starting on %s", *addr)
    log.Fatal(http.ListenAndServe(*addr, nil))
}
