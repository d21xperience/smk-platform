<template>
  <q-drawer v-model="model" side="left" overlay behavior="mobile">
    <q-list>
      <q-item clickable v-ripple to="/assessment">
        <q-item-section avatar>
          <q-icon name="grading" />
        </q-item-section>
        <q-item-section>Penilaian</q-item-section>
      </q-item>
      <q-item clickable v-ripple to="/reporting">
        <q-item-section avatar>
          <q-icon name="bar_chart" />
        </q-item-section>
        <q-item-section>Laporan</q-item-section>
      </q-item>
      <q-item clickable v-ripple to="/inventory">
        <q-item-section avatar>
          <q-icon name="inventory" />
        </q-item-section>
        <q-item-section>Inventaris</q-item-section>
      </q-item>
      <q-item clickable v-ripple to="/finance">
        <q-item-section avatar>
          <q-icon name="payments" />
        </q-item-section>
        <q-item-section>Keuangan Siswa</q-item-section>
      </q-item>
    </q-list>
  </q-drawer>
</template>
<script setup>
const model = defineModel({ default: false });
</script>
