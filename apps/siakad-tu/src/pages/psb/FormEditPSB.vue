<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Formulir -->
    <div class="row items-center justify-between q-mb-md">
      <div class="row items-center q-gutter-sm">
        <q-btn icon="arrow_back" color="grey-8" flat round @click="kembaliKeDaftar" />
        <div>
          <h1 class="text-h5 text-weight-bold q-my-none text-primary">Formulir Detail Calon Siswa</h1>
          <p class="text-caption text-grey-7 q-mb-none">Lengkapi berkas komputerisasi Buku Induk berdasarkan instrumen
            registrasi PPDB</p>
        </div>
      </div>
      <div>
        <q-btn color="primary" icon="save" label="Simpan Perubahan" @click="saveDetailSiswa" />
      </div>
    </div>

    <!-- Konten Formulir Berdasarkan Struktur Database -->
    <q-form class="q-gutter-md">
      <div class="row q-col-gutter-md">

        <!-- BLOCK 1: IDENTITAS POKOK SISWA -->
        <div class="col-12 col-md-6">
          <q-card class="shadow-1 h-100">
            <q-card-section class="bg-indigo-9 text-white q-py-sm">
              <div class="text-subtitle2 text-weight-bold">I. Identitas Pribadi Peserta Didik</div>
            </q-card-section>

            <q-card-section class="q-gutter-sm q-pt-md">
              <q-input v-model="siswa.nm_siswa" label="Nama Lengkap Siswa" outlined dense lazy-rules
                :rules="[val => !!val || 'Nama tidak boleh kosong']" />
              <q-input v-model="siswa.nis" label="Nomor Induk Siswa (NIS)" outlined dense readonly
                hint="NIS otomatis digenerate lewat menu Kelulusan Massal" />

              <div class="row q-col-gutter-sm">
                <div class="col-6">
                  <q-input v-model="siswa.nisn" label="NISN Nasional" outlined dense type="number" />
                </div>
                <div class="col-6">
                  <q-select v-model="siswa.jenis_kelamin"
                    :options="[{ label: 'Laki-laki', value: 'L' }, { label: 'Perempuan', value: 'P' }]" emit-value
                    map-options label="Jenis Kelamin" outlined dense />
                </div>
              </div>

              <div class="row q-col-gutter-sm">
                <div class="col-6">
                  <q-input v-model="siswa.tempat_lahir" label="Tempat Lahir" outlined dense />
                </div>
                <div class="col-6">
                  <q-input v-model="siswa.tanggal_lahir" type="date" stack-label label="Tanggal Lahir" outlined dense />
                </div>
              </div>

              <div class="row q-col-gutter-sm">
                <div class="col-6">
                  <q-select v-model="siswa.agama"
                    :options="['Islam', 'Kristen', 'Katolik', 'Hindu', 'Budha', 'Khonghucu']" label="Agama" outlined
                    dense />
                </div>
                <div class="col-6">
                  <q-input v-model="siswa.telepon_siswa" label="No. Telepon / HP Siswa" outlined dense type="tel" />
                </div>
              </div>

              <q-input v-model="siswa.alamat_siswa" type="textarea" rows="2" label="Alamat Tempat Tinggal Lengkap"
                outlined dense />
              <q-input v-model="siswa.diterima_tanggal" type="date" stack-label label="Tanggal Diterima di Sekolah"
                outlined dense />
            </q-card-section>
          </q-card>
        </div>

        <!-- BLOCK 2: DATA ORANG TUA / WALI -->
        <div class="col-12 col-md-6">
          <q-card class="shadow-1 h-100">
            <q-card-section class="bg-indigo-9 text-white q-py-sm">
              <div class="text-subtitle2 text-weight-bold">II. Data Orang Tua & Wali Kandung</div>
            </q-card-section>

            <q-card-section class="q-gutter-sm q-pt-md">
              <!-- Ayah -->
              <div class="row q-col-gutter-sm">
                <div class="col-6">
                  <q-input v-model="siswa.nm_ayah" label="Nama Lengkap Ayah" outlined dense />
                </div>
                <div class="col-6">
                  <q-input v-model="siswa.pekerjaan_ayah" label="Pekerjaan Ayah" outlined dense />
                </div>
              </div>

              <!-- Ibu -->
              <div class="row q-col-gutter-sm q-mt-xs">
                <div class="col-6">
                  <q-input v-model="siswa.nm_ibu" label="Nama Lengkap Ibu" outlined dense />
                </div>
                <div class="col-6">
                  <q-input v-model="siswa.pekerjaan_ibu" label="Pekerjaan Ibu" outlined dense />
                </div>
              </div>

              <q-separator class="q-my-md" />

              <!-- Wali (Opsional) -->
              <div class="text-caption text-grey-7 text-weight-bold q-mb-xs">Kontak Wali (Jika tidak tinggal bersama
                orang tua):</div>
              <div class="row q-col-gutter-sm">
                <div class="col-6">
                  <q-input v-model="siswa.nm_wali" label="Nama Lengkap Wali" outlined dense />
                </div>
                <div class="col-6">
                  <q-input v-model="siswa.pekerjaan_wali" label="Pekerjaan Wali" outlined dense />
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>
        <!-- ====== TAMBAHKAN BLOK BARU INI (KOLOM UNGGAH DOKUMEN DIGITAL) ====== -->
        <div class="col-12">
          <q-card class="shadow-1">
            <q-card-section class="bg-indigo-9 text-white q-py-sm">
              <div class="text-subtitle2 text-weight-bold">III. Unggah Dokumen & Berkas Digital Pendukung</div>
            </q-card-section>

            <q-card-section class="q-pa-md">
              <div class="row q-col-gutter-md items-center">

                <!-- Sisi Kiri: Pas Foto Resmi Calon Siswa (Dengan Preview) -->
                <div class="col-12 col-sm-4 text-center">
                  <div class="text-caption text-weight-bold text-grey-8 q-mb-xs">Pas Foto Siswa (3x4 / 4x6)</div>

                  <div class="flex flex-center q-mb-sm bg-grey-3 rounded-borders"
                    style="width: 100%; height: 160px; max-width: 130px; margin: 0 auto; overflow: hidden; border: 2px dashed #9e9e9e;">
                    <q-img v-if="fotoPreviewUrl" :src="fotoPreviewUrl" class="fit" />
                    <q-icon v-else name="account_box" size="56px" color="grey-6" />
                  </div>

                  <q-file v-model="berkasDigital.foto" label="Pilih Foto (.jpg, .png)" outlined dense
                    accept=".jpg, .jpeg, .png" max-file-size="2048000" class="q-mt-xs">
                    <template v-slot:prepend><q-icon name="cloud_upload" /></template>
                  </q-file>
                </div>

                <!-- Sisi Kanan: Dokumen Persyaratan Fisik (Ijazah, Akta, KK) -->
                <div class="col-12 col-sm-8 q-gutter-sm">
                  <!-- Unggah Ijazah / SKL -->
                  <div>
                    <div class="text-caption text-weight-bold text-grey-8 q-mb-xs">Saluran Berkas Ijazah / SKL SMP</div>
                    <q-file v-model="berkasDigital.ijazah" label="Unggah Ijazah / SKL Resmi (PDF/Gambar)" outlined dense
                      accept=".pdf, .jpg, .jpeg, .png" max-file-size="5120000">
                      <template v-slot:prepend><q-icon name="picture_as_pdf"
                          :color="berkasDigital.ijazah ? 'green' : 'grey-7'" /></template>
                      <template v-slot:append v-if="berkasDigital.ijazah">
                        <q-icon name="cancel" @click.stop.prevent="berkasDigital.ijazah = null"
                          class="cursor-pointer" />
                      </template>
                    </q-file>
                  </div>

                  <!-- Unggah Akta Kelahiran -->
                  <div class="q-mt-sm">
                    <div class="text-caption text-weight-bold text-grey-8 q-mb-xs">Saluran Berkas Akta Kelahiran</div>
                    <q-file v-model="berkasDigital.akta" label="Unggah Akta Kelahiran Resmi (PDF/Gambar)" outlined dense
                      accept=".pdf, .jpg, .jpeg, .png" max-file-size="5120000">
                      <template v-slot:prepend><q-icon name="description"
                          :color="berkasDigital.akta ? 'green' : 'grey-7'" /></template>
                      <template v-slot:append v-if="berkasDigital.akta">
                        <q-icon name="cancel" @click.stop.prevent="berkasDigital.akta = null" class="cursor-pointer" />
                      </template>
                    </q-file>
                  </div>

                  <!-- Unggah Kartu Keluarga -->
                  <div class="q-mt-sm">
                    <div class="text-caption text-weight-bold text-grey-8 q-mb-xs">Saluran Berkas Kartu Keluarga (KK)
                    </div>
                    <q-file v-model="berkasDigital.kk" label="Unggah Kartu Keluarga Legalitas (PDF/Gambar)" outlined
                      dense accept=".pdf, .jpg, .jpeg, .png" max-file-size="5120000">
                      <template v-slot:prepend><q-icon name="badge"
                          :color="berkasDigital.kk ? 'green' : 'grey-7'" /></template>
                      <template v-slot:append v-if="berkasDigital.kk">
                        <q-icon name="cancel" @click.stop.prevent="berkasDigital.kk = null" class="cursor-pointer" />
                      </template>
                    </q-file>
                  </div>

                </div>

              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </q-form>
  </q-page>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuasar } from 'quasar'

