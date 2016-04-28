package strategy

import (
	"errors"

	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/enigma/jwt"
	"golang.org/x/net/context"
)

type JWTSessionContainer interface {
	// GetTokenClaims returns the claims
	GetTokenClaims() *jwt.JWTClaims

	// GetTokenHeader returns the header
	GetTokenHeader() *jwt.Header
}

// JWTSession Container for the JWT session
type JWTSession struct {
	JWTClaims *jwt.JWTClaims
	JWTHeader *jwt.Header
}

func (j *JWTSession) GetTokenClaims() *jwt.JWTClaims {
	return j.JWTClaims
}

func (j *JWTSession) GetTokenHeader() *jwt.Header {
	return j.JWTHeader
}

// JWTStrategy : Strategy container
type JWTStrategy struct {
	Enigma *jwt.Enigma
}

func (h *JWTStrategy) GenerateAccessToken(_ context.Context, requester fosite.Requester) (token string, signature string, err error) {
	return h.generate(requester)
}

func (h *JWTStrategy) ValidateAccessToken(_ context.Context, _ fosite.Requester, token string) (signature string, err error) {
	return h.validate(token)
}

func (h *JWTStrategy) GenerateRefreshToken(_ context.Context, requester fosite.Requester) (token string, signature string, err error) {
	return h.generate(requester)
}

func (h *JWTStrategy) ValidateRefreshToken(_ context.Context, _ fosite.Requester, token string) (signature string, err error) {
	return h.validate(token)
}

func (h *JWTStrategy) GenerateAuthorizeCode(_ context.Context, requester fosite.Requester) (token string, signature string, err error) {
	return h.generate(requester)
}

func (h *JWTStrategy) ValidateAuthorizeCode(_ context.Context, requester fosite.Requester, token string) (signature string, err error) {
	return h.validate(token)
}

func (h *JWTStrategy) validate(token string) (string, error) {
	t, err := h.Enigma.Decode(token)
	if err != nil {
		return "", err
	}

	claims := jwt.JWTClaimsFromMap(t.Claims)
	if claims.IsNotYetValid() || claims.IsExpired() {
		return "", errors.New("Token claims did not validate")
	}

	return h.Enigma.GetSignature(token)
}

func (h *JWTStrategy) generate(requester fosite.Requester) (string, string, error) {
	if jwtSession, ok := requester.GetSession().(JWTSessionContainer); ok {
		if jwtSession.GetTokenClaims() != nil {
			return h.Enigma.Generate(jwtSession.GetTokenClaims(), jwtSession.GetTokenHeader())
		}
		return "", "", errors.New("GetTokenClaims() must not be nil")
	}
	return "", "", errors.New("Session must be of type JWTSession")

}
