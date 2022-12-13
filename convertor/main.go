package main

import (
	// "github.com/pelletier/go-toml"
	"log"

	"github.com/fasthttp/router"
	"github.com/goccy/go-json"

	"github.com/valyala/fasthttp"

	aheaders "github.com/go-asphyxia/http/headers"
	amethods "github.com/go-asphyxia/http/methods"
	"github.com/go-asphyxia/middlewares/CORS"
)

type (
	Request struct {
		File []byte `json:"file"`
	}

	Server struct {
		router *router.Router
		cors   *CORS.CORS
	}
)

const (
	fs string = "./files"
)

func main() {
	err := New().Run()
	if err != nil {
		log.Panicln(err)
	}
}

func New() (server *Server) {
	server = &Server{
		cors: CORS.NewCORS(&CORS.Configuration{
			Origins: nil,
			Methods: []string{amethods.GET, amethods.POST, amethods.OPTIONS}, //, amethods.PUT, amethods.DELETE}
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

func convert() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		request := Request{}
		err := json.Unmarshal(ctx.Request.Body(), &request)
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}

		// cmd := exec.Command()
	}
}
