<!--
FILE: src/components/common/ReadOnlyContextBanner.vue
STATUS: NEW
STATUS IMPLEMENTASI: COMPLETE
-->

<template>
  <q-banner v-if="isReadOnlyContext" class="bg-warning text-dark q-mb-md">
    <template #avatar>
      <q-icon name="history" />
    </template>

    {{ readOnlyMessage }}
  </q-banner>
</template>

<script setup>
import { useReadOnlyContext } from '../../composables/useReadOnlyContext.js'

const {
  isReadOnlyContext,
  readOnlyMessage,
} = useReadOnlyContext()
</script>
