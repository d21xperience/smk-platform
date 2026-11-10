<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Laporan Kepegawaian Berkala</h1>
        <p class="text-caption text-grey-7 q-mb-none">Laporan Triwulanan ke Dinas & Yayasan</p>
      </div>
      <div class="row q-gutter-sm">
        <q-btn color="secondary" icon="refresh" label="Refresh" @click="fetchData" :loading="isLoading" flat />
        <q-btn color="primary" icon="add" label="Generate Laporan Baru" @click="openGenerateDialog" />
      </div>
    </div>

    <q-card class="shadow-1">
      <q-card-section class="row items-center q-pb-none">
        <q-select v-model="filterStatus" :options="['Semua', 'DRAFT', 'REVIEWED', 'APPROVED', 'SENT', 'ARCHIVED']" label="Filter Status" outlined dense class="col-12 col-md-3" />
      </q-card-section>

      <q-card-section>
        <q-table :rows="filteredData" :columns="columns" row-key="id" flat bordered dense :loading="isLoading">
          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center">
              <q-badge :color="getStatusColor(props.row.status)" class="q-pa-sm">
                {{ props.row.status }}
              </q-badge>
            </q-td>
          </template>
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center">
              <!-- Tombol Generate Dokumen (Jika DRAFT dan belum ada PDF) -->
              <q-btn v-if="props.row.status === 'DRAFT'" flat round dense icon="picture_as_pdf" color="primary" @click="generateDocs(props.row.id)" class="q-mr-sm">
                <q-tooltip>Generate PDF/Excel</q-tooltip>
              </q-btn>

              <!-- Tombol Review (Jika DRAFT) -->
              <q-btn v-if="props.row.status === 'DRAFT'" flat round dense icon="rate_review" color="orange-8" @click="openReviewDialog(props.row.id)" class="q-mr-sm">
                <q-tooltip>Review KTU</q-tooltip>
              </q-btn>

              <!-- Tombol Approve (Jika REVIEWED) -->
              <q-btn v-if="props.row.status === 'REVIEWED'" flat round dense icon="check_circle" color="positive" @click="openApproveDialog(props.row.id)" class="q-mr-sm">
                <q-tooltip>Approve Kepsek</q-tooltip>
              </q-btn>

              <!-- Tombol Send (Jika APPROVED) -->
              <q-btn v-if="props.row.status === 'APPROVED'" flat round dense icon="send" color="secondary" @click="openSendDialog(props.row.id)" class="q-mr-sm">
                <q-tooltip>Kirim ke Dinas</q-tooltip>
              </q-btn>

              <!-- Tombol Download PDF -->
              <q-btn v-if="props.row.pdfUrl" flat round dense icon="download" color="primary" type="a" :href="props.row.pdfUrl" target="_blank">
                <q-tooltip>Download PDF</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- Dialog Generate -->
    <q-dialog v-model="generateDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="bg-primary text-white"><div class="text-h6">Generate Laporan Baru</div></q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-input v-model.number="generateForm.tahun" type="number" label="Tahun" outlined dense />
          <q-select v-model="generateForm.triwulan" :options="[1, 2, 3, 4]" label="Triwulan" outlined dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn flat label="Generate" color="primary" @click="submitGenerate" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Review -->
    <q-dialog v-model="reviewDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="bg-orange-8 text-white"><div class="text-h6">Review Laporan</div></q-card-section>
        <q-card-section>
          <q-input v-model="actionForm.notes" label="Catatan Review" type="textarea" outlined dense :rules="[v => !!v || 'Wajib diisi']" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn flat label="Simpan Review" color="orange-8" @click="submitReview" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Approve -->
    <q-dialog v-model="approveDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="bg-positive text-white"><div class="text-h6">Approve Laporan</div></q-card-section>
        <q-card-section>
          <q-input v-model="actionForm.notes" label="Catatan Persetujuan" type="textarea" outlined dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn flat label="Setujui" color="positive" @click="submitApprove" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Send -->
    <q-dialog v-model="sendDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="bg-secondary text-white"><div class="text-h6">Kirim Laporan</div></q-card-section>
        <q-card-section>
          <q-input v-model="actionForm.email" label="Email Tujuan (Dinas/Yayasan)" outlined dense :rules="[v => !!v || 'Wajib diisi']" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn flat label="Kirim Email" color="secondary" @click="submitSend" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { api } from '@/boot/axios'

