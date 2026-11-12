<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Surat Keterangan Aktif Belajar</h1>
        <p class="text-caption text-grey-7 q-mb-none">Pelayanan penerbitan, verifikasi data, dan pencetakan surat keterangan siswa aktif SMK</p>
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- 1. PANEL KIRI: Form Pembuatan Surat -->
      <div class="col-12 col-md-5">
        <q-card class="shadow-1">
          <q-card-section class="bg-primary text-white">
            <div class="text-subtitle1 text-weight-bold">Buat Surat Keterangan Baru</div>
          </q-card-section>

          <q-card-section class="q-pa-md">
            <q-form @submit.prevent="generateSurat" class="q-gutter-sm">
              <q-select v-model="form.studentId" :options="studentOptions" label="Cari & Pilih Siswa Aktif"
                outlined dense emit-value map-options use-input @filter="filterSiswa"
                :rules="[val => !!val || 'Siswa wajib dipilih']" option-label="label" option-value="value" />

              <q-input v-model="form.letterNumber" label="Nomor Surat Resmi" outlined dense
                placeholder="Contoh: 421.5/102-SMK/2026" :rules="[val => !!val || 'Nomor surat wajib diisi']" />

              <q-input v-model="form.parentName" label="Nama Orang Tua / Wali" outlined dense
                :rules="[val => !!val || 'Nama orang tua wajib diisi']" />

              <q-select v-model="form.purpose" :options="purposeOptions" label="Tujuan / Keperluan Surat" outlined dense
                :rules="[val => !!val || 'Keperluan surat wajib ditentukan']" />

              <div class="row justify-end q-mt-md">
                <q-btn label="Reset" type="reset" color="grey" flat class="q-mr-sm" @click="resetForm" />
                <q-btn label="Terbitkan Surat" type="submit" color="primary" icon="print" :loading="letterStore.isLoading" />
              </div>
            </q-form>
          </q-card-section>
        </q-card>
      </div>

      <!-- 2. PANEL KANAN: Antrean -->
      <div class="col-12 col-md-7">
        <q-card class="shadow-1 fit">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div class="text-subtitle1 text-weight-bold text-grey-9">Antrean Penerbitan Surat</div>
          </q-card-section>
          <q-card-section class="q-pa-none">
            <q-table flat :rows="letterStore.letters" :columns="columns" row-key="ID" :loading="letterStore.isLoading" :pagination="{ rowsPerPage: 7 }">
              <template v-slot:body-cell-status="props">
                <q-td :props="props" class="text-center">
                  <q-chip size="xs" :color="props.row.Status === 'COMPLETED' ? 'green-1' : 'orange-1'"
                    :text-color="props.row.Status === 'COMPLETED' ? 'green-9' : 'orange-9'">
                    {{ props.row.Status === 'COMPLETED' ? 'Selesai' : 'Menunggu Cetak' }}
                  </q-chip>
                </q-td>
              </template>
              <template v-slot:body-cell-aksi="props">
                <q-td :props="props" class="text-center q-gutter-xs">
                  <q-btn size="sm" color="indigo" icon="print" round flat :disable="!props.row.FileURL" @click="printSurat(props.row)">
                    <q-tooltip>Cetak Lembar Surat</q-tooltip>
                  </q-btn>
                  <q-btn size="sm" color="green-8" icon="check_circle" round flat
                    :disable="props.row.Status === 'COMPLETED'" @click="completeSurat(props.row)">
                    <q-tooltip>Tandai Selesai Diambil</q-tooltip>
                  </q-btn>
                </q-td>
              </template>
            </q-table>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useLetterStore } from '@/stores/letterStore'
import { studentApi } from '@/adapters/api/studentApi'
import { useContextStore } from '@/stores/contextStore'

const $q = useQuasar()
const letterStore = useLetterStore()
const contextStore = useContextStore()

const form = ref({
  studentId: null,
  letterNumber: `421.5/0${letterStore.letters.length + 1}-SMK/2026`,
  parentName: '',
  purpose: 'Persyaratan BPJS Kesehatan'
})

const purposeOptions = [
  'Tunjangan Gaji Orang Tua (PNS/BUMN)', 
  'Persyaratan BPJS Kesehatan', 
  'Pengajuan Beasiswa Luar Sekolah', 
  'Pembuatan Paspor Magang / Kunjungan Kerja'
]

const studentOptions = ref([])
const masterStudents = ref([])

const columns = [
  { name: 'letterNumber', align: 'left', label: 'No. Surat', field: 'LetterNumber', sortable: true },
  { name: 'studentName', align: 'left', label: 'Nama Siswa', field: 'StudentName', sortable: true },
  { name: 'className', align: 'left', label: 'Kelas', field: 'ClassName' },
  { name: 'purpose', align: 'left', label: 'Keperluan', field: 'Purpose' },
  { name: 'status', align: 'center', label: 'Status', field: 'Status' },
  { name: 'aksi', align: 'center', label: 'Aksi' }
]

const filterSiswa = (val, update) => {
  if (val === '') {
    update(() => { studentOptions.value = masterStudents.value })
    return
  }
  update(() => {
    const needle = val.toLowerCase()
    studentOptions.value = masterStudents.value.filter(v => v.label.toLowerCase().includes(needle))
  })
}

async function fetchStudents() {
  try {
    const res = await studentApi.getStudents(contextStore.currentContext, { limit: 1000 })
    if (res.success) {
      masterStudents.value = res.data.items.map(s => ({
        label: `${s.fullName} (${s.nis}) - ${s.className || 'Tanpa Kelas'}`,
        value: s.id
      }))
      studentOptions.value = masterStudents.value
    }
  } catch (err) {
    console.error(err)
  }
}

function generateSurat() {
  if (!form.value.studentId) return
  letterStore.createLetter({
    studentId: form.value.studentId,
    letterNumber: form.value.letterNumber,
    parentName: form.value.parentName,
    purpose: form.value.purpose
  }).then(success => {
    if (success) {
      $q.notify({ color: 'positive', message: 'Surat berhasil dibuat dan PDF telah di-generate.', icon: 'check' })
      resetForm()
    } else {
      $q.notify({ color: 'negative', message: letterStore.error || 'Gagal membuat surat', icon: 'error' })
    }
  })
}

function resetForm() {
  form.value = {
    studentId: null,
    letterNumber: `421.5/0${letterStore.letters.length + 1}-SMK/2026`,
    parentName: '',
    purpose: 'Persyaratan BPJS Kesehatan'
  }
}

function printSurat(row) {
  if (row.FileURL) {
    window.open(row.FileURL, '_blank')
  }
}

function completeSurat(row) {
  letterStore.completeLetter(row.ID).then(success => {
    if (success) {
      $q.notify({ color: 'green-9', message: 'Status surat ditandai Selesai.', icon: 'check' })
    }
  })
}

onMounted(() => {
  fetchStudents()
  letterStore.fetchLetters()
})
</script>
