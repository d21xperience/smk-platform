<template>
  <q-header elevated class="bg-primary text-white">
    <q-toolbar>
      <q-btn dense flat round icon="menu" @click="$emit('toggle-sidebar')" />
      <q-toolbar-title>
        SIAKAD Guru
        <span v-if="context" class="text-caption text-grey-3 q-ml-sm">
          {{ context.academicYearCode }} - Semester {{ context.semesterName }}
        </span>
      </q-toolbar-title>

      <div class="row items-center q-gutter-sm">
        <q-badge color="white" text-color="primary" v-if="user">
          {{ user.name }} ({{ user.role }})
        </q-badge>
        {{ contextInfo }}
        <q-btn flat dense round icon="logout" @click="handleLogout" />
      </div>
    </q-toolbar>
  </q-header>
</template>

<script setup>
// import { computed } from 'vue';
import { useAuth } from '@/composables/useAuth';
import { useOperationalContext } from '@/composables/context/useOperationalContext';
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar'
const router = useRouter()
// const emit = defineEmits(['toggle-sidebar']);
const { logout, user } = useAuth();
const { context } = useOperationalContext();
const $q = useQuasar()
const contextInfo = computed(() => {
  if (!context.value) return null;
  console.log("APP HEADER 😁", context.value)

  return {
    academicYearCode: context.value.academicYearCode || '--',
    semesterName: context.value.semesterId || '--',
  };
});

const handleLogout = () => {
  $q.dialog({
    title: 'Konfirmasi',
    message: 'Apakah Anda yakin ingin keluar dari SIAKAD?',
    cancel: true,
    persistent: true,
  }).onOk(() => {
    logout()
    router.push('/login')
  })

}
</script>
