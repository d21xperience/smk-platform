<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Karantina PPDB & Keuangan</h1>
        <p class="text-caption text-grey-7 q-mb-none">Validasi data calon siswa, verifikasi berkas fisik, otomatisasi
          NIS, dan pencatatan angsuran biaya pendaftaran</p>
      </div>
      <div class="q-gutter-sm">
        <q-btn color="orange-9" icon="auto_awesome" label="Otomatisasi NIS Massal" @click="confirmGenerateNisMassal"
          :disable="calonSiswa.length === 0" />
        <q-btn color="primary" icon="person_add" label="Tambah Calon Siswa" @click="openAddDialog" />
      </div>
    </div>

    <!-- Widget Real-time Kuota & Keuangan -->
    <div class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-4">
        <q-card flat class="shadow-1 bg-white text-dark">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase text-weight-bold">Total Calon Siswa</div>
              <div class="text-h5 text-weight-bold text-primary">{{ calonSiswa.length }} Siswa Baru</div>
            </div>
            <q-avatar color="blue-1" text-color="primary" icon="groups" size="40px" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card flat class="shadow-1 bg-white text-dark">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-grey-7 text-uppercase text-weight-bold">Target Biaya Pendaftaran</div>
              <div class="text-h5 text-weight-bold text-grey-9">Rp {{ standarBiaya.toLocaleString('id-ID') }}</div>
            </div>
            <q-avatar color="grey-2" text-color="grey-8" icon="payments" size="40px" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card flat class="shadow-1 bg-teal-1 text-teal-10">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-teal-9 text-uppercase text-weight-bold flex items-center">
                Total Dana PPDB Masuk
                <!-- Tombol Mata Toggle Sembunyikan/Tampilkan -->
                <q-btn flat round dense :icon="isDanaVisible ? 'visibility' : 'visibility_off'" color="teal-10"
                  size="xs" class="q-ml-xs" @click="isDanaVisible = !isDanaVisible">
                  <q-tooltip>{{ isDanaVisible ? 'Sembunyikan Nominal' : 'Tampilkan Nominal' }}</q-tooltip>
                </q-btn>
              </div>

              <!-- Logika Pengkondisian Teks Nominal -->
              <div class="text-h5 text-weight-bold">
                {{ isDanaVisible ? 'Rp ' + totalDanaMasuk.toLocaleString('id-ID') : 'Rp ••••••••' }}
              </div>
            </div>
            <q-avatar color="teal-2" text-color="teal-10" icon="price_check" size="40px" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Tabel Karantina Calon Siswa PPDB -->
    <q-card class="shadow-1">
      <q-card-section class="q-py-sm flex justify-between items-center">
        <div class="text-subtitle2 text-weight-bold text-grey-9">Daftar Karantina Registrasi Siswa Baru</div>
        <q-input v-model="search" dense outlined placeholder="Cari nama calon siswa..." style="max-width: 250px;">
          <template v-slot:prepend><q-icon name="search" /></template>
        </q-input>
      </q-card-section>

      <q-card-section class="q-pa-none">
        <q-table flat :rows="filteredCalonSiswa" :columns="columns" row-key="id" :pagination="{ rowsPerPage: 10 }">
          <!-- Kustomisasi Badge Jurusan Pilihan -->
          <template v-slot:body-cell-jurusan="props">
            <q-td :props="props">
              <q-badge color="indigo-1" text-color="indigo" :label="props.row.jurusan" class="text-weight-bold" />
            </q-td>
          </template>

          <!-- Kustomisasi Kolom Status Keuangan & Angsuran Tunai -->
          <template v-slot:body-cell-keuangan="props">
            <q-td :props="props">
              <div class="text-center">
                <q-badge :color="props.row.bayar === standarBiaya ? 'green' : props.row.bayar > 0 ? 'orange-9' : 'red'"
                  :label="props.row.bayar === standarBiaya ? 'Lunas (Tunai)' : props.row.bayar > 0 ? 'Dicicil' : 'Belum Bayar'"
                  class="q-mb-xs" />
                <div class="text-caption text-grey-7" style="font-size: 11px;">
                  Masuk: Rp {{ props.row.bayar.toLocaleString('id-ID') }}
                </div>
                <div v-if="standarBiaya - props.row.bayar > 0" class="text-caption text-red text-weight-medium"
                  style="font-size: 10px;">
                  Sisa: Rp {{ (standarBiaya - props.row.bayar).toLocaleString('id-ID') }}
                </div>
              </div>
            </q-td>
          </template>

          <!-- Kustomisasi Kolom Ceklis Dokumen Fisik Persyaratan -->
          <template v-slot:body-cell-berkas="props">
            <q-td :props="props" class="text-center">
              <div class="row justify-center q-gutter-xs">
                <q-icon :name="props.row.berkas.ijazah ? 'check_circle' : 'cancel'"
                  :color="props.row.berkas.ijazah ? 'green' : 'red'"
                  size="xs"><q-tooltip>Ijazah/SKL</q-tooltip></q-icon>
                <q-icon :name="props.row.berkas.kk ? 'check_circle' : 'cancel'"
                  :color="props.row.berkas.kk ? 'green' : 'red'" size="xs"><q-tooltip>Kartu
                    Keluarga</q-tooltip></q-icon>
                <q-icon :name="props.row.berkas.akta ? 'check_circle' : 'cancel'"
                  :color="props.row.berkas.akta ? 'green' : 'red'" size="xs"><q-tooltip>Akta Lahir</q-tooltip></q-icon>
              </div>
            </q-td>
          </template>

          <!-- Kustomisasi Kolom NIS -->
          <template v-slot:body-cell-nis="props">
            <q-td :props="props" class="text-center">
              <q-badge v-if="props.row.nis" color="green-8" :label="props.row.nis" />
              <q-badge v-else color="grey-6" label="Belum Generate" />
            </q-td>
          </template>

          <!-- Kustomisasi Kolom Aksi -->
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center q-gutter-xs">
              <q-btn size="sm" color="amber-9" icon="edit" round flat @click="editCalonSiswa(props.row)">
                <q-tooltip>Ubah Data & Angsuran Keuangan</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="red" icon="no_accounts" round flat @click="mundurSiswa(props.row)">
                <q-tooltip>Siswa Mengundurkan Diri</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="green-8" icon="bolt" round flat
                :disable="!props.row.nis || !isLunasDanLengkap(props.row)" @click="aktivasiKeSiswaAktif(props.row)">
                <q-tooltip v-if="isLunasDanLengkap(props.row)">Aktifkan ke Direktori Utama</q-tooltip>
                <q-tooltip v-else>Gagal: Biaya Belum Lunas atau Berkas Fisik Tidak Lengkap!</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- DIALOG POPUP: Tambah / Edit Calon Siswa (Termasuk Keuangan) -->
    <q-dialog v-model="dialogOpen" persistent>
      <q-card style="min-width: 460px;">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">{{ isEdit ? 'Kelola Data & Transaksi Calon Siswa' : 'Registrasi Calon Siswa Baru' }}
          </div>
        </q-card-section>

        <q-card-section class="q-gutter-sm q-pt-md">
          <div class="text-subtitle2 text-grey-8">Biodata Pokok Calon Siswa:</div>
          <q-input v-model="form.nama" label="Nama Lengkap Calon Siswa" outlined dense />
          <q-input v-model="form.asalSekolah" label="Sekolah Asal (SMP / MTs)" outlined dense />

          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-select v-model="form.jurusan" :options="['TKRO', 'RPL', 'AKL']" label="Jurusan Pilihan" outlined
                dense />
            </div>
            <div class="col-6">
              <q-input v-model="form.nisn" label="NISN (Dari SMP)" outlined dense type="number" />
            </div>
          </div>

          <q-separator class="q-my-md" />

          <!-- INTEGRASI MODUL KEUANGAN PPDB -->
          <div class="text-subtitle2 text-grey-8">Administrasi Biaya Pendaftaran (Tunai/Cicil):</div>
          <div class="row q-col-gutter-sm items-center">
            <div class="col-6">
              <q-input v-model.number="form.bayar" label="Jumlah Uang yang Dibayarkan (Rp)" outlined dense type="number"
                prefix="Rp" :rules="[val => val <= standarBiaya || 'Pembayaran melebihi batas tarif sekolah!']" />
            </div>
            <div class="col-6 q-pb-md">
              <div class="text-caption text-grey-7">Tungggakan Sisa Angsuran:</div>
              <div class="text-subtitle2 text-weight-bold"
                :class="standarBiaya - form.bayar === 0 ? 'text-green' : 'text-red'">
                Rp {{ (standarBiaya - form.bayar).toLocaleString('id-ID') }}
              </div>
            </div>
          </div>

          <div class="row q-gutter-xs q-mb-sm">
            <q-btn label="Set Lunas (Tunai)" size="xs" color="green" flat class="bg-green-1"
              @click="form.bayar = standarBiaya" />
            <q-btn label="Set Belum Bayar" size="xs" color="red" flat class="bg-red-1" @click="form.bayar = 0" />
          </div>

          <q-separator class="q-my-sm" />

          <!-- Verifikasi Berkas Fisik PPDB -->
          <div class="text-subtitle2 text-grey-8">Verifikasi Dokumen Fisik Pendaftaran:</div>
          <div class="row q-mt-xs">
            <div class="col-4"><q-checkbox v-model="form.berkas.ijazah" label="Ijazah / SKL" dense color="primary" />
            </div>
            <div class="col-4"><q-checkbox v-model="form.berkas.kk" label="Krt. Keluarga" dense color="primary" /></div>
            <div class="col-4"><q-checkbox v-model="form.berkas.akta" label="Akta Lahir" dense color="primary" /></div>
          </div>

        </q-card-section>
        <q-btn color="primary" label="Simpan Registrasi" @click="saveCalonSiswa" />
      </q-card>
    </q-dialog>
  </q-page>

  <!-- ====== TAMBAHKAN STRUKTUR HTML KUITANSI THERMAL INI ====== -->
  <div v-if="kuitansiAktif" class="area-kuitansi-thermal">
    <div class="text-center text-weight-bold" style="font-size: 14px; margin-bottom: 2px;">
      SMK SWASTA INDONESIA
    </div>
    <div class="text-center text-caption"
      style="font-size: 10px; margin-bottom: 8px; border-bottom: 1px dashed #000; padding-bottom: 4px;">
      Jl. Raya Pendidikan No. 123<br>PANITIA KEUANGAN PPDB
    </div>

    <!-- Info Detail Transaksi -->
    <div style="font-size: 11px; line-height: 1.4; margin-bottom: 6px;">
      <div><strong>No. Kuitansi :</strong> {{ kuitansiAktif.nomorKuitansi }}</div>
      <div><strong>Tanggal :</strong> {{ kuitansiAktif.tanggal }}</div>
      <div><strong>Waktu :</strong> {{ kuitansiAktif.jam }}</div>
    </div>
    <div style="border-bottom: 1px dashed #000; margin-bottom: 6px;"></div>

    <!-- Detail Nominal Pembayaran -->
    <div style="font-size: 11px; line-height: 1.4;">
      <div><strong>Nama Calon :</strong> {{ kuitansiAktif.nama }}</div>
      <div><strong>Jurusan :</strong> {{ kuitansiAktif.jurusan }}</div>
      <div style="margin-top: 4px;"><strong>Jenis Biaya :</strong> Biaya Pendaftaran PPDB</div>
      <div><strong>Jumlah Bayar :</strong> Rp {{ kuitansiAktif.nominal.toLocaleString('id-ID') }}</div>
      <div :style="kuitansiAktif.sisa > 0 ? 'color: red;' : ''">
        <strong>Sisa Tagihan :</strong> Rp {{ kuitansiAktif.sisa.toLocaleString('id-ID') }}
      </div>
    </div>
    <div style="border-bottom: 1px dashed #000; margin-top: 6px; margin-bottom: 6px;"></div>

    <!-- Status Kelulusan Keuangan -->
    <div class="text-center text-weight-bold text-uppercase"
      style="font-size: 12px; margin-top: 4px; margin-bottom: 12px;">
      *** {{ kuitansiAktif.status }} ***
    </div>

    <!-- Footer Kuitansi -->
    <div class="text-center text-caption" style="font-size: 9px; line-height: 1.2;">
      Terima kasih atas pembayaran Anda.<br>Harap simpan struk kuitansi ini<br>sebagai bukti pendaftaran yang sah.
    </div>
    <div class="text-center text-caption" style="font-size: 9px; margin-top: 10px;">
      SIAKAD TU Kesiswaan v1.0
    </div>
  </div>
  <!-- ========================================================== -->
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'
import { useRouter } from 'vue-router' // Taruh import ini di paling atas jika belum ada

