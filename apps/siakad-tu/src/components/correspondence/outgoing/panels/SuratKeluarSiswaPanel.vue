<template>
  <q-table
    :rows="listSiswa"
    :columns="kolomSiswa"
    row-key="id"
    :filter="filterSiswa"
    :loading="isLoading"
    :pagination="{ rowsPerPage: 10 }"
  >
    <template v-slot:top-right>
      <q-input
        dense
        outlined
        v-model="filterSiswa"
        placeholder="Cari..."
        class="q-mr-sm"
      />
      <q-btn color="primary" icon="add" label="Buat Surat Siswa" @click="openDialogSiswa" />
    </template>

    <template v-slot:body-cell-aksi="props">
      <q-td :props="props">
        <q-btn flat icon="print" color="primary" size="sm" @click="openPrintDialog(props.row)" />
      </q-td>
    </template>
  </q-table>

  <!-- Form Dialog dipanggil di sini agar ter-render dalam scope yang sama -->
  <SuratKeluarFormSiswa />
</template>

<script setup>
import { watch } from 'vue'
import { useOutgoingLetter } from '@/composables/useOutgoingLetter.js'
import SuratKeluarFormSiswa from '../forms/SuratKeluarFormSiswa.vue'

const {
  listSiswa,
  isLoading,
  filterSiswa,
  loadList,
  openDialogSiswa,
  openPrintDialog
} = useOutgoingLetter()

const kolomSiswa = [
  { name: 'nomorSurat', label: 'Nomor Surat', field: 'nomorSurat', sortable: true },
  { name: 'jenis', label: 'Jenis', field: 'jenis' },
  { name: 'tujuan', label: 'Tujuan', field: 'tujuan' },
  { name: 'tanggal', label: 'Tanggal', field: 'tanggal' },
  { name: 'status', label: 'Status', field: 'status' },
  { name: 'aksi', label: 'Cetak', align: 'center' }
]

watch(filterSiswa, () => {
  loadList('siswa')
})
</script>
