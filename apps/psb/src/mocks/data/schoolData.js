// ============================================================
// schoolData.js — Satu sumber data untuk seluruh website
// Ganti data di sini untuk update konten tanpa sentuh komponen
// ============================================================

export const SCHOOL_INFO = {
  name: 'SMK Pasundan Jatinaggor',
  tagline: 'Membentuk Generasi Vokasi Unggul & Siap Kerja',
  address: 'Jl. Kol. Achmad Syam, Cikeruh, Kec. Jatinangor, Kab. Sumedang, Jawa Barat',
  phone: '(022) 5432-1098',
  email: 'smkspasundanjatinangor@gmail.som',
  whatsapp: 'https://wa.me/6281234567890',
  instagram: 'https://instagram.com/smkpasja',
  facebook: 'https://facebook.com/smkpasja',
  youtube: 'https://youtube.com/@smkpasja',
  tiktok: 'https://tiktok.com/@smkpasja',
  // brosurUrl: 'https://smkpasja.sch.id/brosur.pdf',
  ppdbUrl: '/ppdb',
}

export const HERO_SLIDES = [
  {
    id: 1,
    img: 'https://images.unsplash.com/photo-1581091226825-a6a2a5aee158?w=1400&q=80',
    eyebrow: 'Selamat Datang',
    title: 'SMK Pasundan Jatinangor',
    subtitle: 'Membentuk Generasi Vokasi Unggul & Siap Kerja di Era Global',
    cta: { label: 'Pelajari Lebih Lanjut', to: '#jurusan' },
  },
  {
    id: 2,
    img: 'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=1400&q=80',
    eyebrow: 'Fasilitas Standar Industri',
    title: 'Praktik Kerja Nyata',
    subtitle: 'Lab modern dan mitra perusahaan internasional untuk pengalaman industri terbaik',
    cta: { label: 'Lihat Fasilitas', to: '#fasilitas' },
  },
  {
    id: 3,
    img: 'https://images.unsplash.com/photo-1540575467063-178a50c2df87?w=1400&q=80',
    eyebrow: 'PPDB 2025/2026 Dibuka!',
    title: 'Daftar Sekarang',
    subtitle: 'Kuota terbatas — amankan tempat putra-putri Anda di sekolah vokasi terbaik',
    cta: { label: 'Daftar PPDB', to: '/ppdb', highlight: true },
  },
]

export const STATS = [
  { icon: 'groups', value: '1.200+', label: 'Siswa Aktif' },
  { icon: 'handshake', value: '45+', label: 'Mitra Industri Aktif' },
  { icon: 'work_history', value: '92%', label: 'Lulusan Langsung Kerja' },
  { icon: 'emoji_events', value: '180+', label: 'Prestasi Kota / Provinsi' },
]

