<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Audit Log</div>
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-3">
        <q-input v-model="search" label="Cari (user, aksi, detail)" outlined dense />
      </div>
      <div class="col-12 col-md-2">
        <q-btn color="negative" label="Hapus Semua" flat @click="clearAll" />
      </div>
    </div>
    <q-table :rows="filteredLogs" :columns="columns" row-key="id" flat bordered dense />
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useAuditLogStore } from 'stores/pengaturan/auditLog'

const $q = useQuasar()
const auditStore = useAuditLogStore()
const search = ref('')

const filteredLogs = computed(() => {
  if (!search.value) return auditStore.list
  const q = search.value.toLowerCase()
  return auditStore.list.filter(l => l.user?.toLowerCase().includes(q) || l.aksi?.toLowerCase().includes(q) || l.detail?.toLowerCase().includes(q))
})
const columns = [
  { name: 'waktu', label: 'Waktu', field: 'waktu' },
  { name: 'user', label: 'User', field: 'user' },
  { name: 'role', label: 'Role', field: 'role' },
  { name: 'aksi', label: 'Aksi', field: 'aksi' },
  { name: 'detail', label: 'Detail', field: 'detail' }
]

function clearAll() {
  $q.dialog({ title: 'Konfirmasi', message: 'Hapus semua log?', cancel: true }).onOk(() => {
    auditStore.clear()
    $q.notify({ type: 'positive', message: 'Semua log dihapus' })
  })
}
onMounted(() => {
  auditStore.loadData()
})
</script>
