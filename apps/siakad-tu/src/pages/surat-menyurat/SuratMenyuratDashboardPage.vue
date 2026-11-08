<template>
  <q-page class="q-pa-md">
    <div class="text-h5 text-weight-bold q-mb-md">Dashboard Persuratan</div>

    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-md-4">
        <q-card class="bg-blue-6 text-white">
          <q-card-section>
            <div class="text-subtitle1">Surat Masuk Bulan Ini</div>
            <div class="text-h3">{{ stats.masuk }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card class="bg-teal-6 text-white">
          <q-card-section>
            <div class="text-subtitle1">Surat Keluar Bulan Ini</div>
            <div class="text-h3">{{ stats.keluar }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card class="bg-amber-8 text-white">
          <q-card-section>
            <div class="text-subtitle1">Disposisi Pending</div>
            <div class="text-h3">{{ stats.disposisi }}</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <q-card flat bordered>
      <q-card-section>
        <div class="text-h6">Tren Surat Masuk & Keluar</div>
        <div class="text-caption text-grey">Grafik akan diintegrasikan pada iterasi berikutnya.</div>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '@/boot/axios'

const stats = ref({
  masuk: 0,
  keluar: 0,
  disposisi: 0
})

async function fetchStats() {
  try {
    const resIncoming = await api.get('/correspondences', { params: { type: 'INCOMING' } })
    const resOutgoing = await api.get('/correspondences', { params: { type: 'OUTGOING' } })
    
    stats.value.masuk = (resIncoming.data || []).length
    stats.value.keluar = (resOutgoing.data || []).length
    stats.value.disposisi = 0 // Placeholder untuk fitur disposisi
  } catch (err) {
    console.error('Gagal memuat stats:', err)
  }
}

onMounted(fetchStats)
</script>
