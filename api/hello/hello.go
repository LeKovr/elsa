package hello

import (
	"log"
	"net/http"
)

// -----------------------------------------------------------------------------

type Args struct {
	Who string
}

type Reply struct {
	Message string
}

type Service struct {
	Log  *log.Logger
	Data string
}

// -----------------------------------------------------------------------------

func New(logger *log.Logger, key string) *Service {
	return &Service{Log: logger, Data: key}
}

// -----------------------------------------------------------------------------

func (srv *Service) Say(r *http.Request, args *Args, reply *Reply) error {
	srv.Log.Printf("debug: Say called for %s", args.Who)
	reply.Message = "Hello, " + args.Who + ", from " + srv.Data + "!"
	return nil
}
