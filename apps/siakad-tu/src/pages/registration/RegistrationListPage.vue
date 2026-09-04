<template>
  <q-page padding>
    <!-- Header -->
    <div class="row items-center justify-between q-mb-lg">
      <div>
        <div class="text-h5 text-weight-bold">Penerimaan Peserta Didik Baru (PPDB)</div>
        <div class="text-grey-7">Kelola pendaftaran siswa baru</div>
      </div>
      <q-btn color="primary" icon="add" label="Pendaftaran Baru" @click="openCreateDialog" />
    </div>

    <!-- Filters -->
    <q-card class="q-mb-md">
      <q-card-section>
        <div class="row q-gutter-md">
          <q-select v-model="filters.status" :options="statusOptions" label="Status" clearable class="col-12 col-sm-3"
            @update:model-value="applyFilters" />
          <q-input v-model="filters.search" label="Cari nama/NISN" clearable class="col-12 col-sm-6"
            @update:model-value="applyFilters">
            <template v-slot:append>
              <q-icon name="search" />
            </template>
          </q-input>
        </div>
      </q-card-section>
    </q-card>

    <!-- Table -->
    <q-card>
      <q-table :rows="registrationStore.registrations" :columns="columns" :loading="registrationStore.isLoading"
        :pagination="pagination" row-key="registrationId" @request="onRequest">
        <!-- Status Badge -->
        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-badge :color="getStatusColor(props.row.status)">
              {{ props.row.status }}
            </q-badge>
          </q-td>
        </template>

        <!-- Actions -->
        <template v-slot:body-cell-actions="props">
          <q-td :props="props">
            <q-btn flat dense round icon="visibility" color="primary" @click="viewDetail(props.row)">
              <q-tooltip>Lihat Detail</q-tooltip>
            </q-btn>

            <q-btn v-if="props.row.status === 'SUBMITTED'" flat dense round icon="check_circle" color="positive"
              @click="openApproveDialog(props.row)">
              <q-tooltip>Approve</q-tooltip>
            </q-btn>

            <q-btn v-if="props.row.status === 'SUBMITTED' || props.row.status === 'VERIFIED'" flat dense round
              icon="cancel" color="negative" @click="openRejectDialog(props.row)">
              <q-tooltip>Reject</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- Create Dialog -->
    <q-dialog v-model="showCreateDialog" persistent>
      <q-card style="width: 600px">
        <q-card-section>
          <div class="text-h6">Pendaftaran Siswa Baru</div>
        </q-card-section>

        <q-card-section>
          <q-form @submit="onCreateSubmit">
            <q-input v-model="createForm.nisn" label="NISN" :rules="[val => !!val || 'NISN wajib diisi']"
              class="q-mb-md" />
            <q-input v-model="createForm.nis" label="NIS" :rules="[val => !!val || 'NIS wajib diisi']"
              class="q-mb-md" />
            <q-input v-model="createForm.fullName" label="Nama Lengkap" :rules="[val => !!val || 'Nama wajib diisi']"
              class="q-mb-md" />
            <q-select v-model="createForm.gender" :options="genderOptions" label="Jenis Kelamin"
              :rules="[val => !!val || 'Jenis kelamin wajib dipilih']" class="q-mb-md" />
            <q-input v-model="createForm.birthDate" label="Tanggal Lahir" type="date"
              :rules="[val => !!val || 'Tanggal lahir wajib diisi']" class="q-mb-md" />

            <div class="row q-gutter-sm justify-end q-mt-md">
              <q-btn flat label="Batal" color="grey" @click="showCreateDialog = false" />
              <q-btn type="submit" label="Simpan Draft" color="primary" :loading="registrationStore.isLoading" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Approve Dialog -->
    <q-dialog v-model="showApproveDialog" persistent>
      <q-card style="width: 500px">
        <q-card-section>
          <div class="text-h6">Approve Pendaftaran</div>
        </q-card-section>

        <q-card-section>
          <q-form @submit="onApproveSubmit">
            <q-select v-model="approveForm.classId" :options="classOptions" label="Kelas Tujuan"
              :rules="[val => !!val || 'Kelas wajib dipilih']" class="q-mb-md" />
            <q-input v-model="approveForm.enrollmentDate" label="Tanggal Mulai" type="date"
              :rules="[val => !!val || 'Tanggal wajib diisi']" class="q-mb-md" />

            <div class="row q-gutter-sm justify-end q-mt-md">
              <q-btn flat label="Batal" color="grey" @click="showApproveDialog = false" />
              <q-btn type="submit" label="Approve" color="positive" :loading="registrationStore.isLoading" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Reject Dialog -->
    <q-dialog v-model="showRejectDialog" persistent>
      <q-card style="width: 500px">
        <q-card-section>
          <div class="text-h6">Reject Pendaftaran</div>
        </q-card-section>

        <q-card-section>
          <q-form @submit="onRejectSubmit">
            <q-input v-model="rejectForm.reason" label="Alasan Penolakan" type="textarea"
              :rules="[val => !!val || 'Alasan wajib diisi']" class="q-mb-md" />

            <div class="row q-gutter-sm justify-end q-mt-md">
              <q-btn flat label="Batal" color="grey" @click="showRejectDialog = false" />
              <q-btn type="submit" label="Reject" color="negative" :loading="registrationStore.isLoading" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Success/Error Notifications -->
    <q-dialog v-model="showNotification">
      <q-card>
        <q-card-section :class="notificationType === 'success' ? 'bg-positive' : 'bg-negative'">
          <div class="text-white">{{ notificationMessage }}</div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="OK" color="primary" @click="showNotification = false" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useRegistrationStore } from '@/stores/registration/registrationStore'

