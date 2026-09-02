<template>
  <q-page class="q-pa-md">
    <!-- BARIS 1: METRIK UTAMA -->
    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-sm-6 col-md-3">
        <q-card flat bordered class="bg-blue-1 text-blue-9">
          <q-card-section class="row items-center no-wrap">
            <div>
              <div class="text-subtitle2 text-grey-7 text-uppercase">Nilai Aset Yayasan</div>
              <div class="text-h5 text-weight-bold">Rp 1.2M</div>
            </div>
            <q-space />
            <q-icon name="account_balance_wallet" size="lg" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card flat bordered class="bg-green-1 text-green-9">
          <q-card-section class="row items-center no-wrap">
            <div>
              <div class="text-subtitle2 text-grey-7 text-uppercase">Aset Kondisi Baik</div>
              <div class="text-h5 text-weight-bold">1,420 <span class="text-caption">Unit</span></div>
            </div>
            <q-space />
            <q-icon name="check_circle" size="lg" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card flat bordered class="bg-amber-1 text-amber-9">
          <q-card-section class="row items-center no-wrap">
            <div>
              <div class="text-subtitle2 text-grey-7 text-uppercase">Butuh Servis</div>
              <div class="text-h5 text-weight-bold">12 <span class="text-caption">Barang</span></div>
            </div>
            <q-space />
            <q-icon name="handyman" size="lg" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-6 col-md-3">
        <q-card flat bordered class="bg-red-1 text-red-9">
          <q-card-section class="row items-center no-wrap">
            <div>
              <div class="text-subtitle2 text-grey-7 text-uppercase">Stok ATK Kritis</div>
              <div class="text-h5 text-weight-bold">4 <span class="text-caption">Item</span></div>
            </div>
            <q-space />
            <q-icon name="warning" size="lg" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- BARIS 2: KONTEN UTAMA DAN TIMELINE -->
    <div class="row q-col-gutter-md">
      <!-- Sisi Kiri: Tiket Kerusakan & Logistik Kritis -->
      <div class="col-12 col-md-8">
        <!-- Komponen Tiket Laporan Guru -->
        <q-card flat bordered class="q-mb-md">
          <q-card-section class="row items-center q-pb-none">
            <div class="text-h6 text-weight-bold text-grey-8">Tiket Laporan Kerusakan Terkini</div>
            <q-space />
            <q-btn flat color="primary" label="Lihat Semua" to="/sarpras/servis/tiket" />
          </q-card-section>

          <q-card-section>
            <q-list separator>
              <q-item v-for="tiket in tiketKerusakan" :key="tiket.id" class="q-py-md">
                <q-item-section avatar>
                  <q-avatar :icon="tiket.icon" :color="tiket.urgensiColor" text-color="white" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-weight-medium">{{ tiket.namaBarang }}</q-item-label>
                  <q-item-label caption>Lokasi: {{ tiket.lokasi }} | Pelapor: {{ tiket.pelapor }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-badge :color="tiket.statusColor">{{ tiket.status }}</q-badge>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>

      <!-- Sisi Kanan: Jadwal Kegiatan & Perawatan -->
      <div class="col-12 col-md-4">
        <q-card flat bordered class="fit">
          <q-card-section>
            <div class="text-h6 text-weight-bold text-grey-8 q-mb-md">Agenda Sarpras Hari Ini</div>
            <q-timeline color="secondary">
              <q-timeline-entry v-for="agenda in agendaHariIni" :key="agenda.id" :title="agenda.judul"
                :subtitle="agenda.jam" :icon="agenda.icon" :color="agenda.color">
                <div>{{ agenda.deskripsi }}</div>
              </q-timeline-entry>
            </q-timeline>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'

// Mock Data Tiket Laporan Kerusakan
const tiketKerusakan = ref([
  { id: 1, namaBarang: 'AC Ruang Kelas 11-B Bocor Air', lokasi: 'Gedung A, Lt 2', pelapor: 'Bu Sri (Wali Kelas)', status: 'Pending', icon: 'ac_unit', urgensiColor: 'red', statusColor: 'grey-7' },
  { id: 2, namaBarang: 'Proyektor BenQ Mati Total', lokasi: 'Lab Fisika', pelapor: 'Pak Bambang', status: 'Proses Servis', icon: 'videocam', urgensiColor: 'amber', statusColor: 'orange' },
  { id: 3, namaBarang: 'Engsel Pintu Kamar Mandi Lepas', lokasi: 'Toilet Siswa Barat', pelapor: 'Pak Joko (Penjaga)', status: 'Selesai', icon: 'door_sliding', urgensiColor: 'green', statusColor: 'green' }
])

// Mock Data Agenda Hari Ini
const agendaHariIni = ref([
  { id: 1, judul: 'Kunjungan Teknisi Internet', jam: '09:00 WIB', icon: 'wifi', color: 'primary', deskripsi: 'Pengecekan bandwidth area Lab Komputer untuk persiapan ujian.' },
  { id: 2, judul: 'Batas Pengembalian Kamera', jam: '14:00 WIB', icon: 'photo_camera', color: 'amber', deskripsi: 'Kamera DSL-02 dipinjam OSIS untuk dokumentasi rapat yayasan.' },
  { id: 3, judul: 'Stok Opname Bulanan', jam: '16:00 WIB', icon: 'inventory', color: 'purple', deskripsi: 'Pengecekan berkala sisa fisik kertas A4 dan spidol di gudang utama.' }
])
</script>
