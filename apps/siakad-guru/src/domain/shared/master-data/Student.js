// domain/shared/master-data/Student.js
import { Entity } from '../Entity.js'

export class Student extends Entity {
  constructor({
    id,
    nis,
    nisn,
    name,
    gender,
    birthPlace,
    birthDate,
    religion,
    address,
    status = 'ACTIVE',
    ...props
  }) {
    super({ id, ...props })
    this.nis = nis
    this.nisn = nisn // 10 digit numeric
    this.name = name
    this.gender = gender // 'L' or 'P'
    this.birthPlace = birthPlace
    this.birthDate = birthDate
    this.religion = religion
    this.address = address
    this.status = status // 'ACTIVE', 'GRADUATED', 'MOVED'
  }
}
