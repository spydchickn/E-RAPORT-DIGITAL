package controllers

import (
	"e_raport_digital/middlewares"
	"e_raport_digital/models"
	"e_raport_digital/utils"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

func SiswaRaport(w http.ResponseWriter, r *http.Request) {
    // get logged-in user id from session
    session, _ := middlewares.Store.Get(r, "session")
    uidv := session.Values["user_id"]
    if uidv == nil {
        utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
            "Title": "Hasil Raport Saya",
            "Content": "<div class='alert alert-danger'>Anda belum login sebagai siswa.</div>",
        })
        return
    }
    // accept int / int64 / float64
    var userID int
    switch v := uidv.(type) {
    case int:
        userID = v
    case int64:
        userID = int(v)
    case float64:
        userID = int(v)
    default:
        userID = 0
    }
    if userID == 0 {
        utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
            "Title": "Hasil Raport Saya",
            "Content": "<div class='alert alert-danger'>Anda belum login sebagai siswa.</div>",
        })
        return
    }
    raport, err := models.GetRaportByUserID(userID)
    if err != nil || raport == nil {
        utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
            "Title": "Hasil Raport Saya",
            "Content": "<div class='alert alert-warning'>Tidak menemukan data raport untuk akun ini.</div>",
        })
        return
    }
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title": "Hasil Raport Saya",
        "Content": template.HTML(utils.RenderPartial("siswa/raport.html", raport)),
    })
}

// SiswaRaportPDF generates and serves PDF for logged-in student
func SiswaRaportPDF(w http.ResponseWriter, r *http.Request) {
    session, _ := middlewares.Store.Get(r, "session")
    uidv := session.Values["user_id"]
    if uidv == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    var userID int
    switch v := uidv.(type) {
    case int:
        userID = v
    case int64:
        userID = int(v)
    case float64:
        userID = int(v)
    default:
        http.Error(w, "Invalid session", http.StatusBadRequest)
        return
    }
    raport, err := models.GetRaportByUserID(userID)
    if err != nil || raport == nil {
        http.Error(w, "No raport data", http.StatusNotFound)
        return
    }
    filename := fmt.Sprintf("raport_%s.pdf", raport.NamaSiswa)
    err = utils.ExportRaportPDF(raport, filename)
    if err != nil {
        http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", "attachment; filename="+filename)
    http.ServeFile(w, r, filepath.Join("exports", filename))
}

func SiswaList(w http.ResponseWriter, r *http.Request) {
    list, err := models.GetAllSiswa()
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }
    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title": "Daftar Siswa",
        "Role":  role,
        "Content": template.HTML(utils.RenderPartial("siswa/list.html", list)),
    })
}

func SiswaCreateForm(w http.ResponseWriter, r *http.Request) {
    // allow pre-linking to an existing user via ?user_id=123
    userID := r.URL.Query().Get("user_id")
    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title": "Tambah Siswa",
        "Role":  role,
        "Content": template.HTML(utils.RenderPartial("siswa/create.html", map[string]string{"UserID": userID})),
    })
}

func SiswaStore(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/siswa/create", http.StatusSeeOther)
        return
    }
    nis := r.FormValue("nis")
    nama := r.FormValue("nama")
    alamat := r.FormValue("alamat")
    foto := r.FormValue("foto")
    ortu_nama := r.FormValue("ortu_nama")
    ortu_email := r.FormValue("ortu_email")
    ortu_telp := r.FormValue("ortu_telp")
    // if a user_id was provided (admin pre-linked), use it; otherwise create a user for siswa
    userIDStr := r.FormValue("user_id")
    var idUser int64
    if userIDStr != "" {
        // parse int
        var parsed int
        fmt.Sscanf(userIDStr, "%d", &parsed)
        idUser = int64(parsed)
    } else {
        id, err := models.CreateUser(nis, nis, "siswa")
        if err != nil {
            http.Error(w, "Unable to create user for siswa", http.StatusInternalServerError)
            return
        }
        idUser = id
    }
    // if admin provided user_id, ensure no siswa already linked
    if userIDStr != "" {
        if s, _ := models.GetSiswaByUserID(int(idUser)); s != nil {
            http.Error(w, "User already linked to a siswa", http.StatusBadRequest)
            return
        }
    }
    _, err := models.CreateSiswa(nis, nama, alamat, foto, ortu_nama, ortu_email, ortu_telp, idUser)
    if err != nil {
        // try to rollback user creation
        if userIDStr == "" {
            models.DeleteUser(int(idUser))
        }
        http.Error(w, "Unable to save siswa", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/siswa/list", http.StatusSeeOther)
}

func SiswaEditForm(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    // Fetch siswa by id (not implemented in model, so just pass id for now)
    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title": "Edit Siswa",
        "Role":  role,
        "Content": template.HTML(utils.RenderPartial("siswa/create.html", map[string]string{"ID": id})),
    })
}

func SiswaUpdate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/siswa/list", http.StatusSeeOther)
        return
    }
    id := r.FormValue("id")
    nis := r.FormValue("nis")
    nama := r.FormValue("nama")
    alamat := r.FormValue("alamat")
    // Convert id to int
    // (error handling omitted for brevity)
    models.UpdateSiswa(stringToInt(id), nis, nama, alamat)
    http.Redirect(w, r, "/siswa/list", http.StatusSeeOther)
}

func SiswaDelete(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    models.DeleteSiswa(stringToInt(id))
    http.Redirect(w, r, "/siswa/list", http.StatusSeeOther)
}

func stringToInt(s string) int {
    var i int
    fmt.Sscanf(s, "%d", &i)
    return i
}
