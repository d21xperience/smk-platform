<template>
  <div>
    <!-- ===== MODE SELECTOR ===== -->
    <div class="row q-gutter-sm q-mb-md">
      <q-btn
        :outline="mode !== 'manual'"
        :color="mode === 'manual' ? 'primary' : 'grey-5'"
        icon="keyboard"
        label="Input Manual"
        size="sm"
        class="col"
        @click="switchMode('manual')"
      />
      <q-btn
        :outline="mode !== 'scan'"
        :color="mode === 'scan' ? 'primary' : 'grey-5'"
        icon="qr_code_scanner"
        label="Scan QR Code"
        size="sm"
        class="col"
        :disable="!hasCamera && !cameraChecked"
        @click="switchMode('scan')"
      />
    </div>

    <!-- ===== MODE MANUAL ===== -->
    <div v-if="mode === 'manual'">
      <q-input
        v-model="serialInput"
        label="Serial Kartu Ujian"
        outlined
        clearable
        :placeholder="placeholder"
        :disable="disable"
        :error="!!inputError"
        :error-message="inputError"
        @update:model-value="onManualInput"
        @keyup.enter="submitManual"
      >
        <template #prepend>
          <q-icon name="credit_card" color="primary" />
        </template>
        <template #append>
          <q-btn
            v-if="serialInput"
            flat
            round
            dense
            icon="send"
            color="primary"
            :loading="loading"
            @click="submitManual"
          />
        </template>
      </q-input>

      <!-- Format hint -->
      <div class="text-caption text-grey-6 q-mt-xs q-ml-sm">
        Format: KU-YYYY-NNNNN (contoh: KU-2024-00123)
      </div>
    </div>

    <!-- ===== MODE SCAN ===== -->
    <div v-else-if="mode === 'scan'">
      <!-- Error kamera -->
      <q-banner v-if="scanError" class="bg-negative text-white q-mb-md" rounded dense>
        <template #avatar><q-icon name="error" /></template>
        {{ scanError }}
        <template #action>
          <q-btn flat label="Coba Lagi" @click="initScan" />
        </template>
      </q-banner>

      <!-- Area Kamera -->
      <div class="camera-wrapper q-mb-md">
        <!-- Video element untuk QrScanner -->
        <video
          ref="videoRef"
          class="camera-video"
          :class="{ 'camera-active': scanning }"
          playsinline
          muted
        />

        <!-- Overlay loading sebelum kamera aktif -->
        <div v-if="!scanning && !scanError" class="camera-overlay">
          <q-spinner-dots color="white" size="40px" />
          <div class="text-white text-caption q-mt-sm">Membuka kamera...</div>
        </div>

        <!-- Kontrol kamera (tampil saat aktif) -->
        <div v-if="scanning" class="camera-controls">
          <!-- Torch -->
          <q-btn
            round
            flat
            :icon="torchOn ? 'flash_on' : 'flash_off'"
            :color="torchOn ? 'yellow' : 'white'"
            size="sm"
            @click="toggleTorch"
          />

          <!-- Switch kamera -->
          <q-btn
            round
            flat
            icon="flip_camera_android"
            color="white"
            size="sm"
            @click="switchCamera"
          />
        </div>

        <!-- Panduan scan -->
        <div v-if="scanning" class="scan-guide">
          <div class="scan-frame" />
          <div class="text-white text-caption text-center q-mt-sm">
            Arahkan QR Code kartu ke dalam kotak
          </div>
        </div>
      </div>

      <!-- Tombol tutup -->
      <q-btn
        outline
        color="negative"
        icon="close"
        label="Tutup Kamera"
        class="full-width"
        @click="switchMode('manual')"
      />
    </div>

    <!-- ===== HASIL (tampil di kedua mode) ===== -->
    <q-slide-transition>
      <q-card v-if="verifiedSerial" flat bordered class="q-mt-md bg-green-1">
        <q-card-section class="q-py-sm row items-center q-gutter-sm">
          <q-icon name="check_circle" color="positive" size="28px" />
          <div>
            <div class="text-caption text-grey-7">Serial Terdeteksi</div>
            <div class="text-body1 text-weight-bold">{{ verifiedSerial }}</div>
          </div>
          <q-space />
          <q-btn flat round dense icon="close" color="grey-6" size="sm" @click="clearResult" />
        </q-card-section>
      </q-card>
    </q-slide-transition>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useCardScanner } from '../composables/useCardScanner'

