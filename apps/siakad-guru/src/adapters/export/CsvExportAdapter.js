import { BaseExportAdapter } from './BaseExportAdapter'

export class CsvExportAdapter extends BaseExportAdapter {
  export(data, filename = 'report.csv') {
    const { headers, rows } = data
    if (!rows || rows.length === 0) {
      console.warn('[CSV Export] No data to export.')
      return
    }

    // Build CSV string
    const headerRow = headers.map((h) => `"${h.label}"`).join(',')
    const csvRows = rows.map((row) => {
      return headers
        .map((h) => {
          const value = row[h.label] !== undefined ? row[h.label] : ''
          return `"${String(value).replace(/"/g, '""')}"`
        })
        .join(',')
    })

    const csvContent = [headerRow, ...csvRows].join('\n')

    // Simulate download (mock-first)
    // In real implementation, use FileSaver or Blob download
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    if (link.download !== undefined) {
      const url = URL.createObjectURL(blob)
      link.setAttribute('href', url)
      link.setAttribute('download', filename)
      link.style.visibility = 'hidden'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      URL.revokeObjectURL(url)
      console.log('[CSV Export] File downloaded:', filename)
    } else {
      console.warn('[CSV Export] Download not supported in this environment.')
      console.log(csvContent)
    }
  }
}
