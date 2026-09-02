<template>
  <q-page padding>
    <div class="text-h5 q-mb-md">Sinkronisasi Dapodik</div>

    <q-card class="q-mb-md">
      <q-card-section>
        <div class="text-subtitle1">Konfigurasi Web Service</div>
        <q-input v-model="urlDapodik" label="URL Endpoint" filled class="q-mb-sm" />
        <q-input v-model="tokenDapodik" label="Token" filled type="password" class="q-mb-sm" />
        <q-btn color="primary" label="Test Koneksi" @click="testKoneksi" :loading="testing" />
        <q-chip
          v-if="statusTest !== null"
          :color="statusTest ? 'green' : 'red'"
          text-color="white"
          class="q-ml-sm"
        >
          {{ statusTest ? 'Berhasil terhubung' : 'Gagal koneksi' }}
        </q-chip>
      </q-card-section>
    </q-card>

    <div class="row q-col-gutter-md">
      <!-- Tarik Siswa -->
      <div class="col-12 col-md-4">
        <q-card>
          <q-card-section>
            <div class="text-h6">Data Siswa</div>
            <q-linear-progress
              v-if="progress.siswa > 0 && progress.siswa < 100"
              :value="progress.siswa / 100"
              color="primary"
              class="q-mt-sm"
            />
            <p class="text-caption">
              {{ progress.siswa > 0 ? `${progress.siswa}%` : 'Belum dimulai' }}
            </p>
            <q-btn
              color="primary"
              label="Tarik Siswa"
              @click="tarikData('siswa')"
              :disable="syncing.siswa"
            >
              <q-spinner v-if="syncing.siswa" size="sm" color="white" class="q-ml-sm" />
            </q-btn>
          </q-card-section>
        </q-card>
      </div>

      <!-- Tarik GTK -->
      <div class="col-12 col-md-4">
        <q-card>
          <q-card-section>
            <div class="text-h6">Data GTK</div>
            <q-linear-progress
              v-if="progress.gtk > 0 && progress.gtk < 100"
              :value="progress.gtk / 100"
              color="secondary"
              class="q-mt-sm"
            />
            <p class="text-caption">
              {{ progress.gtk > 0 ? `${progress.gtk}%` : 'Belum dimulai' }}
            </p>
            <q-btn
              color="secondary"
              label="Tarik GTK"
              @click="tarikData('gtk')"
              :disable="syncing.gtk"
            >
              <q-spinner v-if="syncing.gtk" size="sm" color="white" class="q-ml-sm" />
            </q-btn>
          </q-card-section>
        </q-card>
      </div>

      <!-- Tarik Rombel -->
      <div class="col-12 col-md-4">
        <q-card>
          <q-card-section>
            <div class="text-h6">Rombel</div>
            <q-linear-progress
              v-if="progress.rombel > 0 && progress.rombel < 100"
              :value="progress.rombel / 100"
              color="accent"
              class="q-mt-sm"
            />
            <p class="text-caption">
              {{ progress.rombel > 0 ? `${progress.rombel}%` : 'Belum dimulai' }}
            </p>
            <q-btn
              color="accent"
              label="Tarik Rombel"
              @click="tarikData('rombel')"
              :disable="syncing.rombel"
            >
              <q-spinner v-if="syncing.rombel" size="sm" color="white" class="q-ml-sm" />
            </q-btn>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

const urlDapodik = ref('http://localhost:5774/WebService')
const tokenDapodik = ref('')
const testing = ref(false)
const statusTest = ref(null)

// Progress dan status sinkron
const progress = reactive({ siswa: 0, gtk: 0, rombel: 0 })
const syncing = reactive({ siswa: false, gtk: false, rombel: false })

async function testKoneksi() {
  testing.value = true
  statusTest.value = null
  try {
    // Simulasi fetch
    await new Promise((resolve) => setTimeout(resolve, 1500))
    // Ganti dengan axios.get ke endpoint test
    statusTest.value = true
    $q.notify({ message: 'Koneksi berhasil', color: 'positive' })
  } catch (e) {
    statusTest.value = false
    $q.notify({ message: `Gagal terhubung ${e}`, color: 'negative' })
  } finally {
    testing.value = false
  }
}

async function tarikData(jenis) {
  syncing[jenis] = true
  progress[jenis] = 0
  try {
    for (let i = 1; i <= 10; i++) {
      await new Promise((resolve) => setTimeout(resolve, 500))
      progress[jenis] = i * 10
    }
    $q.notify({ message: `Penarikan data ${jenis.toUpperCase()} selesai`, color: 'positive' })
  } catch (e) {
    $q.notify({ message: `Gagal menarik data ${e}`, color: 'negative' })
  } finally {
    syncing[jenis] = false
  }
}
</script>
