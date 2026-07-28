// src/models/User.js
export class User {
  constructor({ id, name, role, nip }) {
    this.id = id
    this.name = name
    this.role = role    // 'guru', 'admin', dll.
    this.nip = nip
  }
}
