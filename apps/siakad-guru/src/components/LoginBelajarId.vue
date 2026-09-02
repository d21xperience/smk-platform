<template>
  <!-- ✨ AWAL AREA LOGIN ALTERNATIF BELAJAR.ID ✨ -->
  <div class="row items-center q-my-lg text-grey-5">
    <q-separator class="col" />
    <div class="q-px-sm text-caption text-grey-7">Atau masuk dengan</div>
    <q-separator class="col" />
  </div>

  <div>
    <q-btn
      outline
      no-caps
      color="blue-7"
      class="full-width q-py-sm text-weight-bold bg-blue-1"
      :loading="loadingBelajarId"
      :disable="loading"
      @click="handleBelajarIdLogin"
    >
      <!-- Menggunakan logo Google / Belajar.id kustom -->
      <template v-slot:default>
        <div class="row items-center no-wrap">
          <q-icon
            name="img:https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRPR93yd93d7J4Mb2uKdVi73PIcyVCGcrOoKiduON7slQ&s"
            size="xs"
            class="q-mr-sm"
          />
          <span>Masuk dengan Akun Belajar.id</span>
        </div>
      </template>
    </q-btn>
  </div>
  <!-- ✨ AKHIR AREA LOGIN ALTERNATIF BELAJAR.ID ✨ -->
</template>

<script setup>
import { ref } from 'vue'

const loadingBelajarId = ref(false)
const loading = ref(false)
const handleBelajarIdLogin = () => {
  alert('Segera Hadir!!! \nSaat ini belum bisa digunakan.')
}
</script>
