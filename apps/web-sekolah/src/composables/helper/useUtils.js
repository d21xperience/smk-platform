import { computed } from 'vue'

export function useUtils() {
  const currentYear = computed(() => new Date().getFullYear())
  const greetingTime = computed(() => {
    const hour = new Date().getHours()
    if (hour < 11) return 'Pagi'
    if (hour < 15) return 'Siang'
    if (hour < 18) return 'Sore'
    return 'Malam'
  })
  const todayLabel = computed(() => {
    return new Date().toLocaleDateString('id-ID', {
      weekday: 'long',
      day: 'numeric',
      month: 'long',
      year: 'numeric',
    })
  })
  // const roleLabel = computed((role) => {
  //   const roleMap = {
  //     guru: 'Guru Mata Pelajaran',
  //     wali_kelas: 'Wali Kelas',
  //     admin: 'Administrator SIAKAD',
  //   }
  //   return roleMap[role] || 'Guru'
  // })
  return {
    greetingTime,
    currentYear,
    todayLabel,
    // roleLabel,
  }
}
