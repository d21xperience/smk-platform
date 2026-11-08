CREATE TABLE IF NOT EXISTS "tabel_jurusan" (
	"jurusan_id" VARCHAR(25) NOT NULL,
	"nama_jurusan" VARCHAR(100) NOT NULL,
	"untuk_sma" NUMERIC(1,0) NOT NULL DEFAULT 0,
	"untuk_smk" NUMERIC(1,0) NOT NULL DEFAULT 1,
	"untuk_pt" NUMERIC(1,0) NOT NULL DEFAULT 0,
	"untuk_slb" NUMERIC(1,0) NOT NULL DEFAULT 0,
	"untuk_smklb" NUMERIC(1,0) NOT NULL DEFAULT 0,
	"jenjang_pendidikan_id" NUMERIC(2,0) NULL DEFAULT NULL,
	"jurusan_induk" VARCHAR(25) NULL DEFAULT NULL,
	"level_bidang_id" VARCHAR(5) NOT NULL DEFAULT '',
	"create_date" TIMESTAMP NOT NULL DEFAULT NOW(),
	"last_update" TIMESTAMP NOT NULL DEFAULT NOW(),
	"expired_date" TIMESTAMP NULL DEFAULT NULL,
	"last_sync" TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY ("jurusan_id")
);

-- Seed Data Jurusan (DAPODIK)
INSERT INTO tabel_jurusan (jurusan_id, nama_jurusan, level_bidang_id) VALUES 
('TKJ', 'Teknik Komputer dan Jaringan', 'SMK'),
('TKRO', 'Teknik Kendaraan Ringan Otomotif', 'SMK'),
('TBSM', 'Teknik dan Bisnis Sepeda Motor', 'SMK'),
('AKL', 'Akuntansi dan Keuangan Lembaga', 'SMK'),
('DPIB', 'Desain Pemodelan dan Informasi Bangunan', 'SMK'),
('OTKP', 'Otomatisasi dan Tata Kelola Perkantoran', 'SMK')
ON CONFLICT (jurusan_id) DO NOTHING;