const $q = useQuasar()
const laporanList = ref([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const filterStatus = ref('Semua')

const generateDialog = ref(false)
const reviewDialog = ref(false)
const approveDialog = ref(false)
const sendDialog = ref(false)

const generateForm = ref({ tahun: new Date().getFullYear(), triwulan: 1 })
const actionForm = ref({ id: null, notes: '', email: '' })

const columns = [
  { name: 'tahun', label: 'Tahun', field: 'tahun', align: 'center', sortable: true },
  { name: 'triwulan', label: 'Triwulan', field: 'triwulan', align: 'center' },
  { name: 'periode', label: 'Periode', align: 'left' },
  { name: 'status', label: 'Status', align: 'center' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const filteredData = computed(() => {
  if (filterStatus.value === 'Semua') return laporanList.value
  return laporanList.value.filter(l => l.status === filterStatus.value)
})

const getStatusColor = (status) => {
  switch (status) {
    case 'DRAFT': return 'grey-6'
    case 'REVIEWED': return 'orange-8'
    case 'APPROVED': return 'blue-6'
    case 'SENT': return 'green-6'
    case 'ARCHIVED': return 'purple-6'
    default: return 'grey-6'
  }
}

const fetchData = async () => {
  isLoading.value = true
  try {
    const res = await api.get('/kepegawaian/laporan')
    laporanList.value = res.data || []
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Gagal memuat data laporan' })
    console.error(e)
  } finally {
    isLoading.value = false
  }
}

const openGenerateDialog = () => {
  generateForm.value = { tahun: new Date().getFullYear(), triwulan: 1 }
  generateDialog.value = true
}

const submitGenerate = async () => {
  isSubmitting.value = true
  try {
    await api.post('/kepegawaian/laporan/generate', generateForm.value)
    $q.notify({ type: 'positive', message: 'Draft laporan berhasil dibuat' })
    generateDialog.value = false
    await fetchData()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal generate' })
    console.error(e)
  } finally {
    isSubmitting.value = false
  }
}

const generateDocs = async (id) => {
  try {
    await api.post(`/kepegawaian/laporan/${id}/generate-documents`)
    $q.notify({ type: 'positive', message: 'Dokumen PDF & Excel berhasil digenerate' })
    await fetchData()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal generate dokumen' })
    console.error(e)
  }
}

const openReviewDialog = (id) => {
  actionForm.value = { id, notes: '', email: '' }
  reviewDialog.value = true
}

const submitReview = async () => {
  if (!actionForm.value.notes) return
  isSubmitting.value = true
  try {
    await api.post(`/kepegawaian/laporan/${actionForm.value.id}/review`, { notes: actionForm.value.notes })
    $q.notify({ type: 'positive', message: 'Laporan berhasil direview' })
    reviewDialog.value = false
    await fetchData()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal review' })
    console.error(e)
  } finally {
    isSubmitting.value = false
  }
}

const openApproveDialog = (id) => {
  actionForm.value = { id, notes: '', email: '' }
  approveDialog.value = true
}

const submitApprove = async () => {
  isSubmitting.value = true
  try {
    await api.post(`/kepegawaian/laporan/${actionForm.value.id}/approve`, { notes: actionForm.value.notes })
    $q.notify({ type: 'positive', message: 'Laporan berhasil disetujui' })
    approveDialog.value = false
    await fetchData()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal approve' })
    console.error(e)
  } finally {
    isSubmitting.value = false
  }
}

const openSendDialog = (id) => {
  actionForm.value = { id, notes: '', email: '' }
  sendDialog.value = true
}

const submitSend = async () => {
  if (!actionForm.value.email) return
  isSubmitting.value = true
  try {
    await api.post(`/kepegawaian/laporan/${actionForm.value.id}/send`, { email: actionForm.value.email })
    $q.notify({ type: 'positive', message: 'Laporan berhasil dikirim' })
    sendDialog.value = false
    await fetchData()
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal kirim email (cek SMTP)' })
    console.error(e)
  } finally {
    isSubmitting.value = false
  }
}

onMounted(fetchData)
</script>
