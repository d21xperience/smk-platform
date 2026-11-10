<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Usulan Kenaikan</h1>
        <p class="text-caption text-grey-7 q-mb-none">Pusat Usulan Kenaikan Pangkat, Gaji Berkala, & Penghargaan</p>
      </div>
      <q-btn color="primary" icon="refresh" label="Refresh" @click="fetchData" :loading="isLoading" />
    </div>

    <q-linear-progress v-if="isLoading" indeterminate color="primary" class="q-mb-md" />

    <!-- SECTION 1: Saran Usulan -->
    <div class="row items-center q-mb-sm">
      <q-icon name="lightbulb" color="orange" size="sm" class="q-mr-sm" />
      <div class="text-subtitle1 text-weight-bold">Saran Usulan (Dari Buku Kontrol)</div>
      <q-space />
      <q-badge color="orange" class="q-pa-sm">{{ saranUsulan.length }} Pegawai</q-badge>
    </div>

    <div v-if="saranUsulan.length === 0" class="q-mb-lg">
      <q-card class="bg-green-1 text-green-9 q-pa-md">
        <q-icon name="check_circle" class="q-mr-sm" /> Tidak ada pegawai yang perlu diusulkan saat ini. Semua usulan sudah diproses.
      </q-card>
    </div>

    <div v-else class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-md-6" v-for="saran in saranUsulan" :key="saran.ptkId">
        <q-card class="shadow-1 border-left-warning">
          <q-card-section class="row items-center q-pb-none">
            <div class="text-h6">{{ saran.namaPegawai }}</div>
            <q-space />
            <q-badge :color="getStatusColor(saran.statusMonitoring)" class="q-pa-sm">
              {{ saran.statusMonitoring }}
            </q-badge>
          </q-card-section>
          <q-card-section class="q-pt-sm">
            <div class="text-caption text-grey-7">NIP: {{ saran.nip }}</div>
            <div class="text-body2 q-mt-sm">
              <strong>Golongan Sekarang:</strong> {{ saran.golonganSekarang }} <br />
              <strong>Jenis Usulan:</strong> {{ saran.jenisUsulan }} <br />
              <strong>Selisih Hari:</strong> <span class="text-negative">{{ saran.selisihHari }} hari</span>
            </div>
          </q-card-section>
          <q-card-actions align="right">
            <q-btn flat color="primary" label="Lihat Detail" />
            <q-btn color="secondary" label="Buat Usulan" @click="openCreateDialog(saran)" />
          </q-card-actions>
        </q-card>
      </div>
    </div>

    <q-separator class="q-my-lg" />

    <!-- SECTION 2: Daftar Usulan -->
    <div class="text-subtitle1 text-weight-bold q-mb-sm">Riwayat Usulan</div>
    <q-card class="shadow-1">
      <q-card-section>
        <q-table :rows="usulanList" :columns="columns" row-key="id" flat bordered dense>
          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center">
              <q-badge :color="getUsulanStatusColor(props.row.status)" class="q-pa-sm">
                {{ props.row.status }}
              </q-badge>
            </q-td>
          </template>
          <template v-slot:body-cell-tmtBaru="props">
            <q-td :props="props" class="text-center">{{ formatDate(props.row.tmtBaru) }}</q-td>
          </template>
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center">
              <q-btn flat round dense icon="visibility" color="primary" class="q-mr-sm">
                <q-tooltip>Lihat Detail</q-tooltip>
              </q-btn>
              <q-btn v-if="props.row.status === 'DRAFT'" flat round dense icon="send" color="secondary" @click="submitToKepsek(props.row.id)">
                <q-tooltip>Submit ke Kepsek</q-tooltip>
              </q-btn>
              <q-btn v-if="props.row.status === 'SUBMITTED_KEPSEK'" flat round dense icon="check_circle" color="positive" @click="approveUsulan(props.row.id)">
                <q-tooltip>Approve (Kepsek)</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- Dialog Create Usulan Sederhana -->
    <q-dialog v-model="createDialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section class="bg-primary text-white row items-center">
          <div class="text-h6">Buat Usulan Baru</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-input v-model="form.namaPegawai" label="Nama Pegawai" outlined dense readonly />
          <q-input v-model="form.jenisUsulan" label="Jenis Usulan" outlined dense readonly />
          <q-input v-model="form.catatan" label="Catatan KTU" type="textarea" outlined dense />
          <div class="row justify-end q-mt-md">
            <q-btn label="Batal" flat color="negative" v-close-popup class="q-mr-sm" />
            <q-btn label="Buat Usulan" color="primary" @click="submitCreateUsulan" :loading="isSubmitting" />
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { api } from '@/boot/axios'

