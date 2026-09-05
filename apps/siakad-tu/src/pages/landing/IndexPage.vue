<template>
  <q-page padding>
    <!-- Hero Section -->
    <div class="row items-center justify-center q-mb-xl">
      <div class="col-12 text-center">
        <div class="text-h3 text-weight-bold q-mb-md">
          School Digital Platform
        </div>
        <div class="text-h6 text-grey-7 q-mb-lg">
          Sistem Informasi Manajemen Sekolah Terintegrasi
        </div>

        <div v-if="isMaintenanceMode" class="q-mb-lg">
          <q-banner rounded class="bg-warning text-white">
            <template v-slot:avatar>
              <q-icon name="warning" color="white" />
            </template>
            Sistem sedang dalam pemeliharaan. Silakan coba lagi nanti.
          </q-banner>
        </div>

        <div class="row justify-center q-gutter-md">
          <q-btn color="primary" size="lg" label="Masuk ke Sistem" icon="login" @click="goToLogin" />
          <q-btn outline color="primary" size="lg" label="Daftar Akun" icon="person_add" @click="goToRegister" />
        </div>
      </div>
    </div>

    <!-- Features Section -->
    <div class="q-mb-xl">
      <div class="text-h5 text-weight-bold text-center q-mb-lg">
        Fitur Unggulan
      </div>

      <div class="row q-gutter-md">
        <div v-for="(feature, index) in features" :key="index" class="col-12 col-sm-6 col-md-4">
          <q-card class="fit">
            <q-card-section class="text-center">
              <q-icon :name="feature.icon" :color="feature.color" size="4em" class="q-mb-md" />
              <div class="text-h6 q-mb-sm">{{ feature.title }}</div>
              <div class="text-grey-7">{{ feature.description }}</div>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>

    <!-- Modules Section -->
    <div class="q-mb-xl">
      <div class="text-h5 text-weight-bold text-center q-mb-lg">
        Modul Tersedia
      </div>

      <div class="row q-gutter-md justify-center">
        <div v-for="(module, index) in modules" :key="index" class="col-6 col-sm-4 col-md-3 col-lg-2">
          <q-card class="fit cursor-pointer" @click="navigateTo(module.path)">
            <q-card-section class="text-center">
              <q-icon :name="module.icon" color="primary" size="3em" class="q-mb-sm" />
              <div class="text-subtitle1">{{ module.name }}</div>
            </q-card-section>
          </q-card>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div class="text-center text-grey-7 q-mt-xl">
      <div class="text-caption">
        © 2024 School Digital Platform. All rights reserved.
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { useLandingHome } from '@/composables/landing/useLandingHome'

// Destructure semua dari composable
const {
  features,
  modules,
  isMaintenanceMode,
  navigateTo,
  goToLogin,
  goToRegister,
} = useLandingHome()
</script>

<style scoped>
.cursor-pointer {
  cursor: pointer;
  transition: transform 0.2s;
}

.cursor-pointer:hover {
  transform: translateY(-4px);
}
</style>
