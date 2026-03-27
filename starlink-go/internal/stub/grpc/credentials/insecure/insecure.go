package insecure

type Credentials struct{}

func NewCredentials() Credentials {
	return Credentials{}
}