// =====================
// PROPS & EMITS
// =====================
const props = defineProps({
  modelValue: { type: String, default: '' },
  loading: { type: Boolean, default: false },
  disable: { type: Boolean, default: false },
  placeholder: { type: String, default: 'Ketik serial kartu...' },
})

const emit = defineEmits(['update:modelValue', 'submit'])

// =====================
// STATE
// =====================
const mode = ref('manual') // 'manual' | 'scan'
const serialInput = ref(props.modelValue || '')
const verifiedSerial = ref('')
const inputError = ref('')
const cameraChecked = ref(false)
const videoRef = ref(null)

const {
  scanning,
  torchOn,
  hasCamera,
  scanError,
  checkCamera,
  startScan,
  stopScan,
  switchCamera,
  toggleTorch,
} = useCardScanner()

// =====================
// FORMAT VALIDATOR
// =====================
const SERIAL_REGEX = /^KU-\d{4}-\d{5}$/

const validateSerial = (serial) => {
  if (!serial) return 'Serial tidak boleh kosong.'
  if (!SERIAL_REGEX.test(serial.trim().toUpperCase())) {
    return 'Format tidak valid. Contoh: KU-2024-00123'
  }
  return ''
}

// =====================
// MODE MANUAL
// =====================
const onManualInput = (val) => {
  inputError.value = ''
  // Auto format: huruf besar & tambah strip otomatis
  const formatted = formatSerial(val)
  serialInput.value = formatted
}

const formatSerial = (raw) => {
  // Hapus karakter selain huruf & angka
  const clean = raw.toUpperCase().replace(/[^A-Z0-9]/g, '')

  // Sisipkan strip: KU-2024-00123
  if (clean.length <= 2) return clean
  if (clean.length <= 6) return `${clean.slice(0, 2)}-${clean.slice(2)}`
  return `${clean.slice(0, 2)}-${clean.slice(2, 6)}-${clean.slice(6, 11)}`
}

const submitManual = () => {
  const err = validateSerial(serialInput.value)
  if (err) {
    inputError.value = err
    return
  }
  setResult(serialInput.value.trim().toUpperCase())
}

// =====================
// MODE SCAN
// =====================
const initScan = async () => {
  await nextTick()
  await startScan(videoRef.value, onQrResult)
}

const switchMode = async (newMode) => {
  if (mode.value === 'scan') stopScan()

  mode.value = newMode

  if (newMode === 'scan') {
    await nextTick()
    await initScan()
  }
}

const onQrResult = (data) => {
  // QR bisa berisi serial langsung atau URL dengan serial
  // Contoh: "KU-2024-00123" atau "https://sekolah.sch.id/card?s=KU-2024-00123"
  let serial = data

  if (data.includes('?s=')) {
    serial = new URL(data).searchParams.get('s') || data
  }

  const err = validateSerial(serial)
  if (err) {
    scanError.value = `QR tidak valid: ${err}`
    return
  }

  setResult(serial.trim().toUpperCase())
  mode.value = 'manual' // kembali ke manual setelah scan sukses
}

// =====================
// RESULT
// =====================
const setResult = (serial) => {
  verifiedSerial.value = serial
  serialInput.value = serial
  emit('update:modelValue', serial)
  emit('submit', serial)
}

const clearResult = () => {
  verifiedSerial.value = ''
  serialInput.value = ''
  inputError.value = ''
  emit('update:modelValue', '')
}

// =====================
// LIFECYCLE
// =====================
onMounted(async () => {
  await checkCamera()
  cameraChecked.value = true
})
</script>

<style scoped>
.camera-wrapper {
  position: relative;
  width: 100%;
  aspect-ratio: 4/3;
  background: #000;
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.camera-video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  opacity: 0;
  transition: opacity 0.3s;
}

.camera-video.camera-active {
  opacity: 1;
}

.camera-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.6);
}

.camera-controls {
  position: absolute;
  bottom: 12px;
  right: 12px;
  display: flex;
  gap: 8px;
}

.scan-guide {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.scan-frame {
  width: 200px;
  height: 200px;
  border: 3px solid #fff;
  border-radius: 12px;
  box-shadow: 0 0 0 9999px rgba(0, 0, 0, 0.4);
}
</style>
