


// src/models/AcademicPeriod.js
export class AcademicPeriod {
  constructor({ id, semesterId, name, startDate, endDate }) {
    this.id = id
    this.semesterId = semesterId
    this.name = name        // "Tengah Semester 1", "Akhir Semester 1"
    this.startDate = startDate
    this.endDate = endDate
  }
}
