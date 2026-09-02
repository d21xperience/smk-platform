<template>
  <q-dialog :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)" full-width max-width="800px" position="top">
    <q-card class="q-pa-md">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6 text-primary text-weight-bold">Pencarian Data Siswa / Alumni (TU Mode)</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <!-- Form Filter Pencarian -->
      <q-card-section class="q-gutter-md row items-center">
        <div class="col-12 col-sm-3">
          <q-select
            v-model="state.filter.status"
            :options="['Siswa Aktif', 'Alumni']"
            label="Status"
            outlined
            dense
            emit-value
            map-options
          />
        </div>

        <!-- Filter Tahun Lulus (Hanya untuk Alumni) -->
        <div class="col-12 col-sm-4" v-if="state.filter.status === 'Alumni'">
          <div class="text-caption text-grey-7 q-mb-xs">Tahun Lulus: {{ state.filter.yearRange.min }} - {{ state.filter.yearRange.max }}</div>
          <q-range
            v-model="state.filter.yearRange"
            :min="2008"
            :max="2023"
            :step="1"
            label-always
            color="primary"
          />
        </div>

        <!-- Input Kata Kunci -->
        <div class="col-12 col-sm">
          <q-input
            v-model="state.filter.keyword"
            label="Masukkan Nama / NIS / NISN..."
            outlined
            dense
            autofocus
            @keyup.enter="executeSearch"
          >
            <template v-slot:append>
              <q-btn round flat icon="search" @click="executeSearch" />
            </template>
          </q-input>
        </div>
      </q-card-section>

      <!-- Tabel Hasil Pencarian -->
      <q-card-section>
        <q-table
          flat
          bordered
          :rows="state.results"
          :columns="columns"
          row-key="nis"
          :loading="state.loading"
          no-data-label="Tidak ada data ditemukan. Silakan masukkan kata kunci."
        >
          <template v-slot:body-cell-actions="props">
            <q-td :props="props" class="q-gutter-sm">
              <q-btn size="sm" color="secondary" label="Detail" icon="visibility" @click="viewDetail(props.row)" />
              <q-btn size="sm" color="warning" text-color="dark" label="Edit" icon="edit" @click="handleEditAction(props.row)" />
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>
  </q-dialog>

  <!-- SUB-MODAL DETAIL DATA (Terenkapsulasi di dalam komponen search) -->
  <q-dialog v-model="state.isDetailOpen">
    <q-card style="width: 450px;">
      <q-card-section class="bg-primary text-white"><div class="text-h6">Profil Administrasi Siswa</div></q-card-section>
      <q-card-section class="q-pa-md">
        <q-list dense>
          <q-item><q-item-section class="text-weight-bold">Nama:</q-item-section><q-item-section side>{{ state.selectedStudent.nama }}</q-item-section></q-item>
          <q-item><q-item-section class="text-weight-bold">NIS / NISN:</q-item-section><q-item-section side>{{ state.selectedStudent.nis }} / {{ state.selectedStudent.nisn }}</q-item-section></q-item>
          <q-item><q-item-section class="text-weight-bold">Status:</q-item-section><q-item-section side>{{ state.selectedStudent.status }}</q-item-section></q-item>
          <q-item><q-item-section class="text-weight-bold">{{ state.selectedStudent.status === 'Alumni' ? 'Kelas Terakhir:' : 'Kelas Aktif:' }}</q-item-section><q-item-section side>{{ state.selectedStudent.kelas }}</q-item-section></q-item>
          <q-item><q-item-section class="text-weight-bold">Wali Kelas:</q-item-section><q-item-section side>{{ state.selectedStudent.wali_kelas }}</q-item-section></q-item>
          <q-item v-if="state.selectedStudent.status === 'Alumni'"><q-item-section class="text-weight-bold">Tahun Lulus:</q-item-section><q-item-section side>{{ state.selectedStudent.tahun_lulus }}</q-item-section></q-item>
        </q-list>
      </q-card-section>
      <q-card-actions align="right"><q-btn flat label="Tutup" color="primary" v-close-popup /></q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'
import { api } from '@/boot/axios'
import { useQuasar } from 'quasar'

// Deklarasi Props & Emits sesuai standar Vue 3 V-Model Binding
defineProps({
  modelValue: { type: Boolean, default: false }
})
const emit = defineEmits(['update:modelValue'])

const router = useRouter()
const authStore = useAuthStore()
const $q = useQuasar()

const state = reactive({
  loading: false,
  isDetailOpen: false,
  selectedStudent: {},
  filter: {
    status: 'Siswa Aktif',
    keyword: '',
    yearRange: { min: 2015, max: 2020 }
  },
  results: []
})

const columns = [
  { name: 'nama', label: 'Nama Lengkap', field: 'nama', align: 'left', sortable: true },
  { name: 'nis', label: 'NIS', field: 'nis', align: 'left' },
  { name: 'nisn', label: 'NISN', field: 'nisn', align: 'left' },
  { name: 'kelas', label: 'Kelas', field: 'kelas', align: 'center' },
  { name: 'wali_kelas', label: 'Wali Kelas', field: 'wali_kelas', align: 'left' },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' }
]

async function executeSearch() {
  if (!state.filter.keyword.trim()) {
    $q.notify({ type: 'warning', message: 'Masukkan kata kunci pencarian!' })
    return
  }
  state.loading = true
  try {
    const params = {
      status: state.filter.status,
      q: state.filter.keyword,
      ...(state.filter.status === 'Alumni' && {
        year_start: state.filter.yearRange.min,
        year_end: state.filter.yearRange.max
      })
    }
    const response = await api.get('/students/search', { params })
    state.results = response.data.data
  } catch (error) {
    console.log(error)
    $q.notify({ type: 'negative', message: 'Gagal memanggil server internal.' })
  } finally {
    state.loading = false
  }
}

const viewDetail = (row) => {
  state.selectedStudent = row
  state.isDetailOpen = true
}

const handleEditAction = (student) => {
  // Tutup dialog pencarian terlebih dahulu
  emit('update:modelValue', false)

  if (!authStore.isAuthenticated) {
    $q.notify({ type: 'negative', message: 'Akses administrasi ditolak! Silakan login.', position: 'top' })
    router.push({ name: 'auth-login', query: { redirect: `/siswa/ppdb/edit/${student.nis}` } })
    return
  }

  if (['admin', 'kesiswaan'].includes(authStore.userRole)) {
    router.push({ path: `/siswa/ppdb/edit/${student.nis}` })
  } else {
    $q.notify({ type: 'warning', message: 'Akun TU Anda tidak berwenang mengedit modul ini.', position: 'top' })
  }
}
</script>
