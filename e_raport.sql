-- Simple DB template
CREATE DATABASE IF NOT EXISTS e_raport;
USE e_raport;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role ENUM('admin','guru','siswa') NOT NULL
);

-- siswa table: include relational columns id_kelas and id_user
CREATE TABLE IF NOT EXISTS siswa (
    id_siswa INT AUTO_INCREMENT PRIMARY KEY,
    nis VARCHAR(20) NOT NULL UNIQUE,
    nama VARCHAR(100) NOT NULL,
    alamat TEXT,
    id_kelas INT NULL,
    id_user INT NULL
);

-- guru table
CREATE TABLE IF NOT EXISTS guru (
    id_guru INT AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    nip VARCHAR(50),
    alamat TEXT
);

-- kelas table
CREATE TABLE IF NOT EXISTS kelas (
    id_kelas INT AUTO_INCREMENT PRIMARY KEY,
    nama_kelas VARCHAR(100) NOT NULL
);

-- mapel table
CREATE TABLE IF NOT EXISTS mapel (
    id_mapel INT AUTO_INCREMENT PRIMARY KEY,
    nama_mapel VARCHAR(100) NOT NULL
);

-- nilai table (relations to siswa and mapel)
CREATE TABLE IF NOT EXISTS nilai (
    id_nilai INT AUTO_INCREMENT PRIMARY KEY,
    id_siswa INT NOT NULL,
    id_mapel INT NOT NULL,
    nilai INT NOT NULL,
    status VARCHAR(20) DEFAULT 'draft',
    komentar TEXT,
    CONSTRAINT fk_nilai_siswa FOREIGN KEY (id_siswa) REFERENCES siswa(id_siswa) ON DELETE CASCADE,
    CONSTRAINT fk_nilai_mapel FOREIGN KEY (id_mapel) REFERENCES mapel(id_mapel) ON DELETE CASCADE
);

-- Add columns if upgrading existing DB
ALTER TABLE nilai ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'draft';
ALTER TABLE nilai ADD COLUMN IF NOT EXISTS komentar TEXT;

-- guru_mapel junction table for assigning subjects to teachers
CREATE TABLE IF NOT EXISTS guru_mapel (
    id_guru_mapel INT AUTO_INCREMENT PRIMARY KEY,
    id_guru INT NOT NULL,
    id_mapel INT NOT NULL,
    UNIQUE KEY unique_guru_mapel (id_guru, id_mapel),
    CONSTRAINT fk_guru_mapel_guru FOREIGN KEY (id_guru) REFERENCES guru(id_guru) ON DELETE CASCADE,
    CONSTRAINT fk_guru_mapel_mapel FOREIGN KEY (id_mapel) REFERENCES mapel(id_mapel) ON DELETE CASCADE
);

-- If you're upgrading an existing database, make sure siswa has the new columns
ALTER TABLE siswa
  ADD COLUMN IF NOT EXISTS id_kelas INT NULL,
  ADD COLUMN IF NOT EXISTS id_user INT NULL;

INSERT IGNORE INTO users (username, password, role) VALUES
('admin', 'admin123', 'admin'),
('guru01', 'guru123', 'guru'),
('siswa01', 'siswa123', 'siswa');

-- Add id_user to guru if it doesn't exist (safe for existing DBs)
-- Works by checking information_schema; compatible with MySQL
SET @col_exists := (
    SELECT COUNT(*)
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
        AND TABLE_NAME = 'guru'
        AND COLUMN_NAME = 'id_user'
);
-- If not exists, add column
SET @s := IF(@col_exists = 0, 'ALTER TABLE guru ADD COLUMN id_user INT NULL;', 'SELECT "column exists"');
PREPARE stmt FROM @s;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- New tables for expanded features

