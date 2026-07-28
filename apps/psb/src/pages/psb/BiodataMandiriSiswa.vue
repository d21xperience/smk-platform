<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-teal-9">
          Isi Biodata Lengkap Pendaftar
        </h1>
        <p class="text-caption text-grey-7 q-mb-none">
          Silakan lengkapi formulir pendaftaran di bawah ini sesuai dengan dokumen aslinya
        </p>
      </div>
      <div>
        <q-btn
          color="teal-8"
          icon="cloud_upload"
          label="Simpan Biodata Mandiri"
          @click="onSubmit"
        />
      </div>
    </div>

    <!-- Form Konten Utama -->
    <q-form ref="formBiodata" class="q-gutter-md">
      <div class="row q-col-gutter-md">
        <!-- BAGIAN 1: BIODATA PRIBADI SISWA -->
        <div class="col-12 col-md-6">
          <q-card class="shadow-1 fit">
            <q-card-section class="bg-teal-8 text-white q-py-sm">
              <div class="text-subtitle2 text-weight-bold flex items-center">
                <q-icon name="person" class="q-mr-xs" /> I. Data Pribadi Calon Peserta Didik
              </div>
            </q-card-section>

            <q-card-section class="q-gutter-sm q-pt-md">
              <q-input
                v-model="siswa.nm_siswa"
                label="Nama Lengkap Siswa (Sesuai Ijazah SMP)"
                outlined
                dense
                :rules="[(val) => !!val || 'Nama lengkap wajib diisi']"
              />

              <div class="row q-col-gutter-sm">
                <div class="col-6">
                  <q-input
                    v-model="siswa.nisn"
                    label="NISN (10 Digit)"
                    outlined
                    dense
                    type="number"
                    :rules="[
                      (val) => !!val || 'NISN wajib diisi',
                      (val) => val.length === 10 || 'NISN nasional harus berukuran 10 digit',
                    ]"
                  />
                </div>
                <div class="col-6">
                  <q-select
                    v-model="siswa.jenis_kelamin"
                    :options="[
                      { label: 'Laki-laki', value: 'L' },
                      { label: 'Perempuan', value: 'P' },
                    ]"
                    emit-value
                    map-options
                    label="Jenis Kelamin"
                    outlined
                    dense
                    :rules="[(val) => !!val || 'Jenis kelamin wajib dipilih']"
                  />
                </div>
              </div>

              <div class="row q-col-gutter-sm">
                <div class="col-6">
                  <q-input v-model="siswa.tempat_lahir" label="Tempat Lahir" outlined dense />
                </div>
                <div class="col-6">
                  <q-input
                    v-model="siswa.tanggal_lahir"
                    type="date"
                    stack-label
                    label="Tanggal Lahir"
                    outlined
                    dense
                  />
                </div>
              </div>

              <div class="row q-col-gutter-sm">
                <div class="col-6">
                  <q-select
                    v-model="siswa.agama"
                    :options="['Islam', 'Kristen', 'Katolik', 'Hindu', 'Budha', 'Khonghucu']"
                    label="Agama"
                    outlined
                    dense
                  />
                </div>
                <div class="col-6">
                  <q-input
                    v-model="siswa.telepon_siswa"
                    label="No. Handphone Aktif (WhatsApp)"
                    outlined
                    dense
                    type="tel"
                    placeholder="Contoh: 0812..."
                  />
                </div>
              </div>

              <q-input
                v-model="siswa.alamat_siswa"
                type="textarea"
                rows="3"
                label="Alamat Rumah Tinggal Lengkap"
                outlined
                dense
                placeholder="Nama jalan, RT/RW, Desa/Kelurahan, Kecamatan, Kabupaten/Kota..."
              />

              <!-- Sistem mengunci field ini secara read-only untuk diisi oleh TU kesiswaan kelak -->
              <div class="row q-col-gutter-sm q-mt-xs">
                <div class="col-6">
                  <q-input
                    v-model="siswa.nis"
                    label="NIS Sekolah"
                    outlined
                    dense
                    readonly
                    bg-color="grey-3"
                    hint="Otomatis diisi oleh Tata Usaha setelah aktif"
                  />
                </div>
                <div class="col-6">
                  <q-input
                    v-model="siswa.diterima_tanggal"
                    label="Tanggal Masuk Resmi"
                    outlined
                    dense
                    readonly
                    bg-color="grey-3"
                    hint="Diisi oleh sistem saat kelulusan PPDB"
                  />
                </div>
              </div>
            </q-card-section>
          </q-card>
        </div>

        <!-- BAGIAN 2: DATA ORANG TUA KANDUNG & WALI -->
        <div class="col-12 col-md-6">
          <q-card class="shadow-1 fit">
            <q-card-section class="bg-teal-8 text-white q-py-sm">
              <div class="text-subtitle2 text-weight-bold flex items-center">
                <q-icon name="family_restroom" class="q-mr-xs" /> II. Data Orang Tua / Wali Kandung
              </div>
            </q-card-section>

            <q-card-section class="q-gutter-sm q-pt-md">
              <!-- Baris Identitas Ayah -->
              <div class="text-caption text-weight-bold text-teal-9">Data Ayah Kandung:</div>
              <div class="row q-col-gutter-sm">
                <div class="col-6">
                  <q-input v-model="siswa.nm_ayah" label="Nama Lengkap Ayah" outlined dense />
                </div>
                <div class="col-6">
                  <q-input
                    v-model="siswa.pekerjaan_ayah"
                    label="Pekerjaan Ayah"
                    outlined
                    dense
                    placeholder="Contoh: Buruh, Wiraswasta, PNS..."
                  />
                </div>
              </div>

              <q-separator class="q-my-sm" />

              <!-- Baris Identitas Ibu -->
              <div class="text-caption text-weight-bold text-teal-9">Data Ibu Kandung:</div>
              <div class="row q-col-gutter-sm">
                <div class="col-6">
                  <q-input v-model="siswa.nm_ibu" label="Nama Lengkap Ibu" outlined dense />
                </div>
                <div class="col-6">
                  <q-input
                    v-model="siswa.pekerjaan_ibu"
                    label="Pekerjaan Ibu"
                    outlined
                    dense
                    placeholder="Contoh: Ibu Rumah Tangga, Karyawan..."
                  />
                </div>
              </div>

              <q-separator class="q-my-md" />

              <!-- Baris Identitas Wali (Opsional) -->
              <div class="text-caption text-weight-bold text-grey-8">
                Data Wali (Isi hanya jika calon siswa tinggal bersama wali):
              </div>
              <div class="row q-col-gutter-sm">
                <div class="col-6">
                  <q-input v-model="siswa.nm_wali" label="Nama Lengkap Wali" outlined dense />
                </div>
                <div class="col-6">
                  <q-input v-model="siswa.pekerjaan_wali" label="Pekerjaan Wali" outlined dense />
                </div>
              </div>

              <div class="q-mt-lg bg-teal-1 text-teal-10 q-pa-sm rounded-borders text-caption">
                <q-icon name="shield" class="q-mr-xs" /> Pastikan nomor kontak dan e-mail yang
                terdaftar tetap aktif untuk menerima link pengumuman distribusi kelas.
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </q-form>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const formBiodata = ref(null)

