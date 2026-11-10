<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row items-center q-mb-md">
      <q-btn flat round dense icon="arrow_back" @click="$router.push('/kepegawaian/buku-induk')" class="q-mr-sm" />
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Detail Buku Induk</h1>
        <p class="text-caption text-grey-7 q-mb-none">Informasi lengkap data pegawai</p>
      </div>
      <q-space />
      <q-btn color="primary" icon="print" label="Cetak PDF" @click="cetakPDF" />
    </div>

    <q-card v-if="isLoading" class="q-pa-lg text-center">
      <q-spinner color="primary" size="3em" />
      <div class="q-mt-md">Memuat data...</div>
    </q-card>

    <q-card v-else-if="detail" class="shadow-1">
      <q-card-section class="bg-primary text-white row items-center">
        <q-avatar size="80px" class="q-mr-md">
          <img :src="detail.fotoUrl || 'https://cdn.quasar.dev/img/boy-avatar.png'">
        </q-avatar>
        <div>
          <div class="text-h5 text-weight-bold">{{ detail.name }}</div>
          <div class="text-subtitle2">NIK: {{ detail.nik || '-' }} | NIP: {{ detail.nip || '-' }}</div>
        </div>
        <q-space />
        <div class="text-center">
          <q-circular-progress show-value :value="detail.skorKelengkapan" size="60px"
            :color="getBadgeColor(detail.statusKelengkapan)" track-color="grey-8" class="q-ma-md">
            {{ detail.skorKelengkapan }}%
          </q-circular-progress>
          <div class="text-caption">{{ detail.statusKelengkapan }}</div>
        </div>
      </q-card-section>

      <q-tabs v-model="tab" class="text-primary" align="left">
        <q-tab name="identitas" icon="person" label="Identitas" />
        <q-tab name="kepegawaian" icon="work" label="Kepegawaian" />
        <q-tab name="pendidikan" icon="school" label="Pendidikan" />
        <q-tab name="sertifikasi" icon="workspace_premium" label="Sertifikasi" />
        <q-tab name="dokumen" icon="folder" label="Dokumen" />
      </q-tabs>
      <q-separator />

      <q-tab-panels v-model="tab" animated class="bg-white">
        <q-tab-panel name="identitas">
          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6"><q-input v-model="detail.nik" label="NIK" outlined readonly /></div>
            <div class="col-12 col-md-6"><q-input v-model="detail.nuptk" label="NUPTK" outlined readonly /></div>
            <div class="col-12 col-md-6"><q-input v-model="detail.gender" label="Jenis Kelamin" outlined readonly />
            </div>
            <div class="col-12 col-md-6"><q-input v-model="detail.birthPlace" label="Tempat Lahir" outlined readonly />
            </div>
            <div class="col-12 col-md-6"><q-input v-model="tglLahir" label="Tanggal Lahir" outlined readonly /></div>
            <div class="col-12 col-md-6"><q-input v-model="detail.noHp" label="No HP" outlined readonly /></div>
            <div class="col-12 col-md-6"><q-input v-model="detail.email" label="Email" outlined readonly /></div>
            <div class="col-12 col-md-6"><q-input v-model="detail.address" label="Alamat" outlined readonly /></div>
          </div>
        </q-tab-panel>

        <q-tab-panel name="kepegawaian">
          <div class="row q-col-gutter-md">
            <div class="col-12 col-md-6"><q-input v-model="detail.jenisPegawai" label="Jenis Pegawai" outlined
                readonly /></div>
            <div class="col-12 col-md-6"><q-input v-model="detail.statusPegawai" label="Status Pegawai" outlined
                readonly /></div>
            <div class="col-12 col-md-6"><q-input v-model="detail.jabatan" label="Jabatan" outlined readonly /></div>
            <div class="col-12 col-md-6"><q-input v-model="detail.unitKerja" label="Unit Kerja" outlined readonly />
            </div>
            <div class="col-12 col-md-6"><q-input v-model="tmtMasuk" label="TMT Masuk" outlined readonly /></div>
          </div>
        </q-tab-panel>

        <q-tab-panel name="pendidikan">
          <div class="text-grey-7 text-center q-pa-md">Data riwayat pendidikan akan ditampilkan di sini.</div>
        </q-tab-panel>

        <q-tab-panel name="sertifikasi">
          <div class="text-grey-7 text-center q-pa-md">Data sertifikasi akan ditampilkan di sini.</div>
        </q-tab-panel>

        <q-tab-panel name="dokumen">
          <div class="text-grey-7 text-center q-pa-md">Data dokumen arsip akan ditampilkan di sini.</div>
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

const $q = useQuasar()
const route = useRoute()
const detail = ref(null)
const isLoading = ref(false)
const tab = ref('identitas')

const getBadgeColor = (status) => {
  if (status === 'Lengkap') return 'green-6'
  if (status === 'Cukup') return 'orange-6'
  return 'red-6'
}

const formatDate = (val) => {
  if (!val) return '-'
  return new Date(val).toLocaleDateString('id-ID')
}

const tglLahir = computed(() => formatDate(detail.value?.birthDate))
const tmtMasuk = computed(() => formatDate(detail.value?.tmtMasuk))

const fetchDetail = async () => {
  isLoading.value = true
  try {
    const res = await api.get(`/teachers/buku-induk/${route.params.id}`)
    detail.value = res.data
  } catch (err) {
    $q.notify({ type: 'negative', message: `Gagal memuat detail pegawai ${err}` })
  } finally {
    isLoading.value = false
  }
}

const cetakPDF = () => {
  $q.notify({ type: 'info', message: 'Fitur Cetak PDF akan segera tersedia.' })
}

onMounted(fetchDetail)
</script>
eof
