// @/stores/reporting.store.js
import { defineStore } from 'pinia'
import { ReportingEngine } from '@/engine/ReportingEngine'
import { ProjectionReader } from '@/engine/ProjectionReader'
import { CsvExportAdapter } from '@/adapters/export/CsvExportAdapter'
import { PdfExportAdapter } from '@/adapters/export/PdfExportAdapter'
import { ExcelExportAdapter } from '@/adapters/export/ExcelExportAdapter'

export const useReportingStore = defineStore('reporting', {
  state: () => ({
    engine: null,
    exportAdapters: {
      csv: new CsvExportAdapter(),
      pdf: new PdfExportAdapter(),
      excel: new ExcelExportAdapter(),
    },
    currentDefinition: null,
    currentResult: null,
    filterValues: {},
    loading: false,
    error: null,
  }),

  getters: {
    reportData: (state) => state.currentResult?.data || [],
    reportColumns: (state) => state.currentResult?.columns || [],
    reportMetadata: (state) => state.currentResult?.metadata || {},
  },

  actions: {
    /**
     * Inisialisasi engine dengan reader yang terhubung ke read stores.
     * @param {Object} readStores - map projection name → getter function
     */
    init(readStores) {
      const reader = new ProjectionReader(readStores)
      this.engine = new ReportingEngine(reader)
    },

    /**
     * Pilih dan eksekusi definisi laporan.
     * @param {import('@/models/ReportDefinition').ReportDefinition} definition
     * @param {Object} [filters={}]
     */
    async loadReport(definition, filters = {}) {
      if (!this.engine) throw new Error('Engine belum diinisialisasi')
      this.loading = true
      this.error = null
      try {
        this.currentDefinition = definition
        this.filterValues = { ...filters }
        this.currentResult = this.engine.execute(definition, this.filterValues)
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading = false
      }
    },

    /**
     * Terapkan filter dan refresh hasil.
     * @param {Object} filters
     */
    applyFilters(filters) {
      this.filterValues = { ...filters }
      if (this.currentDefinition && this.engine) {
        this.currentResult = this.engine.execute(this.currentDefinition, this.filterValues)
      }
    },

    /**
     * Ekspor hasil laporan ke format tertentu.
     * @param {string} format - 'csv', 'pdf', 'excel'
     * @returns {Promise<Blob>}
     */
    async exportReport(format = 'csv') {
      if (!this.currentResult) throw new Error('Tidak ada hasil laporan')
      const adapter = this.exportAdapters[format]
      if (!adapter) throw new Error(`Format '${format}' tidak didukung`)
      return adapter.export(this.currentResult)
    },
  },
})
