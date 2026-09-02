// @/adapters/export/BaseExportAdapter.js
export class BaseExportAdapter {
  /**
   * Ekspor ReportResult ke format tertentu.
   * @param {import('@/models/ReportResult').ReportResult} reportResult
   * @param {Object} [options] - opsi tambahan
   * @returns {Promise<Blob|string>}
   */
  // eslint-disable-next-line no-unused-vars
  async export(reportResult, options = {}) {
    throw new Error('export() harus diimplementasikan')
  }
}
