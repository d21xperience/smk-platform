import { BaseEntity } from '@/domain/shared/BaseEntity'

export const REPORT_TYPES = {
  CLASS_ATTENDANCE: 'CLASS_ATTENDANCE',
  CLASS_GRADE: 'CLASS_GRADE',
  STUDENT_PROGRESS: 'STUDENT_PROGRESS',
  TEACHER_PERFORMANCE: 'TEACHER_PERFORMANCE',
  DASHBOARD_SUMMARY: 'DASHBOARD_SUMMARY',
}

export class ReportDefinition extends BaseEntity {
  constructor(props) {
    super(props)
    this.type = props.type
    this.name = props.name || ''
    this.description = props.description || ''
    this.columns = props.columns || [] // [{ key, label, type }]
    this.filters = props.filters || {} // { classId, semesterId, teacherId, dateRange }
    this.createdBy = props.createdBy
  }

  toJSON() {
    return {
      ...super.toJSON(),
      type: this.type,
      name: this.name,
      description: this.description,
      columns: this.columns,
      filters: this.filters,
      createdBy: this.createdBy,
    }
  }
}
