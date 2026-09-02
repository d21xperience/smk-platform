<template>
  <q-page class="q-pa-md">
    <!-- INFO STATUS GUDANG (STATISTIK KECIL) -->
    <div class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-4">
        <q-card flat bordered class="bg-blue-1 text-blue-9">
          <q-card-section class="row items-center no-wrap q-py-sm">
            <div>
              <div class="text-caption text-grey-7">Total Jenis Logistik</div>
              <div class="text-h6 text-weight-bold">{{ itemLogistik.length }} Item</div>
            </div>
            <q-space />
            <q-icon name="inventory" size="sm" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card flat bordered class="bg-red-1 text-red-9">
          <q-card-section class="row items-center no-wrap q-py-sm">
            <div>
              <div class="text-caption text-grey-7">Stok Kritis (Wajib Beli)</div>
              <div class="text-h6 text-weight-bold">{{ totalKritis }} Item</div>
            </div>
            <q-space />
            <q-icon name="gpp_bad" size="sm" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card flat bordered class="bg-green-1 text-green-9">
          <q-card-section class="row items-center no-wrap q-py-sm">
            <div>
              <div class="text-caption text-grey-7">Kondisi Stok Aman</div>
              <div class="text-h6 text-weight-bold">{{ itemLogistik.length - totalKritis }} Item</div>
            </div>
            <q-space />
            <q-icon name="check_circle" size="sm" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- TABEL UTAMA LOGISTIK -->
    <q-table title="Katalog Barang Habis Pakai" :rows="itemLogistik" :columns="columns" row-key="id"
      :filter="filterText" flat bordered>
      <!-- Slot Atas Kanan untuk Pencarian & Aksi -->
      <template v-slot:top-right>
        <div class="row q-gutter-sm">
          <q-input borderless dense debounce="300" v-model="filterText" placeholder="Cari ATK/Logistik...">
            <template v-slot:append>
              <q-icon name="search" />
            </template>
          </q-input>
          <q-btn color="primary" icon="add" label="Tambah Item Baru" @click="bukaDialogTambah" />
        </div>
      </template>

      <!-- Kustomisasi Tampilan Kolom Sisa Stok & Indikator Kritis -->
      <template v-slot:body-cell-stok="props">
        <q-td :props="props" :class="props.row.stok <= props.row.minStok ? 'bg-red-1' : ''">
          <div class="row items-center no-wrap justify-center">
            <span class="text-weight-bold q-mr-sm">{{ props.value }} {{ props.row.satuan }}</span>

            <!-- Badge Peringatan jika Stok <= Batas Minimum -->
            <q-badge v-if="props.row.stok <= props.row.minStok" color="red" text-color="white" label="Kritis" />
            <q-badge v-else-if="props.row.stok <= props.row.minStok + 5" color="orange" text-color="white"
              label="Menipis" />
            <q-badge v-else color="green" text-color="white" label="Aman" />
          </div>
        </q-td>
      </template>

      <!-- Kolom Tombol Aksi Cepat (Restock / Ambil Barang) -->
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props" class="text-center">
          <q-btn-group flat>
            <q-btn dense size="sm" color="positive" icon="call_received" label="Restock"
              @click="bukaDialogRestock(props.row)">
              <q-tooltip>Tambah Stok Masuk</q-tooltip>
            </q-btn>
            <q-btn dense size="sm" color="negative" icon="call_made" label="Gunakan" class="q-ml-xs"
              @click="bukaDialogGunakan(props.row)">
              <q-tooltip>Catat Pengambilan Guru/Staf</q-tooltip>
            </q-btn>
          </q-btn-group>
        </q-td>
      </template>
    </q-table>

    <!-- INTERACTIVE DIALOG: CEPAT TAMBAH STOK (RESTOCK) -->
    <q-dialog v-model="dialogRestock.terbuka">
      <q-card style="width: 350px">
        <q-card-section class="bg-positive text-white">
          <div class="text-h6">Restock: {{ dialogRestock.dataItem.nama }}</div>
        </q-card-section>

        <q-card-section class="q-pt-md">
          <q-input outlined dense type="number" v-model.number="dialogRestock.jumlahMasuk" label="Jumlah Barang Masuk"
            :suffix="dialogRestock.dataItem.satuan" autofocus />
        </q-card-section>

        <q-card-actions align="right" class="text-primary">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="positive" label="Konfirmasi" @click="eksekusiRestock" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const filterText = ref('')

// 1. Mock Data Inventaris Habis Pakai Sekolah Swasta
const itemLogistik = ref([
  { id: 'LOG001', name: 'Kertas HVS A4 80gr', category: 'ATK Kelas', stok: 8, minStok: 10, satuan: 'Rim' },
  { id: 'LOG002', name: 'Spidol Whiteboard Hitam Snowman', category: 'ATK Kelas', stok: 45, minStok: 12, satuan: 'Pcs' },
  { id: 'LOG003', name: 'Buku Agenda Kelas Yayasan', category: 'Dokumen', stok: 3, minStok: 5, satuan: 'Buku' },
  { id: 'LOG004', name: 'Sabun Cuci Tangan Cair (Isi Ulang)', category: 'Kebersihan', stok: 15, minStok: 4, satuan: 'Pouch' },
  { id: 'LOG005', name: 'Tinta Printer Epson Hitam 003', category: 'Tinta/Tonner', stok: 1, minStok: 3, satuan: 'Botol' }
])

// 2. Definisi Struktur Kolom Tabel Quasar
const columns = [
  { name: 'id', label: 'ID Logistik', align: 'left', field: 'id', sortable: true },
  { name: 'nama', label: 'Nama Barang Habis Pakai', align: 'left', field: 'name', sortable: true },
  { name: 'category', label: 'Kategori', align: 'left', field: 'category' },
  { name: 'minStok', label: 'Batas Minimum', align: 'center', field: 'minStok' },
  { name: 'stok', label: 'Sisa Stok Saat Ini', align: 'center', field: 'stok', sortable: true },
  { name: 'aksi', label: 'Aset Kontrol Cepat', align: 'center' }
]

// 3. Computed Value untuk Mendeteksi Berapa Item yang Berstatus Kritis
const totalKritis = computed(() => {
  return itemLogistik.value.filter(item => item.stok <= item.minStok).length
})

// 4. Logika State Dialog Kontrol Cepat
const dialogRestock = ref({
  terbuka: false,
  jumlahMasuk: 0,
  dataItem: {}
})

function bukaDialogRestock(item) {
  dialogRestock.value.dataItem = item
  dialogRestock.value.jumlahMasuk = 0
  dialogRestock.value.terbuka = true
}

function eksekusiRestock() {
  const item = itemLogistik.value.find(i => i.id === dialogRestock.value.dataItem.id)
  if (item && dialogRestock.value.jumlahMasuk > 0) {
    item.stok += dialogRestock.value.jumlahMasuk
    $q.notify({
      color: 'green',
      icon: 'done',
      message: `Berhasil menambahkan ${dialogRestock.value.jumlahMasuk} ${item.satuan} ke stok ${item.name}.`
    })
  }
}

function bukaDialogGunakan(item) {
  // Anda bisa membuat struktur dialog yang sama seperti restock, namun untuk mengurangi stok (item.stok -= jumlah)
  $q.notify({
    color: 'orange',
    icon: 'info',
    message: `Fitur pencatatan penggunaan untuk ${item.name} siap dikembangkan.`
  })
}

function bukaDialogTambah() {
  $q.notify({
    message: 'Fitur form pop-up pendaftaran jenis ATK baru.',
    color: 'primary'
  })
}
</script>
