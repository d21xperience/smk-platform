<template>
  <q-card
    flat
    bordered
    class="kpi-card"
    :class="{ 'cursor-pointer hover-shadow': routeTo }"
    @click="handleClick"
  >
    <q-card-section class="row items-center no-wrap">
      <q-avatar :color="color" text-color="white" :icon="icon" size="46px" class="q-mr-md" />
      <div>
        <div class="text-h6 text-weight-bold text-primary">{{ value }}</div>
        <div class="text-caption text-grey-7">{{ label }}</div>
      </div>
    </q-card-section>
    <div class="kpi-trend q-px-md q-pb-sm" v-if="trend !== null && trend !== undefined">
      <q-icon :name="trend > 0 ? 'trending_up' : 'trending_down'"
        :color="trend > 0 ? 'positive' : 'negative'" size="16px" />
      <span class="text-caption q-ml-xs" :class="trend > 0 ? 'text-positive' : 'text-negative'">
        {{ Math.abs(trend) }}% dari bulan lalu
      </span>
    </div>
  </q-card>
</template>

<script setup>
import { useRouter } from 'vue-router'

const props = defineProps({
  label: { type: String, required: true },
  value: { type: [String, Number], required: true },
  icon: { type: String, required: true },
  color: { type: String, default: 'primary' },
  trend: { type: Number, default: null },
  routeTo: { type: [String, Object], default: null } // Bisa berupa string path atau object route
})

const router = useRouter()

const handleClick = () => {
  if (props.routeTo) {
    router.push(props.routeTo)
  }
}
</script>

<style scoped>
.kpi-card {
  border-radius: 8px;
  height: 100%;
  transition: box-shadow 0.3s ease, transform 0.2s ease;
}
.hover-shadow:hover {
  box-shadow: 0 4px 12px rgba(0,0,0,0.1) !important;
  transform: translateY(-2px);
}
.kpi-trend {
  display: flex;
  align-items: center;
}
</style>
