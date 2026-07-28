

// src/models/AcademicYear.js
export class AcademicYear {
  constructor({ id, name, isActive }) {
    this.id = id
    this.name = name        // e.g. "2025/2026"
    this.isActive = isActive
  }
}
