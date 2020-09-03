// gobadges a package for creating go badges
package gobadges

const (
	// user|organization / repository

	// Links to online Go report card
	goReportCardReport = "https://goreportcard.com/report/github.com/%s/%s"

	// Links to bade for embedding in HTML
	goReportCardBadge = "https://goreportcard.com/badge/github.com/%s/%s"

	// Markdown link for embedding
	goReportCardMarkdown = "[![Go Report Card](https://goreportcard.com/badge/github.com/%s/%s)](https://goreportcard.com/report/github.com/%s/%s)"
)

// What type of response is requested
const (
	GoReport = iota
	GoHTMLink
	GoMarkDown
)
