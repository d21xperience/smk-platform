<template>
  <q-page class="q-pa-md">
    <q-linear-progress v-if="store.loading" indeterminate color="teal" class="q-mb-sm" />

    <div class="text-h5 q-mb-md"><q-icon name="fingerprint" class="q-mr-sm" />PIN Tanda Tangan Digital</div>

    <div class="row q-col-gutter-md">
      <!-- STATUS CARD -->
      <div class="col-12 col-md-5">
        <q-card>
          <q-card-section class="bg-teal-1">
            <div class="text-subtitle1 text-weight-medium">Status PIN</div>
          </q-card-section>
          <q-card-section v-if="store.pinStatus">
            <q-list dense>
              <q-item>
                <q-item-section>
                  <q-item-label>Status</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-chip :color="store.pinStatus.is_setup ? 'positive' : 'grey'" text-color="white" dense>
                    {{ store.pinStatus.is_setup ? 'Sudah Setup' : 'Belum Setup' }}
                  </q-chip>
                </q-item-section>
              </q-item>
              <q-item v-if="store.pinStatus.is_setup">
                <q-item-section>
                  <q-item-label>Lock Status</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-chip :color="store.pinStatus.is_locked ? 'negative' : 'positive'" text-color="white" dense>
                    {{ store.pinStatus.is_locked ? 'Terkunci' : 'Aktif' }}
                  </q-chip>
                </q-item-section>
              </q-item>
              <q-item v-if="store.pinStatus.is_setup">
                <q-item-section>
                  <q-item-label>Sisa Percobaan</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <span class="text-weight-bold">
                    {{ store.pinStatus.max_attempt - store.pinStatus.current_attempt }} / {{ store.pinStatus.max_attempt
                    }}
                  </span>
                </q-item-section>
              </q-item>
              <q-item v-if="store.pinStatus.locked_until">
                <q-item-section>
                  <q-item-label class="text-negative">
                    Terkunci s/d {{ formatDateTime(store.pinStatus.locked_until) }}
                  </q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
          <q-card-actions align="right">
            <q-btn flat icon="refresh" label="Refresh" @click="loadStatus" />
            <q-btn v-if="store.pinStatus?.is_locked" color="warning" icon="lock_open" label="Unlock PIN"
              @click="submitUnlock" :loading="store.loading" />
          </q-card-actions>
        </q-card>

        <q-card class="q-mt-md">
          <q-card-section class="bg-blue-1">
            <div class="text-subtitle2">ℹ️ Informasi</div>
          </q-card-section>
          <q-card-section>
            <div class="text-caption text-grey-8">
              <ul class="q-pl-md q-my-none">
                <li>PIN digunakan untuk <strong>tanda tangan digital</strong> pada dokumen keuangan</li>
                <li>PIN terdiri dari <strong>6 digit angka</strong></li>
                <li>Setelah 3× salah, PIN terkunci 15 menit</li>
                <li>Untuk unlock, hubungi admin atau gunakan tombol Unlock di atas</li>
              </ul>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <!-- ACTION CARD -->
      <div class="col-12 col-md-7">
        <!-- Setup PIN (jika belum) -->
        <q-card v-if="!store.pinStatus?.is_setup">
          <q-card-section class="bg-green-1">
            <div class="text-subtitle1 text-weight-medium">
              <q-icon name="add_circle" /> Setup PIN (Pertama Kali)
            </div>
          </q-card-section>
          <q-card-section>
            <q-input v-model="setupPin" label="PIN Baru (6 digit)" outlined dense type="password" maxlength="6"
              inputmode="numeric" class="q-mb-md" />
            <q-input v-model="setupPinConfirm" label="Konfirmasi PIN" outlined dense type="password" maxlength="6"
              inputmode="numeric" />
          </q-card-section>
          <q-card-actions align="right">
            <q-btn color="teal" icon="check" label="Setup PIN" :loading="store.loading" :disable="!isValidSetupPin"
              @click="submitSetup" />
          </q-card-actions>
        </q-card>

        <!-- Verify + Change (jika sudah) -->
        <q-card v-else>
          <q-tabs v-model="tab" dense active-color="teal" indicator-color="teal" align="left">
            <q-tab name="change" icon="edit" label="Ganti PIN" />
            <q-tab name="test" icon="verified" label="Test Verifikasi" />
          </q-tabs>
          <q-separator />

          <q-tab-panels v-model="tab" animated>
            <!-- Tab Ganti PIN -->
            <q-tab-panel name="change">
              <q-input v-model="changePin.old" label="PIN Lama (6 digit)" outlined dense type="password" maxlength="6"
                inputmode="numeric" class="q-mb-sm" />
              <q-input v-model="changePin.newPin" label="PIN Baru (6 digit)" outlined dense type="password"
                maxlength="6" inputmode="numeric" class="q-mb-sm" />
              <q-input v-model="changePin.newPinConfirm" label="Konfirmasi PIN Baru" outlined dense type="password"
                maxlength="6" inputmode="numeric" />
              <q-card-actions align="right">
                <q-btn color="teal" icon="save" label="Ganti PIN" :loading="store.loading" :disable="!isValidChangePin"
                  @click="submitChange" />
              </q-card-actions>
            </q-tab-panel>

            <!-- Tab Test Verify -->
            <q-tab-panel name="test">
              <div class="text-caption text-grey-7 q-mb-sm">
                Test PIN Anda untuk memastikan berfungsi. (Akan increment attempt counter jika salah)
              </div>
              <q-input v-model="testPin" label="PIN (6 digit)" outlined dense type="password" maxlength="6"
                inputmode="numeric" class="q-mb-sm" />
              <q-banner v-if="verifyResult" :class="verifyResult.valid ? 'bg-positive' : 'bg-negative'"
                class="text-white" rounded>
                <template v-slot:avatar>
                  <q-icon :name="verifyResult.valid ? 'check_circle' : 'error'" />
                </template>
                {{ verifyResult.message }}
                <div v-if="!verifyResult.valid" class="text-caption">
                  Sisa percobaan: {{ verifyResult.remaining_attempt }}
                </div>
              </q-banner>
              <q-card-actions align="right">
                <q-btn color="teal" icon="lock" label="Test" :loading="store.loading"
                  :disable="!/^\d{6}$/.test(testPin)" @click="submitVerify" />
              </q-card-actions>
            </q-tab-panel>
          </q-tab-panels>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useTTDStore } from '@/stores/finance/ttdStore'

