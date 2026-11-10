<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row items-center q-mb-md">
      <q-btn flat round dense icon="arrow_back" @click="$router.push('/kepegawaian/sk-pembagian-tugas')" class="q-mr-sm" />
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Buat SK Pembagian Tugas</h1>
        <p class="text-caption text-grey-7 q-mb-none">Wizard Input Beban Mengajar Guru</p>
      </div>
    </div>

    <q-card class="shadow-1">
      <q-stepper v-model="step" ref="stepper" color="primary" animated header-nav>
        <!-- STEP 1: PERIODE -->
        <q-step :name="1" title="Periode" icon="event" :done="step > 1">
          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6">
              <q-select 
                v-model="form.tahunAjaran" 
                :options="tahunAjaranOptions" 
                label="Tahun Ajaran" 
                outlined 
                emit-value 
                map-options
                :rules="[v => !!v || 'Wajib diisi']" 
              />
            </div>
            <div class="col-12 col-md-6">
              <q-select 
                v-model="form.semester" 
                :options="semesterOptions" 
                emit-value 
                map-options 
                label="Semester" 
                outlined 
                :rules="[v => v !== null || 'Wajib diisi']" 
              />
            </div>
          </div>
          <q-stepper-navigation>
            <q-btn @click="step = 2" color="primary" label="Lanjut ke Penugasan" :disable="!form.tahunAjaran || form.semester === null" />
          </q-stepper-navigation>
        </q-step>

        <!-- STEP 2: PENUGASAN MENGAJAR -->
        <q-step :name="2" title="Penugasan" icon="assignment" :done="step > 2">
          <div v-for="(assign, index) in form.assignments" :key="index" class="q-mb-md q-pa-md bg-grey-1 rounded-borders">
            <div class="row items-center q-mb-sm">
              <div class="text-subtitle2 text-primary">Penugasan #{{ index + 1 }}</div>
              <q-space />
              <q-btn flat round dense icon="delete" color="negative" @click="removeAssignment(index)" />
            </div>
            
            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-4">
                <q-select v-model="assign.ptkId" :options="guruOptions" label="Pilih Guru" option-value="id" option-label="nama" emit-value map-options outlined dense :rules="[v => !!v || 'Wajib']" />
              </div>
              <div class="col-12 col-md-4">
                <q-select v-model="assign.mataPelajaranId" :options="mapelOptions" label="Mata Pelajaran" option-value="id" option-label="nama" emit-value map-options outlined dense :rules="[v => !!v || 'Wajib']" />
              </div>
              <div class="col-12 col-md-2">
                <q-input v-model.number="assign.jamPerMinggu" type="number" label="JTM" outlined dense :rules="[v => !!v || 'Wajib']" />
              </div>
            </div>
            <div class="row q-mt-sm">
              <div class="col-12">
                <q-select v-model="assign.rombonganBelajarIds" :options="kelasOptions" label="Pilih Kelas (Bisa lebih dari satu)" option-value="id" option-label="nama" emit-value map-options multiple use-chips outlined dense :rules="[v => v.length > 0 || 'Wajib dipilih']" />
              </div>
            </div>
          </div>

          <q-btn color="secondary" icon="add" label="Tambah Baris Penugasan" @click="addAssignment" flat class="q-mb-md" />

          <q-stepper-navigation>
            <q-btn @click="step = 3" color="primary" label="Lanjut ke Tugas Tambahan" :disable="form.assignments.length === 0" />
            <q-btn flat @click="step = 1" color="primary" label="Kembali" class="q-ml-sm" />
          </q-stepper-navigation>
        </q-step>

        <!-- STEP 3: TUGAS TAMBAHAN -->
        <q-step :name="3" title="Tugas Tambahan" icon="work" :done="step > 3">
          <div v-for="(tugas, index) in form.tugasTambahan" :key="index" class="q-mb-md q-pa-md bg-grey-1 rounded-borders">
            <div class="row items-center q-mb-sm">
              <div class="text-subtitle2 text-primary">Tugas Tambahan #{{ index + 1 }}</div>
              <q-space />
              <q-btn flat round dense icon="delete" color="negative" @click="removeTugasTambahan(index)" />
            </div>
            <div class="row q-col-gutter-sm">
              <div class="col-12 col-md-4">
                <q-select v-model="tugas.ptkId" :options="guruOptions" label="Pilih Guru" option-value="id" option-label="nama" emit-value map-options outlined dense />
              </div>
              <div class="col-12 col-md-4">
                <q-select v-model="tugas.jenisTugas" :options="['Wali Kelas', 'Pembina Ekskul', 'Kepala Lab', 'Lainnya']" label="Jenis Tugas" outlined dense />
              </div>
              <div class="col-12 col-md-4">
                <q-input v-model="tugas.keterangan" label="Keterangan" outlined dense />
              </div>
            </div>
          </div>

          <q-btn color="secondary" icon="add" label="Tambah Tugas Tambahan" @click="addTugasTambahan" flat class="q-mb-md" />

          <q-stepper-navigation>
            <q-btn @click="step = 4" color="primary" label="Lanjut ke Review" />
            <q-btn flat @click="step = 2" color="primary" label="Kembali" class="q-ml-sm" />
          </q-stepper-navigation>
        </q-step>

        <!-- STEP 4: REVIEW & SUBMIT -->
        <q-step :name="4" title="Review" icon="check_circle">
          <div class="text-h6 text-center text-primary q-mb-md">Ringkasan Pembagian Tugas</div>
          <q-markup-table flat bordered>
            <thead>
              <tr>
                <th>No</th>
                <th>Nama Guru</th>
                <th>Mata Pelajaran</th>
                <th>Jumlah Kelas</th>
                <th>JTM / Kelas</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(assign, i) in form.assignments" :key="i">
                <td>{{ i + 1 }}</td>
                <td>{{ getGuruName(assign.ptkId) }}</td>
                <td>{{ getMapelName(assign.mataPelajaranId) }}</td>
                <td>{{ assign.rombonganBelajarIds.length }}</td>
                <td>{{ assign.jamPerMinggu }}</td>
              </tr>
            </tbody>
          </q-markup-table>

          <q-stepper-navigation>
            <q-btn @click="submitWizard" color="positive" label="Simpan Draft SK" :loading="isSubmitting" />
            <q-btn flat @click="step = 3" color="primary" label="Kembali" class="q-ml-sm" />
          </q-stepper-navigation>
        </q-step>
      </q-stepper>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useQuasar } from 'quasar'