-- Wali Kelas (extends guru)
CREATE TABLE IF NOT EXISTS wali_kelas (
    id_wali INT AUTO_INCREMENT PRIMARY KEY,
    id_guru INT NOT NULL,
    id_kelas INT NOT NULL,
    UNIQUE KEY unique_wali_kelas (id_guru, id_kelas),
    CONSTRAINT fk_wali_guru FOREIGN KEY (id_guru) REFERENCES guru(id_guru) ON DELETE CASCADE,
    CONSTRAINT fk_wali_kelas FOREIGN KEY (id_kelas) REFERENCES kelas(id_kelas) ON DELETE CASCADE
);

-- Tahun Ajaran
CREATE TABLE IF NOT EXISTS tahun_ajaran (
    id_tahun INT AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(20) NOT NULL,
    aktif BOOLEAN DEFAULT FALSE
);

-- Semester
CREATE TABLE IF NOT EXISTS semester (
    id_semester INT AUTO_INCREMENT PRIMARY KEY,
    nama VARCHAR(10) NOT NULL,  -- e.g., 'Ganjil', 'Genap'
    id_tahun_ajaran INT NOT NULL,
    CONSTRAINT fk_semester_tahun FOREIGN KEY (id_tahun_ajaran) REFERENCES tahun_ajaran(id_tahun) ON DELETE CASCADE
);

-- Update mapel to add KKM
ALTER TABLE mapel
ADD COLUMN IF NOT EXISTS kkm INT DEFAULT 70;

-- Update siswa to add foto and ortu details
ALTER TABLE siswa
ADD COLUMN IF NOT EXISTS foto VARCHAR(255) NULL,
ADD COLUMN IF NOT EXISTS ortu_nama VARCHAR(100) NULL,
ADD COLUMN IF NOT EXISTS ortu_email VARCHAR(100) NULL,
ADD COLUMN IF NOT EXISTS ortu_telp VARCHAR(20) NULL;

-- Update guru to add foto
ALTER TABLE guru
ADD COLUMN IF NOT EXISTS foto VARCHAR(255) NULL;

-- Update users to add 2FA fields
ALTER TABLE users
ADD COLUMN IF NOT EXISTS otp_secret VARCHAR(255) NULL,
ADD COLUMN IF NOT EXISTS twofa_enabled BOOLEAN DEFAULT FALSE;

-- Update nilai to add status and verifikasi
ALTER TABLE nilai
ADD COLUMN IF NOT EXISTS status ENUM('draft', 'verified', 'final') DEFAULT 'draft',
ADD COLUMN IF NOT EXISTS verifikasi_by INT NULL,
ADD CONSTRAINT fk_verifikasi_guru FOREIGN KEY (verifikasi_by) REFERENCES guru(id_guru) ON DELETE SET NULL;

-- Presensi
CREATE TABLE IF NOT EXISTS presensi (
    id_presensi INT AUTO_INCREMENT PRIMARY KEY,
    id_siswa INT NOT NULL,
    id_mapel INT NOT NULL,
    id_guru INT NOT NULL,
    tanggal DATE NOT NULL,
    status ENUM('hadir', 'alfa', 'izin', 'sakit') NOT NULL,
    catatan TEXT NULL,
    CONSTRAINT fk_presensi_siswa FOREIGN KEY (id_siswa) REFERENCES siswa(id_siswa) ON DELETE CASCADE,
    CONSTRAINT fk_presensi_mapel FOREIGN KEY (id_mapel) REFERENCES mapel(id_mapel) ON DELETE CASCADE,
    CONSTRAINT fk_presensi_guru FOREIGN KEY (id_guru) REFERENCES guru(id_guru) ON DELETE CASCADE,
    UNIQUE KEY unique_presensi (id_siswa, id_mapel, tanggal)
);

-- Catatan Sikap
CREATE TABLE IF NOT EXISTS catatan_sikap (
    id_catatan INT AUTO_INCREMENT PRIMARY KEY,
    id_siswa INT NOT NULL,
    id_semester INT NOT NULL,
    deskripsi TEXT NOT NULL,
    nilai_sikap ENUM('sangat_baik', 'baik', 'cukup', 'kurang') NOT NULL,
    CONSTRAINT fk_catatan_siswa FOREIGN KEY (id_siswa) REFERENCES siswa(id_siswa) ON DELETE CASCADE,
    CONSTRAINT fk_catatan_semester FOREIGN KEY (id_semester) REFERENCES semester(id_semester) ON DELETE CASCADE
);

