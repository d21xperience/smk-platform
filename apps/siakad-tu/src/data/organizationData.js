// src/data/organizationData.js
// Sumber data struktur organisasi sekolah, dipakai oleh StrukturOrganisasi.vue melalui komponen OrgChart.
//
// TODO: ganti dengan hasil GET /api/organisasi/struktur (Go backend).
//
// Setiap simpul (node):
//   id       : unik
//   nama     : nama pejabat/pemegang jabatan (boleh dikosongkan bila belum diisi)
//   jabatan  : nama jabatan/posisi
//   tipe     : salah satu key pada ORG_NODE_TYPES di bawah - menentukan gaya visual
//   children : array simpul di bawahnya (opsional)

export const ORG_NODE_TYPES = {
  pimpinan: { label: 'Pimpinan', icon: 'school', fill: '#0d2b4e', text: '#ffffff', border: '#0d2b4e' },
  wakil: { label: 'Wakil Pimpinan', icon: 'supervisor_account', fill: '#ffffff', text: '#1a4a7a', border: '#1a4a7a' },
  kaprog: { label: 'Kepala Program/Unit', icon: 'engineering', fill: '#e3f2fd', text: '#1a4a7a', border: '#90caf9' },
  staff: { label: 'Staf/Koordinator', icon: 'badge', fill: '#f5f5f5', text: '#424242', border: '#e0e0e0' },
  wali: { label: 'Wali Kelas/Guru', icon: 'person', fill: '#ffffff', text: '#616161', border: '#e0e0e0' },
  koordinatif: { label: 'Mitra Koordinatif', icon: 'diversity_3', fill: '#fff8e1', text: '#8d6e00', border: '#d4af37' }
}

// Simpul yang ditampilkan di samping Kepala Sekolah dengan garis putus-putus
// (hubungan koordinatif, bukan garis komando langsung)
export const sideNode = {
  id: 'komite',
  nama: 'H. Ahmad Suryadi',
  jabatan: 'Komite Sekolah',
  tipe: 'koordinatif'
}

export const organizationTree = {
  id: 'kepsek',
  nama: 'Drs. Wawan Setiawan, M.Pd.',
  jabatan: 'Kepala Sekolah',
  tipe: 'pimpinan',
  children: [
    {
      id: 'waka-kurikulum',
      nama: 'Yeti Sumiati, S.Pd.',
      jabatan: 'Wakasek Kurikulum',
      tipe: 'wakil'
    },
    {
      id: 'waka-kesiswaan',
      nama: 'Dedi Supriadi, S.Kom.',
      jabatan: 'Wakasek Kesiswaan',
      tipe: 'wakil'
    },
    {
      id: 'waka-sarpras',
      nama: 'Iyan Sofyan, S.T.',
      jabatan: 'Wakasek Sarana Prasarana',
      tipe: 'wakil'
    },
    {
      id: 'waka-hubin',
      nama: 'Rina Marlina, S.Pd.',
      jabatan: 'Wakasek Hubungan Industri',
      tipe: 'wakil'
    },
    {
      id: 'ka-tu',
      nama: 'Asep Kurniawan',
      jabatan: 'Kepala Tata Usaha',
      tipe: 'staff'
    },
    {
      id: 'koor-bk',
      nama: 'Siti Nurhaliza, S.Pd.',
      jabatan: 'Koordinator BK',
      tipe: 'staff'
    },
    {
      id: 'kaprog-rpl',
      nama: 'Ahmad Fauzi, S.Kom.',
      jabatan: 'Kaprog Rekayasa Perangkat Lunak',
      tipe: 'kaprog',
      children: [
        { id: 'wali-x-rpl-1', nama: 'Neng Fitriani, S.Pd.', jabatan: 'Wali Kelas X RPL 1', tipe: 'wali' },
        { id: 'wali-xi-rpl-1', nama: 'Fajar Nugraha, S.Kom.', jabatan: 'Wali Kelas XI RPL 1', tipe: 'wali' },
        { id: 'wali-xii-rpl-1', nama: 'Ahmad Fauzi, S.Kom.', jabatan: 'Wali Kelas XII RPL 1', tipe: 'wali' }
      ]
    },
    {
      id: 'kaprog-tkj',
      nama: 'Dedi Supriadi, S.Kom.',
      jabatan: 'Kaprog Teknik Komputer & Jaringan',
      tipe: 'kaprog'
    },
    {
      id: 'kaprog-tbsm',
      nama: 'Iyan Sofyan, S.T.',
      jabatan: 'Kaprog Teknik & Bisnis Sepeda Motor',
      tipe: 'kaprog'
    }
  ]
}
