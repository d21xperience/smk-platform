<template>
  <q-page class="q-pa-md">
    <q-card flat bordered class="q-mx-auto" style="max-width: 800px;">
      <q-card-section class="bg-primary text-white row items-center">
        <q-icon name="add_box" size="sm" class="q-mr-sm" />
        <div class="text-h6 text-weight-bold">Registrasi Inventaris / Aset Baru</div>
      </q-card-section>

      <q-separator />

      <q-card-section class="q-pa-lg">
        <q-form @submit="onSubmit" class="q-gutter-md">

          <!-- SEKSI 1: IDENTITAS BARANG -->
          <div class="text-subtitle2 text-primary text-weight-bold q-mb-sm">I. Identitas & Spesifikasi Barang</div>
          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-6">
              <q-input outlined dense v-model="form.namaBarang" label="Nama Barang *" lazy-rules
                :rules="[val => val && val.length > 0 || 'Nama barang wajib diisi']" />
            </div>
            <div class="col-12 col-sm-6">
              <q-select outlined dense v-model="form.kategori" :options="opsiKategori" label="Kategori Aset *"
                :rules="[val => !!val || 'Kategori wajib dipilih']" />
            </div>
            <div class="col-12 col-sm-6">
              <q-input outlined dense v-model="form.merk" label="Merk / Tipe"
                placeholder="Contoh: BenQ / Epson L3110" />
            </div>
            <div class="col-12 col-sm-6">
              <q-input outlined dense v-model="form.noSeri" label="Nomor Seri Pabrik (S/N)" placeholder="Jika ada" />
            </div>
          </div>

          <q-separator class="q-my-md" />

          <!-- SEKSI 2: LOKASI & PENGADAAN -->
          <div class="text-subtitle2 text-primary text-weight-bold q-mb-sm">II. Lokasi Penempatan & Sumber Dana</div>
          <div class="row q-col-gutter-sm">
            <div class="col-12 col-sm-6">
              <q-select outlined dense v-model="form.lokasi" :options="opsiRuangan" label="Lokasi Penempatan Awal *"
                :rules="[val => !!val || 'Lokasi wajib ditentukan']" />
            </div>
            <div class="col-12 col-sm-6">
              <q-select outlined dense v-model="form.sumberDana"
                :options="['Dana Yayasan', 'Dana BOS', 'Sumbangan/Hibah']" label="Sumber Dana *" />
            </div>
            <div class="col-12 col-sm-6">
              <q-input outlined dense v-model="form.tanggalBeli" label="Tanggal Pembelian" mask="date"
                :rules="['date']">
                <template v-slot:append>
                  <q-icon name="event" class="cursor-pointer">
                    <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                      <q-date v-model="form.tanggalBeli">
                        <div class="row items-center justify-end">
                          <q-btn v-close-popup label="Close" color="primary" flat />
                        </div>
                      </q-date>
                    </q-popup-proxy>
                  </q-icon>
                </template>
              </q-input>
            </div>
            <div class="col-12 col-sm-6">
              <q-input outlined dense v-model.number="form.hargaBeli" type="number" label="Harga Perolehan (Rupiah) *"
                prefix="Rp" :rules="[val => val > 0 || 'Harga harus lebih dari 0']" />
            </div>
          </div>

          <!-- PRATINJAU KODE ASET OTOMATIS -->
          <div class="bg-blue-1 text-blue-9 q-pa-sm rounded-borders text-center text-weight-medium q-mt-md">
            Rencana Kode Unik Aset: <span class="text-underline">{{ generateKodeAset }}</span>
          </div>

          <!-- TOMBOL AKSI -->
          <div class="row justify-end q-gutter-sm q-mt-lg">
            <q-btn label="Reset Form" type="reset" color="amber" flat class="q-ml-sm" @click="resetForm" />
            <q-btn label="Simpan ke Sistem" type="submit" color="primary" icon="save" />
          </div>

        </q-form>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// State Formulir
const form = ref({
  namaBarang: '',
  kategori: null,
  merk: '',
  noSeri: '',
  lokasi: null,
  sumberDana: 'Dana Yayasan',
  tanggalBeli: '2026/06/06', // Menyesuaikan tahun saat ini
  hargaBeli: 0
})

// Opsi Pilihan Dropdown
const opsiKategori = ['Elektronik (ELEK)', 'Mebel & Furnitur (MEBL)', 'Kendaraan (KNDR)', 'Sarana Olahraga (ORAG)']
const opsiRuangan = ['Ruang Kelas 10-A', 'Ruang Kelas 11-A', 'Lab Komputer', 'Ruang Guru', 'Gudang Utama']

// Computed Properti untuk simulasi pembuatan kode unik otomatis di frontend
const generateKodeAset = computed(() => {
  const katKode = form.value.kategori ? form.value.kategori.split('(')[1].replace(')', '') : 'XXXX'
  const tahun = form.value.tanggalBeli ? form.value.tanggalBeli.split('/')[0] : '2026'
  return `YYS/SMA-SW/${katKode}/${tahun}/001` // '001' adalah mock nomor urut
})

// Aksi Submit Form
function onSubmit() {
  $q.notify({
    color: 'green-4',
    textColor: 'white',
    icon: 'cloud_done',
    message: `Aset berhasil didaftarkan dengan Kode: ${generateKodeAset.value}`
  })
  // Di sini nantinya Anda memicu fungsi Pinia store atau Axios POST ke backend
}

function resetForm() {
  form.value = {
    namaBarang: '',
    kategori: null,
    merk: '',
    noSeri: '',
    lokasi: null,
    sumberDana: 'Dana Yayasan',
    tanggalBeli: '2026/06/06',
    hargaBeli: 0
  }
}
</script>
