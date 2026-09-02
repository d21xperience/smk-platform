<template>
  <q-page class="q-pa-md">
    <div class="row justify-between q-mb-md">
      <div class="text-h5">Pendaftaran Siswa Baru (Panitia)</div>
      <q-btn color="primary" icon="add" label="Tambah Pendaftar" @click="openForm" />
    </div>

    <q-table
      :rows="rows"
      :columns="columns"
      row-key="id"
      flat
      bordered
      dense
    >
      <template v-slot:body-cell-pembayaran="props">
        <q-td :props="props">
          <q-badge :color="getPaymentColor(props.row.pembayaran.status)">
            {{ props.row.pembayaran.status }}
          </q-badge>
          <div class="text-caption">Sisa: Rp {{ formatRupiah(props.row.pembayaran.sisa) }}</div>
          <div class="text-caption" v-if="props.row.pembayaran.tanggal_bayar_terakhir">
            Bayar terakhir: {{ formatTanggal(props.row.pembayaran.tanggal_bayar_terakhir) }}
          </div>
        </q-td>
      </template>
      <template v-slot:body-cell-kelengkapan="props">
        <q-td :props="props">
          <q-linear-progress :value="props.row.kelengkapan.uploaded / props.row.kelengkapan.total" />
          <div class="text-caption">{{ props.row.kelengkapan.uploaded }}/{{ props.row.kelengkapan.total }}</div>
        </q-td>
      </template>
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props">
          <q-btn flat round dense icon="edit" color="info" @click="openForm(props.row)" class="q-mr-sm">
            <q-tooltip>Edit Biodata</q-tooltip>
          </q-btn>
          <q-btn flat round dense icon="upload" color="primary" @click="openUploadBerkas(props.row)" class="q-mr-sm">
            <q-tooltip>Upload Berkas</q-tooltip>
          </q-btn>
          <q-btn flat round dense icon="payment" color="secondary" @click="openPembayaran(props.row)" class="q-mr-sm">
            <q-tooltip>Input Pembayaran</q-tooltip>
          </q-btn>
          <q-btn flat round dense icon="receipt" color="grey" @click="lihatRiwayat(props.row)" class="q-mr-sm">
            <q-tooltip>Riwayat Bayar</q-tooltip>
          </q-btn>
          <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row.id)">
            <q-tooltip>Hapus</q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </q-table>

    <!-- Dialog Form Biodata -->
    <q-dialog v-model="formDialog" persistent>
      <q-card style="min-width: 600px">
        <q-card-section><div class="text-h6">{{ editMode ? 'Edit Pendaftar' : 'Tambah Pendaftar' }}</div></q-card-section>
        <q-card-section>
          <q-form @submit="submitForm" class="q-gutter-md">
            <div class="row q-col-gutter-md">
              <div class="col-6"><q-input v-model="form.nama" label="Nama Lengkap" outlined :rules="[v=>!!v]" /></div>
              <div class="col-6"><q-input v-model="form.nisn" label="NISN" outlined /></div>
              <div class="col-6"><q-input v-model="form.tempat_lahir" label="Tempat Lahir" outlined /></div>
              <div class="col-6"><q-input v-model="form.tanggal_lahir" label="Tanggal Lahir" type="date" outlined /></div>
              <div class="col-6"><q-select v-model="form.jenis_kelamin" :options="['Laki-laki', 'Perempuan']" label="Jenis Kelamin" outlined /></div>
              <div class="col-6"><q-input v-model="form.agama" label="Agama" outlined /></div>
              <div class="col-12"><q-input v-model="form.alamat" label="Alamat" type="textarea" rows="2" outlined /></div>
              <div class="col-6"><q-input v-model="form.no_hp" label="No. HP" outlined /></div>
              <div class="col-6"><q-input v-model="form.email" label="Email" type="email" outlined /></div>
              <div class="col-12"><q-input v-model="form.asal_sekolah" label="Asal Sekolah" outlined /></div>
            </div>
            <div class="q-gutter-sm">
              <q-btn label="Simpan" type="submit" color="primary" />
              <q-btn label="Batal" flat color="negative" v-close-popup @click="resetForm" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog Upload Berkas -->
    <q-dialog v-model="berkasDialog">
      <q-card style="min-width: 500px">
        <q-card-section><div class="text-h6">Upload Berkas - {{ selectedCalon?.nama }}</div></q-card-section>
        <q-card-section>
          <div v-for="jenis in jenisBerkas" :key="jenis" class="q-mb-md">
            <div class="text-subtitle2">{{ jenis.toUpperCase() }}</div>
            <q-file v-model="uploadFiles[jenis]" :label="`Upload ${jenis}`" accept=".pdf,.jpg,.png" outlined dense />
            <div v-if="selectedCalon?.berkas[jenis]" class="text-caption">
              File sudah ada: <a href="#" @click.prevent="downloadFile(selectedCalon.berkas[jenis], jenis)">Lihat</a>
            </div>
          </div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn label="Simpan Semua" color="primary" @click="saveAllBerkas" />
          <q-btn label="Tutup" flat color="negative" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Input Pembayaran -->
    <q-dialog v-model="pembayaranDialog">
      <q-card style="min-width: 450px">
        <q-card-section><div class="text-h6">Input Pembayaran - {{ selectedCalon?.nama }}</div></q-card-section>
        <q-card-section>
          <div class="q-mb-sm">Total Tagihan: Rp {{ formatRupiah(selectedCalon?.pembayaran.total_tagihan) }}</div>
          <div class="q-mb-sm">Sudah Dibayar: Rp {{ formatRupiah(selectedCalon?.pembayaran.jumlah_bayar) }}</div>
          <div class="q-mb-sm">Sisa: Rp {{ formatRupiah(selectedCalon?.pembayaran.sisa) }}</div>
          <q-input v-model.number="jumlahBayar" label="Jumlah Bayar (Rp)" type="number" outlined dense :rules="[v=>v>0]" />
          <q-input v-model="catatanBayar" label="Catatan (opsional)" type="textarea" rows="2" outlined dense />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn label="Bayar" color="positive" @click="savePembayaran" />
          <q-btn label="Ubah Total Tagihan" color="warning" @click="editTotalTagihan" />
          <q-btn label="Tutup" flat color="negative" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Riwayat Pembayaran -->
    <q-dialog v-model="riwayatDialog">
      <q-card style="min-width: 500px">
        <q-card-section><div class="text-h6">Riwayat Pembayaran - {{ selectedCalon?.nama }}</div></q-card-section>
        <q-card-section>
          <q-table :rows="selectedCalon?.pembayaran.riwayat || []" :columns="riwayatColumns" dense flat />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn label="Tutup" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Dialog Edit Total Tagihan -->
    <q-dialog v-model="tagihanDialog">
      <q-card>
        <q-card-section><div class="text-h6">Ubah Total Tagihan</div></q-card-section>
        <q-card-section>
          <q-input v-model.number="newTotalTagihan" label="Total Tagihan (Rp)" type="number" outlined dense :rules="[v=>v>0]" />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn label="Simpan" color="primary" @click="saveTotalTagihan" />
          <q-btn label="Batal" flat color="negative" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { usePendaftaranStore } from 'stores/kesiswaan/pendaftaran'

