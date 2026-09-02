<template>
  <q-page padding>
    <div class="text-h5 q-mb-md">Daftar Surat Masuk</div>
    <q-table :rows="list" :columns="columns" row-key="id" :filter="filter" :loading="isLoading"
      :pagination="{ rowsPerPage: 10 }">
      <template v-slot:top-right>
        <q-input dense outlined v-model="filter" placeholder="Cari..." class="q-mr-sm" />
        <q-btn color="primary" icon="add" label="Input Baru" @click="openDialog" />
      </template>
      <template v-slot:body-cell-status="props">
        <q-td :props="props">
          <q-badge :color="statusColor(props.value)">{{ props.value }}</q-badge>
        </q-td>
      </template>
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props">
          <q-btn flat icon="visibility" @click="viewDetail(props.row.id)" />
        </q-td>
      </template>
    </q-table>

    <!-- Dialog Input Surat Masuk -->
    <q-dialog v-model="isDialogOpen" persistent>
      <q-card style="min-width: 600px">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">Input Surat Masuk Baru</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-form @submit.prevent="submitCreate">
            <q-input v-model="form.asal" label="Asal Surat" required />
            <q-input v-model="form.nomorSurat" label="Nomor Surat" required />
            <q-input v-model="form.tanggalSurat" label="Tanggal Surat" type="date" required />
            <q-input v-model="form.perihal" label="Perihal" type="textarea" required />
            <q-select v-model="form.kodeKlasifikasi" :options="klasifikasiOptions" label="Kode Klasifikasi" use-input
              fill-input hide-selected input-debounce="0" @filter="filterKlasifikasi" required />
            <q-uploader v-model="form.fileSurat" label="Unggah Scan PDF" accept=".pdf" max-files="1" class="q-mt-md"
              auto-upload :factory="uploadFile" />
            <div class="q-mt-md text-right">
              <q-btn label="Batal" flat color="grey" v-close-popup />
              <q-btn label="Simpan" type="submit" color="primary" :loading="isLoading" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- Dialog Detail & Timeline -->
    <q-dialog v-model="isDetailDialogOpen" maximized>
      <q-card>
        <q-bar class="bg-primary text-white">
          <div>Detail Surat Masuk</div>
          <q-space />
          <q-btn dense flat icon="close" v-close-popup />
        </q-bar>
        <q-card-section v-if="currentItem">
          <div class="text-h6">{{ currentItem.perihal }}</div>
          <p><strong>Asal:</strong> {{ currentItem.asal }}</p>
          <p><strong>Nomor Surat:</strong> {{ currentItem.nomorSuratAsal }}</p>
          <p><strong>Tanggal Diterima:</strong> {{ currentItem.tanggalDiterima }}</p>
          <p><strong>Klasifikasi:</strong> {{ currentItem.classificationCode }}</p>

          <!-- PDF Viewer Embedded -->
          <q-card class="q-mt-md">
            <q-card-section>
              <div class="text-h6">Berkas PDF</div>
              <iframe v-if="currentItem.fileUrl" :src="currentItem.fileUrl" width="100%" height="500px"
                frameborder="0"></iframe>
            </q-card-section>
          </q-card>

          <!-- Timeline Perjalanan Surat -->
          <q-card class="q-mt-md">
            <q-card-section>
              <div class="text-h6">Riwayat Perjalanan Surat</div>
              <q-timeline color="secondary">
                <q-timeline-entry title="Diterima TU" subtitle="20 Mei 2026" icon="inbox">
                  Surat diterima dan dicatat oleh staf TU.
                </q-timeline-entry>
                <q-timeline-entry title="Disposisi Kepala Sekolah" subtitle="21 Mei 2026" icon="send">
                  Ditujukan kepada Waka Kurikulum dengan instruksi...
                </q-timeline-entry>
                <q-timeline-entry title="Selesai" subtitle="23 Mei 2026" icon="done_all" color="green">
                  Tindak lanjut selesai.
                </q-timeline-entry>
              </q-timeline>
            </q-card-section>
          </q-card>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useIncomingLetter } from '@/composables/useIncomingLetter.js'

const {
  list,
  currentItem,
  isLoading,
  filter,
  form,
  klasifikasiOptions,
  isDialogOpen,
  isDetailDialogOpen,
  loadList,
  loadKlasifikasiOptions,
  submitCreate,
  viewDetail,
  openDialog,
  statusColor
} = useIncomingLetter()

// Columns definition (UI specific)
const columns = [
  { name: 'noAgenda', label: 'No Agenda', field: 'noAgenda', align: 'left', sortable: true },
  { name: 'nomorSuratAsal', label: 'Nomor Surat', field: 'nomorSuratAsal', align: 'left', sortable: true },
  { name: 'asal', label: 'Asal', field: 'asal', align: 'left' },
  { name: 'perihal', label: 'Perihal', field: 'perihal', align: 'left' },
  { name: 'tanggalDiterima', label: 'Tgl Diterima', field: 'tanggalDiterima', align: 'left' },
  { name: 'classificationCode', label: 'Klasifikasi', field: 'classificationCode', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'left' },
  { name: 'aksi', label: 'Aksi', field: 'aksi', align: 'center' }
]

// Local UI logic for classification filter (autocomplete)
const semuaKlasifikasi = ref([])
async function initKlasifikasi() {
  await loadKlasifikasiOptions()
  semuaKlasifikasi.value = [...klasifikasiOptions.value]
}

function filterKlasifikasi(val, update) {
  update(() => {
    const needle = val.toLowerCase()
    klasifikasiOptions.value = semuaKlasifikasi.value.filter(v => v.toLowerCase().indexOf(needle) > -1)
  })
}

// Mock upload file (Nanti akan diganti dengan real upload service)
function uploadFile(files) {
  return new Promise((resolve) => {
    resolve({ files: [files[0]], url: 'https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf' })
  })
}

// Watch filter to reload list (debounced in real app, simple watch here)
watch(filter, () => {
  loadList()
})

onMounted(() => {
  loadList()
  initKlasifikasi()
})
</script>
