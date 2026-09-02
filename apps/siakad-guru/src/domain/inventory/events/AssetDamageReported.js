import { EVENT_ASSET_DAMAGE_REPORTED } from '../../../events/Type.js'

export class AssetDamageReportedEvent {
  constructor(report) {
    this.type = EVENT_ASSET_DAMAGE_REPORTED
    this.payload = {
      reportId: report.id,
      itemId: report.itemId,
      itemName: report.itemName,
      condition: report.condition,
      reportMethod: report.reportMethod,
      date: report.date,
      teacherId: report.teacherId,
      classId: report.classId,
      schoolId: report.schoolId,
      academicYearId: report.academicYearId,
      semesterId: report.semesterId,
      timestamp: new Date().toISOString(),
    }
  }
}
