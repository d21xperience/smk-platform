// boot/services.js
import { attendanceService } from '@/services/AttendanceService'

// Registrasi sebagai plugin (opsional)
export default ({ app }) => {
  // Tidak menggunakan 'new'
  app.provide('attendanceService', attendanceService)
}
