<!-- ============================================================
  components/BkkSection.vue
  Bursa Kerja Khusus — lowongan kerja alumni
============================================================ -->
<template>
  <section id="bkk" :class="['q-py-xl', $q.dark.isActive ? 'bg-grey-10' : 'bg-amber-1']">
    <div class="max-width-center q-mx-auto q-px-md">
      <div class="row items-end justify-between q-mb-lg">
        <SectionHeader
          title="Info Lowongan Kerja BKK"
          subtitle="Pusat penyaluran kerja alumni yang terhubung langsung dengan industri mitra."
        />
        <q-btn
          color="primary"
          no-caps
          label="Hubungi Admin BKK"
          icon="chat"
          class="gt-xs text-weight-bold"
          @click="hubungiBkk"
        />
      </div>

      <div class="row q-col-gutter-lg">
        <div v-for="(job, i) in BKK_JOBS" :key="i" class="col-12 col-sm-6">
          <q-card class="bkk-card shadow-1" flat bordered>
            <q-card-section class="q-pa-lg">
              <!-- Header perusahaan -->
              <div class="row items-center justify-between no-wrap q-mb-md">
                <div class="row items-center no-wrap">
                  <q-avatar
                    :color="job.color"
                    text-color="white"
                    size="44px"
                    class="text-weight-bold shadow-1 q-mr-sm"
                    font-size="18px"
                  >
                    {{ job.initial }}
                  </q-avatar>
                  <div>
                    <div class="text-subtitle2 text-weight-bold lines-1">{{ job.company }}</div>
                    <div class="text-caption text-grey-5 row items-center">
                      <q-icon name="place" size="xs" class="q-mr-xs" />
                      {{ job.location }}
                    </div>
                  </div>
                </div>
                <q-chip
                  dense
                  :color="job.color"
                  text-color="white"
                  class="text-weight-bold text-caption"
                >
                  Khusus {{ job.majorTarget }}
                </q-chip>
              </div>

              <!-- Posisi -->
              <div class="text-h6 text-weight-bold text-primary q-mb-sm lines-2">
                {{ job.position }}
              </div>

              <!-- Requirements -->
              <q-list dense class="req-list">
                <q-item v-for="(req, ri) in job.requirements" :key="ri" class="q-px-none q-py-xs">
                  <q-item-section avatar style="min-width: 20px">
                    <q-icon name="check_circle" :color="job.color" size="xs" />
                  </q-item-section>
                  <q-item-section class="text-caption text-grey-7">{{ req }}</q-item-section>
                </q-item>
              </q-list>
            </q-card-section>

            <q-separator />
            <q-card-actions
              class="q-px-lg q-py-sm justify-between items-center"
              :class="$q.dark.isActive ? 'bg-grey-9' : 'bg-grey-1'"
            >
              <div class="text-caption text-negative text-weight-bold row items-center">
                <q-icon name="schedule" size="xs" class="q-mr-xs" />
                Batas: {{ job.deadline }}
              </div>
              <q-btn
                unelevated
                no-caps
                dense
                :color="job.color"
                label="Lamar Sekarang"
                icon="send"
                class="text-weight-bold q-px-md"
                @click="openUrl(job.applyUrl)"
              />
            </q-card-actions>
          </q-card>
        </div>
      </div>

      <div class="row justify-center q-mt-lg lt-sm">
        <q-btn
          color="primary"
          no-caps
          label="Hubungi Admin BKK"
          icon="chat"
          class="full-width text-weight-bold"
          @click="hubungiBkk"
        />
      </div>
    </div>
  </section>
</template>

<script setup>
import { useQuasar } from 'quasar'
import { BKK_JOBS } from '../../mocks/data/schoolData'
import { useSchool } from '../../composables/useSchool'
import SectionHeader from './SectionHeader.vue'

const $q = useQuasar()
const { hubungiBkk, openUrl } = useSchool()
</script>

<style scoped>
.bkk-card {
  border-radius: 16px;
  transition:
    transform 0.2s,
    box-shadow 0.2s;
  height: 100%;
}

.bkk-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1) !important;
}

.req-list .q-item {
  min-height: 28px;
}

.lines-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.lines-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.max-width-center {
  max-width: 1100px;
}
</style>
