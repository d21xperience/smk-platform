<template>
  <q-page class="q-pa-md q-pa-sm-xl">
    <!-- HEADER CONTROL PANEL -->
    <div class="row items-center justify-between q-mb-xl">
      <div>
        <h1 class="text-h5 text-sm-h4 text-weight-bold text-grey-9 q-my-none">
          Kelola Berita & Prestasi Terbaru
        </h1>
        <p class="text-grey-6 q-mt-xs q-mb-none text-caption text-sm-body2">
          Ruang kontributor Guru & Kaprog. Tulis artikel pengumuman, dokumentasi kegiatan, atau
          publikasi juara prestasi siswa.
        </p>
      </div>
      <!-- TOMBOL UNTUK MEMBUKA DIALOG WRITE ARTICLE -->
      <q-btn
        color="primary"
        icon="create"
        label="Tulis Artikel Baru"
        class="text-weight-bold shadow-2 q-px-md"
        @click="bukaFormArtikel"
      />
    </div>

    <!-- 1. FILTER DAN ARSIP TABEL KONTEN -->
    <q-card flat bordered class="bg-white rounded-borders shadow-1">
      <q-table
        title="Daftar Tulisan Berita Anda"
        :rows="articleRows"
        :columns="tableColumns"
        row-key="id"
        flat
        class="text-grey-9"
        no-data-label="Belum ada artikel yang ditulis. Mulai bagikan berita pertama sekolah Anda!"
      >
        <!-- Custom Tampilan Badge Jenis Konten (Berita / Prestasi) -->
        <template v-slot:body-cell-type="props">
          <q-td :props="props">
            <q-badge
              :color="props.row.type === 'Prestasi' ? 'positive' : 'info'"
              class="text-weight-bold"
            >
              {{ props.row.type }}
            </q-badge>
          </q-td>
        </template>

        <!-- Custom Tampilan Status Persetujuan -->
        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-badge
              :color="props.row.status === 'Published' ? 'primary' : 'grey-7'"
              class="text-weight-bold"
            >
              {{ props.row.status }}
            </q-badge>
          </q-td>
        </template>

        <!-- Tombol Aksi Manipulasi Data -->
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" class="q-gutter-xs">
            <q-btn flat round dense color="primary" icon="edit" @click="editArtikel(props.row)">
              <q-tooltip>Edit Draf</q-tooltip>
            </q-btn>
            <q-btn flat round dense color="negative" icon="delete" @click="hapusArtikel(props.row)">
              <q-tooltip>Hapus</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </q-table>
    </q-card>

    <!-- ====================================================================== -->
    <!-- 2. POP-UP MODAL: RICH TEXT EDITOR UNTUK GURU                          -->
    <!-- ====================================================================== -->
    <q-dialog
      v-model="formDialogOpen"
      persistent
      maximized
      transition-show="slide-up"
      transition-hide="slide-down"
    >
      <q-card class="bg-grey-1">
        <!-- BAR ATAS MODAL -->
        <q-header class="bg-primary text-white">
          <q-toolbar>
            <q-btn flat label="Batal" v-close-popup class="text-weight-bold" />
            <q-toolbar-title class="text-center text-weight-bold"
              >Ruang Kerja Penulis Artikel</q-toolbar-title
            >
            <q-btn
              label="Terbitkan Sekarang"
              color="amber"
              text-color="black"
              class="text-weight-bold"
              :loading="isSaving"
              @click="simpanArtikel"
            />
          </q-toolbar>
        </q-header>

        <!-- WORKSPACE AREA EDITING -->
        <q-page-container class="q-pa-md q-pa-sm-xl flex flex-center">
          <q-card
            style="width: 100%; max-width: 900px"
            class="q-pa-lg shadow-3 bg-white rounded-borders"
          >
            <q-form class="q-gutter-md">
              <div class="row q-col-gutter-md">
                <!-- Input Judul Artikel -->
                <div class="col-12 col-sm-8">
                  <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                    >Judul Berita / Prestasi</q-item-label
                  >
                  <q-input
                    outlined
                    v-model="newArticle.title"
                    placeholder="Contoh: Tim TBSM Raih Podium Utama Lomba Mekanik Nasional"
                    dense
                  />
                </div>

                <!-- Pilihan Tipe Konten -->
                <div class="col-12 col-sm-4">
                  <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                    >Jenis Publikasi</q-item-label
                  >
                  <q-select
                    outlined
                    v-model="newArticle.type"
                    :options="['Berita Kegiatan', 'Prestasi']"
                    dense
                  />
                </div>
              </div>

              <!-- Input Gambar Cover Menggunakan URL Gambar Pintar -->
              <div>
                <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                  >Tautan Gambar Sampul (Cover URL)</q-item-label
                >
                <q-input
                  outlined
                  v-model="newArticle.cover"
                  placeholder="Masukkan URL gambar dokumentasi (jpg/png)"
                  dense
                />
              </div>

              <!-- KOMPONEN RICH TEXT EDITOR (QEDITOR) -->
              <div>
                <q-item-label class="text-weight-bold text-grey-8 q-mb-xs"
                  >Isi Lengkap Artikel / Berita</q-item-label
                >
                <q-editor
                  v-model="newArticle.content"
                  :definitions="editorDefinitions"
                  :toolbar="editorToolbar"
                  min-height="15rem"
                  class="rounded-borders border-light bg-grey-1"
                />
              </div>
            </q-form>
          </q-card>
        </q-page-container>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const isSaving = ref(false)
