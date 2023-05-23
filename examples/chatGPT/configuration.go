package main

type configuration struct {
	ListenAddr string `envconfig:"listen_addr" default:"localhost:3333" desc:"Host to connect to, or ngrok to use tunneling" required:"true"`
	scheme     string
}
