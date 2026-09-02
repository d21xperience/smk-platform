/* eslint-disable no-unused-vars */
// Interface / Base Class
export class BaseExportAdapter {
  export(data, filename) {
    throw new Error('Method export must be implemented by subclass')
  }
}
