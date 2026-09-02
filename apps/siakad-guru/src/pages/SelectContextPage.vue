<template>
  <div class="select-context-page">
    <q-card class="context-card">
      <q-card-section class="text-center">
        <h6>Pilih Tahun Ajaran</h6>
        <p class="text-grey-6">Tentukan tahun ajaran dan semester aktif</p>
      </q-card-section>

      <q-card-section>
        <q-form @submit.prevent="handleSelectContext" class="q-gutter-y-md">
          <q-select v-model="selectedYear" :options="availableYears" label="Tahun Ajaran" option-label="code"
            option-value="id" map-options emit-value outlined dense :loading="isLoading" />

          <q-select v-model="selectedSemester" :options="filteredSemesters" label="Semester" option-label="name"
            option-value="id" map-options emit-value outlined dense :disable="!selectedYear" :loading="isLoading" />

          <q-btn type="submit" label="Lanjutkan" color="primary" class="full-width"
            :disable="!selectedYear || !selectedSemester || isLoading" :loading="isLoading" />
        </q-form>
      </q-card-section>

      <q-card-section class="text-center text-caption text-grey-6">
        <q-btn flat dense @click="logoutHandler" label="Keluar" />
      </q-card-section>
    </q-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useOperationalContext } from '@/composables/context/useOperationalContext';
import { useAuth } from '@/composables/useAuth';

const router = useRouter();
const { logout } = useAuth();
const {
  isContextReady, // ✅ Sekarang tersedia
  availableYears,
  availableSemesters,
  isLoading,
  loadContext,
  switchContext,
} = useOperationalContext();

const selectedYear = ref(null);
const selectedSemester = ref(null);

const filteredSemesters = computed(() => {
  if (!selectedYear.value) return [];
  return availableSemesters.value.filter((s) => s.yearId === selectedYear.value);
});

onMounted(async () => {
  // ✅ Gunakan isContextReady
  if (isContextReady.value) {
    router.push('/dashboard');
    return;
  }
  if (availableYears.value.length === 0) {
    await loadContext();
  }
  console.log('Available years:', availableYears.value);
  console.log('Available semesters:', availableSemesters.value);
});

const handleSelectContext = async () => {
  if (!selectedYear.value || !selectedSemester.value) return;
  try {
    await switchContext({
      academicYearId: selectedYear.value,
      semesterId: selectedSemester.value,
    });
    router.push('/dashboard');
  } catch (err) {
    console.error('Gagal switch context:', err);
  }
};

const logoutHandler = () => {
  logout();
  router.push('/login');
};
</script>

<style scoped>
.select-context-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: #f5f7fa;
}

.context-card {
  width: 450px;
  max-width: 90%;
  padding: 16px;
}
</style>
