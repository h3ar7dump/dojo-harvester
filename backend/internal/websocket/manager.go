package websocket

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/dojo-harvester/backend/internal/config"
	"github.com/dojo-harvester/backend/internal/logger"
	pb "github.com/dojo-harvester/backend/pkg/proto"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Configured via CORS in Gin instead
	},
}

type Client struct {
	hub  *Manager
	conn *websocket.Conn
	send chan []byte
}

type Manager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	cfg        *config.Config
	mu         sync.RWMutex
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 100),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		cfg:        cfg,
	}
}

func (m *Manager) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			m.mu.Lock()
			for client := range m.clients {
				close(client.send)
				delete(m.clients, client)
			}
			m.mu.Unlock()
			return
		case client := <-m.register:
			m.mu.Lock()
			if len(m.clients) < m.cfg.Telemetry.MaxConcurrentConnections {
				m.clients[client] = true
				logger.Get().Info("WebSocket client connected", zap.Int("total_connections", len(m.clients)))
			} else {
				logger.Get().Warn("Max WebSocket connections reached, rejecting")
				client.conn.Close()
			}
			m.mu.Unlock()
		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.send)
				logger.Get().Info("WebSocket client disconnected", zap.Int("total_connections", len(m.clients)))
			}
			m.mu.Unlock()
		case message := <-m.broadcast:
			m.mu.RLock()
			for client := range m.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(m.clients, client)
				}
			}
			m.mu.RUnlock()
		}
	}
}

func (m *Manager) BroadcastTelemetry(telemetry *pb.RobotTelemetry) error {
	data, err := proto.Marshal(telemetry)
	if err != nil {
		return err
	}

	// 1 byte type (0x01 for telemetry), followed by payload
	msg := make([]byte, len(data)+1)
	msg[0] = 0x01
	copy(msg[1:], data)

	m.broadcast <- msg
	return nil
}

func (m *Manager) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Get().Error("Failed to upgrade connection", zap.Error(err))
		return
	}

	client := &Client{
		hub:  m,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(int64(c.hub.cfg.Telemetry.MaxMessageSizeBytes))
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Get().Error("WebSocket error", zap.Error(err))
			}
			break
		}
		// For now, we mainly push data to clients. We can handle client messages here if needed.
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
