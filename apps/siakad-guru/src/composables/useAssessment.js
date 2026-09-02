import { computed } from 'vue'
import { useAssessmentStore } from '../stores/assessmentStore.js'
import { useContextStore } from '../stores/contextStore.js'
import { useRouter } from 'vue-router'

export function useAssessment() {
  const assessmentStore = useAssessmentStore()
  const contextStore = useContextStore()
  const router = useRouter()

  const currentContext = computed(() => contextStore.currentContext)
  const assessment = computed(() => assessmentStore.assessment)
  const result = computed(() => assessmentStore.result)
  const isLoading = computed(() => assessmentStore.loading)
  const isSaving = computed(() => assessmentStore.saving)
  const error = computed(() => assessmentStore.error)
  const isDraft = computed(() => assessmentStore.isDraft)
  const isFinalized = computed(() => assessmentStore.isFinalized)

  const loadAssessment = async ({ classId, subjectId }) => {
    console.log('👀', classId)
    if (!currentContext.value) return
    await assessmentStore.loadAssessment({
      classId,
      subjectId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
    })
  }

  const loadAssessmentById = async (assessmentId) => {
    await assessmentStore.loadAssessmentById({ assessmentId })
  }

  const createAssessment = async ({
    classId,
    className,
    subjectId,
    subjectName,
    date,
    components,
  }) => {
    if (!currentContext.value) return
    await assessmentStore.createAssessment({
      classId,
      className,
      subjectId,
      subjectName,
      date,
      components,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
    })
  }

  const updateScore = async ({ componentId, score }) => {
    if (!assessmentStore.assessment || assessmentStore.isFinalized) return
    await assessmentStore.updateScore({
      assessmentId: assessmentStore.assessment.id,
      componentId,
      score,
    })
  }

  const finalizeAssessment = async () => {
    if (!assessmentStore.assessment) return
    await assessmentStore.finalizeAssessment({
      assessmentId: assessmentStore.assessment.id,
    })
  }

  const previewResult = computed(() => assessmentStore.getPreviewResult())

  const navigateToSession = (assessmentId) => {
    router.push({
      name: 'assessment-session',
      params: { assessmentId },
    })
  }

  return {
    assessment,
    result,
    isLoading,
    isSaving,
    error,
    isDraft,
    isFinalized,
    previewResult,
    loadAssessment,
    loadAssessmentById,
    createAssessment,
    updateScore,
    finalizeAssessment,
    navigateToSession,
  }
}
