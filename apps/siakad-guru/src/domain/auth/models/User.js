export class User {
  constructor({ id, name, email, role, permissions = [] }) {
    this.id = id
    this.name = name
    this.email = email
    this.role = role
    this.permissions = permissions
  }

  hasPermission(requiredPermission) {
    if (this.role === 'admin' || this.role === 'superadmin') {
      return true
    }
    return this.permissions.includes(requiredPermission)
  }
}
