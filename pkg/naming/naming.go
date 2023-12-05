package naming

import "fmt"

func GeneratePrTitle(version string) string {
	title := fmt.Sprintf("Release %s", version)
	return title
}

func CreatePrDescription(version string, changelog string) string {
	return fmt.Sprintf("This is a description for the new pull request for version %s.\n\n## Changelog\n\n%s", version, changelog)
}
