package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (s *Server) requestContent(method string, path string, data []byte, headers map[string]string) ([]byte, error) {
	fullPath := s.contentURL + path

	var req *http.Request
	var err error

	if method == http.MethodGet {
		req, err = http.NewRequest(http.MethodGet, fullPath, nil)
	} else if method == http.MethodPost {
		req, err = http.NewRequest(http.MethodPost, fullPath, bytes.NewBuffer(data))
	} else if method == http.MethodDelete {
		req, err = http.NewRequest(http.MethodDelete, fullPath, bytes.NewBuffer(data))
	} else {
		return nil, fmt.Errorf("method don't supported")
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	var res *http.Response
	res, err = http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	log.Printf("got %d status code from content", res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *Server) getUserFull(username string) (*UserRespFull, error) {
	headers := make(map[string]string)
	headers["Username"] = username

	content, err := s.requestContent(http.MethodGet, "api/internal/user/", nil, headers)
	if err != nil {
		return nil, err
	}
	user := &UserRespFull{}
	err = json.Unmarshal(content, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Server) registerUser(username string, password string) (*UserResp, error) {
	marshal, err := json.Marshal(&AuthRequest{Username: username, Password: password})
	if err != nil {
		return nil, err
	}
	content, err := s.requestContent(http.MethodPost, "api/public/user/", marshal, nil)
	if err != nil {
		return nil, err
	}

	res := &UserResp{}
	err = json.Unmarshal(content, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
