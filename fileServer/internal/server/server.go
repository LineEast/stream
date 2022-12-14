package server

import (
	"github.com/fasthttp/router"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"

	aheaders "github.com/go-asphyxia/http/headers"
	amethods "github.com/go-asphyxia/http/methods"
	"github.com/go-asphyxia/middlewares/CORS"
)

type (
	Server struct {
		router *router.Router
		cors   *CORS.CORS
		db     *pgxpool.Pool
	}
)

const (
	fs string = "./files"
)

func New(db *pgxpool.Pool) (server *Server, err error) {
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

	return
}

func (s *Server) Run() error {
	return fasthttp.ListenAndServe(":8080", s.cors.Middleware(s.router.Handler))
}
