// FILE: src/router/routes.js
// STATUS: MODIFY
// STATUS IMPLEMENTASI: COMPLETE

const routes = [
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    children: [
      {
        path: '',
        name: 'dashboard',
        component: () => import('@/pages/reporting/ReportingPage.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'academic-calender',
        name: 'academic-calender',
        component: () => import('@/pages/academic-calender/KalenderAkademik.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'teaching',
        name: 'teaching',
        component: () => import('@/pages/teaching/TeachingPage.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'teaching/:sessionId',
        name: 'teaching-session',
        component: () => import('@/pages/teaching/TeachingSessionPage.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'attendance',
        name: 'attendance',
        component: () => import('@/pages/attendance/AttendancePage.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'attendance/:teachingSessionId',
        name: 'attendance-session',
        component: () => import('@/pages/attendance/AttendanceSessionPage.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'journal',
        name: 'journal',
        component: () => import('@/pages/journal/JournalPage.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'journal/:teachingSessionId',
        name: 'journal-session',
        component: () => import('@/pages/journal/JournalSessionPage.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'assessment',
        name: 'assessment',
        component: () => import('@/pages/assessment/AssessmentPage.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'assessment/:assessmentId',
        name: 'assessment-session',
        component: () => import('@/pages/assessment/AssessmentSessionPage.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'progress',
        name: 'progress',
        component: () => import('@/pages/progress/ProgressPage.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'progress/:studentId',
        name: 'progress-student',
        component: () => import('@/pages/progress/ProgressStudentPage.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'inventory',
        name: 'inventory',
        component: () => import('@/pages/inventory/InventoryPage.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'inventory/:itemId',
        name: 'inventory-report',
        component: () => import('@/pages/inventory/InventoryReportPage.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'teacher-documents',
        name: 'teacher-documents',
        component: () => import('@/pages/teacher-document/TeacherDocumentPage.vue'),
        meta: {
          requiresAuth: true,
        },
      },
      {
        path: 'homeroom-billing',
        name: 'homeroom-billing',
        component: () => import('@/pages/homeroom-billing/HomeroomBillingPage.vue'),
        meta: {
          requiresAuth: true,
          requiredPermission: 'homeroom.billing.read',
        },
      },
      {
        path: 'homeroom-students',
        name: 'homeroom-students',
        component: () => import('@/pages/homeroom-student/HomeroomStudentListPage.vue'),
        meta: {
          requiresAuth: true,
          requiredPermission: 'homeroom.student.read',
        },
      },
      {
        path: 'homeroom-progress',
        name: 'homeroom-progress',
        component: () => import('@/pages/homeroom-progress/HomeroomProgressPage.vue'),
        meta: {
          requiresAuth: true,
          requiredPermission: 'homeroom.progress.read',
        },
      },
      {
        path: 'homeroom-attendance',
        name: 'homeroom-attendance',
        component: () => import('@/pages/homeroom-attendance/HomeroomAttendancePage.vue'),
        meta: {
          requiresAuth: true,
          requiredPermission: 'homeroom.attendance.read',
        },
      },
      {
        path: 'homeroom-journal',
        name: 'homeroom-journal',
        component: () => import('@/pages/homeroom-journal/HomeroomJournalPage.vue'),
        meta: {
          requiresAuth: true,
          requiredPermission: 'homeroom.journal.read',
        },
      },
      {
        path: 'homeroom-communication',
        name: 'homeroom-communication',
        component: () => import('@/pages/homeroom-communication/HomeroomCommunicationPage.vue'),
        meta: {
          requiresAuth: true,
          requiredPermission: 'homeroom.communication.read',
        },
      },
      {
        path: 'homeroom-print-report',
        name: 'homeroom-print-report',
        component: () => import('@/pages/homeroom-report/HomeroomPrintReportPage.vue'),
        meta: {
          requiresAuth: true,
          requiredPermission: 'homeroom.report.print',
        },
      },
    ],
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/pages/auth/LoginPage.vue'),
    meta: {
      requiresAuth: false,
    },
  },
  {
    path: '/error-server',
    name: 'ErrorServer',
    component: () => import('@/pages/ErrorServerPage.vue'),
    meta: {
      requiresAuth: false,
    },
  },
  {
    path: '/:catchAll(.*)*',
    component: () => import('@/pages/ErrorNotFound.vue'),
    meta: {
      requiresAuth: false,
    },
  },
]

export default routes
