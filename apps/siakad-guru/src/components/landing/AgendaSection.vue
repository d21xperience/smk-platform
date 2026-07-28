<!-- ============================================================
  components/AgendaSection.vue
  Agenda sekolah + info libur terdekat
============================================================ -->
<template>
  <section id="agenda" :class="['q-py-xl', $q.dark.isActive ? 'bg-grey-10' : 'bg-grey-1']">
    <div class="max-width-center q-mx-auto q-px-md">
      <div class="row items-end justify-between q-mb-lg">
        <SectionHeader
          title="Agenda & Informasi Operasional"
          subtitle="Jadwal kegiatan terdekat, ujian, dan info hari libur untuk Siswa & Guru."
        />
        <q-btn
          outline color="primary" no-caps
          label="Semua Agenda" icon="calendar_month"
          to="/agenda" class="gt-xs text-weight-bold q-mb-sm"
        />
      </div>

      <div class="row q-col-gutter-lg">
        <!-- ── Daftar kegiatan ────────────────────────── -->
        <div class="col-12 col-md-7">
          <div class="text-subtitle2 text-weight-bold text-primary q-mb-sm row items-center gap-xs">
            <q-icon name="campaign" size="sm" />
            Kegiatan Mendatang
          </div>
          <q-card flat bordered :class="['rounded-borders overflow-hidden', $q.dark.isActive ? 'bg-grey-9' : 'bg-white']">
            <q-list separator>
              <q-item
                v-for="(a, i) in AGENDA"
                :key="i"
                clickable v-ripple
                class="q-py-md agenda-item"
              >
                <!-- Mini kalender -->
                <q-item-section avatar>
                  <div class="date-box shadow-1 overflow-hidden rounded-borders">
                    <div class="date-month bg-primary text-white text-caption text-weight-bold text-center q-py-xs text-uppercase">
                      {{ a.month }}
                    </div>
                    <div :class="['date-day text-center text-h5 text-weight-black q-py-xs', $q.dark.isActive ? 'bg-grey-8' : 'bg-grey-1']">
                      {{ a.day }}
                    </div>
                  </div>
                </q-item-section>

                <q-item-section>
                  <div class="row items-center q-mb-xs">
                    <q-badge
                      :color="a.category === 'Ujian' ? 'orange-8' : a.category === 'Rapat' ? 'purple-7' : 'blue-7'"
                      class="text-weight-bold text-caption"
                    >
                      {{ a.category }}
                    </q-badge>
                  </div>
                  <q-item-label class="text-subtitle2 text-weight-bold lines-1">
                    {{ a.title }}
                  </q-item-label>
                  <q-item-label caption class="row items-center q-mt-xs">
                    <q-icon name="place" size="xs" class="q-mr-xs" />
                    {{ a.location }}
                  </q-item-label>
                </q-item-section>

                <q-item-section side>
                  <q-icon name="chevron_right" color="grey-4" />
                </q-item-section>
              </q-item>
            </q-list>
          </q-card>
        </div>

        <!-- ── Info libur ─────────────────────────────── -->
        <div class="col-12 col-md-5">
          <div class="text-subtitle2 text-weight-bold text-negative q-mb-sm row items-center gap-xs">
            <q-icon name="beach_access" size="sm" />
            Libur Terdekat
          </div>
          <q-card flat bordered class="holiday-card bg-red-1 full-height flex flex-column justify-between">
            <q-card-section class="q-pa-lg">
              <div class="row items-start no-wrap q-gutter-md">
                <q-avatar color="negative" text-color="white" icon="celebration" size="48px" />
                <div>
                  <div class="text-subtitle2 text-weight-bold q-mb-xs">{{ HOLIDAY.title }}</div>
                  <div class="text-caption text-negative text-weight-bold row items-center">
                    <q-icon name="event" size="xs" class="q-mr-xs" />
                    {{ HOLIDAY.dateRange }}
                  </div>
                </div>
              </div>
              <q-separator class="q-my-md" color="red-2" />
              <p class="text-caption text-grey-7 q-mb-none" style="line-height:1.7">
                {{ HOLIDAY.info }}
              </p>
            </q-card-section>
            <q-card-actions align="right" class="q-pb-md q-px-md">
              <q-btn
                flat dense no-caps color="negative"
                label="Detail Aturan Libur" icon-right="chevron_right"
                class="text-weight-bold text-caption"
                to="/agenda"
              />
            </q-card-actions>
          </q-card>
        </div>
      </div>

      <!-- Mobile: tombol semua agenda -->
      <div class="row justify-center q-mt-lg lt-sm">
        <q-btn
          outline color="primary" no-caps
          label="Lihat Semua Agenda" icon="calendar_month"
          to="/agenda" class="full-width text-weight-bold"
        />
      </div>
    </div>
  </section>
</template>

<script setup>
import { useQuasar } from 'quasar'
import { AGENDA, HOLIDAY } from '../../mocks/data/schoolData'
import SectionHeader from './SectionHeader.vue'

const $q = useQuasar()
</script>

<style scoped>
.date-box { width: 52px; border: 1px solid #e0e0e0; }
.date-month { font-size: 9px; letter-spacing: .5px; }
.date-day { line-height: 1.2; }

.agenda-item { transition: background .15s; }
.holiday-card { border: 1px solid #ffcdd2 !important; border-radius: 14px; }
.lines-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.max-width-center { max-width: 1100px; }
.gap-xs { gap: 4px; }
</style>
