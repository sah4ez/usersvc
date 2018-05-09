package usersvc

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type Endpoints struct {
	PostUserEndpoint  endpoint.Endpoint
	GetUserEndpoint   endpoint.Endpoint
	PatchUserEndpoint endpoint.Endpoint
	GetUsersEndpoint  endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostUserEndpoint:  MakePostUserEndpoint(s),
		GetUserEndpoint:   MakeGetUserEndpoint(s),
		PatchUserEndpoint: MakePatchUserEndpoint(s),
		GetUsersEndpoint:  MakeGetUsersEndpoint(s),
	}
}

func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	return Endpoints{
		PostUserEndpoint:  httptransport.NewClient("POST", tgt, encodePostUserRequest, decodePostUserResponse, options...).Endpoint(),
		GetUserEndpoint:   httptransport.NewClient("GET", tgt, encodeGetUserRequest, decodeGetUserResponse, options...).Endpoint(),
		PatchUserEndpoint: httptransport.NewClient("PATCH", tgt, encodePatchUserRequest, decodePatchUserResponse, options...).Endpoint(),
		GetUsersEndpoint:  httptransport.NewClient("GET", tgt, encodeGetUsersRequest, decodeGetUsersResponse, options...).Endpoint(),
	}, nil
}

func (e Endpoints) PostUser(ctx context.Context, p User) error {
	request := postUserRequest{User: p}
	response, err := e.PostUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(postUserResponse)
	return resp.Err
}

func (e Endpoints) GetUser(ctx context.Context, id string) (User, error) {
	request := getUserRequest{ID: id}
	response, err := e.GetUserEndpoint(ctx, request)
	if err != nil {
		return User{}, err
	}
	resp := response.(getUserResponse)
	return resp.User, resp.Err
}

func (e Endpoints) PatchUser(ctx context.Context, id string, p User) error {
	request := patchUserRequest{ID: id, User: p}
	response, err := e.PatchUserEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(patchUserResponse)
	return resp.Err
}

func MakePostUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postUserRequest)
		e := s.PostUser(ctx, req.User)
		return postUserResponse{Err: e}, nil
	}
}

func MakeGetUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getUserRequest)
		p, e := s.GetUser(ctx, req.ID)
		return getUserResponse{User: p, Err: e}, nil
	}
}

func MakePatchUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(patchUserRequest)
		e := s.PatchUser(ctx, req.ID, req.User)
		return patchUserResponse{Err: e}, nil
	}
}

func MakeGetUsersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//req := request.(patchUserRequest)
		users, e := s.GetUsers(ctx)
		return getUsersResponse{Err: e, Users: users}, nil
	}
}

type postUserRequest struct {
	User User
}

type postUserResponse struct {
	Err error `json:"err,omitempty"`
}

func (r postUserResponse) error() error { return r.Err }

type getUserRequest struct {
	ID string
}

type getUserResponse struct {
	User User  `json:"profile,omitempty"`
	Err  error `json:"err,omitempty"`
}

func (r getUserResponse) error() error { return r.Err }

type patchUserRequest struct {
	ID   string
	User User
}

type patchUserResponse struct {
	Err error `json:"err,omitempty"`
}

func (r patchUserResponse) error() error { return r.Err }

type getUsersRequest struct{}

type getUsersResponse struct {
	Err   error  `json:"err,omitempty"`
	Users []User `json:"users"`
}

func (r getUsersResponse) error() error { return r.Err }
