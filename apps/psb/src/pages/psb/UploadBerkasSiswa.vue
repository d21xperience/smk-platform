<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-teal-9">
          Upload Berkas Digital Pendukung
        </h1>
        <p class="text-caption text-grey-7 q-mb-none">
          Unggah scan dokumen asli untuk keperluan verifikasi dan komputerisasi Buku Induk oleh tim
          TU
        </p>
      </div>
      <div>
        <q-btn
          color="teal-8"
          icon="cloud_done"
          label="Simpan Seluruh Berkas"
          @click="simpanBerkasMandiri"
        />
      </div>
    </div>

    <!-- Sistem Grid Unggah Berkas -->
    <div class="row q-col-gutter-md">
      <!-- KANAN/KIRI: Pas Foto Calon Siswa (Dengan Preview Visual Interaktif) -->
      <div class="col-12 col-md-4">
        <q-card flat class="shadow-1 fit text-center q-pa-md">
          <div class="text-subtitle2 text-weight-bold text-teal-9 q-mb-xs">
            Pas Foto Resmi (3x4 / 4x6)
          </div>
          <div class="text-caption text-grey-6 q-mb-md">
            Format .jpg atau .png dengan ukuran maksimal 2MB
          </div>

          <!-- Box Preview Bingkai Foto -->
          <div
            class="flex flex-center q-mb-md bg-grey-2 rounded-borders"
            style="
              width: 100%;
              height: 180px;
              max-width: 140px;
              margin: 0 auto;
              overflow: hidden;
              border: 2px dashed #009688;
            "
          >
            <q-img v-if="fotoPreviewUrl" :src="fotoPreviewUrl" class="fit" />
            <q-icon v-else name="add_a_photo" size="48px" color="teal-3" />
          </div>

          <q-file
            v-model="berkasDigital.foto"
            label="Pilih File Foto Anda"
            outlined
            dense
            accept=".jpg, .jpeg, .png"
            max-file-size="2048000"
            @rejected="onFileRejected"
          >
            <template v-slot:prepend><q-icon name="photo_camera" color="teal-8" /></template>
          </q-file>
        </q-card>
      </div>

      <!-- KANAN: Dokumen Inti (Ijazah, KK, Akta) -->
      <div class="col-12 col-md-8">
        <q-card flat class="shadow-1 q-pa-md">
          <div class="text-subtitle2 text-weight-bold text-teal-9 q-mb-xs">
            Lampiran Dokumen Persyaratan Pelengkap
          </div>
          <div class="text-caption text-grey-6 q-mb-lg">
            Format dokumen yang diizinkan berupa PDF atau Gambar (Maksimal 5MB per file)
          </div>

          <div class="q-gutter-md">
            <!-- 1. Kolom Ijazah / SKL -->
            <div>
              <div class="text-caption text-weight-bold text-grey-8 q-mb-xs">
                1. Scan Ijazah SMP / Surat Keterangan Lulus (SKL)
              </div>
              <q-file
                v-model="berkasDigital.ijazah"
                label="Pilih file scan Ijazah / SKL"
                outlined
                dense
                accept=".pdf, .jpg, .jpeg, .png"
                max-file-size="5120000"
                @rejected="onFileRejected"
              >
                <template v-slot:prepend
                  ><q-icon name="picture_as_pdf" :color="berkasDigital.ijazah ? 'green' : 'grey-6'"
                /></template>
                <template v-slot:append v-if="berkasDigital.ijazah">
                  <q-icon
                    name="cancel"
                    @click.stop.prevent="berkasDigital.ijazah = null"
                    class="cursor-pointer"
                  />
                </template>
              </q-file>
            </div>

            <!-- 2. Kolom Kartu Keluarga -->
            <div>
              <div class="text-caption text-weight-bold text-grey-8 q-mb-xs">
                2. Scan Kartu Keluarga (KK)
              </div>
              <q-file
                v-model="berkasDigital.kk"
                label="Pilih file scan Kartu Keluarga"
                outlined
                dense
                accept=".pdf, .jpg, .jpeg, .png"
                max-file-size="5120000"
                @rejected="onFileRejected"
              >
                <template v-slot:prepend
                  ><q-icon name="badge" :color="berkasDigital.kk ? 'green' : 'grey-6'"
                /></template>
                <template v-slot:append v-if="berkasDigital.kk">
                  <q-icon
                    name="cancel"
                    @click.stop.prevent="berkasDigital.kk = null"
                    class="cursor-pointer"
                  />
                </template>
              </q-file>
            </div>

            <!-- 3. Kolom Akta Kelahiran -->
            <div>
              <div class="text-caption text-weight-bold text-grey-8 q-mb-xs">
                3. Scan Akta Kelahiran
              </div>
              <q-file
                v-model="berkasDigital.akta"
                label="Pilih file scan Akta Kelahiran"
                outlined
                dense
                accept=".pdf, .jpg, .jpeg, .png"
                max-file-size="5120000"
                @rejected="onFileRejected"
              >
                <template v-slot:prepend
                  ><q-icon name="description" :color="berkasDigital.akta ? 'green' : 'grey-6'"
                /></template>
                <template v-slot:append v-if="berkasDigital.akta">
                  <q-icon
                    name="cancel"
                    @click.stop.prevent="berkasDigital.akta = null"
                    class="cursor-pointer"
                  />
                </template>
              </q-file>
            </div>
          </div>

          <div
            class="q-mt-xl bg-amber-1 text-amber-10 q-pa-sm rounded-borders text-caption flex items-center"
          >
            <q-icon name="warning" class="q-mr-xs" />
            Pastikan hasil pindaian/scan dokumen terlihat jelas, tidak buram, dan teks dapat terbaca
            dengan mudah demi kelancaran validasi sistem.
          </div>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// State Penampung Berkas File Objek Biner
