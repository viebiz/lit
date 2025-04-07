package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/viebiz/lit"
	"github.com/viebiz/lit/monitoring"
	"github.com/viebiz/lit/websocket"
)

func main() {
	ctx := context.TODO()
	if err := run(ctx); err != nil {
		log.Printf("server exit abnormally: %v", err)
		return
	}
}

var (
	pingInterval = websocket.WithPingInterval[Request, Response]
)

func run(ctx context.Context) error {
	m, err := monitoring.New(monitoring.Config{ServerName: "simple-websocket"})
	if err != nil {
		log.Printf("monitoring init failed: %v", err)
	}
	defer m.Flush(monitoring.DefaultFlushWait)

	ctx = monitoring.SetInContext(ctx, m)

	hub := websocket.NewHub(context.Background(), "/message/:client_id", echoHandler,
		pingInterval(10),
	)

	r := lit.NewRouter(ctx)
	registerRouter(hub)(r)

	srv := lit.NewHttpServer("localhost:8080", r.Handler())

	return srv.Run()
}

func registerRouter(hub websocket.Hub[Request, Response]) func(lit.Router) {
	return func(router lit.Router) {
		hub.On(router)
		router.Post("/message", func(ctx lit.Context) error {
			var m struct {
				ClientID string `json:"client_id"`
				Message  string `json:"message"`
			}
			if err := ctx.Bind(&m); err != nil {
				return err
			}

			if err := hub.Emit(ctx, m.ClientID, Response{Message: m.Message}, websocket.EmitOption{}); err != nil {
				return err
			}
			return ctx.String(http.StatusOK, fmt.Sprintf("message sent to client %s successfully", m.ClientID))
		})

		router.Post("/message/broadcast", func(ctx lit.Context) error {
			var m struct {
				Message string `json:"message"`
			}
			if err := ctx.Bind(&m); err != nil {
				return err
			}

			if err := hub.Broadcast(ctx, Response{Message: m.Message}, websocket.EmitOption{}); err != nil {
				return err
			}
			return ctx.String(http.StatusOK, fmt.Sprintf("broadcast message successfully"))
		})
	}
}

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Message string `json:"message"`
}

func echoHandler(ctx context.Context, req Request) (Response, error) {
	return Response{
		Message: "echo:" + req.Message,
	}, nil
}
