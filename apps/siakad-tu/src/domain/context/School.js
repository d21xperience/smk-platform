// apps/siakad-tu/src/domain/context/School.js

export class School {
  constructor(id, name, address, headmasterName) {
    this.id = id;
    this.name = name;
    this.address = address;
    this.headmasterName = headmasterName;
    this._events = [];
  }

  updateProfile(newName, newAddress) {
    if (!newName || newName.trim() === '') {
      throw new Error('Nama sekolah tidak boleh kosong');
    }
    
    this.name = newName;
    this.address = newAddress;
    
    // Record event (akan diproses oleh Service/Projection)
    this._events.push({
      eventName: 'SchoolProfileUpdated',
      aggregateId: this.id,
      payload: { name: newName, address: newAddress }
    });
  }

  getUncommittedEvents() {
    return this._events;
  }

  clearEvents() {
    this._events = [];
  }
}