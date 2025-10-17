package controllers

import (
	"e_raport_digital/middlewares"
	"e_raport_digital/models"
	"e_raport_digital/utils"
	"html/template"
	"net/http"
	"strconv"
)

func InputNilaiForm(w http.ResponseWriter, r *http.Request) {
    session, _ := middlewares.Store.Get(r, "session")
    uidv := session.Values["user_id"]
    if uidv == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    var userID int64
    switch v := uidv.(type) {
    case int:
        userID = int64(v)
    case int64:
        userID = v
    case float64:
        userID = int64(v)
    default:
        http.Error(w, "Invalid session", http.StatusBadRequest)
        return
    }

    guru, err := models.GetGuruByUserID(userID)
    if err != nil || guru == nil {
        http.Error(w, "Guru profile not found", http.StatusInternalServerError)
        return
    }

    assignedMapel, err := models.GetMapelByGuruID(guru.IDGuru)
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    allSiswa, err := models.GetAllSiswa()
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title": "Input Nilai",
        "Role":  role,
        "Content": template.HTML(utils.RenderPartial("guru/input_nilai.html", map[string]interface{}{
            "Siswa":  allSiswa,
            "Mapels": assignedMapel,
        })),
    })
}

// Store nilai from guru input
func StoreNilai(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/guru/input-nilai", http.StatusSeeOther)
        return
    }
    idSiswaStr := r.FormValue("id_siswa")
    idMapelStr := r.FormValue("id_mapel")
    nilaiStr := r.FormValue("nilai")
    komentar := r.FormValue("komentar")
    idSiswa, _ := strconv.Atoi(idSiswaStr)
    idMapel, _ := strconv.Atoi(idMapelStr)
    nilai, _ := strconv.Atoi(nilaiStr)
    // basic validation
    if idSiswa <= 0 || idMapel <= 0 {
        http.Error(w, "Invalid IDs", http.StatusBadRequest)
        return
    }
    _ = models.CreateNilai(idSiswa, idMapel, nilai, komentar)
    http.Redirect(w, r, "/guru/dashboard", http.StatusSeeOther)
}
