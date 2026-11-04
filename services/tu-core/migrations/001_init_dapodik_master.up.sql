CREATE TABLE IF NOT EXISTS "tabel_kelas" (
	"rombongan_belajar_id" UUID NOT NULL PRIMARY KEY,
	"sekolah_id" UUID NOT NULL,
	"semester_id" VARCHAR(5) NOT NULL,
	"jurusan_id" VARCHAR(25) NULL DEFAULT NULL,
	"ptk_id" UUID NULL DEFAULT NULL,
	"nm_kelas" VARCHAR(30) NULL DEFAULT NULL,
	"tingkat_pendidikan_id" NUMERIC(2,0) NULL DEFAULT NULL,
	"jenis_rombel" NUMERIC(2,0) NULL DEFAULT NULL,
	"nama_jurusan_sp" VARCHAR(100) NULL DEFAULT NULL,
	"jurusan_sp_id" UUID NULL DEFAULT NULL,
	"kurikulum_id" SMALLINT NOT NULL
);

CREATE TABLE IF NOT EXISTS "tabel_ptk" (
	"ptk_id" UUID NOT NULL PRIMARY KEY,
	"nama" VARCHAR(100) NOT NULL,
	"nip" VARCHAR(18) NULL DEFAULT NULL,
	"jenis_ptk_id" NUMERIC(2,0) NOT NULL,
	"jenis_kelamin" VARCHAR(1) NOT NULL,
	"tempat_lahir" VARCHAR(32) NOT NULL,
	"tanggal_lahir" DATE NOT NULL,
	"nuptk" VARCHAR(16) NULL DEFAULT NULL,
	"alamat_jalan" VARCHAR(80) NOT NULL,
	"status_keaktifan_id" NUMERIC(2,0) NOT NULL,
	"soft_delete" NUMERIC(1,0) NOT NULL
);

CREATE TABLE IF NOT EXISTS "tabel_siswa" (
	"peserta_didik_id" UUID NOT NULL PRIMARY KEY,
	"nis" VARCHAR(20) NOT NULL UNIQUE,
	"nisn" VARCHAR(13) NULL DEFAULT NULL UNIQUE,
	"nm_siswa" VARCHAR(100) NOT NULL,
	"tempat_lahir" VARCHAR(50) NULL DEFAULT NULL,
	"tanggal_lahir" DATE NULL DEFAULT NULL,
	"jenis_kelamin" VARCHAR(1) NULL DEFAULT NULL,
	"agama" VARCHAR(25) NULL DEFAULT NULL,
	"alamat_siswa" TEXT NULL DEFAULT NULL,
	"telepon_siswa" VARCHAR(20) NULL DEFAULT NULL,
	"diterima_tanggal" DATE NULL DEFAULT NULL,
	"nm_ayah" VARCHAR(100) NULL DEFAULT NULL,
	"nm_ibu" VARCHAR(100) NULL DEFAULT NULL,
	"pekerjaan_ayah" VARCHAR(30) NULL DEFAULT NULL,
	"pekerjaan_ibu" VARCHAR(30) NULL DEFAULT NULL,
	"nm_wali" VARCHAR(100) NULL DEFAULT NULL,
	"pekerjaan_wali" VARCHAR(30) NULL DEFAULT NULL
);
