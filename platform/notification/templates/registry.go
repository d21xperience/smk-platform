package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Template merepresentasikan template notifikasi dengan subject dan body terpisah.
type Template struct {
	ID      string
	Subject string
	Body    string
}

// Registry mengelola template dari multiple sources (memory, file, DB).
type Registry struct {
	mu        sync.RWMutex
	templates map[string]Template
}

// NewRegistry membuat registry baru.
func NewRegistry() *Registry {
	return &Registry{
		templates: make(map[string]Template),
	}
}

// Register menambahkan template ke registry.
func (r *Registry) Register(tmpl Template) error {
	if tmpl.ID == "" {
		return fmt.Errorf("template ID tidak boleh kosong")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.templates[tmpl.ID] = tmpl
	return nil
}

// Get mengambil template dari registry.
func (r *Registry) Get(id string) (Template, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tmpl, exists := r.templates[id]
	return tmpl, exists
}

// List mengembalikan semua template yang terdaftar.
func (r *Registry) List() []Template {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Template, 0, len(r.templates))
	for _, tmpl := range r.templates {
		result = append(result, tmpl)
	}
	return result
}

// Count mengembalikan jumlah template yang terdaftar.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.templates)
}

// LoadFromDir memuat semua template dari direktori.
// Format nama file: {templateID}.tmpl
// File harus berisi 2 baris: Subject di baris pertama, Body di baris kedua (dipisahkan oleh "---")
func (r *Registry) LoadFromDir(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("gagal baca direktori %s: %w", dirPath, err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".tmpl" {
			continue
		}

		id := filepath.Base(file.Name())
		id = id[:len(id)-len(".tmpl")]

		filePath := filepath.Join(dirPath, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("gagal load template dari file %s: %w", filePath, err)
		}

		// Parse: Subject di baris pertama, Body setelah separator "---"
		content := string(data)
		subject := ""
		body := ""

		parts := splitTemplate(content)
		if len(parts) == 2 {
			subject = parts[0]
			body = parts[1]
		} else {
			// Jika tidak ada separator, anggap semuanya adalah body
			body = content
		}

		tmpl := Template{
			ID:      id,
			Subject: subject,
			Body:    body,
		}

		if err := r.Register(tmpl); err != nil {
			return err
		}
	}

	return nil
}

// splitTemplate memisahkan subject dan body berdasarkan separator "---"
func splitTemplate(content string) []string {
	separator := "\n---\n"
	for i := 0; i < len(content)-len(separator)+1; i++ {
		if content[i:i+len(separator)] == separator {
			return []string{content[:i], content[i+len(separator):]}
		}
	}
	return []string{content}
}

