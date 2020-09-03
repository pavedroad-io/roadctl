package gobadges

import "fmt"

func GetGoBadge(badgeType int, orgOrUser, repo string) (badge string) {
	var result string

	switch badgeType {
	case GoReport:
		result = fmt.Sprintf(goReportCardReport,
			orgOrUser, repo)
	case GoHTMLink:
		result = fmt.Sprintf(goReportCardBadge,
			orgOrUser, repo)
	case GoMarkDown:
		result = fmt.Sprintf(goReportCardMarkdown,
			orgOrUser, repo,
			orgOrUser, repo)
	}

	return result
}
