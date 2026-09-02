package value_object

type AssetSpecification struct {
	Brand      string
	Model      string
	Year       int
	Color      string
	Size       string            // misal: "42 inch", "A4"
	Additional map[string]string // fleksibel
}
