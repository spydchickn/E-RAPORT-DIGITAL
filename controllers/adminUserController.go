package controllers

import (
	"e_raport_digital/middlewares"
	"e_raport_digital/models"
	"e_raport_digital/utils"
	"fmt"
	"html/template"
	"net/http"
)

func UsersList(w http.ResponseWriter, r *http.Request) {
    list, err := models.GetAllUsers()
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }
    // enrich users with flags indicating linked siswa/guru records
    type UserView struct {
        ID       int
        Username string
        Role     string
        HasSiswa bool
        HasGuru  bool
    }
    var viewList []UserView
    for _, u := range list {
        hasS := false
        hasG := false
        if s, _ := models.GetSiswaByUserID(u.ID); s != nil {
            hasS = true
        }
        if g, _ := models.GetGuruByUserID(int64(u.ID)); g != nil {
            hasG = true
        }
        viewList = append(viewList, UserView{ID: u.ID, Username: u.Username, Role: u.Role, HasSiswa: hasS, HasGuru: hasG})
    }
    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "base", map[string]interface{}{
        "Title":   "Manajemen Pengguna",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("admin/users_list.html", viewList)),
    })
}

func UserCreateForm(w http.ResponseWriter, r *http.Request) {
    // allow prefilling role via query param ?role=siswa|guru|admin
    role := r.URL.Query().Get("role")
    session, _ := middlewares.Store.Get(r, "session")
    userRole := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "base", map[string]interface{}{
        "Title":   "Tambah Pengguna",
        "Role":    userRole,
        "Content": template.HTML(utils.RenderPartial("admin/user_form.html", map[string]string{"Role": role})),
    })
}

func UserStore(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/users/create", http.StatusSeeOther)
        return
    }
    username := r.FormValue("username")
    password := r.FormValue("password")
    role := r.FormValue("role")
    _, err := models.CreateUser(username, password, role)
    if err != nil {
        http.Error(w, "Unable to save", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/users/list", http.StatusSeeOther)
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    var i int
    fmt.Sscanf(id, "%d", &i)
    models.DeleteUser(i)
    http.Redirect(w, r, "/users/list", http.StatusSeeOther)
}
