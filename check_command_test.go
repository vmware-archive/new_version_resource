package new_version_resource_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"
	resource "github.com/pivotal-cf-experimental/new_version_resource"
)

var _ = Describe("Check Command", func() {
	var (
		command *resource.CheckCommand
		request resource.CheckRequest
		source  resource.Source
		regex   string
	)

	Context("when the source is an http page", func() {
		var (
			httpSource resource.HTTPSource
		)

		BeforeEach(func() {
			command = resource.NewCheckCommand()
			regex = `[\d\.]+[0-9A-Za-z\-]*`

			httpSource = resource.HTTPSource{
				URL:     "https://api.mybiz.com/articles.html",
				CSSPath: "td",
			}

			source = resource.Source{
				Type:  "http",
				HTTP:  httpSource,
				Regex: regex,
			}

			httpmock.Activate()
		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		Context("when this is the first time that the resource has been run", func() {
			Context("when there are no releases", func() {
				BeforeEach(func() {
					httpmock.RegisterResponder("GET", "https://api.mybiz.com/articles.html",
						httpmock.NewStringResponder(200, "<html><body></body></html>"))
				})

				It("returns no versions", func() {
					request = resource.CheckRequest{
						Source:  source,
						Version: resource.Version{},
					}

					versions, err := command.Run(request)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(versions).Should(BeEmpty())
				})
			})

			Context("when there are releases", func() {
				Context("and the releases are ordered randomly and we want to use semver", func() {
					BeforeEach(func() {
						httpmock.RegisterResponder("GET", "https://api.mybiz.com/articles.html",
							httpmock.NewStringResponder(200, `<html><body>
					<table>
					<tr><td>1.0</td></tr>
					<tr><td>1.2</td></tr>
					<tr><td>1.1</td></tr>
					</table>
					</body></html>`))
					})

					It("outputs the most recent version only", func() {
						request = resource.CheckRequest{
							Source:  source,
							Version: resource.Version{},
						}

						response, err := command.Run(request)
						Ω(err).ShouldNot(HaveOccurred())

						Ω(response).Should(HaveLen(1))
						Ω(response[0]).Should(Equal(resource.Version{
							Version: "1.2",
						}))
					})
				})
			})
		})

		Context("when there are prior versions", func() {
			Context("and the releases are ordered randomly and we want to use semver", func() {

				BeforeEach(func() {
					httpmock.RegisterResponder("GET", "https://api.mybiz.com/articles.html",
						httpmock.NewStringResponder(200, `<html><body>
					<table>
					<tr><td>1.0</td></tr>
					<tr><td>1.1</td></tr>
					<tr><td>1.3</td></tr>
					<tr><td>1.2</td></tr>
					</table>
					</body></html>`))
				})

				It("returns an empty list if the latest version has been checked", func() {
					request = resource.CheckRequest{
						Source: source,
						Version: resource.Version{
							Version: "1.3",
						},
					}

					response, err := command.Run(request)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(response).Should(BeEmpty())
				})

				It("returns all of the versions that are newer", func() {
					request = resource.CheckRequest{
						Source: source,
						Version: resource.Version{
							Version: "1.1",
						},
					}

					response, err := command.Run(request)
					Ω(err).ShouldNot(HaveOccurred())

					Ω(response).Should(Equal([]resource.Version{
						{Version: "1.3"},
						{Version: "1.2"},
						{Version: "1.1"},
					}))
				})

				It("returns the latest version if the current version is not found", func() {
					request = resource.CheckRequest{
						Source: source,
						Version: resource.Version{
							Version: "1.7",
						},
					}

					command := resource.NewCheckCommand()

					response, err := command.Run(request)
					Ω(err).ShouldNot(HaveOccurred())

					Ω(response).Should(Equal([]resource.Version{
						{Version: "1.3"},
					}))
				})
			})
		})

		Context("and the releases have additional characters we want to strip", func() {
			BeforeEach(func() {
				httpmock.RegisterResponder("GET", "https://api.mybiz.com/articles.html",
					httpmock.NewStringResponder(200, `<html><body>
					<table>
					<tr><td>1.0/</td></tr>
					<tr><td>1.1/</td></tr>
					<tr><td>1.2/</td></tr>
					<tr><td>1.3/</td></tr>
					</table>
					</body></html>`))
			})

			It("returns an empty list if the latest version has been checked", func() {
				request = resource.CheckRequest{
					Source: source,
					Version: resource.Version{
						Version: "1.3",
					},
				}

				response, err := command.Run(request)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(response).Should(BeEmpty())
			})

			It("returns all of the versions that are newer", func() {
				request = resource.CheckRequest{
					Source: source,
					Version: resource.Version{
						Version: "1.1",
					},
				}

				response, err := command.Run(request)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(response).Should(Equal([]resource.Version{
					{Version: "1.3"},
					{Version: "1.2"},
					{Version: "1.1"},
				}))
			})

			It("returns the latest version if the current version is not found", func() {
				request = resource.CheckRequest{
					Source: source,
					Version: resource.Version{
						Version: "1.5",
					},
				}

				response, err := command.Run(request)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(response).Should(Equal([]resource.Version{
					{Version: "1.3"},
				}))
			})
		})
	})

	Context("when the versions are coming from Github", func() {
		var (
			gitSource resource.GitSource
		)

		BeforeEach(func() {
			command = resource.NewCheckCommand()
			regex = `v[\d\.]+[0-9A-Za-z\-]*`

			gitSource = resource.GitSource{
				Organization: "bundler",
				Repo:         "bundler",
			}

			source = resource.Source{
				Type:  "git",
				Git:   gitSource,
				Regex: regex,
			}
		})

		It("returns an empty list if the latest version has been checked", func() {
			request = resource.CheckRequest{
				Source: source,
				Version: resource.Version{
					Version: "v1.13.7",
				},
			}

			response, err := command.Run(request)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(response).Should(BeEmpty())
		})

		It("returns all of the versions that are newer", func() {
		request = resource.CheckRequest{
				Source: source,
				Version: resource.Version{
					Version: "v1.13.4",
				},
			}

			response, err := command.Run(request)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(response).Should(Equal([]resource.Version{
				{Version: "v1.13.7"},
				{Version: "v1.13.6"},
				{Version: "v1.13.5"},
				{Version: "v1.13.4"},
			}))
		})

		It("returns the latest version if the current version is not found", func() {
			request = resource.CheckRequest{
				Source: source,
				Version: resource.Version{
					Version: "v1.45.45",
				},
			}

			response, err := command.Run(request)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(response).Should(Equal([]resource.Version{
				{Version: "v1.13.7"},
			}))
		})
	})
})
