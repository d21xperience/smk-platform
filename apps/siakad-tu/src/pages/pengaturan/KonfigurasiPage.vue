<template>
  <q-page class="q-pa-md">
    <div class="text-h5 q-mb-md">Konfigurasi Sistem</div>
    <q-card flat bordered>
      <q-card-section>
        <q-form @submit="simpan" class="q-gutter-md">
          <div class="row q-col-gutter-md">
            <div class="col-6">
              <q-select v-model="config.semester_aktif" :options="['Ganjil', 'Genap']" label="Semester Aktif"
                outlined />
            </div>
            <div class="col-6">
              <q-input v-model="config.tahun_ajaran_aktif" label="Tahun Ajaran Aktif" outlined />
            </div>
            <div class="col-6">
              <q-input v-model.number="config.kkm_default" label="KKM Default" type="number" min="0" max="100"
                outlined />
            </div>
            <div class="col-6">
              <q-toggle v-model="config.auto_lock_ujian" label="Auto Lock Ujian (anti cheat)" />
            </div>
            <div class="col-6">
              <q-input v-model.number="config.max_upload_size_mb" label="Max Upload File (MB)" type="number" outlined />
            </div>
            <div class="col-6">
              <q-input v-model.number="config.waktu_minimal_pengerjaan_menit" label="Waktu Minimal Pengerjaan (menit)"
                type="number" outlined />
            </div>
          </div>
          <q-btn label="Simpan Konfigurasi" type="submit" color="primary" />
        </q-form>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { reactive, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useKonfigurasiStore } from 'stores/pengaturan/konfigurasi'

const $q = useQuasar()
const konfigStore = useKonfigurasiStore()
const config = reactive({
  semester_aktif: '',
  tahun_ajaran_aktif: '',
  kkm_default: 75,
  auto_lock_ujian: false,
  max_upload_size_mb: 2,
  waktu_minimal_pengerjaan_menit: 30
})

function simpan() {
  konfigStore.update(config)
  $q.notify({ type: 'positive', message: 'Konfigurasi disimpan' })
}
onMounted(() => {
  konfigStore.loadData()
  Object.assign(config, konfigStore)
})
</script>
