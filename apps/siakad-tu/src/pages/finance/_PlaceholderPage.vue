<template>
  <q-page class="q-pa-md">
    <q-card>
      <q-card-section>
        <div class="text-h5 q-mb-sm">
          <q-icon name="construction" class="q-mr-sm" />
          {{ pageTitle }}
        </div>
        <div class="text-body2 text-grey-7">
          Halaman ini sedang dalam pengembangan. Akan tersedia di rilis berikutnya.
        </div>
      </q-card-section>
      <q-separator />
      <q-card-section>
        <div class="text-caption text-grey-6">
          Route: <code>{{ $route.path }}</code>
        </div>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
const route = useRoute()
const pageTitle = computed(() => route.meta?.title || 'Halaman')
</script>
