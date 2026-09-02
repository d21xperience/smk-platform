<template>
  <q-page class="flex flex-center">
    <div class="text-center">
      <div class="text-h1">404</div>
      <div class="text-h6 q-mb-md">Halaman tidak ditemukan</div>
      <q-btn color="primary" label="Kembali ke Dashboard" :to="{ name: 'dashboard' }" />
    </div>
  </q-page>
</template>

<script>
export default {
  name: 'ErrorNotFound'
}
</script>
