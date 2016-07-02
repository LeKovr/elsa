package elsa

import (
	"net"
	"net/http"

	"github.com/LeKovr/go-base/logger"
)

// -----------------------------------------------------------------------------
type apiHandler struct {
	srv   http.Handler
	log   *logger.Log
	hosts []string
}

// APIServer extends srv.ServeHTTP with OPTIONS support
func APIServer(srv http.Handler, log *logger.Log, host ...string) http.Handler {
	return &apiHandler{srv: srv, log: log, hosts: host}
}

// -----------------------------------------------------------------------------

// ServeHTTP with OPTIONS & Access-Control headers support
func (a *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if r.Header.Get("Content-Type") == "" {
		r.Header.Set("Content-Type", "application/json") // TODO: IE8 only
	}
	a.log.Printf("Got request: (%s) %v", origin, r)

	var host string
	if origin != "" && len(a.hosts) > 0 { // lookup if host is allowed
		for _, h := range a.hosts {
			if origin == h {
				host = h
				break
			}
		}
	} else {
		host = origin
	}
	if origin != "" && host == "" {
		a.log.Warningf("Unregistered request source: %s", origin)
		http.Error(w, "Origin not registered", http.StatusForbidden)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", host)
	w.Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	if r.Method != "OPTIONS" {
		ip, src := userIP(r)
		a.log.Debugf("ip %s source: %s", ip, src)
		r.Header.Set("Client-Ip", ip)
		a.srv.ServeHTTP(w, r)
	}
}

// -----------------------------------------------------------------------------

func userIP(r *http.Request) (ip string, ipSource string) {
	ip = r.Header.Get("X-Real-Ip")
	if ip != "" {
		return ip, "real-ip"
	}
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip, "fwd-for"
	}
	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	return ip, "rem-addr"
}
