package processor

// Represents a license inside the JSON which allows us to hopefully match against
type License struct {
	LicenseText             string   `json:"licenseText"`
	StandardLicenseTemplate string   `json:"standardLicenseTemplate"`
	Name                    string   `json:"name"`
	LicenseId               string   `json:"licenseId"`
	Keywords                []string `json:"keywords"`
	ScorePercentage         float64  `json:"scorePercentage"` // this is used so we don't have a new struct
}