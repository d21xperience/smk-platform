<template>
  <q-layout view="lHh Lpr lFf" class="full-screen-layout">
    <AppHeader />

    <!-- Memanggil Komponen Sidebar -->
    <!-- <SidebarNavigation v-model="leftDrawerOpen" /> -->

    <q-page-container>
      <router-view />
    </q-page-container>
      <div class="absolute-bottom">
        <RunningTextTicker :kehadiran="dataKehadiranHariIni" />
      </div>
  </q-layout>
</template>

<script setup>
// import { ref } from 'vue'
import AppHeader from '@/components/AppHeader.vue'
import RunningTextTicker from '@/components/RunningTextTicker.vue'
import { ref } from 'vue';

// Data rekapitulasi kehadiran dinamis (Bisa Anda hubungkan dengan API backend nantinya)
const dataKehadiranHariIni = ref({
  siswa: { hadir: '1.142 Siswa', sakit: 14, izin: 5, alfa: 2 },
  guru: { hadir: 52, total: 55, absen: 3 }
})

</script>
<style scoped>
.full-screen-layout {
  height: 100vh;
  /* overflow: hidden; */
}

/* .full-page-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
} */

/* :deep(.q-page-container > :last-child) {
  flex: 1;
  min-height: 0;
} */
</style>
