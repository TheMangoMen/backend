package handler

import (
	"fmt"
	"net/http"

	"github.com/TheMangoMen/backend/internal/auth"
	"github.com/TheMangoMen/backend/internal/email"
)

func LogIn(a auth.Auth, ec email.EmailClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uID := r.PathValue("uID")
		token, err := a.NewToken(uID)
		if err != nil {
			http.Error(w, "error signing token", http.StatusInternalServerError)
			return
		}

		// w.Header().Add("Set-Cookie", fmt.Sprintf("__Host-Authorization=Bearer %s; path=/; Secure; SameSite=strict; HttpOnly", token)) // Maybe when live?
		// w.Header().Add("Set-Cookie", fmt.Sprintf("Authorization=Bearer %s; path=/; Secure; SameSite=strict; HttpOnly", token))

		err = ec.Send(
			fmt.Sprintf("%s@uwaterloo.ca", uID),
			"hi there",
			fmt.Sprintf("<p>Here is your log in link: %s</p>", token),
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
