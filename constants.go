package skolengo

type Constants struct {
	OIDCClientID     string
	OIDCClientSecret string
	RedirectURI      string
}

const OIDC_CLIENT_ID = "SkoApp.Prod.0d349217-9a4e-41ec-9af9-df9e69e09494"
const OIDC_CLIENT_SECRET = "7cb4d9a8-2580-4041-9ae8-d5803869183f"
const REDIRECT_URI = "skoapp-prod://sign-in-callback"

var SkolenGoConstants = Constants{
	OIDCClientID:     OIDC_CLIENT_ID,
	OIDCClientSecret: OIDC_CLIENT_SECRET,
	RedirectURI:      REDIRECT_URI,
}
