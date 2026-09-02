const MOCK_HOMEROOM_ASSIGNMENTS = {
  'USR-001': { classId: 'CLS-X-A', className: 'X-A' },
  'USR-002': { classId: 'CLS-XI-B', className: 'XI-B' },
}
const MOCK_BILLING_DATA = {
  'CLS-X-A': {
    classInfo: {
      classId: 'CLS-X-A',
      className: 'X-A',
      academicYearId: 'AY-2026',
      academicYearName: '2026/2027',
      semesterId: '20261',
      semesterName: 'Semester 1',
    },
    bills: [
      {
        studentId: 'STU-001',
        studentName: 'Ahmad Fauzi',
        billId: 'BILL-001',
        billType: 'spp',
        billName: 'SPP Juli 2026',
        amount: 500000,
        paidAmount: 500000,
        outstandingAmount: 0,
        status: 'paid',
        dueDate: '2026-07-10',
      },
      {
        studentId: 'STU-001',
        studentName: 'Ahmad Fauzi',
        billId: 'BILL-002',
        billType: 'spp',
        billName: 'SPP Agustus 2026',
        amount: 500000,
        paidAmount: 250000,
        outstandingAmount: 250000,
        status: 'partial',
        dueDate: '2026-08-10',
      },
      {
        studentId: 'STU-002',
        studentName: 'Budi Hartono',
        billId: 'BILL-003',
        billType: 'spp',
        billName: 'SPP Juli 2026',
        amount: 500000,
        paidAmount: 0,
        outstandingAmount: 500000,
        status: 'unpaid',
        dueDate: '2026-07-10',
      },
      {
        studentId: 'STU-002',
        studentName: 'Budi Hartono',
        billId: 'BILL-004',
        billType: 'building',
        billName: 'Uang Gedung',
        amount: 3000000,
        paidAmount: 1000000,
        outstandingAmount: 2000000,
        status: 'partial',
        dueDate: '2026-09-01',
      },
      {
        studentId: 'STU-003',
        studentName: 'Citra Dewi',
        billId: 'BILL-005',
        billType: 'spp',
        billName: 'SPP Juli 2026',
        amount: 500000,
        paidAmount: 500000,
        outstandingAmount: 0,
        status: 'paid',
        dueDate: '2026-07-10',
      },
      {
        studentId: 'STU-003',
        studentName: 'Citra Dewi',
        billId: 'BILL-006',
        billType: 'activity',
        billName: 'Kegiatan Study Tour',
        amount: 1500000,
        paidAmount: 0,
        outstandingAmount: 1500000,
        status: 'overdue',
        dueDate: '2026-06-15',
      },
      {
        studentId: 'STU-004',
        studentName: 'Diana Putri',
        billId: 'BILL-007',
        billType: 'spp',
        billName: 'SPP Juli 2026',
        amount: 500000,
        paidAmount: 0,
        outstandingAmount: 0,
        status: 'waived',
        dueDate: '2026-07-10',
      },
      {
        studentId: 'STU-005',
        studentName: 'Eko Prasetyo',
        billId: 'BILL-008',
        billType: 'spp',
        billName: 'SPP Juli 2026',
        amount: 500000,
        paidAmount: 500000,
        outstandingAmount: 0,
        status: 'paid',
        dueDate: '2026-07-10',
      },
      {
        studentId: 'STU-005',
        studentName: 'Eko Prasetyo',
        billId: 'BILL-009',
        billType: 'uniform',
        billName: 'Seragam Olahraga',
        amount: 750000,
        paidAmount: 750000,
        outstandingAmount: 0,
        status: 'paid',
        dueDate: '2026-07-20',
      },
      {
        studentId: 'STU-006',
        studentName: 'Fitri Handayani',
        billId: 'BILL-010',
        billType: 'spp',
        billName: 'SPP Juli 2026',
        amount: 500000,
        paidAmount: 0,
        outstandingAmount: 500000,
        status: 'unpaid',
        dueDate: '2026-07-10',
      },
    ],
  },
  'CLS-XI-B': {
    classInfo: {
      classId: 'CLS-XI-B',
      className: 'XI-B',
      academicYearId: 'AY-2026',
      academicYearName: '2026/2027',
      semesterId: '20261',
      semesterName: 'Semester 1',
    },
    bills: [
      {
        studentId: 'STU-011',
        studentName: 'Kartika Sari',
        billId: 'BILL-011',
        billType: 'spp',
        billName: 'SPP Juli 2026',
        amount: 600000,
        paidAmount: 600000,
        outstandingAmount: 0,
        status: 'paid',
        dueDate: '2026-07-10',
      },
      {
        studentId: 'STU-012',
        studentName: 'Lukman Hakim',
        billId: 'BILL-012',
        billType: 'spp',
        billName: 'SPP Juli 2026',
        amount: 600000,
        paidAmount: 300000,
        outstandingAmount: 300000,
        status: 'partial',
        dueDate: '2026-07-10',
      },
      {
        studentId: 'STU-013',
        studentName: 'Maya Anggraini',
        billId: 'BILL-013',
        billType: 'spp',
        billName: 'SPP Juli 2026',
        amount: 600000,
        paidAmount: 0,
        outstandingAmount: 600000,
        status: 'unpaid',
        dueDate: '2026-07-10',
      },
    ],
  },
}
export class HomeroomBillingMockAdapter {
  async fetchHomeroomAssignment({ teacherId }) {
    await new Promise((resolve) => setTimeout(resolve, 300))
    const assignment = MOCK_HOMEROOM_ASSIGNMENTS[teacherId]
    if (!assignment) {
      return null
    }
    return {
      teacherId,
      classId: assignment.classId,
      className: assignment.className,
    }
  }
  async fetchBillingData({ classId, academicYearId, semesterId, teacherId }) {
    await new Promise((resolve) => setTimeout(resolve, 600))
    const assignment = MOCK_HOMEROOM_ASSIGNMENTS[teacherId]
    if (!assignment) {
      throw new Error('No homeroom assignment found for this teacher')
    }
    if (assignment.classId !== classId) {
      throw new Error('You do not have permission to view billing data for this class')
    }
    const billingData = MOCK_BILLING_DATA[classId]
    if (!billingData) {
      throw new Error(`No billing data found for class: ${classId}`)
    }
    if (billingData.classInfo.academicYearId !== academicYearId) {
      throw new Error('Billing data not available for the selected academic year')
    }
    if (billingData.classInfo.semesterId !== semesterId) {
      throw new Error('Billing data not available for the selected semester')
    }
    return {
      classInfo: billingData.classInfo,
      bills: billingData.bills,
    }
  }
}
