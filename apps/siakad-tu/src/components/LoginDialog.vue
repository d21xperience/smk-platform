<template>
  <q-dialog v-model="isVisible" persistent>
    <q-card style="min-width: 350px; max-width: 400px;">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6 text-primary">Login SDP SMK</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup aria-label="Tutup" />
      </q-card-section>

      <q-card-section>
        <q-form @submit="submitLogin" class="q-gutter-md">
          <q-input v-model="form.username" label="Username" outlined dense autocomplete="username"
            :rules="[val => !!val || 'Username wajib diisi']">
            <template v-slot:prepend>
              <q-icon name="person" />
            </template>
          </q-input>

          <q-input v-model="form.password" label="Password" outlined dense type="password"
            autocomplete="current-password" :rules="[val => !!val || 'Password wajib diisi']">
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
const $q = useQuasar()
const authStore = useAuthStore()

const props = defineProps({
  modelValue: { type: Boolean, required: true }
})
const emit = defineEmits(['update:modelValue', 'success'])

const isVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const isLoading = ref(false)
const form = reactive({
  username: '',
  password: ''
})

const submitLogin = async () => {
  isLoading.value = true
  try {
    // Pastikan mengirim objek, bukan dua argumen terpisah
    const success = await authStore.login({
      username: form.username,
      password: form.password
    })

    if (success) {
      $q.notify({ color: 'positive', message: 'Login berhasil!', icon: 'check_circle', position: 'top' })

      // Paksa tutup dialog
      isVisible.value = false

      // Redirect ke halaman yang dituju
      const target = authStore.redirectUrl || router.currentRoute.value.fullPath
      authStore.closeLoginDialog()
      router.push(target)

      // Bersihkan form
      form.username = ''
      form.password = ''
    } else {
      $q.notify({
        color: 'negative',
        message: authStore.error?.message || 'Username atau password salah.',
        icon: 'error',
        position: 'top'
      })
    }
  } catch (error) {
    $q.notify({ color: 'negative', message: `Terjadi kesalahan sistem. ${error}`, icon: 'error', position: 'top' })
  } finally {
    isLoading.value = false
  }
}
</script>
