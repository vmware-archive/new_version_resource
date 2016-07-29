package new_version_resource

type Source struct {
	URL       string `json:"url"`
	CSSPath   string `json:"csspath"`
	UseSemver bool   `json:use_semver`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InRequest struct {
	Source          Source  `json:"source"`
	Version         Version `json:"version"`
	FirstOccurrence bool    `json:"first_occurrence"`
}
type Version struct {
	Version string `json:"version,omitempty"`
}

type InResponse struct {
	Version  Version        `json:"version"`
	Metadata []MetadataPair `json:"metadata"`
}

type MetadataPair struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	URL      string `json:"url"`
	Markdown bool   `json:"markdown"`
}
