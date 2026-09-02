// apps/siakad-tu/src/domain/student/Student.js

import { NISN } from './value-objects/NISN.js'
import { NIS } from './value-objects/NIS.js'
import { FullName } from './value-objects/FullName.js'
import { BirthDate } from './value-objects/BirthDate.js'
import { Gender } from './value-objects/Gender.js'
import { Address } from './value-objects/Address.js'
import { ContactInfo } from './value-objects/ContactInfo.js'
import { GuardianInfo } from './value-objects/GuardianInfo.js'
import { StudentStatus } from './value-objects/StudentStatus.js'
import { Enrollment } from './Enrollment.js'
import {
  StudentRegistered,
  StudentEnrolled,
  StudentTransferred,
  StudentGraduated,
  StudentProfileUpdated,
} from './events/index.js' // index.js yang export semua events

/**
 * Aggregate Root: Student
 *
 * Entry point untuk semua perubahan pada data siswa.
 * Menjaga invariants, menghasilkan domain events.
 *
 * INVARIANTS:
 * 1. NISN unik dan valid (10 digit) — dijaga oleh Value Object
 * 2. NIS unik dalam satu sekolah — dijaga oleh Repository
 * 3. Siswa tidak bisa di-enroll 2x di periode yang sama
 * 4. Siswa "graduated/transferred/dropped" tidak bisa di-enroll lagi
 * 5. Siswa "active" harus punya minimal 1 enrollment aktif
 * 6. Enrollment terikat pada OperationalContext (schoolId + periodId)
 */
export class Student {
  #studentId
  #schoolId // dari OperationalContext saat registrasi
  #nisn
  #nis
  #fullName
  #birthDate
  #gender
  #address
  #contactInfo
  #guardianInfo
  #status
  #enrollments // Entity[]
  #events // DomainEvent[]

