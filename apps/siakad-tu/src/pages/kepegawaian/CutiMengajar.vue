<template>
  <q-page class="q-pa-md">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Pengajuan Cuti / SPT</h1>
        <p class="text-caption text-grey-7 q-mb-none">Manajemen izin dan cuti Guru & Tenaga Kependidikan</p>
      </div>
      <q-btn color="primary" icon="add" label="Ajukan Cuti/SPT" @click="openFormDialog" />
    </div>

    <q-card class="shadow-1">
      <q-card-section class="row items-center q-pb-none">
        <q-select v-model="filterStatus" :options="['semua', 'pending', 'disetujui', 'ditolak']" label="Filter Status" dense
          outlined class="q-mr-md" style="min-width: 200px" />
        <q-btn color="secondary" icon="refresh" label="Refresh" @click="fetchLeaves" :loading="isLoading" flat />
      </q-card-section>

      <q-card-section>
        <q-table :rows="filteredCuti" :columns="columns" row-key="id" flat bordered dense :loading="isLoading">
          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center">
              <q-badge :color="props.row.status === 'pending' ? 'warning' : (props.row.status === 'disetujui' ? 'positive' : 'negative')">
                {{ props.row.status.toUpperCase() }}
              </q-badge>
            </q-td>
          </template>
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center">
              <div v-if="(authStore.user?.role === 'kepsek' || authStore.user?.role === 'admin') && props.row.status === 'pending'">
                <q-btn flat round dense icon="check_circle" color="positive" @click="approve(props.row.id)" class="q-mr-sm">
                  <q-tooltip>Setujui</q-tooltip>
                </q-btn>
                <q-btn flat round dense icon="cancel" color="negative" @click="openRejectDialog(props.row.id)">
                  <q-tooltip>Tolak</q-tooltip>
                </q-btn>
              </div>
              <div v-if="authStore.user?.role === 'kepsek' || authStore.user?.role === 'admin'">
                <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row.id)">
                  <q-tooltip>Hapus</q-tooltip>
                </q-btn>
              </div>
            </q-td>
          </template>
          <template v-slot:body-cell-alasan="props">
            <q-td :props="props">
              <div class="ellipsis" style="max-width: 200px">{{ props.row.alasan || '-' }}</div>
              <q-tooltip v-if="props.row.alasan">{{ props.row.alasan }}</q-tooltip>
            </q-td>
          </template>
          <template v-slot:body-cell-catatan="props">
            <q-td :props="props">
              <div class="ellipsis" style="max-width: 200px">{{ props.row.catatanApproval || '-' }}</div>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- Dialog form pengajuan -->
    <q-dialog v-model="formDialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section class="row items-center bg-primary text-white">
          <div class="text-h6">Form Pengajuan Cuti / SPT</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-form @submit="submitCuti">
            <q-select v-model="form.ptkId" :options="guruOptions" label="Pilih Guru" option-value="id"
              option-label="nama" emit-value map-options outlined dense :rules="[v => !!v || 'Wajib dipilih']" />
            <div class="row q-col-gutter-sm q-mt-sm">
              <div class="col-6">
                <q-input v-model="form.tanggalMulai" label="Tanggal Mulai" type="date" outlined dense :rules="[v => !!v || 'Wajib diisi']" />
              </div>
              <div class="col-6">
                <q-input v-model="form.tanggalSelesai" label="Tanggal Selesai" type="date" outlined dense :rules="[v => !!v || 'Wajib diisi']" />
              </div>
            </div>
            <q-select v-model="form.jenisCuti" :options="['Cuti Tahunan', 'Cuti Sakit', 'SPT (Surat Perintah Tugas)', 'Izin']" label="Jenis Cuti" outlined dense class="q-mt-sm" />
            <q-input v-model="form.alasan" label="Alasan" type="textarea" rows="2" outlined dense class="q-mt-sm" />
            <div class="row justify-end q-mt-md">
              <q-btn label="Batal" flat color="negative" v-close-popup class="q-mr-sm" />
              <q-btn label="Ajukan" type="submit" color="primary" :loading="isSubmitting" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog catatan penolakan -->
    <q-dialog v-model="rejectDialog" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="bg-negative text-white">
          <div class="text-h6">Alasan Penolakan</div>
        </q-card-section>
        <q-card-section>
          <q-input v-model="rejectNote" label="Catatan / Alasan Penolakan" type="textarea" outlined autofocus :rules="[v => !!v || 'Wajib diisi']" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn flat label="Tolak" color="negative" @click="executeReject" :loading="isSubmitting" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { api } from '@/boot/axios'
