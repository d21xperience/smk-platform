<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Ucapan Selamat Datang -->
    <div class="row q-col-gutter-sm items-center q-mb-md">
      <div class="col-12">
        <h1 class="text-h5 text-weight-bold q-my-none text-teal-9">
          Selamat Datang, {{ calonSiswa.nama }}!
        </h1>
        <p class="text-caption text-grey-7 q-mb-none">
          Berikut adalah rangkuman status pendaftaran Anda di SMK Swasta Indonesia
        </p>
      </div>
    </div>

    <!-- 1. Banner Alur Informasi Utama PPDB -->
    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12">
        <q-banner inline-actions class="text-white bg-teal-8 rounded-borders shadow-1">
          <template v-slot:avatar>
            <q-icon name="campaign" size="md" />
          </template>
          <div class="text-subtitle2 text-weight-bold">Informasi Panitia PPDB:</div>
          <div class="text-caption opacity-80">
            Harap segera melengkapi pengisian biodata mandiri dan mengunggah berkas digital.
            Verifikasi berkas fisik akan dilakukan di ruang Tata Usaha setiap hari kerja (08:00 -
            14:00 WIB).
          </div>
        </q-banner>
      </div>
    </div>

    <!-- 2. Row Kartu Informasi Status Pendaftaran -->
    <div class="row q-col-gutter-md q-mb-lg">
      <!-- Status Kelulusan Seleksi -->
      <div class="col-12 col-md-4">
        <q-card flat class="shadow-1 bg-white fit">
          <q-card-section class="q-py-md text-center">
            <q-avatar
              color="teal-1"
              text-color="teal-9"
              icon="assignment_ind"
              size="56px"
              class="q-mb-sm"
            />
            <div class="text-caption text-grey-7 text-uppercase text-weight-bold">
              No. Urut Pendaftaran
            </div>
            <div class="text-h5 text-weight-bold text-teal-9 q-my-xs">
              REG-2627-{{ calonSiswa.id.toString().padStart(3, '0') }}
            </div>

            <q-chip
              :color="calonSiswa.statusSeleksi === 'Lolos Seleksi' ? 'green-1' : 'orange-1'"
              :text-color="calonSiswa.statusSeleksi === 'Lolos Seleksi' ? 'green-9' : 'orange-9'"
              icon="verified"
              class="text-weight-bold q-mt-sm"
            >
              {{ calonSiswa.statusSeleksi }}
            </q-chip>
          </q-card-section>
        </q-card>
      </div>

      <!-- Status Kelengkapan Berkas Fisik -->
      <div class="col-12 col-md-4">
        <q-card flat class="shadow-1 bg-white fit">
          <q-card-section class="q-py-md text-center">
            <q-avatar
              color="indigo-1"
              text-color="indigo-9"
              icon="fact_check"
              size="56px"
              class="q-mb-sm"
            />
            <div class="text-caption text-grey-7 text-uppercase text-weight-bold">
              Kelengkapan Berkas Fisik
            </div>
            <div class="text-h5 text-weight-bold text-grey-9 q-my-xs">
              {{ hitungBerkasLengkap }} / 3 <span class="text-caption text-grey-6">Dokumen</span>
            </div>

            <q-linear-progress
              rounded
              size="8px"
              :value="hitungBerkasLengkap / 3"
              color="indigo-9"
              class="q-mt-md q-mx-auto"
              style="max-width: 180px"
            />
            <div class="text-caption text-grey-6 q-mt-xs">Verifikasi Ruang TU</div>
          </q-card-section>
        </q-card>
      </div>

      <!-- Status Keuangan / Angsuran -->
      <div class="col-12 col-md-4">
        <q-card flat class="shadow-1 bg-white fit">
          <q-card-section class="q-py-md text-center">
            <q-avatar
              :color="calonSiswa.bayar === standarBiaya ? 'green-1' : 'red-1'"
              :text-color="calonSiswa.bayar === standarBiaya ? 'green-9' : 'red-9'"
              icon="payments"
              size="56px"
              class="q-mb-sm"
            />
            <div class="text-caption text-grey-7 text-uppercase text-weight-bold">
              Biaya Pendaftaran (Rp 450k)
            </div>
            <div
              class="text-h5 text-weight-bold q-my-xs"
              :class="calonSiswa.bayar === standarBiaya ? 'text-green-8' : 'text-orange-9'"
            >
              Rp {{ calonSiswa.bayar.toLocaleString('id-ID') }}
            </div>

            <q-badge
              :color="calonSiswa.bayar === standarBiaya ? 'green' : 'orange-9'"
              :label="
                calonSiswa.bayar === standarBiaya ? 'Lunas (Tunai)' : 'Angsuran / Belum Lunas'
              "
              class="q-pa-xs q-mt-sm"
            />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- 3. Detail Checklist Dokumen Fisik Persyaratan -->
    <div class="row q-col-gutter-md">
      <!-- Panel Kiri: Ceklis Dokumen Terverifikasi TU -->
      <div class="col-12 col-md-6">
        <q-card flat class="shadow-1">
          <q-card-section class="bg-grey-1 q-py-sm">
            <div class="text-subtitle2 text-weight-bold text-grey-9">
              Checklist Dokumen Persyaratan Fisik
            </div>
            <div class="text-caption text-grey-6">
              Status validasi berkas fisik asli yang diserahkan ke panitia
            </div>
          </q-card-section>

          <q-card-section class="q-pa-none">
            <q-list separator>
              <!-- Ijazah -->
              <q-item class="q-py-md">
                <q-item-section avatar>
                  <q-icon
                    :name="calonSiswa.berkas.ijazah ? 'check_circle' : 'cancel'"
                    :color="calonSiswa.berkas.ijazah ? 'green' : 'red'"
                    size="md"
                  />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium"
                    >Fotokopi Ijazah / SKL SMP Terlegalisir</q-item-label
                  >
                  <q-item-label caption>{{
                    calonSiswa.berkas.ijazah
                      ? 'Diterima & Sesuai data nasional'
                      : 'Belum diserahkan ke panitia'
                  }}</q-item-label>
                </q-item-section>
              </q-item>

              <!-- Kartu Keluarga -->
              <q-item class="q-py-md">
                <q-item-section avatar>
                  <q-icon
                    :name="calonSiswa.berkas.kk ? 'check_circle' : 'cancel'"
                    :color="calonSiswa.berkas.kk ? 'green' : 'red'"
                    size="md"
                  />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium"
                    >Fotokopi Kartu Keluarga (KK)</q-item-label
                  >
                  <q-item-label caption>{{
                    calonSiswa.berkas.kk
                      ? 'Diterima & Data orang tua valid'
                      : 'Belum diserahkan ke panitia'
                  }}</q-item-label>
                </q-item-section>
              </q-item>

              <!-- Akta Kelahiran -->
              <q-item class="q-py-md">
                <q-item-section avatar>
                  <q-icon
                    :name="calonSiswa.berkas.akta ? 'check_circle' : 'cancel'"
                    :color="calonSiswa.berkas.akta ? 'green' : 'red'"
                    size="md"
                  />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">Fotokopi Akta Kelahiran</q-item-label>
                  <q-item-label caption>{{
                    calonSiswa.berkas.akta
                      ? 'Diterima & Tanggal lahir sesuai'
                      : 'Belum diserahkan ke panitia'
                  }}</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>

      <!-- Panel Kanan: Rincian Keuangan & Panduan Langkah Selanjutnya -->
      <div class="col-12 col-md-6">
        <q-card flat class="shadow-1 bg-white">
          <q-card-section class="bg-grey-1 q-py-sm">
            <div class="text-subtitle2 text-weight-bold text-grey-9">
              Informasi Tunggakan Keuangan PPDB
            </div>
          </q-card-section>

          <q-card-section class="q-gutter-sm">
            <div class="row justify-between items-center bg-grey-1 q-pa-md rounded-borders">
              <span class="text-body2 text-grey-7">Total Biaya Pendaftaran:</span>
              <span class="text-subtitle1 text-weight-bold"
                >Rp {{ standarBiaya.toLocaleString('id-ID') }}</span
              >
            </div>

            <div
              class="row justify-between items-center bg-teal-1 text-teal-10 q-pa-md rounded-borders"
            >
              <span class="text-body2">Total yang Sudah Dibayar:</span>
              <span class="text-subtitle1 text-weight-bold"
                >Rp {{ calonSiswa.bayar.toLocaleString('id-ID') }}</span
              >
            </div>

            <div
              class="row justify-between items-center bg-red-1 text-red-9 q-pa-md rounded-borders"
            >
              <span class="text-body2">Sisa Kekurangan Angsuran:</span>
              <span class="text-subtitle1 text-weight-bold"
                >Rp {{ sisaTunggakan.toLocaleString('id-ID') }}</span
              >
            </div>

            <div class="text-caption text-grey-6 q-mt-md flex items-center">
              <q-icon name="info" class="q-mr-xs" color="orange-9" />
              *Akun Anda baru dapat diaktifkan ke menu utama sekolah jika kekurangan angsuran di
              atas bernilai Rp 0 (Lunas).
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'

