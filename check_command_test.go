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
	)

	BeforeEach(func() {
		command = resource.NewCheckCommand()
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
				request := resource.CheckRequest{
					Source: resource.Source{
						URL:     "https://api.mybiz.com/articles.html",
						CSSPath: "td",
					},
					Version: resource.Version{},
				}
				versions, err := command.Run(request)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(versions).Should(BeEmpty())
			})
		})

		Context("when there are releases", func() {
			Context("and the releases are ordered from newest to oldest", func() {
				BeforeEach(func() {
					httpmock.RegisterResponder("GET", "https://api.mybiz.com/articles.html",
						httpmock.NewStringResponder(200, `<html><body>
					<table>
					<tr><td>1.2</td></tr>
					<tr><td>1.1</td></tr>
					<tr><td>1.0</td></tr>
					</table>
					</body></html>`))
				})

				It("outputs the most recent version only", func() {
					command := resource.NewCheckCommand()
					request := resource.CheckRequest{
						Source: resource.Source{
							URL:     "https://api.mybiz.com/articles.html",
							CSSPath: "td",
						},
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

			Context("and the releases are ordered randomly and we want to use semver", func() {
				BeforeEach(func() {
					httpmock.RegisterResponder("GET", "https://api.mybiz.com/articles.html",
						httpmock.NewStringResponder(200, `<html><body>
					<table>
					<tr><td>1.0</td></tr>
					<tr><td>1.1</td></tr>
					<tr><td>1.2</td></tr>
					</table>
					</body></html>`))
				})

				It("outputs the most recent version only", func() {
					command := resource.NewCheckCommand()
					request := resource.CheckRequest{
						Source: resource.Source{
							URL:       "https://api.mybiz.com/articles.html",
							CSSPath:   "td",
							UseSemver: true,
						},
						Version: resource.Version{},
					}
					response, err := command.Run(request)
					Ω(err).ShouldNot(HaveOccurred())

					Ω(response).Should(HaveLen(1))
					Ω(response[0]).Should(Equal(resource.Version{
						Version: "1.2.0",
					}))
				})
			})
		})
	})

	Context("when there are prior versions", func() {
		Context("and the releases are ordered from newest to oldest", func() {

			BeforeEach(func() {
				httpmock.RegisterResponder("GET", "https://api.mybiz.com/articles.html",
					httpmock.NewStringResponder(200, `<html><body>
					<table>
					<tr><td>1.3</td></tr>
					<tr><td>1.2</td></tr>
					<tr><td>1.1</td></tr>
					<tr><td>1.0</td></tr>
					</table>
					</body></html>`))
			})

			It("returns an empty list if the lastet version has been checked", func() {
				command := resource.NewCheckCommand()

				response, err := command.Run(resource.CheckRequest{
					Source: resource.Source{
						URL:     "https://api.mybiz.com/articles.html",
						CSSPath: "td",
					},
					Version: resource.Version{
						Version: "1.3",
					},
				})
				Ω(err).ShouldNot(HaveOccurred())
				Ω(response).Should(BeEmpty())
			})

			It("returns all of the versions that are newer", func() {
				command := resource.NewCheckCommand()

				response, err := command.Run(resource.CheckRequest{
					Source: resource.Source{
						URL:     "https://api.mybiz.com/articles.html",
						CSSPath: "td",
					},
					Version: resource.Version{
						Version: "1.1",
					},
				})
				Ω(err).ShouldNot(HaveOccurred())

				Ω(response).Should(Equal([]resource.Version{
					{Version: "1.3"},
					{Version: "1.2"},
					{Version: "1.1"},
				}))
			})

			It("returns the latest version if the current version is not found", func() {
				command := resource.NewCheckCommand()

				response, err := command.Run(resource.CheckRequest{
					Source: resource.Source{
						URL:     "https://api.mybiz.com/articles.html",
						CSSPath: "td",
					},
					Version: resource.Version{
						Version: "1.4",
					},
				})
				Ω(err).ShouldNot(HaveOccurred())

				Ω(response).Should(Equal([]resource.Version{
					{Version: "1.3"},
				}))
			})
		})

		Context("and the releases are ordered randomly and we want to use semver", func() {

			BeforeEach(func() {
				httpmock.RegisterResponder("GET", "https://api.mybiz.com/articles.html",
					httpmock.NewStringResponder(200, `<html><body>
					<table>
					<tr><td>1.0</td></tr>
					<tr><td>1.1</td></tr>
					<tr><td>1.2</td></tr>
					<tr><td>1.3</td></tr>
					</table>
					</body></html>`))
			})

			It("returns an empty list if the latest version has been checked", func() {
				command := resource.NewCheckCommand()

				response, err := command.Run(resource.CheckRequest{
					Source: resource.Source{
						URL:       "https://api.mybiz.com/articles.html",
						CSSPath:   "td",
						UseSemver: true,
					},
					Version: resource.Version{
						Version: "1.3.0",
					},
				})
				Ω(err).ShouldNot(HaveOccurred())
				Ω(response).Should(BeEmpty())
			})

			It("returns all of the versions that are newer", func() {
				command := resource.NewCheckCommand()

				response, err := command.Run(resource.CheckRequest{
					Source: resource.Source{
						URL:       "https://api.mybiz.com/articles.html",
						CSSPath:   "td",
						UseSemver: true,
					},
					Version: resource.Version{
						Version: "1.1.0",
					},
				})
				Ω(err).ShouldNot(HaveOccurred())

				Ω(response).Should(Equal([]resource.Version{
					{Version: "1.3.0"},
					{Version: "1.2.0"},
					{Version: "1.1.0"},
				}))
			})

			It("returns the latest version if the current version is not found", func() {
				command := resource.NewCheckCommand()

				response, err := command.Run(resource.CheckRequest{
					Source: resource.Source{
						URL:       "https://api.mybiz.com/articles.html",
						CSSPath:   "td",
						UseSemver: true,
					},
					Version: resource.Version{
						Version: "1.4.0",
					},
				})
				Ω(err).ShouldNot(HaveOccurred())

				Ω(response).Should(Equal([]resource.Version{
					{Version: "1.3.0"},
				}))
			})
		})
	})
})
