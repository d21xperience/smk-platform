// src/adapters/export/PdfExportAdapter.js
import { BaseExportAdapter } from './BaseExportAdapter'

export class PdfExportAdapter extends BaseExportAdapter {
  // eslint-disable-next-line no-unused-vars
  async export(reportResult, options = {}) {
    // Simulasi: buat HTML sederhana untuk PDF (bisa pakai jsPDF atau print)
    const html = this._buildHtml(reportResult)
    return new Blob([html], { type: 'text/html' }) // placeholder
  }

  _buildHtml(reportResult) {
    const { reportName, columns, data, metadata } = reportResult
    let html = `<h1>${reportName}</h1>
      <p>Generated: ${metadata.generatedAt}</p>
      <table border="1">
        <thead><tr>${columns.map((c) => `<th>${c.header}</th>`).join('')}</tr></thead>
        <tbody>`
    data.forEach((row) => {
      html += '<tr>'
      columns.forEach((col) => {
        html += `<td>${row[col.field] ?? ''}</td>`
      })
      html += '</tr>'
    })
    html += `</tbody></table>`
    if (metadata.aggregates) {
      html += '<h3>Agregasi</h3><ul>'
      Object.entries(metadata.aggregates).forEach(([key, val]) => {
        html += `<li>${key}: ${val}</li>`
      })
      html += '</ul>'
    }
    return `<!DOCTYPE html><html><body>${html}</body></html>`
  }
}
