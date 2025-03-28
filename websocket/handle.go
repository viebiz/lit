package websocket

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/viebiz/lit"
)

func (i *impl) Handle(ctx lit.Context) {
	// Identify client id
	clientID := ctx.Query("client_id")
	if clientID == "" {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error_message": "client id is mandatory",
		})
		return
	}

	i.mutex.Lock()
	if i.count >= MaxConnections {
		i.mutex.Unlock()
		ctx.JSON(http.StatusTooManyRequests, "Too many connections")
		return
	}
	i.count++
	i.mutex.Unlock()

	conn, err := websocket.Accept(ctx.Writer(), ctx.Request(), &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	i.clients.Store(clientID, conn)
	log.Println("Client connected:", clientID)
	timeoutCtx, cancel := context.WithTimeout(i.ctx, PongWait)

	defer func() {
		cancel()
		conn.Close(websocket.StatusNormalClosure, "Closing connection")
		i.clients.Delete(clientID)
		i.mutex.Lock()
		i.count--
		i.mutex.Unlock()
		log.Println("Client disconnected:", clientID)
	}()
	go i.pingClient(timeoutCtx, clientID, conn)
	//conn.SetReadLimit(512)
	for {
		msgType, msgBytes, err := conn.Read(context.Background())
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("message type %d", msgType)
		// TODO something

		// Send msg to client back
		i.broadcast <- event{msgType, msgBytes}
	}
}

func (i *impl) sendToClientAPI(ctx lit.Context) {
	var req struct {
		ClientID string `json:"client_id"`
		Message  string `json:"message"`
	}

	if err := ctx.Bind(&req); err != nil {
		//ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if client, ok := i.clients.Load(req.ClientID); ok {
		if err := client.(*websocket.Conn).Write(nil, websocket.MessageText, []byte(req.Message)); err != nil {
			log.Println("Write error:", err)
			client.(*websocket.Conn).Close(websocket.StatusGoingAway, "Client disconnected")
			i.clients.Delete(req.ClientID)
		}
		ctx.JSON(http.StatusOK, "Message sent")
	} else {
		ctx.JSON(http.StatusNotFound, "Client not found")
	}
}

func (i *impl) pingClient(ctx context.Context, clientID string, conn *websocket.Conn) {
	ticker := time.NewTicker(PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := conn.Ping(ctx); err != nil {
				log.Println("Ping error, closing connection:", err)
				conn.Close(websocket.StatusGoingAway, "Ping failed")
				i.clients.Delete(clientID)
				return
			}
		}
	}
}