export const JURUSAN = [
  {
    id: 'tbsm',
    name: 'TBSM',
    fullName: 'Teknik & Bisnis Sepeda Motor',
    color: 'red-7',
    colorHex: '#c62828',
    bgLight: '#fff5f5',
    icon: 'two_wheeler',
    badge: '⭐ Kelas Khusus Astra',
    badgeColor: 'red-7',
    akreditasi: 'A',
    img: 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=600&q=80',
    mitra: 'PT Astra Honda Motor (AHM)',
    desc: 'Program unggulan yang mempersiapkan mekanik profesional. Disupervisi langsung oleh PT Astra Honda Motor, menguasai teknologi mesin injeksi terbaru, kelistrikan modern, dan manajemen bisnis bengkel digital standar AHASS.',
    keunggulan: ['Sertifikasi resmi AHM', 'Magang di bengkel AHASS', 'Kelas kewirausahaan bengkel'],
    to: '/jurusan/tbsm',
  },
  {
    id: 'tkj',
    name: 'TKJ',
    fullName: 'Teknik Komputer & Jaringan',
    color: 'blue-7',
    colorHex: '#1565c0',
    bgLight: '#f0f4ff',
    icon: 'router',
    badge: '🌐 Cisco & MikroTik Academy',
    badgeColor: 'blue-7',
    akreditasi: 'A',
    img: 'https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=600&q=80',
    mitra: 'PT Telkom & MikroTik',
    desc: 'Fokus pada perancangan infrastruktur jaringan berskala luas, administrasi server cloud, dan sistem pertahanan keamanan siber. Tersertifikasi langsung oleh Cisco CCNA dan MikroTik MTCNA.',
    keunggulan: ['Sertifikasi Cisco CCNA', 'Lab jaringan fiber optik', 'Proyek real cloud server'],
    to: '/jurusan/tkj',
  },
  {
    id: 'tkr',
    name: 'TKR',
    fullName: 'Teknik Kendaraan Ringan',
    color: 'orange-8',
    colorHex: '#e65100',
    bgLight: '#fff8f0',
    icon: 'directions_car',
    badge: '🔧 Standar Roda Empat',
    badgeColor: 'orange-8',
    akreditasi: 'A',
    img: 'https://images.unsplash.com/photo-1486262715619-67b85e0b08d3?w=600&q=80',
    mitra: 'Daihatsu & Isuzu Motors',
    desc: 'Membekali kompetensi perawatan mobil bensin dan diesel. Siswa menguasai sistem sasis, pemindah tenaga, kelistrikan bodi, serta komputer scanner diagnostic untuk mobil EFI modern.',
    keunggulan: [
      'Scanner diagnostic EFI',
      'Magang di dealer resmi',
      'Sertifikasi mekanik otomotif',
    ],
    to: '/jurusan/tkr',
  },
  {
    id: 'akl',
    name: 'AKL',
    fullName: 'Akuntansi & Keuangan Lembaga',
    color: 'green-7',
    colorHex: '#2e7d32',
    bgLight: '#f0fff4',
    icon: 'account_balance',
    badge: '💹 Fintech & Digital Accounting',
    badgeColor: 'green-7',
    akreditasi: 'A',
    img: 'https://images.unsplash.com/photo-1554224155-6726b3ff858f?w=600&q=80',
    mitra: 'Bank Syariah & Perbankan Daerah',
    desc: 'Fokus pada pembukuan keuangan korporasi, administrasi perpajakan, audit manufaktur, dan penguasaan piranti lunak MYOB Accounting dan Accurate Professional.',
    keunggulan: [
      'Sertifikat MYOB/Accurate',
      'PKL di bank & perusahaan',
      'Simulasi perpajakan online',
    ],
    to: '/jurusan/akl',
  },
  {
    id: 'otkp',
    name: 'OTKP',
    fullName: 'Otomatisasi Tata Kelola Perkantoran',
    color: 'green-7',
    colorHex: '#2e7d32',
    bgLight: '#f0fff4',
    icon: 'account_balance',
    badge: '💹 Fintech & Digital Accounting',
    badgeColor: 'green-7',
    akreditasi: 'A',
    img: 'https://images.unsplash.com/photo-1554224155-6726b3ff858f?w=600&q=80',
    mitra: 'Bank Syariah & Perbankan Daerah',
    desc: 'Fokus pada pembukuan keuangan korporasi, administrasi perpajakan, audit manufaktur, dan penguasaan piranti lunak MYOB Accounting dan Accurate Professional.',
    keunggulan: [
      'Sertifikat MYOB/Accurate',
      'PKL di bank & perusahaan',
      'Simulasi perpajakan online',
    ],
    to: '/jurusan/akl',
  },
]

