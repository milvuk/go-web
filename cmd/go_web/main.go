package main

import (
	_ "github.com/lib/pq"
	"github.com/milvuk/go-web/internal"
)

func main() {
	store := internal.NewPostgresStore()
	listenAddr := internal.ViperEnvVariable("API_LISTEN_ADDR")

	srv := internal.APIServer{
		ListenAddr: listenAddr,
		Store:      store,
	}

	srv.Run()
}
