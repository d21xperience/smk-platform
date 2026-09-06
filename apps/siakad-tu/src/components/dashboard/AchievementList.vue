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
          <q-avatar color="amber-2" text-color="amber-9" icon="emoji_events" size="38px" />
        </q-item-section>
        <q-item-section>
          <q-item-label class="text-weight-medium">{{ item.judul || item.title }}</q-item-label>
          <q-item-label caption>{{ item.siswa || item.student }} • {{ item.tanggal || item.date }}</q-item-label>
        </q-item-section>
      </q-item>
    </q-list>
  </q-card>
</template>

<script setup>
defineProps({
  items: { type: Array, required: true, default: () => [] },
  title: { type: String, default: 'Prestasi Terbaru' },
})

defineEmits(['click:item'])
</script>

<style scoped>
.content-card {
  border-radius: 8px;
}
</style>