const $q = useQuasar()
const store = useTTDStore()

const tab = ref('change')
const setupPin = ref('')
const setupPinConfirm = ref('')
const changePin = ref({ old: '', newPin: '', newPinConfirm: '' })
const testPin = ref('')
const verifyResult = ref(null)

const isValidSetupPin = computed(() =>
  /^\d{6}$/.test(setupPin.value) && setupPin.value === setupPinConfirm.value
)

const isValidChangePin = computed(() =>
  /^\d{6}$/.test(changePin.value.old) &&
  /^\d{6}$/.test(changePin.value.newPin) &&
  changePin.value.newPin === changePin.value.newPinConfirm &&
  changePin.value.old !== changePin.value.newPin
)

function formatDateTime(s) {
  if (!s) return '-'
  try { return new Date(s).toLocaleString('id-ID') } catch { return s }
}

async function loadStatus() {
  try { await store.getPinStatus() } catch {
    console.log('PinTTDPage.vue')
  }
}

async function submitSetup() {
  if (!isValidSetupPin.value) {
    $q.notify({ type: 'warning', message: 'PIN harus 6 digit & konfirmasi cocok' })
    return
  }
  try {
    await store.setupPin(setupPin.value)
    $q.notify({ type: 'positive', message: 'PIN berhasil disetup' })
    setupPin.value = ''
    setupPinConfirm.value = ''
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
  }
}

async function submitChange() {
  if (!isValidChangePin.value) {
    $q.notify({ type: 'warning', message: 'Periksa input (6 digit, konfirmasi cocok, beda dari lama)' })
    return
  }
  try {
    await store.changePin(changePin.value.old, changePin.value.newPin)
    $q.notify({ type: 'positive', message: 'PIN berhasil diganti' })
    changePin.value = { old: '', newPin: '', newPinConfirm: '' }
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message, timeout: 5000 })
  }
}

async function submitVerify() {
  if (!/^\d{6}$/.test(testPin.value)) return
  verifyResult.value = null
  try {
    const r = await store.verifyPin(testPin.value)
    verifyResult.value = r
    testPin.value = ''
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
  }
}

async function submitUnlock() {
  try {
    await store.unlockPin()
    $q.notify({ type: 'positive', message: 'PIN berhasil di-unlock' })
  } catch (e) {
    $q.notify({ type: 'negative', message: e.response?.data?.error || e.message })
  }
}

onMounted(loadStatus)
</script>
