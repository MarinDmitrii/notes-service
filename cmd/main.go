package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/MarinDmitrii/notes-service/internal/common"
	noteBuilder "github.com/MarinDmitrii/notes-service/internal/note/builder"
	notePorts "github.com/MarinDmitrii/notes-service/internal/note/ports"
	userBuilder "github.com/MarinDmitrii/notes-service/internal/user/builder"
	userPorts "github.com/MarinDmitrii/notes-service/internal/user/ports"
)

type HttpServer struct {
	notePorts.HttpNoteHandler
	userPorts.HttpUserHandler
}

type Application struct {
	httpServer *http.Server
}

func (a *Application) Run(addr string, debug bool) error {
	router := http.NewServeMux()

	ctx := context.Background()

	userApp, userCleanup := userBuilder.NewApplication(ctx)
	userHttpHandler := userPorts.NewHttpUserHandler(userApp)
	userPorts.CustomRegisterHandlers(router, userHttpHandler)
	defer userCleanup()

	authMiddleware := common.NewAuthMiddleware(userApp.GetUserByEmail)

	noteApp, noteCleanup := noteBuilder.NewApplication(ctx)
	noteHttpHandler := notePorts.NewHttpNoteHandler(noteApp)
	notePorts.CustomRegisterHandlers(router, noteHttpHandler, authMiddleware)
	defer noteCleanup()

	a.httpServer = &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    180 * time.Second,
		WriteTimeout:   180 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Server is running...")

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func main() {
	app := &Application{}
	err := app.Run(":9090", false)
	if err != nil {
		panic(err)
	}

}
