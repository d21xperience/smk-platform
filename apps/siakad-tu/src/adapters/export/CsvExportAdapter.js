// src/adapters/export/CsvExportAdapter.js
import { BaseExportAdapter } from './BaseExportAdapter.js'

export class CsvExportAdapter extends BaseExportAdapter {
  // eslint-disable-next-line no-unused-vars
  async export(reportResult, options = {}) {
    const { columns, data } = reportResult

    // Header
    const headerRow = columns.map((col) => this._escapeCsv(col.header)).join(',')

    // Data rows
    const dataRows = data.map((row) => {
      return columns
        .map((col) => {
          const value = row[col.field]
          return this._escapeCsv(this._formatValue(value, col.type))
        })
        .join(',')
    })

    const csv = [headerRow, ...dataRows].join('\n')
    return new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  }

  _escapeCsv(value) {
    if (value == null) return ''
    const str = String(value)
    if (str.includes(',') || str.includes('"') || str.includes('\n')) {
      return `"${str.replace(/"/g, '""')}"`
    }
    return str
  }

  _formatValue(value, type) {
    if (value == null) return ''
    switch (type) {
      case 'percentage':
        return `${value}%`
      default:
        return value
    }
  }
}
