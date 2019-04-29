package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/kfreiman/imageservice/pkg/logging"
)

// Server allows to access the application over the network.
// This can be implemented in HTTP, gRPC or another stack.
type Server interface {
	Start(port int)
}

type server struct {
	logger logging.Logger
	router http.Handler
}

// NewHTTPServer return Server's HTTP implementation
func NewHTTPServer(logger logging.Logger, router http.Handler) Server {
	return &server{
		logger: logger,
		router: router,
	}
}

func (s *server) Start(port int) {
	listenErr := make(chan error, 1)

	httpServer := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		ReadTimeout:  time.Millisecond * 200,
		WriteTimeout: time.Second * 2,
		IdleTimeout:  time.Second * 2,
		Handler:      s.router,
	}
	go func() {
		s.logger.Infow("Server started", "port", port)
		listenErr <- httpServer.ListenAndServe()
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-listenErr:
			if err != nil {
				s.logger.Error(err)
				return
			}
			os.Exit(0)
		case <-osSignals:
			httpServer.SetKeepAlivesEnabled(false)
			timeout := time.Second * 5
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			if err := httpServer.Shutdown(ctx); err != nil {
				s.logger.Error(err)
			}
		}
	}
}
