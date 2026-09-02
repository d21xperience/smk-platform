let dummyKategori = [
  {
    id: 1,
    nama: 'Uang Praktikum',
    deskripsi: 'Biaya praktikum per semester',
    nominalDefault: 150000,
  },
  { id: 2, nama: 'Uang Gedung', deskripsi: 'Iuran pembangunan gedung', nominalDefault: 500000 },
  { id: 3, nama: 'Ujian', deskripsi: 'Biaya ujian semester', nominalDefault: 100000 },
  { id: 4, nama: 'Kegiatan Sekolah', deskripsi: 'Donasi kegiatan sekolah', nominalDefault: 50000 },
]
let nextId = 5
const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

export const kategoriPenerimaanService = {
  async getAll() {
    await delay(200)
    return [...dummyKategori]
  },
  async create(data) {
    await delay(300)
    const newKategori = { id: nextId++, ...data }
    dummyKategori.push(newKategori)
    return newKategori
  },
  // async update(id, data) {
  //   /* update logic */
  // },
  // async delete(id) {
  //   /* delete logic */
  // },
}
