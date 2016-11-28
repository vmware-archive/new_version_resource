package new_version_resource

type Source struct {
	Type  string     `json:"type"`
	Regex string     `json:"regex"`
	Git   GitSource  `json:"git",omitempty`
	HTTP  HTTPSource `json:"http",omitempty`
}

type HTTPSource struct {
	URL     string `json:"url"`
	CSSPath string `json:"csspath"`
}

type GitSource struct {
	Organization string `json:"organization"`
	Repo         string `json:"repo"`
	AccessToken  string `json:"access_token",omitempty`
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
