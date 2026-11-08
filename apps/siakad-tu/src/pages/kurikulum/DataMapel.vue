<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row justify-between items-center q-mb-md">
      <h1 class="text-h5 text-weight-bold q-my-none text-primary">Master Mata Pelajaran</h1>
      <q-btn color="primary" icon="add" label="Tambah Mapel" @click="openDialog" />
    </div>

    <q-card class="shadow-1">
      <q-table :rows="mapelStore.list" :columns="columns" row-key="id" :loading="mapelStore.isLoading" flat bordered
        dense>
        <template v-slot:body-cell-aksi="props">
          <q-td :props="props" class="text-center">
            <q-btn flat round dense icon="delete" color="negative" @click="confirmDelete(props.row.id)" />
          </q-td>
        </template>
      </q-table>
    </q-card>

    <q-dialog v-model="dialogOpen" persistent>
      <q-card style="min-width: 400px">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">Tambah Mata Pelajaran</div>
        </q-card-section>
        <q-card-section>
          <q-form @submit="onSubmit" class="q-gutter-md">
            <q-input v-model="form.name" label="Nama Mapel" filled required />
            <q-input v-model="form.kelompok" label="Kelompok" hint="Wajib/Umum/Kejuruan" filled />
            <q-input v-model.number="form.semester" label="Semester" type="number" filled />
            <div class="row justify-end q-mt-md">
              <q-btn label="Batal" color="grey" flat v-close-popup />
              <q-btn label="Simpan" color="primary" type="submit" :loading="mapelStore.isLoading" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useMapelStore } from '@/stores/kesiswaan/mapelStore'

const $q = useQuasar()
const mapelStore = useMapelStore()
const dialogOpen = ref(false)
const form = ref({ name: '', kelompok: '', semester: 1 })

const columns = [
  { name: 'name', label: 'Nama Mapel', field: 'name', align: 'left', sortable: true },
  { name: 'kelompok', label: 'Kelompok', field: 'kelompok', align: 'left' },
  { name: 'semester', label: 'Semester', field: 'semester', align: 'center' },
  { name: 'aksi', label: 'Aksi', align: 'center' }
]

onMounted(() => mapelStore.fetchSubjects())

const openDialog = () => {
  form.value = { name: '', kelompok: '', semester: 1 }
  dialogOpen.value = true
}

const onSubmit = async () => {
  try {
    await mapelStore.createSubject(form.value)
    $q.notify({ color: 'positive', message: 'Mapel ditambahkan' })
    dialogOpen.value = false
  } catch (e) {
    $q.notify({ color: 'negative', message: `Gagal menambah mapel ${e}` })
  }
}

const confirmDelete = (id) => {
  $q.dialog({ title: 'Hapus', message: 'Yakin hapus mapel ini?', cancel: true })
    .onOk(async () => {
      await mapelStore.deleteSubject(id)
      $q.notify({ color: 'positive', message: 'Mapel dihapus' })
    })
}
</script>
