<template>
  <q-page class="q-pa-md">
    <div class="row justify-between q-mb-md">
      <div class="text-h5">Beban Mengajar Guru</div>
      <q-btn color="primary" icon="add" label="Tambah Jadwal" @click="openDialog()" />
    </div>

    <q-table :rows="jadwalStore.list" :columns="columns" row-key="id" flat bordered dense>
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props">
          <q-btn flat round dense icon="edit" color="info" @click="openDialog(props.row)" class="q-mr-sm" />
          <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row.id)" />
        </q-td>
      </template>
    </q-table>

    <!-- Dialog form -->
    <q-dialog v-model="dialog" persistent>
      <q-card style="min-width: 600px">
        <q-card-section>
          <div class="text-h6">{{ editMode ? 'Edit Jadwal' : 'Tambah Jadwal' }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitForm" class="q-gutter-md">
            <div class="row q-col-gutter-md">
              <div class="col-12 col-md-6">
                <q-select v-model="form.guru_id" :options="guruOptions" label="Guru" option-value="id"
                  option-label="nama" outlined :rules="[v => !!v]" @update:model-value="setGuruNama" />
              </div>
              <div class="col-12 col-md-6">
                <q-input v-model="form.mapel" label="Mata Pelajaran" outlined :rules="[v => !!v]" />
              </div>
              <div class="col-12 col-md-6">
                <q-input v-model="form.kelas" label="Kelas" outlined :rules="[v => !!v]" />
              </div>
              <div class="col-12 col-md-6">
                <q-select v-model="form.hari" :options="hariOptions" label="Hari" outlined />
              </div>
              <div class="col-12 col-md-6">
                <q-input v-model="form.jam_mulai" label="Jam Mulai" type="time" outlined />
              </div>
              <div class="col-12 col-md-6">
                <q-input v-model="form.jam_selesai" label="Jam Selesai" type="time" outlined />
              </div>
              <div class="col-12 col-md-6">
                <q-input v-model="form.ruang" label="Ruang" outlined />
              </div>
            </div>
            <div v-if="conflicts.length" class="text-negative">
              <q-icon name="warning" /> Konflik terdeteksi:
              <ul>
                <li v-for="(c, idx) in conflicts" :key="idx">{{ c }}</li>
              </ul>
            </div>
            <div class="q-gutter-sm">
              <q-btn label="Simpan" type="submit" color="primary" :disable="conflicts.length > 0" />
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
import { useJadwalStore } from '@/stores/kepegawaian/jadwal'
import { useGuruStore } from '@/stores/kepegawaian/guru'

const $q = useQuasar()
const jadwalStore = useJadwalStore()
const guruStore = useGuruStore()

const dialog = ref(false)
const editMode = ref(false)
const editId = ref(null)
const conflicts = ref([])

const form = ref({
  guru_id: null,
  guru_nama: '',
  mapel: '',
  kelas: '',
  hari: 'Senin',
  jam_mulai: '07:30',
  jam_selesai: '09:00',
  ruang: ''
})

const hariOptions = ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu']
const columns = [
  { name: 'guru_nama', label: 'Guru', field: 'guru_nama', align: 'left' },
  { name: 'mapel', label: 'Mata Pelajaran', field: 'mapel', align: 'left' },
  { name: 'kelas', label: 'Kelas', field: 'kelas', align: 'left' },
  { name: 'hari', label: 'Hari', field: 'hari', align: 'left' },
  { name: 'jam_mulai', label: 'Jam Mulai', field: 'jam_mulai', align: 'center' },
  { name: 'jam_selesai', label: 'Jam Selesai', field: 'jam_selesai', align: 'center' },
  { name: 'ruang', label: 'Ruang', field: 'ruang', align: 'left' },
  { name: 'aksi', label: 'Aksi', field: 'aksi', align: 'center' }
]

const guruOptions = computed(() => guruStore.list.map(g => ({ id: g.id, nama: g.nama })))

function setGuruNama(val) {
  const guru = guruStore.list.find(g => g.id === val)
  form.value.guru_nama = guru ? guru.nama : ''
}

// function cekKonflik() {
//   if (!form.value.guru_id || !form.value.hari || !form.value.jam_mulai || !form.value.jam_selesai) return
//   const jadwalBaru = {
//     guru_id: form.value.guru_id,
//     guru_nama: form.value.guru_nama,
//     mapel: form.value.mapel,
//     kelas: form.value.kelas,
//     hari: form.value.hari,
//     jam_mulai: form.value.jam_mulai,
//     jam_selesai: form.value.jam_selesai,
//     ruang: form.value.ruang
//   }
//   conflicts.value = jadwalStore.cekKonflik(jadwalBaru, editMode.value ? editId.value : null)
// }

function openDialog(row = null) {
  if (row) {
    editMode.value = true
    editId.value = row.id
    form.value = { ...row }
  } else {
    editMode.value = false
    editId.value = null
    resetForm()
  }
  conflicts.value = []
  dialog.value = true
}

function resetForm() {
  form.value = {
    guru_id: null,
    guru_nama: '',
    mapel: '',
    kelas: '',
    hari: 'Senin',
    jam_mulai: '07:30',
    jam_selesai: '09:00',
    ruang: ''
  }
  conflicts.value = []
}

function submitForm() {
  if (conflicts.value.length) {
    $q.notify({ type: 'negative', message: 'Ada konflik jadwal, tidak bisa disimpan' })
    return
  }
  if (editMode.value) {
    jadwalStore.update(editId.value, form.value)
    $q.notify({ type: 'positive', message: 'Jadwal diupdate' })
  } else {
    jadwalStore.tambah({ ...form.value })
    $q.notify({ type: 'positive', message: 'Jadwal ditambahkan' })
  }
  dialog.value = false
  resetForm()
}

function confirmDelete(id) {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Hapus jadwal ini?',
    cancel: true,
    persistent: true
  }).onOk(() => {
    jadwalStore.hapus(id)
    $q.notify({ type: 'positive', message: 'Jadwal dihapus' })
  })
}

onMounted(() => {
  guruStore.loadData()
  jadwalStore.loadData()
})
</script>
