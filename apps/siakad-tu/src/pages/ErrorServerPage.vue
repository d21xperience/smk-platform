<template>
  <div class="text-center q-pa-md">
    <h1>Oops, Terjadi Kesalahan!</h1>
    <p class="text-red">{{ route.query.message }}</p>
    <q-btn to="/" label="Kembali ke Beranda" color="primary" />
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
const route = useRoute()
</script>
