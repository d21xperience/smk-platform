<!-- FILE: src/pages/homeroom-communication/HomeroomCommunicationPage.vue -->
<!-- STATUS: NEW -->
<!-- STATUS IMPLEMENTASI: COMPLETE -->

<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div class="row items-center">
        <div class="text-h5">Komunikasi Orang Tua</div>
        <q-badge v-if="className" color="primary" class="q-ml-md" style="font-size: 14px; padding: 4px 12px;">
          Kelas {{ className }}
        </q-badge>
      </div>
      <q-btn color="primary" icon="add" label="Buat Log" @click="openCreateDialog" />
    </div>

    <q-banner v-if="error" class="bg-negative text-white q-mb-md">
      {{ error }}
    </q-banner>

    <div v-if="loading" class="flex flex-center" style="min-height: 200px;">
      <q-spinner-dots size="60px" color="primary" />
    </div>

    <template v-if="!loading">
      <div v-if="summary" class="row q-col-gutter-md q-mb-md">
        <div class="col-6 col-sm-3">
          <q-card flat bordered>
            <q-card-section class="text-center q-py-sm">
              <div class="text-h4 text-primary">{{ summary.totalCommunications }}</div>
              <div class="text-caption text-grey-7">Total Log</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-6 col-sm-3">
          <q-card flat bordered>
            <q-card-section class="text-center q-py-sm">
              <div class="text-h4 text-warning">{{ summary.pendingCount }}</div>
              <div class="text-caption text-grey-7">Perlu Follow Up</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-6 col-sm-3">
          <q-card flat bordered>
            <q-card-section class="text-center q-py-sm">
              <div class="text-h4 text-info">{{ summary.followedUpCount }}</div>
              <div class="text-caption text-grey-7">Sudah Ditindaklanjuti</div>
            </q-card-section>
          </q-card>
        </div>
        <div class="col-6 col-sm-3">
          <q-card flat bordered>
            <q-card-section class="text-center q-py-sm">
              <div class="text-h4 text-positive">{{ summary.closedCount }}</div>
              <div class="text-caption text-grey-7">Selesai</div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <q-table flat bordered :rows="communications" :columns="columns" row-key="communicationId"
        :pagination="{ rowsPerPage: 10, sortBy: 'communicationDate', ascending: false }">
        <template #body-cell-category="props">
          <q-td :props="props" class="text-center">
            <q-badge :color="getCategoryColor(props.row.category)">
              <q-icon :name="getCategoryIcon(props.row.category)" size="xs" class="q-mr-xs" />
              {{ getCategoryLabel(props.row.category) }}
            </q-badge>
          </q-td>
        </template>

        <template #body-cell-direction="props">
          <q-td :props="props" class="text-center">
            <q-badge :color="getDirectionColor(props.row.direction)" outline>
              <q-icon :name="getDirectionIcon(props.row.direction)" size="xs" class="q-mr-xs" />
              {{ getDirectionLabel(props.row.direction) }}
            </q-badge>
          </q-td>
        </template>

        <template #body-cell-status="props">
          <q-td :props="props" class="text-center">
            <q-badge :color="getStatusColor(props.row.status)">
              {{ getStatusLabel(props.row.status) }}
            </q-badge>
          </q-td>
        </template>

        <template #body-cell-actions="props">
          <q-td :props="props" class="text-center">
            <q-btn v-if="props.row.status === 'PENDING_FOLLOW_UP'" flat dense icon="check_circle" color="info" size="sm"
              title="Tandai sudah ditindaklanjuti" :loading="updating"
              @click="openUpdateDialog(props.row, 'FOLLOWED_UP')" />
            <q-btn v-if="props.row.status === 'PENDING_FOLLOW_UP' || props.row.status === 'FOLLOWED_UP'" flat dense
              icon="close" color="positive" size="sm" title="Tandai selesai" :loading="updating"
              @click="openUpdateDialog(props.row, 'CLOSED')" />
            <q-btn flat dense icon="visibility" color="grey" size="sm" title="Lihat detail"
              @click="viewCommunication(props.row)" />
          </q-td>
        </template>
      </q-table>
    </template>

    <q-dialog v-model="showCreateDialog" persistent>
      <q-card style="min-width: 520px;">
        <q-card-section>
          <div class="text-h6">Buat Log Komunikasi</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-select v-model="formData.studentId" :options="studentOptions" label="Siswa" outlined dense emit-value
            map-options class="q-mb-md" @update:model-value="onStudentSelected" />

          <div class="row q-col-gutter-md q-mb-md">
            <div class="col-6">
              <q-select v-model="formData.category" :options="categoryOptions" label="Kategori" outlined dense
                emit-value map-options />
            </div>
            <div class="col-6">
              <q-select v-model="formData.direction" :options="directionOptions" label="Arah Komunikasi" outlined dense
                emit-value map-options />
            </div>
          </div>

          <q-input v-model="formData.title" label="Judul" outlined dense class="q-mb-md" :maxlength="200" counter />

          <q-input v-model="formData.summary" label="Ringkasan Komunikasi" type="textarea" outlined dense rows="4"
            class="q-mb-md" :maxlength="2000" counter />

          <q-input v-model="formData.communicationDate" label="Tanggal Komunikasi" outlined dense mask="####-##-##"
            hint="Format: YYYY-MM-DD">
            <template #append>
              <q-icon name="event" class="cursor-pointer">
                <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                  <q-date v-model="formData.communicationDate" mask="YYYY-MM-DD">
                    <div class="row items-center justify-end">
                      <q-btn v-close-popup label="Tutup" color="primary" flat />
                    </div>
                  </q-date>
                </q-popup-proxy>
              </q-icon>
            </template>
          </q-input>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" color="grey" @click="showCreateDialog = false" />
          <q-btn label="Simpan Log" color="primary" :loading="creating" @click="handleCreateCommunication" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog v-model="showUpdateDialog" persistent>
      <q-card v-if="updateData.communication" style="min-width: 450px;">
        <q-card-section>
          <div class="text-h6">
            {{ updateData.newStatus === 'FOLLOWED_UP' ? 'Tandai Sudah Ditindaklanjuti' : 'Tandai Selesai' }}
          </div>
          <div class="text-caption text-grey q-mt-sm">
            {{ updateData.communication.title }} — {{ updateData.communication.studentName }}
          </div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-input v-model="updateData.followUpNote" label="Catatan Tindak Lanjut" type="textarea" outlined dense
            rows="3" :maxlength="1000" counter />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" color="grey" @click="showUpdateDialog = false" />
          <q-btn :label="updateData.newStatus === 'FOLLOWED_UP' ? 'Tandai Follow Up' : 'Tandai Selesai'"
            :color="updateData.newStatus === 'FOLLOWED_UP' ? 'info' : 'positive'" :loading="updating"
            @click="handleUpdateStatus" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog v-model="showDetailDialog">
      <q-card v-if="selectedCommunication" style="min-width: 480px;">
        <q-card-section>
          <div class="text-h6">{{ selectedCommunication.title }}</div>
          <div class="text-caption text-grey">
            {{ selectedCommunication.studentName }} — {{ selectedCommunication.communicationDate }}
          </div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <div class="row q-col-gutter-sm q-mb-md">
            <div class="col-4">
              <q-badge :color="getCategoryColor(selectedCommunication.category)">
                {{ getCategoryLabel(selectedCommunication.category) }}
              </q-badge>
            </div>
            <div class="col-4">
              <q-badge :color="getDirectionColor(selectedCommunication.direction)" outline>
                {{ getDirectionLabel(selectedCommunication.direction) }}
              </q-badge>
            </div>
            <div class="col-4">
              <q-badge :color="getStatusColor(selectedCommunication.status)">
                {{ getStatusLabel(selectedCommunication.status) }}
              </q-badge>
            </div>
          </div>

          <div class="text-subtitle2 q-mb-xs">Ringkasan:</div>
          <p class="q-mb-md">{{ selectedCommunication.summary }}</p>

          <template v-if="selectedCommunication.followUpNote">
            <div class="text-subtitle2 q-mb-xs">Catatan Tindak Lanjut:</div>
            <p>{{ selectedCommunication.followUpNote }}</p>
          </template>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Tutup" color="grey" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import { useQuasar } from 'quasar'
