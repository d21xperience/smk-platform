<template>
  <q-card
    flat
    bordered
    class="news-card"
    :class="{ 'news-card--featured': isFeatured }"
    clickable
    v-ripple
    role="link"
    tabindex="0"
    :aria-label="`Baca berita: ${news.title}`"
    @click="navigate"
    @keyup.enter="navigate"
  >
    <q-img :src="news.img" loading="lazy" :ratio="isFeatured ? 16 / 9 : 4 / 3" class="news-img">
      <div v-if="isFeatured" class="absolute-top-left q-ma-sm">
        <q-badge color="accent" class="text-weight-bold q-px-sm news-badge"> 🔥 Featured </q-badge>
      </div>
      <div class="absolute-bottom news-overlay" />
    </q-img>

    <q-card-section class="q-pa-md">
      <div class="news-meta row items-center q-mb-sm text-caption text-grey-6">
        <q-icon name="calendar_today" size="xs" class="q-mr-xs" />
        <time>{{ news.date }}</time>
        <span class="q-mx-xs">•</span>
        <q-icon name="person" size="xs" class="q-mr-xs" />
        <span>{{ news.author }}</span>
        <span v-if="news.category" class="q-ml-sm">
          <q-badge outline color="primary" size="sm">{{ news.category }}</q-badge>
        </span>
      </div>
      <div class="news-title text-subtitle2 text-weight-bold">
        {{ news.title }}
      </div>
    </q-card-section>

    <q-card-actions class="q-pt-none q-px-md q-pb-md justify-end">
      <q-btn
        flat
        dense
        no-caps
        color="primary"
        label="Baca Selengkapnya"
        icon-right="arrow_forward"
        class="text-caption text-weight-bold read-more-btn"
        :to="news.to"
      />
    </q-card-actions>
  </q-card>
</template>

<script setup>
import { useRouter } from 'vue-router'

const props = defineProps({
  news: { type: Object, required: true },
  isFeatured: { type: Boolean, default: false },
})

const router = useRouter()
const navigate = () => router.push(props.news.to)
</script>

<style scoped>
.news-card {
  border-radius: 16px;
  overflow: hidden;
  height: 100%;
  background: #fff;
  border: 1px solid rgba(0, 0, 0, 0.08);
  transition:
    transform 0.3s cubic-bezier(0.4, 0, 0.2, 1),
    box-shadow 0.3s ease;
}

.news-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.12) !important;
}

.news-card--featured {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.news-img {
  transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.news-card:hover .news-img {
  transform: scale(1.05);
}

.news-overlay {
  background: linear-gradient(to top, rgba(0, 0, 0, 0.35) 0%, transparent 100%);
  height: 35%;
}

.news-badge {
  backdrop-filter: blur(8px);
  background: rgba(255, 82, 82, 0.85) !important;
}

.news-title {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.45;
  transition: color 0.2s;
}

.news-card:hover .news-title {
  color: var(--q-primary);
}

.read-more-btn {
  opacity: 0.85;
  transition:
    opacity 0.2s,
    transform 0.2s;
}

.read-more-btn:hover {
  opacity: 1;
  transform: translateX(4px);
}

@media (max-width: 599px) {
  .news-title {
    -webkit-line-clamp: 3;
  }
}
</style>
