DELETE FROM buku_kas_umum WHERE referensi_tabel IN ('penerimaan_iuran', 'pengeluaran');
DELETE FROM pengeluaran;
DELETE FROM pengajuan_pengeluaran;
DELETE FROM penerimaan_iuran;
DELETE FROM tagihan_siswa;
