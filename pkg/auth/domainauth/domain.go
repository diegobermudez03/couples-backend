package domainauth

type AuthService interface {
	RegisterUser(email string, password string, device string, os string) (refreshToken string, accessToken string, err error)
}