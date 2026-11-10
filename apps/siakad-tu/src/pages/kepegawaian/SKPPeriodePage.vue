<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Periode Penilaian SKP</h1>
        <p class="text-caption text-grey-7 q-mb-none">Manajemen Siklus Penilaian Kinerja Pegawai</p>
      </div>
      <q-btn color="primary" icon="add" label="Buat Periode Baru" @click="openDialog" />
    </div>

    <q-linear-progress v-if="isLoading" indeterminate color="primary" class="q-mb-md" />

    <q-card class="shadow-1">
      <q-card-section>
        <q-table :rows="periodes" :columns="columns" row-key="id" flat bordered dense :loading="isLoading">
          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center">
              <q-badge :color="props.row.status === 'OPEN' ? 'green-6' : 'grey-6'" class="q-pa-sm">
                {{ props.row.status }}
              </q-badge>
            </q-td>
          </template>
          <template v-slot:body-cell-tanggal="props">
            <q-td :props="props" class="text-center">
              {{ formatDate(props.row.tanggalMulai) }} - {{ formatDate(props.row.tanggalSelesai) }}
            </q-td>
          </template>
          <template v-slot:body-cell-deadline="props">
            <q-td :props="props" class="text-center text-negative text-weight-bold">
              {{ formatDate(props.row.deadlinePengisian) }}
            </q-td>
          </template>
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center">
              <q-btn flat round dense icon="visibility" color="primary" @click="lihatPenilaian(props.row.id)" class="q-mr-sm">
                <q-tooltip>Lihat Daftar Pegawai</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- Dialog Form Periode -->
    <q-dialog v-model="formDialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section class="row items-center bg-primary text-white">
          <div class="text-h6">Buat Periode Penilaian</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-form @submit="submitPeriode">
            <q-input v-model="form.namaPeriode" label="Nama Periode (Cth: SKP Tahunan 2026)" outlined dense :rules="[v => !!v || 'Wajib diisi']" />
            
            <div class="row q-col-gutter-sm q-mt-sm">
              <div class="col-6">
                <q-select v-model="form.jenisPeriode" :options="['TAHUNAN', 'SEMESTERAN', 'BULANAN']" label="Jenis Periode" outlined dense :rules="[v => !!v || 'Wajib diisi']" />
              </div>
              <div class="col-6">
                <q-input v-model.number="form.tahunReferensi" label="Tahun Referensi" type="number" outlined dense :rules="[v => !!v || 'Wajib diisi']" />
              </div>
            </div>

            <div class="row q-col-gutter-sm q-mt-sm">
              <div class="col-6">
                <q-input v-model="form.tanggalMulai" label="Tanggal Mulai" type="date" outlined dense :rules="[v => !!v || 'Wajib diisi']" />
              </div>
              <div class="col-6">
                <q-input v-model="form.tanggalSelesai" label="Tanggal Selesai" type="date" outlined dense :rules="[v => !!v || 'Wajib diisi']" />
              </div>
            </div>

            <q-input v-model="form.deadlinePengisian" label="Batas Waktu Pengisian" type="date" outlined dense class="q-mt-sm" :rules="[v => !!v || 'Wajib diisi']" />

            <div class="row justify-end q-mt-md">
              <q-btn label="Batal" flat color="negative" v-close-popup class="q-mr-sm" />
              <q-btn label="Simpan Periode" type="submit" color="primary" :loading="isSubmitting" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useRouter } from 'vue-router'
import { api } from '@/boot/axios'

const $q = useQuasar()
const router = useRouter()

const periodes = ref([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const formDialog = ref(false)

const form = ref({
  namaPeriode: '',
  jenisPeriode: 'TAHUNAN',
  tahunReferensi: new Date().getFullYear(),
  tanggalMulai: '',
  tanggalSelesai: '',
  deadlinePengisian: ''
})

const columns = [
  { name: 'namaPeriode', label: 'Nama Periode', field: 'namaPeriode', align: 'left', sortable: true },
  { name: 'jenisPeriode', label: 'Jenis', field: 'jenisPeriode', align: 'center' },
  { name: 'tahunReferensi', label: 'Tahun', field: 'tahunReferensi', align: 'center' },
  { name: 'tanggal', label: 'Periode Penilaian', align: 'center' },
  { name: 'deadline', label: 'Deadline', align: 'center' },
  { name: 'status', label: 'Status', align: 'center' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const formatDate = (val) => {
  if (!val) return '-'
  return new Date(val).toLocaleDateString('id-ID')
}

const fetchPeriodes = async () => {
  isLoading.value = true
  try {
    const res = await api.get('/kepegawaian/skp/periodes')
    periodes.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal memuat data periode SKP' })
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

const openDialog = () => {
  form.value = { namaPeriode: '', jenisPeriode: 'TAHUNAN', tahunReferensi: new Date().getFullYear(), tanggalMulai: '', tanggalSelesai: '', deadlinePengisian: '' }
  formDialog.value = true
}

const submitPeriode = async () => {
  isSubmitting.value = true
  try {
    await api.post('/kepegawaian/skp/periodes', form.value)
    $q.notify({ type: 'positive', message: 'Periode berhasil dibuat' })
    formDialog.value = false
    await fetchPeriodes()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Gagal membuat periode' })
    console.error(err)
  } finally {
    isSubmitting.value = false
  }
}

const lihatPenilaian = (id) => {
  router.push(`/kepegawaian/skp/periode/${id}`)
}

onMounted(fetchPeriodes)
</script>
