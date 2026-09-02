<template>
  <div class="ticker-wrapper bg-dark text-white row items-center no-wrap">
    <!-- Label Indikator Kiri (Tetap Diam/Tidak Ikut Berjalan) -->
    <div class="ticker-label bg-amber-9 text-dark text-weight-bolder row items-center q-px-md z-top">
      <q-icon name="campaign" size="50px" class="q-mr-xs" />
      INFO HARI INI
    </div>

    <!-- Area Teks Berjalan (Menggunakan CSS Animation agar Mulus) -->
    <div class="ticker-content-container col row items-center">
      <div class="ticker-text text-weight-medium">
        <!-- Kehadiran Siswa -->
        <span class="q-mr-xl">
          <q-icon name="groups" color="blue-4" size="18px" class="q-mr-xs" />
          Kehadiran Siswa:
          <span class="text-green-4 text-weight-bold">Hadir ({{ props.kehadiran.siswa.hadir }})</span> |
          <span class="text-amber-4 text-weight-bold">Sakit ({{ props.kehadiran.siswa.sakit }})</span> |
          <span class="text-orange-4 text-weight-bold">Izin ({{ props.kehadiran.siswa.izin }})</span> |
          <span class="text-red-4 text-weight-bold">Alfa ({{ props.kehadiran.siswa.alfa }})</span>
        </span>

        <!-- Pembatas Antar Informasi -->
        <span class="text-grey-6 q-mr-xl">•</span>

        <!-- Kehadiran Guru -->
        <span class="q-mr-xl">
          <q-icon name="co_present" color="teal-4" size="48px" class="q-mr-xs" />
          Kehadiran Guru & Staf:
          <span class="text-green-4 text-weight-bold">Hadir ({{ props.kehadiran.guru.hadir }}/{{
            props.kehadiran.guru.total
            }})</span> |
          <span class="text-red-4 text-weight-bold">Absen ({{ props.kehadiran.guru.absen }})</span>
        </span>

        <!-- Duplikasi Teks untuk Efek Loop Tanpa Putus (Seamless) -->
        <span class="text-grey-6 q-mr-xl">•</span>
        <span class="q-mr-xl">
          <q-icon name="groups" color="blue-4" size="48px" class="q-mr-xs" />
          Kehadiran Siswa: <span class="text-green-4 text-weight-bold">Hadir ({{ props.kehadiran.siswa.hadir
            }})</span>...
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>

// Definisikan props agar data kehadiran bisa dikirim dinamis dari database/parent
const props = defineProps({
  kehadiran: {
    type: Object,
    default: () => ({
      siswa: { hadir: '94.2%', sakit: 12, izin: 8, alfa: 3 },
      guru: { hadir: 42, total: 45, absen: 3 }
    })
  }
})
</script>

<style scoped>
.ticker-wrapper {
  height: 65px;
  overflow: hidden;
  border-bottom: 2px solid var(--q-primary);
  box-shadow: inset 0 -5px 10px rgba(0, 0, 0, 0.2);
}

.ticker-label {
  height: 100%;
  box-shadow: 4px 0 10px rgba(0, 0, 0, 0.5);
  white-space: nowrap;
  font-size: 12px;
  letter-spacing: 0.5px;
}

.ticker-content-container {
  overflow: hidden;
  position: relative;
  height: 100%;
}

/* Animasi CSS Bergerak dari Kanan ke Kiri secara Smooth */
.ticker-text {
  display: inline-block;
  white-space: nowrap;
  padding-left: 100%;
  animation: marquee-animation 30s linear infinite;
  font-size: 30px;
  letter-spacing: 0.3px;
}

/* Efek Berhenti Berjalan saat Kursor Mouse Diarahkan ke Teks */
.ticker-content-container:hover .ticker-text {
  animation-play-state: paused;
  cursor: pointer;
}

@keyframes marquee-animation {
  0% {
    transform: translate3d(0, 0, 0);
  }

  100% {
    transform: translate3d(-100%, 0, 0);
  }
}
</style>
