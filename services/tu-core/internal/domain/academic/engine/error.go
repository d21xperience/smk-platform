package engine

type errString string

func (e errString) Error() string { return string(e) }

var (
	errAcademicYearNotFound = errString("tahun ajaran tidak ditemukan")
	errClassNotFound        = errString("kelas tidak ditemukan")
)
