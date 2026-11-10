<template>
  <q-page class="q-pa-md">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">SK Pembagian Tugas</h1>
        <p class="text-caption text-grey-7 q-mb-none">Manajemen Surat Keputusan Beban Mengajar</p>
      </div>
      <q-btn color="primary" icon="add" label="Buat SK Baru" @click="$router.push('/kepegawaian/sk-pembagian-tugas/create')" />
    </div>

    <q-card class="shadow-1">
      <q-card-section class="row items-center q-pb-none">
        <q-select v-model="filterStatus" :options="['Semua', 'DRAFT', 'SUBMITTED', 'ISSUED', 'REJECTED', 'REVISION']" label="Filter Status" outlined dense class="col-12 col-md-3 q-mr-md" />
        <q-btn color="secondary" icon="refresh" label="Refresh" @click="fetchData" :loading="isLoading" flat />
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
              <!-- Tombol Approve/Reject (Hanya jika SUBMITTED dan role kepsek/admin) -->
              <div v-if="canApprove && props.row.status === 'SUBMITTED'">
                <q-btn flat round dense icon="check_circle" color="positive" @click="approveSK(props.row.id)" class="q-mr-sm">
                  <q-tooltip>Setujui & Terbitkan</q-tooltip>
                </q-btn>
                <q-btn flat round dense icon="cancel" color="negative" @click="openRejectDialog(props.row.id)" class="q-mr-sm">
                  <q-tooltip>Tolak</q-tooltip>
                </q-btn>
              </div>

              <!-- Tombol Submit (Hanya jika DRAFT/REVISION) -->
              <q-btn v-if="props.row.status === 'DRAFT' || props.row.status === 'REVISION'" flat round dense icon="send" color="primary" @click="submitSK(props.row.id)" class="q-mr-sm">
                <q-tooltip>Ajukan Approval</q-tooltip>
              </q-btn>

              <!-- Tombol Download PDF & Excel (Jika sudah ISSUED) -->
              <template v-if="props.row.status === 'ISSUED'">
                <q-btn flat round dense icon="picture_as_pdf" color="red-6" type="a" :href="props.row.finalPdfUrl" target="_blank" class="q-mr-sm">
                  <q-tooltip>Download PDF SK</q-tooltip>
                </q-btn>
                <q-btn flat round dense icon="table_view" color="green-6" @click="exportExcel(props.row.id)">
                  <q-tooltip>Export Rekap Excel</q-tooltip>
                </q-btn>
              </template>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- Dialog Reject -->
    <q-dialog v-model="rejectDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="bg-negative text-white"><div class="text-h6">Alasan Penolakan</div></q-card-section>
        <q-card-section>
          <q-input v-model="rejectReason" label="Alasan Penolakan" type="textarea" outlined autofocus :rules="[v => !!v || 'Wajib diisi']" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn flat label="Tolak SK" color="negative" @click="executeReject" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { api } from '@/boot/axios'
import { useAuthStore } from '@/stores/authStore'

const $q = useQuasar()
const authStore = useAuthStore()

const skList = ref([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const filterStatus = ref('Semua')
const rejectDialog = ref(false)
const rejectReason = ref('')
let selectedRejectId = null

const canApprove = computed(() => {
  const role = authStore.user?.role
  return role === 'admin' || role === 'kepsek'
})

const columns = [
  { name: 'nomorSk', label: 'Nomor SK', field: 'nomorSk', align: 'left', sortable: true },
  { name: 'tahunAjaran', label: 'Tahun Ajaran', field: 'tahunAjaran', align: 'center' },
  { name: 'semester', label: 'Semester', field: 'semester', align: 'center' },
  { name: 'status', label: 'Status', align: 'center' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const filteredData = computed(() => {
  if (filterStatus.value === 'Semua') return skList.value
  return skList.value.filter(sk => sk.status === filterStatus.value)
})

const getStatusColor = (status) => {
  switch (status) {
    case 'DRAFT': return 'grey-6'
    case 'SUBMITTED': return 'orange-6'
    case 'ISSUED': return 'green-6'
    case 'REJECTED': return 'red-6'
    case 'REVISION': return 'blue-6'
    default: return 'grey-6'
  }
}

const fetchData = async () => {
  isLoading.value = true
  try {
    const res = await api.get('/kepegawaian/sk-pembagian-tugas')
    skList.value = res.data || []
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Gagal memuat data SK' })
    console.error(e)
  } finally {
    isLoading.value = false
  }
}

const submitSK = async (id) => {
  try {
    await api.put(`/kepegawaian/sk-pembagian-tugas/${id}/submit`)
    $q.notify({ type: 'positive', message: 'SK berhasil diajukan' })
    await fetchData()
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Gagal mengajukan SK' })
    console.error(e)
  }
}

const approveSK = (id) => {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Setujui dan terbitkan SK ini? PDF akan di-generate otomatis.',
    cancel: true
  }).onOk(async () => {
    isSubmitting.value = true
    try {
      await api.put(`/kepegawaian/sk-pembagian-tugas/${id}/approve`)
      $q.notify({ type: 'positive', message: 'SK berhasil diterbitkan' })
      await fetchData()
    } catch (e) {
      $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal menerbitkan SK' })
      console.error(e)
    } finally {
      isSubmitting.value = false
    }
  })
}

const openRejectDialog = (id) => {
  selectedRejectId = id
  rejectReason.value = ''
  rejectDialog.value = true
}

const executeReject = async () => {
  if (!rejectReason.value.trim()) return
  isSubmitting.value = true
  try {
    await api.put(`/kepegawaian/sk-pembagian-tugas/${selectedRejectId}/reject`, { reason: rejectReason.value })
    $q.notify({ type: 'negative', message: 'SK ditolak' })
    rejectDialog.value = false
    await fetchData()
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Gagal menolak SK' })
    console.error(e)
  } finally {
    isSubmitting.value = false
  }
}

const exportExcel = async (id) => {
  try {
    const res = await api.get(`/kepegawaian/sk-pembagian-tugas/export/${id}`)
    window.open(res.data.url, '_blank')
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Gagal export Excel' })
    console.error(e)
  }
}

onMounted(fetchData)
</script>
