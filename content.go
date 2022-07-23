package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (s *Server) doRequest(method string, path string, data []byte, headers map[string]string) ([]byte, int, error) {

	var req *http.Request
	var err error

	if method == http.MethodGet {
		req, err = http.NewRequest(http.MethodGet, path, nil)
	} else if method == http.MethodPost {
		req, err = http.NewRequest(http.MethodPost, path, bytes.NewBuffer(data))
	} else if method == http.MethodDelete {
		req, err = http.NewRequest(http.MethodDelete, path, bytes.NewBuffer(data))
	} else {
		return nil, 0, fmt.Errorf("method don't supported")
	}
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	var res *http.Response
	res, err = http.DefaultClient.Do(req)

	if err != nil {
		return nil, 0, err
	}
	log.Printf("got %d status code from %s", res.StatusCode, path)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, res.StatusCode, nil
}

func (s *Server) requestContent(method string, path string, data []byte, headers map[string]string) ([]byte, int, error) {
	fullPath := s.contentURL + path
	return s.doRequest(method, fullPath, data, headers)
}
