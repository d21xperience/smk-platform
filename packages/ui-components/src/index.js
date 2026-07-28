// =============================================================================
// SHARED UI COMPONENTS - MAIN EXPORT
// =============================================================================
// File ini mengekspor semua shared components dan composables
// yang bisa digunakan oleh aplikasi frontend (web-sekolah, siakad-guru, ppdb)
//
// CARA PENGGUNAAN DI APP:
// import { AppButton, useNotification } from '@smk-platform/ui-components'
// =============================================================================

// =============================================================================
// COMMON COMPONENTS
// =============================================================================
// export { default as AppButton } from './components/common/AppButton.vue'
// export { default as AppTable } from './components/common/AppTable.vue'
// export { default as AppModal } from './components/common/AppModal.vue'
// export { default as AppCard } from './components/common/AppCard.vue'
// export { default as AppEmptyState } from './components/common/AppEmptyState.vue'
// export { default as AppLoader } from './components/common/AppLoader.vue'

// =============================================================================
// FORM COMPONENTS
// =============================================================================
// export { default as AppDatePicker } from './components/forms/AppDatePicker.vue'
// export { default as AppFileUpload } from './components/forms/AppFileUpload.vue'
// export { default as AppInput } from './components/forms/AppInput.vue'
// export { default as AppSelect } from './components/forms/AppSelect.vue'

// =============================================================================
// LAYOUT COMPONENTS
// =============================================================================
// export { default as AppFooter } from './components/layouts/AppFooter.vue'
export { default as AppHeader } from "./components/layouts/AppHeader.vue";
export { default as AppLayout } from "./components/layouts/AppLayout.vue";
export { default as AppSidebar } from "./components/layouts/AppSidebar.vue";

// =============================================================================
// COMPOSABLES
// =============================================================================
// export { useDateFormat } from './composables/useDateFormat.js'
// export { useLoading } from './composables/useLoading.js'
// export { useNotification } from './composables/useNotification.js'
// export { useValidation } from './composables/useValidation.js'
// =============================================================================
// UTILITIES
// =============================================================================
export * from "./utils/constants.js";
export * from "./utils/formatters.js";
export * from "./utils/validators.js";
