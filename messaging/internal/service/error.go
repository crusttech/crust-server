package service

type (
	serviceError string
)

func (e serviceError) Error() string {
	return string(e)
}

const (
	ErrNoPermissions serviceError = "You don't have permissions for this operation"
)
