package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sah4ez/usersvc/pkg/endpoints"
)

var (
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
	ErrEmptyToken = errors.New("token is empty")
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := endpoints.MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST    /users/                          adds another user
	// GET     /users/:id                       retrieves the given user by id
	// PATCH   /users/:id                       partial updated user information
	// GET	   /users/							get all users

	r.Methods("POST").Path("/users/").Handler(httptransport.NewServer(
		e.PostUserEndpoint,
		decodePostUserRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/users/{id}").Handler(httptransport.NewServer(
		e.GetUserEndpoint,
		decodeGetUserRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PATCH").Path("/users/{id}").Handler(httptransport.NewServer(
		e.PatchUserEndpoint,
		decodePatchUserRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/users/").Handler(httptransport.NewServer(
		e.GetUsersEndpoint,
		decodeGetUsersRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodePostUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postUserRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	err = Validate(r)
	if err != nil {
		return nil, err
	}
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getUserRequest{ID: id}, nil
}

func decodePatchUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	err = Validate(r)
	if err != nil {
		return nil, err
	}
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, err
	}
	return patchUserRequest{
		ID:   id,
		User: user,
	}, nil
}

func decodeGetUsersRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return getUsersRequest{}, nil
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
