export class ReportDefinition {
  /**
   * @param {Object} options
   * @param {string} options.id - unique identifier (e.g., 'student-attendance')
   * @param {string} options.name - nama laporan
   * @param {string} options.projectionSource - nama proyeksi di Read Store
   * @param {Array<ColumnDefinition>} options.columns - definisi kolom
   * @param {Array<FilterDefinition>} [options.filters] - filter yang tersedia
   * @param {Array<SortDefinition>} [options.defaultSort] - default sorting
   * @param {GroupDefinition} [options.grouping] - grouping
   * @param {Array<AggregateDefinition>} [options.aggregates] - agregasi (total, rata-rata)
   * @param {FooterDefinition} [options.footer] - footer laporan
   */
  constructor({
    id,
    name,
    projectionSource,
    columns,
    filters = [],
    defaultSort = [],
    grouping = null,
    aggregates = [],
    footer = null,
  }) {
    this.id = id
    this.name = name
    this.projectionSource = projectionSource // e.g., 'studentProgress', 'classAttendance'
    this.columns = columns
    this.filters = filters
    this.defaultSort = defaultSort
    this.grouping = grouping
    this.aggregates = aggregates
    this.footer = footer
  }
}

export class ColumnDefinition {
  /**
   * @param {string} field - path ke data di proyeksi (e.g., 'studentName', 'attendanceSummary.persentaseKehadiran')
   * @param {string} header - label kolom
   * @param {string} [type='string'] - tipe data ('string', 'number', 'percentage', 'date')
   * @param {boolean} [sortable=true]
   * @param {string} [width] - lebar kolom (CSS)
   */
  constructor({ field, header, type = 'string', sortable = true, width = 'auto' }) {
    this.field = field
    this.header = header
    this.type = type
    this.sortable = sortable
    this.width = width
  }
}

export class FilterDefinition {
  /**
   * @param {string} field - field yang difilter
   * @param {string} label - label UI
   * @param {string} type - 'text', 'select', 'date-range', 'number-range'
   * @param {Array<{ label: string, value: any }>} [options] - untuk tipe select
   */
  constructor({ field, label, type = 'text', options = [] }) {
    this.field = field
    this.label = label
    this.type = type
    this.options = options
  }
}

export class SortDefinition {
  /**
   * @param {string} field
   * @param {'asc'|'desc'} direction
   */
  constructor({ field, direction = 'asc' }) {
    this.field = field
    this.direction = direction
  }
}

export class GroupDefinition {
  /**
   * @param {string} field - field untuk grouping
   * @param {string} [labelField] - label untuk grup
   */
  constructor({ field, labelField = null }) {
    this.field = field
    this.labelField = labelField || field
  }
}

export class AggregateDefinition {
  /**
   * @param {string} field - field yang diagregasi
   * @param {string} method - 'sum', 'average', 'count', 'min', 'max'
   * @param {string} label - label untuk hasil agregasi
   */
  constructor({ field, method = 'sum', label = '' }) {
    this.field = field
    this.method = method
    this.label = label
  }
}

export class FooterDefinition {
  /**
   * @param {Array<AggregateDefinition>} aggregates - agregasi untuk footer
   */
  constructor({ aggregates = [] }) {
    this.aggregates = aggregates
  }
}
