<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-md">
      <div class="text-h5 text-weight-bold">Surat Masuk</div>
      <q-space />
      <q-btn color="primary" icon="add" label="Input Surat Masuk" @click="openForm()" />
    </div>

    <!-- Toolbar Filter -->
    <q-card flat bordered class="q-mb-md bg-grey-1">
      <q-card-section class="row q-col-gutter-sm items-center">
        <div class="col-12 col-md-3">
          <q-input v-model="filter.search" label="Cari Perihal/Nomor" outlined dense clearable @update:model-value="fetchSurat">
            <template v-slot:prepend><q-icon name="search" /></template>
          </q-input>
        </div>
        <div class="col-6 col-md-2">
          <q-select v-model="filter.jenis_surat" :options="letterTypeStore.flatOptions" option-group="group" label="Jenis Surat" outlined dense clearable option-label="label" option-value="value" emit-value map-options @update:model-value="fetchSurat" />
        </div>
        <div class="col-6 col-md-2">
          <q-select v-model="filter.status" :options="['DRAFT', 'PROCESSED', 'ARCHIVED']" label="Status" outlined dense clearable @update:model-value="fetchSurat" />
        </div>
        <div class="col-6 col-md-2">
          <q-input v-model="filter.date_from" label="Dari Tanggal" type="date" outlined dense @update:model-value="fetchSurat" />
        </div>
        <div class="col-6 col-md-2">
          <q-input v-model="filter.date_to" label="Sampai Tanggal" type="date" outlined dense @update:model-value="fetchSurat" />
        </div>
      </q-card-section>
    </q-card>

    <!-- Tabel Data -->
    <q-table :rows="suratStore.suratMasuk" :columns="columns" row-key="id" flat bordered :loading="suratStore.isLoading">
      <template v-slot:body-cell-agenda="props">
        <q-td :props="props">
          <span class="text-weight-bold text-primary">{{ props.row.nomor_agenda || '-' }}</span>
        </q-td>
      </template>
      <template v-slot:body-cell-jenis="props">
        <q-td :props="props">
          <q-badge color="blue-6">{{ getLabelJenis(props.row.jenis_surat) }}</q-badge>
        </q-td>
      </template>
      <template v-slot:body-cell-status="props">
        <q-td :props="props">
          <q-badge :color="getStatusColor(props.row.status)">{{ props.row.status }}</q-badge>
        </q-td>
      </template>
      <template v-slot:body-cell-actions="props">
        <q-td :props="props" class="q-gutter-x-xs">
          <q-btn flat round dense color="primary" icon="visibility" @click="printSurat(props.row)">
            <q-tooltip>Lihat PDF</q-tooltip>
          </q-btn>
          <q-btn flat round dense color="orange" icon="forward_to_inbox" @click="openDisposisi(props.row)">
            <q-tooltip>Disposisi</q-tooltip>
          </q-btn>
          <q-btn v-if="!props.row.ttd_url" flat round dense color="teal" icon="draw" @click="openTtd(props.row)">
            <q-tooltip>Tanda Tangan</q-tooltip>
          </q-btn>
          <q-btn v-if="props.row.status === 'DRAFT'" flat round dense color="blue" icon="edit" @click="openForm(props.row)">
            <q-tooltip>Edit</q-tooltip>
          </q-btn>
          <q-btn flat round dense color="negative" icon="delete" @click="confirmDelete(props.row)">
            <q-tooltip>Hapus</q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </q-table>

    <!-- Dialog Form Input/Edit -->
    <q-dialog v-model="showForm" persistent>
      <q-card style="min-width: 600px; max-width: 90vw">
        <q-card-section class="row items-center">
          <div class="text-h6">{{ isEdit ? 'Edit Surat Masuk' : 'Input Surat Masuk' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-none q-gutter-y-md scroll" style="max-height: 70vh">
          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-select v-model="form.jenis_surat" :options="letterTypeStore.flatOptions" option-group="group" label="Jenis Surat" outlined dense option-label="label" option-value="value" emit-value map-options :disable="isEdit" />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.kode_klasifikasi" label="Kode Klasifikasi (Opsional)" outlined dense :disable="isEdit" />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.number" label="Nomor Surat" outlined dense :disable="isEdit" />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.date" label="Tanggal Surat" type="date" outlined dense :disable="isEdit" />
            </div>
            <div class="col-12">
              <q-input v-model="form.subject" label="Perihal" outlined dense />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.from" label="Pengirim" outlined dense />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="form.to" label="Penerima" outlined dense />
            </div>
            <div class="col-12">
              <q-input v-model="form.description" label="Keterangan" type="textarea" outlined dense />
            </div>
          </div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Simpan" :loading="isSubmitting" @click="submitForm" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Disposisi -->
    <q-dialog v-model="showDispo" persistent>
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">Disposisi Surat</div>
          <div class="text-caption text-grey-7">{{ selectedRow?.subject }}</div>
        </q-card-section>
        <q-card-section class="q-pt-none q-gutter-y-md">
          <q-select v-model="dispoForm.to_user_id" :options="teachers" label="Pilih Guru / Staf" outlined dense option-label="name" option-value="id" emit-value map-options use-input input-debounce="300" @filter="filterTeachers" />
          <q-input v-model="dispoForm.note" label="Catatan / Instruksi" type="textarea" outlined dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="orange" label="Kirim Disposisi" :loading="isSubmitting" @click="submitDisposisi" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog TTD Digital -->
    <q-dialog v-model="showTtd" persistent>
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">Tanda Tangan Digital</div>
        </q-card-section>
        <q-card-section class="q-pt-none q-gutter-y-md text-center">
          <div v-if="previewTtd" class="q-mb-md">
            <img :src="previewTtd" style="max-height: 150px; border: 1px dashed #ccc;" />
          </div>
          <q-file v-model="ttdFile" label="Pilih Gambar TTD (PNG)" outlined dense accept=".png" @update:model-value="handleFileUpload">
            <template v-slot:prepend><q-icon name="attach_file" /></template>
          </q-file>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="teal" label="Tanda Tangani" :loading="isSubmitting" :disable="!previewTtd" @click="submitTtd" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Notify, Dialog } from 'quasar'
import { api } from '@/boot/axios'
import { useSuratStore } from '@/stores/surat/surat'
import { useLetterTypeStore } from '@/stores/letter-type'

const suratStore = useSuratStore()
const letterTypeStore = useLetterTypeStore()

// State Tabel & Filter
const filter = ref({ search: '', jenis_surat: null, status: null, date_from: '', date_to: '' })

// State Form
const showForm = ref(false)
const isEdit = ref(false)
const isSubmitting = ref(false)
const form = ref({ id: '', type: 'INCOMING', jenis_surat: '', kode_klasifikasi: '', number: '', date: '', subject: '', from: '', to: '', description: '' })

// State Disposisi
const showDispo = ref(false)
const selectedRow = ref(null)
const dispoForm = ref({ to_user_id: '', note: '' })
const teachers = ref([])
const allTeachers = ref([])

// State TTD
const showTtd = ref(false)
const ttdFile = ref(null)
const previewTtd = ref(null)

function getLabelJenis(val) {
  const item = letterTypeStore.list.find(i => i.value === val)
  return item ? item.label : val
}

const columns = [
  { name: 'agenda', label: 'No. Agenda', align: 'left' },
  { name: 'number', label: 'Nomor Surat', field: 'number', align: 'left', sortable: true },
  { name: 'jenis', label: 'Jenis', align: 'left' },
  { name: 'date', label: 'Tanggal', field: r => new Date(r.date).toLocaleDateString('id-ID'), align: 'left', sortable: true },
  { name: 'subject', label: 'Perihal', field: 'subject', align: 'left' },
  { name: 'from', label: 'Pengirim', field: 'from', align: 'left' },
  { name: 'status', label: 'Status', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center' }
]

function getStatusColor(status) {
  if (status === 'DRAFT') return 'grey-6'
  if (status === 'PROCESSED') return 'blue-6'
  if (status === 'ARCHIVED') return 'green-6'
  return 'grey'
}

async function fetchSurat() {
  try {
    await suratStore.fetchSuratMasuk(filter.value)
  } catch {
    Notify.create({ type: 'negative', message: 'Gagal memuat surat masuk' })
  }
}

function openForm(row = null) {
  if (row) {
    isEdit.value = true
    form.value = { ...row, date: row.date.split('T')[0] }
  } else {
    isEdit.value = false
    form.value = {
      id: '', type: 'INCOMING', jenis_surat: '', kode_klasifikasi: '', number: '',
      date: new Date().toISOString().slice(0, 10), subject: '', from: '', to: 'Kepala Sekolah', description: ''
    }
  }
  showForm.value = true
}

async function submitForm() {
  if (!form.value.subject || !form.value.from || !form.value.jenis_surat) {
    return Notify.create({ type: 'warning', message: 'Jenis, Perihal, dan Pengirim wajib diisi' })
  }
  isSubmitting.value = true
  try {
    if (isEdit.value) {
      await suratStore.updateSurat(form.value.id, {
        kategori: form.value.kategori,
        kode_klasifikasi: form.value.kode_klasifikasi,
        subject: form.value.subject,
        description: form.value.description
      })
      Notify.create({ type: 'positive', message: 'Surat masuk berhasil diperbarui' })
    } else {
      await suratStore.createSurat(form.value)
      Notify.create({ type: 'positive', message: 'Surat masuk berhasil dicatat' })
    }
    showForm.value = false
    await fetchSurat()
  } catch (err) {
    Notify.create({ type: 'negative', message: err.response?.data?.error || 'Gagal menyimpan surat' })
  } finally {
    isSubmitting.value = false
  }
}

function confirmDelete(row) {
  Dialog.create({
    title: 'Konfirmasi Hapus',
    message: `Apakah Anda yakin ingin menghapus surat "${row.subject}"? Tindakan ini tidak dapat dibatalkan.`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await suratStore.deleteSurat(row.id)
      Notify.create({ type: 'positive', message: 'Surat berhasil dihapus' })
      await fetchSurat()
    } catch {
      Notify.create({ type: 'negative', message: 'Gagal menghapus surat' })
    }
  })
}

// --- LOGIKA DISPOSISI ---
async function fetchTeachers() {
  try {
    const res = await api.get('/teachers')
    allTeachers.value = res.data?.data || []
    teachers.value = allTeachers.value
  } catch (err) {
    console.error('Gagal load guru', err)
  }
}

function filterTeachers(val, update) {
  update(() => {
    if (!val) {
      teachers.value = allTeachers.value
    } else {
      const needle = val.toLowerCase()
      teachers.value = allTeachers.value.filter(v => v.name.toLowerCase().includes(needle))
    }
  })
}

function openDisposisi(row) {
  selectedRow.value = row
  dispoForm.value = { to_user_id: '', note: '' }
  showDispo.value = true
}

async function submitDisposisi() {
  if (!dispoForm.value.to_user_id) {
    return Notify.create({ type: 'warning', message: 'Pilih guru/staf tujuan disposisi' })
  }
  isSubmitting.value = true
  try {
    await suratStore.disposisiSurat(selectedRow.value.id, dispoForm.value)
    Notify.create({ type: 'positive', message: 'Disposisi berhasil dikirim' })
    showDispo.value = false
  } catch {
    Notify.create({ type: 'negative', message: 'Gagal mengirim disposisi' })
  } finally {
    isSubmitting.value = false
  }
}

// --- LOGIKA TTD DIGITAL ---
function openTtd(row) {
  selectedRow.value = row
  ttdFile.value = null
  previewTtd.value = null
  showTtd.value = true
}

function handleFileUpload(file) {
  if (!file) return
  const reader = new FileReader()
  reader.onload = (e) => {
    previewTtd.value = e.target.result
  }
  reader.readAsDataURL(file)
}

async function submitTtd() {
  isSubmitting.value = true
  try {
    await suratStore.signSurat(selectedRow.value.id, {
      signer_role: 'Kepala Sekolah',
      image_base64: previewTtd.value
    })
    Notify.create({ type: 'positive', message: 'Tanda tangan digital berhasil disimpan' })
    showTtd.value = false
    await fetchSurat()
  } catch {
    Notify.create({ type: 'negative', message: 'Gagal menyimpan tanda tangan' })
  } finally {
    isSubmitting.value = false
  }
}

// --- LOGIKA PRINT PDF ---
async function printSurat(row) {
  try {
    const response = await api.get(`/correspondences/${row.id}/pdf`, { responseType: 'blob' })
    const file = new Blob([response.data], { type: 'application/pdf' })
    const fileURL = URL.createObjectURL(file)
    window.open(fileURL, '_blank')
  } catch {
    Notify.create({ type: 'negative', message: 'Gagal mengunduh PDF surat.' })
  }
}

onMounted(() => {
  fetchSurat()
  fetchTeachers()
  if (letterTypeStore.list.length === 0) {
    letterTypeStore.fetchList()
  }
})
</script>
