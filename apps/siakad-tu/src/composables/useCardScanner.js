// src/composables/useCardScanner.js
import { ref, onUnmounted } from 'vue'
import QrScanner from 'qr-scanner'

export function useCardScanner() {
  const scanning = ref(false)
  const torchOn = ref(false)
  const hasCamera = ref(false)
  const scanError = ref(null)
  const cameras = ref([]) // daftar kamera tersedia
  const activeCamera = ref('environment') // 'environment' = belakang

  let qrScanner = null

  // =====================
  // CEK KETERSEDIAAN KAMERA
  // =====================
  const checkCamera = async () => {
    hasCamera.value = await QrScanner.hasCamera()
    if (hasCamera.value) {
      cameras.value = await QrScanner.listCameras(true)
    }
    return hasCamera.value
  }

  // =====================
  // MULAI SCAN
  // =====================
  const startScan = async (videoElement, onResult) => {
    scanError.value = null

    if (!videoElement) {
      scanError.value = 'Elemen video tidak ditemukan.'
      return
    }

    try {
      qrScanner = new QrScanner(
        videoElement,
        (result) => {
          // result.data = string dari QR code
          onResult(result.data)
          stopScan()
        },
        {
          preferredCamera: activeCamera.value,
          highlightScanRegion: true, // kotak panduan scan
          highlightCodeOutline: true, // outline QR yang terdeteksi
          maxScansPerSecond: 10,
        },
      )

      await qrScanner.start()
      scanning.value = true
    } catch (err) {
      scanError.value = err.message?.includes('permission')
        ? 'Akses kamera ditolak. Izinkan akses kamera di browser.'
        : 'Kamera tidak dapat dibuka. Coba refresh halaman.'
    }
  }

  // =====================
  // STOP SCAN
  // =====================
  const stopScan = () => {
    qrScanner?.stop()
    qrScanner?.destroy()
    qrScanner = null
    scanning.value = false
    torchOn.value = false
  }

  // =====================
  // GANTI KAMERA
  // =====================
  const switchCamera = async () => {
    if (!qrScanner) return
    activeCamera.value = activeCamera.value === 'environment' ? 'user' : 'environment'
    await qrScanner.setCamera(activeCamera.value)
  }

  // =====================
  // TORCH / FLASH
  // =====================
  const toggleTorch = async () => {
    if (!qrScanner) return
    await qrScanner.toggleFlash()
    torchOn.value = !torchOn.value
  }

  // Cleanup otomatis saat komponen unmount
  onUnmounted(() => stopScan())

  return {
    scanning,
    torchOn,
    hasCamera,
    scanError,
    cameras,
    activeCamera,
    checkCamera,
    startScan,
    stopScan,
    switchCamera,
    toggleTorch,
  }
}
