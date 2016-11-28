package new_version_resource

import (
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-github/github"
	semver "github.com/hashicorp/go-version"
	"golang.org/x/oauth2"
)

type CheckCommand struct {
}

func NewCheckCommand() *CheckCommand {
	return &CheckCommand{}
}

type SortVersion struct {
	Version string
	Semver  *semver.Version
}

type VersionSorter struct {
	sortVersions []SortVersion
}

func (v *VersionSorter) Len() int {
	return len(v.sortVersions)
}

func (v *VersionSorter) Swap(i, j int) {
	v.sortVersions[i], v.sortVersions[j] = v.sortVersions[j], v.sortVersions[i]
}

func (v *VersionSorter) Less(i, j int) bool {
	return v.sortVersions[i].Semver.LessThan(v.sortVersions[j].Semver)
}

func (c *CheckCommand) getVersionsFromHttp(source HTTPSource, regex string) ([]SortVersion, error) {
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
	versions := make([]SortVersion, 0)

	doc.Find(source.CSSPath).Each(func(i int, s *goquery.Selection) {
		matches := re.FindAllString(strings.TrimSpace(s.Text()), -1)

		if matches != nil {
			v, err := semver.NewVersion(matches[0])

			if err != nil {
				return
			}

			versions = append(versions, SortVersion{
				Version: matches[0],
				Semver:  v,
			})
		}
	})
	return versions, nil
}

func (c *CheckCommand) getVersionsFromGithub(source GitSource, regex string) ([]SortVersion, error) {
	var client *github.Client

	if source.AccessToken == "" {
		client = github.NewClient(nil)
	} else {
		tokenSource := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: source.AccessToken},
		)

		tokenClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

		client = github.NewClient(tokenClient)
	}

	re := regexp.MustCompile(regex)

	versions := make([]SortVersion, 0)

	options := &github.ListOptions{PerPage: 100}

	tags, _, err := client.Repositories.ListTags(source.Organization, source.Repo, options)

	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		matches := re.FindAllString(*tag.Name, -1)

		if matches != nil {
			v, err := semver.NewVersion(matches[0])

			if err != nil {
				continue
			}

			versions = append(versions, SortVersion{
				Version: matches[0],
				Semver:  v,
			})
		}
	}

	return versions, nil
}

func (c *CheckCommand) Run(request CheckRequest) ([]Version, error) {
	var sortVersions []SortVersion
	var err error

	if request.Source.Type == "http" {
		sortVersions, err = c.getVersionsFromHttp(request.Source.HTTP, request.Source.Regex)
	} else if request.Source.Type == "git" {
		sortVersions, err = c.getVersionsFromGithub(request.Source.Git, request.Source.Regex)
	}

	if err != nil {
		return nil, err
	}

	sorter := new(VersionSorter)
	sorter.sortVersions = sortVersions

	sort.Sort(sort.Reverse(sorter))

	outputVersions := make([]Version, len(sortVersions))

	for i, _ := range sortVersions {
		outputVersions[i] = Version{Version: sortVersions[i].Version}
	}

	if len(outputVersions) == 0 {
		return outputVersions, nil
	}

	if request.Version.Version == "" {
		return outputVersions[:1], nil
	}

	if request.Version.Version == outputVersions[0].Version {
		return nil, nil
	}

	for i, version := range outputVersions {
		if request.Version.Version == version.Version {
			return outputVersions[:i+1], nil
		}
	}

	return outputVersions[:1], nil
}
