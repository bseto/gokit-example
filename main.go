package main

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"strings"

	"encoding/json"
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

type lowercaseRequest struct {
	S string `json:"s"`
}

type lowercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

type stringService struct{}

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

func (stringService) Lowercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToLower(s), nil
}

// ErrEmpty is returned when input string is empty
var ErrEmpty = errors.New("Empty string")

type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
	Lowercase(string) (string, error)
}

func makeLowercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(lowercaseRequest)
		v, err := svc.Lowercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return lowercaseResponse{v, ""}, nil
	}
}

func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.Count(req.S)
		return countResponse{v}, nil
	}
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeLowercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request lowercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func main() {
	svc := stringService{}

	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	lowercaseHandler := httptransport.NewServer(
		makeLowercaseEndpoint(svc),
		decodeLowercaseRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/lowercase", lowercaseHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
