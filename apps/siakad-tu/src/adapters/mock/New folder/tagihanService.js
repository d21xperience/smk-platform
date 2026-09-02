// services/kategoriTagihanService.js
let dummyKategori = [
  { id: 1, nama: 'SPP', nominalDefault: 200000, deskripsi: 'Sumbangan Pembinaan Pendidikan' },
  { id: 2, nama: 'Uang Gedung', nominalDefault: 500000, deskripsi: 'Iuran pembangunan gedung' },
  { id: 3, nama: 'Ujian', nominalDefault: 100000, deskripsi: 'Biaya ujian semester' },
  { id: 4, nama: 'Kegiatan Sekolah', nominalDefault: 50000, deskripsi: 'Kegiatan ekstrakurikuler' },
  { id: 5, nama: 'Praktikum', nominalDefault: 150000, deskripsi: 'Biaya praktikum laboratorium' },
]

let nextId = 6
const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

export const tagihanService = {
  async getAll() {
    await delay(200)
    return [...dummyKategori]
  },
  async getById(id) {
    await delay(100)
    return dummyKategori.find((k) => k.id === id)
  },
  async create(data) {
    await delay(300)
    const newKategori = { id: nextId++, ...data }
    dummyKategori.push(newKategori)
    return newKategori
  },
  async update(id, data) {
    await delay(300)
    const index = dummyKategori.findIndex((k) => k.id === id)
    if (index !== -1) {
      dummyKategori[index] = { ...dummyKategori[index], ...data }
      return dummyKategori[index]
    }
    throw new Error('Kategori not found')
  },
  async delete(id) {
    await delay(300)
    const index = dummyKategori.findIndex((k) => k.id === id)
    if (index !== -1) {
      dummyKategori.splice(index, 1)
      return true
    }
    throw new Error('Kategori not found')
  },
}
