package server

import (
	"net/http"
)

// routes returns the routes that the service instance supports
// The code below demostrates middleware wrapping
func (s *Server) routes() {
	s.Router.HandlerFunc(http.MethodGet, "/unprotectedAPI",
		s.rlog(s.handleUnprotectedAPI()))
	s.Router.HandlerFunc(http.MethodGet, "/protectedAPI",
		s.rlog(s.authorize(s.handleProtectedAPI())))
	s.Router.HandlerFunc(http.MethodGet, "/admin",
		s.rlog(s.adminOnly(s.handleProtectedAPI())))
}
