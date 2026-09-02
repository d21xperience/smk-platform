import { defineStore, acceptHMRUpdate } from 'pinia'

export const useJadwalStore = defineStore('jadwal', {
  state: () => ({
    list: [],
  }),
  actions: {
    loadData() {
      const stored = localStorage.getItem('jadwal_mengajar')
      if (stored) {
        this.list = JSON.parse(stored)
      } else {
        // Data dummy
        this.list = [
          {
            id: 1,
            guru_id: 1,
            guru_nama: 'Budi Santoso',
            mapel: 'Matematika',
            kelas: '10 IPA 1',
            hari: 'Senin',
            jam_mulai: '07:30',
            jam_selesai: '09:00',
            ruang: 'Ruang 101',
          },
          {
            id: 2,
            guru_id: 2,
            guru_nama: 'Dewi Lestari',
            mapel: 'Bahasa Inggris',
            kelas: '10 IPA 1',
            hari: 'Senin',
            jam_mulai: '09:15',
            jam_selesai: '10:45',
            ruang: 'Ruang 102',
          },
        ]
        this.saveData()
      }
    },
    saveData() {
      localStorage.setItem('jadwal_mengajar', JSON.stringify(this.list))
    },
    tambah(jadwal) {
      jadwal.id = Date.now()
      this.list.push(jadwal)
      this.saveData()
    },
    update(id, updated) {
      const index = this.list.findIndex((j) => j.id === id)
      if (index !== -1) {
        this.list[index] = { ...this.list[index], ...updated }
        this.saveData()
      }
    },
    hapus(id) {
      this.list = this.list.filter((j) => j.id !== id)
      this.saveData()
    },
    // Cek konflik: return array konflik
    cekKonflik(jadwalBaru, excludeId = null) {
      const conflicts = []
      for (const j of this.list) {
        if (excludeId && j.id === excludeId) continue
        if (j.hari !== jadwalBaru.hari) continue
        // cek overlap jam
        const mulaiBaru = jadwalBaru.jam_mulai
        const selesaiBaru = jadwalBaru.jam_selesai
        const mulaiExist = j.jam_mulai
        const selesaiExist = j.jam_selesai
        const overlap = mulaiBaru < selesaiExist && selesaiBaru > mulaiExist
        if (overlap) {
          if (j.guru_id === jadwalBaru.guru_id) {
            conflicts.push(
              `Guru ${j.guru_nama} sudah mengajar ${j.mapel} di ${j.kelas} pada jam yang sama.`,
            )
          }
          if (j.ruang === jadwalBaru.ruang) {
            conflicts.push(
              `Ruang ${j.ruang} sudah digunakan untuk ${j.mapel} (${j.kelas}) pada jam tersebut.`,
            )
          }
          if (j.kelas === jadwalBaru.kelas) {
            conflicts.push(`Kelas ${j.kelas} sudah memiliki jadwal ${j.mapel} pada jam tersebut.`)
          }
        }
      }
      return conflicts
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useJadwalStore, import.meta.hot))
}
