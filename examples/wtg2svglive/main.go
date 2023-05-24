package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/kelseyhightower/envconfig"
	svgmap "github.com/owulveryck/wardleyToGo/encoding/svg"
	"github.com/owulveryck/wardleyToGo/parser/wtg"
	"nhooyr.io/websocket"
)

type configuration struct {
	Width  int    `default:"1500"`
	Height int    `default:"900"`
	Port   string `default:"8080"`
}

var config configuration

//go:embed assets/*
var assets embed.FS

var fullPath string

func main() {
	err := envconfig.Process("WTGLIVE", &config)
	if err != nil {
		envconfig.Usage("WTGLIVE", &config)
		log.Fatal(err.Error())
	}

	if len(os.Args) != 2 {
		envconfig.Usage("WTGLIVE", &config)
		log.Fatalf("usage: %v [wtg file to watch]", os.Args[0])
	}
	fileToWatch := os.Args[1]
	// Get the directory
	fullPath, err = filepath.Abs(fileToWatch)
	if err != nil {
		log.Fatal(err)
	}
	pathToWatch, err := filepath.Abs(filepath.Dir(fileToWatch))
	if err != nil {
		log.Fatal(err)
	}
	commC := make(chan string)

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		// send initial image
		commC <- fullPath
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					if event.Name == fullPath {
						log.Println(event.Name)
						commC <- event.Name
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	log.Println("watching " + fullPath)
	err = watcher.Add(pathToWatch)
	if err != nil {
		log.Fatal(err)
	}

	ws := &wsWriter{
		C: commC,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", ws.handler)

	myFs, err := fs.Sub(assets, "assets")
	if err != nil {
		log.Fatal(err)
	}
	assetsFs := http.FileServer(http.FS(myFs))

	mux.Handle("/", http.StripPrefix("/", assetsFs))
	log.Println("listening on " + config.Port + ". Use the PORT env var to change it")
	openbrowser("http://localhost:" + config.Port)
	err = http.ListenAndServe(":"+config.Port, mux)
	log.Fatal(err)
}

type wsWriter struct {
	C <-chan string
}

// This handler demonstrates how to correctly handle a write only WebSocket connection.
// i.e you only expect to write messages and do not expect to read any messages.
func (ws *wsWriter) handler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	ctx := c.CloseRead(r.Context())

	for {
		select {
		case <-ctx.Done():
			c.Close(websocket.StatusNormalClosure, "")
			return
		case file := <-ws.C:
			svg, err := generateSVG(file)
			if err != nil {
				svg = []byte(err.Error())
				log.Println(err)
			}
			w, err := c.Writer(ctx, websocket.MessageText)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Fprintf(w, "%s", svg)
			w.Close()
		}
	}
}

func generateSVG(filePath string) ([]byte, error) {
	p := wtg.NewParser()
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	err = p.Parse(f)
	if err != nil {
		return nil, err
	}
	if len(p.InvalidEntries) != 0 {
		for _, err := range p.InvalidEntries {
			log.Println(err)
		}
	}

	imgArea := (p.ImageSize.Max.X - p.ImageSize.Min.X) * (p.ImageSize.Max.X - p.ImageSize.Min.Y)
	canvasArea := (p.MapSize.Max.X - p.MapSize.Min.X) * (p.MapSize.Max.X - p.MapSize.Min.Y)
	if imgArea == 0 || canvasArea == 0 {
		p.ImageSize = image.Rect(0, 0, config.Width, config.Height)
		p.MapSize = image.Rect(30, 50, config.Width-30, config.Height-50)
	}
	var output bytes.Buffer
	e, err := svgmap.NewEncoder(&output, p.ImageSize, p.MapSize)
	if err != nil {
		return nil, err
	}
	defer e.Close()
	style := svgmap.NewOctoStyle(p.EvolutionStages)
	style.WithControls = true
	e.Init(style)
	err = e.Encode(p.WMap)
	if err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func openbrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err

}
