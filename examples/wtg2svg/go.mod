module github.com/owulveryck/wardleyToGo/wtg2svg

go 1.19

require (
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/owulveryck/wardleyToGo v0.0.0
)

replace github.com/owulveryck/wardleyToGo v0.0.0 => ../../

require (
	golang.org/x/image v0.3.0 // indirect
	gonum.org/v1/gonum v0.12.0 // indirect
)
