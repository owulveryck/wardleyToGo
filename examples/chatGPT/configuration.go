package main

type configuration struct {
	ListenAddr  string `envconfig:"listen_addr" default:"localhost:3333" desc:"Host to connect to, or ngrok to use tunneling" required:"true"`
	StoragePath string `envconfig:"storage" default:"memory" desc:"the path of the persistent storage or memory"`
	scheme      string
}
