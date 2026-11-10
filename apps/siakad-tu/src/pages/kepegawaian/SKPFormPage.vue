<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row items-center q-mb-md">
      <q-btn flat round dense icon="arrow_back" @click="$router.back()" class="q-mr-sm" />
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Form Penilaian SKP</h1>
        <p class="text-caption text-grey-7 q-mb-none" v-if="penilaian">Pegawai: {{ penilaian.teacherName }}</p>
      </div>
      <q-space />
      <q-badge :color="getStatusColor(penilaian?.status)" class="q-pa-md text-subtitle2">
        STATUS: {{ penilaian?.status }}
      </q-badge>
    </div>

    <q-linear-progress v-if="isLoading" indeterminate color="primary" class="q-mb-md" />

    <q-card v-if="penilaian" class="shadow-1">
      <q-tabs v-model="tab" class="text-primary" align="left">
        <q-tab name="perilaku" icon="psychology" label="Perilaku Kerja (50%)" />
        <q-tab name="hasil" icon="work" label="Hasil Kerja" />
        <q-tab name="hasil_akhir" icon="assessment" label="Hasil Akhir & Aksi" />
      </q-tabs>
      <q-separator />

      <q-tab-panels v-model="tab" animated class="bg-white">
        <!-- TAB 1: PERILAKU -->
        <q-tab-panel name="perilaku">
          <div class="text-subtitle1 text-weight-bold text-primary q-mb-md">A. Aspek Perilaku Kerja</div>
          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6" v-for="(aspek, index) in aspekPerilaku" :key="index">
              <q-input 
                v-model.number="detailPerilaku[index].nilai" 
                :label="aspek" 
                type="number" 
                outlined 
                dense 
                :disable="!bisaEdit"
                :rules="[v => (v >= 0 && v <= 100) || 'Nilai 0-100']"
              />
            </div>
          </div>
        </q-tab-panel>

        <!-- TAB 2: HASIL KERJA -->
        <q-tab-panel name="hasil">
          <div class="text-subtitle1 text-weight-bold text-primary q-mb-md">B. Hasil Kerja Utama</div>
          <div v-for="(item, index) in detailHasilUtama" :key="'utama-'+index" class="row q-col-gutter-sm q-mb-sm">
            <div class="col-8">
              <q-input v-model="item.aspek" label="Nama Aspek / Indikator" outlined dense :disable="!bisaEdit" />
            </div>
            <div class="col-3">
              <q-input v-model.number="item.nilai" label="Nilai (0-100)" type="number" outlined dense :disable="!bisaEdit" />
            </div>
            <div class="col-1" v-if="bisaEdit">
              <q-btn flat round dense icon="delete" color="negative" @click="detailHasilUtama.splice(index, 1)" />
            </div>
          </div>
          <q-btn v-if="bisaEdit" flat color="secondary" icon="add" label="Tambah Aspek Utama" @click="tambahUtama" class="q-mb-md" />

          <q-separator class="q-my-md" />

          <div class="text-subtitle1 text-weight-bold text-primary q-mb-md">C. Hasil Kerja Tambahan (Opsional)</div>
          <div v-for="(item, index) in detailHasilTambahan" :key="'tambahan-'+index" class="row q-col-gutter-sm q-mb-sm">
            <div class="col-8">
              <q-input v-model="item.aspek" label="Nama Aspek / Indikator" outlined dense :disable="!bisaEdit" />
            </div>
            <div class="col-3">
              <q-input v-model.number="item.nilai" label="Nilai (0-100)" type="number" outlined dense :disable="!bisaEdit" />
            </div>
            <div class="col-1" v-if="bisaEdit">
              <q-btn flat round dense icon="delete" color="negative" @click="detailHasilTambahan.splice(index, 1)" />
            </div>
          </div>
          <q-btn v-if="bisaEdit" flat color="secondary" icon="add" label="Tambah Aspek Tambahan" @click="tambahTambahan" />
        </q-tab-panel>

        <!-- TAB 3: HASIL AKHIR -->
        <q-tab-panel name="hasil_akhir">
          <div v-if="hasil && hasil.nilaiAkhir > 0" class="row q-col-gutter-md q-mb-lg">
            <div class="col-6">
              <q-card class="bg-primary text-white text-center q-pa-md">
                <div class="text-subtitle2">Nilai Akhir</div>
                <div class="text-h3 text-weight-bold">{{ hasil.nilaiAkhir.toFixed(2) }}</div>
              </q-card>
            </div>
            <div class="col-6">
              <q-card class="bg-green text-white text-center q-pa-md">
                <div class="text-subtitle2">Predikat Akhir</div>
                <div class="text-h3 text-weight-bold">{{ hasil.predikatAkhir }}</div>
              </q-card>
            </div>
          </div>

          <div class="text-subtitle1 text-weight-bold text-primary q-mb-md">Catatan Evaluasi</div>
          <q-input 
            v-model="catatanPenilai" 
            label="Catatan Penilai Utama" 
            type="textarea" 
            rows="3" 
            outlined 
            :disable="!bisaEdit"
          />
          <q-input 
            v-if="penilaian.status === 'FINAL'" 
            v-model="catatanBerikat" 
            label="Catatan Atasan Berikat" 
            type="textarea" 
            rows="2" 
            outlined 
            readonly 
            class="q-mt-md bg-grey-2"
          />

          <!-- Tombol Aksi -->
          <div class="row justify-end q-mt-lg">
            <q-btn v-if="bisaEdit" label="Simpan & Submit" color="primary" @click="submitPenilaian" :loading="isSubmitting" />
            
            <template v-if="penilaian.status === 'SUBMITTED' && canApprove">
              <q-btn label="Kembalikan (Return)" color="orange" class="q-mr-sm" @click="promptReturn" :loading="isSubmitting" />
              <q-btn label="Setujui (Approve)" color="green" @click="promptApprove" :loading="isSubmitting" />
            </template>

            <q-btn v-if="penilaian.status === 'FINAL'" label="Download PDF SKP" color="secondary" icon="download" @click="downloadPdf" />
          </div>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useQuasar } from 'quasar'
