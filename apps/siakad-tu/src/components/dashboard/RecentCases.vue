<template>
  <q-card flat bordered class="content-card">
    <q-card-section>
      <div class="text-subtitle1 text-weight-bold text-primary">{{ title }}</div>
    </q-card-section>
    <q-separator />
    <q-list separator>
      <q-item
        v-for="item in items"
        :key="item.id"
        clickable
        @click="$emit('click:item', item)"
      >
        <q-item-section avatar>
          <q-avatar :color="item.avatarColor || 'blue-grey-6'" text-color="white" size="38px">
            {{ (item.kategori || 'R').charAt(0) }}
          </q-avatar>
        </q-item-section>
        <q-item-section>
          <q-item-label class="text-weight-medium">{{ item.judul || item.title }}</q-item-label>
          <q-item-label caption>
            {{ item.siswa || item.student }} •
            {{ item.kelas || item.class }} •
            {{ item.waktu || item.time }}
          </q-item-label>
        </q-item-section>
        <q-item-section side>
          <q-chip
            dense
            :color="item.chipColor || 'blue-grey-6'"
            text-color="white"
            class="text-caption"
          >
            {{ item.kategori || item.category }}
          </q-chip>
        </q-item-section>
      </q-item>
    </q-list>
  </q-card>
</template>

<script setup>
defineProps({
  items: { type: Array, required: true, default: () => [] },
  title: { type: String, default: 'Catatan Kasus Terbaru' },
})

defineEmits(['click:item'])
</script>

<style scoped>
.content-card {
  border-radius: 8px;
}
</style>
