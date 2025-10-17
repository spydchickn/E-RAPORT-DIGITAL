package middlewares

import (
    "net/http"

    "github.com/gorilla/sessions"
)

// single exported Store used by controllers and middleware
var Store = sessions.NewCookieStore([]byte("super-secret-key"))

// AuthMiddleware ensures the user is logged in and has the required role
func AuthMiddleware(role string, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, _ := Store.Get(r, "session")
        v := session.Values["role"]
        if v == nil {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        userRole, ok := v.(string)
        if !ok {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        // If role is empty, allow any authenticated role
        if role != "" && userRole != role {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        next(w, r)
    }
}
