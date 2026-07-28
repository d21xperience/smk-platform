<template>
  <section id="berita" class="q-py-xl bg-body">
    <div class="news-container">
      <div class="row items-end justify-between q-mb-lg">
        <SectionHeader title="Berita & Prestasi Terbaru"
          subtitle="Kabar terkini dari kegiatan akademik, prestasi siswa, dan informasi sekolah." />
        <q-btn flat color="primary" no-caps label="Semua Berita" icon-right="arrow_forward" to="/berita"
          class="gt-xs text-weight-bold q-mb-sm modern-btn" />
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="row q-col-gutter-md">
        <div v-for="n in 3" :key="n" class="col-12 col-sm-4">
          <q-card flat bordered class="skeleton-card" style="height: 320px">
            <q-img :ratio="4 / 3" class="skeleton-img" />
            <q-card-section class="q-pa-md">
              <div class="skeleton-line q-mb-sm" style="width: 60%"></div>
              <div class="skeleton-line" style="width: 80%"></div>
            </q-card-section>
          </q-card>
        </div>
      </div>

      <!-- Error State -->
      <q-banner v-else-if="error" class="bg-negative text-white q-mt-md" dense>
        {{ error }}
      </q-banner>

      <!-- News Grid -->
      <div v-else class="row q-col-gutter-md">
        <div v-for="(item, index) in displayNews" :key="item.id" class="col-12 col-sm-4">
          <NewsCard :news="item" :is-featured="index === 0" />
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed } from 'vue'
import { useNews } from '../composables/useNews'
// import SectionHeader from 'components/SectionHeader.vue'
import NewsCard from 'components/NewsCard.vue'

const { newsList, loading, error } = useNews()
const displayNews = computed(() => newsList.value.slice(0, 3))
</script>

<style scoped>
.news-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.modern-btn {
  transition: transform 0.2s, background 0.2s;
}

.modern-btn:hover {
  transform: translateX(4px);
  background: rgba(var(--q-primary-rgb), 0.08) !important;
}

.skeleton-card {
  background: #f8f9fa;
  border: none;
}

.skeleton-img {
  background: #e9ecef;
  height: 200px;
}

.skeleton-line {
  height: 12px;
  background: #e9ecef;
  border-radius: 6px;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {

  0%,
  100% {
    opacity: 0.6;
  }

  50% {
    opacity: 1;
  }
}
</style>
