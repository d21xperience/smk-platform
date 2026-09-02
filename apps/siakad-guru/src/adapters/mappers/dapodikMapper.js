// // adapters/mappers/dapodikMapper.js
// import { Student } from '../../domain/shared/master-data/Student.js'
// import { Class } from '../../domain/shared/master-data/Class.js'
// // import { Teacher } from '../../domain/shared/master-data/Teacher.js' // Belum dibuat, tapi sama patternnya

// export const dapodikMapper = {
//   // Map raw JSON dari Dapodik ke Domain Student
//   toStudent: (raw) => {
//     // Asumsikan raw dari Dapodik punya field 'peserta_didik_id', 'nama', dll.
//     return new Student({
//       id: raw.peserta_didik_id || raw.id, // Pastikan pakai UUID internal
//       nis: raw.nis || raw.nomor_induk,
//       nisn: raw.nisn,
//       name: raw.nama_lengkap || raw.nama,
//       gender: raw.jenis_kelamin === 'L' ? 'L' : 'P',
//       birthPlace: raw.tempat_lahir,
//       birthDate: raw.tanggal_lahir,
//       religion: raw.agama,
//       address: raw.alamat,
//       status: raw.status_aktif === 1 ? 'ACTIVE' : 'INACTIVE',
//       sourceSystem: 'DAPODIK',
//       metadata: { lastSync: new Date().toISOString() },
//     })
//   },

//   toClass: (raw) => {
//     return new Class({
//       id: raw.rombongan_belajar_id,
//       academicYearId: raw.tahun_ajaran_id,
//       grade: raw.tingkat_kelas, // 10, 11, 12
//       major: raw.jurusan_kode, // 'TKJ'
//       parallelClass: raw.rombel_ke, // 1, 2
//       homeroomTeacherId: raw.ptk_id,
//       roomId: raw.ruang_id,
//       sourceSystem: 'DAPODIK',
//     })
//   },

//   // toTeacher, toAcademicYear, dll.
// }
