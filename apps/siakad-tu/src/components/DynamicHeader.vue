<template>
  <!-- <q-page class="q-pa-md bg-grey-1"> -->
  <!-- Judul Halaman Dinamis via Props -->
  <div class="row items-center q-mb-md">
    <q-btn flat round icon="arrow_back" color="primary" @click="$router.back()" />
    <h5 class="text-h5 text-weight-bold q-my-none q-ml-sm text-primary">{{ title }}</h5>
  </div>
  <!-- </q-page> -->
</template>

<script setup>
// 1. Mendefinisikan PROPS agar komponen flexibel
defineProps({
  title: String,               // Contoh: 'Guru' atau 'Pegawai'
  // columns: Array,             // Struktur kolom pratinjau tabel
  // previewRows: Array,         // Data baris hasil baca excel untuk pratinjau
  // pastePlaceholderHint: String, // Petunjuk teks format kolom copy-paste
  // examplePasteText: String     // Contoh isi data teks di dalam textarea
})

</script>
