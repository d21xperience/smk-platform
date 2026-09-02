export const HomeroomBillingContract = {
  fetchBillingData: {
    request: {
      classId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
      teacherId: 'string',
    },
    response: {
      classInfo: {
        classId: 'string',
        className: 'string',
        academicYearId: 'string',
        academicYearName: 'string',
        semesterId: 'string',
        semesterName: 'string',
      },
      bills: [
        {
          studentId: 'string',
          studentName: 'string',
          billId: 'string',
          billType: 'string',
          billName: 'string',
          amount: 'number',
          paidAmount: 'number',
          outstandingAmount: 'number',
          status: 'string',
          dueDate: 'string',
        },
      ],
    },
  },
}
