package elsa

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bakins/net-http-recover"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	rpc "github.com/gorilla/rpc/v2"
	json "github.com/gorilla/rpc/v2/json2"
	"github.com/justinas/alice"

	"github.com/LeKovr/go-base/database"
	"github.com/LeKovr/go-base/logger"
)

// Version of elsa package
const Version = "1.1"

// Flags is a package flags sample
// in form ready for use with github.com/jessevdk/go-flags
type Flags struct {
	Addr  string   `long:"http_addr" default:"localhost:8080"  description:"Http listen address"`
	Hosts []string `long:"http_origin" description:"Allowed http origin(s)"`
}

// -----------------------------------------------------------------------------

// Server struct handles all server onjects
type Server struct {
	Listener net.Listener
	Router   *mux.Router
	RPC      *rpc.Server
	Chain    alice.Chain
	DB       *database.DB
	Log      *logger.Log
}

// DB is a database attr setter
func DB(db *database.DB) func(*Server) error {
	return func(s *Server) error {
		return s.setDB(db)
	}
}
func (s *Server) setDB(db *database.DB) error {
	s.DB = db
	return nil
}

// -----------------------------------------------------------------------------

// New returns initialized Server object
func New(addr string, log *logger.Log, options ...func(*Server) error) (*Server, error) {

	log.Printf("Listening http://%s", addr)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	srv := Server{Listener: l, Log: log}

	srv.Router = mux.NewRouter()
	rpc := rpc.NewServer()
	rpc.RegisterCodec(json.NewCodec(), "application/json")
	srv.Chain = alice.New(
		func(h http.Handler) http.Handler {
			return handlers.CombinedLoggingHandler(os.Stdout, h)
		},
		// handlers.CompressHandler,
		func(h http.Handler) http.Handler {
			return recovery.Handler(os.Stderr, h, true)
		})

	srv.RPC = rpc

	for _, option := range options {
		err := option(&srv)
		if err != nil {
			return nil, err
		}
	}

	return &srv, nil
}

// -----------------------------------------------------------------------------

// Handle adds chained route handler
func (s *Server) Handle(uri string, handler http.Handler) {
	s.Router.Handle(uri, s.Chain.Then(handler))
}

// -----------------------------------------------------------------------------

// Cleanup server objects
func (s *Server) Cleanup() {

}

// -----------------------------------------------------------------------------

// RunServer starts server's ListenAndServe
func (s *Server) RunServer() {

	// http://stackoverflow.com/questions/18106749/golang-catch-signals
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		s.Log.Infof("Got signal %v", sig)
		s.Cleanup()
		os.Exit(0)
	}()

	s.Log.Fatal(http.Serve(s.Listener, s.Router))
	return
}
