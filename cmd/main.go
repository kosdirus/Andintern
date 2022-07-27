package main

import (
	"github.com/kosdirus/andintern/internal/andintern"
	"github.com/kosdirus/andintern/internal/api/http/handler"
	"github.com/kosdirus/andintern/internal/config"
	"github.com/kosdirus/andintern/internal/database"
	"github.com/kosdirus/andintern/internal/database/dataprovider/pg"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("can't create config: %s", err.Error())
	}

	db, err := database.NewClient(*cfg)
	if err != nil {
		log.Fatalf("can't construct database: %s", err.Error())
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("error while closing database connection: %s", err.Error())
		}
	}()

	if err := db.Migrate(); err != nil {
		log.Fatalf("can't migrate the database")
	}

	txer := pg.NewTxManager(db)
	carStore := pg.NewCarStore(db, txer)
	core := andintern.NewCore(cfg, db, carStore, txer)

	apiServer, err := handler.NewServer(cfg, core)
	if err != nil {
		log.Fatalf("can't construct api http server: %s", err.Error())
	}

	apiServer.ListenAndServe()
}
