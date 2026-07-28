export class JournalEntry {
  constructor({
    id,
    sessionId,
    date,
    subject,
    className,
    material, // materi yang diajarkan
    attendance, // jumlah siswa hadir / total
    obstacles, // kendala
    notes, // catatan tambahan
  }) {
    this.id = id
    this.sessionId = sessionId
    this.date = date
    this.subject = subject
    this.className = className
    this.material = material
    this.attendance = attendance
    this.obstacles = obstacles
    this.notes = notes
  }
}
