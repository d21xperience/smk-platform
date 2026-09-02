import { BaseEntity } from '@/domain/shared/BaseEntity'

export class ReportResult extends BaseEntity {
  constructor(props) {
    super(props)
    this.reportDefinitionId = props.reportDefinitionId
    this.data = props.data || [] // array of rows
    this.summary = props.summary || {} // { total, average, etc }
    this.generatedAt = props.generatedAt || new Date().toISOString()
    this.generatedBy = props.generatedBy
    this.format = props.format || 'JSON' // JSON, CSV, EXCEL, PDF
  }

  toJSON() {
    return {
      ...super.toJSON(),
      reportDefinitionId: this.reportDefinitionId,
      data: this.data,
      summary: this.summary,
      generatedAt: this.generatedAt,
      generatedBy: this.generatedBy,
      format: this.format,
    }
  }
}
