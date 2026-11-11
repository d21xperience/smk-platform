<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-md">
      <div class="text-h5 text-weight-bold">Master Jenis Surat</div>
      <q-space />
      <q-btn color="primary" icon="add" label="Tambah Jenis Surat" @click="openForm()" />
    </div>

    <!-- Tabel Data -->
    <q-table
      :rows="store.list"
      :columns="columns"
      row-key="id"
      flat
      bordered
      :loading="store.isLoading"
      :filter="filter"
    >
      <template v-slot:top-right>
        <q-input v-model="filter" placeholder="Cari jenis surat..." outlined dense clearable>
          <template v-slot:prepend>
            <q-icon name="search" />
          </template>
        </q-input>
      </template>

      <template v-slot:body-cell-peruntukan="props">
        <q-td :props="props">
          <q-badge :color="getPeruntukanColor(props.row.peruntukan)">
            {{ props.row.peruntukan }}
          </q-badge>
        </q-td>
      </template>

      <template v-slot:body-cell-actions="props">
        <q-td :props="props" class="q-gutter-x-xs">
          <q-btn flat round dense color="blue" icon="edit" @click="openForm(props.row)" />
          <q-btn flat round dense color="negative" icon="delete" @click="confirmDelete(props.row)" />
        </q-td>
      </template>
    </q-table>

    <!-- Dialog Form CRUD -->
    <q-dialog v-model="showForm" persistent>
      <q-card style="min-width: 500px; max-width: 90vw">
        <q-card-section class="row items-center">
          <div class="text-h6">{{ isEdit ? 'Edit Jenis Surat' : 'Tambah Jenis Surat' }}</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-none q-gutter-y-md">
          <q-input v-model="form.kategori" label="Nama Kelompok/Kategori" outlined dense hint="Contoh: 1. Kelompok Surat Berdasarkan Aturan Resmi Dinas" />
          
          <q-select
            v-model="form.peruntukan"
            :options="['SISWA', 'GURU', 'UMUM']"
            label="Peruntukan"
            outlined
            dense
          />
          
          <q-input v-model="form.label" label="Label Surat (Tampil di Dropdown)" outlined dense hint="Contoh: Surat Keputusan (SK)" />
          
          <q-input v-model="form.value" label="Kode Unik (Value)" outlined dense hint="Contoh: SK atau PERMOHONAN_PKL" />
          
          <q-input v-model="form.kode_klasifikasi" label="Kode Klasifikasi" outlined dense hint="Contoh: 421.5" />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Simpan" :loading="store.isSubmitting" @click="submitForm" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Notify, Dialog } from 'quasar'
import { useLetterTypeStore } from '@/stores/letter-type'

const store = useLetterTypeStore()

const filter = ref('')
const showForm = ref(false)
const isEdit = ref(false)
const form = ref({ id: '', kategori: '', peruntukan: 'UMUM', label: '', value: '', kode_klasifikasi: '' })

const columns = [
  { name: 'kategori', label: 'Kategori', field: 'kategori', align: 'left', sortable: true },
  { name: 'peruntukan', label: 'Peruntukan', field: 'peruntukan', align: 'center', sortable: true },
  { name: 'label', label: 'Nama Surat', field: 'label', align: 'left', sortable: true },
  { name: 'value', label: 'Kode Sistem', field: 'value', align: 'left' },
  { name: 'kode_klasifikasi', label: 'Kode Klasifikasi', field: 'kode_klasifikasi', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center' }
]

function getPeruntukanColor(peruntukan) {
  if (peruntukan === 'SISWA') return 'orange-6'
  if (peruntukan === 'GURU') return 'blue-6'
  return 'green-6'
}

function openForm(row = null) {
  if (row) {
    isEdit.value = true
    form.value = { ...row }
  } else {
    isEdit.value = false
    form.value = { id: '', kategori: '', peruntukan: 'UMUM', label: '', value: '', kode_klasifikasi: '' }
  }
  showForm.value = true
}

async function submitForm() {
  if (!form.value.kategori || !form.value.label || !form.value.value || !form.value.kode_klasifikasi) {
    return Notify.create({ type: 'warning', message: 'Semua kolom wajib diisi' })
  }
  try {
    if (isEdit.value) {
      await store.update(form.value.id, form.value)
      Notify.create({ type: 'positive', message: 'Jenis surat berhasil diperbarui' })
    } else {
      await store.create(form.value)
      Notify.create({ type: 'positive', message: 'Jenis surat berhasil ditambahkan' })
    }
    showForm.value = false
  } catch (err) {
    Notify.create({ type: 'negative', message: err.response?.data?.error || 'Gagal menyimpan data' })
  }
}

function confirmDelete(row) {
  Dialog.create({
    title: 'Konfirmasi Hapus',
    message: `Apakah Anda yakin ingin menghapus jenis surat "${row.label}"?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    try {
      await store.delete(row.id)
      Notify.create({ type: 'positive', message: 'Jenis surat berhasil dihapus' })
    } catch {
      Notify.create({ type: 'negative', message: 'Gagal menghapus data' })
    }
  })
}

onMounted(() => {
  store.fetchList()
})
</script>