const standarBiaya = 450000

// Data Mock / Dummy Sesi Calon Siswa yang sedang Login Mandiri
// Data ini merepresentasikan satu pendaftar (contoh diambil kondisi mengangsur cicilan)
const calonSiswa = ref({
  id: 2,
  nama: 'Aditya Pratama',
  jurusan: 'TKRO',
  statusSeleksi: 'Karantina PPDB', // Status awal di sistem panitia
  bayar: 200000, // Siswa mencicil Rp 200.000 (belum lunas)
  berkas: {
    ijazah: true,
    kk: true,
    akta: false, // Akta belum dikumpulkan ke sekolah
  },
})

// Perhitungan Sisa Tunggakan Keuangan secara otomatis di Frontend
const sisaTunggakan = computed(() => {
  return standarBiaya - calonSiswa.value.bayar
})

// Menghitung jumlah berkas fisik yang sudah divalidasi oleh panitia TU kesiswaan
const hitungBerkasLengkap = computed(() => {
  let count = 0
  if (calonSiswa.value.berkas.ijazah) count++
  if (calonSiswa.value.berkas.kk) count++
  if (calonSiswa.value.berkas.akta) count++
  return count
})
</script>

<style scoped>
.opacity-80 {
  opacity: 0.85;
}
</style>
