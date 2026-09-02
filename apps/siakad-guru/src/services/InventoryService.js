import { InventoryEngine } from '../domain/inventory/engine/InventoryEngine.js'
import { InventoryItem } from '../domain/inventory/models/InventoryItem.js'
import { DamageReport } from '../domain/inventory/models/DamageReport.js'
import { AssetDamageReportedEvent } from '../domain/inventory/events/AssetDamageReported.js'

export class InventoryService {
  constructor({ inventoryAdapter, eventDispatcher }) {
    this.inventoryAdapter = inventoryAdapter
    this.eventDispatcher = eventDispatcher
  }

  async loadItemsByClass({ classId }) {
    const rawData = await this.inventoryAdapter.fetchItemsByClass({ classId })
    return rawData.map((item) => new InventoryItem(item))
  }

  async loadItemById({ itemId }) {
    const rawData = await this.inventoryAdapter.fetchItemById({ itemId })

    if (!rawData) {
      return null
    }

    const item = new InventoryItem(rawData)
    InventoryEngine.validateItem(item)
    return item
  }

  async submitDamageReport({
    itemId,
    itemName,
    condition,
    description,
    reportMethod,
    date,
    teacherId,
    teacherName,
    classId,
    className,
    schoolId,
    academicYearId,
    semesterId,
  }) {
    // Validate input using domain engine
    InventoryEngine.validateReportInput({ condition, description, reportMethod })

    const rawData = await this.inventoryAdapter.submitDamageReport({
      itemId,
      itemName,
      condition,
      description,
      reportMethod,
      date,
      teacherId,
      teacherName,
      classId,
      className,
      schoolId,
      academicYearId,
      semesterId,
    })

    const report = new DamageReport(rawData)
    InventoryEngine.validateReport(report)

    if (this.eventDispatcher) {
      this.eventDispatcher.dispatch(new AssetDamageReportedEvent(report))
    }

    return report
  }
}
