package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/TheMangoMen/backend/internal/auth"
	"github.com/TheMangoMen/backend/internal/email"
	"github.com/TheMangoMen/backend/internal/service"
)

func LogIn(a auth.Auth, us service.UserService, emailer email.Emailer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uID := r.PathValue("uID")

		// TODO: Export to Waterloo package
		// TODO: more validation
		if len(uID) > 8 || len(uID) == 0 {
			http.Error(w, "invalid uID", http.StatusBadRequest)
			return
		}

		if err := us.CreateUser(uID); err != nil {
			http.Error(w, "error creating user", http.StatusInternalServerError)
			return
		}

		admin, err := us.GetIsAdmin(uID)
		if err != nil {
			http.Error(w, "error checking user admin status", http.StatusInternalServerError)
			return
		}

		token, err := a.NewToken(auth.User{UID: uID, Admin: admin})
		if err != nil {
			http.Error(w, "error signing token", http.StatusInternalServerError)
			return
		}

		encoded := base64.URLEncoding.EncodeToString([]byte(token))
		err = emailer.Send(
			fmt.Sprintf("%s@uwaterloo.ca", uID),
			"WatRank Login Link",
			fmt.Sprintf("<p>Here is your <a href=\"http://watrank.com/callback?code=%s\">login link.</a></p>", encoded),
		)
		if err != nil {
			http.Error(w, "error sending email", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}
