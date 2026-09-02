let dummyKategori = [
  { id: 1, nama: 'ATK' },
  { id: 2, nama: 'Perbaikan Gedung' },
  { id: 3, nama: 'Internet' },
  { id: 4, nama: 'Listrik' },
  { id: 5, nama: 'Honor' },
]
let nextId = 6
const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

export const kategoriPengeluaranService = {
  async getAll() {
    await delay(200)
    return [...dummyKategori]
  },
  async create(data) {
    await delay(300)
    const baru = { id: nextId++, ...data }
    dummyKategori.push(baru)
    return baru
  },
  // eslint-disable-next-line no-unused-vars
  async update(id, data) {
    /* update */
  },
  // eslint-disable-next-line no-unused-vars
  async delete(id) {
    /* delete */
  },
}
