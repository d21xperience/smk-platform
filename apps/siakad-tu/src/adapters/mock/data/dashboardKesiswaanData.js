/**
 * Data statis untuk Dashboard Kesiswaan
 *
 * Data ini digunakan oleh Mock Adapter untuk mensimulasikan response.
 * Saat migrasi ke real API, data ini akan digantikan oleh response dari backend.
 */

export const dashboardMockData = {
  // ---- KPI ----
  totalSiswaAktif: 1248,
  totalPelanggaranBulanIni: 17,
  totalBerprestasi: 32,
  totalEkskul: 486,

  // ---- Siswa Perlu Perhatian Khusus ----
  siswaPerhatian: [
    {
      nama: 'Rian Firmansyah',
      kelas: 'XI TBSM 2',
      catatan: 'Alpa berulang & poin pelanggaran tinggi',
      poin: 85,
    },
    {
      nama: 'Dedi Kurniawan',
      kelas: 'X TKJ 1',
      catatan: 'Terlibat perkelahian di lingkungan sekolah',
      poin: 70,
    },
  ],

  // ---- Tren Pelanggaran (6 Bulan Terakhir) ----
  trenPelanggaran: [
    { label: 'Feb', jumlah: 12 },
    { label: 'Mar', jumlah: 15 },
    { label: 'Apr', jumlah: 9 },
    { label: 'Mei', jumlah: 14 },
    { label: 'Jun', jumlah: 8 },
    { label: 'Jul', jumlah: 17 },
  ],

  // ---- Distribusi Poin Pelanggaran per Tingkat Kelas ----
  poinPerTingkat: [
    { label: 'Kelas X', total: 420, warna: '#0a192f' },
    { label: 'Kelas XI', total: 610, warna: '#c9a227' },
    { label: 'Kelas XII', total: 260, warna: '#00796b' },
  ],

  // ---- Kasus Terbaru ----
  kasusTerbaru: [
    {
      id: 1,
      judul: 'Terlambat masuk sekolah lebih dari 3 kali dalam sepekan',
      siswa: 'Ahmad Fauzan',
      kelas: 'X TKJ 2',
      waktu: '2 jam lalu',
      kategori: 'Ringan',
    },
    {
      id: 2,
      judul: 'Tidak mengenakan atribut seragam lengkap',
      siswa: 'Siti Nurhaliza',
      kelas: 'XI Akuntansi',
      waktu: 'Kemarin, 09:15',
      kategori: 'Ringan',
    },
    {
      id: 3,
      judul: 'Membawa rokok elektrik ke lingkungan sekolah',
      siswa: 'Rian Firmansyah',
      kelas: 'XI TBSM 2',
      waktu: '2 hari lalu',
      kategori: 'Berat',
    },
    {
      id: 4,
      judul: 'Bolos pada jam pelajaran ke-5 dan ke-6',
      siswa: 'Dedi Kurniawan',
      kelas: 'X TKJ 1',
      waktu: '3 hari lalu',
      kategori: 'Sedang',
    },
  ],

  // ---- Ekstrakurikuler ----
  ekstrakurikuler: [
    {
      nama: 'Pramuka',
      pembina: 'Bpk. Sutrisno',
      peserta: 142,
      icon: 'terrain',
      warna: 'green-8',
    },
    {
      nama: 'Futsal',
      pembina: 'Bpk. Andi Wijaya',
      peserta: 68,
      icon: 'sports_soccer',
      warna: 'indigo-8',
    },
    {
      nama: 'PMR',
      pembina: 'Ibu Ratna Sari',
      peserta: 54,
      icon: 'medical_services',
      warna: 'red-6',
    },
    {
      nama: 'Robotika',
      pembina: 'Bpk. Yusuf Hidayat',
      peserta: 38,
      icon: 'precision_manufacturing',
      warna: 'deep-purple-6',
    },
  ],

  // ---- Prestasi Terbaru ----
  prestasiList: [
    {
      judul: 'Juara 1 Mekanik Edukasi Tingkat Provinsi',
      siswa: 'Bayu Aji Santoso - XII TBSM 1',
      tanggal: '30 Mei 2026',
    },
    {
      judul: 'Juara 2 LKS Bidang IT Network Systems Administration',
      siswa: 'Nadia Putri - XI TKJ 1',
      tanggal: '18 Mei 2026',
    },
    {
      judul: 'Juara 3 Lomba Debat Bahasa Inggris Kabupaten',
      siswa: 'Fajar Ramadhan - XI Akuntansi',
      tanggal: '05 Mei 2026',
    },
  ],

  // ---- Agenda Kesiswaan ----
  agendaList: [
    {
      tanggal: '13',
      bulan: 'Jul',
      judul: 'Pembinaan Karakter & Kedisiplinan Siswa Kelas X',
      waktu: '08:00 - Aula Sekolah',
    },
    {
      tanggal: '19',
      bulan: 'Jul',
      judul: 'Seleksi Ketua OSIS Periode 2026/2027',
      waktu: '09:00 - Ruang OSIS',
    },
    {
      tanggal: '25',
      bulan: 'Jul',
      judul: 'Rapat Evaluasi Tata Tertib Semester Ganjil',
      waktu: '13:00 - Ruang Guru',
    },
  ],
}
