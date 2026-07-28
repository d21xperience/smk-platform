<!-- ============================================================
  components/SectionHeader.vue
  Header generik untuk setiap section (pill + judul + sub)
  Props: title, subtitle, align (left|center|right)
============================================================ -->
<template>
  <div :class="['section-header', `text-${align}`]">
    <div :class="['section-pill', { 'mx-auto': align === 'center' }]" />
    <div class="text-h5 text-sm-h4 text-weight-bold text-primary q-mb-xs">
      {{ title }}
    </div>
    <p v-if="subtitle" class="text-body2 text-grey-6 q-mb-none section-subtitle">
      {{ subtitle }}
    </p>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, required: true },
  subtitle: { type: String, default: '' },
  align: { type: String, default: 'left' },
})
</script>

<style scoped>
.section-header {
  margin-bottom: 28px;
}

.section-pill {
  width: 36px;
  height: 4px;
  background: var(--q-primary);
  border-radius: 2px;
  margin-bottom: 12px;
}
.section-subtitle {
  max-width: 560px;
  line-height: 1.7;
}
</style>
