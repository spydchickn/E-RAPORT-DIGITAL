package models

import "e_raport_digital/config"

// Simple helpers to get counts for dashboard
func CountSiswa() (int, error) {
    var cnt int
    err := config.DB.QueryRow("SELECT COUNT(*) FROM siswa").Scan(&cnt)
    return cnt, err
}

func CountGuru() (int, error) {
    var cnt int
    err := config.DB.QueryRow("SELECT COUNT(*) FROM guru").Scan(&cnt)
    return cnt, err
}

func CountKelas() (int, error) {
    var cnt int
    err := config.DB.QueryRow("SELECT COUNT(*) FROM kelas").Scan(&cnt)
    return cnt, err
}

func CountMapel() (int, error) {
    var cnt int
    err := config.DB.QueryRow("SELECT COUNT(*) FROM mapel").Scan(&cnt)
    return cnt, err
}