const router = useRouter()
const $q = useQuasar()
// Konfigurasi Biaya Resmi Sekolah Swasta
const standarBiaya = 450000
// State Komponen
const search = ref('')
const dialogOpen = ref(false)
const isEdit = ref(false)
const form = ref({ id: null, nama: '', nisn: '', asalSekolah: '', jurusan: 'TKRO', gender: 'L', bayar: 0, berkas: { ijazah: false, kk: false, akta: false } })
// Struktur Judul Kolom Tabel Karantina PPDB Swasta yang Baru
const columns = [{ name: 'nama', align: 'left', label: 'Nama Calon Siswa', field: 'nama', sortable: true }, { name: 'jurusan', align: 'left', label: 'Jurusan', field: 'jurusan', sortable: true }, { name: 'keuangan', align: 'center', label: 'Status Biaya Masuk (Rp 450k)' }, { name: 'berkas', align: 'center', label: 'Berkas Fisik' }, { name: 'nis', align: 'center', label: 'NIS', field: 'nis' }, { name: 'aksi', align: 'center', label: 'Aksi Karantina' }]
// Data Dummy Calon Siswa Baru PPDB dengan Record Angsuran Keuangan
const calonSiswa = ref([{ id: 1, nama: 'Bagus Setiawan', nisn: '0091122331', asalSekolah: 'SMPN 1 Soreang', jurusan: 'TKRO', gender: 'L', nis: '', bayar: 450000, berkas: { ijazah: true, kk: true, akta: true } }, { id: 2, nama: 'Aditya Pratama', nisn: '0091122332', asalSekolah: 'SMPN 2 Bandung', jurusan: 'TKRO', gender: 'L', nis: '', bayar: 200000, berkas: { ijazah: true, kk: true, akta: true } },
// Dicicil
{ id: 3, nama: 'Chandra Kirana', nisn: '0091122333', asalSekolah: 'MTs As-Syifa', jurusan: 'RPL', gender: 'P', nis: '', bayar: 450000, berkas: { ijazah: true, kk: true, akta: true } }, { id: 4, nama: 'Deni Sukma', nisn: '0091122334', asalSekolah: 'SMP Pasundan 1', jurusan: 'AKL', gender: 'L', nis: '', bayar: 0, berkas: { ijazah: false, kk: false, akta: false } }
  // Belum Bayar
])
const isDanaVisible = ref(false)
const urutanJurusanSekolah = ['TKRO', 'RPL', 'AKL']
// Perhitungan Kumulatif Total Dana Masuk Kas PPDB Sekolah untuk Isian Widget Atas
const totalDanaMasuk = computed(() => { return calonSiswa.value.reduce((total, siswa) => total + siswa.bayar, 0) })
// Helper Validasi Ganda Kelayakan Aktivasi (Wajib Lunas + Berkas Lengkap)
const isLunasDanLengkap = (siswa) => {
  const lunas = siswa.bayar === standarBiaya
  const berkasLengkap = siswa.berkas.ijazah && siswa.berkas.kk && siswa.berkas.akta
  return lunas && berkasLengkap
}
// Logika Filter Pencarian & Urutan Baris Abjad/Jurusan Sisi Klien
const filteredCalonSiswa = computed(() => {
  let result = [...calonSiswa.value]
  if (search.value) { result = result.filter(siswa => siswa.nama.toLowerCase().includes(search.value.toLowerCase())) } return result.sort((a, b) => {
    const indeksA = urutanJurusanSekolah.indexOf(a.jurusan)
    const indeksB = urutanJurusanSekolah.indexOf(b.jurusan)
    if (indeksA !== indeksB) return indeksA - indeksB
    return a.nama.localeCompare(b.nama)
  })
})
const openAddDialog = () => {
  isEdit.value = false
  form.value = { id: null, nama: '', nisn: '', asalSekolah: '', jurusan: 'TKRO', gender: 'L', bayar: 0, berkas: { ijazah: false, kk: false, akta: false } }
  dialogOpen.value = true
}
const editCalonSiswa = (row) => {
  // isEdit.value = true
  // form.value = JSON.parse(JSON.stringify(row))
  // dialogOpen.value = true
  router.push(`/siswa/ppdb/edit/${row.id}`)
}
const saveCalonSiswa = () => {
  if (!form.value.nama || form.value.bayar === undefined) return
  // if (isEdit.value) {
  //   const target = calonSiswa.value.find(item => item.id === form.value.id)
  //   if (target) { Object.assign(target, form.value) }
  //   $q.notify({ color: 'green-8', message: 'Data registrasi & nominal pembayaran sukses diperbarui.', icon: 'check' })
  // } else {
  calonSiswa.value.push({ id: Date.now(), ...form.value, nis: '' })
  $q.notify({ color: 'primary', message: 'Data pendaftar baru disimpan ke list PPDB.', icon: 'person_add' })
  // }
  dialogOpen.value = false
  triggerKonfirmasiCetak(form.value)



}
const confirmGenerateNisMassal = () => {
  $q.dialog({ title: 'Generate NIS Massal', message: 'Proses otomatisasi pembuatan nomor induk siswa massal berurutan untuk angkatan baru tahun ajaran 2026/2027.', cancel: true, persistent: true }).onOk(() => {
    const sortedList = filteredCalonSiswa.value
    let counterUrut = 1
    const tahunAjaranPrefix = '2627'
    sortedList.forEach((siswa) => {
      const nomorUrutPadded = String(counterUrut).padStart(4, '0')
      const originalSiswa = calonSiswa.value.find(item => item.id === siswa.id)
      if (originalSiswa) { originalSiswa.nis = `${tahunAjaranPrefix}${nomorUrutPadded}` }
      counterUrut++
    })
    $q.notify({ color: 'green-9', message: 'NIS Massal sukses diterbitkan.', icon: 'auto_awesome' })
  })
}
const mundurSiswa = (row) => {
  $q.dialog({ title: 'Siswa Mengundurkan Diri', message: `Hapus ${row.nama} ? Uang pendaftaran yang masuk sebesar Rp ${row.bayar.toLocaleString('id-ID')} wajib dikembalikan sesuai aturan administrasi yayasan.`, cancel: true, color: 'red' }).onOk(() => {
    calonSiswa.value = calonSiswa.value.filter(item => item.id !== row.id)
    $q.notify({ color: 'negative', message: 'Siswa berhasil dihapus dari sistem kesiswaan.', icon: 'delete' })
  })
}
const aktivasiKeSiswaAktif = (row) => {
  $q.notify({ color: 'green-9', message: `Sukses! ${row.nama} dinyatakan LUNAS & LENGKAP.Resmi dipindahkan ke Direktori Utama Siswa Aktif SMK.`, icon: 'done_all' })
  calonSiswa.value = calonSiswa.value.filter(item => item.id !== row.id)
}

