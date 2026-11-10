DROP TABLE IF EXISTS riwayat_kenaikan;
DROP TABLE IF EXISTS buku_kontrol;
ALTER TABLE ptk_buku_induk_ext DROP COLUMN IF EXISTS tmt_gaji_berkala;
ALTER TABLE ptk_buku_induk_ext DROP COLUMN IF EXISTS tmt_pangkat;
ALTER TABLE ptk_buku_induk_ext DROP COLUMN IF EXISTS golongan_id;
DROP TABLE IF EXISTS master_golongan;
