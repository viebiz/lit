package main

import (
	"context"
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

func run(ctx context.Context) error {
	m, err := monitoring.New(monitoring.Config{ServerName: "websocket"})
	if err != nil {
		log.Printf("monitoring init failed: %v", err)
	}
	defer m.Flush(monitoring.DefaultFlushWait)

	ctx = monitoring.SetInContext(ctx, m)

	corsCfg := lit.NewCORSConfig([]string{"*"})
	srv := lit.NewHttpServer("localhost:8080", lit.Handler(ctx, corsCfg, registerRouter))

	return srv.Run()
}

func registerRouter(router lit.Router) {
	router.Get("/ping", func(ctx lit.Context) error {
		ctx.JSON(http.StatusOK, "pong")
		return nil
	})

	hub := websocket.NewPublicHub(context.Background(), "/echo", echoHandler)
	hub.On(router)
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
