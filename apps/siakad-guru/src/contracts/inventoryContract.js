export const InventoryContract = {
  fetchItemsByClass: {
    request: {
      classId: 'string',
      schoolId: 'string',
    },
    response: [
      {
        id: 'string',
        name: 'string',
        category: 'string',
        location: 'string',
        classId: 'string',
        className: 'string',
        hasQRCode: 'boolean',
        schoolId: 'string',
      },
    ],
  },

  fetchItemById: {
    request: {
      itemId: 'string',
    },
    response: {
      id: 'string',
      name: 'string',
      category: 'string',
      location: 'string',
      classId: 'string',
      className: 'string',
      hasQRCode: 'boolean',
      schoolId: 'string',
    },
  },

  submitDamageReport: {
    request: {
      itemId: 'string',
      itemName: 'string',
      condition: 'string', // good, damaged, broken, missing
      description: 'string',
      reportMethod: 'string', // qr, manual
      date: 'string',
      teacherId: 'string',
      teacherName: 'string',
      classId: 'string',
      className: 'string',
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
    },
    response: {
      id: 'string',
      itemId: 'string',
      itemName: 'string',
      condition: 'string',
      description: 'string',
      reportMethod: 'string',
      date: 'string',
      teacherId: 'string',
      teacherName: 'string',
      classId: 'string',
      className: 'string',
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
    },
  },
}
