package utils

import (
    "fmt"
    "os"
    "path/filepath"
    "e_raport_digital/models"
    "github.com/jung-kurt/gofpdf"
)

// ExportRaportPDF generates a PDF report for the given raport data
func ExportRaportPDF(raport *models.Raport, filename string) error {
    if raport == nil {
        return fmt.Errorf("nil raport")
    }

    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)
    pdf.Cell(190.0, 10.0, "Laporan Raport Siswa")
    pdf.Ln(15.0)

    pdf.SetFont("Arial", "", 12)
    pdf.Cell(50.0, 10.0, fmt.Sprintf("Nama Siswa: %s", raport.NamaSiswa))
    pdf.Ln(7.0)
    pdf.Cell(50.0, 10.0, fmt.Sprintf("Kelas: %s", raport.NamaKelas))
    pdf.Ln(15.0)

    // Table header
    pdf.SetFont("Arial", "B", 12)
    pdf.Cell(60.0, 10.0, "Mata Pelajaran")
    pdf.Cell(30.0, 10.0, "Nilai")
    pdf.Cell(30.0, 10.0, "Grade")
    pdf.Ln(10.0)

    // Table rows
    pdf.SetFont("Arial", "", 12)
    for _, nilai := range raport.NilaiList {
        pdf.Cell(60.0, 10.0, nilai.NamaMapel)
        pdf.Cell(30.0, 10.0, fmt.Sprintf("%d", nilai.Nilai))
        pdf.Cell(30.0, 10.0, nilai.Grade)
        pdf.Ln(10.0)
    }

    // Average
    pdf.Ln(10.0)
    pdf.SetFont("Arial", "B", 12)
    pdf.Cell(120.0, 10.0, fmt.Sprintf("Rata-rata: %.2f", raport.RataRata))

    // Save file
    // Ensure exports directory exists
    _ = os.MkdirAll("exports", 0o755)
    fullPath := filepath.Join("exports", filename)
    err := pdf.OutputFileAndClose(fullPath)
    if err != nil {
        return err
    }
    return nil
}
