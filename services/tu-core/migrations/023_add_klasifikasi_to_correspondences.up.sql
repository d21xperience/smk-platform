ALTER TABLE correspondences ADD COLUMN kategori VARCHAR(50);
ALTER TABLE correspondences ADD COLUMN kode_klasifikasi VARCHAR(50);

CREATE INDEX idx_corr_kategori ON correspondences(kategori);
CREATE INDEX idx_corr_kode_klasifikasi ON correspondences(kode_klasifikasi);
