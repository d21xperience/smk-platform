<template>
  <q-page class="q-pa-md">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Buku Induk Pegawai</h1>
        <p class="text-caption text-grey-7 q-mb-none">Pusat data administratif pegawai (Guru & Tendik)</p>
      </div>
      <div class="row q-gutter-sm">
        <q-btn color="secondary" icon="download" label="Export Excel" @click="exportExcel" outline />
        <q-btn color="primary" icon="refresh" label="Refresh" @click="fetchData" :loading="isLoading" />
      </div>
    </div>

    <q-card class="shadow-1">
      <q-card-section class="row items-center q-pb-none">
        <q-input v-model="filterKeyword" label="Cari Nama / NIK / NIP" outlined dense class="col-12 col-md-4" clearable>
          <template v-slot:prepend><q-icon name="search" /></template>
        </q-input>
        <q-space />
        <q-select v-model="filterKelengkapan" :options="['Semua', 'Lengkap', 'Cukup', 'Belum Lengkap']"
          label="Filter Kelengkapan" outlined dense class="col-12 col-md-3" />
      </q-card-section>

      <q-card-section>
        <q-table :rows="filteredData" :columns="columns" row-key="ptkId" flat bordered dense :loading="isLoading"
          :pagination="{ rowsPerPage: 10 }">
          <template v-slot:body-cell-foto="props">
            <q-td :props="props" class="text-center">
              <q-avatar size="40px">
                <img :src="props.row.fotoUrl || 'https://cdn.quasar.dev/img/boy-avatar.png'"
                  @error="(e) => e.target.src = 'https://cdn.quasar.dev/img/boy-avatar.png'">
              </q-avatar>
            </q-td>
          </template>
          <template v-slot:body-cell-kelengkapan="props">
            <q-td :props="props" class="text-center">
              <q-badge :color="getBadgeColor(props.row.statusKelengkapan)" class="q-pa-sm">
                {{ props.row.statusKelengkapan }} ({{ props.row.skorKelengkapan }}%)
              </q-badge>
            </q-td>
          </template>
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center">
              <q-btn flat round dense icon="visibility" color="primary"
                @click="$router.push(`/kepegawaian/buku-induk/${props.row.ptkId}`)" class="q-mr-sm">
                <q-tooltip>Lihat Detail</q-tooltip>
              </q-btn>
              <q-btn flat round dense icon="print" color="secondary" @click="cetakPDF(props.row)">
                <q-tooltip>Cetak PDF</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { api } from '@/boot/axios'

const $q = useQuasar()
const rawData = ref([])
const isLoading = ref(false)
const filterKeyword = ref('')
const filterKelengkapan = ref('Semua')

const columns = [
  { name: 'foto', label: 'Foto', align: 'center' },
  { name: 'name', label: 'Nama Lengkap & Gelar', field: 'name', align: 'left', sortable: true },
  { name: 'nik', label: 'NIK', field: 'nik', align: 'left' },
  { name: 'nip', label: 'NIP / NUPTK', field: 'nip', align: 'left' },
  { name: 'jenisPegawai', label: 'Jenis', field: 'jenisPegawai', align: 'center' },
  { name: 'statusPegawai', label: 'Status', field: 'statusPegawai', align: 'center' },
  { name: 'kelengkapan', label: 'Kelengkapan Data', align: 'center', sortable: true },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

const filteredData = computed(() => {
  return rawData.value.filter(item => {
    const matchKeyword = item.name.toLowerCase().includes(filterKeyword.value.toLowerCase()) ||
      (item.nik && item.nik.includes(filterKeyword.value)) ||
      (item.nip && item.nip.includes(filterKeyword.value))

    const matchKelengkapan = filterKelengkapan.value === 'Semua' || item.statusKelengkapan === filterKelengkapan.value

    return matchKeyword && matchKelengkapan
  })
})

const getBadgeColor = (status) => {
  if (status === 'Lengkap') return 'green-6'
  if (status === 'Cukup') return 'orange-6'
  return 'red-6'
}

const fetchData = async () => {
  isLoading.value = true
  try {
    const res = await api.get('/teachers/buku-induk')
    rawData.value = res.data || []
  } catch (err) {
    $q.notify({ type: 'negative', message: `Gagal memuat data Buku Induk ${err}` })
  } finally {
    isLoading.value = false
  }
}

const cetakPDF = (row) => {
  $q.notify({ type: 'info', message: `Fitur Cetak PDF untuk ${row.name} akan segera tersedia.` })
}

const exportExcel = () => {
  $q.notify({ type: 'info', message: 'Fitur Export Excel akan segera tersedia.' })
}

onMounted(fetchData)
</script>
