package router

import (
	"context"
	"crypto/tls"
	"net/http"
	"tem/internal/srv"

	"github.com/gin-gonic/gin"
)

type App struct {
	app    *gin.Engine
	srv    *srv.Srv
	server *http.Server
}

func NewRouter(srv *srv.Srv) *App {
	app := gin.New()

	app.Use(gin.Recovery())
	return &App{app: app, srv: srv}
}

func (app *App) Run() error {
	addr := app.srv.Config.Gin.Addr + ":" + app.srv.Config.Gin.Port
	app.server = &http.Server{Addr: addr}
	if app.srv.Config.Gin.TlsCert == "" || app.srv.Config.Gin.TlsKey == "" {
		return app.server.ListenAndServe()
	} else {
		app.server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{

			}
		}
		app.app.RunTLS(
			app.srv.Config.Gin.Addr+":"+app.srv.Config.Gin.Port,
			app.srv.Config.Gin.TlsCert,
			app.srv.Config.Gin.TlsKey,
		)
	}
}

func (app *App) Stop() {
	app.server.Shutdown(context.Background())
}
