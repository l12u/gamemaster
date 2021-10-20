// Package auth will later be used to connect to another
// authentication service to make sure the token we get is valid and
// belongs to an authenticated client/other service.
package auth

func IsAuthenticated(token string) bool {
	return true
}