const formDialogOpen = ref(false)

// Struktur Data Konstruksi Artikel Baru
const newArticle = ref({
  id: null,
  title: '',
  type: 'Berita Kegiatan',
  cover: '',
  content: '',
})

// Skema Kolom Tabel Arsip
const tableColumns = [
  { name: 'title', align: 'left', label: 'Judul Publikasi', field: 'title', sortable: true },
  { name: 'type', align: 'center', label: 'Kategori', field: 'type' },
  { name: 'date', align: 'center', label: 'Tanggal Buat', field: 'date' },
  { name: 'author', align: 'center', label: 'Penulis', field: 'author' },
  { name: 'status', align: 'center', label: 'Status', field: 'status' },
  { name: 'actions', align: 'center', label: 'Aksi', field: 'actions' },
]

// Mock Data Sinkronisasi dengan Baris Berita di IndexPage Utama
const articleRows = ref([
  {
    id: 1,
    title:
      'Siswa TBSM SMK Pasundan Jatinangor Raih Juara 1 Kompetensi Mekanik Astra Honda Tingkat Provinsi',
    type: 'Prestasi',
    date: '30-05-2026',
    author: 'Kaprog TBSM',
    status: 'Published',
  },
  {
    id: 2,
    title: 'Pelaksanaan Servis Motor Gratis oleh Komunitas Siswa TBSM bagi Warga Sekitar Sekolah',
    type: 'Berita Kegiatan',
    date: '25-05-2026',
    author: 'Humas Hub',
    status: 'Published',
  },
])

// Konfigurasi Tombol Toolbar Pengolah Kata QEditor
const editorToolbar = [
  ['bold', 'italic', 'underline', 'strike'],
  ['quote', 'unordered', 'ordered'],
  [
    {
      label: $q.lang.editor.align,
      icon: $q.iconSet.editor.align,
      fixedLabel: true,
      options: ['left', 'center', 'right', 'justify'],
    },
  ],
  ['undo', 'redo'],
]

function bukaFormArtikel() {
  newArticle.value = {
    id: null,
    title: '',
    type: 'Berita Kegiatan',
    cover: '',
    content: 'Tulis isi pengumuman atau berita di sini...',
  }
  formDialogOpen.value = true
}

function simpanArtikel() {
  if (!newArticle.value.title) {
    $q.notify({ type: 'negative', message: 'Judul tulisan tidak boleh kosong!' })
    return
  }

  isSaving.value = true
  setTimeout(() => {
    isSaving.value = false
    formDialogOpen.value = false

    // Logika Simpan Data Baru / Update
    articleRows.value.unshift({
      id: Date.now(),
      title: newArticle.value.title,
      type: newArticle.value.type,
      date: '31-05-2026',
      author: localStorage.getItem('user_name') || 'Guru Kontributor',
      status: 'Published',
    })

    $q.notify({
      type: 'positive',
      message: 'Artikel Berita berhasil disimpan dan otomatis terbit ke halaman depan!',
      position: 'top',
    })
  }, 1000)
}

function editArtikel(row) {
  newArticle.value = { ...row, content: 'Draf konten lama berhasil di-load kembali ke editor.' }
  formDialogOpen.value = true
}

function hapusArtikel(row) {
  $q.dialog({
    title: 'Hapus Berita',
    message: `Hapus publikasi artikel "${row.title}" dari database sekolah?`,
    cancel: true,
  }).onOk(() => {
    articleRows.value = articleRows.value.filter((a) => a.id !== row.id)
    $q.notify({ type: 'info', message: 'Artikel berhasil dihapus.' })
  })
}
</script>

<style scoped>
.border-light {
  border: 1px solid #e0e0e0;
}
</style>