const $q = useQuasar()
const isLoading = ref(false)
const isSubmitting = ref(false)
const saranUsulan = ref([])
const usulanList = ref([])
const createDialog = ref(false)
const form = ref({ ptkId: '', namaPegawai: '', jenisUsulan: '', catatan: '', golSekarangId: 0, golBaruId: 0 })

const columns = [
  { name: 'nomorUsulan', label: 'Nomor Usulan', field: 'nomorUsulan', align: 'left', sortable: true },
  { name: 'namaPegawai', label: 'Nama Pegawai', field: 'namaPegawai', align: 'left', sortable: true },
  { name: 'jenis', label: 'Jenis', field: 'jenis', align: 'center' },
  { name: 'golonganSekarang', label: 'Gol. Sekarang', field: 'golonganSekarang', align: 'center' },
  { name: 'golonganBaru', label: 'Gol. Baru', field: 'golonganBaru', align: 'center' },
  { name: 'tmtBaru', label: 'TMT Baru', align: 'center' },
  { name: 'status', label: 'Status', align: 'center' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const formatDate = (val) => {
  if (!val) return '-'
  return new Date(val).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric' })
}

const getStatusColor = (status) => {
  switch (status) {
    case 'AMAN': return 'green-6'
    case 'PERLU_PERSIAPAN': return 'yellow-8'
    case 'SEGERA': return 'orange-6'
    case 'URGENT': return 'red-6'
    case 'OVERDUE': return 'red-9'
    default: return 'grey-6'
  }
}

const getUsulanStatusColor = (status) => {
  switch (status) {
    case 'DRAFT': return 'grey-6'
    case 'LENGKAP': return 'blue-6'
    case 'SUBMITTED_KEPSEK': return 'yellow-8'
    case 'APPROVED_KEPSEK': return 'green-4'
    case 'SUBMITTED_DINAS': return 'indigo-6'
    case 'SK_TERBIT': return 'green-6'
    case 'DITOLAK': return 'red-6'
    case 'REVISI': return 'orange-6'
    default: return 'grey-6'
  }
}

const fetchData = async () => {
  isLoading.value = true
  try {
    const [resSaran, resList] = await Promise.all([
      api.get('/kepegawaian/usulan-kenaikan/saran'),
      api.get('/kepegawaian/usulan-kenaikan')
    ])
    saranUsulan.value = resSaran.data || []
    usulanList.value = resList.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal memuat data usulan' })
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

const openCreateDialog = (saran) => {
  form.value = { ptkId: saran.ptkId, namaPegawai: saran.namaPegawai, jenisUsulan: saran.jenisUsulan, catatan: '', golSekarangId: 0, golBaruId: 0 }
  createDialog.value = true
}

const submitCreateUsulan = async () => {
  isSubmitting.value = true
  try {
    await api.post('/kepegawaian/usulan-kenaikan/create', form.value)
    $q.notify({ type: 'positive', message: 'Usulan berhasil dibuat' })
    createDialog.value = false
    await fetchData()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Gagal membuat usulan' })
    console.error(err)
  } finally {
    isSubmitting.value = false
  }
}

const submitToKepsek = async (id) => {
  try {
    await api.put(`/kepegawaian/usulan-kenaikan/${id}/status`, { status: 'SUBMITTED_KEPSEK', catatan: 'Disubmit ke Kepala Sekolah' })
    $q.notify({ type: 'positive', message: 'Usulan disubmit ke Kepala Sekolah' })
    await fetchData()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal submit' })
    console.error(err)
  }
}

const approveUsulan = async (id) => {
  try {
    await api.put(`/kepegawaian/usulan-kenaikan/${id}/status`, { status: 'APPROVED_KEPSEK', catatan: 'Disetujui oleh Kepala Sekolah' })
    $q.notify({ type: 'positive', message: 'Usulan disetujui' })
    await fetchData()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal approve' })
    console.error(err)
  }
}

onMounted(fetchData)
</script>

<style scoped>
.border-left-warning {
  border-left: 4px solid #FFA726;
}
</style>
