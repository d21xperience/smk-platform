<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Karantina PPDB & Keuangan</h1>
        <p class="text-caption text-grey-7 q-mb-none">Validasi data calon siswa, verifikasi berkas fisik, otomatisasi NIS, dan pencatatan angsuran biaya pendaftaran</p>
      </div>
      <div class="q-gutter-sm">
        <q-btn color="orange-9" icon="auto_awesome" label="Otomatisasi NIS Massal" @click="confirmGenerateNisMassal"
          :disable="registrations.length === 0 || isLoading" :loading="isLoading" />
        <q-btn color="primary" icon="person_add" label="Tambah Calon Siswa" @click="openAddDialog" />
      </div>
    </div>

    <!-- Widget Real-time Kuota & Keuangan -->
    <div class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-4">
        <q-card flat class="shadow-1 bg-white text-dark">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase text-weight-bold">Total Calon Siswa</div>
              <div class="text-h5 text-weight-bold text-primary">{{ registrations.length }} Siswa Baru</div>
            </div>
            <q-avatar color="blue-1" text-color="primary" icon="groups" size="40px" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card flat class="shadow-1 bg-white text-dark">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase text-weight-bold">Target Biaya Pendaftaran</div>
              <div class="text-h5 text-weight-bold text-grey-9">Rp {{ standarBiaya.toLocaleString('id-ID') }}</div>
            </div>
            <q-avatar color="grey-2" text-color="grey-8" icon="payments" size="40px" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card flat class="shadow-1 bg-teal-1 text-teal-10">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-teal-9 text-uppercase text-weight-bold flex items-center">
                Total Dana PPDB Masuk
                <q-btn flat round dense :icon="isDanaVisible ? 'visibility' : 'visibility_off'" color="teal-10"
                  size="xs" class="q-ml-xs" @click="isDanaVisible = !isDanaVisible">
                  <q-tooltip>{{ isDanaVisible ? 'Sembunyikan Nominal' : 'Tampilkan Nominal' }}</q-tooltip>
                </q-btn>
              </div>
              <div class="text-h5 text-weight-bold">
                {{ isDanaVisible ? 'Rp ' + totalDanaMasuk.toLocaleString('id-ID') : 'Rp ••••••••' }}
              </div>
            </div>
            <q-avatar color="teal-2" text-color="teal-10" icon="price_check" size="40px" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Tabel Karantina Calon Siswa PPDB -->
    <q-card class="shadow-1">
      <q-card-section class="q-py-sm flex justify-between items-center">
        <div class="text-subtitle2 text-weight-bold text-grey-9">Daftar Karantina Registrasi Siswa Baru</div>
        <q-input v-model="search" dense outlined placeholder="Cari nama calon siswa..." style="max-width: 250px;" @update:model-value="onSearch">
          <template v-slot:prepend><q-icon name="search" /></template>
        </q-input>
      </q-card-section>

      <q-card-section class="q-pa-none">
        <q-table flat :rows="registrations" :columns="columns" row-key="id" :pagination="{ rowsPerPage: 10 }" :loading="isLoading">
          
          <template v-slot:body-cell-jurusan="props">
            <q-td :props="props">
              <q-badge color="indigo-1" text-color="indigo" :label="props.row.jurusanDipilih" class="text-weight-bold" />
            </q-td>
          </template>

          <template v-slot:body-cell-keuangan="props">
            <q-td :props="props">
              <div class="text-center">
                <q-badge :color="props.row.biayaPendaftaran >= props.row.biayaStandar ? 'green' : props.row.biayaPendaftaran > 0 ? 'orange-9' : 'red'"
                  :label="props.row.biayaPendaftaran >= props.row.biayaStandar ? 'Lunas (Tunai)' : props.row.biayaPendaftaran > 0 ? 'Dicicil' : 'Belum Bayar'"
                  class="q-mb-xs" />
                <div class="text-caption text-grey-7" style="font-size: 11px;">
                  Masuk: Rp {{ props.row.biayaPendaftaran.toLocaleString('id-ID') }}
                </div>
                <div v-if="props.row.biayaStandar - props.row.biayaPendaftaran > 0" class="text-caption text-red text-weight-medium"
                  style="font-size: 10px;">
                  Sisa: Rp {{ (props.row.biayaStandar - props.row.biayaPendaftaran).toLocaleString('id-ID') }}
                </div>
              </div>
            </q-td>
          </template>

          <template v-slot:body-cell-berkas="props">
            <q-td :props="props" class="text-center">
              <div class="row justify-center q-gutter-xs">
                <q-icon :name="props.row.berkasIjazah ? 'check_circle' : 'cancel'" :color="props.row.berkasIjazah ? 'green' : 'red'" size="xs"><q-tooltip>Ijazah/SKL</q-tooltip></q-icon>
                <q-icon :name="props.row.berkasKK ? 'check_circle' : 'cancel'" :color="props.row.berkasKK ? 'green' : 'red'" size="xs"><q-tooltip>Kartu Keluarga</q-tooltip></q-icon>
                <q-icon :name="props.row.berkasAkta ? 'check_circle' : 'cancel'" :color="props.row.berkasAkta ? 'green' : 'red'" size="xs"><q-tooltip>Akta Lahir</q-tooltip></q-icon>
              </div>
            </q-td>
          </template>

          <template v-slot:body-cell-nis="props">
            <q-td :props="props" class="text-center">
              <q-badge v-if="props.row.nis" color="green-8" :label="props.row.nis" />
              <q-badge v-else color="grey-6" label="Belum Generate" />
            </q-td>
          </template>

          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center q-gutter-xs">
              <q-btn size="sm" color="green-8" icon="bolt" round flat
                :disable="!props.row.nis || !isLunasDanLengkap(props.row)" @click="aktivasiKeSiswaAktif(props.row)">
                <q-tooltip v-if="isLunasDanLengkap(props.row)">Aktifkan ke Direktori Utama</q-tooltip>
                <q-tooltip v-else>Gagal: Biaya Belum Lunas atau Berkas Fisik Tidak Lengkap!</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- DIALOG POPUP: Tambah Calon Siswa -->
    <q-dialog v-model="dialogOpen" persistent>
      <q-card style="min-width: 460px;">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">Registrasi Calon Siswa Baru</div>
        </q-card-section>

        <q-card-section class="q-gutter-sm q-pt-md">
          <div class="text-subtitle2 text-grey-8">Biodata Pokok Calon Siswa:</div>
          <q-input v-model="form.calonNama" label="Nama Lengkap Calon Siswa" outlined dense />
          <q-input v-model="form.asalSekolah" label="Sekolah Asal (SMP / MTs)" outlined dense />

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-select v-model="form.jurusanDipilih" :options="['TKRO', 'RPL', 'AKL']" label="Jurusan Pilihan" outlined dense />
            </div>
            <div class="col-6">
              <q-input v-model="form.calonNisn" label="NISN (Dari SMP)" outlined dense type="number" />
            </div>
          </div>

          <q-separator class="q-my-md" />

          <div class="text-subtitle2 text-grey-8">Administrasi Biaya Pendaftaran (Tunai/Cicil):</div>
          <div class="row q-col-gutter-sm items-center">
            <div class="col-6">
              <q-input v-model.number="form.biayaPendaftaran" label="Jumlah Uang yang Dibayarkan (Rp)" outlined dense type="number"
                prefix="Rp" :rules="[val => val <= standarBiaya || 'Pembayaran melebihi batas tarif sekolah!']" />
            </div>
            <div class="col-6 q-pb-md">
              <div class="text-caption text-grey-7">Tunggakan Sisa Angsuran:</div>
              <div class="text-subtitle2 text-weight-bold" :class="standarBiaya - form.biayaPendaftaran === 0 ? 'text-green' : 'text-red'">
                Rp {{ (standarBiaya - form.biayaPendaftaran).toLocaleString('id-ID') }}
              </div>
            </div>
          </div>

          <div class="row q-gutter-xs q-mb-sm">
            <q-btn label="Set Lunas (Tunai)" size="xs" color="green" flat class="bg-green-1" @click="form.biayaPendaftaran = standarBiaya" />
            <q-btn label="Set Belum Bayar" size="xs" color="red" flat class="bg-red-1" @click="form.biayaPendaftaran = 0" />
          </div>

          <q-separator class="q-my-sm" />

          <div class="text-subtitle2 text-grey-8">Verifikasi Dokumen Fisik Pendaftaran:</div>
          <div class="row q-mt-xs">
            <div class="col-4"><q-checkbox v-model="form.berkasIjazah" label="Ijazah / SKL" dense color="primary" /></div>
            <div class="col-4"><q-checkbox v-model="form.berkasKK" label="Krt. Keluarga" dense color="primary" /></div>
            <div class="col-4"><q-checkbox v-model="form.berkasAkta" label="Akta Lahir" dense color="primary" /></div>
          </div>
        </q-card-section>
        
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="grey" v-close-popup />
          <q-btn color="primary" label="Simpan Registrasi" @click="saveCalonSiswa" :loading="isLoading" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useRegistrationStore } from '@/stores/registration/registrationStore'

