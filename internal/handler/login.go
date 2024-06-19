package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/TheMangoMen/backend/internal/auth"
	"github.com/TheMangoMen/backend/internal/email"
)

func LogIn(a auth.Auth, emailer email.Emailer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uID := r.PathValue("uID")
		// TODO: uID validation

		token, err := a.NewToken(uID)
		if err != nil {
			http.Error(w, "error signing token", http.StatusInternalServerError)
			return
		}

		// w.Header().Add("Set-Cookie", fmt.Sprintf("__Host-Authorization=Bearer %s; path=/; Secure; SameSite=strict; HttpOnly", token)) // Maybe when live?
		// w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=Bearer %s; path=/; Secure; SameSite=strict; HttpOnly", token))

		encoded := base64.URLEncoding.EncodeToString([]byte(token))
		err = emailer.Send(
			fmt.Sprintf("%s@uwaterloo.ca", uID),
			"WatRank Login Link",
			fmt.Sprintf("<p>Here is your <a href=\"http://localhost:3000/callback?code=%s\">login link.</a></p>", encoded),
		)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error sending email", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}
