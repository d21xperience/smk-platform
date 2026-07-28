// src/engine/ReportingEngine.js
import { ReportResult } from '@/models/ReportResult'
// import { SortDefinition, GroupDefinition, AggregateDefinition } from '@/models/ReportDefinition'

export class ReportingEngine {
  /**
   * @param {import('./ProjectionReader').ProjectionReader} reader
   */
  constructor(reader) {
    this.reader = reader
  }

  /**
   * Eksekusi laporan berdasarkan definisi dan parameter filter/sort.
   * @param {import('src/models/ReportDefinition').ReportDefinition} definition
   * @param {Object} [filterValues={}] - nilai filter { field: value }
   * @param {Array<SortDefinition>} [customSort] - sort kustom dari UI
   * @returns {ReportResult}
   */
  execute(definition, filterValues = {}, customSort = null) {
    // 1. Ambil data mentah dari proyeksi
    let data = this._readData(definition)

    // 2. Terapkan filter
    data = this._applyFilters(data, definition.filters, filterValues)

    // 3. Terapkan grouping (jika ada)
    if (definition.grouping) {
      data = this._applyGrouping(data, definition.grouping)
    }

    // 4. Terapkan sorting
    const sort = customSort || definition.defaultSort
    data = this._applySort(data, sort)

    // 5. Hitung agregasi
    const aggregates = this._calculateAggregates(data, definition.aggregates)

    // 6. Hitung footer
    const footer = definition.footer
      ? this._calculateAggregates(data, definition.footer.aggregates)
      : null

    // 7. Pilih kolom yang ditampilkan
    data = this._selectColumns(data, definition.columns)

    return new ReportResult(definition, data, { aggregates, footer })
  }

  // --- Private methods ---

  _readData(definition) {
    const raw = this.reader.read(definition.projectionSource)
    // Jika objek (map), ubah jadi array
    if (!Array.isArray(raw)) {
      return Object.values(raw)
    }
    return [...raw]
  }

  _applyFilters(data, filterDefs, filterValues) {
    if (!filterDefs || filterDefs.length === 0) return data
    return data.filter((row) => {
      return filterDefs.every((filter) => {
        const value = filterValues[filter.field]
        if (value === undefined || value === null || value === '') return true
        const fieldValue = this._getNestedValue(row, filter.field)
        switch (filter.type) {
          case 'text':
            return String(fieldValue).toLowerCase().includes(String(value).toLowerCase())
          case 'select':
            return fieldValue === value
          case 'date-range':
            if (value.start && fieldValue < value.start) return false
            if (value.end && fieldValue > value.end) return false
            return true
          case 'number-range':
            if (value.min !== undefined && fieldValue < value.min) return false
            if (value.max !== undefined && fieldValue > value.max) return false
            return true
          default:
            return true
        }
      })
    })
  }

  _applyGrouping(data, groupDef) {
    const groups = {}
    data.forEach((row) => {
      const key = this._getNestedValue(row, groupDef.field)
      if (!groups[key]) {
        groups[key] = {
          _groupKey: key,
          _groupLabel: groupDef.labelField ? this._getNestedValue(row, groupDef.labelField) : key,
          _items: [],
        }
      }
      groups[key]._items.push(row)
    })
    return Object.values(groups)
  }

  _applySort(data, sortDefs) {
    if (!sortDefs || sortDefs.length === 0) return data
    const sorted = [...data]
    sorted.sort((a, b) => {
      for (const sort of sortDefs) {
        const valA = this._getNestedValue(a, sort.field)
        const valB = this._getNestedValue(b, sort.field)
        let cmp = 0
        if (typeof valA === 'string') cmp = valA.localeCompare(valB)
        else cmp = valA - valB
        if (cmp !== 0) return sort.direction === 'asc' ? cmp : -cmp
      }
      return 0
    })
    return sorted
  }

  _calculateAggregates(data, aggDefs) {
    if (!aggDefs || aggDefs.length === 0) return {}
    const result = {}
    aggDefs.forEach((agg) => {
      const values = data
        .map((row) => this._getNestedValue(row, agg.field))
        .filter((v) => v != null)
      switch (agg.method) {
        case 'sum':
          result[agg.label || agg.field] = values.reduce((s, v) => s + v, 0)
          break
        case 'average':
          result[agg.label || agg.field] =
            values.length > 0 ? values.reduce((s, v) => s + v, 0) / values.length : 0
          break
        case 'count':
          result[agg.label || agg.field] = values.length
          break
        case 'min':
          result[agg.label || agg.field] = values.length > 0 ? Math.min(...values) : 0
          break
        case 'max':
          result[agg.label || agg.field] = values.length > 0 ? Math.max(...values) : 0
          break
      }
    })
    return result
  }

  _selectColumns(data, columns) {
    return data.map((row) => {
      const selected = {}
      columns.forEach((col) => {
        selected[col.field] = this._getNestedValue(row, col.field)
      })
      return selected
    })
  }

  _getNestedValue(obj, path) {
    return path.split('.').reduce((current, key) => {
      return current !== undefined && current !== null ? current[key] : undefined
    }, obj)
  }
}
