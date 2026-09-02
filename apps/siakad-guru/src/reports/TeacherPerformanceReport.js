// @/reports/TeacherPerformanceReport.js
import {
  ReportDefinition,
  ColumnDefinition,
  // FilterDefinition,
  SortDefinition,
  AggregateDefinition,
} from '@/models/ReportDefinition'

export const teacherPerformanceReport = new ReportDefinition({
  id: 'teacher-performance',
  name: 'Laporan Performa Guru',
  projectionSource: 'teacherDashboard',
  columns: [
    new ColumnDefinition({ field: 'teacherName', header: 'Nama Guru' }),
    new ColumnDefinition({ field: 'totalSessions', header: 'Total Sesi', type: 'number' }),
    new ColumnDefinition({ field: 'onTime', header: 'Tepat Waktu', type: 'number' }),
    new ColumnDefinition({ field: 'late', header: 'Terlambat', type: 'number' }),
    new ColumnDefinition({ field: 'totalMinutes', header: 'Total Menit', type: 'number' }),
  ],
  filters: [],
  defaultSort: [new SortDefinition({ field: 'totalSessions', direction: 'desc' })],
  aggregates: [
    new AggregateDefinition({
      field: 'totalSessions',
      method: 'sum',
      label: 'Total Sesi Semua Guru',
    }),
    new AggregateDefinition({ field: 'totalMinutes', method: 'average', label: 'Rata-rata Menit' }),
  ],
})