// ====== TAMBAHKAN LOGIKA CETAK KUITANSI THERMAL ======
// State penampung data kuitansi yang sedang aktif dicetak
const kuitansiAktif = ref(null)

const cetakKuitansiThermal = (siswaBaru) => {
  console.log(siswaBaru)
  // 1. Masukkan data siswa baru ke objek kuitansi aktif
  kuitansiAktif.value = {
    nomorKuitansi: 'KW-PPDB-' + Date.now().toString().slice(-6),
    tanggal: new Date().toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' }),
    jam: new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) + ' WIB',
    nama: siswaBaru.nama,
    jurusan: siswaBaru.jurusan,
    nominal: siswaBaru.bayar,
    sisa: 450000 - siswaBaru.bayar, // Menghitung sisa tunggakan dari tarif Rp 450.000
    status: siswaBaru.bayar === 450000 ? 'LUNAS (TUNAI)' : 'ANGSURAN / CICIL'
  }

  // 2. Beri jeda beberapa milidetik agar Vue selesai merender komponen kuitansi tersembunyi, lalu pemicu printer berjalan
  setTimeout(() => {
    window.print()
  }, 300)
}

// Modifikasi fungsi saveCalonSiswa lama Anda pada bagian "Tambah Siswa Baru"
// Agar ketika sukses menyimpan data siswa baru, sistem langsung memicu dialog konfirmasi cetak kuitansi
// const saveCalonSiswaLama = saveCalonSiswa; // Referensi fungsi lama Anda jika diperlukan

