<template>
  <q-page padding>
    <div class="row q-col-gutter-md">
      <!-- Ringkasan Statistik -->
      <div class="col-12 col-md-4">
        <q-card class="bg-primary text-white">
          <q-card-section>
            <div class="text-h6">Surat Masuk Bulan Ini</div>
            <div class="text-h2">{{ stats.masuk }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card class="bg-secondary text-white">
          <q-card-section>
            <div class="text-h6">Surat Keluar Bulan Ini</div>
            <div class="text-h2">{{ stats.keluar }}</div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-md-4">
        <q-card class="bg-warning text-black">
          <q-card-section>
            <div class="text-h6">Disposisi Pending</div>
            <div class="text-h2">{{ stats.disposisiPending }}</div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Grafik Tren -->
    <q-card class="q-mt-md">
      <q-card-section>
        <div class="text-h6">Tren Surat Masuk & Keluar</div>
        <apexchart type="line" height="300" :options="chartOptions" :series="chartSeries" />
      </q-card-section>
    </q-card>

    <!-- Notifikasi/ Antrean Tugas (Role-Based) -->
    <q-card class="q-mt-md" v-if="context.role === 'kepsek'">
      <q-card-section>
        <div class="text-h6">Surat Masuk Menunggu Disposisi</div>
        <q-list separator>
          <q-item v-for="surat in pendingDisposisi" :key="surat.id" clickable @click="goToDisposisi(surat.id)">
            <q-item-section avatar>
              <q-icon name="mail" color="primary" />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ surat.perihalSurat || surat.perihal }}</q-item-label>
              <q-item-label caption>{{ surat.asal }} - {{ surat.tanggalDiterima }}</q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card-section>
    </q-card>

    <q-card class="q-mt-md" v-if="context.role === 'waka' || context.role === 'guru'">
      <q-card-section>
        <div class="text-h6">Instruksi Disposisi Baru</div>
        <q-list separator>
          <q-item v-for="item in tugasDisposisi" :key="item.id" clickable>
            <q-item-section avatar>
              <q-icon name="assignment" color="orange" />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ item.instruksi }}</q-item-label>
              <q-item-label caption>Dari: {{ item.dari }} - {{ item.tanggal }}</q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </q-card-section>
    </q-card>

    <!-- Tombol Pintas -->
    <div class="q-mt-md row q-col-gutter-sm">
      <div class="col-6">
        <q-btn color="positive" icon="add" label="Input Surat Masuk Baru"
          @click="$router.push('/persuratan/surat-masuk')" />
      </div>
      <div class="col-6">
        <q-btn color="purple" icon="edit" label="Buat Nomor Surat Keluar"
          @click="$router.push('/persuratan/surat-keluar')" />
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useDashboardStats } from '@/composables/useDashboardStats.js'
import { useOperationalContext } from '@/composables/useOperationalContext.js'

const router = useRouter()
const context = useOperationalContext()

const {
  stats,
  chartOptions,
  chartSeries,
  pendingDisposisi,
  tugasDisposisi,
  loadDashboardData
} = useDashboardStats()

function goToDisposisi(id) {
  console.log(id)
  router.push(`/persuratan/disposisi`)
}

onMounted(() => {
  loadDashboardData()
})
</script>
