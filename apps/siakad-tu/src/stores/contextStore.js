import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { ContextService } from '../../../../../packages/context/src/services/ContextService.js';

export const useContextStore = defineStore('context', () => {
  const availablePeriods = ref([]);
  const selectedSemesterId = ref('');
  const userId = ref('user-123'); 
  const schoolId = ref('school-abc'); 
  const isLoading = ref(false);
  const error = ref(null);

  const currentContext = computed(() => ({
    userId: userId.value,
    schoolId: schoolId.value,
    semesterId: selectedSemesterId.value
  }));

  const selectedYearLabel = computed(() => {
    const selected = availablePeriods.value.find(p => p.semester_id === selectedSemesterId.value);
    return selected ? `${selected.tahun_ajaran_id}/${parseInt(selected.tahun_ajaran_id) + 1} - ${selected.nama_semester}` : '-';
  });

  async function initialize() {
    isLoading.value = true;
    error.value = null;
    try {
      const ctx = currentContext.value;
      const cacheKey = `ctx:${ctx.userId}:${ctx.schoolId}`;
      
      availablePeriods.value = await ContextService.fetchPeriods(ctx);
      const pref = await ContextService.fetchPreference(ctx);
      
      if (pref && pref.semester_id) {
        selectedSemesterId.value = pref.semester_id;
        localStorage.setItem(cacheKey, pref.semester_id);
      } else {
        const active = availablePeriods.value.find(p => p.periode_aktif === '1');
        selectedSemesterId.value = active ? active.semester_id : (availablePeriods.value[0]?.semester_id || '');
      }
    } catch (err) {
      error.value = err.message;
    } finally {
      isLoading.value = false;
    }
  }

  async function setSemester(semesterId) {
    isLoading.value = true;
    try {
      const ctx = currentContext.value;
      await ContextService.updatePreference(ctx, semesterId);
      selectedSemesterId.value = semesterId;
      localStorage.setItem(`ctx:${ctx.userId}:${ctx.schoolId}`, semesterId);
    } catch (err) {
      error.value = err.message;
      throw err;
    } finally {
      isLoading.value = false;
    }
  }

  return { availablePeriods, selectedSemesterId, isLoading, error, currentContext, selectedYearLabel, initialize, setSemester };
});
