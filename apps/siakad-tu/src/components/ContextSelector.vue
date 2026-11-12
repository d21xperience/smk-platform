<template>
  <div>
    <!-- TOMBOL SEMESTER DENGAN GLOW JIKA PSB AKTIF -->
    <q-btn 
      flat no-caps no-wrap icon="calendar_today" 
      :label="selectedYearLabel" 
      @click="openDialog"
      :class="{'psb-glow': isPsbActive}"
    >
      <q-badge v-if="isPsbActive" color="orange" floating rounded>
        <q-icon name="campaign" size="12px" />
      </q-badge>
    </q-btn>

    <q-dialog v-model="showDialog">
      <q-card style="min-width: 350px">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6">Pilih Tahun Ajaran</div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>
        <q-card-section>
          <q-select v-model="tempYear" :options="yearOptions" label="Tahun Ajaran" dense outlined emit-value map-options class="q-mb-md" />
          <q-select v-model="tempSemester" :options="semesterOptions" label="Semester" dense outlined emit-value map-options class="q-mb-md" />
          
          <!-- TOGGLE PSB (Hanya muncul jika Semester 1 / Ganjil) -->
          <q-toggle
            v-if="tempSemester === '1'"
            v-model="tempIsPsbActive"
            label="Aktifkan Masa PSB (Penerimaan Siswa Baru)"
            color="orange"
          />
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
import { api } from '@/boot/axios';

const contextStore = useContextStore();
const { availablePeriods, selectedSemesterId, isLoading, selectedYearLabel } = storeToRefs(contextStore);

const showDialog = ref(false);
const tempYear = ref('');
const tempSemester = ref('');
const tempIsPsbActive = ref(false);

const isPsbActive = computed(() => {
  const current = availablePeriods.value.find(p => String(p.semester_id) === String(selectedSemesterId.value));
  return current?.is_psb_active || false;
});

onMounted(() => {
  if (availablePeriods.value.length === 0) contextStore.initialize();
});

const yearOptions = computed(() => {
  const years = [...new Set(availablePeriods.value.map(p => String(p.tahun_ajaran_id)))];
  return years.map(y => ({ label: `${y}/${parseInt(y) + 1}`, value: y }));
});

const semesterOptions = [
  { label: 'Semester 1 (Ganjil)', value: '1' },
  { label: 'Semester 2 (Genap)', value: '2' }
];

const openDialog = () => {
  const current = availablePeriods.value.find(p => String(p.semester_id) === String(selectedSemesterId.value));
  if (current) {
    tempYear.value = String(current.tahun_ajaran_id);
    tempSemester.value = current.semester_id ? String(current.semester_id).slice(-1) : '1';
    tempIsPsbActive.value = current.is_psb_active || false;
  }
  showDialog.value = true;
};

const save = async () => {
  if (!tempYear.value || !tempSemester.value) return;
  const newSemesterId = `${tempYear.value}${tempSemester.value}`;

  try {
    // Simpan PSB Toggle ke Backend (contoh endpoint)
    await api.post('/context/toggle-psb', { 
      semester_id: newSemesterId, 
      is_active: tempIsPsbActive.value 
    });
    
    await contextStore.setSemester(newSemesterId);
    await contextStore.initialize(); // Refresh data is_psb_active
    showDialog.value = false;
  } catch (err) {
    console.error('[ContextSelector] Save error:', err);
  }
};
</script>

<style scoped>
.psb-glow {
  animation: glow 1.5s infinite alternate;
  border-radius: 8px;
}
@keyframes glow {
  from { box-shadow: 0 0 5px #ff9800; }
  to { box-shadow: 0 0 15px #ff9800, 0 0 20px #ff9800; }
}
</style>
