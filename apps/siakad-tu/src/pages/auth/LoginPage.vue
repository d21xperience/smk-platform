<template>
  <q-page class="flex flex-center bg-grey-2">
    <q-card style="width: 420px; max-width: 90vw">
      <!-- Header -->
      <q-card-section class="bg-primary text-white text-center">
        <q-icon name="school" size="3em" class="q-mb-sm" />
        <div class="text-h5 text-weight-bold">School Digital Platform</div>
        <div class="text-caption">Sistem Informasi Manajemen Sekolah</div>
      </q-card-section>

      <!-- Form -->
      <q-card-section class="q-pt-lg">
        <div class="text-h6 text-center q-mb-md">Login Pegawai TU</div>

        <!-- Redirect Info -->
        <q-banner v-if="redirectTo" rounded class="bg-info text-white q-mb-md">
          <template v-slot:avatar>
            <q-icon name="info" color="white" />
          </template>
          Silakan login untuk mengakses halaman yang dituju.
        </q-banner>

        <!-- Error Banner -->
        <q-banner v-if="error" rounded class="bg-negative text-white q-mb-md">
          <template v-slot:avatar>
            <q-icon name="error" color="white" />
          </template>
          {{ error.message }}
        </q-banner>

        <q-form @submit="handleLogin">
          <q-input v-model="formData.username" label="Username / NIP" outlined
            :rules="[val => !!val || 'Username wajib diisi']" class="q-mb-md">
            <template v-slot:prepend>
              <q-icon name="person" />
            </template>
          </q-input>

          <q-input v-model="formData.password" :type="showPassword ? 'text' : 'password'" label="Password" outlined
            :rules="[val => !!val || 'Password wajib diisi']" class="q-mb-md">
            <template v-slot:prepend>
              <q-icon name="lock" />
            </template>
            <template v-slot:append>
              <q-icon :name="showPassword ? 'visibility_off' : 'visibility'" class="cursor-pointer"
                @click="togglePassword" />
            </template>
          </q-input>

          <q-btn type="submit" color="primary" label="Login" class="full-width q-mt-md" size="lg" :loading="isLoading"
            :disable="!isFormValid" />
        </q-form>
      </q-card-section>

      <!-- Footer -->
      <q-card-section class="text-center text-grey-7">
        <q-btn flat label="Kembali ke Beranda" icon="arrow_back" @click="goToLanding" />
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { useAuthLogin } from '@/composables/auth/useAuthLogin'

const {
  formData,
  showPassword,
  isLoading,
  error,
  redirectTo,
  isFormValid,
  handleLogin,
  togglePassword,
  goToLanding,
} = useAuthLogin()
</script>
