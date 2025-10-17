package controllers

import (
    "database/sql"
    "html/template"
    "log"
    "net/http"
    "strconv"

    "e_raport_digital/middlewares"
    "e_raport_digital/models"
    "e_raport_digital/utils"
)

func NotifikasiList(w http.ResponseWriter, r *http.Request) {
    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"].(string)
    userID := session.Values["user_id"].(int)

    var notifs []models.Notifikasi
    var err error

    if role == "admin" {
        notifs, err = models.GetAllNotifikasi()
    } else {
        notifs, err = models.GetUnreadByUser(userID)
    }

    if err != nil {
        log.Println("Error getting notifikasi:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    data := map[string]interface{}{
        "Notifikasi": notifs,
        "Role":       role,
    }

    tmpl := template.Must(template.ParseFiles("views/layouts/base.html", "views/notifikasi_list.html"))
    tmpl.Execute(w, data)
}

func NotifikasiMarkRead(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    id, _ := strconv.Atoi(r.FormValue("id_notif"))

    err := models.MarkRead(id)
    if err != nil {
        log.Println("Error marking read:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/notifikasi/list", http.StatusSeeOther)
}

func NotifikasiBroadcastForm(w http.ResponseWriter, r *http.Request) {
    siswa, _ := models.GetAllSiswa()

    data := map[string]interface{}{
        "Siswa": siswa,
    }

    tmpl := template.Must(template.ParseFiles("views/layouts/base.html", "views/notifikasi_broadcast.html"))
    tmpl.Execute(w, data)
}

func NotifikasiBroadcastStore(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    target := r.FormValue("target")
    pesan := r.FormValue("pesan")

    var err error
    if target == "role" {
        role := r.FormValue("role")
        err = models.CreateNotif(sql.NullInt64{}, sql.NullString{String: role, Valid: true}, pesan)
    } else {
        idSiswa, _ := strconv.Atoi(r.FormValue("id_siswa"))
        err = models.CreateNotif(sql.NullInt64{Int64: int64(idSiswa), Valid: true}, sql.NullString{}, pesan)
    }

    if err != nil {
        log.Println("Error broadcasting notifikasi:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Log action
    userID := utils.GetUserIDFromSession(r)
    models.LogAction(userID, "broadcast", "notifikasi", sql.NullInt64{}, sql.NullString{String: "Notifikasi broadcasted", Valid: true})

    http.Redirect(w, r, "/notifikasi/list", http.StatusSeeOther)
}
