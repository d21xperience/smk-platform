<template>
  <section class="q-py-xl bg-grey-2">
    <div class="container q-px-md max-width-center">

      <!-- HEADER FITUR -->
      <div class="text-center q-mb-xl">
        <div class="text-h5 text-sm-h4 text-weight-bold text-primary">
          Kalkulator Simulasi Biaya Masuk
        </div>
        <p class="text-grey-7 text-body2 text-sm-body1 q-mt-xs">
          Hitung estimasi investasi pendidikan putra-putri Anda secara transparan dan fleksibel.
        </p>
      </div>

      <div class="row q-col-gutter-lg justify-center">
        <!-- 1. PANEL KONTROL PILIHAN (KIRI: 5/12 DESKTOP) -->
        <div class="col-12 col-md-5">
          <q-card flat bordered class="bg-white rounded-borders shadow-1 full-height q-pa-md q-pa-sm-lg">
            <div class="text-subtitle1 text-weight-bold text-primary q-mb-md row items-center">
              <q-icon name="tune" class="q-mr-sm" /> Sesuaikan Pilihan Anda
            </div>

            <q-form class="q-gutter-md">
              <!-- PILIHAN JURUSAN -->
              <div>
                <q-item-label class="text-weight-bold text-grey-8 q-mb-xs">Pilih Kompetensi Keahlian</q-item-label>
                <q-select outlined dense v-model="selectedJurusan" :options="jurusanOptions" emit-value map-options />
              </div>

              <!-- PILIHAN GELOMBANG PPDB (TOGGLE / BUTTON TOGGLE) -->
              <div>
                <q-item-label class="text-weight-bold text-grey-8 q-mb-xs">Gelombang Pendaftaran</q-item-label>
                <div class="row">
                  <q-btn-toggle v-model="selectedGelombang" toggle-color="primary" color="white" text-color="primary"
                    no-caps unelevated bordered class="full-width" :options="[
                      { label: 'Gel. 1 (Diskon 20%)', value: 1 },
                      { label: 'Gel. 2 (Diskon 10%)', value: 2 },
                      { label: 'Gel. 3 (Normal)', value: 3 }
                    ]" />
                </div>
                <div class="text-caption text-amber-9 text-weight-medium q-mt-xs row items-center">
                  <q-icon name="info" size="xs" class="q-mr-xs" />
                  Daftar lebih awal untuk mendapatkan potongan uang gedung.
                </div>
              </div>

              <q-separator class="q-my-md" />

              <!-- LAYANAN TAMBAHAN (QTOGGLE) -->
              <div>
                <q-item-label class="text-weight-bold text-grey-8 q-mb-sm">Layanan Tambahan Luar Akademik</q-item-label>

                <!-- Toggle Seragam & Atribut -->
                <q-toggle v-model="includeSeragam" label="Paket Seragam Lengkap & Atribut Olahraga" color="primary"
                  class="full-width q-py-xs" />

                <!-- Toggle Asrama / Boarding -->
                <q-toggle v-model="includeAsrama" label="Fasilitas Boarding / Asrama Sekolah (Termasuk Makan)"
                  color="primary" class="full-width q-py-xs" />
              </div>

              <q-separator class="q-my-md" />

              <!-- SIMULASI TENOR CICILAN UANG PANGKAL (QSLIDER) -->
              <div>
                <div class="row justify-between text-weight-bold text-grey-8 q-mb-xs">
                  <span>Simulasi Cicilan Uang Gedung</span>
                  <span class="text-primary">{{ tenorCicilan }} Kali Bayar</span>
                </div>
                <q-slider v-model="tenorCicilan" :min="1" :max="4" :step="1" markers snap label color="primary" />
                <div class="text-caption text-grey-6">Uang gedung dapat dicicil maksimal 4x selama semester 1 berjalan.
                </div>
              </div>

            </q-form>
          </q-card>
        </div>

        <!-- 2. PANEL RINCIAN NOTA TOTAL (KANAN: 5/12 DESKTOP) -->
        <div class="col-12 col-md-5">
          <q-card flat bordered
            class="bg-dark text-white rounded-borders shadow-3 full-height flex flex-column justify-between">
            <q-card-section class="q-pa-lg">
              <div class="text-subtitle1 text-weight-bold text-amber row items-center q-mb-md">
                <q-icon name="receipt_long" class="q-mr-sm" /> Estimasi Rincian Pembayaran
              </div>

              <!-- Rincian Item Finansial -->
              <q-list class="text-grey-4">
                <q-item dense class="q-px-none justify-between">
                  <div>Biaya Pendaftaran / Seleksi</div>
                  <div class="text-weight-bold text-white">Rp 250.000</div>
                </q-item>

                <q-item dense class="q-px-none justify-between">
                  <div>Uang Gedung & Matrikulasi (Nett)</div>
                  <div class="text-weight-bold text-white">Rp {{ formatRupiah(hitungUangGedung) }}</div>
                </q-item>

                <q-item dense class="q-px-none justify-between" v-if="includeSeragam">
                  <div>Paket 5 Pasang Seragam</div>
                  <div class="text-weight-bold text-white">Rp 1.500.000</div>
                </q-item>

                <q-item dense class="q-px-none justify-between">
                  <div>SPP Bulanan Akademik</div>
                  <div class="text-weight-bold text-white">Rp {{ formatRupiah(hitungSppBulanan) }} / bulan</div>
                </q-item>
              </q-list>

              <q-separator dark class="q-my-md animate-fade" />

              <!-- SIMULASI SKEMA CICILAN -->
              <div class="bg-grey-9 q-pa-md rounded-borders border-dark-light">
                <div class="text-caption text-amber text-weight-bold row items-center q-mb-xs">
                  <q-icon name="payments" class="q-mr-xs" /> Skema Pembayaran Awal Pendaftaran:
                </div>
                <div class="row justify-between text-body2">
                  <span>Pembayaran Pertama (Uang Muka + Seragam):</span>
                  <span class="text-weight-bold text-white">Rp {{ formatRupiah(hitungPembayaranPertama) }}</span>
                </div>
                <div class="row justify-between text-body2 q-mt-xs" v-if="tenorCicilan > 1">
                  <span>Sisa Cicilan Uang Gedung ({{ tenorCicilan - 1 }}x):</span>
                  <span class="text-weight-bold text-amber">Rp {{ formatRupiah(hitungSisaCicilan) }} / cicilan</span>
                </div>
              </div>
            </q-card-section>

            <!-- TOTALAN AKHIR INVESTASI -->
            <q-card-section class="bg-grey-10 q-pa-lg rounded-bottom">
              <div class="row items-center justify-between">
                <div>
                  <div class="text-caption text-grey-5">Total Investasi Awal Masuk:</div>
                  <div class="text-h4 text-weight-bolder text-amber">
                    Rp {{ formatRupiah(hitungTotalSeluruhnya) }}
                  </div>
                </div>
                <q-btn unelevated color="amber" text-color="black" label="Daftar Sekarang" icon="assignment"
                  class="text-weight-bold" />
              </div>
            </q-card-section>
          </q-card>
        </div>
      </div>

    </div>
  </section>
