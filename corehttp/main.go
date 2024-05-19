package main

import (
	"log/slog"
	"os"

	"github.com/chaitanyamaili/learn-golang/corehttp/api"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(log)

	srv := api.NewAppServer(":8080", log)
	err := srv.Run()
	if err != nil {
		log.Error(err.Error())
	}
}