// DefaultTemplates mengembalikan template default untuk notifikasi sekolah.
func DefaultTemplates() []Template {
	return []Template{
		{
			ID:      "invoice_created_wa",
			Subject: "",
			Body:    "Yth. {{.parentName}},\n\nTagihan {{.componentType}} untuk {{.studentName}} sebesar Rp {{.amount}} jatuh tempo pada {{.dueDate}}.\n\nMohon segera melakukan pembayaran.\n\nTerima kasih.",
		},
		{
			ID:      "invoice_created_email",
			Subject: "Tagihan {{.componentType}} - {{.studentName}}",
			Body:    "<h1>Tagihan Sekolah</h1><p>Yth. {{.parentName}},</p><p>Berikut adalah tagihan untuk {{.studentName}}:</p><ul><li>Komponen: {{.componentType}}</li><li>Jumlah: Rp {{.amount}}</li><li>Jatuh Tempo: {{.dueDate}}</li></ul><p>Mohon segera melakukan pembayaran.</p><p>Terima kasih.</p>",
		},
		{
			ID:      "payment_received_wa",
			Subject: "",
			Body:    "Yth. {{.parentName}},\n\nPembayaran sebesar Rp {{.amount}} untuk {{.studentName}} telah kami terima pada {{.paymentDate}} via {{.paymentMethod}}.\n\nTerima kasih.",
		},
		{
			ID:      "payment_received_email",
			Subject: "Konfirmasi Pembayaran - {{.studentName}}",
			Body:    "<h1>Konfirmasi Pembayaran</h1><p>Yth. {{.parentName}},</p><p>Pembayaran Anda telah kami terima:</p><ul><li>Jumlah: Rp {{.amount}}</li><li>Tanggal: {{.paymentDate}}</li><li>Metode: {{.paymentMethod}}</li></ul><p>Terima kasih.</p>",
		},
		{
			ID:      "document_issued_wa",
			Subject: "",
			Body:    "Yth. {{.parentName}},\n\nDokumen {{.documentType}} (No: {{.documentNumber}}) untuk {{.studentName}} telah diterbitkan pada {{.issueDate}}.\n\nSilakan ambil di Tata Usaha.",
		},
		{
			ID:      "document_issued_email",
			Subject: "Dokumen {{.documentType}} Telah Diterbitkan",
			Body:    "<h1>Dokumen Telah Diterbitkan</h1><p>Yth. {{.parentName}},</p><p>Dokumen berikut telah diterbitkan untuk {{.studentName}}:</p><ul><li>Jenis: {{.documentType}}</li><li>Nomor: {{.documentNumber}}</li><li>Tanggal: {{.issueDate}}</li></ul><p>Silakan ambil di Tata Usaha.</p>",
		},
		{
			ID:      "attendance_absent_wa",
			Subject: "",
			Body:    "Yth. {{.parentName}},\n\n{{.studentName}} tidak hadir di {{.className}} pada {{.date}}.\n\nMohon konfirmasi jika ada keperluan mendesak.",
		},
		{
			ID:      "attendance_absent_email",
			Subject: "Ketidakhadiran {{.studentName}} - {{.date}}",
			Body:    "<h1>Pemberitahuan Ketidakhadiran</h1><p>Yth. {{.parentName}},</p><p>Kami informasikan bahwa {{.studentName}} tidak hadir di {{.className}} pada {{.date}}.</p><p>Mohon konfirmasi jika ada keperluan mendesak.</p>",
		},
		// Tambahkan ke dalam slice return di fungsi DefaultTemplates()

		{
			ID:      "finance_payment_received_wa",
			Subject: "",
			Body:    "Yth. {{.parentName}},\n\nPembayaran {{.componentType}} bulan {{.month}} {{.year}} untuk {{.studentName}} sebesar Rp {{.amount}} telah kami terima pada {{.paymentDate}}.\n\nTerima kasih atas kerjasamanya.\n- {{.schoolName}}",
		},
		{
			ID:      "finance_payment_received_email",
			Subject: "Konfirmasi Pembayaran {{.componentType}} - {{.studentName}}",
			Body:    "<h1>Konfirmasi Pembayaran</h1><p>Yth. {{.parentName}},</p><p>Pembayaran Anda telah kami terima:</p><ul><li>Siswa: {{.studentName}}</li><li>Komponen: {{.componentType}} (Bulan {{.month}} {{.year}})</li><li>Jumlah: Rp {{.amount}}</li><li>Tanggal: {{.paymentDate}}</li></ul><p>Terima kasih.</p>",
		},
		{
			ID:      "finance_payment_reminder_wa",
			Subject: "",
			Body:    "Yth. {{.parentName}},\n\nPengingat: Tagihan {{.componentType}} bulan {{.month}} {{.year}} untuk {{.studentName}} sebesar Rp {{.amount}} belum kami terima.\n\nJatuh tempo: {{.dueDate}}.\n\nMohon segera melakukan pembayaran untuk menghindari denda.\n- {{.schoolName}}",
		},
		{
			ID:      "finance_payment_reminder_email",
			Subject: "Pengingat Pembayaran {{.componentType}} - {{.studentName}}",
			Body:    "<h1>Pengingat Pembayaran</h1><p>Yth. {{.parentName}},</p><p>Kami ingin mengingatkan bahwa tagihan berikut belum dibayarkan:</p><ul><li>Siswa: {{.studentName}}</li><li>Komponen: {{.componentType}} (Bulan {{.month}} {{.year}})</li><li>Jumlah: Rp {{.amount}}</li><li>Jatuh Tempo: {{.dueDate}}</li></ul><p>Mohon segera melakukan pembayaran.</p>",
		},
		// Tambahkan ke dalam slice return di fungsi DefaultTemplates()

		{
			ID:      "finance_payment_received_wa",
			Subject: "",
			Body:    "Yth. {{.parentName}},\n\nPembayaran {{.componentType}} bulan {{.month}} {{.year}} untuk {{.studentName}} sebesar Rp {{.amount}} telah kami terima pada {{.paymentDate}}.\n\nTerima kasih atas kerjasamanya.\n- {{.schoolName}}",
		},
		{
			ID:      "finance_payment_received_email",
			Subject: "Konfirmasi Pembayaran {{.componentType}} - {{.studentName}}",
			Body:    "<h1>Konfirmasi Pembayaran</h1><p>Yth. {{.parentName}},</p><p>Pembayaran Anda telah kami terima:</p><ul><li>Siswa: {{.studentName}}</li><li>Komponen: {{.componentType}} (Bulan {{.month}} {{.year}})</li><li>Jumlah: Rp {{.amount}}</li><li>Tanggal: {{.paymentDate}}</li></ul><p>Terima kasih.</p>",
		},
		{
			ID:      "finance_payment_reminder_wa",
			Subject: "",
			Body:    "Yth. {{.parentName}},\n\nPengingat: Tagihan {{.componentType}} bulan {{.month}} {{.year}} untuk {{.studentName}} sebesar Rp {{.amount}} belum kami terima.\n\nJatuh tempo: {{.dueDate}}.\n\nMohon segera melakukan pembayaran untuk menghindari denda.\n- {{.schoolName}}",
		},
		{
			ID:      "finance_payment_reminder_email",
			Subject: "Pengingat Pembayaran {{.componentType}} - {{.studentName}}",
			Body:    "<h1>Pengingat Pembayaran</h1><p>Yth. {{.parentName}},</p><p>Kami ingin mengingatkan bahwa tagihan berikut belum dibayarkan:</p><ul><li>Siswa: {{.studentName}}</li><li>Komponen: {{.componentType}} (Bulan {{.month}} {{.year}})</li><li>Jumlah: Rp {{.amount}}</li><li>Jatuh Tempo: {{.dueDate}}</li></ul><p>Mohon segera melakukan pembayaran.</p>",
		},
	}
}
