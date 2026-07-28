export class ReportResult {
  /**
   * @param {ReportDefinition} definition
   * @param {Array<Object>} data - baris data
   * @param {Object} [metadata] - informasi tambahan (total, waktu generate, dll.)
   */
  constructor(definition, data = [], metadata = {}) {
    this.reportId = definition.id
    this.reportName = definition.name
    this.columns = definition.columns
    this.data = data
    this.metadata = {
      generatedAt: new Date().toISOString(),
      totalRows: data.length,
      ...metadata,
    }
  }
}
