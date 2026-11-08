<template>
  <q-page class="q-pa-md">
    <q-banner class="bg-info text-white q-mb-md" rounded>
      <template v-slot:avatar><q-icon name="info" color="white" /></template>
      Pilih aset dengan mencentang kotak di kiri tabel, lalu klik tombol cetak untuk mencetak stiker label.
    </q-banner>

    <q-card>
      <q-card-section class="row items-center justify-between">
        <div class="text-h6">Generator Stiker QR Code</div>
        <div class="row q-gutter-sm">
          <q-input v-model="search" dense outlined placeholder="Cari Kode/Nama Aset..." style="width: 250px">
            <template v-slot:append><q-icon name="search" /></template>
          </q-input>
          <q-btn color="primary" icon="print" label="CETAK STIKER TERPILIH" @click="cetakStiker" :loading="isPrinting"
            :disable="selected.length === 0" no-caps />
        </div>
      </q-card-section>

      <q-table :rows="filteredItems" :columns="columns" row-key="id" selection="multiple" v-model:selected="selected"
        flat bordered :loading="kirStore.isLoading" :rows-per-page-options="[0]">
        <template v-slot:body-cell-condition="props">
          <q-td :props="props">
            <q-badge :color="props.value === 'Bagus' ? 'green' : props.value === 'Rusak Ringan' ? 'orange' : 'red'">
              {{ props.value }}
            </q-badge>
          </q-td>
        </template>
      </q-table>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useKirStore } from '@/stores/sarpras/kir'

const $q = useQuasar()
const kirStore = useKirStore()

const search = ref('')
const selected = ref([])
const isPrinting = ref(false)

const columns = [
  { name: 'id', label: 'Kode Unik Aset', align: 'left', field: 'id', sortable: true },
  { name: 'name', label: 'Nama Aset', align: 'left', field: 'name', sortable: true },
  { name: 'room_name', label: 'Penempatan Ruang', align: 'left', field: 'room_name' },
  { name: 'condition', label: 'Kondisi', align: 'center', field: 'condition' }
]

const filteredItems = computed(() => {
  if (!search.value) return kirStore.items
  const lowerSearch = search.value.toLowerCase()
  return kirStore.items.filter(i =>
    i.id.toLowerCase().includes(lowerSearch) ||
    i.name.toLowerCase().includes(lowerSearch)
  )
})

onMounted(() => { kirStore.fetchAllItems() })

async function cetakStiker() {
  isPrinting.value = true
  try {
    const ids = selected.value.map(item => item.id)
    await kirStore.printQR(ids)
    $q.notify({ type: 'positive', message: 'File PDF berhasil diunduh' })
  } catch (err) {
    $q.notify({ type: 'negative', message: `Gagal mengunduh PDF ${err}` })
  } finally {
    isPrinting.value = false
  }
}
</script>
