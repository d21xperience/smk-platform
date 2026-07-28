<template>
  <q-page padding>
    <div class="text-h5 q-mb-md">Riwayat Aktivitas</div>

    <q-table :rows="riwayat" :columns="columns" row-key="id" :filter="filter" :pagination="{ rowsPerPage: 15 }">
      <template v-slot:top-right>
        <q-input dense outlined v-model="filter" placeholder="Cari..." />
      </template>
      <template v-slot:body-cell-status="props">
        <q-td :props="props">
          <q-badge :color="statusColor(props.value)">{{ props.value }}</q-badge>
        </q-td>
      </template>
      <template v-slot:body-cell-catatan="props">
        <q-td :props="props">
          <span v-if="props.value">{{ props.value }}</span>
          <span v-else class="text-grey">-</span>
        </q-td>
      </template>
    </q-table>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'

const filter = ref('')

// Kolom tabel
const columns = [
  { name: 'waktu', label: 'Waktu', field: 'waktu', sortable: true },
  { name: 'tugas', label: 'Tugas', field: 'tugas' },
  { name: 'kategori', label: 'Kategori', field: 'kategori' },
  { name: 'status', label: 'Status', field: 'status' },
  { name: 'catatan', label: 'Catatan', field: 'catatan' },
  { name: 'pengguna', label: 'Pengguna', field: 'pengguna' }
]

// Data dummy
const riwayat = ref([
  { id: 1, waktu: '2026-05-31 10:30', tugas: 'Mutasi Masuk', kategori: 'Mutasi', status: 'Selesai', catatan: '', pengguna: 'Operator' },
  { id: 2, waktu: '2026-05-31 08:15', tugas: 'Pembaruan Biodata', kategori: 'Biodata', status: 'Ditolak', catatan: 'Dokumen tidak lengkap', pengguna: 'Operator' },
  { id: 3, waktu: '2026-05-30 14:20', tugas: 'Sinkronisasi GTK', kategori: 'GTK', status: 'Selesai', catatan: '', pengguna: 'Sistem' }
])

function statusColor(status) {
  if (status === 'Selesai') return 'positive'
  if (status === 'Ditolak') return 'negative'
  return 'grey'
}

// Auto-update bisa dilakukan dengan watch terhadap store jika ada
</script>
