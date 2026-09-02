<template>
  <div class="login-page">
    <q-card class="login-card">
      <q-card-section class="text-center">
        <h5>SIAKAD Guru</h5>
        <p class="text-grey-6">Masuk ke platform sekolah</p>
      </q-card-section>

      <q-card-section>
        <q-form @submit.prevent="handleLogin" class="q-gutter-y-md">
          <q-input v-model="email" label="Email" type="email" outlined dense :error="!!validationErrors.email"
            :error-message="validationErrors.email" autofocus />
          <q-input v-model="password" label="Password" type="password" outlined dense
            :error="!!validationErrors.password" :error-message="validationErrors.password"
            @keyup.enter="handleLogin" />

          <div v-if="loginError" class="text-negative text-caption">
            {{ loginError }}
          </div>

          <q-btn type="submit" label="Masuk" color="primary" class="full-width" :loading="loading" :disable="loading"
            size="lg" />
        </q-form>
      </q-card-section>

      <q-card-section class="text-center text-caption text-grey-6">
        <p>Gunakan <strong>guru@smk.sch.id</strong> / <strong>password123</strong></p>
      </q-card-section>
    </q-card>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth.store';
// import { useQuasar } from 'quasar';

const router = useRouter();
// const $q = useQuasar();
const authStore = useAuthStore();

const email = ref('');
const password = ref('');
const loginError = ref('');
const validationErrors = ref({});
const loading = ref(false);

const handleLogin = async () => {
  // Reset error
  loginError.value = '';
  validationErrors.value = {};
  loading.value = true;

  try {
    await authStore.login(email.value, password.value);
    // Login sukses, redirect ke select context
    router.push('/select-context');
  } catch (err) {
    // Coba parsing error dari Engine (validasi)
    if (err.errors && Array.isArray(err.errors)) {
      err.errors.forEach(e => {
        if (e.field === 'email') validationErrors.value.email = e.message;
        if (e.field === 'password') validationErrors.value.password = e.message;
      });
      if (!validationErrors.value.email && !validationErrors.value.password) {
        loginError.value = err.message;
      }
    } else {
      loginError.value = err.message || 'Login gagal';
    }
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: #f5f7fa;
}

.login-card {
  width: 400px;
  max-width: 90%;
  padding: 16px;
}
</style>
