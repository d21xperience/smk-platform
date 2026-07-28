<template>
  <div v-if="!auth.initialized" class="window-height window-width flex flex-center bg-grey-1">
    <div class="text-center">
      <q-spinner color="primary" size="4em" />
      <div class="text-h6 q-mt-md">Memuat sesi...</div>
    </div>
  </div>
  <router-view v-else />
</template>

<script setup>
import { useAuth } from '@/composables/useAuth'
const auth = useAuth()
</script>
