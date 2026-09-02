// application/attendance/FetchTeacherClasses.js
import { AdapterRegistry } from '../../adapters/registry.js'
import { useContextStore } from '../../stores/context.store.js'

export const FetchTeacherClassesUseCase = {
  execute: async () => {
    console.log('[UseCase] FetchTeacherClasses.execute called')
    const contextStore = useContextStore()
    const { academicYearId, userId } = contextStore
    console.log('[UseCase] Context:', { academicYearId, userId })

    const classAdapter = AdapterRegistry.class
    if (!classAdapter) {
      console.error('[UseCase] Class adapter is null!')
      throw new Error('Class adapter not initialized')
    }
    console.log('[UseCase] classAdapter:', classAdapter)

    const classes = await classAdapter.getClassesByTeacher(userId, academicYearId)
    console.log('[UseCase] Raw classes from adapter:', classes)

    const filtered = classes.filter((cls) => cls.academicYearId === academicYearId)
    console.log('[UseCase] Filtered classes:', filtered)
    return filtered
  },
}
