// apps/siakad-tu/src/composables/landing/useLandingLayout.js

import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'

/**
 * Composable untuk LandingLayout.
 *
 * Tanggung jawab:
 * - Mengelola state drawer (mobile menu)
 * - Mengelola state search dialog
 * - Mengelola state login required dialog
 * - Menyediakan menu items untuk navigation
 * - Menangani logic "login required" untuk akses detail
 */
export function useLandingLayout() {
  const router = useRouter()
  const authStore = useAuthStore()

  // === STATE ===
  const leftDrawerOpen = ref(false)
  const isSearchOpen = ref(false)
  const showLoginRequiredDialog = ref(false)
  const pendingNavigation = ref(null) // Menyimpan route yang ingin dituju setelah login

  // === MENU ITEMS ===
  const menuItems = ref([
    { name: 'Beranda', routeName: 'landing-home', icon: 'home' },
    { name: 'Kurikulum', routeName: 'landing-kurikulum', icon: 'school' },
    { name: 'Kesiswaan', routeName: 'landing-kesiswaan', icon: 'groups' },
    { name: 'Sarpras', routeName: 'landing-sarpras', icon: 'inventory' },
    { name: 'Hubin', routeName: 'landing-hubin', icon: 'handshake' },
    { name: 'PSB', routeName: 'landing-psb', icon: 'person_add' },
  ])

  const itemPortal = ref([
    { title: 'Our Website', url: 'https://kemdikbud.go.id', target: '_blank' },
    { title: 'CBT', url: 'https://kemdikbud.go.id', target: '_blank' },
    { title: 'e-Rapor', url: 'https:e-rapor.smkpasjat.my.id', target: '_blank' },
    { title: 'Ruang GTK', url: 'https://guru.kemendikdasmen.go.id/', target: '_blank' },
    {
      title: 'Info GTK',
      url: 'https://info.gtk.kemendikdasmen.go.id/',
      target: '_blank',
    },
    {
      title: 'SIMPKB',
      url: 'https://paspor-gtk.simpkb.id/casgpo/login?service=https%3A%2F%2Fppg-backend.simpkb.id%2Fauth%2Flogin',
      target: '_blank',
    },
  ])

  // === COMPUTED ===
  const isAuthenticated = computed(() => authStore.isAuthenticated)
  const currentUser = computed(() => authStore.user)

  // === ACTIONS ===
  function toggleLeftDrawer() {
    leftDrawerOpen.value = false //!leftDrawerOpen.value
  }

  function openSearch() {
    isSearchOpen.value = true
  }

  function closeSearch() {
    isSearchOpen.value = false
  }

  /**
   * Handle klik pada portal dropdown
   */
  function handlePortalClick(portalType) {
    console.log(`Portal ${portalType} clicked`)
    // Implementasi sesuai kebutuhan
  }

  /**
   * Handle klik pada tombol "Lihat Detail" di halaman preview.
   * Jika belum login, tampilkan dialog login required.
   * Jika sudah login, redirect ke halaman detail.
   *
   * @param {string} targetRoute - Nama route tujuan (e.g., 'student-list')
   * @param {Object} params - Parameter untuk route (optional)
   */
  function requestAccess(targetRoute, params = {}) {
    if (isAuthenticated.value) {
      // Sudah login, langsung redirect
      router.push({ name: targetRoute, params })
    } else {
      // Belum login, simpan route tujuan dan tampilkan dialog
      pendingNavigation.value = { name: targetRoute, params }
      showLoginRequiredDialog.value = true
    }
  }

  /**
   * Handle klik "Login" pada dialog
   */
  function proceedToLogin() {
    showLoginRequiredDialog.value = false
    // Redirect ke halaman login dengan query parameter untuk kembali setelah login
    router.push({
      name: 'auth-login',
      query: {
        redirectTo: pendingNavigation.value?.name,
        ...pendingNavigation.value?.params,
      },
    })
  }

  /**
   * Handle klik "Batal" pada dialog
   */
  function cancelLogin() {
    showLoginRequiredDialog.value = false
    pendingNavigation.value = null
  }

  /**
   * Handle setelah login berhasil (dipanggil dari halaman login)
   */
  function handlePostLoginRedirect() {
    if (pendingNavigation.value) {
      router.push(pendingNavigation.value)
      pendingNavigation.value = null
    }
  }
  async function handleLogout() {
    try {
      await authStore.logout()
      router.push({ name: 'landing-home' })
    } catch (err) {
      console.error('[useLandingLayout] Logout error:', err)
    }
  }
  // === LIFECYCLE ===
  onMounted(() => {
    // Cek apakah ada pending navigation dari query parameter
    const query = router.currentRoute.value.query
    if (query.redirectTo) {
      pendingNavigation.value = {
        name: query.redirectTo,
        params: query,
      }
    }
  })

  // === RETURN ===
  return {
    // State
    leftDrawerOpen,
    isSearchOpen,
    showLoginRequiredDialog,
    pendingNavigation,

    // Computed
    isAuthenticated,
    currentUser,

    // Data
    menuItems,
    itemPortal,

    // Actions
    toggleLeftDrawer,
    openSearch,
    closeSearch,
    handlePortalClick,
    requestAccess,
    proceedToLogin,
    cancelLogin,
    handlePostLoginRedirect,
    handleLogout,
  }
}
