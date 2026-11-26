package collectors

import (
	"os"
	"strings"
)

type InstalledApplication struct {
	Name string
}

func GetInstalledApplications() ([]InstalledApplication, error) {
	var applications []InstalledApplication
	var newApplication InstalledApplication
	result, err := os.ReadDir("/Applications")
	if err != nil {
		return nil, err
	}
	for _, application := range result {
		if application.IsDir() && strings.HasSuffix(application.Name(), ".app") {
			newApplication.Name = application.Name()
			applications = append(applications, newApplication)
		}
	}
	return applications, nil
}
