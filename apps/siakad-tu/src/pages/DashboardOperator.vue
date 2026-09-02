<template>
  <div>
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-4">
        <q-card class="bg-primary text-white">
          <q-card-section>
            <div class="text-h6">Status Sinkron Dapodik</div>
            <div class="text-h3">{{ syncStatus }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card class="bg-secondary text-white">
          <q-card-section>
            <div class="text-h6">Backup Terakhir</div>
            <div class="text-h3">{{ lastBackup }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card class="bg-warning text-black">
          <q-card-section>
            <div class="text-h6">Error Ticket</div>
            <div class="text-h3">{{ errorTickets.length }}</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <div class="col-12">
        <q-card>
          <q-card-section>
            <div class="text-h6">Daftar Error Ticket</div>
            <q-table :rows="errorTickets" :columns="errorColumns" dense flat />
          </q-card-section>
        </q-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const syncStatus = ref('Terakhir sinkron: 2 jam lalu')
const lastBackup = ref('2025-06-10 03:00')
const errorTickets = ref([
  { id: 1, judul: 'Gagal sinkron siswa', tanggal: '2025-06-15', status: 'open' },
  { id: 2, judul: 'Backup gagal karena storage penuh', tanggal: '2025-06-14', status: 'closed' }
])
const errorColumns = [
  { name: 'judul', label: 'Judul', field: 'judul' },
  { name: 'tanggal', label: 'Tanggal', field: 'tanggal' },
  { name: 'status', label: 'Status', field: 'status' }
]
</script>
