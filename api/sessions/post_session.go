package sessions

import (
	"net/http"

	"github.com/keratin/authn-server/api"
	"github.com/keratin/authn-server/services"
)

func postSession(app *api.App) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Check the password
		account, errors := services.CredentialsVerifier(
			app.AccountStore,
			app.Config,
			req.FormValue("username"),
			req.FormValue("password"),
		)
		if errors != nil {
			api.WriteErrors(w, errors)
			return
		}

		err := api.RevokeSession(app.RefreshTokenStore, app.Config, req)
		if err != nil {
			// TODO: alert but continue
		}

		sessionToken, identityToken, err := api.NewSession(app.RefreshTokenStore, app.Config, account.Id)
		if err != nil {
			panic(err)
		}

		// Return the signed session in a cookie
		api.SetSession(app.Config, w, sessionToken)

		// Return the signed identity token in the body
		api.WriteData(w, http.StatusCreated, map[string]string{
			"id_token": identityToken,
		})
	}
}
