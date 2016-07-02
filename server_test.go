package elsa

import (
	json "github.com/gorilla/rpc/v2/json2"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeKovr/elsa/elclient"
	"github.com/LeKovr/go-base/logger"
)

type Service1Request struct {
	A int
	B int
}

type Service1Response struct {
	Result int
}

type Service1 struct{}

func (t *Service1) Multiply(r *http.Request, req *Service1Request, res *Service1Response) error {
	res.Result = req.A * req.B
	return nil
}

// -----------------------------------------------------------------------------
// Конструктор сервера API
func NewApp() *Service1 {
	a := &Service1{}
	return a
}

// -----------------------------------------------------------------------------
func TestService(t *testing.T) {

	log, err := logger.New()
	if err != nil {
		t.Error("Expected err to be nil, but got:", err)
	}

	addr := "localhost:8080"
	s, _ := New(addr, log)
	s.Handle("/api", APIServer(s.RPC, log))

	//	s := NewServer()
	//	s.Router.Handle("/api", s.Chain.Then(APIServer(s.RPC)))

	app := NewApp()
	s.RPC.RegisterService(app, "App")

	// стартуем сразу, чтобы успел развернуться до второго теста
	go s.RunServer()

	req, err := elclient.Request("http://"+addr+"/api", "App.Multiply", &Service1Request{4, 2})
	if err != nil {
		t.Error("Expected err to be nil, but got:", err)
	}
	req.RemoteAddr = addr // "127.0.0.0"
	w := httptest.NewRecorder()
	s.Router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected http response code 200, but got %v", w.Code)
	}

	var res Service1Response

	// log.Printf("Req: %+v  BODY: %+v", req, w.Body)

	err = json.DecodeClientResponse(w.Body, &res)
	if err != nil {
		t.Errorf("Couldn't decode response. %s", err)
	}
	if res.Result != 8 {
		t.Errorf("Wrong response: %v.", res.Result)
	}

	// Now network mode

	req, err = elclient.Request("http://"+addr+"/api", "App.Multiply", &Service1Request{4, 2})
	if err != nil {
		t.Error("Expected err to be nil, but got:", err)
	}

	var res1 Service1Response

	_, err = elclient.Call(req, &res1)
	if err != nil {
		t.Errorf("Call error: %s", err)
	}
	if res1.Result != 8 {
		t.Errorf("Wrong response: %v.", res1.Result)
	}

}
