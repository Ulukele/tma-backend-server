package main

import (
	"encoding/json"
	"net/http"
)

func (s *Server) requestAuth(method string, path string, data []byte, headers map[string]string) ([]byte, int, error) {
	fullPath := s.authURL + path
	return s.doRequest(method, fullPath, data, headers)
}

func (s *Server) requestValidation(jwt string) (*ValidationResp, error) {
	req := ValidationRequest{JWT: jwt}
	marshal, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	auth, _, err := s.requestAuth(http.MethodPost, "api/v1/auth/validate/", marshal, nil)
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