const route = useRoute()
const router = useRouter()
const $q = useQuasar()

// Inisialisasi skema data objek reaktif kosong agar sesuai dengan struktur kolom database
const siswa = ref({
  peserta_didik_id: '',
  nis: '',
  nisn: '',
  nm_siswa: '',
  tempat_lahir: '',
  tanggal_lahir: '',
  jenis_kelamin: 'L',
  agama: 'Islam',
  alamat_siswa: '',
  telepon_siswa: '',
  diterima_tanggal: '',
  nm_ayah: '',
  nm_ibu: '',
  pekerjaan_ayah: '',
  pekerjaan_ibu: '',
  nm_wali: '',
  pekerjaan_wali: ''
})

onMounted(() => {
  // Mengambil ID parameter siswa dari URL rute (Misal: /siswa/ppdb/edit/1)
  const idSiswa = route.params.id

  // Karena saat ini masih fokus frontend, kita simulasikan penarikan data dummy berdasarkan id
  if (idSiswa) {
    siswa.value.peserta_didik_id = 'b9c231a4-9273-4f92-b062-819a3b817cba' // Contoh Format UUID
    siswa.value.nm_siswa = idSiswa === '1' ? 'Bagus Setiawan' : 'Aditya Pratama'
    siswa.value.nisn = idSiswa === '1' ? '0091122331' : '0091122332'
    siswa.value.nis = idSiswa === '1' ? '26270001' : ''
    siswa.value.telepon_siswa = '08123456789'
    siswa.value.nm_ayah = 'Heri Setiawan'
    siswa.value.nm_ibu = 'Siti Rahma'
    siswa.value.pekerjaan_ayah = 'Wiraswasta'
    siswa.value.pekerjaan_ibu = 'Ibu Rumah Tangga'
  }
})

const saveDetailSiswa = () => {
  if (!siswa.value.nm_siswa) {
    $q.notify({ color: 'negative', message: 'Kolom Nama Lengkap wajib terisi!', icon: 'warning' })
    return
  }

  // Umpan balik keberhasilan pengisian data
  $q.notify({
    color: 'green-8',
    message: `Arsip data kesiswaan ${siswa.value.nm_siswa} berhasil disimpan ke database lokal!`,
    icon: 'cloud_done'
  })

  kembaliKeDaftar()
}

const kembaliKeDaftar = () => {
  router.push({ name: 'tambah-siswa' })
}

// ====== TAMBAHKAN STATE BERKAS DIGITAL INI ======
const berkasDigital = ref({
  foto: null,
  ijazah: null,
  akta: null,
  kk: null
})

// Fungsi bantu untuk membuat link preview gambar instan di frontend (Khusus Pas Foto)
const fotoPreviewUrl = computed(() => {
  if (!berkasDigital.value.foto) return null
  return URL.createObjectURL(berkasDigital.value.foto)
})
// ================================================
</script>
