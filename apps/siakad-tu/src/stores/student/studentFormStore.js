// apps/siakad-tu/src/stores/student/studentFormStore.js

import { defineStore } from 'pinia'

/**
 * studentFormStore — Mengelola state form create/edit siswa.
 *
 * STATE ONLY. Validasi ada di Engine.
 *
 * Mode:
 * - 'create': form kosong untuk siswa baru
 * - 'edit': form terisi dengan data siswa existing
 */
export const useStudentFormStore = defineStore('studentForm', {
  state: () => ({
    /** @type {'create'|'edit'|null} Mode form */
    mode: null,
    /** @type {Object|null} Data form saat ini */
    formData: null,
    /** @type {string|null} studentId yang sedang diedit (mode 'edit') */
    editingStudentId: null,
    /** @type {boolean} Loading state (saat submit) */
    loading: false,
    /** @type {Object|null} Error dari submit */
    error: null,
    /** @type {Array} Validation errors dari Engine */
    validationErrors: [],
    /** @type {boolean} Apakah form sudah di-submit dan sukses */
    submitted: false,
    /** @type {boolean} Apakah form dirty (ada perubahan) */
    dirty: false,
  }),

  getters: {
    /**
     * Apakah form dalam mode create
     */
    isCreateMode: (state) => state.mode === 'create',

    /**
     * Apakah form dalam mode edit
     */
    isEditMode: (state) => state.mode === 'edit',

    /**
     * Default form data untuk create mode
     */
    defaultFormData: () => ({
      studentId: '',
      nisn: '',
      nis: '',
      fullName: { firstName: '', middleName: '', lastName: '' },
      birthDate: '',
      gender: '',
      address: {
        street: '',
        rtRw: '',
        village: '',
        district: '',
        city: '',
        postalCode: '',
      },
      contactInfo: { phone: '', email: '' },
      guardianInfo: {
        name: '',
        relation: '',
        phone: '',
        occupation: '',
      },
    }),

    /**
     * Apakah form bisa di-submit
     */
    canSubmit: (state) => {
      return state.mode !== null && !state.loading && state.dirty
    },
  },

  actions: {
    /**
     * Inisialisasi form untuk mode create
     */
    initCreate() {
      this.mode = 'create'
      this.formData = { ...this.defaultFormData }
      this.editingStudentId = null
      this.error = null
      this.validationErrors = []
      this.submitted = false
      this.dirty = false
    },

    /**
     * Inisialisasi form untuk mode edit dengan data existing
     */
    initEdit(studentData) {
      if (!studentData) {
        throw new Error('studentFormStore.initEdit: studentData wajib.')
      }

      this.mode = 'edit'
      this.editingStudentId = studentData.studentId
      this.formData = {
        studentId: studentData.studentId,
        nisn: studentData.nisn,
        nis: studentData.nis,
        fullName: { ...studentData.fullName },
        birthDate: studentData.birthDate ? studentData.birthDate.split('T')[0] : '',
        gender: studentData.gender,
        address: { ...studentData.address },
        contactInfo: { ...studentData.contactInfo },
        guardianInfo: { ...studentData.guardianInfo },
      }
      this.error = null
      this.validationErrors = []
      this.submitted = false
      this.dirty = false
    },

    /**
     * Update field di form data
     * Support nested path dengan dot notation: 'fullName.firstName'
     */
    updateField(path, value) {
      const keys = path.split('.')
      let target = this.formData

      for (let i = 0; i < keys.length - 1; i++) {
        target = target[keys[i]]
      }

      target[keys[keys.length - 1]] = value
      this.dirty = true
    },

    /**
     * Set loading state
     */
    setLoading(loading) {
      this.loading = loading
      if (loading) {
        this.error = null
        this.validationErrors = []
      }
    },

    /**
     * Set error dari submit
     */
    setError(error) {
      this.error = error
      this.loading = false
    },

    /**
     * Set validation errors dari Engine
     */
    setValidationErrors(errors) {
      this.validationErrors = errors
      this.loading = false
    },

    /**
     * Tandai form sebagai submitted (sukses)
     */
    markSubmitted() {
      this.submitted = true
      this.loading = false
      this.dirty = false
    },

    /**
     * Reset form (kembali ke state awal)
     */
    reset() {
      this.mode = null
      this.formData = null
      this.editingStudentId = null
      this.loading = false
      this.error = null
      this.validationErrors = []
      this.submitted = false
      this.dirty = false
    },
  },
})
