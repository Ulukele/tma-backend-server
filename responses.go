package main

type UserResp struct {
	UserId   uint   `json:"id"`
	Username string `json:"username"`
}

type UserRespFull struct {
	UserResp
	Password string `json:"password"`
}

type AuthResp struct {
	Id  uint   `json:"id"`
	JWT string `json:"jwt"`
}

type ValidationResp struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}
