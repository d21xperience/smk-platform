import { BaseExportAdapter } from './BaseExportAdapter'

export class PdfExportAdapter extends BaseExportAdapter {
  export(data, filename = 'report.pdf') {
    console.log('[PDF Export] Mock export:', { data, filename })
    // Nanti implementasi dengan jsPDF atau library lain.
    // Untuk mock, kita log data ke console dan tampilkan preview sederhana.
    console.log('--- PDF PREVIEW ---')
    console.table(data.rows)
    console.log('--- SUMMARY ---', data.summary)
    alert(`PDF Export Mock: ${filename} generated. Check console for data.`)
  }
}
