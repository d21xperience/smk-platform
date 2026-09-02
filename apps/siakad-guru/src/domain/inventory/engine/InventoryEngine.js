import { DamageReport } from '../models/DamageReport.js'
import { InventoryItem } from '../models/InventoryItem.js'
import { isValidCondition, isValidReportMethod } from '../models/InventoryCondition.js'

export class InventoryEngine {
  static validateItem(item) {
    if (!item || !(item instanceof InventoryItem)) {
      throw new Error('Invalid inventory item object')
    }
    if (!item.isValid()) {
      throw new Error('Inventory item is missing required fields')
    }
    return true
  }

  static validateReport(report) {
    if (!report || !(report instanceof DamageReport)) {
      throw new Error('Invalid damage report object')
    }
    if (!report.isValid()) {
      throw new Error('Damage report is missing required fields')
    }
    if (!isValidCondition(report.condition)) {
      throw new Error(`Invalid condition: ${report.condition}`)
    }
    if (!isValidReportMethod(report.reportMethod)) {
      throw new Error(`Invalid report method: ${report.reportMethod}`)
    }
    return true
  }

  static validateReportInput({ condition, description, reportMethod }) {
    if (!isValidCondition(condition)) {
      throw new Error(`Invalid condition: ${condition}`)
    }
    if (!isValidReportMethod(reportMethod)) {
      throw new Error(`Invalid report method: ${reportMethod}`)
    }
    if (!description || description.trim().length === 0) {
      throw new Error('Description is required for damage report')
    }
    return true
  }
}
