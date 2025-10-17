package models

type CatatanSikap struct {
    IDCatatan   int
    IDSiswa     int
    IDSemester  int
    Deskripsi   string
    NilaiSikap  string  // 'sangat_baik', 'baik', 'cukup', 'kurang'
}
