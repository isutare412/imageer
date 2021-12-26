package http

import "github.com/isutare412/imageer/api/pkg/core/user"

type getGreetingRes struct {
	Message string `yaml:"msg" json:"msg"`
}

type createUserReq struct {
	GivenName  string `yaml:"givenName" json:"givenName"`
	FamilyName string `yaml:"familyName" json:"familyName"`
	Email      string `yaml:"email" json:"email"`
	Password   string `yaml:"password" json:"password"`
}

func (req *createUserReq) into() *user.User {
	return &user.User{
		GivenName:  req.GivenName,
		FamilyName: req.FamilyName,
		Email:      req.Email,
		Password:   req.Password + "-fake-hash", // TODO: Hash password
	}
}

type createUserRes struct {
	ID         int64  `yaml:"id" json:"id"`
	GivenName  string `yaml:"givenName" json:"givenName"`
	FamilyName string `yaml:"familyName" json:"familyName"`
	Email      string `yaml:"email" json:"email"`
	Credit     int64  `yaml:"credit" json:"credit"`
}

func (resp *createUserRes) from(user *user.User) {
	resp.ID = user.ID
	resp.GivenName = user.GivenName
	resp.FamilyName = user.FamilyName
	resp.Email = user.Email
	resp.Credit = user.Credit
}

type getUserRes struct {
	ID         int64  `yaml:"id" json:"id"`
	GivenName  string `yaml:"givenName" json:"givenName"`
	FamilyName string `yaml:"familyName" json:"familyName"`
	Email      string `yaml:"email" json:"email"`
	Credit     int64  `yaml:"credit" json:"credit"`
}

func (resp *getUserRes) from(user *user.User) {
	resp.ID = user.ID
	resp.GivenName = user.GivenName
	resp.FamilyName = user.FamilyName
	resp.Email = user.Email
	resp.Credit = user.Credit
}
