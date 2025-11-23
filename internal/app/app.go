package app

import (
	"context"
	"errors"
	"net"
	"net/http"

	"pr-reviewer-assigment-service/internal/closer"
	"pr-reviewer-assigment-service/internal/config"

	"github.com/go-chi/chi/v5"
)

type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
	router          chi.Router
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.startHTTP()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initRouter,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) startHTTP() error {
	//a.logger.Info("HTTP Server initializing")
	//
	listener, err := net.Listen("tcp", a.serviceProvider.HTTPConfig().Address())
	//if err != nil {
	//	a.logger.Error("failed to create listener")
	//}
	//
	//a.logger.Info("CORS initializing")

	a.httpServer = &http.Server{
		Handler:      a.router,
		WriteTimeout: a.serviceProvider.HTTPConfig().Timeout(),
		ReadTimeout:  a.serviceProvider.HTTPConfig().Timeout(),
		IdleTimeout:  a.serviceProvider.HTTPConfig().IdleTimeout(),
	}

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			//log.err("server shutdown")
		default:
			//a.logger.Error("failed to start server")
		}
	}

	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		//a.logger.Error("failed to shutdown server")
	}

	return err
}
