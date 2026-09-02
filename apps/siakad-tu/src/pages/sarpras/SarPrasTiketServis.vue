<template>
  <q-page class="q-pa-md">
    <!-- PANEL FILTER TAB STATUS -->
    <div class="row items-center justify-between q-mb-md">
      <q-tabs v-model="tabAktif" dense class="text-grey shadow-1 rounded-borders bg-white" active-color="primary"
        indicator-color="primary" align="left" narrow-indicator>
        <q-tab name="Semua" label="Semua Laporan" />
        <q-tab name="Pending" class="text-negative">
          <div class="row items-center no-wrap">
            <span>Baru (Pending)</span>
            <q-badge color="red" floating class="q-ml-xs">{{ hitungTiket('Pending') }}</q-badge>
          </div>
        </q-tab>
        <q-tab name="Proses" class="text-warning" label="Diproses" />
        <q-tab name="Selesai" class="text-positive" label="Selesai" />
      </q-tabs>

      <q-btn color="primary" icon="add_comment" label="Buat Laporan Baru" @click="dialogBaruTerbuka = true" />
    </div>

    <!-- LIST KARTU TIKET LAPORAN -->
    <div class="row q-col-gutter-md">
      <div v-for="tiket in tiketTersaring" :key="tiket.id" class="col-12 col-md-6">
        <q-card flat bordered class="tiket-card">
          <!-- Batas Samping Berwarna Berdasarkan Tingkat Urgensi -->
          <div :class="`urgensi-bar bg-${getUrgensiColor(tiket.urgensi)}`"></div>

          <q-card-section class="q-pl-lg">
            <div class="row items-center no-wrap justify-between">
              <div class="row items-center">
                <q-badge :color="getStatusColor(tiket.status)" class="q-mr-sm text-uppercase">
                  {{ tiket.status }}
                </q-badge>
                <span class="text-caption text-grey-7">#{{ tiket.id }} - {{ tiket.tanggal }}</span>
              </div>

              <q-chip dense :color="getUrgensiColor(tiket.urgensi) + '-1'"
                :text-color="getUrgensiColor(tiket.urgensi) + '-9'" class="text-weight-bold">
                Urgensi: {{ tiket.urgensi }}
              </q-chip>
            </div>

            <div class="text-h6 text-grey-9 q-mt-sm">{{ tiket.namaBarang }}</div>
            <div class="text-body2 text-grey-7 q-mt-xs">
              <q-icon name="room" size="xs" /> <strong>Lokasi:</strong> {{ tiket.lokasi }}
            </div>
            <div class="text-body2 text-grey-8 q-mt-sm bg-grey-2 q-pa-sm rounded-borders">
              "{{ tiket.deskripsiMasalah }}"
            </div>

            <q-separator class="q-my-md" />

            <div class="row items-center justify-between">
              <div class="text-caption text-grey-7">
                Pelapor: <strong>{{ tiket.pelapor }}</strong>
                <div v-if="tiket.biaya > 0" class="text-positive text-weight-bold q-mt-xs">
                  Biaya Servis: Rp {{ tiket.biaya.toLocaleString('id-ID') }}
                </div>
              </div>

              <!-- Tombol Respons Tindakan Staf TU -->
              <q-btn outline dense color="primary" icon="edit_note" label="Update Tindakan"
                @click="bukaDialogTindakan(tiket)" />
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- TAMPILAN JIKA TIDAK ADA TIKET -->
      <div v-if="tiketTersaring.length === 0" class="col-12 text-center q-pa-xl text-grey-5">
        <q-icon name="assignment_turned_in" size="64px" class="q-mb-sm" />
        <div class="text-h6">Tidak ada tiket laporan dalam kategori ini.</div>
      </div>
    </div>

    <!-- DIALOG POP-UP: UPDATE STATUS DAN BIAYA SERVIS -->
    <q-dialog v-model="dialogTindakanTerbuka">
      <q-card style="width: 450px; max-width: 900px;">
        <q-card-section class="bg-primary text-white row items-center">
          <div class="text-h6">Tindak Lanjut Laporan #{{ FormTindakan.id }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <div class="text-weight-bold text-grey-8">{{ FormTindakan.namaBarang }}</div>

          <!-- Ubah Status Pekerjaan -->
          <q-select outlined dense v-model="FormTindakan.status" :options="['Pending', 'Proses', 'Selesai']"
            label="Ubah Status Penanganan *" />

          <!-- Input Biaya Nyata (Sangat penting untuk Yayasan) -->
          <q-input outlined dense type="number" v-model.number="FormTindakan.biaya"
            label="Realisasi Biaya Perbaikan (Rp)" prefix="Rp" hint="Isi 0 jika diperbaiki mandiri / masih garansi" />

          <!-- Catatan Teknisi -->
          <q-input outlined dense type="textarea" v-model="FormTindakan.catatan"
            label="Catatan Solusi Perbaikan / Nama Vendor" rows="3" />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Simpan Pembaruan" @click="simpanTindakan" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const tabAktif = ref('Semua')
