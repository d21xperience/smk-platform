<template>
  <q-page class="q-pa-md">
    <div class="row justify-between q-mb-md">
      <div class="text-h5">Struktur Kurikulum (CP/TP, KKM, Alokasi Jam)</div>
      <q-btn color="primary" icon="add" label="Tambah" @click="openForm" />
    </div>

    <q-table :rows="strukturStore.list" :columns="columns" row-key="id" flat bordered dense>
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props">
          <q-btn flat round dense icon="edit" color="info" @click="openForm(props.row)" class="q-mr-sm" />
          <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row.id)" />
        </q-td>
      </template>
    </q-table>

    <q-dialog v-model="formDialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section>
          <div class="text-h6">{{ editMode ? 'Edit Struktur' : 'Tambah Struktur' }}</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="submitForm" class="q-gutter-md">
            <q-select v-model="form.kelas_id" :options="kelasOptions" label="Kelas" option-value="id"
              option-label="nama" outlined :rules="[v => !!v]" />
            <q-select v-model="form.mapel_id" :options="mapelOptions" label="Mata Pelajaran" option-value="id"
              option-label="nama" outlined :rules="[v => !!v]" />
            <q-input v-model.number="form.kkm" label="KKM" type="number" min="0" max="100" outlined
              :rules="[v => v >= 0 && v <= 100]" />
            <q-input v-model.number="form.alokasi_jam" label="Alokasi Jam per Minggu" type="number" min="1" outlined />
            <q-select v-model="form.guru_id" :options="guruOptions" label="Guru Pengajar" option-value="id"
              option-label="nama" outlined />
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
import { useStrukturStore } from 'stores/kurikulum/struktur'
import { useKelasStore } from '@/stores/kesiswaan/kelasStore'
import { useMapelStore } from 'stores/akademik/mapel'
import { useGuruStore } from 'stores/kepegawaian/guru'

const $q = useQuasar()
const strukturStore = useStrukturStore()
const kelasStore = useKelasStore()
const mapelStore = useMapelStore()
const guruStore = useGuruStore()

const formDialog = ref(false)
const editMode = ref(false)
const editId = ref(null)
const form = ref({ kelas_id: null, mapel_id: null, kkm: 75, alokasi_jam: 4, guru_id: null })

const columns = [
  { name: 'kelas_nama', label: 'Kelas', field: 'kelas_nama' },
  { name: 'mapel_nama', label: 'Mata Pelajaran', field: 'mapel_nama' },
  { name: 'kkm', label: 'KKM', field: 'kkm' },
  { name: 'alokasi_jam', label: 'Jam/Minggu', field: 'alokasi_jam' },
  { name: 'guru_nama', label: 'Guru', field: 'guru_nama' },
  { name: 'aksi', label: 'Aksi', field: 'aksi' }
]

const kelasOptions = computed(() => kelasStore.list)
const mapelOptions = computed(() => mapelStore.list)
const guruOptions = computed(() => guruStore.list)

function openForm(row = null) {
  if (row) {
    editMode.value = true
    editId.value = row.id
    form.value = { ...row }
  } else {
    editMode.value = false
    resetForm()
  }
  formDialog.value = true
}
function resetForm() {
  form.value = { kelas_id: null, mapel_id: null, kkm: 75, alokasi_jam: 4, guru_id: null }
}
function submitForm() {
  const kelas = kelasStore.list.find(k => k.id === form.value.kelas_id)
  const mapel = mapelStore.list.find(m => m.id === form.value.mapel_id)
  const guru = guruStore.list.find(g => g.id === form.value.guru_id)
  const data = {
    ...form.value,
    kelas_nama: kelas?.nama,
    mapel_nama: mapel?.nama,
    guru_nama: guru?.nama
  }
  if (editMode.value) {
    strukturStore.update(editId.value, data)
    $q.notify({ type: 'positive', message: 'Data diupdate' })
  } else {
    strukturStore.tambah(data)
    $q.notify({ type: 'positive', message: 'Data ditambahkan' })
  }
  formDialog.value = false
}
function confirmDelete(id) {
  $q.dialog({ title: 'Konfirmasi', message: 'Hapus data ini?', cancel: true }).onOk(() => {
    strukturStore.hapus(id)
    $q.notify({ type: 'positive', message: 'Data dihapus' })
  })
}
onMounted(() => {
  strukturStore.loadData()
  kelasStore.loadData()
  mapelStore.loadData()
  guruStore.loadData()
})
</script>
