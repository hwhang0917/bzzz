package main

import (
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "0.0.0.0:8080", "HTTP service address")
var dist = flag.String("dist", "../client/dist", "Static files")

var upgrader = websocket.Upgrader{}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Mutex to protect clients map
	mutex sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			log.Printf("Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mutex.Unlock()
			log.Printf("Client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
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

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Printf("recv: %s", message)

		// Broadcast to all other clients (excluding sender)
		c.broadcastToOthers(message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}

// broadcastToOthers sends message to all clients except the sender
func (c *Client) broadcastToOthers(message []byte) {
	c.hub.mutex.RLock()
	defer c.hub.mutex.RUnlock()

	for client := range c.hub.clients {
		if client != c { // Skip the sender
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(c.hub.clients, client)
			}
		}
	}
}

// serveWS handles websocket requests from the peer.
func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func fileServerWithMIME(root http.FileSystem) http.Handler {
	fs := http.FileServer(root)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set appropriate MIME types
		if r.URL.Path != "/" {
			switch {
			case r.URL.Path[len(r.URL.Path)-3:] == ".js":
				w.Header().Set("Content-Type", "application/javascript")
			case r.URL.Path[len(r.URL.Path)-4:] == ".css":
				w.Header().Set("Content-Type", "text/css")
			case r.URL.Path[len(r.URL.Path)-5:] == ".html":
				w.Header().Set("Content-Type", "text/html")
			case r.URL.Path[len(r.URL.Path)-4:] == ".png":
				w.Header().Set("Content-Type", "image/png")
			case r.URL.Path[len(r.URL.Path)-4:] == ".jpg" || r.URL.Path[len(r.URL.Path)-5:] == ".jpeg":
				w.Header().Set("Content-Type", "image/jpeg")
			case r.URL.Path[len(r.URL.Path)-4:] == ".gif":
				w.Header().Set("Content-Type", "image/gif")
			case r.URL.Path[len(r.URL.Path)-4:] == ".svg":
				w.Header().Set("Content-Type", "image/svg+xml")
			}
		}
		fs.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	hub := newHub()
	go hub.run()

	fs := fileServerWithMIME(http.Dir(*dist))
	http.Handle("/", fs)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})

	log.Printf("WebSocket server starting on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
