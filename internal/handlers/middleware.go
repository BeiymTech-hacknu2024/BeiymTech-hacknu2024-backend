package handlers

import "net/http"

func (userh *UserHandler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := userh.store.Get(r, "session-name")
		_, ok := session.Values["user_id"]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (userh *UserHandler) RequirePermissions(permissons ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// code
			next.ServeHTTP(w, r)
		})
	}
}