// Pastikan di akhir blok logika penambahan siswa baru (bukan saat edit), Anda memanggil fungsi konfirmasi ini:
const triggerKonfirmasiCetak = (siswaBaru) => {
  $q.dialog({
    title: 'Pendaftaran Sukses',
    message: `Calon siswa ${siswaBaru.nama} berhasil didaftarkan. Apakah Anda ingin langsung mencetak kuitansi pembayaran digital sekarang?`,
    cancel: { label: 'Batal', flat: true, color: 'grey' },
    ok: { label: 'Cetak Kuitansi', color: 'primary', icon: 'print' },
    persistent: true
  }).onOk(() => {
    cetakKuitansiThermal(siswaBaru)
  })
}
// =====================================================




</script>


<style scoped>
/* Pengaturan standar tampilan layout halaman pendaftaran */
.text-h4 {
  line-height: 1.1;
}

/* ====== TAMBAHKAN LOGIKA STYLING PRINTER THERMAL INI ====== */
/* Secara default, sembunyikan struk kuitansi saat aplikasi dibuka di layar monitor */
.area-kuitansi-thermal {
  display: none;
}

/* Perintah CSS khusus saat tombol cetak browser aktif (@media print) */
@media print {

  /* Sembunyikan seluruh elemen dashboard utama sekolah, header, sidebar, dan tabel */
  body *,
  .q-layout,
  .q-page,
  .q-drawer,
  .q-header,
  .q-card {
    display: none !important;
  }

  /* Hanya tampilkan area komponen kuitansi pendaftaran */
  .area-kuitansi-thermal {
    display: block !important;
    position: absolute;
    left: 0;
    top: 0;
    width: 58mm;
    /* Sesuaikan menjadi 80mm jika Anda menggunakan mesin thermal printer ukuran besar */
    padding: 2mm;
    font-family: 'Courier New', Courier, monospace;
    /* Font khas mesin kasir agar rapi */
    color: #000;
    background: #fff;
  }
}

/* ========================================================== */
</style>
