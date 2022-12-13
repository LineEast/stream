package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/fasthttp/router"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pelletier/go-toml"
	"github.com/valyala/fasthttp"

	aheaders "github.com/go-asphyxia/http/headers"
	amethods "github.com/go-asphyxia/http/methods"
	"github.com/go-asphyxia/middlewares/CORS"
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

const (
	fs string = "./files"
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

	server, err := New(configuration)
	if err != nil {
		log.Panicln(err)
	}

	server.Run()
}

func New(configuration *Configuration) (server *Server, err error) {
	server = &Server{
		cors: CORS.NewCORS(&CORS.Configuration{
			Origins: nil,
			Methods: []string{amethods.GET, amethods.POST, amethods.OPTIONS, amethods.PUT}, //, amethods.DELETE}
			Headers: []string{aheaders.Accept, aheaders.ContentType},
		}),
	}

	r := router.New()
	r.OPTIONS("/{everything:*}", server.cors.Handler())
	r.ServeFiles("/", fs)
	server.router = r

	server.db, err = Conn(configuration.DSN)

	return
}

func (s *Server) Run() error {
	return fasthttp.ListenAndServe(":8080", s.cors.Middleware(s.router.Handler))
}

func Conn(DSN string) (db *pgxpool.Pool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err = pgxpool.New(ctx, DSN)
	if err != nil {
		panic(err)
	}

	return
}

func PutFile(db *pgxpool.Pool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

	}
}
