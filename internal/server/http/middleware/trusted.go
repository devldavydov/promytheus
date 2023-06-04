package middleware

import (
	"net"
	"net/http"
)

// Trusted is a middleware to check RemoteAddr according to trusted network.
type Trusted struct {
	trustedNetwork *net.IPNet
}

func NewTrusted(trustedNetwork *net.IPNet) *Trusted {
	return &Trusted{trustedNetwork: trustedNetwork}
}

func (t *Trusted) Handle(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if t.trustedNetwork != nil {
			remoteIP := net.ParseIP(r.RemoteAddr)
			if !t.trustedNetwork.Contains(remoteIP) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
