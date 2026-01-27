package domain

// Identity represents the authentication credentials of a client.
type Identity interface {
	isIdentity()
}

type UserTokenIdentity struct {
	Payload UserTokenPayload
}

func NewUserTokenIdentity(payload UserTokenPayload) UserTokenIdentity {
	return UserTokenIdentity{
		Payload: payload,
	}
}

func (UserTokenIdentity) isIdentity() {}

type ServiceAccountIdentity struct {
	ServiceAccount ServiceAccount
}

func NewServiceAccountIdentity(account ServiceAccount) ServiceAccountIdentity {
	return ServiceAccountIdentity{
		ServiceAccount: account,
	}
}

func (ServiceAccountIdentity) isIdentity() {}
