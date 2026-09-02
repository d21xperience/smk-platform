// apps/siakad-tu/src/stores/contextStore.js

import { defineStore } from 'pinia'
// eslint-disable-next-line no-unused-vars
import { OperationalContext } from '../domain/context/OperationalContext.js'

/**
 * contextStore — HANYA mengelola state.
 * TIDAK BOLEH ada logika bisnis, validasi, atau derivasi di sini.
 * Semua logika ada di ContextEngine.
 */
export const useContextStore = defineStore('operationalContext', {
  state: () => ({
    /** @type {OperationalContext | null} */
    _context: null,
    /** @type {boolean} */
    isContextReady: false,
  }),

  getters: {
    /**
     * Getter reaktif untuk digunakan di Vue components
     */
    current: (state) => state._context,
    schoolId: (state) => state._context?.schoolId ?? null,
    schoolName: (state) => state._context?.schoolName ?? '',
    academicYear: (state) => state._context?.academicYear ?? '',
    semester: (state) => state._context?.semester ?? null,
    periodId: (state) => state._context?.periodId ?? null,
    displayLabel: (state) => state._context?.displayLabel ?? 'Pilih Context',
    hasContext: (state) => state._context !== null && state._context.isValid(),
  },

  actions: {
    /**
     * Set context — dipanggil oleh Service/Composable SETELAH Engine memvalidasi.
     * Store TIDAK melakukan validasi sendiri.
     */
    setContext(operationalContext) {
      this._context = operationalContext
      this.isContextReady = true
    },

    /**
     * Clear context — dipanggil saat logout
     */
    clearContext() {
      this._context = null
      this.isContextReady = false
    },
  },
})
