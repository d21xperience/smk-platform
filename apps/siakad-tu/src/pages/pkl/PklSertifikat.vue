<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Generate Sertifikat PKL</div>
    <q-select v-model="selectedPenempatan" :options="penempatanOptions" label="Pilih Siswa" outlined />
    <q-btn color="primary" label="Cetak Sertifikat" @click="cetak" :disable="!selectedPenempatan" class="q-mt-md" />
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { usePenempatanStore } from '@/stores/pkl/penempatan'
import { usePenilaianStore } from '@/stores/pkl/penilaian'

const penempatanStore = usePenempatanStore()
const penilaianStore = usePenilaianStore()
const selectedPenempatan = ref(null)

const penempatanOptions = computed(() => {
  return penempatanStore.list.map(p => ({
    id: p.id,
    label: `${p.siswa_nama} - ${p.mitra_nama}`,
    ...p
  }))
})

function cetak() {
  const nilai = penilaianStore.getRataRata(selectedPenempatan.value.siswa_id, selectedPenempatan.value.id)
  const content = `
    <div style="text-align:center; padding:50px; border:2px solid #ccc;">
      <h1>SERTIFIKAT PKL</h1>
      <p>Diberikan kepada:</p>
      <h2>${selectedPenempatan.value.siswa_nama}</h2>
      <p>Telah menyelesaikan Praktek Kerja Lapangan di</p>
      <h3>${selectedPenempatan.value.mitra_nama}</h3>
      <p>Periode: ${selectedPenempatan.value.tanggal_mulai} s.d. ${selectedPenempatan.value.tanggal_selesai}</p>
      <p>Dengan nilai rata-rata: ${nilai ? nilai.toFixed(1) : 'Belum dinilai'}</p>
      <p>Bandung, ${new Date().toLocaleDateString('id-ID')}</p>
      <p>Kepala Sekolah</p>
    </div>
  `
  const win = window.open('', '_blank')
  win.document.write(`<html><head><title>Sertifikat PKL</title></head><body>${content}</body></html>`)
  win.document.close()
  win.print()
}

onMounted(() => {
  penempatanStore.loadData()
  penilaianStore.loadData()
})
</script>
