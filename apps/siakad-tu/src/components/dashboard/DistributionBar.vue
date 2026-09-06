<template>
  <q-card flat bordered class="content-card">
    <q-card-section>
      <div class="text-subtitle1 text-weight-bold text-primary">{{ title }}</div>
    </q-card-section>
    <q-separator />
    <q-card-section>
      <div v-for="item in items" :key="item.label" class="q-mb-md">
        <div class="row items-center justify-between q-mb-xs">
          <div class="text-body2 text-weight-medium">{{ item.label }}</div>
          <div class="text-caption text-grey-7">{{ item.total }} poin akumulasi</div>
        </div>
        <div class="bar-track">
          <div
            class="bar-fill"
            :style="{
              width: barWidth(item.total) + '%',
              backgroundColor: item.color || '#0a192f'
            }"
          />
        </div>
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  items: { type: Array, required: true, default: () => [] },
  title: { type: String, default: 'Distribusi Poin Pelanggaran per Tingkat Kelas' },
})

const maxValue = computed(() => {
  if (!props.items.length) return 1
  return Math.max(...props.items.map((i) => i.total || 0))
})

function barWidth(value) {
  const max = maxValue.value
  return Math.round((value / max) * 100)
}
</script>

<style scoped>
.content-card {
  border-radius: 8px;
}
.bar-track {
  width: 100%;
  height: 10px;
  border-radius: 6px;
  background-color: #eceff1;
  overflow: hidden;
}
.bar-fill {
  height: 100%;
  border-radius: 6px;
  transition: width 0.4s ease;
}
</style>
