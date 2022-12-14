package main

import (
	"log"
	"os"

	"github.com/fasthttp/router"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pelletier/go-toml"

	"github.com/go-asphyxia/middlewares/CORS"

	"github.com/LineEast/stream/fileServer/internal/database"
	"github.com/LineEast/stream/fileServer/internal/server"
)

type (
	Configuration struct {
		DSN  string `toml:"DSN"`
		Host string `toml:"host"`
	}

	Server struct {
		router *router.Router
		cors   *CORS.CORS
		db     *pgxpool.Pool
	}
)

func main() {
	file, err := os.OpenFile(".configuration", os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}

	configuration := &Configuration{}
	err = toml.NewDecoder(file).Decode(configuration)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Conn(configuration.DSN)
	if err != nil {
		panic(err)
	}

	server, err := server.New(db)
	if err != nil {
		log.Panicln(err)
	}

	server.Run()
}