  constructor({
    studentId,
    schoolId,
    nisn,
    nis,
    fullName,
    birthDate,
    gender,
    address,
    contactInfo,
    guardianInfo,
    status = StudentStatus.ACTIVE,
    enrollments = [],
  }) {
    if (!studentId) throw new Error('Student: studentId wajib.')
    if (!schoolId) throw new Error('Student: schoolId wajib (dari OperationalContext).')

    this.#studentId = studentId
    this.#schoolId = schoolId
    this.#nisn = nisn instanceof NISN ? nisn : new NISN(nisn)
    this.#nis = nis instanceof NIS ? nis : new NIS(nis)
    this.#fullName =
      fullName instanceof FullName
        ? fullName
        : new FullName(fullName.firstName, fullName.middleName, fullName.lastName)
    this.#birthDate = birthDate instanceof BirthDate ? birthDate : new BirthDate(birthDate)
    this.#gender = gender instanceof Gender ? gender : gender
    this.#address = address instanceof Address ? address : new Address(address)
    this.#contactInfo =
      contactInfo instanceof ContactInfo ? contactInfo : new ContactInfo(contactInfo)
    this.#guardianInfo =
      guardianInfo instanceof GuardianInfo ? guardianInfo : new GuardianInfo(guardianInfo)
    this.#status = status instanceof StudentStatus ? status : status
    this.#enrollments = enrollments.map((e) =>
      e instanceof Enrollment ? e : Enrollment.fromJSON(e),
    )
    this.#events = []
  }

  // === GETTERS ===
  get studentId() {
    return this.#studentId
  }
  get schoolId() {
    return this.#schoolId
  }
  get nisn() {
    return this.#nisn
  }
  get nis() {
    return this.#nis
  }
  get fullName() {
    return this.#fullName
  }
  get birthDate() {
    return this.#birthDate
  }
  get gender() {
    return this.#gender
  }
  get address() {
    return this.#address
  }
  get contactInfo() {
    return this.#contactInfo
  }
  get guardianInfo() {
    return this.#guardianInfo
  }
  get status() {
    return this.#status
  }
  get enrollments() {
    return [...this.#enrollments]
  }

  // === BEHAVIORS ===

  /**
   * Register siswa baru
   * Dipanggil saat siswa pertama kali didaftarkan di sekolah
   */
  static register({
    studentId,
    schoolId,
    nisn,
    nis,
    fullName,
    birthDate,
    gender,
    address,
    contactInfo,
    guardianInfo,
  }) {
    const student = new Student({
      studentId,
      schoolId,
      nisn,
      nis,
      fullName,
      birthDate,
      gender,
      address,
      contactInfo,
      guardianInfo,
      status: StudentStatus.ACTIVE,
    })

    // Record event
    student.#events.push(
      new StudentRegistered(studentId, schoolId, {
        nisn: student.#nisn.value,
        nis: student.#nis.value,
        fullName: student.#fullName.displayName,
      }),
    )

    return student
  }

  /**
   * Enroll siswa ke kelas pada periode tertentu
   *
   * INVARIANTS:
   * - Status harus ACTIVE
   * - Tidak boleh ada enrollment aktif di periode yang sama
   */
  enroll({ enrollmentId, periodId, classId, enrollmentDate }) {
    // Invariant 4: Hanya ACTIVE yang bisa di-enroll
    if (!this.#status.canBeEnrolled()) {
      throw new Error(`Siswa dengan status "${this.#status.label}" tidak bisa di-enroll.`)
    }

    // Invariant 3: Tidak boleh ada enrollment aktif di periode yang sama
    const existingActive = this.#enrollments.find((e) => e.isForPeriod(periodId) && e.isActive())
    if (existingActive) {
      throw new Error(`Siswa sudah memiliki enrollment aktif di periode ${periodId}.`)
    }

    const enrollment = new Enrollment({
      enrollmentId,
      schoolId: this.#schoolId,
      periodId,
      classId,
      enrollmentDate,
    })

    this.#enrollments.push(enrollment)

    // Record event
    this.#events.push(
      new StudentEnrolled(this.#studentId, this.#schoolId, {
        enrollmentId,
        periodId,
        classId,
      }),
    )
  }

  /**
   * Update profil siswa (address, contact, guardian)
   */
  updateProfile({ address, contactInfo, guardianInfo }) {
    if (address) {
      this.#address = address instanceof Address ? address : new Address(address)
    }
    if (contactInfo) {
      this.#contactInfo =
        contactInfo instanceof ContactInfo ? contactInfo : new ContactInfo(contactInfo)
    }
    if (guardianInfo) {
      this.#guardianInfo =
        guardianInfo instanceof GuardianInfo ? guardianInfo : new GuardianInfo(guardianInfo)
    }

    this.#events.push(
      new StudentProfileUpdated(this.#studentId, this.#schoolId, {
        updatedFields: Object.keys({ address, contactInfo, guardianInfo }).filter(
          (k) =>
            [address, contactInfo, guardianInfo][
              ['address', 'contactInfo', 'guardianInfo'].indexOf(k)
            ],
        ),
      }),
    )
  }

  /**
   * Tandai siswa sebagai lulus
   */
  graduate() {
    if (!this.#status.isActive()) {
      throw new Error(`Siswa dengan status "${this.#status.label}" tidak bisa diluluskan.`)
    }

    // Complete semua enrollment aktif
    this.#enrollments.filter((e) => e.isActive()).forEach((e) => e.complete())

    this.#status = StudentStatus.GRADUATED

    this.#events.push(
      new StudentGraduated(this.#studentId, this.#schoolId, {
        graduationDate: new Date().toISOString(),
      }),
    )
  }

  /**
   * Tandai siswa sebagai pindah
   */
  transfer({ reason, targetSchool }) {
    if (!this.#status.isActive()) {
      throw new Error(`Siswa dengan status "${this.#status.label}" tidak bisa dipindahkan.`)
    }

    this.#enrollments.filter((e) => e.isActive()).forEach((e) => e.complete())

    this.#status = StudentStatus.TRANSFERRED

    this.#events.push(
      new StudentTransferred(this.#studentId, this.#schoolId, {
        reason,
        targetSchool,
        transferDate: new Date().toISOString(),
      }),
    )
  }

  /**
   * Cek apakah siswa memiliki enrollment aktif di periode tertentu
   */
  hasActiveEnrollmentInPeriod(periodId) {
    return this.#enrollments.some((e) => e.isForPeriod(periodId) && e.isActive())
  }

  /**
   * Dapatkan enrollment aktif saat ini (terbaru)
   */
  getCurrentEnrollment() {
    return (
      this.#enrollments
        .filter((e) => e.isActive())
        .sort((a, b) => b.enrollmentDate - a.enrollmentDate)[0] || null
    )
  }

  // === EVENT MANAGEMENT ===

  getUncommittedEvents() {
    return [...this.#events]
  }

  clearEvents() {
    this.#events = []
  }

  // === SERIALIZATION ===

  toJSON() {
    return {
      studentId: this.#studentId,
      schoolId: this.#schoolId,
      nisn: this.#nisn.toJSON(),
      nis: this.#nis.toJSON(),
      fullName: this.#fullName.toJSON(),
      birthDate: this.#birthDate.toJSON(),
      gender: this.#gender.toJSON(),
      address: this.#address.toJSON(),
      contactInfo: this.#contactInfo.toJSON(),
      guardianInfo: this.#guardianInfo.toJSON(),
      status: this.#status.toJSON(),
      enrollments: this.#enrollments.map((e) => e.toJSON()),
    }
  }

  static fromJSON(data) {
    return new Student({
      studentId: data.studentId,
      schoolId: data.schoolId,
      nisn: NISN.fromJSON(data.nisn),
      nis: NIS.fromJSON(data.nis),
      fullName: FullName.fromJSON(data.fullName),
      birthDate: BirthDate.fromJSON(data.birthDate),
      gender: Gender.fromJSON(data.gender),
      address: Address.fromJSON(data.address),
      contactInfo: ContactInfo.fromJSON(data.contactInfo),
      guardianInfo: GuardianInfo.fromJSON(data.guardianInfo),
      status: StudentStatus.fromJSON(data.status),
      enrollments: data.enrollments || [],
    })
  }
}