export const TEACHERS = [
  {
    name: 'Haryono, S.T., M.T.',
    role: 'Kepala Program TBSM',
    expertise: 'S2 Teknik Mesin — Universitas Negeri',
    cert: 'Certified AHM Instructor',
    photo: 'https://api.dicebear.com/7.x/avataaars/svg?seed=haryono&backgroundColor=b6e3f4',
  },
  {
    name: 'Siti Rahmawati, S.Kom.',
    role: 'Kepala Program TKJ',
    expertise: 'S1 Teknik Informatika',
    cert: 'MikroTik MTCNA & Cisco CCNA',
    photo: 'https://api.dicebear.com/7.x/avataaars/svg?seed=siti&backgroundColor=ffd5dc',
  },
  {
    name: 'Budi Santoso, S.Pd.',
    role: 'Waka Kesiswaan',
    expertise: 'S1 Pendidikan Karakter',
    cert: 'Asesor Kompetensi BNSP',
    photo: 'https://api.dicebear.com/7.x/avataaars/svg?seed=budi&backgroundColor=c0aede',
  },
  {
    name: 'Dr. H. Ahmad Subarjo, M.Pd.',
    role: 'Kepala Sekolah',
    expertise: 'S3 Manajemen Pendidikan',
    cert: 'Pengawas Sekolah Berprestasi',
    photo: 'https://api.dicebear.com/7.x/avataaars/svg?seed=ahmad&backgroundColor=d1f4e0',
  },
]

export const AGENDA = [
  {
    day: '08',
    month: 'Jun',
    category: 'Ujian',
    title: 'Penilaian Akhir Semester (PAS) Genap',
    location: 'Ruang Kelas & Lab Komputer',
  },
  {
    day: '19',
    month: 'Jun',
    category: 'Kegiatan',
    title: 'Classmeeting & Turnamen Futsal OSIS Cup',
    location: 'GOR Indoor',
  },
  {
    day: '25',
    month: 'Jun',
    category: 'Rapat',
    title: 'Rapat Pleno Kenaikan Kelas TP 2025/2026',
    location: 'Aula Utama',
  },
]

export const HOLIDAY = {
  title: 'Hari Libur Nasional: Idul Adha 1447 H',
  dateRange: '27 – 28 Juni 2026',
  info: 'KBM tetap berjalan secara daring pada tanggal 29 Juni 2026.',
}

export const BKK_JOBS = [
  {
    company: 'PT Astra Honda Motor (AHM)',
    initial: 'A',
    position: 'Mekanik / Operator Assembly Line',
    location: 'Cikarang, Bekasi',
    majorTarget: 'TBSM',
    color: 'red-7',
    deadline: '15 Juni 2026',
    applyUrl: 'https://smkpasja.sch.id/loker/ahm',
    requirements: [
      'Alumni SMK Pasundan Jatinangor Jurusan TBSM',
      'Usia maksimal 22 tahun',
      'Menguasai troubleshooting mesin PGM-FI',
    ],
  },
  {
    company: 'PT Telekomunikasi Selular (Telkomsel)',
    initial: 'T',
    position: 'Junior Network Technical Support',
    location: 'Bandung Area',
    majorTarget: 'TKJ',
    color: 'blue-7',
    deadline: '20 Juni 2026',
    applyUrl: 'https://smkpasja.sch.id/loker/telkomsel',
    requirements: [
      'Alumni Jurusan TKJ (Pria/Wanita)',
      'Sertifikat MTCNA / CCNA menjadi nilai plus',
      'Paham konfigurasi Routing & Fiber Optik dasar',
    ],
  },
]

export const TESTIMONIES = [
  {
    name: 'Rian Hidayat',
    role: 'Alumni TBSM 2023 — Pemilik Bengkel Motor Mandiri',
    statement:
      'Berkat kurikulum Kelas Khusus Astra di SMK ini, saya tidak hanya menguasai mesin injeksi terbaru, tetapi juga diajari cara mengelola keuangan bengkel. Setahun lulus, saya langsung berani buka usaha sendiri.',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=rian&backgroundColor=ffd5dc',
  },
  {
    name: 'Ibu Ratna Sari',
    role: 'Orang Tua dari Bagas — Siswa TKJ Kelas XII',
    statement:
      'Sangat bersyukur menyekolahkan anak di sini. Kedisiplinannya luar biasa, gurunya sangat komunikatif. Perkembangan anak bisa dipantau langsung. Fasilitas lab komputernya sangat lengkap dan modern.',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=ratna&backgroundColor=b6e3f4',
  },
  {
    name: 'Deni Saputra',
    role: 'Alumni TKJ 2022 — Network Engineer PT Telkom',
    statement:
      'Sertifikasi CCNA yang difasilitasi sekolah membuat saya langsung diterima kerja sebelum wisuda. Pengalaman PKL di PT Telkom membuka koneksi yang sangat berharga untuk karier saya sekarang.',
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=deni&backgroundColor=c0aede',
  },
]

