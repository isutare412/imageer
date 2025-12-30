package domain

// Passport represents the authentication credentials of a client.
type Passport interface {
	isPassport()
}

type UserTokenPassport struct {
	Payload UserTokenPayload
}

func NewUserTokenPassport(payload UserTokenPayload) UserTokenPassport {
	return UserTokenPassport{
		Payload: payload,
	}
}

func (UserTokenPassport) isPassport() {}

type ServiceAccountPassport struct {
	ServiceAccount ServiceAccount
}

func NewServiceAccountPassport(account ServiceAccount) ServiceAccountPassport {
	return ServiceAccountPassport{
		ServiceAccount: account,
	}
}

func (ServiceAccountPassport) isPassport() {}
