// Package fossa CLI wrapper
//
//   Support is limited to:
//		Badge generation
package fossa

import (
	"fmt"
)

func GetBadge(outputType, badgeSize int, projectName string) string {
	var result string

	switch outputType {
	case MarkDown:
		var MarkDownURL = "[![FOSSA Status](" + APIURL + ")](" + RefURL + ")"
		result = fmt.Sprintf(MarkDownURL,
			projectName,
			MetricName[badgeSize],
			projectName,
			MetricName[badgeSize])
	case HTML:
		var HTMLURL = "<a href=\"" + RefURL + "\" alt=\"FOSSA Status\"><img src=\"" + APIURL + "\"/></a>"
		result = fmt.Sprintf(HTMLURL,
			projectName,
			MetricName[badgeSize],
			projectName,
			MetricName[badgeSize])
	case Link:
		var LinkURL = APIURL
		result = fmt.Sprintf(LinkURL,
			projectName,
			MetricName[badgeSize])
	}
	return result
}
