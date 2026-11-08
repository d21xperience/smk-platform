<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row justify-between items-center q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Penugasan Mengajar</h1>
        <p class="text-caption text-grey-7 q-mb-none">Petakan Mata Pelajaran & Guru ke Kelas</p>
      </div>
      <q-btn color="primary" icon="add" label="Tambah Penugasan" @click="openDialog" />
    </div>

    <q-card class="shadow-1">
      <q-table :rows="pembelajaranStore.list" :columns="columns" row-key="id" :loading="pembelajaranStore.isLoading"
        flat bordered dense>
        <template v-slot:top-left>
          <q-select v-model="filterClassId" :options="kelasOptions" label="Filter Kelas" option-value="id"
            option-label="name" emit-value map-options outlined dense clearable style="min-width: 200px"
            @update:model-value="loadData" />
        </template>

        <template v-slot:body-cell-teacherName="props">
          <q-td :props="props">
            <span v-if="props.value" class="text-weight-medium text-primary">{{ props.value }}</span>
            <span v-else class="text-grey-5 italic">Belum ada guru</span>
          </q-td>
        </template>

        <template v-slot:body-cell-aksi="props">
          <q-td :props="props" class="text-center">
            <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row)" />
          </q-td>
        </template>
      </q-table>
    </q-card>

    <q-dialog v-model="dialogOpen" persistent>
      <q-card style="min-width: 500px">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">Tambah Penugasan Mengajar</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="onSubmit" class="q-gutter-md">
            <q-select v-model="form.classId" :options="kelasOptions" label="Pilih Kelas" option-value="id"
              option-label="name" emit-value map-options filled required />
            <q-select v-model="form.subjectId" :options="mapelOptions" label="Pilih Mata Pelajaran" option-value="id"
              option-label="name" emit-value map-options filled required />
            <q-select v-model="form.teacherId" :options="guruOptions" label="Pilih Guru Pengajar" option-value="id"
              option-label="name" emit-value map-options filled clearable />
            <div class="row justify-end q-mt-lg">
              <q-btn label="Batal" color="grey" flat v-close-popup />
              <q-btn label="Simpan" color="primary" type="submit" :loading="isSaving" />
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
import { usePembelajaranStore } from '@/stores/kesiswaan/pembelajaranStore'
import { useKelasStore } from '@/stores/kesiswaan/kelasStore'
import { useMapelStore } from '@/stores/kesiswaan/mapelStore'
import { useGuruStore } from '@/stores/kepegawaian/guru'

const $q = useQuasar()
const pembelajaranStore = usePembelajaranStore()
const kelasStore = useKelasStore()
const mapelStore = useMapelStore()
const guruStore = useGuruStore()

const dialogOpen = ref(false)
const isSaving = ref(false)
const filterClassId = ref(null)
const form = ref({ classId: '', subjectId: '', teacherId: '' })

const kelasOptions = computed(() => kelasStore.list)
const mapelOptions = computed(() => mapelStore.list)
const guruOptions = computed(() => guruStore.list)

const columns = [
  { name: 'className', label: 'Kelas', field: 'className', align: 'left', sortable: true },
  { name: 'subjectName', label: 'Mata Pelajaran', field: 'subjectName', align: 'left', sortable: true },
  { name: 'teacherName', label: 'Guru Pengajar', field: 'teacherName', align: 'left' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

onMounted(async () => {
  await Promise.all([
    kelasStore.fetchKelas({ limit: 100 }),
    mapelStore.fetchSubjects(),
    guruStore.fetchTeachers()
  ])
  loadData()
})

const loadData = () => {
  pembelajaranStore.fetchAll(filterClassId.value || '')
}

const openDialog = () => {
  form.value = { classId: filterClassId.value || '', subjectId: '', teacherId: '' }
  dialogOpen.value = true
}

const onSubmit = async () => {
  isSaving.value = true
  try {
    await pembelajaranStore.assignTeacher(form.value)
    $q.notify({ color: 'positive', message: 'Penugasan berhasil disimpan' })
    dialogOpen.value = false
    filterClassId.value = form.value.classId // Set filter ke kelas yang baru ditambahkan
    loadData()
  } catch (err) {
    $q.notify({ color: 'negative', message: err.response?.data?.error || 'Gagal menyimpan penugasan' })
  } finally {
    isSaving.value = false
  }
}

const confirmDelete = (row) => {
  $q.dialog({
    title: 'Konfirmasi Hapus',
    message: `Hapus penugasan ${row.subjectName} di kelas ${row.className}?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await pembelajaranStore.removeAssignment(row.id, row.classId)
      $q.notify({ color: 'positive', message: 'Penugasan dihapus' })
      loadData()
    } catch (err) {
      $q.notify({ color: 'negative', message: `Gagal menghapus penugasan ${err}` })
    }
  })
}
</script>
