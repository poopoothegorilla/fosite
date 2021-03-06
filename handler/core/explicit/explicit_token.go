package explicit

import (
	"net/http"
	"time"

	"github.com/go-errors/errors"
	"golang.org/x/net/context"
	"github.com/ory-am/fosite"
)

// HandleTokenEndpointRequest implements
// * https://tools.ietf.org/html/rfc6749#section-4.1.3 (everything)
func (c *AuthorizeExplicitGrantTypeHandler) HandleTokenEndpointRequest(ctx context.Context, r *http.Request, request fosite.AccessRequester) error {
	// grant_type REQUIRED.
	// Value MUST be set to "authorization_code".
	if !request.GetGrantTypes().Exact("authorization_code") {
		return errors.New(fosite.ErrUnknownRequest)
	}

	if !request.GetClient().GetGrantTypes().Has("authorization_code") {
		return errors.New(fosite.ErrInvalidGrant)
	}

	// The authorization server MUST verify that the authorization code is valid
	signature, err := c.AuthorizeCodeStrategy.ValidateAuthorizeCode(ctx, request, r.PostForm.Get("code"))
	if err != nil {
		return errors.New(fosite.ErrInvalidRequest)
	}

	authorizeRequest, err := c.AuthorizeCodeGrantStorage.GetAuthorizeCodeSession(ctx, signature, request.GetSession())
	if errors.Is(err, fosite.ErrNotFound) {
		return errors.New(fosite.ErrInvalidRequest)
	} else if err != nil {
		return errors.New(fosite.ErrServerError)
	}

	// Override scopes
	request.SetScopes(authorizeRequest.GetScopes())

	// The authorization server MUST ensure that the authorization code was issued to the authenticated
	// confidential client, or if the client is public, ensure that the
	// code was issued to "client_id" in the request,
	if authorizeRequest.GetClient().GetID() != request.GetClient().GetID() {
		return errors.New(fosite.ErrInvalidRequest)
	}

	// ensure that the "redirect_uri" parameter is present if the
	// "redirect_uri" parameter was included in the initial authorization
	// request as described in Section 4.1.1, and if included ensure that
	// their values are identical.
	forcedRedirectURI := authorizeRequest.GetRequestForm().Get("redirect_uri")
	if forcedRedirectURI != "" && forcedRedirectURI != r.PostForm.Get("redirect_uri") {
		return errors.New(fosite.ErrInvalidRequest)
	}

	// If no lifespan has been set, reset to default lifespan
	if c.AuthCodeLifespan <= 0 {
		c.AuthCodeLifespan = authCodeDefaultLifespan
	}

	// https://tools.ietf.org/html/rfc6819#section-5.1.5.3]
	// A short expiration time for tokens is a means of protection against
	// the following threats: replay, token leak, online guessing
	if authorizeRequest.GetRequestedAt().Add(c.AuthCodeLifespan).Before(time.Now()) {
		return errors.New(fosite.ErrInvalidRequest)
	}

	// Checking of POST client_id skipped, because:
	// If the client type is confidential or the client was issued client
	// credentials (or assigned other authentication requirements), the
	// client MUST authenticate with the authorization server as described
	// in Section 3.2.1.
	request.SetSession(authorizeRequest.GetSession())
	return nil
}

func (c *AuthorizeExplicitGrantTypeHandler) PopulateTokenEndpointResponse(ctx context.Context, req *http.Request, requester fosite.AccessRequester, responder fosite.AccessResponder) error {
	// grant_type REQUIRED.
	// Value MUST be set to "authorization_code".
	if !requester.GetGrantTypes().Exact("authorization_code") {
		return errors.New(fosite.ErrUnknownRequest)
	}

	signature, err := c.AuthorizeCodeStrategy.ValidateAuthorizeCode(ctx, requester, req.PostForm.Get("code"))
	if err != nil {
		return errors.New(fosite.ErrInvalidRequest)
	}

	access, accessSignature, err := c.AccessTokenStrategy.GenerateAccessToken(ctx, requester)
	if err != nil {
		return errors.New(fosite.ErrServerError)
	}

	refresh, refreshSignature, err := c.RefreshTokenStrategy.GenerateRefreshToken(ctx, requester)
	if err != nil {
		return errors.New(fosite.ErrServerError)
	}

	if err := c.AuthorizeCodeGrantStorage.PersistAuthorizeCodeGrantSession(ctx, signature, accessSignature, refreshSignature, requester); err != nil {
		return errors.New(fosite.ErrServerError)
	}

	responder.SetAccessToken(access)
	responder.SetTokenType("bearer")
	responder.SetExpiresIn(c.AccessTokenLifespan / time.Second)
	responder.SetExtra("refresh_token", refresh)
	responder.SetScopes(requester.GetGrantedScopes())
	return nil
}
