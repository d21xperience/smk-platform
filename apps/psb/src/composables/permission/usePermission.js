// src/composables/permission/usePermission.js

export function usePermission() {
  async function checkPagePermission(page, context) {
    // Logika sederhana: hanya guru yang bisa akses attendance
    if (page === 'attendance') {
      return context.user?.role === 'teacher'
    }
    return true
  }

  return {
    checkPagePermission,
  }
}
