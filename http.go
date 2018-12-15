// Package http contains an Endpoint implementation for the Git HTTP(s) protocol.
package http

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mughub/mughub/bare"
	"net"
	"net/http"

	"github.com/spf13/viper"
)

type endpoint struct {
	s *http.Server
}

// ListenAndServe serves the Git HTTP(s) Protocol.
func (e *endpoint) ListenAndServe(ctx context.Context) (err error) {
	return
}

// NewEndpoint returns an Endpoint implementation, which
// serves the Git HTTP(s) protocol.
//
func NewEndpoint(cfg *viper.Viper) (bare.Endpoint, bare.Router) {
	r := mux.NewRouter()

	s := &http.Server{
		Handler: r,
	}
	return &endpoint{s: s}, r
}

func getTCPListener(cfg *viper.Viper) net.Listener {
	addr := fmt.Sprintf("%s:%d", cfg.GetString("addr"), cfg.GetInt("port"))
	l, err := net.Listen("tcp", addr)
	if err != nil {
		l.Close()
		panic(err)
	}
	return l
}
