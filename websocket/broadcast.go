package websocket

import (
	"context"
	"log"

	"github.com/coder/websocket"
)

func (i *impl) Broadcast(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping broadcast...")
			return
		case e := <-i.broadcast:
			i.clients.Range(func(key, value interface{}) bool {
				client := value.(*websocket.Conn)
				if err := client.Write(ctx, e.Type, e.Content); err != nil {
					log.Println("Write error:", err)
					client.Close(websocket.StatusGoingAway, "Client disconnected")
					i.clients.Delete(key)
				}
				return true
			})
		}
	}
}
