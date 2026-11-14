DELETE FROM tagihan_siswa
WHERE nis IN ('25260001', '25260297', '25260003')
  AND jenis_iuran_id IN (
      'f1000000-0000-0000-0000-000000000002',
      'f1000000-0000-0000-0000-000000000003'
  )
  AND periode_bulan = 11
  AND periode_tahun = 2026
  AND status IN ('BELUM_LUNAS', 'SEBAGIAN');
