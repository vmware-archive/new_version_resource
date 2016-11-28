package new_version_resource

import (
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-github/github"
	version "github.com/hashicorp/go-version"
)

type CheckCommand struct {
}

func NewCheckCommand() *CheckCommand {
	return &CheckCommand{}
}

type VersionSorter struct {
	versions []Version
}

func (v *VersionSorter) Len() int {
	return len(v.versions)
}

func (v *VersionSorter) Swap(i, j int) {
	v.versions[i], v.versions[j] = v.versions[j], v.versions[i]
}

func (v *VersionSorter) Less(i, j int) bool {
	semI, _ := version.NewVersion(v.versions[i].Version)
	semJ, _ := version.NewVersion(v.versions[j].Version)

	return semI.LessThan(semJ)
}

func (c *CheckCommand) getVersionsFromHttp(source HTTPSource, regex string) ([]Version, error) {
	res, err := http.Get(source.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(regex)
	versions := make([]Version, 0)

	doc.Find(source.CSSPath).Each(func(i int, s *goquery.Selection) {
		versions = append(versions, Version{
			Version: re.FindAllString(strings.TrimSpace(s.Text()), -1)[0],
		})
	})
	return versions, nil
}

func (c *CheckCommand) getVersionsFromGithub(source GitSource, regex string) ([]Version, error) {
	client := github.NewClient(nil)

	re := regexp.MustCompile(regex)

	tags, _, err := client.Repositories.ListTags(source.Organization, source.Repo, nil)

	if err != nil {
		return nil, err
	}

	versions := make([]Version, 0)

	for _, tag := range tags {
		versions = append(versions, Version{
			Version: re.FindAllString(*tag.Name, -1)[0],
		})
	}

	return versions, nil
}

func (c *CheckCommand) Run(request CheckRequest) ([]Version, error) {
	var versions []Version
	var err error

	if request.Source.Type == "http" {
		versions, err = c.getVersionsFromHttp(request.Source.HTTP, request.Source.Regex)
	} else if request.Source.Type == "git" {
		versions, err = c.getVersionsFromGithub(request.Source.Git, request.Source.Regex)
	} else {

	}

	if err != nil {
		return nil, err
	}

	sorter := new(VersionSorter)
	sorter.versions = versions

	sort.Sort(sort.Reverse(sorter))

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
