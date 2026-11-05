<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-md">
      <div class="text-h5 text-weight-bold">Pengaturan Tahun Ajaran</div>
      <q-space />
      <q-btn color="primary" label="Tambah Tahun Ajaran" icon="add" @click="showAddDialog = true" />
    </div>

    <q-table :rows="periods" :columns="columns" row-key="semester_id" flat bordered :loading="isLoading">
      <template v-slot:body-cell-status="props">
        <q-td :props="props">
          <q-badge :color="props.row.periode_aktif ? 'green' : 'grey'" class="q-px-sm q-py-xs">
            {{ props.row.periode_aktif ? 'Aktif' : 'Tidak Aktif' }}
          </q-badge>
        </q-td>
      </template>

      <template v-slot:body-cell-actions="props">
        <q-td :props="props" class="q-gutter-x-sm">
          <q-btn v-if="!props.row.periode_aktif" color="positive" label="Aktifkan" size="sm" unelevated
            @click="activatePeriod(props.row.semester_id)" />
          <span v-else class="text-grey-6 text-italic text-caption q-mr-sm">Sedang Digunakan</span>

          <!-- Tombol Hapus (Hanya muncul jika tidak aktif) -->
          <q-btn v-if="!props.row.periode_aktif" color="negative" icon="delete" size="sm" unelevated
            @click="confirmDelete(props.row)">
            <q-tooltip>Hapus Periode</q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </q-table>

    <!-- Dialog Tambah -->
    <q-dialog v-model="showAddDialog" persistent>
      <q-card style="min-width: 350px">
        <q-card-section class="row items-center">
          <div class="text-h6">Tambah Tahun Ajaran Baru</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-input v-model="form.year" label="Tahun Awal (Cth: 2027)" mask="####" fill-mask="0"
            hint="Sistem akan otomatis membuat format 2027/2028 beserta Semester 1 & 2" outlined dense />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Simpan" :loading="isSubmitting" @click="submitPeriod"
            :disable="form.year.length !== 4" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { api } from '@/boot/axios'
import { Notify } from 'quasar'

const $q = useQuasar()

const periods = ref([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const showAddDialog = ref(false)
const form = ref({ year: '' })

const columns = [
  { name: 'semester_id', label: 'ID', field: 'semester_id', align: 'left' },
  { name: 'tahun_ajaran_label', label: 'Tahun Ajaran', field: 'tahun_ajaran_label', align: 'left' },
  { name: 'nama_semester', label: 'Semester', field: 'nama_semester', align: 'left' },
  { name: 'status', label: 'Status', field: 'periode_aktif', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center' }
]

async function fetchPeriods() {
  isLoading.value = true
  try {
    const res = await api.get('/context/periods')
    periods.value = res.data
  } catch (err) {
    Notify.create({ type: 'negative', message: `Gagal memuat data periode ${err}` })
  } finally {
    isLoading.value = false
  }
}

async function submitPeriod() {
  if (!form.value.year || form.value.year.length !== 4) return
  isSubmitting.value = true
  try {
    await api.post('/academic/periods', { year: form.value.year })
    Notify.create({ type: 'positive', message: 'Tahun ajaran dan 2 semester berhasil ditambahkan' })
    showAddDialog.value = false
    form.value = { year: '' }
    await fetchPeriods()
  } catch (err) {
    Notify.create({ type: 'negative', message: err.response?.data?.error || 'Gagal menambah periode' })
  } finally {
    isSubmitting.value = false
  }
}

async function activatePeriod(id) {
  try {
    await api.put(`/academic/periods/${id}/activate`)
    Notify.create({ type: 'positive', message: 'Periode berhasil diaktifkan' })
    await fetchPeriods()
    window.location.reload()
  } catch (err) {
    Notify.create({ type: 'negative', message: err.response?.data?.error || 'Gagal mengaktifkan' })
  }
}

// Fungsi Konfirmasi Hapus
function confirmDelete(period) {
  $q.dialog({
    title: 'Konfirmasi Hapus',
    message: `Apakah Anda yakin ingin menghapus ${period.tahun_ajaran_label} - ${period.nama_semester}?`,
    cancel: true,
    persistent: true,
    color: 'negative'
  }).onOk(async () => {
    try {
      await api.delete(`/academic/periods/${period.semester_id}`)
      Notify.create({ type: 'positive', message: 'Periode berhasil dihapus' })
      await fetchPeriods()
    } catch (err) {
      // Menangkap pesan error dari backend jika semester sedang aktif
      Notify.create({
        type: 'negative',
        message: err.response?.data?.error || 'Gagal menghapus periode',
        timeout: 4000
      })
    }
  })
}

onMounted(fetchPeriods)
</script>