const berkasDigital = ref({
  foto: null,
  ijazah: null,
  akta: null,
  kk: null,
})

// Fungsi Reactive Komputasi Link Preview Gambar untuk Pas Foto Profil
const fotoPreviewUrl = computed(() => {
  if (!berkasDigital.value.foto) return null
  return URL.createObjectURL(berkasDigital.value.foto)
})

// Penanganan Jika File Ditolak (Melebihi batasan Max-File-Size atau Ekstensi Salah)
const onFileRejected = (rejectedEntries) => {
  rejectedEntries.forEach((entry) => {
    if (entry.failedPropValidation === 'max-file-size') {
      $q.notify({
        color: 'negative',
        message: 'Gagal memuat! Ukuran berkas terlalu besar melebihi batas ketentuan.',
        icon: 'error',
      })
    } else if (entry.failedPropValidation === 'accept') {
      $q.notify({
        color: 'negative',
        message: 'Format berkas tidak didukung! Sila gunakan PDF atau ekstensi gambar standar.',
        icon: 'warning',
      })
    }
  })
}

// Fungsi Simpan Berkas Keseluruhan di Frontend
const simpanBerkasMandiri = () => {
  // Cek minimal kelayakan draf (misal wajib mengunggah Pas Foto terlebih dahulu)
  if (!berkasDigital.value.foto) {
    $q.notify({
      color: 'orange-9',
      message: 'Perhatian: Harap unggah berkas Pas Foto resmi Anda terlebih dahulu.',
      icon: 'add_a_photo',
    })
    return
  }

  // Umpan balik kesuksesan proses upload di sisi klien
  $q.notify({
    color: 'green-8',
    message: 'Selamat! Dokumen persyaratan digital Anda telah berhasil diunggah dan disimpan.',
    icon: 'cloud_done',
  })
}
</script>

<style scoped>
.text-h5 {
  line-height: 1.2;
}
</style>
