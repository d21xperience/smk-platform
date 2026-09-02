export function createExpedition(data = {}) {
  return {
    id: data.id || null,
    schoolId: data.schoolId || null,
    outgoingLetterId: data.outgoingLetterId || null,
    courier: data.courier || '',
    sendDate: data.sendDate || null,
    receiptProof: data.receiptProof || '',
    createdAt: data.createdAt || null,
  }
}
