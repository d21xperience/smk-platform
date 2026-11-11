<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-md">
      <div class="text-h5 text-weight-bold">Surat Keluar</div>
      <q-space />
      <q-btn color="primary" icon="add" label="Buat Surat Baru" @click="openForm()" />
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
    <q-table :rows="suratStore.suratKeluar" :columns="columns" row-key="id" flat bordered :loading="suratStore.isLoading">
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
          <q-btn flat round dense color="primary" icon="print" :loading="printingId === props.row.id" @click="printSurat(props.row)">
            <q-tooltip>Cetak PDF</q-tooltip>
          </q-btn>
          <q-btn flat round dense color="orange" icon="local_shipping" @click="openEkspedisi(props.row)">
            <q-tooltip>Buat Ekspedisi</q-tooltip>
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

    <!-- Dialog Form Stepper -->
    <q-dialog v-model="showForm" persistent maximized transition-show="slide-up" transition-hide="slide-down">
      <q-card class="column full-height">
        <q-card-section class="row items-center q-pb-none bg-primary text-white">
          <div class="text-h6">{{ isEdit ? 'Edit Surat Keluar' : 'Buat Surat Keluar Baru' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-stepper v-model="step" vertical color="primary" animated class="col scroll">
          <!-- LANGKAH 1: Klasifikasi & Penomoran -->
          <q-step :name="1" title="Klasifikasi & Penomoran" icon="folder" :done="step > 1">
            <div class="row q-col-gutter-md">
              <div class="col-12 col-md-6">
                <q-select v-model="form.peruntukan" :options="PERUNTUKAN_OPTIONS" label="Peruntukan Surat *" outlined dense emit-value map-options :disable="isEdit" @update:model-value="onPeruntukanChange" />
              </div>
              <div class="col-12 col-md-6">
                <q-select v-model="form.jenis_surat" :options="filteredJenisSurat" option-group="group" label="Pilih Jenis Surat *" outlined dense emit-value map-options option-label="label" option-value="value" :disable="isEdit || !form.peruntukan" @update:model-value="onJenisSuratChange" />
              </div>
              <div class="col-12 col-md-6">
                <q-input v-model="form.date" label="Tanggal Surat *" type="date" outlined dense :disable="isEdit" />
              </div>
              <div class="col-12 col-md-6">
                <q-input v-model="form.kode_klasifikasi" label="Kode Klasifikasi (Otomatis)" outlined dense readonly bg-color="grey-2">
                  <template v-slot:prepend><q-icon name="lock" size="xs" /></template>
                </q-input>
              </div>
              <div class="col-12 col-md-6">
                <q-input :model-value="isEdit ? form.number : 'AUTO-GENERATE'" label="Nomor Surat" outlined dense readonly bg-color="grey-2">
                  <template v-slot:prepend><q-icon name="lock" size="xs" /></template>
                </q-input>
              </div>
              <div class="col-12">
                <div class="text-subtitle2 q-mb-sm">Sifat Surat</div>
                <q-option-group v-model="form.konten.sifat" :options="sifatOptions" color="primary" inline />
              </div>
            </div>
            <q-stepper-navigation>
              <q-btn @click="step = 2" color="primary" label="Lanjut ke Pihak Terkait" :disable="!form.jenis_surat || !form.date" />
            </q-stepper-navigation>
          </q-step>

          <!-- LANGKAH 2: Pengirim & Penerima -->
          <q-step :name="2" title="Pengirim & Penerima" icon="people" :done="step > 2">
            <div class="row q-col-gutter-md">
              <div class="col-12 col-md-6">
                <q-input v-model="form.from" label="Pengirim / Penandatangan *" outlined dense readonly />
              </div>
              <div class="col-12 col-md-6">
                <q-input v-model="form.to" label="Penerima / Tujuan Surat *" outlined dense hint="Contoh: Kepala Dinas Pendidikan / PT Maju Jaya" />
              </div>
            </div>
            <q-stepper-navigation>
              <q-btn flat @click="step = 1" color="primary" label="Kembali" class="q-mr-sm" />
              <q-btn @click="step = 3" color="primary" label="Lanjut ke Isi Surat" :disable="!form.to" />
            </q-stepper-navigation>
          </q-step>

          <!-- LANGKAH 3: Editor Konten -->
          <q-step :name="3" title="Isi Surat" icon="edit_document" :done="step > 3">
            <div class="row q-col-gutter-md">
              <div class="col-12">
                <q-input v-model="form.subject" label="Perihal / Subjek Surat *" outlined dense />
              </div>
              <div class="col-12" v-if="dynamicFields.length > 0">
                <q-card flat bordered class="bg-blue-1 q-pa-md q-mb-md">
                  <div class="text-subtitle2 text-primary q-mb-sm">Data Spesifik ({{ getLabelJenis(form.jenis_surat) }})</div>
                  <div class="row q-col-gutter-sm">
                    <div class="col-12 col-md-6" v-for="field in dynamicFields" :key="field.key">
                      <q-input v-model="form.konten[field.key]" :label="field.label" outlined dense :type="field.type || 'text'" />
                    </div>
                  </div>
                </q-card>
              </div>
              <div class="col-12">
                <div class="text-subtitle2 q-mb-sm">Badan Surat (Editor)</div>
                <q-editor v-model="form.konten.body" min-height="10rem" :toolbar="[['bold', 'italic', 'underline'], ['quote', 'unordered', 'ordered'], ['undo', 'redo']]" />
              </div>
            </div>
            <q-stepper-navigation>
              <q-btn flat @click="step = 2" color="primary" label="Kembali" class="q-mr-sm" />
              <q-btn @click="step = 4" color="primary" label="Lanjut ke Validasi" :disable="!form.subject" />
            </q-stepper-navigation>
          </q-step>

          <!-- LANGKAH 4: Lampiran & Validasi -->
          <q-step :name="4" title="Lampiran & Validasi" icon="fact_check">
            <div class="row q-col-gutter-md">
              <div class="col-12">
                <q-file v-model="form.attachment_file" label="Unggah Dokumen Pendukung (Opsional)" outlined dense accept=".pdf,.jpg,.jpeg,.png">
                  <template v-slot:prepend><q-icon name="attach_file" /></template>
                </q-file>
              </div>
              <div class="col-12">
                <q-input v-model="form.description" label="Catatan Tambahan / Keterangan" type="textarea" outlined dense rows="3" />
              </div>
              <div class="col-12">
                <q-card flat bordered>
                  <q-card-section class="bg-grey-2">
                    <div class="text-subtitle1 text-weight-bold">Ringkasan Surat</div>
                  </q-card-section>
                  <q-list dense>
                    <q-item><q-item-section><strong>Jenis:</strong> {{ getLabelJenis(form.jenis_surat) }}</q-item-section></q-item>
                    <q-item><q-item-section><strong>Perihal:</strong> {{ form.subject }}</q-item-section></q-item>
                    <q-item><q-item-section><strong>Penerima:</strong> {{ form.to }}</q-item-section></q-item>
                    <q-item><q-item-section><strong>Tanggal:</strong> {{ form.date }}</q-item-section></q-item>
                  </q-list>
                </q-card>
              </div>
            </div>
            <q-stepper-navigation>
              <q-btn flat @click="step = 3" color="primary" label="Kembali" class="q-mr-sm" />
              <q-btn color="positive" label="Simpan sebagai Draft" :loading="isSubmitting" @click="submitForm" icon="save" />
              <q-btn color="primary" label="Ajukan Persetujuan" :loading="isSubmitting" @click="submitAndProcess" icon="send" class="q-ml-sm" />
            </q-stepper-navigation>
          </q-step>
        </q-stepper>
      </q-card>
    </q-dialog>

    <!-- Dialog TTD Digital -->
    <q-dialog v-model="showTtd" persistent>
      <q-card style="min-width: 400px">
        <q-card-section><div class="text-h6">Tanda Tangan Digital</div></q-card-section>
        <q-card-section class="q-pt-none q-gutter-y-md text-center">
          <div v-if="previewTtd" class="q-mb-md"><img :src="previewTtd" style="max-height: 150px; border: 1px dashed #ccc;" /></div>
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

    <!-- Dialog Ekspedisi -->
    <q-dialog v-model="showEkspedisi" persistent>
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">Buat Bukti Ekspedisi</div>
          <div class="text-caption text-grey-7">Surat: {{ selectedSurat?.subject }}</div>
        </q-card-section>
        <q-card-section class="q-pt-none q-gutter-y-md">
          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-input v-model="ekspedisiForm.penerima_nama" label="Nama Penerima *" outlined dense />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="ekspedisiForm.penerima_instansi" label="Instansi Penerima" outlined dense />
            </div>
            <div class="col-12">
              <q-input v-model="ekspedisiForm.penerima_alamat" label="Alamat Penerima" outlined dense />
            </div>
            <div class="col-12 col-md-6">
              <q-select v-model="ekspedisiForm.metode_kirim" :options="['KURIR', 'POS', 'EMAIL', 'TANGAN']" label="Metode Kirim *" outlined dense />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="ekspedisiForm.nomor_resi" label="Nomor Resi (Opsional)" outlined dense />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model="ekspedisiForm.tanggal_kirim" label="Tanggal Kirim *" type="date" outlined dense />
            </div>
            <div class="col-12 col-md-6">
              <q-input v-model.number="ekspedisiForm.jumlah_lembar" label="Jumlah Lembar" type="number" outlined dense />
            </div>
          </div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="orange" label="Simpan & Cetak Bukti" :loading="isSubmittingEkspedisi" @click="submitEkspedisi" icon="print" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { Notify, Dialog } from 'quasar'
import { api } from '@/boot/axios'
import { useSuratStore } from '@/stores/surat/surat'
import { useLetterTypeStore } from '@/stores/letter-type'
import { PERUNTUKAN_OPTIONS } from '@/constants/surat'

const suratStore = useSuratStore()
const letterTypeStore = useLetterTypeStore()

const printingId = ref(null)
const filter = ref({ search: '', jenis_surat: null, status: null, date_from: '', date_to: '' })

const showForm = ref(false)
const step = ref(1)
const isEdit = ref(false)
const isSubmitting = ref(false)
const form = ref({
  id: '', type: 'OUTGOING', peruntukan: '', jenis_surat: '', kode_klasifikasi: '', number: '', 
  date: '', subject: '', from: 'Kepala Sekolah', to: '', description: '',
  konten: { sifat: 'Biasa', body: '', formData: {} }
})

const showTtd = ref(false)
const selectedRow = ref(null)
const ttdFile = ref(null)
const previewTtd = ref(null)

const sifatOptions = [
  { label: 'Biasa', value: 'Biasa' },
  { label: 'Penting', value: 'Penting' },
  { label: 'Rahasia', value: 'Rahasia' }
]

// Ekspedisi State
const showEkspedisi = ref(false)
const selectedSurat = ref(null)
const isSubmittingEkspedisi = ref(false)
const ekspedisiForm = ref({
  penerima_nama: '', penerima_instansi: '', penerima_alamat: '', 
  metode_kirim: 'KURIR', nomor_resi: '', tanggal_kirim: '', jumlah_lembar: 1
})

function getLabelJenis(val) {
  const item = letterTypeStore.list.find(i => i.value === val)
  return item ? item.label : val
}

const filteredJenisSurat = computed(() => {
  return letterTypeStore.flatOptions.filter(opt => opt.peruntukan === form.value.peruntukan)
})

const dynamicFields = computed(() => {
  switch (form.value.jenis_surat) {
    case 'PERMOHONAN_PKL': return [ { key: 'nama_industri', label: 'Nama Perusahaan / DU/DI' }, { key: 'alamat_industri', label: 'Alamat Perusahaan' }, { key: 'jumlah_siswa', label: 'Jumlah Siswa', type: 'number' } ]
    case 'UNDANGAN_ORTU': return [ { key: 'hari_tanggal', label: 'Hari / Tanggal Kegiatan' }, { key: 'waktu', label: 'Waktu (Jam)' }, { key: 'tempat', label: 'Tempat Kegiatan' } ]
    case 'SPPD': return [ { key: 'tujuan', label: 'Kota Tujuan' }, { key: 'lama_hari', label: 'Lama Perjalanan (Hari)', type: 'number' }, { key: 'tgl_berangkat', label: 'Tanggal Berangkat', type: 'date' } ]
    default: return []
  }
})

const columns = [
  { name: 'agenda', label: 'No. Agenda', align: 'left' },
  { name: 'number', label: 'Nomor Surat', field: 'number', align: 'left', sortable: true },
  { name: 'jenis', label: 'Jenis', align: 'left' },
  { name: 'date', label: 'Tanggal', field: r => new Date(r.date).toLocaleDateString('id-ID'), align: 'left', sortable: true },
  { name: 'subject', label: 'Perihal', field: 'subject', align: 'left' },
  { name: 'to', label: 'Penerima', field: 'to', align: 'left' },
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
  try { await suratStore.fetchSuratKeluar(filter.value) } catch { Notify.create({ type: 'negative', message: 'Gagal memuat surat' }) }
}

function onPeruntukanChange() { form.value.jenis_surat = ''; form.value.kode_klasifikasi = '' }
function onJenisSuratChange(val) {
  form.value.konten.formData = {}
  const selected = letterTypeStore.flatOptions.find(s => s.value === val)
  if (selected) form.value.kode_klasifikasi = selected.kode_klasifikasi
  else form.value.kode_klasifikasi = ''
}

function openForm(row = null) {
  step.value = 1
  if (row) {
    isEdit.value = true
    form.value = { ...row, date: row.date.split('T')[0], konten: row.konten || { sifat: 'Biasa', body: '', formData: {} } }
    if(!form.value.konten.body) form.value.konten.body = ''
    if(!form.value.konten.sifat) form.value.konten.sifat = 'Biasa'
    const opt = letterTypeStore.flatOptions.find(o => o.value === row.jenis_surat)
    form.value.peruntukan = opt ? opt.peruntukan : 'UMUM'
  } else {
    isEdit.value = false
    form.value = { id: '', type: 'OUTGOING', peruntukan: '', jenis_surat: '', kode_klasifikasi: '', number: '', date: new Date().toISOString().slice(0, 10), subject: '', from: 'Kepala Sekolah', to: '', description: '', konten: { sifat: 'Biasa', body: '', formData: {} } }
  }
  showForm.value = true
}

async function submitForm(isProcess = false) {
  if (!form.value.peruntukan || !form.value.jenis_surat || !form.value.subject || !form.value.to) {
    return Notify.create({ type: 'warning', message: 'Peruntukan, Jenis, Perihal, dan Penerima wajib diisi' })
  }
  isSubmitting.value = true
  try {
    const payload = { ...form.value }
    if (payload.attachment_file) { payload.attachment_url = payload.attachment_file.name; delete payload.attachment_file }

    if (isEdit.value) {
      await suratStore.updateSurat(form.value.id, payload)
      Notify.create({ type: 'positive', message: 'Surat berhasil diperbarui' })
    } else {
      const res = await suratStore.createSurat(payload)
      Notify.create({ type: 'positive', message: 'Surat berhasil dibuat' })
      if (isProcess && res?.id) {
        await api.put(`/correspondences/${res.id}/archive`)
        Notify.create({ type: 'info', message: 'Surat diajukan untuk persetujuan' })
      }
    }
    showForm.value = false
    await fetchSurat()
  } catch (err) {
    Notify.create({ type: 'negative', message: err.response?.data?.error || 'Gagal menyimpan surat' })
  } finally { isSubmitting.value = false }
}

function submitAndProcess() { submitForm(true) }

function confirmDelete(row) {
  Dialog.create({ title: 'Konfirmasi Hapus', message: `Apakah Anda yakin ingin menghapus surat "${row.subject}"?`, cancel: true, persistent: true }).onOk(async () => {
    try { await suratStore.deleteSurat(row.id); Notify.create({ type: 'positive', message: 'Surat berhasil dihapus' }); await fetchSurat() } catch { Notify.create({ type: 'negative', message: 'Gagal menghapus surat' }) }
  })
}

// --- LOGIKA TTD DIGITAL ---
function openTtd(row) { selectedRow.value = row; ttdFile.value = null; previewTtd.value = null; showTtd.value = true }
function handleFileUpload(file) { if (!file) return; const reader = new FileReader(); reader.onload = (e) => { previewTtd.value = e.target.result }; reader.readAsDataURL(file) }
async function submitTtd() {
  isSubmitting.value = true
  try {
    await suratStore.signSurat(selectedRow.value.id, { signer_role: 'Kepala Sekolah', image_base64: previewTtd.value })
    Notify.create({ type: 'positive', message: 'Tanda tangan berhasil disimpan' })
    showTtd.value = false
    await fetchSurat()
  } catch { Notify.create({ type: 'negative', message: 'Gagal menyimpan tanda tangan' }) } finally { isSubmitting.value = false }
}

// --- LOGIKA EKSPEDISI ---
function openEkspedisi(row) {
  selectedSurat.value = row
  ekspedisiForm.value = { penerima_nama: row.to, penerima_instansi: '', penerima_alamat: '', metode_kirim: 'KURIR', nomor_resi: '', tanggal_kirim: new Date().toISOString().slice(0, 10), jumlah_lembar: 1 }
  showEkspedisi.value = true
}

async function submitEkspedisi() {
  if (!ekspedisiForm.value.penerima_nama || !ekspedisiForm.value.tanggal_kirim) {
    return Notify.create({ type: 'warning', message: 'Nama Penerima dan Tanggal Kirim wajib diisi' })
  }
  isSubmittingEkspedisi.value = true
  try {
    const res = await suratStore.createEkspedisi(selectedSurat.value.id, ekspedisiForm.value)
    Notify.create({ type: 'positive', message: 'Ekspedisi berhasil dibuat' })
    showEkspedisi.value = false
    
    // Buka PDF bukti ekspedisi
    if (res.bukti_pdf_url) {
      const url = `http://localhost:8081${res.bukti_pdf_url}` // Sesuaikan base URL backend
      window.open(url, '_blank')
    }
  } catch (err) {
    Notify.create({ type: 'negative', message: err.response?.data?.error || 'Gagal membuat ekspedisi' })
  } finally {
    isSubmittingEkspedisi.value = false
  }
}

// --- LOGIKA PRINT PDF ---
async function printSurat(row) {
  printingId.value = row.id
  try {
    const response = await api.get(`/correspondences/${row.id}/pdf`, { responseType: 'blob' })
    const file = new Blob([response.data], { type: 'application/pdf' })
    window.open(URL.createObjectURL(file), '_blank')
  } catch { Notify.create({ type: 'negative', message: 'Gagal mengunduh PDF surat.' }) } finally { printingId.value = null }
}

onMounted(() => {
  fetchSurat()
  if (letterTypeStore.list.length === 0) letterTypeStore.fetchList()
})
</script>
