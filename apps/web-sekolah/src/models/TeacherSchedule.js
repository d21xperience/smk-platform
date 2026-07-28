export class TeacherSchedule {
  constructor({ id, dayOfWeek, startTime, endTime, subject, className, room }) {
    this.id = id
    this.dayOfWeek = dayOfWeek // 'Senin', 'Selasa', ...
    this.startTime = startTime // '07:00'
    this.endTime = endTime // '09:30'
    this.subject = subject // nama mata pelajaran
    this.className = className // kelas yang diajar, misal 'XII RPL 1'
    this.room = room // ruangan (opsional)
  }
}