import { api } from '@/boot/axios'
import { useAuthStore } from '@/stores/authStore'

const $q = useQuasar()
const route = useRoute()
const authStore = useAuthStore()
const penilaianId = route.params.id

const isLoading = ref(false)
const isSubmitting = ref(false)
const tab = ref('perilaku')

const penilaian = ref(null)
const hasil = ref(null)
const catatanPenilai = ref('')
const catatanBerikat = ref('')

const aspekPerilaku = [
  'Orientasi Pelayanan', 'Akuntabilitas', 'Kompetensi', 
  'Harmonisasi', 'Kolaborasi', 'Adaptif'
]

const detailPerilaku = ref(aspekPerilaku.map(a => ({ kelompok: 'PERILAKU', aspek: a, nilai: 0 })))
const detailHasilUtama = ref([{ kelompok: 'HASIL_UTAMA', aspek: '', nilai: 0 }])
const detailHasilTambahan = ref([{ kelompok: 'HASIL_TAMBAHAN', aspek: '', nilai: 0 }])

const canApprove = computed(() => {
  const role = authStore.user?.role
  return role === 'admin' || role === 'kepsek'
})

const bisaEdit = computed(() => {
  if (!penilaian.value) return false
  return penilaian.value.status === 'PENDING' || penilaian.value.status === 'RETURNED'
})

const getStatusColor = (status) => {
  switch (status) {
    case 'PENDING': return 'orange-6'
    case 'SUBMITTED': return 'blue-6'
    case 'FINAL': return 'green-6'
    case 'RETURNED': return 'red-6'
    default: return 'grey-6'
  }
}

