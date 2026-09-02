<template>
  <q-page class="q-pa-md">
    <!-- PANEL HEADER & PENCARI -->
    <div class="row q-col-gutter-sm q-mb-md items-center justify-between">
      <div class="col-12 col-sm-6 col-md-4">
        <q-input outlined dense v-model="filterText" placeholder="Cari vendor atau spesialisasi (ex: AC)...">
          <template v-slot:append>
            <q-icon name="search" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-sm-auto">
        <q-btn color="primary" icon="person_add" label="Tambah Mitra Vendor" @click="bukaDialogTambah" />
      </div>
    </div>

    <!-- GRID DATA KARTU KONTAK VENDOR -->
    <div class="row q-col-gutter-md">
      <div v-for="vendor in vendorTersaring" :key="vendor.id" class="col-12 col-sm-6 col-md-4">
        <q-card flat bordered class="vendor-card">
          <!-- Header Kartu: Nama & Badge Spesialisasi -->
          <q-card-section class="bg-blue-grey-1 row items-center no-wrap">
            <q-avatar icon="store" color="primary" text-color="white" size="md" class="q-mr-sm" />
            <div class="ellipsis">
              <div class="text-subtitle1 text-weight-bold text-grey-9">{{ vendor.nama }}</div>
              <span class="text-caption text-grey-7">ID: {{ vendor.id }}</span>
            </div>
          </q-card-section>

          <q-separator />

          <!-- Konten Utama: Kontak & Alamat -->
          <q-card-section class="q-gutter-xs text-grey-8">
            <div class="row items-center">
              <q-icon name="build" class="q-mr-xs text-primary" />
              <span><strong>Keahlian:</strong>
                <q-badge color="teal" class="q-ml-xs">{{ vendor.spesialisasi }}</q-badge>
              </span>
            </div>

            <div class="row items-center q-mt-sm">
              <q-icon name="phone" class="q-mr-xs text-green" />
              <span><strong>Telepon:</strong> {{ vendor.telepon }}</span>
            </div>

            <div class="row items-center q-mt-sm">
              <q-icon name="person" class="q-mr-xs text-orange" />
              <span><strong>Kontak Person (PIC):</strong> {{ vendor.pic }}</span>
            </div>

            <div class="row items-top q-mt-sm no-wrap">
              <q-icon name="place" class="q-mr-xs text-red q-mt-xs" />
              <div class="ellipsis-2-lines"><strong>Alamat:</strong> {{ vendor.alamat }}</div>
            </div>
          </q-card-section>

          <q-separator />

          <!-- Tombol Aksi Cepat Hubungi -->
          <q-card-actions align="right">
            <q-btn flat round color="green-7" icon="chat" @click="hubungiWhatsApp(vendor.telepon)">
              <q-tooltip>Hubungi via WhatsApp</q-tooltip>
            </q-btn>
            <q-btn flat round color="blue" icon="phone_forwarded" :href="`tel:${vendor.telepon}`">
              <q-tooltip>Panggil Telepon</q-tooltip>
            </q-btn>
            <q-btn flat round color="negative" icon="delete" @click="hapusVendor(vendor)" />
          </q-card-actions>
        </q-card>
      </div>

      <!-- TAMPILAN JIKA TIDAK ADA VENDOR YANG COCOK -->
      <div v-if="vendorTersaring.length === 0" class="col-12 text-center q-pa-xl text-grey-5">
        <q-icon name="person_search" size="64px" class="q-mb-sm" />
        <div class="text-h6">Vendor tidak ditemukan. Sila periksa kata kunci Anda.</div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const filterText = ref('')

// 1. Data Mock Rekanan Vendor / Teknisi Luar Sekolah
const daftarVendor = ref([
  { id: 'VND-01', nama: 'CV Cool Mandiri Jaya', spesialisasi: 'Perbaikan AC', telepon: '081234567890', pic: 'Pak Rahmat', alamat: 'Jl. Ahmad Yani No. 45, Kota Bandung' },
  { id: 'VND-02', nama: 'Abadi Komputer Sistem', spesialisasi: 'IT Hardware & Networking', telepon: '087798765432', pic: 'Ko Alung', alamat: 'Pusat Elektronik Jaya Lt. 2 Blok C' },
  { id: 'VND-03', nama: 'Bengkel Mobil Auto Berkah', spesialisasi: 'Mekanik Kendaraan', telepon: '082111223344', pic: 'Montir Sugeng', alamat: 'Ruko Utama Bypass No. 12B' },
  { id: 'VND-04', nama: 'Toko Bangunan Sumber Utama', spesialisasi: 'Konstruksi & Sipil', telepon: '085244556677', pic: 'Ibu Linda', alamat: 'Jl. Raya Barat No. 102' }
])

// 2. Computed Filter untuk Fitur Pencarian Dinamis di Frontend
const vendorTersaring = computed(() => {
  const query = filterText.value.toLowerCase().trim()
  if (!query) return daftarVendor.value

  return daftarVendor.value.filter(vendor =>
    vendor.nama.toLowerCase().includes(query) ||
    vendor.spesialisasi.toLowerCase().includes(query) ||
    vendor.pic.toLowerCase().includes(query)
  )
})

// 3. Fitur Akses Komunikasi Eksternal
// function hubungiWhatsApp(noTelp) {
function hubungiWhatsApp() {
  // Format nomor untuk link API WhatsApp Indonesia (mengubah 0 jadi 62)
  // const formatWa = noTelp.replace(/^0/, '62')
  const url = `https://whatsapp.com{formatWa}&text=Halo,%20kami%20dari%20Bagian%20Sarpras%20Sekolah...`
  window.open(url, '_blank')
}

// 4. Aksi Tambah & Hapus Data Mock
function bukaDialogTambah() {
  $q.notify({
    message: 'Membuka formulir pendaftaran kemitraan teknisi/vendor baru.',
    color: 'primary',
    icon: 'assignment_ind'
  })
}

function hapusVendor(vendor) {
  $q.dialog({
    title: 'Putuskan Kemitraan',
    message: `Apakah Anda yakin ingin menghapus data "${vendor.nama}" dari daftar rekanan sarpras?`,
    cancel: true,
    persistent: true
  }).onOk(() => {
    daftarVendor.value = daftarVendor.value.filter(v => v.id !== vendor.id)
    $q.notify({
      color: 'negative',
      message: 'Data vendor telah dihapus dari sistem internal.',
      icon: 'delete'
    })
  })
}
</script>

<style scoped>
.vendor-card {
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.vendor-card:hover {
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}
</style>
