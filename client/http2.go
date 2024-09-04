package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/u1f35c/grpc-test/testservice"
)

func http2Call(client *http.Client, port int, rpcReq *testservice.JSONRequest) (*testservice.JSONResponse, error) {
	marshalled, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("http://localhost:%d/statusstore", port),
		bytes.NewReader(marshalled))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	var rpcResp testservice.JSONResponse
	err = decoder.Decode(&rpcResp)
	if err != nil {
		return nil, err
	}

	return &rpcResp, nil
}

// HTTP2Connect connects to an HTTP2 server on the given port
func HTTP2Connect(port int) error {
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}

	rpcReq := &testservice.JSONRequest{Action: "GetStatus"}
	rpcResp, err := http2Call(client, port, rpcReq)
	if err != nil {
		return err
	}
	fmt.Printf("Status: %d\n", rpcResp.Value)

	rpcReq = &testservice.JSONRequest{Action: "SetStatus", Value: 4}
	rpcResp, err = http2Call(client, port, rpcReq)
	if err != nil {
		return err
	}
	fmt.Printf("Status: %d\n", rpcResp.Value)

	rpcReq = &testservice.JSONRequest{Action: "GetStatus"}
	rpcResp, err = http2Call(client, port, rpcReq)
	if err != nil {
		return err
	}
	fmt.Printf("Status: %d\n", rpcResp.Value)

	return nil
}
