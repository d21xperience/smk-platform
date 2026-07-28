// src/reports/StudentAttendanceReport.js
import {
  ReportDefinition,
  ColumnDefinition,
  FilterDefinition,
  SortDefinition,
  AggregateDefinition,
  FooterDefinition,
} from 'src/models/ReportDefinition'

export const studentAttendanceReport = new ReportDefinition({
  id: 'student-attendance',
  name: 'Laporan Absensi Siswa',
  projectionSource: 'studentProgress',
  columns: [
    new ColumnDefinition({ field: 'studentName', header: 'Nama Siswa', width: '200px' }),
    new ColumnDefinition({ field: 'className', header: 'Kelas' }),
    new ColumnDefinition({ field: 'attendanceSummary.hadir', header: 'Hadir', type: 'number' }),
    new ColumnDefinition({ field: 'attendanceSummary.izin', header: 'Izin', type: 'number' }),
    new ColumnDefinition({ field: 'attendanceSummary.sakit', header: 'Sakit', type: 'number' }),
    new ColumnDefinition({ field: 'attendanceSummary.alpha', header: 'Alpha', type: 'number' }),
    new ColumnDefinition({
      field: 'attendanceSummary.persentaseKehadiran',
      header: 'Persentase',
      type: 'percentage',
    }),
  ],
  filters: [
    new FilterDefinition({ field: 'className', label: 'Kelas', type: 'select', options: [] }),
    new FilterDefinition({ field: 'studentName', label: 'Nama Siswa', type: 'text' }),
  ],
  defaultSort: [new SortDefinition({ field: 'studentName', direction: 'asc' })],
  aggregates: [],
  footer: new FooterDefinition({
    aggregates: [
      new AggregateDefinition({
        field: 'attendanceSummary.hadir',
        method: 'sum',
        label: 'Total Hadir',
      }),
      new AggregateDefinition({
        field: 'attendanceSummary.izin',
        method: 'sum',
        label: 'Total Izin',
      }),
    ],
  }),
})
