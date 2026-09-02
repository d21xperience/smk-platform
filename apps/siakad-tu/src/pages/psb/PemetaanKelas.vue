<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Otomatisasi Pembagian Kelas X</h1>
        <p class="text-caption text-grey-7 q-mb-none">Algoritma pembagian rombel berimbang (Balanced Split) dengan
          kapasitas maksimal 36 siswa per kelas</p>
      </div>
      <div>
        <q-btn color="orange-9" icon="gavel" label="Proses Ploting Kelas" @click="prosesPembagianKelasMassal"
          :disable="calonSiswaTerdaftar.length === 0" />
      </div>
    </div>

    <!-- Ringkasan Kuota Rombel -->
    <div class="row q-col-gutter-sm q-mb-md">
      <div v-for="jurusan in ['TKRO', 'RPL', 'AKL']" :key="jurusan" class="col-12 col-sm-4">
        <q-card flat class="shadow-1 bg-white">
          <q-card-section class="q-py-sm">
            <div class="text-caption text-grey-7 text-uppercase text-weight-bold">Simulasi Kelas {{ jurusan }}</div>
            <div class="text-h6 text-weight-bold text-grey-9 q-my-xs">
              {{ hitungTotalSiswa(jurusan) }} <span class="text-caption text-grey-6">Pendaftar</span>
            </div>
            <!-- Estimasi Rombel Terbentuk -->
            <q-badge :color="hitungKebutuhanKelas(jurusan) > 0 ? 'indigo' : 'grey-7'" class="q-pa-xs">
              Terbentuk: {{ hitungKebutuhanKelas(jurusan) }} Rombel
            </q-badge>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Hasil Pembagian Rombel Otomatis -->
    <div class="row q-col-gutter-md">
      <div class="col-12 col-md-4" v-for="(kelompok, namaRombel) in hasilRombel" :key="namaRombel">
        <q-card class="shadow-1">
          <q-card-section class="bg-indigo-9 text-white q-py-sm flex justify-between items-center">
            <div class="text-subtitle2 text-weight-bold">Kelas X {{ namaRombel }}</div>
            <q-badge color="white" text-color="indigo-9" :label="kelompok.length + ' Siswa'" class="text-weight-bold" />
          </q-card-section>

          <q-card-section class="q-pa-none" style="max-height: 350px; overflow-y: auto;">
            <q-list separator dense>
              <q-item v-for="(siswa, idx) in kelompok" :key="siswa.id" class="q-py-sm">
                <q-item-section avatar class="text-caption text-grey-6" style="min-width: 20px;">
                  {{ idx + 1 }}
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium text-grey-9">{{ siswa.nama }}</q-item-label>
                  <q-item-label caption>NIS: {{ siswa.nis }} | {{ siswa.gender }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>

      <!-- Kondisi jika belum diproses -->
      <div class="col-12 text-center q-pa-xl text-grey-6" v-if="Object.keys(hasilRombel).length === 0">
        <q-icon name="group_work" size="64px" color="grey-4" />
        <div class="text-subtitle1 q-mt-md">Klik tombol "Proses Ploting Kelas" di kanan atas untuk menyimulasikan
          pembagian kelas berimbang.</div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const MAX_KAPASITAS = 36 // Batas maksimal siswa per kelas standar Dapodik [1]

// Data Dummy Calon Siswa Terdaftar (Contoh: Jurusan TKRO sengaja dibuat 50 orang)
const calonSiswaTerdaftar = ref([
  // Anggap panitia sudah melakukan generate NIS dan verifikasi lunas keuangan
  ...Array.from({ length: 50 }, (_, i) => ({ id: `TKRO-${i}`, nama: `Siswa TKRO Abjad-${String.fromCharCode(65 + (i % 26))}${i}`, jurusan: 'TKRO', gender: 'L', nis: `262700${String(i + 1).padStart(2, '0')}` })),
  ...Array.from({ length: 37 }, (_, i) => ({ id: `RPL-${i}`, nama: `Siswa RPL Abjad-${String.fromCharCode(65 + (i % 26))}${i}`, jurusan: 'RPL', gender: 'P', nis: `262700${String(i + 51).padStart(2, '0')}` })),
  { id: 'AKL-1', nama: 'Deni Sukma', jurusan: 'AKL', gender: 'L', nis: '26270088' }
])

// State Penampung Hasil Akhir Pembagian Kelas
const hasilRombel = ref({})

// Fungsi Helper menghitung total pendaftar per jurusan
const hitungTotalSiswa = (jurusan) => {
  return calonSiswaTerdaftar.value.filter(s => s.jurusan === jurusan).length
}

// Fungsi menghitung kebutuhan jumlah kelas (Pembulatan ke atas)
const hitungKebutuhanKelas = (jurusan) => {
  const total = hitungTotalSiswa(jurusan)
  return Math.ceil(total / MAX_KAPASITAS) // Misal 50 / 36 = 1.38 -> dibulatkan jadi 2 kelas
}

// ALGORITMA UTAMA: BALANCED SPLIT AUTOMATION
const prosesPembagianKelasMassal = () => {
  const listJurusan = ['TKRO', 'RPL', 'AKL']
  const tempHasil = {}

  listJurusan.forEach((jurusan) => {
    // 1. Ambil semua siswa di jurusan tersebut dan urutkan berdasarkan Nama (A-Z)
    const siswaJurusan = calonSiswaTerdaftar.value
      .filter(s => s.jurusan === jurusan)
      .sort((a, b) => a.nama.localeCompare(b.nama))

    const totalSiswa = siswaJurusan.length
    if (totalSiswa === 0) return

    // 2. Tentukan berapa jumlah kelas yang dibutuhkan
    const jumlahKelas = Math.ceil(totalSiswa / MAX_KAPASITAS)

    // 3. Hitung jumlah dasar siswa per kelas (Metode Balanced Split)
    const baseSiswaPerKelas = Math.floor(totalSiswa / jumlahKelas)
    let sisaSiswa = totalSiswa % jumlahKelas // Sisa yang tidak habis dibagi rata

    // 4. Distribusikan siswa ke dalam nama kelas alfabetis (A, B, C, dst.)
    let indexSiswa = 0
    for (let i = 0; i < jumlahKelas; i++) {
      const labelAbjad = String.fromCharCode(65 + i) // 0 -> 'A', 1 -> 'B'
      const namaKelasLengkap = `${jurusan} ${labelAbjad}`

      // Jika masih ada sisa pembagian murni, kelas ini mendapat +1 siswa tambahan agar adil
      const targetJumlahSiswaKelasIni = baseSiswaPerKelas + (sisaSiswa > 0 ? 1 : 0)
      sisaSiswa--

      // Potong array siswa sesuai porsi kuota berimbang kelas ini
      tempHasil[namaKelasLengkap] = siswaJurusan.slice(indexSiswa, indexSiswa + targetJumlahSiswaKelasIni)
      indexSiswa += targetJumlahSiswaKelasIni
    }
  })

  hasilRombel.value = tempHasil
  $q.notify({
    color: 'green-9',
    message: 'Simulasi pembagian rombel berimbang berhasil diterapkan!',
    icon: 'done_all'
  })
}
</script>
