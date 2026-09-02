<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Kelola Soal</div>

    <q-tabs v-model="tab" dense class="q-mb-md" align="left">
      <q-tab name="bank" label="Bank Soal" />
      <q-tab name="form" label="Buat/Edit Soal" />
      <q-tab name="kategori" label="Kategori Soal" />
    </q-tabs>

    <q-tab-panels v-model="tab" animated>
      <!-- Tab Bank Soal -->
      <q-tab-panel name="bank" class="q-pa-none">
        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12 col-md-3">
            <q-input dense outlined label="Cari soal..." v-model="searchQuery" debounce="300">
              <template v-slot:append><q-icon name="search" /></template>
            </q-input>
          </div>
          <div class="col-12 col-md-2">
            <q-select dense outlined label="Kategori" v-model="filterKategori" :options="kategoriOptions" clearable />
          </div>
          <div class="col-12 col-md-2">
            <q-select dense outlined label="Tingkat Kesulitan" v-model="filterKesulitan"
              :options="['Mudah', 'Sedang', 'Sulit']" clearable />
          </div>
          <div class="col-12 col-md-1">
            <q-btn color="primary" label="Filter" @click="loadSoal" />
          </div>
        </div>

        <q-table :rows="filteredSoal" :columns="soalColumns" row-key="id" flat bordered dense :loading="loading">
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props">
              <q-btn flat round dense icon="edit" color="info" @click="editSoal(props.row)" class="q-mr-sm" />
              <q-btn flat round dense icon="delete" color="negative" @click="confirmDeleteSoal(props.row)" />
              <q-btn flat round dense icon="visibility" color="primary" @click="previewSoal(props.row)" />
            </q-td>
          </template>
          <template v-slot:body-cell-pertanyaan="props">
            <q-td :props="props" class="text-left" style="max-width: 300px;">
              <div v-html="props.row.pertanyaan" class="ellipsis-2-lines"></div>
            </q-td>
          </template>
        </q-table>
      </q-tab-panel>

      <!-- Tab Buat/Edit Soal -->
      <q-tab-panel name="form" class="q-pa-none">
        <q-card flat bordered>
          <q-card-section>
            <div class="text-h6">{{ isEdit ? 'Edit Soal' : 'Buat Soal Baru' }}</div>
          </q-card-section>
          <q-card-section>
            <q-form @submit="submitSoal" class="q-gutter-md">
              <q-select v-model="soalForm.kategori_id" :options="kategoriList" label="Kategori" option-value="id"
                option-label="nama" outlined :rules="[val => !!val || 'Pilih kategori']" />
              <q-input v-model="soalForm.pertanyaan" label="Pertanyaan" type="textarea" rows="3" outlined
                :rules="[val => !!val || 'Pertanyaan harus diisi']" />
              <q-input v-model="soalForm.pilihan_a" label="Pilihan A" outlined :rules="[val => !!val]" />
              <q-input v-model="soalForm.pilihan_b" label="Pilihan B" outlined :rules="[val => !!val]" />
              <q-input v-model="soalForm.pilihan_c" label="Pilihan C" outlined />
              <q-input v-model="soalForm.pilihan_d" label="Pilihan D" outlined />
              <q-input v-model="soalForm.pilihan_e" label="Pilihan E (opsional)" outlined />
              <q-select v-model="soalForm.jawaban_benar" :options="jawabanOptions" label="Jawaban Benar" outlined
                :rules="[val => !!val]" />
              <q-select v-model="soalForm.tingkat_kesulitan" :options="['Mudah', 'Sedang', 'Sulit']"
                label="Tingkat Kesulitan" outlined :rules="[val => !!val]" />
              <q-input v-model="soalForm.pembahasan" label="Pembahasan (opsional)" type="textarea" rows="2" outlined />
              <div class="q-gutter-sm">
                <q-btn label="Simpan" type="submit" color="primary" />
                <q-btn label="Batal" flat color="negative" @click="resetForm" />
              </div>
            </q-form>
          </q-card-section>
        </q-card>
      </q-tab-panel>

      <!-- Tab Kategori Soal -->
      <q-tab-panel name="kategori" class="q-pa-none">
        <div class="q-mb-md">
          <q-btn color="primary" icon="add" label="Tambah Kategori" @click="openFormKategori(null)" />
        </div>
        <q-table :rows="kategoriList" :columns="kategoriColumns" row-key="id" flat bordered dense>
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props">
              <q-btn flat round dense icon="edit" color="info" @click="openFormKategori(props.row)" class="q-mr-sm" />
              <q-btn flat round dense icon="delete" color="negative" @click="confirmDeleteKategori(props.row)" />
            </q-td>
          </template>
        </q-table>
      </q-tab-panel>
    </q-tab-panels>

    <!-- Dialog Preview Soal -->
    <q-dialog v-model="previewDialog">
      <q-card style="min-width: 500px;">
        <q-card-section>
          <div class="text-h6">Preview Soal</div>
        </q-card-section>
        <q-card-section v-if="previewData">
          <div><strong>Kategori:</strong> {{ previewData.kategori_nama }}</div>
          <div><strong>Pertanyaan:</strong></div>
          <div v-html="previewData.pertanyaan" class="q-mb-md"></div>
          <div><strong>Pilihan:</strong></div>
          <div>A. {{ previewData.pilihan_a }}</div>
          <div>B. {{ previewData.pilihan_b }}</div>
          <div>C. {{ previewData.pilihan_c || '-' }}</div>
          <div>D. {{ previewData.pilihan_d || '-' }}</div>
          <div>E. {{ previewData.pilihan_e || '-' }}</div>
          <div><strong>Jawaban Benar:</strong> {{ previewData.jawaban_benar }}</div>
          <div><strong>Tingkat Kesulitan:</strong> {{ previewData.tingkat_kesulitan }}</div>
          <div v-if="previewData.pembahasan"><strong>Pembahasan:</strong> {{ previewData.pembahasan }}</div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Tutup" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Form Kategori -->
    <q-dialog v-model="kategoriDialog" persistent>
      <q-card style="min-width: 400px;">
        <q-card-section>
          <div class="text-h6">{{ kategoriFormTitle }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitKategori" class="q-gutter-md">
            <q-input v-model="kategoriForm.nama" label="Nama Kategori" lazy-rules
              :rules="[val => !!val || 'Nama kategori harus diisi']" outlined autofocus />
            <q-input v-model="kategoriForm.deskripsi" label="Deskripsi (opsional)" type="textarea" rows="2" outlined />
            <div class="q-gutter-sm">
              <q-btn label="Simpan" type="submit" color="primary" />
              <q-btn label="Batal" flat color="negative" v-close-popup />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog Konfirmasi Hapus -->
    <q-dialog v-model="deleteDialog" persistent>
      <q-card>
        <q-card-section>
          <div class="text-h6">Konfirmasi Hapus</div>
        </q-card-section>
        <q-card-section>
          Apakah Anda yakin ingin menghapus <strong>{{ deleteItem?.nama || deleteItem?.pertanyaan }}</strong>?
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" color="negative" v-close-popup />
          <q-btn flat label="Hapus" color="primary" @click="executeDelete" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// Data dummy
const soalList = ref([
  { id: 1, kategori_id: 1, kategori_nama: 'Matematika', pertanyaan: 'Hasil dari 2 + 2 adalah?', pilihan_a: '3', pilihan_b: '4', pilihan_c: '5', pilihan_d: '6', pilihan_e: '', jawaban_benar: 'B', tingkat_kesulitan: 'Mudah', pembahasan: '2+2=4' },
  { id: 2, kategori_id: 2, kategori_nama: 'Fisika', pertanyaan: 'Satuan gaya dalam SI adalah?', pilihan_a: 'Joule', pilihan_b: 'Watt', pilihan_c: 'Newton', pilihan_d: 'Pascal', pilihan_e: '', jawaban_benar: 'C', tingkat_kesulitan: 'Sedang', pembahasan: 'Gaya satuannya Newton' }
])
const kategoriList = ref([
  { id: 1, nama: 'Matematika', deskripsi: 'Soal-soal matematika' },
  { id: 2, nama: 'Fisika', deskripsi: 'Soal-soal fisika' },
  { id: 3, nama: 'Bahasa Indonesia', deskripsi: 'Soal bahasa' }
])

const tab = ref('bank')
const searchQuery = ref('')
const filterKategori = ref(null)
const filterKesulitan = ref(null)
const loading = ref(false)

// Form state
const isEdit = ref(false)
const editId = ref(null)
const soalForm = ref({
  kategori_id: null,
  pertanyaan: '',
  pilihan_a: '',
  pilihan_b: '',
  pilihan_c: '',
  pilihan_d: '',
  pilihan_e: '',
  jawaban_benar: null,
  tingkat_kesulitan: null,
  pembahasan: ''
})
const jawabanOptions = ['A', 'B', 'C', 'D', 'E']

// Kategori form
const kategoriDialog = ref(false)
const isEditKategori = ref(false)
const editKategoriId = ref(null)
const kategoriForm = ref({ nama: '', deskripsi: '' })
const kategoriFormTitle = ref('Tambah Kategori')

// Delete & preview
const deleteDialog = ref(false)
const deleteType = ref('') // 'soal' atau 'kategori'
const deleteItem = ref(null)
const previewDialog = ref(false)
const previewData = ref(null)

const kategoriOptions = computed(() => kategoriList.value.map(k => ({ label: k.nama, value: k.id })))

const filteredSoal = computed(() => {
  let result = soalList.value
  if (searchQuery.value) {
    result = result.filter(s => s.pertanyaan.toLowerCase().includes(searchQuery.value.toLowerCase()))
  }
  if (filterKategori.value) {
    result = result.filter(s => s.kategori_id === filterKategori.value.value)
  }
  if (filterKesulitan.value) {
    result = result.filter(s => s.tingkat_kesulitan === filterKesulitan.value)
  }
  return result
})

const soalColumns = [
  { name: 'id', label: 'ID', field: 'id', align: 'center' },
  { name: 'kategori_nama', label: 'Kategori', field: 'kategori_nama', align: 'left' },
  { name: 'pertanyaan', label: 'Pertanyaan', field: 'pertanyaan', align: 'left' },
  { name: 'tingkat_kesulitan', label: 'Kesulitan', field: 'tingkat_kesulitan', align: 'center' },
  { name: 'jawaban_benar', label: 'Jawaban', field: 'jawaban_benar', align: 'center' },
  { name: 'aksi', label: 'Aksi', field: 'aksi', align: 'center' }
]

const kategoriColumns = [
  { name: 'nama', label: 'Nama Kategori', field: 'nama', align: 'left' },
  { name: 'deskripsi', label: 'Deskripsi', field: 'deskripsi', align: 'left' },
  { name: 'aksi', label: 'Aksi', field: 'aksi', align: 'center' }
]

function loadSoal() { loading.value = false } // dummy

function editSoal(soal) {
  isEdit.value = true
  editId.value = soal.id
  soalForm.value = { ...soal, kategori_id: soal.kategori_id }
  tab.value = 'form'
}

function resetForm() {
  isEdit.value = false
  editId.value = null
  soalForm.value = {
    kategori_id: null, pertanyaan: '', pilihan_a: '', pilihan_b: '', pilihan_c: '', pilihan_d: '', pilihan_e: '',
    jawaban_benar: null, tingkat_kesulitan: null, pembahasan: ''
  }
  tab.value = 'bank'
}

function submitSoal() {
  if (isEdit.value) {
    const index = soalList.value.findIndex(s => s.id === editId.value)
    if (index !== -1) {
      soalList.value[index] = { ...soalList.value[index], ...soalForm.value, kategori_nama: kategoriList.value.find(k => k.id === soalForm.value.kategori_id)?.nama }
      $q.notify({ type: 'positive', message: 'Soal berhasil diupdate' })
    }
  } else {
    const newId = Date.now()
    const newSoal = {
      id: newId,
      ...soalForm.value,
      kategori_nama: kategoriList.value.find(k => k.id === soalForm.value.kategori_id)?.nama
    }
    soalList.value.push(newSoal)
    $q.notify({ type: 'positive', message: 'Soal berhasil ditambahkan' })
  }
  resetForm()
}

function confirmDeleteSoal(soal) {
  deleteType.value = 'soal'
  deleteItem.value = soal
  deleteDialog.value = true
}

function confirmDeleteKategori(kategori) {
  // cek apakah kategori memiliki soal
  const hasSoal = soalList.value.some(s => s.kategori_id === kategori.id)
  if (hasSoal) {
    $q.notify({ type: 'warning', message: 'Kategori ini memiliki soal, tidak bisa dihapus' })
    return
  }
  deleteType.value = 'kategori'
  deleteItem.value = kategori
  deleteDialog.value = true
}

function executeDelete() {
  if (deleteType.value === 'soal') {
    const index = soalList.value.findIndex(s => s.id === deleteItem.value.id)
    if (index !== -1) soalList.value.splice(index, 1)
    $q.notify({ type: 'positive', message: 'Soal berhasil dihapus' })
  } else if (deleteType.value === 'kategori') {
    const index = kategoriList.value.findIndex(k => k.id === deleteItem.value.id)
    if (index !== -1) kategoriList.value.splice(index, 1)
    $q.notify({ type: 'positive', message: 'Kategori berhasil dihapus' })
  }
  deleteDialog.value = false
  deleteItem.value = null
}

function previewSoal(soal) {
  previewData.value = soal
  previewDialog.value = true
}

// Kategori methods
function openFormKategori(kategori) {
  if (kategori) {
    isEditKategori.value = true
    editKategoriId.value = kategori.id
    kategoriForm.value = { nama: kategori.nama, deskripsi: kategori.deskripsi || '' }
    kategoriFormTitle.value = 'Edit Kategori'
  } else {
    isEditKategori.value = false
    editKategoriId.value = null
    kategoriForm.value = { nama: '', deskripsi: '' }
    kategoriFormTitle.value = 'Tambah Kategori'
  }
  kategoriDialog.value = true
}

function submitKategori() {
  if (isEditKategori.value) {
    const index = kategoriList.value.findIndex(k => k.id === editKategoriId.value)
    if (index !== -1) {
      kategoriList.value[index].nama = kategoriForm.value.nama
      kategoriList.value[index].deskripsi = kategoriForm.value.deskripsi
      $q.notify({ type: 'positive', message: 'Kategori berhasil diupdate' })
      // update kategori_nama di soal yang terkait
      soalList.value.forEach(s => {
        if (s.kategori_id === editKategoriId.value) s.kategori_nama = kategoriForm.value.nama
      })
    }
  } else {
    const newId = Date.now()
    kategoriList.value.push({ id: newId, nama: kategoriForm.value.nama, deskripsi: kategoriForm.value.deskripsi })
    $q.notify({ type: 'positive', message: 'Kategori berhasil ditambahkan' })
  }
  kategoriDialog.value = false
}

onMounted(() => {
  loadSoal()
})
</script>

<style scoped>
.ellipsis-2-lines {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
