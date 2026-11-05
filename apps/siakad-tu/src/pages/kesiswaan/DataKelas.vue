<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Manajemen Data Kelas</h1>
        <p class="text-caption text-grey-7 q-mb-none">Kelola data rombongan belajar sesuai semester aktif</p>
      </div>
      <div class="q-gutter-x-sm">
        <q-btn color="primary" icon="add" label="Tambah Kelas" @click="openDialog()" />
      </div>
    </div>

    <!-- Tabel Kelas -->
    <q-card class="shadow-1">
      <q-table :rows="kelasStore.list" :columns="columns" row-key="id" :loading="kelasStore.isLoading" flat bordered>
        <template v-slot:top-right>
          <q-input v-model="filterKeyword" placeholder="Cari kelas..." dense outlined clearable
            @update:model-value="searchData">
            <template v-slot:append><q-icon name="search" /></template>
          </q-input>
        </template>

        <template v-slot:body-cell-level="props">
          <q-td :props="props">
            <q-badge color="blue-6">Kelas {{ props.row.level || props.row.grade }}</q-badge>
          </q-td>
        </template>

        <template v-slot:body-cell-totalStudents="props">
          <q-td :props="props">
            <q-badge :color="props.row.currentStudents >= props.row.capacity ? 'red-6' : 'green-6'">
              {{ props.row.currentStudents || 0 }} / {{ props.row.capacity }}
            </q-badge>
          </q-td>
        </template>

        <!-- KOLOM BARU: Wali Kelas -->
        <template v-slot:body-cell-homeroom="props">
          <q-td :props="props">
            <span v-if="props.row.homeroomId" class="text-weight-medium text-primary">
              {{ getTeacherName(props.row.homeroomId) }}
            </span>
            <span v-else class="text-grey-5 italic">Belum diatur</span>
          </q-td>
        </template>

        <template v-slot:body-cell-actions="props">
          <q-td :props="props" class="q-gutter-x-sm">
            <q-btn flat round dense color="teal" icon="people" @click="openMemberDialog(props.row)">
              <q-tooltip>Kelola Anggota Kelas</q-tooltip>
            </q-btn>
            <q-btn flat round dense color="orange" icon="edit" @click="openDialog(props.row)">
              <q-tooltip>Edit Kelas</q-tooltip>
            </q-btn>
            <q-btn flat round dense color="negative" icon="delete" @click="confirmDelete(props.row)">
              <q-tooltip>Hapus Kelas</q-tooltip>
            </q-btn>
          </q-td>
        </template>

        <template v-slot:no-data>
          <div class="full-width row flex-center text-grey-6 q-gutter-y-sm q-pa-lg">
            <q-icon size="48px" name="domain_disabled" />
            <span class="text-subtitle1">Belum ada data kelas di semester ini.</span>
          </div>
        </template>
      </q-table>
    </q-card>

    <!-- Dialog Form Tambah/Edit -->
    <q-dialog v-model="dialogOpen" persistent>
      <q-card style="min-width: 450px">
        <q-card-section class="row items-center q-pb-none bg-primary text-white">
          <div class="text-h6">{{ isEdit ? 'Edit Data Kelas' : 'Tambah Kelas Baru' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section>
          <q-form @submit="onSubmit" class="q-gutter-md">
            <q-input v-model="form.name" label="Nama Kelas" hint="Contoh: X TKJ 1" filled
              :rules="[val => !!val || 'Wajib diisi']" />
            <div class="row q-col-gutter-sm">
              <div class="col-6">
                <q-input v-model.number="form.grade" label="Tingkat" type="number" filled
                  :rules="[val => !!val || 'Wajib diisi']" />
              </div>
              <div class="col-6">
                <q-input v-model.number="form.capacity" label="Kapasitas" type="number" filled
                  :rules="[val => val > 0 || 'Minimal 1']" />
              </div>
            </div>
            <q-input v-model="form.major" label="Jurusan (Opsional)" filled />

            <!-- DROPDOWN WALI KELAS -->
            <q-select v-model="form.homeroomId" :options="teacherOptions" label="Wali Kelas" option-value="id"
              option-label="name" emit-value map-options filled clearable :loading="guruStore.isLoading" />

            <div class="row justify-end q-mt-lg">
              <q-btn label="Batal" color="grey" flat v-close-popup class="q-mr-sm" />
              <q-btn label="Simpan Data" color="primary" type="submit" :loading="kelasStore.isLoading" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog Kelola Anggota -->
    <KelolaAnggotaDialog v-model="memberDialogOpen" :class-data="selectedClassForMembers" />
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useQuasar } from 'quasar'
import { useKelasStore } from '@/stores/kesiswaan/kelasStore'
import { useGuruStore } from '@/stores/kepegawaian/guru' // Import store guru
import KelolaAnggotaDialog from '@/components/kesiswaan/KelolaAnggotaDialog.vue'

const $q = useQuasar()
const kelasStore = useKelasStore()
const guruStore = useGuruStore()

const filterKeyword = ref('')
const dialogOpen = ref(false)
const isEdit = ref(false)
const form = ref({ id: '', name: '', grade: 10, major: '', capacity: 36, homeroomId: null })

const memberDialogOpen = ref(false)
const selectedClassForMembers = ref(null)

// Opsi untuk Dropdown Wali Kelas
const teacherOptions = computed(() => guruStore.list)

const columns = [
  { name: 'name', label: 'Nama Kelas', field: 'name', align: 'left', sortable: true },
  { name: 'level', label: 'Tingkat', field: 'grade', align: 'center' },
  { name: 'major', label: 'Jurusan', field: 'major', align: 'left' },
  { name: 'homeroom', label: 'Wali Kelas', field: 'homeroomId', align: 'left' }, // Kolom Baru
  { name: 'totalStudents', label: 'Total Siswa', field: 'currentStudents', align: 'center', sortable: true },
  { name: 'actions', label: 'Aksi', align: 'center' }
]

onMounted(async () => {
  loadData()
  try {
    await guruStore.fetchTeachers() // Ambil data guru saat halaman dimuat
  } catch (e) {
    console.error("Gagal load guru", e)
  }
})

// Helper untuk menampilkan nama wali kelas di tabel
const getTeacherName = (id) => {
  const teacher = guruStore.list.find(t => t.id === id)
  return teacher ? teacher.name : id
}

const loadData = async (keyword = '') => {
  try {
    await kelasStore.fetchKelas({ keyword, limit: 100 })
  } catch (err) {
    $q.notify({ color: 'negative', message: `Gagal memuat data kelas ${err}` })
  }
}

const searchData = () => loadData(filterKeyword.value)

const openDialog = (row = null) => {
  isEdit.value = !!row
  if (row) {
    form.value = { ...row, homeroomId: row.homeroomId || null }
  } else {
    form.value = { id: '', name: '', grade: 10, major: '', capacity: 36, homeroomId: null }
  }
  dialogOpen.value = true
}

const onSubmit = async () => {
  try {
    const payload = {
      ...form.value,
      level: form.value.level || '',
      homeroomId: form.value.homeroomId || "" // Kirim string kosong jika null
    }
    if (isEdit.value) {
      await kelasStore.updateKelas(form.value.id, payload)
      $q.notify({ color: 'positive', message: 'Kelas diperbarui' })
    } else {
      await kelasStore.createKelas(payload)
      $q.notify({ color: 'positive', message: 'Kelas ditambahkan' })
    }
    dialogOpen.value = false
  } catch (err) {
    $q.notify({ color: 'negative', message: err.response?.data?.error || 'Terjadi kesalahan' })
  }
}

const confirmDelete = (row) => {
  $q.dialog({
    title: 'Konfirmasi Hapus',
    message: `Hapus kelas ${row.name}?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await kelasStore.deleteKelas(row.id)
      $q.notify({ color: 'positive', message: 'Kelas dihapus' })
    } catch (err) {
      $q.notify({ color: 'negative', message: `Gagal menghapus kelas ${err}` })
    }
  })
}

const openMemberDialog = (row) => {
  selectedClassForMembers.value = row
  memberDialogOpen.value = true
}
</script>
