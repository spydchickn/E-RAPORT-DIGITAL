package controllers

import (
    "database/sql"
    "html/template"
    "log"
    "net/http"
    "strconv"
    "strings"

    "e_raport_digital/config"
    "e_raport_digital/middlewares"
    "e_raport_digital/models"
    "e_raport_digital/utils"
)

func PresensiList(w http.ResponseWriter, r *http.Request) {
    rows, err := config.DB.Query(`
        SELECT p.id_presensi, p.id_siswa, p.id_mapel, p.tanggal, p.status, p.catatan,
               s.nama AS nama_siswa, m.nama_mapel
        FROM presensi p
        JOIN siswa s ON p.id_siswa = s.id_siswa
        JOIN mapel m ON p.id_mapel = m.id_mapel
        ORDER BY p.tanggal DESC
    `)
    if err != nil {
        log.Println("Error querying presensi:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var presensi []struct {
        IDPresensi int
        IDSiswa    int
        IDMapel    int
        Tanggal    string
        Status     string
        Catatan    string
        NamaSiswa  string
        NamaMapel  string
    }

    for rows.Next() {
        var p struct {
            IDPresensi int
            IDSiswa    int
            IDMapel    int
            Tanggal    string
            Status     string
            Catatan    string
            NamaSiswa  string
            NamaMapel  string
        }
        err := rows.Scan(&p.IDPresensi, &p.IDSiswa, &p.IDMapel, &p.Tanggal, &p.Status, &p.Catatan, &p.NamaSiswa, &p.NamaMapel)
        if err != nil {
            log.Println("Error scanning presensi:", err)
            continue
        }
        presensi = append(presensi, p)
    }

    siswa, _ := models.GetAllSiswa()
    mapel, _ := models.GetAllMapel()
    kelas, _ := models.GetAllKelas()

    // Use global templates and base layout
    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   "Daftar Presensi",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("admin/presensi_list.html", map[string]interface{}{
            "Presensi": presensi,
            "Siswa":    siswa,
            "Mapel":    mapel,
            "Kelas":    kelas,
        })),
    })
}

func PresensiCreateForm(w http.ResponseWriter, r *http.Request) {
    siswa, _ := models.GetAllSiswa()
    mapel, _ := models.GetAllMapel()

    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   "Tambah Presensi",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("admin/presensi_create.html", map[string]interface{}{
            "Siswa": siswa,
            "Mapel": mapel,
        })),
    })
}

func PresensiStore(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    idSiswa, _ := strconv.Atoi(r.FormValue("id_siswa"))
    idMapel, _ := strconv.Atoi(r.FormValue("id_mapel"))
    tanggal := r.FormValue("tanggal")
    status := strings.ToLower(r.FormValue("status"))
    catatan := r.FormValue("catatan")

    err := models.CreatePresensi(idSiswa, idMapel, 0, tanggal, status, catatan)
    if err != nil {
        log.Println("Error creating presensi:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Log action
    userID := utils.GetUserIDFromSession(r)
    models.LogAction(userID, "create", "presensi", sql.NullInt64{Int64: int64(idSiswa), Valid: true}, sql.NullString{String: "Presensi created", Valid: true})

    http.Redirect(w, r, "/admin/presensi/list", http.StatusSeeOther)
}

func PresensiEditForm(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, _ := strconv.Atoi(idStr)

    presensi, err := models.GetPresensiByID(id)
    if err != nil || presensi == nil {
        http.NotFound(w, r)
        return
    }

    siswa, _ := models.GetAllSiswa()
    mapel, _ := models.GetAllMapel()

    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   "Edit Presensi",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("admin/presensi_edit.html", map[string]interface{}{
            "Presensi": presensi,
            "Siswa":    siswa,
            "Mapel":    mapel,
        })),
    })
}

func PresensiUpdate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    id, _ := strconv.Atoi(r.FormValue("id_presensi"))
    idSiswa, _ := strconv.Atoi(r.FormValue("id_siswa"))
    idMapel, _ := strconv.Atoi(r.FormValue("id_mapel"))
    tanggal := r.FormValue("tanggal")
    status := strings.ToLower(r.FormValue("status"))
    catatan := r.FormValue("catatan")

    err := models.UpdatePresensi(id, idSiswa, idMapel, 0, tanggal, status, catatan)
    if err != nil {
        log.Println("Error updating presensi:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Log action
    userID := utils.GetUserIDFromSession(r)
    models.LogAction(userID, "update", "presensi", sql.NullInt64{Int64: int64(id), Valid: true}, sql.NullString{String: "Presensi updated", Valid: true})

    http.Redirect(w, r, "/admin/presensi/list", http.StatusSeeOther)
}

func PresensiDelete(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, _ := strconv.Atoi(idStr)

    err := models.DeletePresensi(id)
    if err != nil {
        log.Println("Error deleting presensi:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Log action
    userID := utils.GetUserIDFromSession(r)
    models.LogAction(userID, "delete", "presensi", sql.NullInt64{Int64: int64(id), Valid: true}, sql.NullString{String: "Presensi deleted", Valid: true})

    http.Redirect(w, r, "/admin/presensi/list", http.StatusSeeOther)
}

func GuruPresensiInputForm(w http.ResponseWriter, r *http.Request) {
    siswa, _ := models.GetAllSiswa()
    mapel, _ := models.GetAllMapel()

    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   "Input Presensi",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("guru/presensi_input.html", map[string]interface{}{
            "Siswa": siswa,
            "Mapel": mapel,
        })),
    })
}

func GuruPresensiStore(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    idSiswa, _ := strconv.Atoi(r.FormValue("id_siswa"))
    idMapel, _ := strconv.Atoi(r.FormValue("id_mapel"))
    tanggal := r.FormValue("tanggal")
    status := r.FormValue("status")
    catatan := r.FormValue("catatan")

    err := models.CreatePresensi(idSiswa, idMapel, 0, tanggal, status, catatan)
    if err != nil {
        log.Println("Error creating presensi:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Log action
    userID := utils.GetUserIDFromSession(r)
    models.LogAction(userID, "create", "presensi", sql.NullInt64{Int64: int64(idSiswa), Valid: true}, sql.NullString{String: "Presensi created by guru", Valid: true})

    http.Redirect(w, r, "/guru/dashboard", http.StatusSeeOther)
}


