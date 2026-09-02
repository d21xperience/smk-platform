<template>
  <q-table :rows="listKedinasan" :columns="kolomKedinasan" row-key="id" :filter="filterKedinasan" :loading="isLoading"
    :pagination="{ rowsPerPage: 10 }">
    <template v-slot:top-right>
      <q-input dense outlined v-model="filterKedinasan" placeholder="Cari..." class="q-mr-sm" />
      <q-btn color="primary" icon="add" label="Buat Surat Kedinasan" @click="openDialogKedinasan" />
    </template>
  </q-table>

  <SuratKeluarFormKedinasan />
</template>

<script setup>
import { watch } from 'vue'
import { useOutgoingLetter } from '@/composables/useOutgoingLetter.js'
import SuratKeluarFormKedinasan from '../forms/SuratKeluarFormKedinasan.vue'

const {
  listKedinasan,
  isLoading,
  filterKedinasan,
  loadList,
  openDialogKedinasan
} = useOutgoingLetter()

const kolomKedinasan = [
  { name: 'nomorSurat', label: 'Nomor Surat', field: 'nomorSurat' },
  { name: 'jenis', label: 'Jenis', field: 'jenis' },
  { name: 'perihal', label: 'Perihal', field: 'perihal' },
  { name: 'tujuan', label: 'Tujuan', field: 'tujuan' },
  { name: 'tanggal', label: 'Tanggal', field: 'tanggal' },
  { name: 'status', label: 'Status', field: 'status' }
]

watch(filterKedinasan, () => {
  loadList('kedinasan')
})
</script>
