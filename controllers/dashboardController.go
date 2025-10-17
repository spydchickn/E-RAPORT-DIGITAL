package controllers

import (
	"e_raport_digital/middlewares"
	"e_raport_digital/models"
	"e_raport_digital/utils"
	"html/template"
	"net/http"
)

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
    // gather stats
    sCount, _ := models.CountSiswa()
    gCount, _ := models.CountGuru()
    kCount, _ := models.CountKelas()
    mCount, _ := models.CountMapel()
    uCount, _ := models.CountUsers()

    recent, _ := models.GetAllSiswa()

    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]

    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title": "Dashboard Admin",
        "Role":  role,
        "Content": template.HTML(utils.RenderPartial("admin/dashboard.html", map[string]interface{}{
            "Siswa":  sCount,
            "Guru":   gCount,
            "Kelas":  kCount,
            "Mapel":  mCount,
            "Users":  uCount,
            "Recent": recent,
        })),
    })
}

func GuruDashboard(w http.ResponseWriter, r *http.Request) {
    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]

    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   "Dashboard Guru",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("guru/dashboard.html", nil)),
    })
}

func SiswaDashboard(w http.ResponseWriter, r *http.Request) {
    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]

    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   "Dashboard Siswa",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("siswa/dashboard.html", nil)),
    })
}