const fetchData = async () => {
  isLoading.value = true
  try {
    const res = await api.get(`/kepegawaian/skp/penilaians/${penilaianId}`)
    penilaian.value = res.data.penilaian
    hasil.value = res.data.hasil
    catatanPenilai.value = res.data.hasil?.catatanPenilai || ''
    catatanBerikat.value = res.data.hasil?.catatanBerikat || ''

    if (res.data.detail && res.data.detail.length > 0) {
      // Kelompokkan detail dari backend
      detailPerilaku.value = aspekPerilaku.map(a => {
        const found = res.data.detail.find(d => d.kelompok === 'PERILAKU' && d.aspek === a)
        return { kelompok: 'PERILAKU', aspek: a, nilai: found ? found.nilai : 0 }
      })
      const utama = res.data.detail.filter(d => d.kelompok === 'HASIL_UTAMA')
      if (utama.length > 0) detailHasilUtama.value = utama
      
      const tambahan = res.data.detail.filter(d => d.kelompok === 'HASIL_TAMBAHAN')
      if (tambahan.length > 0) detailHasilTambahan.value = tambahan
    }
  } catch (err) {
    $q.notify({ type: 'negative', message: 'Gagal memuat form penilaian' })
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

const tambahUtama = () => detailHasilUtama.value.push({ kelompok: 'HASIL_UTAMA', aspek: '', nilai: 0 })
const tambahTambahan = () => detailHasilTambahan.value.push({ kelompok: 'HASIL_TAMBAHAN', aspek: '', nilai: 0 })

const submitPenilaian = async () => {
  isSubmitting.value = true
  try {
    const payload = {
      catatanPenilai: catatanPenilai.value,
      detail: [
        ...detailPerilaku.value,
        ...detailHasilUtama.value.filter(d => d.aspek !== ''),
        ...detailHasilTambahan.value.filter(d => d.aspek !== '')
      ]
    }
    await api.post(`/kepegawaian/skp/penilaians/${penilaianId}/submit`, payload)
    $q.notify({ type: 'positive', message: 'Penilaian berhasil disubmit' })
    await fetchData()
    tab.value = 'hasil_akhir'
  } catch (err) {
    $q.notify({ type: 'negative', message: err.response?.data?.error || 'Gagal submit' })
    console.error(err)
  } finally {
    isSubmitting.value = false
  }
}

const promptApprove = () => {
  $q.dialog({
    title: 'Setujui Penilaian',
    message: 'Masukkan catatan persetujuan:',
    prompt: { model: '', type: 'textarea' },
    cancel: true,
    persistent: true
  }).onOk(async (catatan) => {
    isSubmitting.value = true
    try {
      await api.put(`/kepegawaian/skp/penilaians/${penilaianId}/approve`, { catatan })
      $q.notify({ type: 'positive', message: 'Penilaian disetujui & Final' })
      await fetchData()
    } catch (err) {
      $q.notify({ type: 'negative', message: 'Gagal approve' })
    } finally {
      isSubmitting.value = false
    }
  })
}

const promptReturn = () => {
  $q.dialog({
    title: 'Kembalikan Penilaian',
    message: 'Masukkan alasan pengembalian:',
    prompt: { model: '', type: 'textarea' },
    cancel: true,
    persistent: true
  }).onOk(async (catatan) => {
    isSubmitting.value = true
    try {
      await api.put(`/kepegawaian/skp/penilaians/${penilaianId}/return`, { catatan })
      $q.notify({ type: 'negative', message: 'Penilaian dikembalikan untuk revisi' })
      await fetchData()
    } catch (err) {
      $q.notify({ type: 'negative', message: 'Gagal return' })
    } finally {
      isSubmitting.value = false
    }
  })
}

const downloadPdf = () => {
  $q.notify({ type: 'info', message: 'Fitur download PDF akan diarahkan ke MinIO' })
}

onMounted(fetchData)
</script>
