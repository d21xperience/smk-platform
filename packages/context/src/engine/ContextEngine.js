export const ContextEngine = {
  validateContext(context) {
    if (!context.userId) throw new Error('userId is required');
    if (!context.schoolId) throw new Error('schoolId is required');
    return true;
  }
};
