<template>
  <q-page class="q-pa-md">
    <div class="text-h5">Pengajuan Mutasi (Pindah/Keluar)</div>
    <q-card flat bordered class="q-mt-md">
      <q-card-section>
        <q-form @submit="submitMutasi" class="q-gutter-md">
          <q-select v-model="form.jenis_mutasi" :options="['Pindah Sekolah', 'Keluar (Drop Out)']" label="Jenis Mutasi"
            outlined />
          <q-input v-model="form.alasan" label="Alasan" type="textarea" rows="3" outlined />
          <q-file v-model="form.dokumen" label="Upload Dokumen Pendukung (SK Pindah, dll)" accept=".pdf,.jpg,.png"
            outlined />
          <div class="q-gutter-sm">
            <q-btn label="Ajukan" type="submit" color="primary" />
            <q-btn label="Batal" flat color="negative" @click="resetForm" />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'
import { useMutasiStore } from '@/stores/kesiswaan/mutasiStore'
// import { useAuthStore } from 'stores/auth-storeth'

const $q = useQuasar()
const mutasiStore = useMutasiStore()
// const authStore = useAuthStore()

const form = ref({ jenis_mutasi: '', alasan: '', dokumen: null })

function submitMutasi() {
  // ambil data siswa dari auth (misal siswa login)
  // const siswa = { id: 1, nama: 'Ahmad Faizal' } // sementara
  let dokumenBase64 = ''
  if (form.value.dokumen) {
    const reader = new FileReader()
    reader.onload = (e) => {
      dokumenBase64 = e.target.result
      simpan(dokumenBase64)
    }
    reader.readAsDataURL(form.value.dokumen)
  } else {
    simpan('')
  }
}

function simpan(base64) {
  mutasiStore.tambah({
    siswa_id: 1,
    siswa_nama: 'Ahmad Faizal',
    tanggal_pengajuan: new Date().toISOString(),
    jenis_mutasi: form.value.jenis_mutasi,
    alasan: form.value.alasan,
    dokumen: base64
  })
  $q.notify({ type: 'positive', message: 'Pengajuan mutasi dikirim' })
  resetForm()
}
function resetForm() {
  form.value = { jenis_mutasi: '', alasan: '', dokumen: null }
}
</script>
