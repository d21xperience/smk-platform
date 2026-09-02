<!-- components/TikTokPlayer.vue -->
<template>
  <div class="tiktok-player-container">
    <iframe
      ref="playerIframe"
      :src="currentVideoSrc"
      class="tiktok-iframe"
      allow="autoplay; encrypted-media; fullscreen"
      allowfullscreen
      @load="onIframeLoad"
    />
  </div>
</template>

<script setup>
import { ref, computed, onBeforeUnmount } from 'vue'

// 1. DEKLARASI ref untuk iframe
const playerIframe = ref(null)

// Props
const props = defineProps({
  videoIds: {
    type: Array,
    required: true,
    validator: (value) => value.every(id => typeof id === 'string')
  },
  autoplay: {
    type: Boolean,
    default: true
  },
  muted: {
    type: Boolean,
    default: true
  },
  // Durasi estimasi per video (dalam milidetik) untuk auto-next
  estimatedDuration: {
    type: Number,
    default: 30000 // 30 detik
  }
})

// State
const currentIndex = ref(0)
const isFirstLoad = ref(true)
let autoNextTimer = null

// Computed: URL pemutar TikTok
const currentVideoSrc = computed(() => {
  const videoId = props.videoIds[currentIndex.value]
  return `https://www.tiktok.com/player/v1/${videoId}?autoplay=${props.autoplay ? 1 : 0}&muted=${props.muted ? 1 : 0}&loop=0&controls=1`
})

// Fungsi untuk berpindah ke video berikutnya
const playNext = () => {
  if (props.videoIds.length === 0) return
  currentIndex.value = (currentIndex.value + 1) % props.videoIds.length

  // Perbarui src iframe secara manual untuk memuat video baru
  const iframe = playerIframe.value
  if (iframe) {
    iframe.src = currentVideoSrc.value
  }

  // Reset timer untuk auto-next (jika digunakan)
  resetAutoNextTimer()
}

// Timer untuk auto-next (sebagai pengganti event 'ended' yang tidak bisa dideteksi)
const startAutoNextTimer = () => {
  if (props.estimatedDuration <= 0) return
  clearTimeout(autoNextTimer)
  autoNextTimer = setTimeout(() => {
    playNext()
  }, props.estimatedDuration)
}

const resetAutoNextTimer = () => {
  clearTimeout(autoNextTimer)
  if (props.videoIds.length > 1) {
    startAutoNextTimer()
  }
}

// Event handler saat iframe selesai dimuat
const onIframeLoad = () => {
  if (isFirstLoad.value) {
    isFirstLoad.value = false
    // Jika autoplay aktif, mulai timer setelah video dimuat
    if (props.autoplay && props.videoIds.length > 1) {
      startAutoNextTimer()
    }
    return
  }

  // Untuk pemuatan berikutnya (setelah ganti video), restart timer
  if (props.autoplay && props.videoIds.length > 1) {
    startAutoNextTimer()
  }
}

// Ekspos metode ke komponen induk
defineExpose({
  playNext,
  currentIndex,
  // Opsional: tambahkan metode untuk mengatur indeks secara manual
  goTo: (index) => {
    if (index >= 0 && index < props.videoIds.length) {
      currentIndex.value = index
      const iframe = playerIframe.value
      if (iframe) {
        iframe.src = currentVideoSrc.value
      }
      resetAutoNextTimer()
    }
  }
})

// Bersihkan timer saat komponen di-unmount
onBeforeUnmount(() => {
  clearTimeout(autoNextTimer)
})
</script>

<style scoped>
.tiktok-player-container {
  width: 100%;
  height: 100%;
  position: relative;
  background: #000;
}

.tiktok-iframe {
  width: 100%;
  height: 100%;
  border: none;
  display: block;
}
</style>
