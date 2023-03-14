module github.com/owulveryck/wardleyToGo/wtg2svglive

go 1.19

require (
	github.com/fsnotify/fsnotify v1.6.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/owulveryck/wardleyToGo v0.0.0
	nhooyr.io/websocket v1.8.7
)

replace github.com/owulveryck/wardleyToGo v0.0.0 => ../../

require (
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/klauspost/compress v1.10.3 // indirect
	golang.org/x/image v0.3.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
	gonum.org/v1/gonum v0.12.0 // indirect
)
