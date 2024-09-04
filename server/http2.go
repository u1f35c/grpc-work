package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/u1f35c/grpc-test/testservice"
)

type testServiceHTTPServer struct {
	*http.Server
}

func getStatus(ctx context.Context, _ testservice.JSONRequest) *testservice.JSONResponse {
	ss := ctx.Value(statusStoreCtxKey{})
	if ss == nil {
		return nil
	}

	value := ss.(*StatusStore).GetStatus()

	return &testservice.JSONResponse{Value: value}
}

func setStatus(ctx context.Context, req testservice.JSONRequest) *testservice.JSONResponse {
	ss := ctx.Value(statusStoreCtxKey{})
	if ss == nil {
		return nil
	}

	value := ss.(*StatusStore).SetStatus(int(req.Value))

	return &testservice.JSONResponse{Value: value}
}

func (s *testServiceHTTPServer) statusStore(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var rpcReq testservice.JSONRequest
	err := decoder.Decode(&rpcReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not decode JSON request"))
		return
	}

	ctx := req.Context()

	resp := &testservice.JSONResponse{}
	if rpcReq.Action == "GetStatus" {
		resp = getStatus(ctx, rpcReq)
	} else if rpcReq.Action == "SetStatus" {
		resp = setStatus(ctx, rpcReq)
	}

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	err = encoder.Encode(resp)
	if err != nil {
		w.Write([]byte("Could not encode JSON request"))
		return
	}
}

// setHTTP2Context sets up a StatusStore context for a new connection
func setHTTP2Context(ctx context.Context, c net.Conn) context.Context {
	return context.WithValue(ctx, statusStoreCtxKey{}, &StatusStore{})
}

// HTTP2Serve starts an HTTP/2 server on the given port
func HTTP2Serve(port int) error {
	s := &testServiceHTTPServer{}

	router := http.NewServeMux()
	router.HandleFunc("/statusstore", s.statusStore)

	s.Server = &http.Server{
		Addr:        fmt.Sprintf(":%d", port),
		Handler:     router,
		ConnContext: setHTTP2Context,
	}

	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	fmt.Printf("Starting HTTPS server\n")

	return s.Serve(ln)
}
