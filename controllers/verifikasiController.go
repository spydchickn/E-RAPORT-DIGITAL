package controllers

import (
	"database/sql"

	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"e_raport_digital/config"
	"e_raport_digital/middlewares"
	"e_raport_digital/utils"
)

// VerifikasiItem represents a single pending nilai row for the template
type VerifikasiItem struct {
	IDNilai    int
	IDSiswa    int
	IDMapel    int
	Nilai      int
	Status     string
	NamaSiswa  string
	NamaMapel  string
}

func VerifikasiList(w http.ResponseWriter, r *http.Request) {
	// fetch pending nilai (status = 'pending')
	rows, err := config.DB.Query(`
        SELECT n.id_nilai, s.nama, m.nama_mapel, n.nilai_angka, n.nilai_huruf, n.semester, n.tahun_ajaran
    FROM nilai n
    JOIN siswa s ON s.id_siswa = n.id_siswa
    JOIN mapel m ON m.id_mapel = n.id_mapel
        WHERE n.status = 'pending'
        ORDER BY n.id_nilai DESC
    `)
	if err != nil {
		log.Println("Error querying nilai pending:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []VerifikasiItem
	for rows.Next() {
		var v VerifikasiItem
		if err := rows.Scan(&v.IDNilai, &v.IDSiswa, &v.IDMapel, &v.Nilai, &v.Status, &v.NamaSiswa, &v.NamaMapel); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		list = append(list, v)
	}

	// render into layout
	session, _ := middlewares.Store.Get(r, "session")
	role := ""
	if session.Values["role"] != nil {
		if rr, ok := session.Values["role"].(string); ok {
			role = rr
		}
	}

	content := template.HTML(utils.RenderPartial("admin/verifikasi_list.html", map[string]interface{}{
		"List": list,
	}))

	utils.Templates.ExecuteTemplate(w, "layouts/base.html", map[string]interface{}{
		"Title":   "Verifikasi Nilai",
		"Role":    role,
		"Alert":   "",
		"Content": content,
	})
}

func VerifikasiStore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idNilaiStr := r.FormValue("id_nilai")
	idNilai, err := strconv.Atoi(idNilaiStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	// Update nilai status to 'verified'
	res, err := config.DB.Exec("UPDATE nilai SET status = 'verified' WHERE id_nilai = ?", idNilai)
	if err != nil {
		log.Println("Error updating nilai status:", err)
		http.Error(w, "Unable to verify nilai", http.StatusInternalServerError)
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		http.Error(w, "Nilai not found", http.StatusNotFound)
		return
	}

	// Get siswa id and user_id to notify
	var idSiswa, toUserID int
	err = config.DB.QueryRow("SELECT id_siswa, id_user FROM siswa WHERE id_siswa = (SELECT id_siswa FROM nilai WHERE id_nilai = ?)", idNilai).Scan(&idSiswa, &toUserID)
	if err != nil {
		// If cannot get user id, just log and continue
		if err != sql.ErrNoRows {
			log.Println("Error fetching siswa for notification:", err)
		}
	} else {
		// Insert notification to notifikasi table
		_, err = config.DB.Exec(`INSERT INTO notifikasi (to_user_id, to_role, pesan, dibaca, tanggal) VALUES (?, ?, ?, ?, ?)`,
			toUserID, "siswa", "Nilai Anda telah diverifikasi.", false, time.Now())
		if err != nil {
			log.Println("Error creating notification:", err)
		}
	}

	// Log action in logs_aktivitas
	userID := utils.GetUserIDFromSession(r)
	_, err = config.DB.Exec(`INSERT INTO logs_aktivitas (user_id, action, table_name, record_id, details, tanggal) VALUES (?, ?, ?, ?, ?, ?)`,
		userID, "verify", "nilai", idNilai, "Nilai verified", time.Now())
	if err != nil {
		log.Println("Error logging action:", err)
	}

	http.Redirect(w, r, "/admin/verifikasi/list", http.StatusSeeOther)
}