import { useHomeroomCommunication } from '@/composables/useHomeroomCommunication.js'
import { useHomeroomStudent } from '@/composables/useHomeroomStudent.js'

const $q = useQuasar()

const {
  communications,
  summary,
  loading,
  creating,
  updating,
  error,
  loadCommunications,
  createCommunication,
  updateStatus,
} = useHomeroomCommunication()

const {
  students,
  className,
  loadStudents,
} = useHomeroomStudent()

const showCreateDialog = ref(false)
const showUpdateDialog = ref(false)
const showDetailDialog = ref(false)
const selectedCommunication = ref(null)

const formData = ref({
  studentId: null,
  studentName: '',
  category: 'PHONE_CALL',
  direction: 'OUTBOUND',
  title: '',
  summary: '',
  communicationDate: '',
})

const updateData = ref({
  communication: null,
  newStatus: '',
  followUpNote: '',
})

const columns = ref([
  { name: 'communicationDate', label: 'Tanggal', field: 'communicationDate', align: 'center', sortable: true },
  { name: 'studentName', label: 'Siswa', field: 'studentName', align: 'left', sortable: true },
  { name: 'category', label: 'Kategori', field: 'category', align: 'center' },
  { name: 'direction', label: 'Arah', field: 'direction', align: 'center' },
  { name: 'title', label: 'Judul', field: 'title', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'center', sortable: true },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' },
])

