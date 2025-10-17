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

    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   "Notifikasi",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("notifikasi_list.html", map[string]interface{}{
            "Notifikasi": notifs,
            "Role":       role,
        })),
    })
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
    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   "Kirim Notifikasi",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("notifikasi_broadcast.html", map[string]interface{}{"Siswa": siswa})),
    })
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
        // Translate siswa id to linked users.id (id_user)
        idSiswa, _ := strconv.Atoi(r.FormValue("id_siswa"))
        siswa, _ := models.GetSiswaByID(idSiswa)
        if siswa != nil && siswa.IDUser.Valid {
            err = models.CreateNotif(sql.NullInt64{Int64: siswa.IDUser.Int64, Valid: true}, sql.NullString{}, pesan)
        } else {
            // Fallback: no linked user, treat as error
            err = sql.ErrNoRows
        }
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
