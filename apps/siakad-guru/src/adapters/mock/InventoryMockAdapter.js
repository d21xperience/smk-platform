// Mock inventory items data
const mockInventoryItems = {
  'CLS-X-A': [
    {
      id: 'INV-001',
      name: 'Proyektor Epson X51',
      category: 'Elektronik',
      location: 'Plafon Depan',
      classId: 'CLS-X-A',
      className: 'X-A',
      hasQRCode: false,
      schoolId: 'SCH-001',
    },
    {
      id: 'INV-002',
      name: 'AC Daikin 1.5 PK',
      category: 'Elektronik',
      location: 'Dinding Kanan',
      classId: 'CLS-X-A',
      className: 'X-A',
      hasQRCode: false,
      schoolId: 'SCH-001',
    },
    {
      id: 'INV-003',
      name: 'Papan Tulis',
      category: 'Furnitur',
      location: 'Depan Kelas',
      classId: 'CLS-X-A',
      className: 'X-A',
      hasQRCode: true,
      schoolId: 'SCH-001',
    },
    {
      id: 'INV-004',
      name: 'Lampu LED 40W',
      category: 'Elektrikal',
      location: 'Plafon',
      classId: 'CLS-X-A',
      className: 'X-A',
      hasQRCode: false,
      schoolId: 'SCH-001',
    },
    {
      id: 'INV-005',
      name: 'Komputer Desktop HP',
      category: 'Elektronik',
      location: 'Meja Guru',
      classId: 'CLS-X-A',
      className: 'X-A',
      hasQRCode: true,
      schoolId: 'SCH-001',
    },
    {
      id: 'INV-006',
      name: 'Colokan Listrik',
      category: 'Elektrikal',
      location: 'Dinding Belakang',
      classId: 'CLS-X-A',
      className: 'X-A',
      hasQRCode: false,
      schoolId: 'SCH-001',
    },
  ],
  'CLS-XI-B': [
    {
      id: 'INV-007',
      name: 'Proyektor BenQ MH560',
      category: 'Elektronik',
      location: 'Plafon Depan',
      classId: 'CLS-XI-B',
      className: 'XI-B',
      hasQRCode: false,
      schoolId: 'SCH-001',
    },
    {
      id: 'INV-008',
      name: 'AC Panasonic 2 PK',
      category: 'Elektronik',
      location: 'Dinding Kiri',
      classId: 'CLS-XI-B',
      className: 'XI-B',
      hasQRCode: false,
      schoolId: 'SCH-001',
    },
    {
      id: 'INV-009',
      name: 'Speaker Aktif',
      category: 'Elektronik',
      location: 'Sudut Kelas',
      classId: 'CLS-XI-B',
      className: 'XI-B',
      hasQRCode: true,
      schoolId: 'SCH-001',
    },
  ],
}

// Mock damage reports storage
let mockDamageReports = []
let reportCounter = 1

export class InventoryMockAdapter {
  async fetchItemsByClass({ classId }) {
    await new Promise((resolve) => setTimeout(resolve, 500))
    return mockInventoryItems[classId] || []
  }

  async fetchItemById({ itemId }) {
    await new Promise((resolve) => setTimeout(resolve, 300))

    for (const classId in mockInventoryItems) {
      const item = mockInventoryItems[classId].find((i) => i.id === itemId)
      if (item) {
        return item
      }
    }

    return null
  }

  async submitDamageReport({
    itemId,
    itemName,
    condition,
    description,
    reportMethod,
    date,
    teacherId,
    teacherName,
    classId,
    className,
    schoolId,
    academicYearId,
    semesterId,
  }) {
    await new Promise((resolve) => setTimeout(resolve, 400))

    const newReport = {
      id: `RPT-${String(reportCounter++).padStart(3, '0')}`,
      itemId,
      itemName,
      condition,
      description,
      reportMethod,
      date,
      teacherId,
      teacherName,
      classId,
      className,
      schoolId,
      academicYearId,
      semesterId,
    }

    mockDamageReports.push(newReport)
    return newReport
  }
}