-- Ekstrakurikuler
CREATE TABLE IF NOT EXISTS ekstrakurikuler (
    id_ekstra INT AUTO_INCREMENT PRIMARY KEY,
    id_siswa INT NOT NULL,
    nama_ekstra VARCHAR(100) NOT NULL,
    nilai ENUM('sangat_baik', 'baik', 'cukup', 'kurang') NOT NULL,
    CONSTRAINT fk_ekstra_siswa FOREIGN KEY (id_siswa) REFERENCES siswa(id_siswa) ON DELETE CASCADE
);

-- Komentar
CREATE TABLE IF NOT EXISTS komentar (
    id_komentar INT AUTO_INCREMENT PRIMARY KEY,
    id_nilai INT NULL,
    id_siswa INT NOT NULL,
    dari_role ENUM('guru', 'wali_kelas', 'admin') NOT NULL,
    id_pengirim INT NOT NULL,  -- user_id
    pesan TEXT NOT NULL,
    tanggal TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_komentar_nilai FOREIGN KEY (id_nilai) REFERENCES nilai(id_nilai) ON DELETE CASCADE,
    CONSTRAINT fk_komentar_siswa FOREIGN KEY (id_siswa) REFERENCES siswa(id_siswa) ON DELETE CASCADE
);

-- Logs Aktivitas
CREATE TABLE IF NOT EXISTS logs_aktivitas (
    id_log INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    action VARCHAR(100) NOT NULL,  -- e.g., 'create_siswa', 'update_nilai'
    table_name VARCHAR(50) NOT NULL,
    record_id INT NULL,
    details TEXT NULL,
    tanggal TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_log_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Notifikasi
CREATE TABLE IF NOT EXISTS notifikasi (
    id_notif INT AUTO_INCREMENT PRIMARY KEY,
    to_user_id INT NULL,  -- specific user
    to_role ENUM('admin', 'guru', 'siswa', 'wali_kelas') NULL,  -- or role-based
    pesan TEXT NOT NULL,
    dibaca BOOLEAN DEFAULT FALSE,
    tanggal TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_notif_user FOREIGN KEY (to_user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Sample data for new tables
INSERT IGNORE INTO tahun_ajaran (nama, aktif) VALUES ('2024/2025', TRUE);
SET @tahun_id = (SELECT id_tahun FROM tahun_ajaran WHERE aktif = TRUE LIMIT 1);
INSERT IGNORE INTO semester (nama, id_tahun_ajaran) VALUES ('Ganjil', @tahun_id);

-- Sample wali_kelas (assign guru01 to kelas 1, assume kelas exists)
INSERT IGNORE INTO kelas (nama_kelas) VALUES ('X IPA 1');
INSERT IGNORE INTO guru (nama, nip, alamat) VALUES ('Guru Sample', '123456', 'Alamat Sample');
SET @guru_id = (SELECT id_guru FROM guru ORDER BY id_guru DESC LIMIT 1);
SET @kelas_id = (SELECT id_kelas FROM kelas WHERE nama_kelas = 'X IPA 1' LIMIT 1);
INSERT IGNORE INTO wali_kelas (id_guru, id_kelas) VALUES (@guru_id, @kelas_id);

-- Sample mapel with KKM
INSERT IGNORE INTO mapel (nama_mapel, kkm) VALUES ('Matematika', 75), ('Bahasa Indonesia', 70);

-- Sample presensi (assume siswa/mapel/guru exist)
-- INSERT IGNORE INTO presensi (id_siswa, id_mapel, id_guru, tanggal, status) VALUES (1, 1, 1, CURDATE(), 'hadir');
