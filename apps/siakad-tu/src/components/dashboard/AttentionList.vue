<template>
  <q-card flat bordered class="content-card border-alert">
    <q-card-section>
      <div class="row items-center">
        <q-icon :name="icon || 'report'" :color="color || 'red'" size="22px" class="q-mr-sm" />
        <div class="text-subtitle1 text-weight-bold" :class="color ? `text-${color}-9` : 'text-red-9'">
          {{ title || 'Siswa Memerlukan Perhatian Khusus' }}
        </div>
      </div>
    </q-card-section>
    <q-separator />
    <q-list separator>
      <q-item v-for="item in items" :key="item.nama || item.name" clickable @click="$emit('click:item', item)">
        <q-item-section avatar>
          <q-avatar :color="color ? `${color}-1` : 'red-1'" text-color="red-9" icon="person" size="38px" />
        </q-item-section>
        <q-item-section>
          <q-item-label class="text-weight-medium">{{ item.nama || item.name }}</q-item-label>
          <q-item-label caption>{{ item.kelas || item.class }} • {{ item.catatan || item.note }}</q-item-label>
        </q-item-section>
        <q-item-section side>
          <q-badge :color="color || 'red'" text-color="white" class="q-px-sm">
            {{ item.poin || item.points }} poin
          </q-badge>
        </q-item-section>
      </q-item>
    </q-list>
  </q-card>
</template>

<script setup>
defineProps({
  items: { type: Array, required: true, default: () => [] },
  title: { type: String, default: 'Siswa Memerlukan Perhatian Khusus' },
  icon: { type: String, default: 'report' },
  color: { type: String, default: 'red' },
})

defineEmits(['click:item'])
</script>

<style scoped>
.content-card {
  border-radius: 8px;
}

.border-alert {
  border-left: 4px solid #c62828;
}
</style>
