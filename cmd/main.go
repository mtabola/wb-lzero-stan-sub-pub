package main

import (
	"flag"
	"log"
	"os"
	"stan/config"
	"stan/internal/clients"
)

func main() {
	usage := `
	Please, enter STAN client mode like
		go run main.go -m [--mode] (pub | sub)
	If you use pub mode, please enter uploaded file like:
		go run main.go -m [--mode] pub "dir/example.json"
	`

	var mode string

	flag.StringVar(&mode, "m", "mode", "Setting working mode")
	flag.StringVar(&mode, "mode", "mode", "Setting working mode")

	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 && mode == "pub" {
		log.Fatal(usage)
	}

	os.Setenv("CONFIG_PATH", "../config/config.yaml")
	cfg := config.MustLoadConfig()

	switch mode {
	case "sub":
		cfg.User = "subscriber"
		clients.NewSubscriberListener(cfg)
	case "pub":
		cfg.User = "publisher"
		clients.NewPublisher(cfg, args[0])
	default:
		log.Fatal(usage)
	}

}
