package main

import (
	"context"
	"template-subscriber-go/config"
	"template-subscriber-go/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.Info("Starting ...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	var s server.Server

	if err := s.Create(ctx, config); err != nil {
		log.Fatal(err.Error())
	}

	errc := make(chan error, 1)

	go func(errc chan error) {
		log.Fatal(<-errc)
	}(errc)

	s.Serve(ctx, errc)
}