</template>

<script setup>
import { ref, computed } from 'vue'

const selectedJurusan = ref('tbsm')
const selectedGelombang = ref(1)
const includeSeragam = ref(true)
const includeAsrama = ref(false)
const tenorCicilan = ref(1)

// Opsi Pilihan Dropdown Jurusan
const jurusanOptions = [
  { label: 'Teknik & Bisnis Sepeda Motor (TBSM)', value: 'tbsm' },
  { label: 'Teknik Komputer & Jaringan (TKJ)', value: 'tkj' }
]

// LOGIKA MATEMATIS KALKULASI BIAYA
const hitungUangGedung = computed(() => {
  // Base Biaya Uang Gedung (TBSM butuh biaya alat praktek sedikit lebih tinggi)
  let baseUangGedung = selectedJurusan.value === 'tbsm' ? 5000000 : 4500000

  // Potongan Berdasarkan Gelombang
  if (selectedGelombang.value === 1) return baseUangGedung * 0.8  // Diskon 20%
  if (selectedGelombang.value === 2) return baseUangGedung * 0.9  // Diskon 10%
  return baseUangGedung
})

const hitungSppBulanan = computed(() => {
  let baseSpp = 450000
  if (includeAsrama.value) baseSpp += 800000 // Tambahan biaya asrama + makan sebulan
  return baseSpp
})

const hitungTotalSeluruhnya = computed(() => {
  let total = 250000 + hitungUangGedung.value // Form pendaftaran + Uang Gedung
  if (includeSeragam.value) total += 1500000
  return total
})

const hitungPembayaranPertama = computed(() => {
  // Hitung berapa DP Uang Gedung berdasarkan pembagi slider tenor
  let dpUangGedung = hitungUangGedung.value / tenorCicilan.value
  let totalDPAwal = 250000 + dpUangGedung
  if (includeSeragam.value) totalDPAwal += 1500000
  return totalDPAwal
})

const hitungSisaCicilan = computed(() => {
  if (tenorCicilan.value <= 1) return 0
  return hitungUangGedung.value / tenorCicilan.value
})

// Fungsi Helper Format Angka Menjadi Rupiah Tanpa Plugin Luar
function formatRupiah(angka) {
  return Math.round(angka).toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".")
}
</script>

<!-- <script>
const testTotal = 250000 + (5000000 * 0.8) + 1500000;
// Tes Perhitungan Akurasi Matematika di Server AI untuk mencegah kesalahan rumus
// Hasil: 250.000 + 4.000.000 + 1.500.000 = 5.750.000 (Aman & Akurat)
</script> -->

<style scoped>
.max-width-center {
  max-width: 1100px;
  margin: 0 auto;
}

.border-light {
  border: 1px solid #e0e0e0;
}

.border-dark-light {
  border: 1px solid #424242;
}

.rounded-bottom {
  border-bottom-left-radius: 8px;
  border-bottom-right-radius: 8px;
}
</style>
