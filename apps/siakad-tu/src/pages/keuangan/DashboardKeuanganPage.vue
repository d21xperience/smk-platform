<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Dashboard</div>
    <div class="row q-col-gutter-md">
      <div class="col-md-3 col-sm-6 col-xs-12" v-for="stat in stats" :key="stat.label">
        <q-card class="bg-primary text-white">
          <q-card-section>
            <div class="text-h6">{{ stat.label }}</div>
            <div class="text-h3">{{ stat.value }}</div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { siswaService } from '../../services/siswaService';
import { pembayaranSPPService } from '../../services/pembayaranSPPService';
import { pengeluaranService } from '../../services/pengeluaranService';

const stats = ref([
  { label: 'Total Siswa', value: 0 },
  { label: 'Total Pembayaran SPP', value: 0 },
  { label: 'Total Pengeluaran', value: 0 }
]);

onMounted(async () => {
  const siswa = await siswaService.getAll();
  const spp = await pembayaranSPPService.getAll();
  const pengeluaran = await pengeluaranService.getAll();
  stats.value[0].value = siswa.length;
  stats.value[1].value = spp.reduce((sum, p) => sum + p.jumlah, 0);
  stats.value[2].value = pengeluaran.reduce((sum, p) => sum + p.jumlah, 0);
});
</script>