import { useGuruStore } from '@/stores/kepegawaian/guru'
import { useAuthStore } from '@/stores/authStore'

const $q = useQuasar()
const guruStore = useGuruStore()
const authStore = useAuthStore()

const filterStatus = ref('semua')
const formDialog = ref(false)
const rejectDialog = ref(false)
const currentRejectId = ref(null)
const rejectNote = ref('')
const isLoading = ref(false)
const isSubmitting = ref(false)
const leaves = ref([])

const form = ref({
  ptkId: null,
  tanggalMulai: '',
  tanggalSelesai: '',
  jenisCuti: 'Cuti Tahunan',
  alasan: ''
})

const columns = [
  { name: 'teacherName', label: 'Nama Guru', field: 'teacherName', align: 'left', sortable: true },
  { name: 'tanggalMulai', label: 'Mulai', field: 'tanggalMulai', align: 'left', format: val => new Date(val).toLocaleDateString('id-ID') },
  { name: 'tanggalSelesai', label: 'Selesai', field: 'tanggalSelesai', align: 'left', format: val => new Date(val).toLocaleDateString('id-ID') },
  { name: 'jenisCuti', label: 'Jenis', field: 'jenisCuti', align: 'left' },
  { name: 'alasan', label: 'Alasan', field: 'alasan', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'catatan', label: 'Catatan', field: 'catatanApproval', align: 'left' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const guruOptions = computed(() => guruStore.list.map(g => ({ id: g.id, nama: g.name })))

const filteredCuti = computed(() => {
  if (filterStatus.value === 'semua') return leaves.value
  return leaves.value.filter(c => c.status === filterStatus.value)
})

const fetchLeaves = async () => {
  isLoading.value = true
  try {
    const res = await api.get('/teachers/leaves')
    leaves.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal memuat data cuti' })
  } finally {
    isLoading.value = false
  }
}

const openFormDialog = () => {
  form.value = { ptkId: null, tanggalMulai: '', tanggalSelesai: '', jenisCuti: 'Cuti Tahunan', alasan: '' }
  formDialog.value = true
}

const submitCuti = async () => {
  isSubmitting.value = true
  try {
    await api.post('/teachers/leaves', form.value)
    $q.notify({ type: 'positive', message: 'Pengajuan cuti berhasil dikirim' })
    formDialog.value = false
    await fetchLeaves()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Gagal mengirim pengajuan' })
  } finally {
    isSubmitting.value = false
  }
}

const approve = (id) => {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Setujui pengajuan cuti ini?',
    cancel: true
  }).onOk(async () => {
    try {
      await api.put(`/teachers/leaves/${id}/status`, { status: 'disetujui', note: 'Disetujui oleh Kepala Sekolah/Admin' })
      $q.notify({ type: 'positive', message: 'Pengajuan disetujui' })
      await fetchLeaves()
    } catch (err) {
      $q.notify({ type: 'negative', message: 'Gagal menyetujui pengajuan' })
    }
  })
}

const openRejectDialog = (id) => {
  currentRejectId.value = id
  rejectNote.value = ''
  rejectDialog.value = true
}

const executeReject = async () => {
  if (!rejectNote.value.trim()) return
  isSubmitting.value = true
  try {
    await api.put(`/teachers/leaves/${currentRejectId.value}/status`, { status: 'ditolak', note: rejectNote.value })
    $q.notify({ type: 'negative', message: 'Pengajuan ditolak' })
    rejectDialog.value = false
    await fetchLeaves()
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal menolak pengajuan' })
  } finally {
    isSubmitting.value = false
  }
}

const confirmDelete = (id) => {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Hapus data pengajuan cuti ini?',
    cancel: true
  }).onOk(async () => {
    try {
      await api.delete(`/teachers/leaves/${id}`)
      $q.notify({ type: 'positive', message: 'Data berhasil dihapus' })
      await fetchLeaves()
    } catch (err) {
      $q.notify({ type: 'negative', message: 'Gagal menghapus data' })
    }
  })
}

onMounted(async () => {
  await guruStore.fetchTeachers()
  await fetchLeaves()
})
</script>