import { useRouter } from 'vue-router'
import { api } from '@/boot/axios'
import { useContextStore } from '@/stores/contextStore'

const $q = useQuasar()
const router = useRouter()
const contextStore = useContextStore()

const step = ref(1)
const isSubmitting = ref(false)

// Ambil nilai default dari contextStore (Tahun Ajaran aktif)
const form = ref({
  tahunAjaran: contextStore.academicYear || '2026/2027',
  semester: contextStore.semester || 1,
  assignments: [],
  tugasTambahan: []
})

const teachers = ref([])
const subjects = ref([])
const classes = ref([])

const guruOptions = computed(() => teachers.value.map(t => ({ id: t.id, nama: t.name })))
const mapelOptions = computed(() => subjects.value.map(s => ({ id: s.id, nama: s.name })))
const kelasOptions = computed(() => classes.value.map(c => ({ id: c.rombongan_belajar_id || c.ID, nama: c.nm_kelas || c.Name })))

// Opsi Tahun Ajaran Dinamis (dari contextStore)
const tahunAjaranOptions = computed(() => {
  const years = new Set(contextStore.availablePeriods.map(p => `${p.tahun_ajaran_id}/${parseInt(p.tahun_ajaran_id) + 1}`))
  return Array.from(years).sort()
})

// Opsi Semester (0 = Semua Semester)
const semesterOptions = [
  { label: 'Semua Semester (1 Tahun Penuh)', value: 0 },
  { label: 'Ganjil', value: 1 },
  { label: 'Genap', value: 2 }
]

const getGuruName = (id) => {
  const g = teachers.value.find(t => t.id === id)
  return g ? g.name : id
}
const getMapelName = (id) => {
  const m = subjects.value.find(s => s.id === id)
  return m ? m.name : id
}

const addAssignment = () => {
  form.value.assignments.push({ ptkId: null, mataPelajaranId: null, rombonganBelajarIds: [], jamPerMinggu: null })
}

const removeAssignment = (index) => {
  form.value.assignments.splice(index, 1)
}

const addTugasTambahan = () => {
  form.value.tugasTambahan.push({ ptkId: null, jenisTugas: 'Wali Kelas', keterangan: '' })
}

const removeTugasTambahan = (index) => {
  form.value.tugasTambahan.splice(index, 1)
}

const fetchOptions = async () => {
  try {
    const [resTeachers, resSubjects, resClasses] = await Promise.all([
      api.get('/teachers'),
      api.get('/subjects'),
      api.get('/schedules/options') // Mengambil data classes dari endpoint jadwal
    ])
    teachers.value = resTeachers.data.data || []
    subjects.value = resSubjects.data.data || []
    classes.value = resClasses.data.classes || []
  } catch (e) {
    $q.notify({ type: 'negative', message: 'Gagal memuat opsi data' })
    console.error(e)
  }
}

const submitWizard = async () => {
  isSubmitting.value = true
  try {
    const payload = {
      tahunAjaran: form.value.tahunAjaran,
      semester: form.value.semester,
      assignments: form.value.assignments,
      tugasTambahan: form.value.tugasTambahan
    }
    await api.post('/kepegawaian/sk-pembagian-tugas', payload)
    $q.notify({ type: 'positive', message: 'Draft SK berhasil dibuat!' })
    router.push('/kepegawaian/sk-pembagian-tugas')
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || 'Gagal menyimpan draft SK' })
    console.error(e)
  } finally {
    isSubmitting.value = false
  }
}

onMounted(() => {
  fetchOptions()
  addAssignment() // Inisialisasi 1 baris kosong
})
</script>
