package main

import (
	_ "github.com/lib/pq"
	"github.com/milvuk/go-web/internal"
)

func main() {
	db := internal.GetDbHandle()

	listenAddr := internal.ViperEnvVariable("API_LISTEN_ADDR")

	srv := internal.APIServer{
		ListenAddr: listenAddr,
		DB:         db,
	}

	srv.Run()
}
