package new_version_resource

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type CheckCommand struct {
}

func NewCheckCommand() *CheckCommand {
	return &CheckCommand{}
}

func (c *CheckCommand) getVersions(source Source) ([]Version, error) {
	res, err := http.Get(source.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	versions := make([]Version, 0)
	doc.Find(source.CSSPath).Each(func(i int, s *goquery.Selection) {
		versions = append(versions, Version{
			Version: strings.TrimSpace(s.Text()),
		})
	})
	return versions, nil
}

func (c *CheckCommand) Run(request CheckRequest) ([]Version, error) {
	versions, err := c.getVersions(request.Source)
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return versions, nil
	}

	if request.Version.Version == "" {
		return versions[:1], nil
	}

	if request.Version.Version == versions[0].Version {
		return nil, nil
	}

	for i, version := range versions {
		if request.Version.Version == version.Version {
			return versions[:i+1], nil
		}
	}

	return versions[:1], nil
}