const $q = useQuasar()
const pendaftaranStore = usePendaftaranStore()

// State
const formDialog = ref(false)
const editMode = ref(false)
const editId = ref(null)
const form = ref({
  nama: '', nisn: '', tempat_lahir: '', tanggal_lahir: '', jenis_kelamin: '', agama: '',
  alamat: '', no_hp: '', email: '', asal_sekolah: ''
})
const berkasDialog = ref(false)
const selectedCalon = ref(null)
const uploadFiles = ref({})
const jenisBerkas = ['ijazah', 'skl', 'kk', 'akte', 'foto']
const pembayaranDialog = ref(false)
const riwayatDialog = ref(false)
const tagihanDialog = ref(false)
const jumlahBayar = ref(0)
const catatanBayar = ref('')
const newTotalTagihan = ref(0)

const columns = [
  { name: 'nama', label: 'Nama', field: 'nama' },
  { name: 'nisn', label: 'NISN', field: 'nisn' },
  { name: 'asal_sekolah', label: 'Asal Sekolah', field: 'asal_sekolah' },
  { name: 'pembayaran', label: 'Pembayaran', field: 'pembayaran' },
  { name: 'kelengkapan', label: 'Berkas', field: 'kelengkapan' },
  { name: 'aksi', label: 'Aksi', field: 'aksi' }
]
const riwayatColumns = [
  { name: 'tanggal', label: 'Tanggal', field: 'tanggal' },
  { name: 'jumlah', label: 'Jumlah (Rp)', field: 'jumlah' },
  { name: 'sisa_setelah', label: 'Sisa', field: 'sisa_setelah' },
  { name: 'catatan', label: 'Catatan', field: 'catatan' }
]

