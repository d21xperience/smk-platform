CREATE TABLE IF NOT EXISTS "tabel_siswa_pelengkap" (
	"pelengkap_siswa_id" UUID NOT NULL PRIMARY KEY,
	"peserta_didik_id" UUID NOT NULL REFERENCES "tabel_siswa"("peserta_didik_id") ON DELETE CASCADE,
	"status_dalam_kel" VARCHAR(30) NULL DEFAULT NULL,
	"anak_ke" NUMERIC(3,0) NULL DEFAULT NULL,
	"sekolah_asal" VARCHAR(100) NULL DEFAULT NULL,
	"diterima_kelas" VARCHAR(20) NULL DEFAULT NULL,
	"alamat_ortu" TEXT NULL DEFAULT NULL,
	"telepon_ortu" VARCHAR(20) NULL DEFAULT NULL,
	"alamat_wali" TEXT NULL DEFAULT NULL,
	"telepon_wali" VARCHAR(20) NULL DEFAULT NULL,
	"foto_siswa" VARCHAR(100) NULL DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "tabel_anggotakelas" (
	"anggota_rombel_id" UUID NOT NULL PRIMARY KEY,
	"peserta_didik_id" UUID NOT NULL REFERENCES "tabel_siswa"("peserta_didik_id") ON DELETE CASCADE,
	"rombongan_belajar_id" UUID NOT NULL REFERENCES "tabel_kelas"("rombongan_belajar_id") ON DELETE CASCADE,
	"semester_id" VARCHAR(5) NOT NULL,
	CONSTRAINT unique_siswa_per_semester UNIQUE ("peserta_didik_id", "semester_id")
);

CREATE TABLE IF NOT EXISTS "tabel_kehadiran" (
	"kd_kehadiran" UUID NOT NULL PRIMARY KEY,
	"semester_id" VARCHAR(5) NOT NULL,
	"anggota_rombel_id" UUID NOT NULL REFERENCES "tabel_anggotakelas"("anggota_rombel_id") ON DELETE CASCADE,
	"sakit" NUMERIC(3,0) NULL DEFAULT NULL,
	"izin" NUMERIC(3,0) NULL DEFAULT NULL,
	"tanpa_keterangan" NUMERIC(3,0) NULL DEFAULT NULL,
	"peserta_didik_id" UUID NULL DEFAULT NULL,
	"semester" NUMERIC(1,0) NULL DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "tabel_pembelajaran" (
	"pembelajaran_id" UUID NOT NULL PRIMARY KEY,
	"rombongan_belajar_id" UUID NOT NULL REFERENCES "tabel_kelas"("rombongan_belajar_id") ON DELETE CASCADE,
	"mata_pelajaran_id" INTEGER NOT NULL,
	"semester_id" VARCHAR(5) NOT NULL,
	"ptk_terdaftar_id" UUID NULL DEFAULT NULL,
	"status_di_kurikulum" NUMERIC(2,0) NULL DEFAULT NULL,
	"nama_mata_pelajaran" VARCHAR(50) NULL DEFAULT NULL,
	"induk_pembelajaran" UUID NULL DEFAULT NULL,
	"is_dapo" NUMERIC(1,0) NULL DEFAULT '1'
);
