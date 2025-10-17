package routes

import (
	"net/http"

	"e_raport_digital/controllers"
	"e_raport_digital/middlewares"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
    r := mux.NewRouter()
    r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

    r.HandleFunc("/", controllers.ShowLogin).Methods("GET")
    r.HandleFunc("/login", controllers.ShowLogin).Methods("GET")
    r.HandleFunc("/login", controllers.Login).Methods("POST")
    r.HandleFunc("/logout", controllers.Logout).Methods("GET")

    // Admin routes (CRUD)
    r.HandleFunc("/admin/dashboard", middlewares.AuthMiddleware("admin", controllers.AdminDashboard)).Methods("GET")
    r.HandleFunc("/siswa/list", middlewares.AuthMiddleware("admin", controllers.SiswaList)).Methods("GET")
    r.HandleFunc("/siswa/create", middlewares.AuthMiddleware("admin", controllers.SiswaCreateForm)).Methods("GET")
    r.HandleFunc("/siswa/store", middlewares.AuthMiddleware("admin", controllers.SiswaStore)).Methods("POST")
    r.HandleFunc("/siswa/edit", middlewares.AuthMiddleware("admin", controllers.SiswaEditForm)).Methods("GET")
    r.HandleFunc("/siswa/update", middlewares.AuthMiddleware("admin", controllers.SiswaUpdate)).Methods("POST")
    r.HandleFunc("/siswa/delete", middlewares.AuthMiddleware("admin", controllers.SiswaDelete)).Methods("GET")

    // Admin Presensi routes
    r.HandleFunc("/admin/presensi/list", middlewares.AuthMiddleware("admin", controllers.PresensiList)).Methods("GET")
    r.HandleFunc("/admin/presensi/create", middlewares.AuthMiddleware("admin", controllers.PresensiCreateForm)).Methods("GET")
    r.HandleFunc("/admin/presensi/store", middlewares.AuthMiddleware("admin", controllers.PresensiStore)).Methods("POST")
    r.HandleFunc("/admin/presensi/edit", middlewares.AuthMiddleware("admin", controllers.PresensiEditForm)).Methods("GET")
    r.HandleFunc("/admin/presensi/update", middlewares.AuthMiddleware("admin", controllers.PresensiUpdate)).Methods("POST")
    r.HandleFunc("/admin/presensi/delete", middlewares.AuthMiddleware("admin", controllers.PresensiDelete)).Methods("GET")

    // Admin Verifikasi routes
    r.HandleFunc("/admin/verifikasi/list", middlewares.AuthMiddleware("admin", controllers.VerifikasiList)).Methods("GET")
    r.HandleFunc("/admin/verifikasi/store", middlewares.AuthMiddleware("admin", controllers.VerifikasiStore)).Methods("POST")

    // Admin CRUD for Guru, Mapel, Kelas
    r.HandleFunc("/guru/list", middlewares.AuthMiddleware("admin", controllers.GuruList)).Methods("GET")
    r.HandleFunc("/guru/create", middlewares.AuthMiddleware("admin", controllers.GuruCreateForm)).Methods("GET")
    r.HandleFunc("/guru/store", middlewares.AuthMiddleware("admin", controllers.GuruStore)).Methods("POST")
    r.HandleFunc("/guru/delete", middlewares.AuthMiddleware("admin", controllers.GuruDelete)).Methods("GET")
    r.HandleFunc("/guru/assign/{id}", middlewares.AuthMiddleware("admin", controllers.GuruAssignMapelForm)).Methods("GET")
    r.HandleFunc("/guru/assign-mapel", middlewares.AuthMiddleware("admin", controllers.StoreGuruMapel)).Methods("POST")

    r.HandleFunc("/mapel/list", middlewares.AuthMiddleware("admin", controllers.MapelList)).Methods("GET")
    r.HandleFunc("/mapel/create", middlewares.AuthMiddleware("admin", controllers.MapelCreateForm)).Methods("GET")
    r.HandleFunc("/mapel/store", middlewares.AuthMiddleware("admin", controllers.MapelStore)).Methods("POST")
    r.HandleFunc("/mapel/delete", middlewares.AuthMiddleware("admin", controllers.MapelDelete)).Methods("GET")

    r.HandleFunc("/kelas/list", middlewares.AuthMiddleware("admin", controllers.KelasList)).Methods("GET")
    r.HandleFunc("/kelas/create", middlewares.AuthMiddleware("admin", controllers.KelasCreateForm)).Methods("GET")
    r.HandleFunc("/kelas/store", middlewares.AuthMiddleware("admin", controllers.KelasStore)).Methods("POST")
    r.HandleFunc("/kelas/delete", middlewares.AuthMiddleware("admin", controllers.KelasDelete)).Methods("GET")

    // User management (admin)
    r.HandleFunc("/users/list", middlewares.AuthMiddleware("admin", controllers.UsersList)).Methods("GET")
    r.HandleFunc("/users/create", middlewares.AuthMiddleware("admin", controllers.UserCreateForm)).Methods("GET")
    r.HandleFunc("/users/store", middlewares.AuthMiddleware("admin", controllers.UserStore)).Methods("POST")
    r.HandleFunc("/users/delete", middlewares.AuthMiddleware("admin", controllers.UserDelete)).Methods("GET")

    // Guru routes (input nilai)
    r.HandleFunc("/guru/dashboard", middlewares.AuthMiddleware("guru", controllers.GuruDashboard)).Methods("GET")
    r.HandleFunc("/guru/input-nilai", middlewares.AuthMiddleware("guru", controllers.InputNilaiForm)).Methods("GET")
    r.HandleFunc("/guru/store-nilai", middlewares.AuthMiddleware("guru", controllers.StoreNilai)).Methods("POST")

    // Guru Presensi routes
    r.HandleFunc("/guru/presensi/input", middlewares.AuthMiddleware("guru", controllers.GuruPresensiInputForm)).Methods("GET")
    r.HandleFunc("/guru/presensi/store", middlewares.AuthMiddleware("guru", controllers.GuruPresensiStore)).Methods("POST")

    // Siswa routes (view raport)
    r.HandleFunc("/siswa/dashboard", middlewares.AuthMiddleware("siswa", controllers.SiswaDashboard)).Methods("GET")
    r.HandleFunc("/siswa/raport", middlewares.AuthMiddleware("siswa", controllers.SiswaRaport)).Methods("GET")
    r.HandleFunc("/siswa/raport/pdf", middlewares.AuthMiddleware("siswa", controllers.SiswaRaportPDF)).Methods("GET")

    // Notifikasi routes (all roles)
    r.HandleFunc("/notifikasi/list", middlewares.AuthMiddleware("", controllers.NotifikasiList)).Methods("GET")
    r.HandleFunc("/notifikasi/mark", middlewares.AuthMiddleware("", controllers.NotifikasiMarkRead)).Methods("POST")
    r.HandleFunc("/notifikasi/broadcast", middlewares.AuthMiddleware("admin", controllers.NotifikasiBroadcastForm)).Methods("GET")
    r.HandleFunc("/notifikasi/broadcast", middlewares.AuthMiddleware("admin", controllers.NotifikasiBroadcastStore)).Methods("POST")

    return r
}
