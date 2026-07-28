import { defineStore } from 'pinia'
import { AttendanceEngine } from '@/engine/AttendanceEngine'
import { AttendanceRecord } from '@/models/Attendance'

export const useAttendanceStore = defineStore('attendance', {
  state: () => ({
    // Service diinjeksi nanti
    attendanceService: null,

    // Record absensi yang sedang aktif (draft atau submitted)
    currentRecord: null, // AttendanceRecord

    // Timestamp terakhir autosave (untuk indikator UI)
    lastAutoSave: null,

    // Status operasi
    loading: {
      init: false,
      submit: false,
    },
    saving: false, // true saat autosave/manual save berlangsung
    error: null,
  }),

  getters: {
    // Daftar siswa di record aktif
    students: (state) => state.currentRecord?.students || [],

    // Apakah record bisa diedit (hanya draft)
    canEdit: (state) => {
      if (!state.currentRecord) return false
      return AttendanceEngine.canEdit(state.currentRecord.status)
    },

    // Apakah semua siswa sudah ditandai
    isAllMarked: (state) => {
      if (!state.currentRecord?.students) return false
      return AttendanceEngine.isAllMarked(state.currentRecord.students)
    },

    // Ringkasan absensi (total, hadir, izin, dll.)
    summary: (state) => {
      if (!state.currentRecord?.students) return null
      return AttendanceEngine.calculateSummary(state.currentRecord.students)
    },

    // Validasi untuk submit
    submitValidation: (state) => {
      if (!state.currentRecord) return { valid: false, errors: ['Belum ada data absensi.'] }
      return AttendanceEngine.validateForSubmit(state.currentRecord)
    },

    // Status record
    recordStatus: (state) => state.currentRecord?.status || null,
  },

  actions: {
    /**
     * Injeksi AttendanceService.
     */
    setService(service) {
      this.attendanceService = service
    },

    // --- Inisialisasi ---

    /**
     * Load atau buat record absensi untuk sesi tertentu.
     * Dipanggil saat guru membuka halaman absensi.
     * @param {Object} session - data sesi dari TeachingSession
     */
    async loadOrCreateAttendance(session) {
      if (!this.attendanceService) throw new Error('AttendanceService belum diinisialisasi')
      if (!session?.id) throw new Error('Data sesi tidak valid')

      this.loading.init = true
      this.error = null

      try {
        // Coba ambil record yang sudah ada
        let record = await this.attendanceService.fetchRecord(session.id)

        if (record) {
          // Sudah ada draft/submitted
          this.currentRecord = record
        } else {
          // Belum ada: ambil daftar siswa, buat record baru
          const rawStudents = await this.attendanceService.fetchStudents(
            session.id,
            session.className,
          )
          console.log('Attendan.store', rawStudents)

          // const studentItems = rawStudents.map((s) => AttendanceEngine.createStudentItem(s))
          const studentItems = rawStudents.map((s) => AttendanceEngine.createStudentItem(s))
          console.log('Attendan.store', studentItems)

          this.currentRecord = new AttendanceRecord({
            id: 0,
            sessionId: session.id,
            date: session.date,
            className: session.className,
            subject: session.subject,
            students: studentItems,
            status: 'draft',
          })
        }
      } catch (err) {
        this.error = err.message || 'Gagal memuat data absensi'
        throw err
      } finally {
        this.loading.init = false
      }
    },

    // --- Manipulasi Status Siswa ---

    /**
     * Toggle status siswa (hadir/izin/sakit/alpha/unmarked).
     * @param {number|string} studentId
     * @param {string} newStatus
     */
    updateStudentStatus(studentId, newStatus) {
      if (!this.currentRecord) return
      if (!AttendanceEngine.canEdit(this.currentRecord.status)) return

      const student = this.currentRecord.students.find((s) => s.studentId === studentId)
      if (!student) return

      // Gunakan engine untuk toggle
      student.status = AttendanceEngine.toggleStudentStatus(student, newStatus)

      // Trigger autosave
      this._scheduleAutoSave()
    },

    /**
     * Tandai semua siswa unmarked dengan status tertentu.
     * @param {string} status
     */
    markAllAs(status) {
      if (!this.currentRecord) return
      if (!AttendanceEngine.canEdit(this.currentRecord.status)) return

      this.currentRecord.students = AttendanceEngine.markAllUnmarkedAs(
        this.currentRecord.students,
        status,
      )

      // Trigger autosave
      this._scheduleAutoSave()
    },

    /**
     * Update catatan pada siswa tertentu.
     * @param {number|string} studentId
     * @param {string} note
     */
    updateStudentNote(studentId, note) {
      if (!this.currentRecord) return
      if (!AttendanceEngine.canEdit(this.currentRecord.status)) return

      const student = this.currentRecord.students.find((s) => s.studentId === studentId)
      if (!student) return

      student.note = note
      this._scheduleAutoSave()
    },

    // --- Autosave ---

    // Timer untuk debounce autosave
    _autoSaveTimer: null,

    /**
     * Jadwalkan autosave setelah delay 2 detik.
     * Jika ada perubahan baru, timer di-reset.
     */
    _scheduleAutoSave() {
      if (this._autoSaveTimer) clearTimeout(this._autoSaveTimer)
      this._autoSaveTimer = setTimeout(() => {
        this._executeAutoSave()
      }, 2000)
    },

    /**
     * Eksekusi autosave draft.
     */
    async _executeAutoSave() {
      if (!this.currentRecord) return
      if (this.currentRecord.status !== 'draft') return
      if (!this.attendanceService) return

      this.saving = true
      try {
        const session = {
          id: this.currentRecord.sessionId,
          date: this.currentRecord.date,
          className: this.currentRecord.className,
          subject: this.currentRecord.subject,
        }
        const saved = await this.attendanceService.saveDraft(session, this.currentRecord.students)
        this.currentRecord = saved
        this.lastAutoSave = new Date().toISOString()
      } catch (err) {
        this.error = 'Gagal menyimpan draft otomatis'
        console.error('Autosave failed:', err)
      } finally {
        this.saving = false
      }
    },

    /**
     * Simpan draft secara manual (jika guru ingin simpan eksplisit).
     */
    async saveDraftManually() {
      if (!this.currentRecord) return
      if (this.currentRecord.status !== 'draft') return
      if (!this.attendanceService) throw new Error('Service tidak tersedia')

      // Batalkan timer autosave yang sedang berjalan
      if (this._autoSaveTimer) {
        clearTimeout(this._autoSaveTimer)
        this._autoSaveTimer = null
      }

      this.saving = true
      this.error = null
      try {
        const session = {
          id: this.currentRecord.sessionId,
          date: this.currentRecord.date,
          className: this.currentRecord.className,
          subject: this.currentRecord.subject,
        }
        const saved = await this.attendanceService.saveDraft(session, this.currentRecord.students)
        this.currentRecord = saved
        this.lastAutoSave = new Date().toISOString()
      } catch (err) {
        this.error = err.message || 'Gagal menyimpan draft'
        throw err
      } finally {
        this.saving = false
      }
    },

    // --- Submit Final ---

    /**
     * Submit absensi (final, tidak bisa diedit lagi).
     * @returns {Promise<Object>} hasil submit
     */
    async submitAttendance() {
      if (!this.currentRecord) throw new Error('Tidak ada data absensi')
      if (!this.attendanceService) throw new Error('Service tidak tersedia')

      // Validasi lewat engine
      const validation = AttendanceEngine.validateForSubmit(this.currentRecord)
      if (!validation.valid) {
        throw new Error(validation.errors.join(' '))
      }

      this.loading.submit = true
      this.error = null
      try {
        const submitted = await this.attendanceService.submit(this.currentRecord.sessionId)
        this.currentRecord = submitted
        return submitted
      } catch (err) {
        this.error = err.message || 'Gagal submit absensi'
        throw err
      } finally {
        this.loading.submit = false
      }
    },

    // --- Summary ---

    /**
     * Load ringkasan absensi (setelah submit atau untuk laporan).
     * @param {number} sessionId
     */
    async loadSummary(sessionId) {
      if (!this.attendanceService) throw new Error('Service tidak tersedia')

      this.loading.init = true
      this.error = null
      try {
        const data = await this.attendanceService.getSummary(sessionId)
        if (data) {
          this.currentRecord = new AttendanceRecord(data)
        }
      } catch (err) {
        this.error = err.message || 'Gagal memuat ringkasan'
      } finally {
        this.loading.init = false
      }
    },

    // --- Reset ---

    /**
     * Bersihkan state (saat guru keluar halaman absensi).
     */
    reset() {
      if (this._autoSaveTimer) {
        clearTimeout(this._autoSaveTimer)
        this._autoSaveTimer = null
      }
      this.currentRecord = null
      this.lastAutoSave = null
      this.loading = { init: false, submit: false }
      this.saving = false
      this.error = null
    },
  },
})
