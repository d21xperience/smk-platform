import { BaseExportAdapter } from './BaseExportAdapter'
import { CsvExportAdapter } from './CsvExportAdapter'

export class ExcelExportAdapter extends BaseExportAdapter {
  export(data, filename = 'report.xlsx') {
    console.log('[Excel Export] Mock export:', { data, filename })
    // Nanti implementasi dengan library xlsx
    // Untuk mock, kita panggil CSV dulu atau simulasikan download.
    // Kita tawarkan output sebagai CSV dulu untuk testing.
    const csvAdapter = new CsvExportAdapter()
    csvAdapter.export(data, filename.replace('.xlsx', '.csv'))
    console.log('[Excel Export] Mock: Excel export simulated as CSV for now.')
  }
}