export const OSIS = [
  {
    name: 'Fathur Rahman',
    role: 'Ketua OSIS',
    major: 'Kelas XI — TKJ',
    photo: 'https://api.dicebear.com/7.x/avataaars/svg?seed=fathur&backgroundColor=b6e3f4',
  },
  {
    name: 'Andini Putri',
    role: 'Wakil Ketua OSIS',
    major: 'Kelas XI — TBSM',
    photo: 'https://api.dicebear.com/7.x/avataaars/svg?seed=andini&backgroundColor=ffd5dc',
  },
  {
    name: 'Rizky Pratama',
    role: 'Sekretaris Umum',
    major: 'Kelas X — TKJ',
    photo: 'https://api.dicebear.com/7.x/avataaars/svg?seed=rizky&backgroundColor=d1f4e0',
  },
]

export const NEWS = [
  {
    id: 1,
    title: 'Siswa TBSM Raih Juara 1 LKS Mekanik Sepeda Motor Tingkat Provinsi Jawa Barat 2026',
    date: '30 Mei 2026',
    author: 'Kesiswaan',
    img: 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=200&q=80',
    to: '/berita/1',
  },
  {
    id: 2,
    title: 'Tim TKJ SMK Pasundan Jatinangor Juarai Kompetisi Cyber Security Tingkat Nasional',
    date: '22 Mei 2026',
    author: 'Humas',
    img: 'https://images.unsplash.com/photo-1550751827-4bd374c3f58b?w=200&q=80',
    to: '/berita/2',
  },
  {
    id: 3,
    title:
      'Program PKL Kolaborasi dengan 5 Perusahaan Multinasional Resmi Dimulai Tahun Ajaran Baru',
    date: '15 Mei 2026',
    author: 'Humas',
    img: 'https://images.unsplash.com/photo-1521737852567-6949f3f9f2b5?w=200&q=80',
    to: '/berita/3',
  },
]

// Database dummy siswa berdasarkan kelas
export const mockStudentsDB = {
  101: [
    // X RPL 1
    { id: 1001, nis: '2627001', name: 'Ahmad Fauzi' },
    { id: 1002, nis: '2627002', name: 'Siti Aminah' },
    { id: 1003, nis: '2627003', name: 'Budi Santoso' },
    { id: 1004, nis: '2627004', name: 'Dewi Lestari' },
    { id: 1005, nis: '2627005', name: 'Eko Prasetyo' },
    { id: 1006, nis: '2627006', name: 'Fani Rahmawati' },
  ],
  102: [
    // XI TKR 2
    { id: 2001, nis: '2527001', name: 'Gilang Ramadhan' },
    { id: 2002, nis: '2527002', name: 'Hana Pertiwi' },
    { id: 2003, nis: '2527003', name: 'Indra Lesmana' },
    { id: 2004, nis: '2527004', name: 'Joko Anwar' },
  ],
}

// Data dummy jadwal mingguan
export const dummySchedule = [
  {
    day: 'Senin',
    class_id: 101,
    class_name: 'X RPL 1',
    subject: 'Pemrograman Dasar',
    start_time: '07:00',
    end_time: '08:30',
    room: 'Lab Komputer 1',
  },
  {
    day: 'Senin',
    class_id: 102,
    class_name: 'XI TKR 2',
    subject: 'Basis Data',
    start_time: '09:00',
    end_time: '10:30',
    room: 'Ruang 204',
  },
  {
    day: 'Rabu',
    class_id: 101,
    class_name: 'X RPL 1',
    subject: 'Pemrograman Web',
    start_time: '07:00',
    end_time: '09:30',
    room: 'Lab Komputer 1',
  },
  {
    day: 'Jumat',
    class_id: 103,
    class_name: 'XII AKL 1',
    subject: 'Komputer Akuntansi',
    start_time: '08:00',
    end_time: '10:00',
    room: 'Lab Akuntansi',
  },
]
