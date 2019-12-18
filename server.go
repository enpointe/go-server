// Package server a basic implementation of a http server with
// protected and unprotected enpoint. Endpoints are protected using
// JWT Bear Token Authentication.
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

// Server http server instance
type Server struct {
	Router      *httprouter.Router
	authHandler *jwtmiddleware.JWTMiddleware
	config      *Config
}

// NewServer create a new server instance
func NewServer(config *Config) (*Server, error) {
	s := &Server{}
	s.Router = httprouter.New()
	s.config = config
	s.routes()
	s.authHandler = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(s.config.JWTKey), nil
		},
		// The middleware verifies that tokens are signed with the specific signing algorithm
		// Important to avoid security issues described here:
		// https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})
	return s, nil
}

// Decode data sent on the request
func (s *Server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	// Future proof decode
	return json.NewDecoder(r.Body).Decode(v)
}

// respond send a json encoded response if data interface provided.
// Otherwise simply return status
func (s *Server) respond(w http.ResponseWriter, r *http.Request,
	data interface{}, status int) {

	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			fmt.Fprintf(w, "{ message: %s code: %d}",
				err, http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(status)
	}
}

// responeError send a json encode response with the error and error status code
func (s *Server) respondError(w http.ResponseWriter, r *http.Request, err error, status int) {
	type response struct {
		Message string `json:"message,omitempty"`
		Code    int    `json:"code,ommitempty"`
	}
	s.respond(w, r, &response{err.Error(), status}, status)
}

// responseErrCode send http error code back
func (s *Server) responseErrCode(w http.ResponseWriter, r *http.Request, status int) {
	s.respond(w, r, nil, status)
}

// ServeHTTP implement method to allow Sever to become and Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// handleUnprotectedAPI endpoint example for unprotected data.
func (s *Server) handleUnprotectedAPI() http.HandlerFunc {
	type response struct {
		Payload string `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		payload := response{"some unprotected data to send back"}
		s.respond(w, r, payload, http.StatusOK)
	}
}

// handleProtectedAPI enpointe example for protected data. The route
// wraps this method with authorization handler.
func (s *Server) handleProtectedAPI() http.HandlerFunc {
	type response struct {
		Payload  string `json:"data"`
		Username string `json:"username"`
		IsAdmin  bool   `json:"is_admin"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Extrac the JWT Token from the request
		tokenStr, err := jwtmiddleware.FromAuthHeader(r)
		if err != nil {
			s.responseErrCode(w, r, http.StatusForbidden)
			return
		}
		// Retrieve our custom claims object
		claims, err := ParseClaims(tokenStr, []byte(s.config.JWTKey))
		payload := response{
			Payload:  "some data to send back",
			Username: claims.Username,
			IsAdmin:  claims.IsAdmin}
		// Return a json payload back to the caller
		s.respond(w, r, payload, http.StatusOK)
	}
}

// adminOnly handler function that ensure enpoint can only
// be accessed by a authorized admin user.
func (s *Server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.authHandler.CheckJWT(w, r)
		// If there was an error, do not continue.
		if err != nil {
			s.responseErrCode(w, r, http.StatusForbidden)
			return
		}

		// Extrac the JWT Token from the request
		tokenStr, err := jwtmiddleware.FromAuthHeader(r)
		if err != nil {
			s.responseErrCode(w, r, http.StatusForbidden)
			return
		}
		// Retrieve our custom claims object
		claims, err := ParseClaims(tokenStr, []byte(s.config.JWTKey))
		if err != nil || !claims.IsAdmin {
			s.responseErrCode(w, r, http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

// authorize ensure user is authorized to access enpoint
func (s *Server) authorize(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := s.authHandler.CheckJWT(w, r)
		// If there was an error, do not continue.
		if err != nil {
			s.responseErrCode(w, r, http.StatusForbidden)
			return
		}
		h(w, r)
	})
}

func (s *Server) rlog(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"%s %s called",
			r.Method,
			r.RequestURI)
		start := time.Now()

		// defer log closing message so it prints even if we panic
		defer log.Printf(
			"%s %s executed in %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
		h(w, r)
	}
}
