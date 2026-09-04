const penilaianRoutes = {
  path: '/penilaian',
  meta: {
    requiresAuth: true,
    allowedRoles: ['admin', 'guru', 'kurikulum'],
  },
  children: [
    {
      path: '',
      redirect: '/penilaian/nilai',
    },
    // {
    //   path: 'nilai',
    //   name: 'penilaian-nilai',
    //   component: () => import('@/pages/penilaian/AssessmentListPage.vue'),
    // },
    // {
    //   path: 'nilai/input',
    //   name: 'penilaian-input',
    //   component: () => import('@/pages/penilaian/AssessmentFormPage.vue'),
    // },
    // {
    //   path: 'nilai/:id/edit',
    //   name: 'penilaian-edit',
    //   component: () => import('@/pages/penilaian/AssessmentFormPage.vue'),
    //   props: true,
    // },
    {
      path: 'rapor/:studentId',
      name: 'penilaian-rapor',
      component: () => import('@/pages/penilaian/ReportCardPage.vue'),
      props: true,
    },
    {
      path: 'rekap-kelas/:classId/:subjectId',
      name: 'penilaian-rekap-kelas',
      component: () => import('@/pages/penilaian/ClassGradeSummaryPage.vue'),
      props: true,
    },
  ],
}

export default penilaianRoutes
