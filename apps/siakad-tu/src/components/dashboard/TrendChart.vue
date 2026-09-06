<template>
  <q-card flat bordered class="content-card">
    <q-card-section>
      <div class="row items-center justify-between">
        <div class="text-subtitle1 text-weight-bold text-primary">{{ title }}</div>
        <q-badge color="grey-3" text-color="grey-8" class="q-px-sm">{{ badgeLabel }}</q-badge>
      </div>
    </q-card-section>
    <q-separator />
    <q-card-section>
      <div class="trend-chart row items-end no-wrap">
        <div
          v-for="item in items"
          :key="item.label"
          class="trend-bar-wrap col text-center"
        >
          <div class="text-caption text-grey-7 q-mb-xs">{{ item.value }}</div>
          <div class="trend-bar" :style="{ height: trendHeight(item.value) + 'px' }" />
          <div class="text-caption text-grey-6 q-mt-xs">{{ item.label }}</div>
        </div>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  items: { type: Array, required: true, default: () => [] },
  title: { type: String, default: 'Tren Pelanggaran Siswa (6 Bulan Terakhir)' },
  badgeLabel: { type: String, default: '' },
})

const maxValue = computed(() => {
  if (!props.items.length) return 1
  return Math.max(...props.items.map((i) => i.value || i.jumlah || 0))
})

function trendHeight(value) {
  const max = maxValue.value
  return Math.max(8, Math.round((value / max) * 120))
}
</script>

<style scoped>
.content-card {
  border-radius: 8px;
}
.trend-chart {
  height: 170px;
  gap: 12px;
}
.trend-bar-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
}
.trend-bar {
  width: 100%;
  max-width: 32px;
  background: linear-gradient(180deg, #123a63 0%, #0a192f 100%);
  border-radius: 4px 4px 0 0;
  transition: height 0.4s ease;
}
</style>
