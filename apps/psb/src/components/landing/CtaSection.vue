<!-- ============================================================
  components/CtaSection.vue
  Call-to-action banner bawah halaman
============================================================ -->
<template>
  <section class="cta-section text-white text-center relative-position overflow-hidden">
    <!-- Decorative circles -->
    <div class="cta-circle cta-circle--1" />
    <div class="cta-circle cta-circle--2" />

    <div class="max-width-center q-mx-auto q-px-md q-py-xl relative-position" style="z-index: 2">
      <q-badge
        color="amber-4"
        text-color="black"
        class="text-weight-bold q-px-md q-py-xs q-mb-lg text-caption"
      >
        PPDB 2025/2026 — Kuota Terbatas!
      </q-badge>

      <h2 class="text-h4 text-sm-h3 text-weight-black text-amber-3 q-mb-md">
        Mulai Masa Depan Cemerlang<br class="gt-xs" />
        di Sini
      </h2>

      <p class="text-subtitle1 text-blue-2 q-mb-xl cta-sub">
        Jangan lewatkan kesempatan bergabung dengan Kelas Khusus Industri binaan Astra dan jejaring
        teknologi terbaik. Unduh brosur atau daftar sekarang juga.
      </p>

      <div class="row justify-center q-gutter-md">
        <q-btn
          unelevated
          no-caps
          color="amber"
          text-color="black"
          label="Daftar PPDB Online"
          icon="assignment"
          class="text-weight-bold q-px-xl cta-btn"
          size="lg"
          to="/ppdb"
        />
        <q-btn
          outline
          no-caps
          color="white"
          label="Unduh Brosur Digital"
          icon="file_download"
          class="text-weight-bold q-px-lg"
          size="lg"
          @click="downloadBrosur"
        />
      </div>

      <!-- Trust indicators -->
      <div class="row justify-center q-mt-xl q-gutter-lg">
        <div
          v-for="(t, i) in trustItems"
          :key="i"
          class="row items-center text-blue-2 text-caption"
        >
          <q-icon :name="t.icon" color="amber-4" size="xs" class="q-mr-xs" />
          {{ t.label }}
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { useSchool } from '../../composables/useSchool'

const { downloadBrosur } = useSchool()

const trustItems = [
  { icon: 'verified', label: 'Terakreditasi A' },
  { icon: 'school', label: 'Kurikulum Merdeka' },
  { icon: 'handshake', label: '45+ Mitra Industri' },
  { icon: 'work_history', label: '92% Lulusan Langsung Kerja' },
]
</script>

<style scoped>
.cta-section {
  background: linear-gradient(135deg, #0d1b3e 0%, #162451 50%, #1a237e 100%);
  position: relative;
}

.cta-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.04);
  pointer-events: none;
}

.cta-circle--1 {
  width: 400px;
  height: 400px;
  top: -150px;
  right: -100px;
}

.cta-circle--2 {
  width: 300px;
  height: 300px;
  bottom: -100px;
  left: -50px;
}

.cta-sub {
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
  line-height: 1.7;
}

.cta-btn {
  border-radius: 50px !important;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
}

.max-width-center {
  max-width: 1100px;
}
</style>
