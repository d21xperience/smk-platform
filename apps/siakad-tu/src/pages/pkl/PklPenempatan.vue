<template>
  <q-page class="q-pa-md">
    <div class="row justify-between q-mb-md">
      <div class="text-h5">Penempatan PKL</div>
      <q-btn color="primary" icon="add" label="Penempatan Baru" @click="openForm" />
    </div>

    <q-table :rows="penempatanStore.list" :columns="columns" row-key="id" flat bordered dense>
      <template v-slot:body-cell-status="props">
        <q-td :props="props">
          <q-badge :color="props.row.status === 'ditempatkan' ? 'positive' : 'warning'">{{ props.row.status }}</q-badge>
        </q-td>
      </template>
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props">
          <q-btn flat round dense icon="edit" color="info" @click="openForm(props.row)" class="q-mr-sm" />
          <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row.id)" />
          <q-btn flat round dense icon="print" color="primary" @click="generateSurat(props.row)">
            <q-tooltip>Cetak Surat Penempatan</q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </q-table>

    <!-- Form dialog -->
    <q-dialog v-model="formDialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">{{ editMode ? 'Edit Penempatan' : 'Penempatan Baru' }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitPenempatan" class="q-gutter-md">
            <q-select v-model="form.siswa_id" :options="siswaOptions" label="Siswa" option-value="id"
              option-label="nama" outlined :rules="[v => !!v]" />
            <q-select v-model="form.mitra_id" :options="mitraOptions" label="Mitra Industri" option-value="id"
              option-label="nama" outlined :rules="[v => !!v]" />
            <q-input v-model="form.tanggal_mulai" label="Tanggal Mulai" type="date" outlined />
            <q-input v-model="form.tanggal_selesai" label="Tanggal Selesai" type="date" outlined />
            <q-input v-model="form.pembimbing_industri" label="Pembimbing Industri (Nama)" outlined />
            <q-input v-model="form.catatan" label="Catatan" type="textarea" rows="2" outlined />
            <div class="q-gutter-sm">
              <q-btn label="Simpan" type="submit" color="primary" />
              <q-btn label="Batal" flat color="negative" v-close-popup @click="resetForm" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { usePenempatanStore } from 'stores/pkl/penempatan'
import { useMitraStore } from 'stores/pkl/mitra'
import { useSiswaStore } from 'stores/kesiswaan/siswa' // asumsi kita punya store siswa

const $q = useQuasar()
const penempatanStore = usePenempatanStore()
const mitraStore = useMitraStore()
const siswaStore = useSiswaStore()

const formDialog = ref(false)
const editMode = ref(false)
const editId = ref(null)
const form = ref({
  siswa_id: null,
  mitra_id: null,
  tanggal_mulai: '',
  tanggal_selesai: '',
  pembimbing_industri: '',
  catatan: ''
})

const columns = [
  { name: 'siswa_nama', label: 'Siswa', field: 'siswa_nama', align: 'left' },
  { name: 'mitra_nama', label: 'Mitra', field: 'mitra_nama', align: 'left' },
  { name: 'tanggal_mulai', label: 'Mulai', field: 'tanggal_mulai', align: 'left' },
  { name: 'tanggal_selesai', label: 'Selesai', field: 'tanggal_selesai', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'aksi', label: 'Aksi', field: 'aksi', align: 'center' }
]

const siswaOptions = computed(() => siswaStore.list.map(s => ({ id: s.id, nama: s.nama })))
const mitraOptions = computed(() => mitraStore.list)

function openForm(row = null) {
  if (row) {
    editMode.value = true
    editId.value = row.id
    form.value = {
      siswa_id: row.siswa_id,
      mitra_id: row.mitra_id,
      tanggal_mulai: row.tanggal_mulai,
      tanggal_selesai: row.tanggal_selesai,
      pembimbing_industri: row.pembimbing_industri,
      catatan: row.catatan
    }
  } else {
    editMode.value = false
    resetForm()
  }
  formDialog.value = true
}

function resetForm() {
  form.value = {
    siswa_id: null,
    mitra_id: null,
    tanggal_mulai: '',
    tanggal_selesai: '',
    pembimbing_industri: '',
    catatan: ''
  }
}

function submitPenempatan() {
  if (!form.value.siswa_id || !form.value.mitra_id) {
    $q.notify({ type: 'negative', message: 'Lengkapi data siswa dan mitra' })
    return
  }
  const siswa = siswaStore.list.find(s => s.id === form.value.siswa_id)
  const mitra = mitraStore.list.find(m => m.id === form.value.mitra_id)
  const data = {
    siswa_id: form.value.siswa_id,
    siswa_nama: siswa.nama,
    mitra_id: form.value.mitra_id,
    mitra_nama: mitra.nama,
    tanggal_mulai: form.value.tanggal_mulai,
    tanggal_selesai: form.value.tanggal_selesai,
    pembimbing_industri: form.value.pembimbing_industri,
    catatan: form.value.catatan,
    status: 'ditempatkan'
  }
  if (editMode.value) {
    penempatanStore.update(editId.value, data)
    $q.notify({ type: 'positive', message: 'Penempatan diupdate' })
  } else {
    penempatanStore.tambah(data)
    $q.notify({ type: 'positive', message: 'Penempatan ditambahkan' })
  }
  formDialog.value = false
  resetForm()
}

function confirmDelete(id) {
  $q.dialog({ title: 'Konfirmasi', message: 'Hapus penempatan ini?', cancel: true }).onOk(() => {
    penempatanStore.hapus(id)
    $q.notify({ type: 'positive', message: 'Data dihapus' })
  })
}

function generateSurat(penempatan) {
  // Simulasi cetak surat penempatan
  const content = `
    SURAT PENEMPATAN PKL
    Kepada Yth. ${penempatan.mitra_nama}
    Kami menempatkan siswa: ${penempatan.siswa_nama}
    Periode: ${penempatan.tanggal_mulai} s.d. ${penempatan.tanggal_selesai}
    Pembimbing Industri: ${penempatan.pembimbing_industri}
  `
  const win = window.open('', '_blank')
  win.document.write(`<html><head><title>Surat Penempatan</title></head><body><pre>${content}</pre></body></html>`)
  win.document.close()
  win.print()
}

onMounted(() => {
  mitraStore.loadData()
  siswaStore.loadData()
  penempatanStore.loadData()
})
</script>
