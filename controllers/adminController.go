package controllers

import (
	"e_raport_digital/middlewares"
	"e_raport_digital/models"
	"e_raport_digital/utils"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// Guru CRUD
func GuruList(w http.ResponseWriter, r *http.Request) {
    list, err := models.GetAllGuru()
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }
    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   "Daftar Guru",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("admin/guru_list.html", list)),
    })
}

func GuruCreateForm(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   "Tambah Guru",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("admin/guru_form.html", map[string]string{"UserID": userID})),
    })
}

func GuruStore(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/guru/create", http.StatusSeeOther)
        return
    }
    nama := r.FormValue("nama")
    nip := r.FormValue("nip")
    alamat := r.FormValue("alamat")
    foto := r.FormValue("foto")
    userIDStr := r.FormValue("user_id")
    var err error
    if userIDStr != "" {
        var parsed int
        fmt.Sscanf(userIDStr, "%d", &parsed)
        // prevent duplicate guru for same user
        if g, _ := models.GetGuruByUserID(int64(parsed)); g != nil {
            http.Error(w, "User already linked to a guru", http.StatusBadRequest)
            return
        }
        _, err = models.CreateGuruWithUser(nama, nip, alamat, foto, int64(parsed))
    } else {
        err = models.CreateGuru(nama, nip, alamat, foto)
    }
    if err != nil {
        http.Error(w, "Unable to save", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/guru/list", http.StatusSeeOther)
}

func GuruDelete(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    var i int
    fmt.Sscanf(id, "%d", &i)
    models.DeleteGuru(i)
    http.Redirect(w, r, "/guru/list", http.StatusSeeOther)
}

// GuruAssignMapelForm shows form to assign/unassign mapel to a guru
func GuruAssignMapelForm(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    guruIDStr := vars["id"]
    var guruID int
    fmt.Sscanf(guruIDStr, "%d", &guruID)

    guru, err := models.GetGuruByID(guruID)
    if err != nil || guru == nil {
        http.Error(w, "Guru not found", http.StatusNotFound)
        return
    }

    allMapel, err := models.GetAllMapel()
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    assignedMapel, err := models.GetMapelByGuruID(guruID)
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    // Prepare data for template: all mapel with assigned flag
    type MapelAssign struct {
        Mapel    models.Mapel
        Assigned bool
    }
    var mapelAssigns []MapelAssign
    for _, m := range allMapel {
        assigned := false
        for _, am := range assignedMapel {
            if am.IDMapel == m.IDMapel {
                assigned = true
                break
            }
        }
        mapelAssigns = append(mapelAssigns, MapelAssign{Mapel: m, Assigned: assigned})
    }

    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]

    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   fmt.Sprintf("Assign Mapel for %s", guru.Nama),
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("admin/guru_assign.html", map[string]interface{}{
            "Guru":    guru,
            "Mapels":  mapelAssigns,
            "GuruID":  guruID,
        })),
    })
}

// StoreGuruMapel handles POST to assign/unassign mapel
func StoreGuruMapel(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/guru/list", http.StatusSeeOther)
        return
    }

    guruIDStr := r.FormValue("guru_id")
    var guruID int
    fmt.Sscanf(guruIDStr, "%d", &guruID)

    // Get all mapel IDs
    allMapel, err := models.GetAllMapel()
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    // First, remove all assignments for this guru
    for _, m := range allMapel {
        if err := models.RemoveGuruMapel(guruID, m.IDMapel); err != nil {
            // Ignore errors for non-existent assignments
        }
    }

    // Then, assign selected ones
    for _, m := range allMapel {
        if r.FormValue(fmt.Sprintf("mapel_%d", m.IDMapel)) == "on" {
            if err := models.AssignGuruMapel(guruID, m.IDMapel); err != nil {
                http.Error(w, "Unable to save assignments", http.StatusInternalServerError)
                return
            }
        }
    }

    http.Redirect(w, r, fmt.Sprintf("/guru/assign/%d", guruID), http.StatusSeeOther)
}

// Mapel CRUD
func MapelList(w http.ResponseWriter, r *http.Request) {
    list, err := models.GetAllMapel()
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }
    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   "Daftar Mapel",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("admin/mapel_list.html", list)),
    })
}

func MapelCreateForm(w http.ResponseWriter, r *http.Request) {
    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   "Tambah Mapel",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("admin/mapel_form.html", nil)),
    })
}

func MapelStore(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/mapel/create", http.StatusSeeOther)
        return
    }
    nama := r.FormValue("nama")
    kkm := r.FormValue("kkm")
    err := models.CreateMapel(nama, kkm)
    if err != nil {
        http.Error(w, "Unable to save", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/mapel/list", http.StatusSeeOther)
}

func MapelDelete(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    var i int
    fmt.Sscanf(id, "%d", &i)
    models.DeleteMapel(i)
    http.Redirect(w, r, "/mapel/list", http.StatusSeeOther)
}

// Kelas CRUD
func KelasList(w http.ResponseWriter, r *http.Request) {
    list, err := models.GetAllKelas()
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }
    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   "Daftar Kelas",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("admin/kelas_list.html", list)),
    })
}

func KelasCreateForm(w http.ResponseWriter, r *http.Request) {
    session, _ := middlewares.Store.Get(r, "session")
    role := session.Values["role"]
    utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
        "Title":   "Tambah Kelas",
        "Role":    role,
        "Content": template.HTML(utils.RenderPartial("admin/kelas_form.html", nil)),
    })
}

func KelasStore(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/kelas/create", http.StatusSeeOther)
        return
    }
    nama := r.FormValue("nama")
    err := models.CreateKelas(nama)
    if err != nil {
        http.Error(w, "Unable to save", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/kelas/list", http.StatusSeeOther)
}

func KelasDelete(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    var i int
    fmt.Sscanf(id, "%d", &i)
    models.DeleteKelas(i)
    http.Redirect(w, r, "/kelas/list", http.StatusSeeOther)
}
