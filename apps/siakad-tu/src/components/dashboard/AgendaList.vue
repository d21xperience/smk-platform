<template>
  <q-card flat bordered class="content-card">
    <q-card-section>
      <div class="text-subtitle1 text-weight-bold text-primary">{{ title }}</div>
    </q-card-section>
    <q-separator />
    <q-list separator>
      <q-item
        v-for="item in items"
        :key="item.judul || item.title"
        clickable
        @click="$emit('click:item', item)"
      >
        <q-item-section avatar>
          <div class="agenda-date text-center">
            <div class="text-weight-bold text-primary">{{ item.tanggal || item.date }}</div>
            <div class="text-caption text-grey-7">{{ item.bulan || item.month }}</div>
          </div>
        </q-item-section>
        <q-item-section>
          <q-item-label class="text-weight-medium">{{ item.judul || item.title }}</q-item-label>
          <q-item-label caption>{{ item.waktu || item.time }}</q-item-label>
        </q-item-section>
      </q-item>
    </q-list>
  </q-card>
</template>

<script setup>
defineProps({
  items: { type: Array, required: true, default: () => [] },
  title: { type: String, default: 'Agenda Kesiswaan' },
})

defineEmits(['click:item'])
</script>

<style scoped>
.content-card {
  border-radius: 8px;
}
.agenda-date {
  min-width: 44px;
}
</style>
