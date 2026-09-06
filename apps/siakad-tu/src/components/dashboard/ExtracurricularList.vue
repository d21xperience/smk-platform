<template>
  <q-card flat bordered class="content-card">
    <q-card-section>
      <div class="text-subtitle1 text-weight-bold text-primary">{{ title }}</div>
    </q-card-section>
    <q-separator />
    <q-list separator>
      <q-item
        v-for="item in items"
        :key="item.nama || item.name"
        clickable
        @click="$emit('click:item', item)"
      >
        <q-item-section avatar>
          <q-avatar :color="item.warna || item.color" text-color="white" :icon="item.icon" size="38px" />
        </q-item-section>
        <q-item-section>
          <q-item-label class="text-weight-medium">{{ item.nama || item.name }}</q-item-label>
          <q-item-label caption>Pembina: {{ item.pembina || item.coach }}</q-item-label>
        </q-item-section>
        <q-item-section side>
          <div class="text-weight-bold text-primary">{{ item.peserta || item.participants }}</div>
        </q-item-section>
      </q-item>
    </q-list>
  </q-card>
</template>

<script setup>
defineProps({
  items: { type: Array, required: true, default: () => [] },
  title: { type: String, default: 'Peserta Ekstrakurikuler' },
})

defineEmits(['click:item'])
</script>

<style scoped>
.content-card {
  border-radius: 8px;
}
</style>