const $q = useQuasar()
const registrationStore = useRegistrationStore()

const standarBiaya = 450000
const search = ref('')
const dialogOpen = ref(false)
const isDanaVisible = ref(false)

const form = ref({
  calonNama: '',
  calonNisn: '',
  asalSekolah: '',
  jurusanDipilih: 'TKRO',
  calonJenisKelamin: 'L',
  biayaPendaftaran: 0,
  biayaStandar: standarBiaya,
  berkasIjazah: false,
  berkasKK: false,
  berkasAkta: false,
})

const columns = [
  { name: 'nama', align: 'left', label: 'Nama Calon Siswa', field: 'calonNama', sortable: true },
  { name: 'jurusan', align: 'left', label: 'Jurusan', field: 'jurusanDipilih', sortable: true },
  { name: 'keuangan', align: 'center', label: 'Status Biaya Masuk (Rp 450k)' },
  { name: 'berkas', align: 'center', label: 'Berkas Fisik' },
  { name: 'nis', align: 'center', label: 'NIS', field: 'nis' },
  { name: 'aksi', align: 'center', label: 'Aksi Karantina' },
]

const registrations = computed(() => registrationStore.registrations)
const isLoading = computed(() => registrationStore.isLoading)

const totalDanaMasuk = computed(() => {
  return registrations.value.reduce((total, siswa) => total + (siswa.biayaPendaftaran || 0), 0)
})

