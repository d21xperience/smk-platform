<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row items-center q-mb-md">
      <q-btn flat round dense icon="arrow_back" @click="$router.push('/kepegawaian/skp')" class="q-mr-sm" />
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Daftar Penilaian Pegawai</h1>
        <p class="text-caption text-grey-7 q-mb-none">Periode ID: {{ periodeId }}</p>
      </div>
      <q-space />
      <q-btn color="primary" icon="refresh" label="Refresh" @click="fetchData" :loading="isLoading" />
    </div>

    <q-linear-progress v-if="isLoading" indeterminate color="primary" class="q-mb-md" />

    <q-card class="shadow-1">
      <q-card-section>
        <q-table :rows="penilaians" :columns="columns" row-key="id" flat bordered dense :loading="isLoading">
          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center">
              <q-badge :color="getStatusColor(props.row.status)" class="q-pa-sm">
                {{ props.row.status }}
              </q-badge>
            </q-td>
          </template>
          <template v-slot:body-cell-nilaiAkhir="props">
            <q-td :props="props" class="text-center">
              <div v-if="props.row.nilaiAkhir > 0" class="text-weight-bold">
                {{ props.row.nilaiAkhir }}
              </div>
              <span v-else class="text-grey-5">-</span>
            </q-td>
          </template>
          <template v-slot:body-cell-predikat="props">
            <q-td :props="props" class="text-center">
              <q-badge v-if="props.row.predikatAkhir" :color="getPredikatColor(props.row.predikatAkhir)" class="q-pa-sm">
                {{ props.row.predikatAkhir }}
              </q-badge>
              <span v-else class="text-grey-5">-</span>
            </q-td>
          </template>
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center">
              <q-btn 
                flat 
                round 
                dense 
                icon="edit" 
                color="primary" 
                @click="bukaFormPenilaian(props.row.id)" 
                class="q-mr-sm"
              >
                <q-tooltip>Isi / Lihat Penilaian</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { api } from '@/boot/axios'

const $q = useQuasar()
const route = useRoute()
const router = useRouter()

const periodeId = ref(route.params.id)
const penilaians = ref([])
const isLoading = ref(false)

const columns = [
  { name: 'teacherName', label: 'Nama Guru / Pegawai', field: 'teacherName', align: 'left', sortable: true },
  { name: 'status', label: 'Status', field: 'status', align: 'center', sortable: true },
  { name: 'nilaiAkhir', label: 'Nilai Akhir', field: 'nilaiAkhir', align: 'center' },
  { name: 'predikat', label: 'Predikat', field: 'predikatAkhir', align: 'center' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const getStatusColor = (status) => {
  switch (status) {
    case 'PENDING': return 'orange-6'
    case 'SUBMITTED': return 'blue-6'
    case 'FINAL': return 'green-6'
    case 'RETURNED': return 'red-6'
    default: return 'grey-6'
  }
}

const getPredikatColor = (predikat) => {
  switch (predikat) {
    case 'SANGAT_BAIK': return 'green-8'
    case 'BAIK': return 'green-6'
    case 'BUTUH_BIMBINGAN': return 'orange-6'
    case 'KURANG': return 'red-6'
    case 'BURUK': return 'red-8'
    default: return 'grey-6'
  }
}

const fetchData = async () => {
  isLoading.value = true
  try {
    const res = await api.get(`/kepegawaian/skp/periodes/${periodeId.value}/penilaians`)
    penilaians.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal memuat daftar penilaian' })
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

const bukaFormPenilaian = (penilaianId) => {
  router.push({ name: 'skp-form', params: { id: penilaianId } })
}

onMounted(fetchData)
</script>
