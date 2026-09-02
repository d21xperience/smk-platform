// PURE JS - Value Object
export class Material {
  constructor(props) {
    this.title = props.title || ''
    this.content = props.content || '' // bisa berupa teks panjang atau markdown
    this.attachments = props.attachments || [] // array of { name, url, type }
  }

  isEmpty() {
    return !this.content || this.content.trim() === ''
  }

  toJSON() {
    return {
      title: this.title,
      content: this.content,
      attachments: this.attachments,
    }
  }

  static fromJSON(json) {
    return new Material(json)
  }
}
