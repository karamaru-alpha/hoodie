package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/karamaru-alpha/days/pkg/domain/config"
)

func NewServer(
	cfg *config.APIConfig,
	handler http.Handler,
) *http.Server {
	h2s := &http2.Server{}
	h1s := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Port),
		Handler:           h2c.NewHandler(handler, h2s),
		ReadHeaderTimeout: 30 * time.Second,
	}
	return h1s
}

func Serve(s *http.Server) (start, stop func(ctx context.Context) error) {
	start = func(ctx context.Context) error {
		ln, err := net.Listen("tcp", s.Addr)
		if err != nil {
			return derrors.Wrap(err, derrors.Internal, err.Error())
		}

		go func() {
			log.Printf("starting serve... addr=%s", s.Addr)
			if err := s.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal("failed to serve.", err)
			}
		}()
		return nil
	}
	stop = func(ctx context.Context) error {
		if err := s.Shutdown(ctx); err != nil {
			return derrors.Wrap(err, derrors.Internal, err.Error())
		}
		log.Printf("finish serve.")
		return nil
	}

	return start, stop
}