const router = useRouter()
const registrationStore = useRegistrationStore()

// === STATE ===
const filters = ref({ status: null, search: '' })
const pagination = ref({ page: 1, limit: 20 })

const showCreateDialog = ref(false)
const showApproveDialog = ref(false)
const showRejectDialog = ref(false)
const showNotification = ref(false)
const notificationMessage = ref('')
const notificationType = ref('success')

const createForm = ref({
  nisn: '',
  nis: '',
  fullName: '',
  gender: null,
  birthDate: '',
})

const approveForm = ref({
  registrationId: null,
  classId: null,
  enrollmentDate: new Date().toISOString().split('T')[0],
})

const rejectForm = ref({
  registrationId: null,
  reason: '',
})

// === OPTIONS ===
const statusOptions = [
  { label: 'Draft', value: 'DRAFT' },
  { label: 'Submitted', value: 'SUBMITTED' },
  { label: 'Verified', value: 'VERIFIED' },
  { label: 'Approved', value: 'APPROVED' },
  { label: 'Rejected', value: 'REJECTED' },
]

const genderOptions = [
  { label: 'Laki-laki', value: 'MALE' },
  { label: 'Perempuan', value: 'FEMALE' },
]

const classOptions = [
  { label: 'X RPL 1', value: 'class-1' },
  { label: 'X RPL 2', value: 'class-2' },
  { label: 'X TKJ 1', value: 'class-3' },
]

// === COLUMNS ===
const columns = [
  { name: 'nisn', label: 'NISN', field: 'nisn', align: 'left', sortable: true },
  { name: 'fullName', label: 'Nama Lengkap', field: 'fullName', align: 'left', sortable: true },
  { name: 'gender', label: 'L/P', field: 'gender', align: 'center' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'createdAt', label: 'Tanggal Daftar', field: row => new Date(row.createdAt).toLocaleDateString('id-ID'), align: 'center', sortable: true },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' },
]

// === COMPUTED ===
// const tablePagination = computed(() => ({
//   page: pagination.value.page,
//   rowsPerPage: pagination.value.limit,
//   rowsNumber: registrationStore.pagination.total,
// }))

// === METHODS ===
function getStatusColor(status) {
  const colors = {
    DRAFT: 'grey',
    SUBMITTED: 'orange',
    VERIFIED: 'blue',
    APPROVED: 'positive',
    REJECTED: 'negative',
  }
  return colors[status] || 'grey'
}

function applyFilters() {
  registrationStore.updateFilters(filters.value)
  registrationStore.fetchRegistrations({ page: 1 })
}

function onRequest(props) {
  const { page, rowsPerPage } = props.pagination

  pagination.value.page = page
  pagination.value.limit = rowsPerPage

  registrationStore.fetchRegistrations({
    page,
    limit: rowsPerPage,
  })
}

function viewDetail(registration) {
  router.push(`/registration/${registration.registrationId}`)
}

function openCreateDialog() {
  createForm.value = {
    nisn: '',
    nis: '',
    fullName: '',
    gender: null,
    birthDate: '',
  }
  showCreateDialog.value = true
}

function openApproveDialog(registration) {
  approveForm.value.registrationId = registration.registrationId
  showApproveDialog.value = true
}

function openRejectDialog(registration) {
  rejectForm.value.registrationId = registration.registrationId
  showRejectDialog.value = true
}

async function onCreateSubmit() {
  const result = await registrationStore.createDraft(createForm.value)

  if (result) {
    showCreateDialog.value = false
    notificationMessage.value = 'Draft pendaftaran berhasil dibuat'
    notificationType.value = 'success'
    showNotification.value = true
    registrationStore.fetchRegistrations()
  } else {
    notificationMessage.value = registrationStore.error?.message || 'Gagal membuat draft'
    notificationType.value = 'error'
    showNotification.value = true
  }
}

async function onApproveSubmit() {
  const result = await registrationStore.approveRegistration(
    approveForm.value.registrationId,
    {
      classId: approveForm.value.classId,
      enrollmentDate: approveForm.value.enrollmentDate,
    }
  )

  if (result) {
    showApproveDialog.value = false
    notificationMessage.value = 'Pendaftaran berhasil di-approve'
    notificationType.value = 'success'
    showNotification.value = true
    registrationStore.fetchRegistrations()
  } else {
    notificationMessage.value = registrationStore.error?.message || 'Gagal approve'
    notificationType.value = 'error'
    showNotification.value = true
  }
}

async function onRejectSubmit() {
  const result = await registrationStore.rejectRegistration(
    rejectForm.value.registrationId,
    rejectForm.value.reason
  )

  if (result) {
    showRejectDialog.value = false
    notificationMessage.value = 'Pendaftaran berhasil di-reject'
    notificationType.value = 'success'
    showNotification.value = true
    registrationStore.fetchRegistrations()
  } else {
    notificationMessage.value = registrationStore.error?.message || 'Gagal reject'
    notificationType.value = 'error'
    showNotification.value = true
  }
}

// === LIFECYCLE ===
onMounted(() => {
  registrationStore.fetchRegistrations()
})
</script>
