<template>
  <q-layout view="lHh Lpr lFf">
    <q-page-container>
      <q-page class="flex flex-center bg-primary">
        <q-card style="width: 420px; max-width: 90vw">
          <q-card-section class="text-center q-pb-none">
            <div class="text-h5">SIAKAD Guru</div>
            <div class="text-subtitle2 text-grey">School Digital Platform</div>
          </q-card-section>

          <q-card-section>
            <q-banner v-if="isExpired" class="bg-warning text-white q-mb-md">
              <template #avatar>
                <q-icon name="warning" />
              </template>
              Sesi Anda telah berakhir. Silakan login kembali.
            </q-banner>

            <q-banner v-if="error" class="bg-negative text-white q-mb-md">
              <template #avatar>
                <q-icon name="error" />
              </template>
              {{ error }}
            </q-banner>

            <q-form @submit="handleLogin" class="q-gutter-md">
              <q-input v-model="username" label="Username" outlined dense :disable="isLoading">
                <template #prepend>
                  <q-icon name="person" />
                </template>
              </q-input>

              <q-input v-model="password" label="Password" type="password" outlined dense :disable="isLoading">
                <template #prepend>
                  <q-icon name="lock" />
                </template>
              </q-input>

              <q-btn type="submit" label="Login" color="primary" :loading="isLoading" unelevated class="full-width" />
            </q-form>

            <q-separator class="q-my-md" />

            <div class="text-center">
              <q-btn v-if="isPasskeyAvailable" flat color="secondary" icon="fingerprint" label="Login dengan Passkey"
                :loading="isLoading" @click="handlePasskeyLogin" />
              <div v-else class="text-caption text-grey">
                Passkey tidak tersedia di browser ini
              </div>
            </div>
          </q-card-section>

          <q-card-section class="text-center text-caption text-grey">
            <p class="q-mb-none">Demo Accounts:</p>
            <p class="q-mb-none">Username: guru | Password: password</p>
            <p>Username: wali | Password: password</p>
          </q-card-section>
        </q-card>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref } from 'vue'
import { useAuth } from '../../composables/useAuth.js'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const {
  login,
  loginWithPasskey,
  isLoading,
  error,
  isExpired,
  isPasskeyAvailable
} = useAuth()

const username = ref('')
const password = ref('')

const handleLogin = async () => {
  if (!username.value || !password.value) {
    $q.notify({
      type: 'warning',
      message: 'Username dan password wajib diisi'
    })
    return
  }

  try {
    await login(username.value, password.value)
  } catch (err) {
    console.log(err)
    // Error handled by store
  }
}

const handlePasskeyLogin = async () => {
  try {
    await loginWithPasskey('mock-challenge-' + Date.now())
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: err.message
    })
  }
}
</script>
