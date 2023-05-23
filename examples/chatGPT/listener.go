package main

import (
	"context"
	"net"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func setupListener(ctx context.Context, s *configuration) (net.Listener, error) {
	if s.ListenAddr == "ngrok" {
		l, err := ngrok.Listen(ctx,
			config.HTTPEndpoint(),
			ngrok.WithAuthtokenFromEnv(),
		)
		s.scheme = "https"
		s.ListenAddr = l.Addr().String()
		return l, err
	}
	s.scheme = "http"
	return net.Listen("tcp", s.ListenAddr)

}
