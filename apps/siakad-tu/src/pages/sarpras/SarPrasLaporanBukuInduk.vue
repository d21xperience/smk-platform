<template>
  <q-page class="q-pa-md">
    <q-card flat bordered>
      <q-card-section class="row items-center justify-between bg-primary text-white">
        <div class="text-h6"><q-icon name="menu_book" class="q-mr-sm" /> Buku Induk Inventaris</div>
        <q-btn color="white" text-color="primary" icon="picture_as_pdf" label="Cetak PDF Buku Induk"
          @click="cetakLaporan" :loading="isDownloading" no-caps />
      </q-card-section>

      <q-card-section>
        <div class="text-subtitle1 q-mb-sm">Daftar Seluruh Aset Inventaris</div>
        <q-table :rows="kirStore.items" :columns="columns" row-key="id" flat bordered :loading="kirStore.isLoading"
          :rows-per-page-options="[10, 25, 50, 0]">

          <template v-slot:body-cell-status_aset="props">
            <q-td :props="props" class="text-center">
              <q-badge :color="props.value === 'AKTIF' ? 'green' : props.value === 'DIHAPUS' ? 'red' : 'orange'">
                {{ props.value }}
              </q-badge>
            </q-td>
          </template>

          <template v-slot:body-cell-nilai_buku="props">
            <q-td :props="props" class="text-right">
              Rp {{ props.value.toLocaleString('id-ID') }}
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useKirStore } from '@/stores/sarpras/kir'

const $q = useQuasar()
const kirStore = useKirStore()
const isDownloading = ref(false)

const columns = [
  { name: 'nomor_registrasi', label: 'No. Registrasi', field: 'nomor_registrasi', align: 'left', sortable: true },
  { name: 'id', label: 'Kode Aset', field: 'id', align: 'left', sortable: true },
  { name: 'name', label: 'Nama Barang', field: 'name', align: 'left', sortable: true },
  { name: 'brand', label: 'Merk', field: 'brand', align: 'left' },
  { name: 'purchase_date', label: 'Tahun', field: 'purchase_date', align: 'center', format: val => new Date(val).getFullYear() },
  { name: 'nilai_buku', label: 'Nilai Buku', field: 'nilai_buku', align: 'right' },
  { name: 'condition', label: 'Kondisi', field: 'condition', align: 'center' },
  { name: 'status_aset', label: 'Status', field: 'status_aset', align: 'center' }
]

onMounted(() => {
  kirStore.fetchAllItems()
})

async function cetakLaporan() {
  isDownloading.value = true
  try {
    await kirStore.downloadReport()
    $q.notify({ type: 'positive', message: 'Laporan PDF berhasil diunduh' })
  } catch (e) {
    $q.notify({ type: 'negative', message: `Gagal mengunduh laporan ${e}` })
  } finally {
    isDownloading.value = false
  }
}
</script>