const dialogTindakanTerbuka = ref(false)
const dialogBaruTerbuka = ref(false)

// 1. Data Mock Tiket Kerusakan Fasilitas Sekolah
const daftarTiket = ref([
  { id: 'TKT-001', tanggal: '2026/06/05', namaBarang: 'AC Split Sharp 1 PK', lokasi: 'Ruang Kelas 11-B', pelapor: 'Bu Sri (Wali Kelas)', deskripsiMasalah: 'AC mengeluarkan air bocor deras ke lantai kelas saat jam pelajaran.', urgensi: 'Tinggi', status: 'Pending', biaya: 0, catatan: '' },
  { id: 'TKT-002', tanggal: '2026/06/04', namaBarang: 'Proyektor BenQ Wall Mount', lokasi: 'Ruang Lab Fisika', pelapor: 'Pak Bambang', deskripsiMasalah: 'Lampu proyektor berkedip merah lalu mati sendiri setelah 5 menit dinyalakan.', urgensi: 'Tinggi', status: 'Proses', biaya: 0, catatan: '' },
  { id: 'TKT-003', tanggal: '2026/06/01', namaBarang: 'Engsel Pintu Utama Aluminium', lokasi: 'Toilet Siswa Barat', pelapor: 'Pak Joko (Penjaga)', deskripsiMasalah: 'Baut engsel lepas sehingga pintu agak miring dan sulit dikunci dari dalam.', urgensi: 'Rendah', status: 'Selesai', biaya: 75000, catatan: 'Ganti baut baja baru oleh tim kebersihan internal.' }
])

// 2. State Form untuk Penanganan Tiket
const FormTindakan = ref({
  id: '',
  namaBarang: '',
  status: '',
  biaya: 0,
  catatan: ''
})

// 3. Computed Properti untuk Memfilter Tiket berdasarkan Tab aktif
const tiketTersaring = computed(() => {
  if (tabAktif.value === 'Semua') return daftarTiket.value
  return daftarTiket.value.filter(tiket => tiket.status === tabAktif.value)
})

function hitungTiket(status) {
  return daftarTiket.value.filter(tiket => tiket.status === status).length
}

// 4. Helper Warna Responsif UI
function getUrgensiColor(urgensi) {
  if (urgensi === 'Tinggi') return 'red'
  if (urgensi === 'Sedang') return 'amber'
  return 'green'
}

function getStatusColor(status) {
  if (status === 'Pending') return 'grey-7'
  if (status === 'Proses') return 'orange'
  return 'green'
}

// 5. Manajemen Aksi Dialog
function bukaDialogTindakan(tiket) {
  FormTindakan.value = {
    id: tiket.id,
    namaBarang: tiket.namaBarang,
    status: tiket.status,
    biaya: tiket.biaya,
    catatan: tiket.catatan
  }
  dialogTindakanTerbuka.value = true
}

function simpanTindakan() {
  const tiket = daftarTiket.value.find(t => t.id === FormTindakan.value.id)
  if (tiket) {
    tiket.status = FormTindakan.value.status
    tiket.biaya = FormTindakan.value.biaya
    tiket.catatan = FormTindakan.value.catatan

    $q.notify({
      color: 'positive',
      icon: 'cloud_done',
      message: `Status Tiket Laporan #${tiket.id} berhasil diperbarui.`
    })
  }
}
</script>

<style scoped>
.tiket-card {
  position: relative;
  overflow: hidden;
  transition: transform 0.2s;
}

.tiket-card:hover {
  transform: translateY(-2px);
}

.urgensi-bar {
  position: absolute;
  top: 0;
  left: 0;
  bottom: 0;
  width: 6px;
}
</style>
