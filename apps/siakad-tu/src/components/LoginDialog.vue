<!-- @/components/LoginDialog.vue -->
<template>
  <q-dialog v-model="isVisible" persistent @hide="handleDialogHide">
    <q-card style="min-width: 350px; max-width: 400px;">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6 text-primary">Login SDP SMK</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup aria-label="Tutup" />
      </q-card-section>
      <!-- ✨ BANNER INFORMASI KONTEKSTUAL ✨ -->
      <q-card-section v-if="authStore.loginPromptMessage" class="q-pt-md">
        <q-banner dense rounded class="bg-info text-white">
          <template v-slot:avatar>
            <q-icon name="info" color="white" />
          </template>
          <div class="text-body2">{{ authStore.loginPromptMessage }}</div>
        </q-banner>
      </q-card-section>
      <q-card-section>
        <q-form @submit="submitLogin" class="q-gutter-md">
          <q-input v-model="form.username" label="Username" outlined dense autocomplete="username"
            :rules="[val => !!val || 'Field wajib diisi']">
            <template v-slot:prepend>
              <q-icon name="person" />
            </template>
          </q-input>

          <q-input v-model="form.password" label="Password" outlined dense type="password"
            autocomplete="current-password" :rules="[val => !!val || 'Field wajib diisi']">
            <template v-slot:prepend>
              <q-icon name="lock" />
            </template>
          </q-input>

          <div class="row justify-end q-mt-md">
            <q-btn label="Batal" color="grey" flat v-close-popup class="q-mr-sm" />
            <q-btn label="Masuk" type="submit" color="primary" :loading="isLoading" unelevated />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useQuasar } from 'quasar'
import { useAuthStore } from '@/stores/authStore'
import { useRouter } from 'vue-router'
const router = useRouter()
// Props & Emits untuk v-model support
const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true
  }
})

const emit = defineEmits(['update:modelValue', 'success'])

const $q = useQuasar()
const authStore = useAuthStore()

// Computed property untuk sinkronisasi v-model dengan parent
const isVisible = computed({
  get: () => props.modelValue,
  set: (val) => {
    emit('update:modelValue', val)
    if (!val) {
      authStore.closeLoginDialog()
    }
  }
})

const isLoading = ref(false)
const form = reactive({
  username: '',
  password: ''
})

const resetForm = () => {
  form.username = ''
  form.password = ''
  isLoading.value = false
}

const submitLogin = async () => {
  isLoading.value = true
  try {
    const success = await authStore.login(form.username, form.password)

    if (success) {
      $q.notify({ color: 'positive', message: 'Login berhasil!', icon: 'check_circle', position: 'top' })
      console.log(authStore.redirectUrl)
      // 1. Prioritaskan `redirectUrl` (jika user sebelumnya ditolak akses ke halaman proteksi)
      // 2. Jika `null`, gunakan halaman saat ini (`fullPath`) agar user tetap di halaman publik yang sama
      const target = authStore.redirectUrl || router.currentRoute.value.fullPath
      authStore.closeLoginDialog() // Ini akan menutup dialog DAN membersihkan redirectUrl

      router.push(target)
      resetForm()
    } else {
      $q.notify({ color: 'negative', message: 'Username atau password salah.', icon: 'error', position: 'top' })
    }
  } catch (error) {
    console.log(error)
    $q.notify({ color: 'negative', message: 'Terjadi kesalahan sistem.', icon: 'error', position: 'top' })
  } finally {
    isLoading.value = false
  }
}

// Handler tambahan untuk event @hide (backup safety)
const handleDialogHide = () => {
  // Dipastikan state selalu bersih saat dialog benar-benar hilang
  resetForm()
  if (authStore.showLoginDialog) {
    authStore.closeLoginDialog()
  }
}
</script>