// Penulisan Model Reaktif yang sama persis dengan atribut skema CREATE TABLE database Anda
const siswa = ref({
  peserta_didik_id: 'b9c231a4-9273-4f92-b062-819a3b817cba', // UUID simulasi sesi user login
  nis: '', // Kosong di fase pendaftaran awal
  nisn: '',
  nm_siswa: '',
  tempat_lahir: '',
  tanggal_lahir: '',
  jenis_kelamin: 'L',
  agama: 'Islam',
  alamat_siswa: '',
  telepon_siswa: '',
  diterima_tanggal: '', // Diisi nanti oleh sistem TU
  nm_ayah: '',
  nm_ibu: '',
  pekerjaan_ayah: '',
  pekerjaan_ibu: '',
  nm_wali: '',
  pekerjaan_wali: '',
})

const onSubmit = async () => {
  // Memicu validasi aturan isian bawaan q-form Quasar
  const isValid = await formBiodata.value.validate()

  if (!isValid) {
    $q.notify({
      color: 'negative',
      message: 'Gagal Menyimpan! Harap periksa kembali kolom isian wajib yang bertanda merah.',
      icon: 'warning',
    })
    return
  }

  // Simulasi pengiriman payload JSON murni yang siap ditembak ke server Golang Anda kelak
  console.log('Payload Database Siswa:', JSON.stringify(siswa.value))

  $q.notify({
    color: 'green-8',
    message: 'Alhamdulillah, draf biodata mandiri Anda berhasil disimpan ke cloud pendaftaran!',
    icon: 'cloud_done',
  })
}
</script>

<style scoped>
/* Pengaturan jarak h1 */
.text-h5 {
  line-height: 1.2;
}
</style>
