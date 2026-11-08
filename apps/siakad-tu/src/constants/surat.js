export const KATEGORI_SURAT = [
  {
    label: '1. Kelompok Surat Berdasarkan Aturan Resmi Dinas',
    value: 'RULES_DINAS',
    children: [
      { label: 'Surat Keputusan (SK)', value: 'SK' },
      { label: 'Surat Tugas (ST)', value: 'ST' },
      { label: 'Surat Perintah Perjalanan Dinas (SPPD)', value: 'SPPD' },
      { label: 'Surat Pengantar (SP)', value: 'SP' },
      { label: 'Surat Edaran (SE)', value: 'SE' },
      { label: 'Surat Undangan Dinas (SU)', value: 'SU' },
      { label: 'Surat Pengumuman', value: 'SPENGUMUMAN' },
      { label: 'Berita Acara (BA)', value: 'BA' }
    ]
  },
  {
    label: '2. Kelompok Surat Khusus Hubungan Industri (Khas SMK)',
    value: 'HUBUNGAN_INDUSTRI',
    children: [
      { label: 'Nota Kesepahaman / MoU', value: 'MOU' },
      { label: 'Permohonan Tempat PKL', value: 'PERMOHONAN_PKL' },
      { label: 'Tugas Pembimbing PKL', value: 'TUGAS_PEMBIMBING' },
      { label: 'Pengantar Penyerahan Siswa PKL', value: 'PENGANTAR_PKL' },
      { label: 'Penarikan Siswa PKL', value: 'PENARIKAN_PKL' }
    ]
  },
  {
    label: '3. Kelompok Surat Hubungan Orang Tua & Kesiswaan',
    value: 'ORTU_KESISWAAN',
    children: [
      { label: 'Undangan Orang Tua/Wali', value: 'UNDANGAN_ORTU' },
      { label: 'Surat Keterangan Aktif Sekolah', value: 'AKTIF_SEKOLAH' },
      { label: 'Surat Keterangan Lulus (SKL)', value: 'SKL' },
      { label: 'Surat Keterangan Kelakuan Baik (SKKB)', value: 'SKKB' },
      { label: 'Panggilan Siswa/Orang Tua', value: 'PANGGILAN_ORTU' },
      { label: 'Izin / Rekomendasi Kegiatan', value: 'REKOMENDASI_LOMBA' }
    ]
  },
  {
    label: '4. Kelompok Surat Administrasi Internal & Logistik',
    value: 'ADMIN_INTERNAL',
    children: [
      { label: 'Peminjaman Sarana dan Prasarana', value: 'PEMINJAMAN_SARPRAS' },
      { label: 'Permohonan Barang / Nota Dinas Intern', value: 'NOTA_DINAS' },
      { label: 'Laporan Pelaksanaan Tugas', value: 'LAPORAN_TUGAS' }
    ]
  }
]

// Fungsi untuk membuat Flat Array yang dikenali Quasar
export const FLAT_SURAT = KATEGORI_SURAT.flatMap(kategori => 
  kategori.children.map(item => ({
    label: item.label,
    value: item.value,
    group: kategori.label // Ini akan menjadi header grup di dropdown
  }))
)

// Helper untuk mencari nama label dari value
export function formatJenisSurat(val) {
  for (const grup of KATEGORI_SURAT) {
    const item = grup.children.find(c => c.value === val)
    if (item) return item.label
  }
  return val ? val.replace(/_/g, ' ') : '-'
}
