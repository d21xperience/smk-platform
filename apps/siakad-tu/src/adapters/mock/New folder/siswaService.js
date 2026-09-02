const dummySiswa = [
  { id: 1, nis: '1001', nama: 'Ahmad Fauzi', kelasId: 1, tahunAjaranId: 1 },
  { id: 2, nis: '1002', nama: 'Budi Santoso', kelasId: 1, tahunAjaranId: 1 },
]
const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

export const siswaService = {
  async getAll() {
    await delay(300)
    return [...dummySiswa]
  },
  async getById(id) {
    await delay(200)
    return dummySiswa.find((s) => s.id === id)
  },
  async create(data) {
    await delay(400)
    const newId = dummySiswa.length + 1
    const newSiswa = { id: newId, ...data }
    dummySiswa.push(newSiswa)
    return newSiswa
  },
  async update(id, data) {
    await delay(400)
    const index = dummySiswa.findIndex((s) => s.id === id)
    if (index !== -1) dummySiswa[index] = { ...dummySiswa[index], ...data }
    return dummySiswa[index]
  },
  async delete(id) {
    await delay(400)
    const index = dummySiswa.findIndex((s) => s.id === id)
    if (index !== -1) dummySiswa.splice(index, 1)
    return true
  },
}
