package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (s *Server) requestContent(method string, path string, data []byte, usernameHeader string) ([]byte, error) {
	fullPath := s.contentURL + path

	var req *http.Request
	var err error

	if method == http.MethodGet {
		req, err = http.NewRequest(http.MethodGet, fullPath, nil)
	} else if method == http.MethodPost {
		req, err = http.NewRequest(http.MethodPost, fullPath, bytes.NewBuffer(data))
	} else if method == http.MethodDelete {
		req, err = http.NewRequest(http.MethodDelete, fullPath, bytes.NewBuffer(data))
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if usernameHeader != "" {
		req.Header.Set("Username", usernameHeader)
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

	if res.StatusCode != 200 {
		log.Printf("got %d status code from content. body: %s", res.StatusCode, body)
		return nil, err
	}

	return body, nil
}

func (s *Server) getUserFull(username string) (*UserRespFull, error) {
	content, err := s.requestContent("GET", "/internal/user/", nil, username)
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
	type UserReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	marshal, err := json.Marshal(&UserReq{Username: username, Password: password})
	if err != nil {
		return nil, err
	}
	content, err := s.requestContent("POST", "/public/user/", marshal, "")
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
