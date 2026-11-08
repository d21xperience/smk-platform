<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Buku Induk & Klaper</h1>
        <p class="text-caption text-grey-7 q-mb-none">Pencatatan data kronologis arsip statis kesiswaan SMK</p>
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- 1. PANEL KIRI: Indeks Buku Klaper (Daftar Alfabet & Siswa) -->
      <div class="col-12 col-md-4">
        <q-card class="shadow-1">
          <q-card-section class="q-pb-none">
            <div class="text-subtitle2 text-weight-bold text-grey-9 q-mb-sm">Indeks Klaper (Abjad)</div>
            <!-- Baris Navigasi Alfabet -->
            <div class="flex q-gutter-xs justify-center q-mb-md">
              <q-btn v-for="letter in alphabets" :key="letter" :color="selectedLetter === letter ? 'primary' : 'grey-3'"
                :text-color="selectedLetter === letter ? 'white' : 'dark'" size="xs" padding="xs"
                style="min-width: 26px" @click="selectLetter(letter)" dense>
                {{ letter }}
              </q-btn>
            </div>

            <q-input v-model="searchKlaper" dense outlined placeholder="Cari nama di abjad ini..." clearable
              class="q-mb-sm">
              <template v-slot:prepend><q-icon name="search" /></template>
            </q-input>
          </q-card-section>

          <!-- Daftar Siswa Hasil Klaper -->
          <q-card-section class="q-pa-none" style="max-height: 400px; overflow-y: auto;">
            <q-list separator>
              <q-item v-for="siswa in klaperList" :key="siswa.nis" clickable v-ripple
                :active="activeSiswa?.nis === siswa.nis" active-class="bg-blue-1 text-primary text-weight-bold"
                @click="loadBukuInduk(siswa)">
                <q-item-section avatar>
                  <q-avatar color="primary" text-color="white" size="32px">
                    {{ siswa.nama.charAt(0) }}
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label>{{ siswa.nama }}</q-item-label>
                  <q-item-label caption>NIS: {{ siswa.nis }} | Kelas: {{ siswa.kelas }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-icon name="chevron_right" />
                </q-item-section>
              </q-item>
              <q-item v-if="klaperList.length === 0" class="text-center text-grey-6 q-pa-md">
                Tidak ada data siswa berawalan huruf "{{ selectedLetter }}"
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>

      <!-- 2. PANEL KANAN: Detail Lembar Buku Induk Siswa -->
      <div class="col-12 col-md-8">
        <q-card class="shadow-1 fit flex flex-center text-center q-pa-xl" v-if="!activeSiswa">
          <div>
            <q-icon name="menu_book" size="64px" color="grey-4" />
            <div class="text-subtitle1 text-grey-6 q-mt-md">Silakan pilih nama siswa di panel Klaper sebelah kiri untuk
              melihat lembar Buku Induk.</div>
          </div>
        </q-card>

        <q-card class="shadow-1" v-else>
          <!-- Header Profil Buku Induk -->
          <q-card-section class="bg-primary text-white row items-center justify-between">
            <div class="row items-center q-gutter-md">
              <q-avatar size="60px" class="bg-white text-primary text-weight-bold">
                {{ activeSiswa.nama.charAt(0) }}
              </q-avatar>
              <div>
                <div class="text-h6 text-weight-bold">{{ activeSiswa.nama }}</div>
                <div class="text-caption text-blue-2">No. Buku Induk: BI/{{ activeSiswa.nis }}/SMK</div>
              </div>
            </div>
            <q-btn icon="print" color="white" text-color="primary" label="Cetak Lembar Induk" @click="printBukuInduk" />
          </q-card-section>

          <!-- Navigasi Tab Data Lembar Buku Induk -->
          <q-tabs v-model="activeTab" dense class="text-grey" active-color="primary" indicator-color="primary"
            align="justify" narrow-indicator>
            <q-tab name="pribadi" label="Data Pribadi" />
            <q-tab name="akademik" label="Pendidikan" />
            <q-tab name="keluarga" label="Orang Tua/Wali" />
            <q-tab name="fisik" label="Fisik & Kesehatan" />
          </q-tabs>

          <q-separator />

          <!-- Konten Lembar Buku Induk (Read Only - Standard Arsip) -->
          <q-card-section class="q-pa-md">
            <q-tab-panels v-model="activeTab" animated>
              <!-- TAB 1: DATA PRIBADI -->
              <q-tab-panel name="pribadi" class="q-pa-none">
                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Nama Lengkap"
                      v-model="detailData.nama" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="NISN" v-model="detailData.nisn"
                      dense />
                  </div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Tempat, Tanggal Lahir"
                      v-model="detailData.ttl" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Jenis Kelamin"
                      v-model="detailData.jk" dense /></div>
                  <div class="col-12"><q-input readonly stack-label type="textarea" label="Alamat Tempat Tinggal"
                      v-model="detailData.alamat" dense rows="2" /></div>
                </div>
              </q-tab-panel>

              <!-- TAB 2: AKADEMIK / ASAL SEKOLAH -->
              <q-tab-panel name="akademik" class="q-pa-none">
                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Kompetensi Keahlian (Jurusan)"
                      v-model="detailData.jurusan" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Rombel / Kelas Saat Ini"
                      v-model="detailData.kelas" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Asal Sekolah (SMP/MTs)"
                      v-model="detailData.asalSekolah" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="No. Ijazah Kelulusan SMP"
                      v-model="detailData.noIjazahSmp" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Tanggal Diterima di SMK"
                      v-model="detailData.tglDiterima" dense /></div>
                </div>
              </q-tab-panel>

              <!-- TAB 3: KELUARGA / ORANG TUA -->
              <q-tab-panel name="keluarga" class="q-pa-none">
                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Nama Ayah Kandung"
                      v-model="detailData.namaAyah" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Pekerjaan Ayah"
                      v-model="detailData.kerjaAyah" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Nama Ibu Kandung"
                      v-model="detailData.namaIbu" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Pekerjaan Ibu"
                      v-model="detailData.kerjaIbu" dense /></div>
                  <div class="col-12"><q-input readonly stack-label label="No. Telepon / HP Orang Tua"
                      v-model="detailData.telpOrangTua" dense /></div>
                </div>
              </q-tab-panel>

              <!-- TAB 4: KONDISI FISIK (Khas Syarat Masuk DUDI SMK) -->
              <q-tab-panel name="fisik" class="q-pa-none">
                <div class="row q-col-gutter-sm">
                  <div class="col-12 col-sm-4"><q-input readonly stack-label label="Golongan Darah"
                      v-model="detailData.golDarah" dense /></div>
                  <div class="col-12 col-sm-4"><q-input readonly stack-label label="Tinggi Badan (cm)"
                      v-model="detailData.tinggi" dense /></div>
                  <div class="col-12 col-sm-4"><q-input readonly stack-label label="Berat Badan (kg)"
                      v-model="detailData.berat" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Status Buta Warna"
                      v-model="detailData.butaWarna" dense /></div>
                  <div class="col-12 col-sm-6"><q-input readonly stack-label label="Riwayat Penyakit Berat"
                      v-model="detailData.riwayatPenyakit" dense /></div>
                  <div class="col-12 q-mt-sm">
                    <q-banner class="bg-blue-1 text-primary rounded-borders dense">
                      <q-icon name="info" class="q-mr-xs" /> Data fisik diperbarui berkala untuk keperluan pemetaan
                      kelayakan seleksi kerja industri (DUDI).
                    </q-banner>
                  </div>
                </div>
              </q-tab-panel>
            </q-tab-panels>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// State Navigasi Klaper Alfabet
const alphabets = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'.split('')
const selectedLetter = ref('A')
const searchKlaper = ref('')

// State Detail Tampilan Buku Induk
const activeSiswa = ref(null)
const activeTab = ref('pribadi')

// Data Dummy List Siswa untuk Indeks Klaper
const masterSiswaKlaper = [
  { nis: '22230101', nama: 'Ahmad Fauzi', kelas: 'XII RPL 1' },
  { nis: '22230112', nama: 'Anisa Rahmawati', kelas: 'XII AKL 2' },
  { nis: '22230190', nama: 'Andi Wijaya', kelas: 'XII TKRO 3' },
  { nis: '22230102', nama: 'Budi Santoso', kelas: 'XII TKRO 3' },
  { nis: '23240201', nama: 'Chandra Wijaya', kelas: 'XI RPL 2' },
  { nis: '23240202', nama: 'Dinda Lestari', kelas: 'XI AKL 1' }]

// Data Dummy Record Buku Induk Detail yang Di-load ketika Nama Siswa Diklik
const mockDatabaseBukuInduk = { '22230101': { nama: 'Ahmad Fauzi', nisn: '0061234561', ttl: 'Bandung, 14 April 2008', jk: 'Laki-laki', alamat: 'Jl. Raya Soreang No. 45, Kab. Bandung', jurusan: 'Rekayasa Perangkat Lunak', kelas: 'XII RPL 1', asalSekolah: 'SMP Negeri 1 Soreang', noIjazahSmp: 'DN-01/DI/00998877', tglDiterima: '17 Juli 2024', namaAyah: 'Hendra Fauzi', kerjaAyah: 'Wiraswasta', namaIbu: 'Siti Aminah', kerjaIbu: 'Ibu Rumah Tangga', telpOrangTua: '081234567890', golDarah: 'O', tinggi: '172', berat: '65', butaWarna: 'Tidak Buta Warna', riwayatPenyakit: 'Tidak Ada' }, '22230112': { nama: 'Anisa Rahmawati', nisn: '0061234599', ttl: 'Soreang, 20 Agustus 2008', jk: 'Perempuan', alamat: 'Perum Gading Tutuka Blok C3, Bandung', jurusan: 'Akuntansi & Keuangan Lembaga', kelas: 'XII AKL 2', asalSekolah: 'MTs Al-Jawami', noIjazahSmp: 'DN-01/MTs/00223344', tglDiterima: '17 Juli 2024', namaAyah: 'Supriatna', kerjaAyah: 'PNS', namaIbu: 'Iis Dahlia', kerjaIbu: 'Guru', telpOrangTua: '085799887766', golDarah: 'A', tinggi: '158', berat: '50', butaWarna: 'Tidak Buta Warna', riwayatPenyakit: 'Asma (Ringan)' } }
// State Penampung Detail Data Aktif di Form
const detailData = ref({})
// Filter Pintar untuk List Klaper (Abjad Huruf Depan + Input Teks Pencarian)
const klaperList = computed(() => {
  return masterSiswaKlaper.filter(siswa => {
    const startWithLetter = siswa.nama.toUpperCase().startsWith(selectedLetter.value)
    const matchSearch = searchKlaper.value ? siswa.nama.toLowerCase().includes(searchKlaper.value.toLowerCase()) || siswa.nis.includes(searchKlaper.value) : true
    return startWithLetter && matchSearch
  })
})
// Fungsi Eksekusi Pengalihan Huruf Abjad Klaper
const selectLetter = (letter) => {
  selectedLetter.value = letter
  selectLetter.value = ''// Reset pencarian teks jika abjad berpindah
}
// // Fungsi Mengambil Data Lengkap dari 'Database' Frontend
const loadBukuInduk = (siswa) => {
  activeSiswa.value = siswa
  activeTab.value = 'pribadi' // Kembalikan ke tab pertama setiap ganti siswa
  // Ambil data jika tersedia di mock db, jika tidak buat objek kosong
  if (mockDatabaseBukuInduk[siswa.nis]) {
    detailData.value = { ...mockDatabaseBukuInduk[siswa.nis] }
  } else { detailData.value = { nama: siswa.nama, nisn: '-', ttl: '-', jk: '-', alamat: '-', jurusan: '-', kelas: siswa.kelas, asalSekolah: '-', noIjazahSmp: '-', tglDiterima: '-', namaAyah: '-', kerjaAyah: '-', namaIbu: '-', kerjaIbu: '-', telpOrangTua: '-', golDarah: '-', tinggi: '-', berat: '-', butaWarna: '-', riwayatPenyakit: '-' } }
}
// Fungsi Cetak Berkas Dokumen
const printBukuInduk = () => {
  $q.notify({ color: 'primary', message: `Menyiapkan tata letak cetak dokumen Buku Induk untuk: ${activeSiswa.value.nama}`, icon: 'print' })
}
</script>
<style scoped>
/* Membuat list klaper terlihat lebih rapi */
.q-item--active {
  border-left: 4px solid var(--q - primary);
}
</style>
