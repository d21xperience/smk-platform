// src/models/Semester.js
export class Semester {
  constructor({ id, academicYearId, name, type }) {
    this.id = id
    this.academicYearId = academicYearId
    this.name = name // "Semester 1", "Semester 2"
    this.type = type // "ganjil" / "genap"
  }
}
