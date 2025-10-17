package controllers

import (
	"e_raport_digital/middlewares"
	"e_raport_digital/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// ShowLogin serves the standalone login page
func ShowLogin(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "views/auth/login.html")
}

// Login processes POST /login
func Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    username := r.FormValue("username")
    password := r.FormValue("password")

    user, err := models.GetUserByUsername(username)
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }
    if user == nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    // First try comparing as a bcrypt hash
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        // If bcrypt compare failed, it might be because the stored password is plaintext
        // In that case, accept plaintext match and migrate to a bcrypt hash
        if user.Password == password {
            if hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err == nil {
                // best-effort update; ignore error but log if needed
                _ = models.SetUserPasswordHash(user.ID, string(hashed))
            }
        } else {
            // not a match, redirect back to login
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
    }

    session, _ := middlewares.Store.Get(r, "session")
    session.Values["username"] = user.Username
    session.Values["role"] = user.Role
    // store numeric user ID so other handlers can map to siswa/guru records
    session.Values["user_id"] = user.ID
    session.Save(r, w)

    switch user.Role {
    case "admin":
        http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
    case "guru":
        http.Redirect(w, r, "/guru/dashboard", http.StatusSeeOther)
    case "siswa":
        http.Redirect(w, r, "/siswa/dashboard", http.StatusSeeOther)
    default:
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func Logout(w http.ResponseWriter, r *http.Request) {
    session, _ := middlewares.Store.Get(r, "session")
    session.Options.MaxAge = -1
    session.Save(r, w)
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}
