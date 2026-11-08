<template>
  <q-page class="q-pa-md">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Arsip Dokumen GTK</h1>
        <p class="text-caption text-grey-7 q-mb-none">Penyimpanan dokumen penting (SK, Kontrak, dll)</p>
      </div>
      <q-btn color="primary" icon="add" label="Upload Dokumen" @click="openFormDialog()" />
    </div>

    <!-- Alert dokumen mendekati kadaluarsa -->
    <q-banner v-if="expiringDocs.length" rounded class="bg-warning text-black q-mb-md">
      <template v-slot:avatar><q-icon name="warning" /></template>
      Terdapat {{ expiringDocs.length }} dokumen yang akan kadaluarsa dalam 30 hari ke depan.
      <q-btn flat label="Lihat" color="black" @click="expiryDialog = true" />
    </q-banner>

    <q-card class="shadow-1">
      <q-card-section class="row items-center q-pb-none">
        <q-btn color="secondary" icon="refresh" label="Refresh" @click="fetchArchives" :loading="isLoading" flat
          class="q-ml-auto" />
      </q-card-section>

      <q-card-section>
        <q-table :rows="archives" :columns="columns" row-key="id" flat bordered dense :loading="isLoading">
          <template v-slot:body-cell-tanggalKadaluarsa="props">
            <q-td :props="props" class="text-center">
              <span :class="{ 'text-negative text-weight-bold': isExpiringSoon(props.row.tanggalKadaluarsa) }">
                {{ formatDate(props.row.tanggalKadaluarsa) }}
              </span>
            </q-td>
          </template>
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center">
              <q-btn flat round dense icon="edit" color="info" @click="openFormDialog(props.row)" class="q-mr-sm">
                <q-tooltip>Edit Data</q-tooltip>
              </q-btn>
              <q-btn flat round dense icon="visibility" color="primary" type="a" :href="props.row.fileUrl"
                target="_blank" class="q-mr-sm">
                <q-tooltip>Lihat File</q-tooltip>
              </q-btn>
              <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row.id)">
                <q-tooltip>Hapus</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- Dialog form tambah/edit -->
    <q-dialog v-model="formDialog" persistent>
      <q-card style="min-width: 550px">
        <q-card-section class="row items-center bg-primary text-white">
          <div class="text-h6">{{ editMode ? 'Edit Dokumen' : 'Upload Dokumen' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-form @submit="submitForm">
            <q-select v-model="form.ptkId" :options="guruOptions" label="Pilih Guru" option-value="id"
              option-label="nama" emit-value map-options outlined dense :rules="[v => !!v || 'Wajib dipilih']"
              :disable="editMode" />
            <q-input v-model="form.jenisDokumen" label="Jenis Dokumen (SK, Kontrak, dll)" outlined dense
              :rules="[v => !!v || 'Wajib diisi']" />
            <q-input v-model="form.nomorDokumen" label="Nomor Dokumen" outlined dense />
            <div class="row q-col-gutter-sm">
              <div class="col-6"><q-input v-model="form.tanggalTerbit" label="Tanggal Terbit" type="date" outlined
                  dense />
              </div>
              <div class="col-6"><q-input v-model="form.tanggalKadaluarsa" label="Tanggal Kadaluarsa" type="date"
                  outlined dense hint="Kosongkan jika tidak ada" /></div>
            </div>
            <q-input v-model="form.catatan" label="Catatan" type="textarea" rows="2" outlined dense />

            <q-file v-model="fileUpload" label="Upload File (PDF/Gambar)" accept=".pdf,.jpg,.jpeg,.png" outlined dense>
              <template v-slot:prepend><q-icon name="attach_file" /></template>
            </q-file>
            <div v-if="editMode && form.fileUrl" class="text-caption text-green-8">
              <q-icon name="check_circle" /> File sudah ada. Abaikan input di atas jika tidak ingin mengganti file.
            </div>

            <div class="row justify-end q-mt-md">
              <q-btn label="Batal" flat color="negative" v-close-popup class="q-mr-sm" />
              <q-btn label="Simpan" type="submit" color="primary" :loading="isSubmitting" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog daftar dokumen kadaluarsa -->
    <q-dialog v-model="expiryDialog">
      <q-card style="min-width: 500px">
        <q-card-section class="bg-warning text-black">
          <div class="text-h6">Dokumen Mendekati Kadaluarsa</div>
        </q-card-section>
        <q-card-section>
          <q-list dense>
            <q-item v-for="doc in expiringDocs" :key="doc.id">
              <q-item-section>
                <q-item-label>{{ doc.jenisDokumen }} - {{ doc.teacherName }}</q-item-label>
                <q-item-label caption class="text-negative">Kadaluarsa: {{ formatDate(doc.tanggalKadaluarsa)
                }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Tutup" color="primary" v-close-popup />
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

const $q = useQuasar()
const guruStore = useGuruStore()

const archives = ref([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const formDialog = ref(false)
const expiryDialog = ref(false)
const editMode = ref(false)
const editId = ref(null)
const fileUpload = ref(null)

const form = ref({
  ptkId: null,
  jenisDokumen: '',
  nomorDokumen: '',
  tanggalTerbit: '',
  tanggalKadaluarsa: '',
  catatan: '',
  fileUrl: ''
})

const columns = [
  { name: 'teacherName', label: 'Nama Guru', field: 'teacherName', align: 'left', sortable: true },
  { name: 'jenisDokumen', label: 'Jenis', field: 'jenisDokumen', align: 'left' },
  { name: 'nomorDokumen', label: 'Nomor', field: 'nomorDokumen', align: 'left' },
  { name: 'tanggalTerbit', label: 'Terbit', field: 'tanggalTerbit', align: 'center', format: val => val ? new Date(val).toLocaleDateString('id-ID') : '-' },
  { name: 'tanggalKadaluarsa', label: 'Kadaluarsa', field: 'tanggalKadaluarsa', align: 'center' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const guruOptions = computed(() => guruStore.list.map(g => ({ id: g.id, nama: g.name })))

const expiringDocs = computed(() => {
  return archives.value.filter(a => isExpiringSoon(a.tanggalKadaluarsa))
})

const isExpiringSoon = (dateStr) => {
  if (!dateStr) return false
  const expiry = new Date(dateStr)
  const today = new Date()
  const diffDays = Math.ceil((expiry - today) / (1000 * 60 * 60 * 24))
  return diffDays <= 30 && diffDays >= 0
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID')
}

const fetchArchives = async () => {
  isLoading.value = true
  try {
    const res = await api.get('/teachers/archives')
    archives.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: `Gagal memuat arsip dokumen${err}` })
  } finally {
    isLoading.value = false
  }
}

const openFormDialog = (row = null) => {
  fileUpload.value = null
  if (row) {
    editMode.value = true
    editId.value = row.id
    form.value = {
      ptkId: row.ptkId,
      jenisDokumen: row.jenisDokumen,
      nomorDokumen: row.nomorDokumen,
      tanggalTerbit: row.tanggalTerbit ? new Date(row.tanggalTerbit).toISOString().split('T')[0] : '',
      tanggalKadaluarsa: row.tanggalKadaluarsa ? new Date(row.tanggalKadaluarsa).toISOString().split('T')[0] : '',
      catatan: row.catatan,
      fileUrl: row.fileUrl
    }
  } else {
    editMode.value = false
    editId.value = null
    form.value = { ptkId: null, jenisDokumen: '', nomorDokumen: '', tanggalTerbit: '', tanggalKadaluarsa: '', catatan: '', fileUrl: '' }
  }
  formDialog.value = true
}

const submitForm = async () => {
  isSubmitting.value = true
  try {
    let finalFileUrl = form.value.fileUrl

    // Upload file jika ada file baru dipilih
    if (fileUpload.value) {
      const formData = new FormData()
      formData.append('document', fileUpload.value)
      const uploadRes = await api.post('/teachers/archives/upload', formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      finalFileUrl = uploadRes.data.url
    }

    if (!finalFileUrl) {
      $q.notify({ type: 'warning', message: 'Silakan upload file terlebih dahulu' })
      isSubmitting.value = false
      return
    }

    const payload = { ...form.value, fileUrl: finalFileUrl }

    if (editMode.value) {
      await api.put(`/teachers/archives/${editId.value}`, payload)
      $q.notify({ type: 'positive', message: 'Dokumen berhasil diupdate' })
    } else {
      await api.post('/teachers/archives', payload)
      $q.notify({ type: 'positive', message: 'Dokumen berhasil diupload' })
    }

    formDialog.value = false
    await fetchArchives()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Gagal menyimpan dokumen' })
  } finally {
    isSubmitting.value = false
  }
}

const confirmDelete = (id) => {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Hapus dokumen ini? File di MinIO tidak akan terhapus otomatis.',
    cancel: true
  }).onOk(async () => {
    try {
      await api.delete(`/teachers/archives/${id}`)
      $q.notify({ type: 'positive', message: 'Dokumen berhasil dihapus' })
      await fetchArchives()
    } catch (err) {
      $q.notify({ type: 'negative', message: `Gagal menghapus dokumen ${err}` })
    }
  })
}

onMounted(async () => {
  await guruStore.fetchTeachers()
  await fetchArchives()
})
</script>
