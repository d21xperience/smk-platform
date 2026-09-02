<template>
  <q-page class="q-pa-md">
    <!-- Header dengan sapaan -->
    <div class="text-h5 q-mb-md">Dashboard {{ roleName }}</div>
    <div class="text-subtitle1">Selamat datang, {{ user.nama }}</div>

    <!-- Konten berdasarkan role -->
    <div v-if="role === 'admin' || role === 'kepala_sekolah'">
      <DashboardKepsek />
    </div>
    <div v-else-if="role === 'guru'">
      <DashboardGuru />
    </div>
    <div v-else-if="role === 'wali_kelas'">
      <DashboardWaliKelas />
    </div>
    <div v-else-if="role === 'operator'">
      <DashboardOperator />
    </div>
    <div v-else>
      <div class="text-center">Dashboard tidak tersedia untuk role ini</div>
    </div>
  </q-page>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from 'stores/auth'
import DashboardKepsek from 'components/dashboard/DashboardKepsek.vue'
import DashboardGuru from 'components/dashboard/DashboardGuru.vue'
import DashboardWaliKelas from 'components/dashboard/DashboardWaliKelas.vue'
import DashboardOperator from 'components/dashboard/DashboardOperator.vue'

const authStore = useAuthStore()
const user = computed(() => authStore.user || { nama: 'User', role: 'guru' })
const role = computed(() => user.value.role)
const roleName = computed(() => {
  const map = { admin: 'Admin', kepala_sekolah: 'Kepala Sekolah', guru: 'Guru', wali_kelas: 'Wali Kelas', operator: 'Operator' }
  return map[role.value] || role.value
})
</script>