const rows = computed(() => {
  return pendaftaranStore.list.map(calon => ({
    ...calon,
    kelengkapan: pendaftaranStore.getKelengkapanBerkas(calon)
  }))
})

function formatRupiah(value) {
  if (!value && value !== 0) return '0'
  return new Intl.NumberFormat('id-ID').format(value)
}
function formatTanggal(iso) {
  if (!iso) return '-'
  const d = new Date(iso)
  return d.toLocaleDateString('id-ID')
}
function getPaymentColor(status) {
  if (status === 'lunas') return 'positive'
  if (status === 'sebagian') return 'warning'
  return 'negative'
}

// Biodata
function openForm(calon = null) {
  if (calon) {
    editMode.value = true
    editId.value = calon.id
    form.value = { ...calon }
  } else {
    editMode.value = false
    resetForm()
  }
  formDialog.value = true
}
function resetForm() {
  form.value = {
    nama: '', nisn: '', tempat_lahir: '', tanggal_lahir: '', jenis_kelamin: '', agama: '',
    alamat: '', no_hp: '', email: '', asal_sekolah: ''
  }
}
function submitForm() {
  if (editMode.value) {
    pendaftaranStore.update(editId.value, form.value)
    $q.notify({ type: 'positive', message: 'Data pendaftar diupdate' })
  } else {
    pendaftaranStore.tambah(form.value)
    $q.notify({ type: 'positive', message: 'Pendaftar ditambahkan' })
  }
  formDialog.value = false
}
function confirmDelete(id) {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Hapus pendaftar ini?',
    cancel: true
  }).onOk(() => {
    pendaftaranStore.hapus(id)
    $q.notify({ type: 'positive', message: 'Pendaftar dihapus' })
  })
}

// Berkas
function openUploadBerkas(calon) {
  selectedCalon.value = calon
  uploadFiles.value = {}
  berkasDialog.value = true
}
function saveAllBerkas() {
  for (const jenis of jenisBerkas) {
    const file = uploadFiles.value[jenis]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        pendaftaranStore.updateBerkas(selectedCalon.value.id, jenis, e.target.result)
      }
      reader.readAsDataURL(file)
    }
  }
  $q.notify({ type: 'positive', message: 'Berkas diupload' })
  berkasDialog.value = false
}
function downloadFile(base64, filename) {
  const link = document.createElement('a')
  link.href = base64
  link.download = filename
  link.click()
}

// Pembayaran
function openPembayaran(calon) {
  selectedCalon.value = calon
  jumlahBayar.value = 0
  catatanBayar.value = ''
  pembayaranDialog.value = true
}
function savePembayaran() {
  if (jumlahBayar.value <= 0) {
    $q.notify({ type: 'negative', message: 'Jumlah bayar harus > 0' })
    return
  }
  pendaftaranStore.tambahPembayaran(selectedCalon.value.id, jumlahBayar.value, catatanBayar.value)
  $q.notify({ type: 'positive', message: 'Pembayaran dicatat' })
  pembayaranDialog.value = false
}
function lihatRiwayat(calon) {
  selectedCalon.value = calon
  riwayatDialog.value = true
}
function editTotalTagihan() {
  newTotalTagihan.value = selectedCalon.value.pembayaran.total_tagihan
  tagihanDialog.value = true
}
function saveTotalTagihan() {
  if (newTotalTagihan.value <= 0) {
    $q.notify({ type: 'negative', message: 'Total tagihan harus > 0' })
    return
  }
  pendaftaranStore.setTotalTagihan(selectedCalon.value.id, newTotalTagihan.value)
  $q.notify({ type: 'positive', message: 'Total tagihan diupdate' })
  tagihanDialog.value = false
  pembayaranDialog.value = false
}

onMounted(() => {
  pendaftaranStore.loadData()
})
</script>
