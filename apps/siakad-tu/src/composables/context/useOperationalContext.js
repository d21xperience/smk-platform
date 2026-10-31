import { computed } from 'vue';
import { useContextStore } from '../../stores/contextStore.js';
import { storeToRefs } from 'pinia';

export function useOperationalContext() {
  const contextStore = useContextStore();
  const { availablePeriods, selectedSemesterId, isLoading, error, selectedYearLabel, currentContext } = storeToRefs(contextStore);

  const selectedYear = computed(() => {
     const p = availablePeriods.value.find(x => x.semester_id === selectedSemesterId.value);
     return p ? p.tahun_ajaran_id : '';
  });
  const selectedSemester = computed(() => {
     const p = availablePeriods.value.find(x => x.semester_id === selectedSemesterId.value);
     return p ? p.semester : '';
  });
  const schoolLabel = computed(() => currentContext.value.schoolId);
  const hasYearChangeAuthority = computed(() => true); 
  
  const onYearChange = async (year, semester) => {
     const newSemesterId = `${year}${semester}`;
     await contextStore.setSemester(newSemesterId);
  };

  return {
    availablePeriods, selectedSemesterId, isLoading, error, selectedYearLabel, currentContext,
    // Map langsung ke currentContext untuk menghindari Proxy berlapis
    current: currentContext, 
    context: currentContext, 
    selectedYear, selectedSemester, schoolLabel, hasYearChangeAuthority, onYearChange,
    initialize: contextStore.initialize
  };
}

export const setOperationalContext = () => {};
