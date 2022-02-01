package http

type errorRes struct {
	Code    int    `yaml:"code" json:"code,omitempty" example:"400"`
	Message string `yaml:"msg" json:"msg,omitempty" example:"simple error message"`
}

type signInReq struct {
	Email    string `yaml:"email" json:"email"`
	Password string `yaml:"password" json:"password"`
}

type signInRes struct {
	Token string `yaml:"token" json:"token"`
}

type getGreetingRes struct {
	Message string `yaml:"msg" json:"msg"`
}

type createUserReq struct {
	GivenName  string `yaml:"givenName" json:"givenName"`
	FamilyName string `yaml:"familyName" json:"familyName"`
	Email      string `yaml:"email" json:"email"`
	Password   string `yaml:"password" json:"password"`
}

type createUserRes struct {
	ID         int64  `yaml:"id" json:"id"`
	GivenName  string `yaml:"givenName" json:"givenName"`
	FamilyName string `yaml:"familyName" json:"familyName"`
	Email      string `yaml:"email" json:"email"`
	Credit     int64  `yaml:"credit" json:"credit"`
}

type getUserRes struct {
	ID         int64  `yaml:"id" json:"id"`
	Privilege  string `yaml:"privilege" json:"privilege"`
	GivenName  string `yaml:"givenName" json:"givenName"`
	FamilyName string `yaml:"familyName" json:"familyName"`
	Email      string `yaml:"email" json:"email"`
	Credit     int64  `yaml:"credit" json:"credit"`
}
