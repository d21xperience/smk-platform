<template>
  <q-page class="q-pa-md">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Kinerja GTK (SKP)</h1>
        <p class="text-caption text-grey-7 q-mb-none">Penilaian Sasaran Kinerja Pegawai</p>
      </div>
      <q-btn color="primary" icon="add" label="Tambah Penilaian" @click="openFormDialog" />
    </div>

    <q-card class="shadow-1">
      <q-card-section class="row items-center q-pb-none">
        <q-btn color="secondary" icon="refresh" label="Refresh" @click="fetchPerformances" :loading="isLoading" flat
          class="q-ml-auto" />
      </q-card-section>

      <q-card-section>
        <q-table :rows="performances" :columns="columns" row-key="id" flat bordered dense :loading="isLoading">
          <template v-slot:body-cell-nilaiSkp="props">
            <q-td :props="props" class="text-center">
              <div class="text-weight-bold">{{ props.row.nilaiSkp }}</div>
            </q-td>
          </template>
          <template v-slot:body-cell-predikat="props">
            <q-td :props="props" class="text-center">
              <q-badge :color="getPredikatColor(props.row.predikat)" class="q-pa-sm">
                {{ props.row.predikat }}
              </q-badge>
            </q-td>
          </template>
          <template v-slot:body-cell-catatan="props">
            <q-td :props="props">
              <div class="ellipsis" style="max-width: 200px">{{ props.row.catatan || '-' }}</div>
              <q-tooltip v-if="props.row.catatan">{{ props.row.catatan }}</q-tooltip>
            </q-td>
          </template>
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center">
              <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row.id)">
                <q-tooltip>Hapus</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- Dialog Form Penilaian -->
    <q-dialog v-model="formDialog" persistent>
      <q-card style="min-width: 500px">
        <q-card-section class="row items-center bg-primary text-white">
          <div class="text-h6">Form Penilaian Kinerja</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section class="q-gutter-y-md">
          <q-form @submit="submitPerformance">
            <q-select v-model="form.ptkId" :options="guruOptions" label="Pilih Guru" option-value="id"
              option-label="nama" emit-value map-options outlined dense :rules="[v => !!v || 'Wajib dipilih']" />

            <div class="row q-col-gutter-sm q-mt-sm">
              <div class="col-6">
                <q-input v-model="form.tahun" label="Tahun" outlined dense mask="####"
                  :rules="[v => !!v || 'Wajib diisi']" />
              </div>
              <div class="col-6">
                <q-select v-model="form.semester" :options="['Ganjil', 'Genap']" label="Semester" outlined dense
                  :rules="[v => !!v || 'Wajib diisi']" />
              </div>
            </div>

            <div class="row q-col-gutter-sm q-mt-sm">
              <div class="col-6">
                <q-input v-model.number="form.nilaiSkp" label="Nilai SKP" type="number" step="0.01" outlined dense
                  :rules="[v => !!v || 'Wajib diisi']" @update:model-value="autoPredikat" />
              </div>
              <div class="col-6">
                <q-select v-model="form.predikat" :options="['Sangat Baik', 'Baik', 'Cukup', 'Kurang']" label="Predikat"
                  outlined dense :rules="[v => !!v || 'Wajib diisi']" />
              </div>
            </div>

            <q-input v-model="form.evaluator" label="Nama Evaluator (Penilai)" outlined dense class="q-mt-sm" />
            <q-input v-model="form.catatan" label="Catatan Evaluasi" type="textarea" rows="2" outlined dense
              class="q-mt-sm" />

            <div class="row justify-end q-mt-md">
              <q-btn label="Batal" flat color="negative" v-close-popup class="q-mr-sm" />
              <q-btn label="Simpan" type="submit" color="primary" :loading="isSubmitting" />
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
import { api } from '@/boot/axios'
import { useGuruStore } from '@/stores/kepegawaian/guru'

const $q = useQuasar()
const guruStore = useGuruStore()

const performances = ref([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const formDialog = ref(false)

const form = ref({
  ptkId: null,
  tahun: new Date().getFullYear().toString(),
  semester: 'Ganjil',
  nilaiSkp: null,
  predikat: 'Baik',
  evaluator: 'Kepala Sekolah',
  catatan: ''
})

const columns = [
  { name: 'teacherName', label: 'Nama Guru', field: 'teacherName', align: 'left', sortable: true },
  { name: 'tahun', label: 'Tahun', field: 'tahun', align: 'center' },
  { name: 'semester', label: 'Semester', field: 'semester', align: 'center' },
  { name: 'nilaiSkp', label: 'Nilai SKP', field: 'nilaiSkp', align: 'center' },
  { name: 'predikat', label: 'Predikat', field: 'predikat', align: 'center' },
  { name: 'evaluator', label: 'Evaluator', field: 'evaluator', align: 'left' },
  { name: 'catatan', label: 'Catatan', field: 'catatan', align: 'left' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const guruOptions = computed(() => guruStore.list.map(g => ({ id: g.id, nama: g.name })))

const getPredikatColor = (predikat) => {
  switch (predikat.toLowerCase()) {
    case 'sangat baik': return 'green-6'
    case 'baik': return 'blue-6'
    case 'cukup': return 'orange-6'
    case 'kurang': return 'red-6'
    default: return 'grey-6'
  }
}

const fetchPerformances = async () => {
  isLoading.value = true
  try {
    const res = await api.get('/teachers/performances')
    performances.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: `Gagal memuat data kinerja ${err}` })
  } finally {
    isLoading.value = false
  }
}

const autoPredikat = (val) => {
  if (val >= 90) form.value.predikat = 'Sangat Baik'
  else if (val >= 80) form.value.predikat = 'Baik'
  else if (val >= 70) form.value.predikat = 'Cukup'
  else if (val > 0) form.value.predikat = 'Kurang'
}

const openFormDialog = () => {
  form.value = { ptkId: null, tahun: new Date().getFullYear().toString(), semester: 'Ganjil', nilaiSkp: null, predikat: 'Baik', evaluator: 'Kepala Sekolah', catatan: '' }
  formDialog.value = true
}

const submitPerformance = async () => {
  isSubmitting.value = true
  try {
    await api.post('/teachers/performances', form.value)
    $q.notify({ type: 'positive', message: 'Penilaian berhasil disimpan' })
    formDialog.value = false
    await fetchPerformances()
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Gagal menyimpan penilaian' })
  } finally {
    isSubmitting.value = false
  }
}

const confirmDelete = (id) => {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Hapus data penilaian kinerja ini?',
    cancel: true
  }).onOk(async () => {
    try {
      await api.delete(`/teachers/performances/${id}`)
      $q.notify({ type: 'positive', message: 'Data berhasil dihapus' })
      await fetchPerformances()
    } catch (err) {
      $q.notify({ type: 'negative', message: `Gagal menghapus data ${err}` })
    }
  })
}

onMounted(async () => {
  await guruStore.fetchTeachers()
  await fetchPerformances()
})
</script>