const isLunasDanLengkap = (siswa) => {
  const lunas = siswa.biayaPendaftaran >= siswa.biayaStandar
  const berkasLengkap = siswa.berkasIjazah && siswa.berkasKK && siswa.berkasAkta
  return lunas && berkasLengkap
}

// Search debounce
let searchTimeout
function onSearch() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    registrationStore.fetchRegistrations({ search: search.value })
  }, 500)
}

const openAddDialog = () => {
  form.value = {
    calonNama: '',
    calonNisn: '',
    asalSekolah: '',
    jurusanDipilih: 'TKRO',
    calonJenisKelamin: 'L',
    biayaPendaftaran: 0,
    biayaStandar: standarBiaya,
    berkasIjazah: false,
    berkasKK: false,
    berkasAkta: false,
  }
  dialogOpen.value = true
}

const saveCalonSiswa = async () => {
  if (!form.value.calonNama) return
  
  const success = await registrationStore.createRegistration(form.value)
  if (success) {
    $q.notify({ color: 'positive', message: 'Data pendaftar baru disimpan.', icon: 'check' })
    dialogOpen.value = false
  } else {
    $q.notify({ color: 'negative', message: registrationStore.error?.message || 'Gagal menyimpan', icon: 'error' })
  }
}

const confirmGenerateNisMassal = () => {
  $q.dialog({
    title: 'Generate NIS Massal',
    message: 'Proses otomatisasi pembuatan nomor induk siswa massal berurutan untuk angkatan baru tahun ajaran 2026/2027.',
    cancel: true,
    persistent: true,
  }).onOk(async () => {
    const prefix = '2627'
    const success = await registrationStore.generateNisMassal(prefix)
    if (success) {
      $q.notify({ color: 'green-9', message: 'NIS Massal sukses diterbitkan.', icon: 'auto_awesome' })
    } else {
      $q.notify({ color: 'negative', message: registrationStore.error?.message || 'Gagal generate NIS', icon: 'error' })
    }
  })
}

const aktivasiKeSiswaAktif = async (row) => {
  $q.dialog({
    title: 'Aktivasi Siswa',
    message: `Aktifkan ${row.calonNama} menjadi siswa aktif? Pastikan biaya lunas dan berkas lengkap.`,
    cancel: true,
    persistent: true,
  }).onOk(async () => {
    const success = await registrationStore.activateRegistration(row.id)
    if (success) {
      $q.notify({ color: 'green-9', message: `${row.calonNama} berhasil dipindahkan ke Direktori Utama Siswa Aktif.`, icon: 'done_all' })
    } else {
      $q.notify({ color: 'negative', message: registrationStore.error?.message || 'Gagal aktivasi', icon: 'error' })
    }
  })
}

onMounted(() => {
  registrationStore.fetchRegistrations()
})
</script>
