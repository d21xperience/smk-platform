<template>
  <q-table :rows="listPtk" :columns="kolomPTK" row-key="id" :filter="filterPtk" :loading="isLoading"
    :pagination="{ rowsPerPage: 10 }">
    <template v-slot:top-right>
      <q-input dense outlined v-model="filterPtk" placeholder="Cari..." class="q-mr-sm" />
      <q-btn color="primary" icon="add" label="Buat Surat PTK" @click="openDialogPtk" />
    </template>
  </q-table>

  <SuratKeluarFormPtk />
</template>

<script setup>
import { watch } from 'vue'
import { useOutgoingLetter } from '@/composables/useOutgoingLetter.js'
import SuratKeluarFormPtk from '../forms/SuratKeluarFormPtk.vue'

const {
  listPtk,
  isLoading,
  filterPtk,
  loadList,
  openDialogPtk
} = useOutgoingLetter()

const kolomPTK = [
  { name: 'nomorSurat', label: 'Nomor Surat', field: 'nomorSurat' },
  { name: 'jenis', label: 'Jenis', field: 'jenis' },
  { name: 'tujuan', label: 'Tujuan/Deskripsi', field: 'tujuan' },
  { name: 'tanggal', label: 'Tanggal', field: 'tanggal' },
  { name: 'status', label: 'Status', field: 'status' }
]

watch(filterPtk, () => {
  loadList('ptk')
})
</script>
