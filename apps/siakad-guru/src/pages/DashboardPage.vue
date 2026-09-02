<template>
  <q-page padding>
    <div class="text-h5 q-mb-md">Dashboard</div>

    <div v-if="currentContext" class="q-mb-md">
      <q-card>
        <q-card-section>
          <div class="text-h6">Konteks Operasional Aktif</div>
          <q-separator class="q-my-md" />
          <q-list dense>
            <q-item>
              <q-item-section>
                <q-item-label>School ID</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-item-label caption>{{ currentContext.schoolId }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item>
              <q-item-section>
                <q-item-label>Academic Year ID</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-item-label caption>{{ currentContext.academicYearId }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item>
              <q-item-section>
                <q-item-label>Semester ID</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-item-label caption>{{ currentContext.semesterId }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item>
              <q-item-section>
                <q-item-label>User ID</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-item-label caption>{{ currentContext.userId }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-item>
              <q-item-section>
                <q-item-label>Role</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-item-label caption>{{ currentContext.role }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </q-card>
    </div>

    <q-banner v-else class="bg-warning text-white">
      Tidak ada konteks operasional yang dipilih.
    </q-banner>
  </q-page>
</template>

<script>
import { useContext } from '../composables/useContext.js'

export default {
  name: 'DashboardPage',

  setup() {
    const { currentContext } = useContext()

    return {
      currentContext
    }
  }
}
</script>