const categoryOptions = ref([
  { label: 'Telepon', value: 'PHONE_CALL' },
  { label: 'WhatsApp', value: 'WHATSAPP' },
  { label: 'Pertemuan', value: 'MEETING' },
  { label: 'Kunjungan Rumah', value: 'HOME_VISIT' },
  { label: 'Lainnya', value: 'OTHER' },
])

const directionOptions = ref([
  { label: 'Masuk (dari Ortu)', value: 'INBOUND' },
  { label: 'Keluar (ke Ortu)', value: 'OUTBOUND' },
  { label: 'Dua Arah', value: 'BIDIRECTIONAL' },
])

const studentOptions = computed(() => {
  return students.value.map(s => ({
    label: `${s.seatNumber}. ${s.fullName}`,
    value: s.id,
  }))
})

function onStudentSelected(studentId) {
  const found = students.value.find(s => s.id === studentId)
  if (found) {
    formData.value.studentName = found.fullName
  }
}

function openCreateDialog() {
  formData.value = {
    studentId: null,
    studentName: '',
    category: 'PHONE_CALL',
    direction: 'OUTBOUND',
    title: '',
    summary: '',
    communicationDate: new Date().toISOString().split('T')[0],
  }
  showCreateDialog.value = true
}

async function handleCreateCommunication() {
  try {
    await createCommunication({ ...formData.value })
    showCreateDialog.value = false
    $q.notify({
      type: 'positive',
      message: 'Log komunikasi berhasil dibuat.',
      position: 'top',
      timeout: 3000,
    })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message || 'Gagal membuat log komunikasi.',
      position: 'top',
      timeout: 3000,
    })
  }
}

function openUpdateDialog(communication, newStatus) {
  updateData.value = {
    communication: communication,
    newStatus: newStatus,
    followUpNote: communication.followUpNote || '',
  }
  showUpdateDialog.value = true
}

async function handleUpdateStatus() {
  try {
    await updateStatus(
      updateData.value.communication,
      updateData.value.newStatus,
      updateData.value.followUpNote,
    )
    showUpdateDialog.value = false
    $q.notify({
      type: 'positive',
      message: 'Status log komunikasi berhasil diupdate.',
      position: 'top',
      timeout: 3000,
    })
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message || 'Gagal update status.',
      position: 'top',
      timeout: 3000,
    })
  }
}

function viewCommunication(communication) {
  selectedCommunication.value = communication
  showDetailDialog.value = true
}

function getCategoryIcon(category) {
  switch (category) {
    case 'PHONE_CALL': return 'phone'
    case 'WHATSAPP': return 'chat'
    case 'MEETING': return 'groups'
    case 'HOME_VISIT': return 'home'
    case 'OTHER': return 'more_horiz'
    default: return 'more_horiz'
  }
}

function getCategoryColor(category) {
  switch (category) {
    case 'PHONE_CALL': return 'primary'
    case 'WHATSAPP': return 'positive'
    case 'MEETING': return 'secondary'
    case 'HOME_VISIT': return 'deep-orange'
    case 'OTHER': return 'grey'
    default: return 'grey'
  }
}

function getCategoryLabel(category) {
  switch (category) {
    case 'PHONE_CALL': return 'Telepon'
    case 'WHATSAPP': return 'WhatsApp'
    case 'MEETING': return 'Pertemuan'
    case 'HOME_VISIT': return 'Kunjungan'
    case 'OTHER': return 'Lainnya'
    default: return category
  }
}

function getDirectionIcon(direction) {
  switch (direction) {
    case 'INBOUND': return 'call_received'
    case 'OUTBOUND': return 'call_made'
    case 'BIDIRECTIONAL': return 'swap_vert'
    default: return 'swap_vert'
  }
}

function getDirectionColor(direction) {
  switch (direction) {
    case 'INBOUND': return 'info'
    case 'OUTBOUND': return 'primary'
    case 'BIDIRECTIONAL': return 'secondary'
    default: return 'grey'
  }
}

function getDirectionLabel(direction) {
  switch (direction) {
    case 'INBOUND': return 'Masuk'
    case 'OUTBOUND': return 'Keluar'
    case 'BIDIRECTIONAL': return 'Dua Arah'
    default: return direction
  }
}

function getStatusColor(status) {
  switch (status) {
    case 'PENDING_FOLLOW_UP': return 'warning'
    case 'FOLLOWED_UP': return 'info'
    case 'CLOSED': return 'positive'
    default: return 'grey'
  }
}

function getStatusLabel(status) {
  switch (status) {
    case 'PENDING_FOLLOW_UP': return 'Perlu Follow Up'
    case 'FOLLOWED_UP': return 'Ditindaklanjuti'
    case 'CLOSED': return 'Selesai'
    default: return status
  }
}

onMounted(() => {
  loadCommunications()
  loadStudents()
})
</script>
