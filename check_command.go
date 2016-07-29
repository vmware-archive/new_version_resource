package new_version_resource

import (
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	version "github.com/hashicorp/go-version"
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

func getSemverVersions(resourceVersions []Version) []*version.Version {
	versions := make([]*version.Version, len(resourceVersions))
	re := regexp.MustCompile(`[\d\.]+.*`)
	for idx, resource_version := range resourceVersions {
		cleaned_resource_version := re.FindAllString(resource_version.Version, -1)[0]
		v, _ := version.NewVersion(cleaned_resource_version)
		versions[idx] = v
	}
	return versions
}

func versionFromSemverVersions(semverVersions []*version.Version) []Version {
	versions := make([]Version, len(semverVersions))
	for idx, semverVersion := range semverVersions {
		versions[idx] = Version{Version: semverVersion.String()}
	}
	return versions
}

func (c *CheckCommand) Run(request CheckRequest) ([]Version, error) {
	versions, err := c.getVersions(request.Source)
	if request.Source.UseSemver {
		//get version semvers
		semverVersions := getSemverVersions(versions)
		//sort them
		sort.Sort(sort.Reverse(version.Collection(semverVersions)))
		//cast back to regular version types
		versions = versionFromSemverVersions(semverVersions)
	}

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
