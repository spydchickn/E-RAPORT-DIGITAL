package controllers

import (
    "database/sql"
    "html/template"
    "log"
    "net/http"
    "strconv"

    "e_raport_digital/config"
    "e_raport_digital/middlewares"
    "e_raport_digital/models"
    "e_raport_digital/utils"
)

func VerifikasiList(w http.ResponseWriter, r *http.Request) {
    rows, err := config.DB.Query(`
        SELECT n.id_nilai, n.id_siswa, n.id_mapel, n.nilai, n.status,
               s.nama AS nama_siswa, m.nama_mapel
        FROM nilai n
        JOIN siswa s ON n.id_siswa = s.id_siswa
        JOIN mapel m ON n.id_mapel = m.id_mapel
        WHERE n.status = 'draft'
        ORDER BY n.id_nilai DESC
    `)
    if err != nil {
        log.Println("Error querying nilai pending:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var nilai []struct {
        IDNilai    int
        IDSiswa    int
        IDMapel    int
        Nilai      int
        Status     string
        NamaSiswa  string
        NamaMapel  string
    }

    for rows.Next() {
        var n struct {
            IDNilai    int
            IDSiswa    int
            IDMapel    int
            Nilai      int
            Status     string
            NamaSiswa  string
            NamaMapel  string
        }
        err := rows.Scan(&n.IDNilai, &n.IDSiswa, &n.IDMapel, &n.Nilai, &n.Status, &n.NamaSiswa, &n.NamaMapel)
        if err != nil {
            log.Println("Error scanning nilai:", err)
            continue
        }
        nilai = append(nilai, n)
    }

    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   "Verifikasi Nilai",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("admin/verifikasi_list.html", map[string]interface{}{
            "Nilai": nilai,
        })),
    })
}

func VerifikasiStore(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    idNilai, _ := strconv.Atoi(r.FormValue("id_nilai"))

    userID := utils.GetUserIDFromSession(r)
    err := models.UpdateNilaiStatus(idNilai, "verified")
    if err != nil {
        log.Println("Error verifying nilai:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Create notifikasi to siswa
    nilai, err := models.GetNilaiByID(idNilai)
    if err != nil || nilai == nil {
        log.Println("Error getting nilai for notification:", err)
        http.Redirect(w, r, "/admin/verifikasi/list", http.StatusSeeOther)
        return
    }
    siswa, err := models.GetSiswaByID(nilai.IDSiswa)
    if err != nil || siswa == nil {
        log.Println("Error getting siswa for notification:", err)
        http.Redirect(w, r, "/admin/verifikasi/list", http.StatusSeeOther)
        return
    }
    err = models.CreateNotif(siswa.IDUser, sql.NullString{String: "siswa", Valid: true}, "Nilai Anda telah diverifikasi.")
    if err != nil {
        log.Println("Error creating notification:", err)
    }

    // Log action
    err = models.LogAction(userID, "verify", "nilai", sql.NullInt64{Int64: int64(idNilai), Valid: true}, sql.NullString{String: "Nilai verified", Valid: true})
    if err != nil {
        log.Println("Error logging action:", err)
    }

    http.Redirect(w, r, "/admin/verifikasi/list", http.StatusSeeOther)
}
