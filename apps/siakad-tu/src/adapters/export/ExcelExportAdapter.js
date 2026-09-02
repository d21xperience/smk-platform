// src/adapters/export/ExcelExportAdapter.js
import { CsvExportAdapter } from './CsvExportAdapter.js'

export class ExcelExportAdapter extends CsvExportAdapter {
  async export(reportResult, options = {}) {
    // Placeholder: nanti pakai library xlsx
    const csvBlob = await super.export(reportResult, options)
    return new Blob([csvBlob], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    })
  }
}
