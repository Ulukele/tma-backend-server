package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (s *Server) requestAuth(path string, data []byte) ([]byte, error) {
	fullPath := s.authURL + path

	res, err := http.Post(fullPath, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	log.Printf("got %d status code from auth", res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *Server) requestSignIn(username string, userId uint) (*AuthResp, error) {
	req := TokenRequest{Username: username, Id: userId}
	marshal, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	auth, err := s.requestAuth("api/auth/sign-in/", marshal)
	if err != nil {
		return nil, err
	}

	res := &AuthResp{Id: userId}
	err = json.Unmarshal(auth, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Server) requestValidation(jwt string) (*ValidationResp, error) {
	req := ValidationRequest{JWT: jwt}
	marshal, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	auth, err := s.requestAuth("api/auth/validate/", marshal)
	if err != nil {
		return nil, err
	}

	res := &ValidationResp{}
	err = json.Unmarshal(auth, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
