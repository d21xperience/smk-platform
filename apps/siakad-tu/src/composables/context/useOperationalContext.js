import { computed } from 'vue';
import { useContextStore } from '../../stores/contextStore.js';
import { useAuthStore } from '../../stores/authStore.js';
import { storeToRefs } from 'pinia';

export function useOperationalContext() {
  const contextStore = useContextStore();
  const authStore = useAuthStore();
  
  const { availablePeriods, selectedSemesterId, isLoading, error, selectedYearLabel, currentContext } = storeToRefs(contextStore);

  const selectedYear = computed(() => {
     const p = availablePeriods.value.find(x => x.semester_id === selectedSemesterId.value);
     return p ? p.tahun_ajaran_id : '';
  });
  
  const selectedSemester = computed(() => {
     const p = availablePeriods.value.find(x => x.semester_id === selectedSemesterId.value);
     return p ? p.semester : '';
  });
  
  // PERBAIKAN: Ambil school_name dari authStore, fallback ke schoolId jika kosong
  const schoolLabel = computed(() => authStore.user?.school_name || currentContext.value.schoolId);
  
  const hasYearChangeAuthority = computed(() => true); 
  
  const onYearChange = async (year, semester) => {
     const newSemesterId = `${year}${semester}`;
     await contextStore.setSemester(newSemesterId);
  };

  return {
    availablePeriods, selectedSemesterId, isLoading, error, selectedYearLabel, currentContext,
    current: currentContext, 
    context: currentContext, 
    selectedYear, selectedSemester, schoolLabel, hasYearChangeAuthority, onYearChange,
    initialize: contextStore.initialize
  };
}

export const setOperationalContext = () => {};
