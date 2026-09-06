<template>
  <div class="dashboard-header q-px-md q-py-lg">
    <div class="container max-width-center row items-center justify-between">
      <div>
        <div class="text-overline text-secondary text-weight-medium letter-spacing-2">
          {{ subtitle || 'SISTEM INFORMASI AKADEMIK' }}
        </div>
        <div class="text-h5 text-weight-bold text-white q-mt-xs">
          {{ title }}
        </div>
        <div class="text-caption text-grey-4 q-mt-xs">
          {{ dateLabel }}
        </div>
      </div>
      <div class="row items-center q-gutter-sm">
        <q-btn
          v-if="actionLabel"
          unelevated
          :color="actionColor || 'secondary'"
          text-color="dark"
          :label="actionLabel"
          no-caps
          :icon="actionIcon"
          @click="$emit('click:action')"
        />
        <q-btn round flat color="white" icon="notifications" @click="$emit('click:notification')">
          <q-badge color="secondary" text-color="dark" floating rounded>{{ notificationCount }}</q-badge>
        </q-btn>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, required: true },
  subtitle: { type: String, default: 'SISTEM INFORMASI AKADEMIK' },
  dateLabel: { type: String, default: '' },
  notificationCount: { type: [String, Number], default: 0 },
  actionLabel: { type: String, default: null },
  actionIcon: { type: String, default: 'add' },
  actionColor: { type: String, default: 'secondary' },
})

defineEmits(['click:action', 'click:notification'])
</script>

<style scoped>
.max-width-center {
  max-width: 1200px;
  margin: 0 auto;
}
.letter-spacing-2 {
  letter-spacing: 0.12em;
}
.dashboard-header {
  background: linear-gradient(135deg, #0a192f 0%, #123a63 100%);
}
</style>
