<template>
  <div>
    <q-btn flat no-caps icon="calendar_today" :label="selectedYearLabel" @click="openDialog" />
    <q-dialog v-model="showDialog">
      <q-card style="min-width: 350px">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">Pilih Tahun Ajaran</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section>
          <q-select v-model="tempYear" :options="yearOptions" label="Tahun Ajaran" dense outlined emit-value map-options class="q-mb-md" />
          <q-select v-model="tempSemester" :options="semesterOptions" label="Semester" dense outlined emit-value map-options />
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="BATAL" v-close-popup />
          <q-btn color="primary" label="SIMPAN" :loading="isLoading" @click="save" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useContextStore } from '../stores/contextStore.js';

const contextStore = useContextStore();
const { availablePeriods, selectedSemesterId, isLoading, selectedYearLabel } = storeToRefs(contextStore);

const showDialog = ref(false);
const tempYear = ref('');
const tempSemester = ref('');

onMounted(() => {
  if (availablePeriods.value.length === 0) contextStore.initialize();
});

const yearOptions = computed(() => {
  const years = [...new Set(availablePeriods.value.map(p => p.tahun_ajaran_id))];
  return years.map(y => ({ label: `${y}/${parseInt(y) + 1}`, value: y }));
});

const semesterOptions = [
  { label: 'Semester 1', value: '1' },
  { label: 'Semester 2', value: '2' }
];

const openDialog = () => {
  const current = availablePeriods.value.find(p => p.semester_id === selectedSemesterId.value);
  if (current) {
    tempYear.value = current.tahun_ajaran_id;
    tempSemester.value = current.semester;
  }
  showDialog.value = true;
};

const save = async () => {
  if (!tempYear.value || !tempSemester.value) return;
  const newSemesterId = `${tempYear.value}${tempSemester.value}`;
  try {
    await contextStore.setSemester(newSemesterId);
    showDialog.value = false;
  } catch (err) {
    console.error(err);
  }
};
</script>
