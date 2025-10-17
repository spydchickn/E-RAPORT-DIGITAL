package models

type Ekstrakurikuler struct {
    IDEkstra    int
    IDSiswa     int
    NamaEkstra  string
    Nilai       string  // 'sangat_baik', 'baik', 'cukup', 'kurang'
}
