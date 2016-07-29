# new_version_resource
Version tracking based on JQuery CSS Paths

For when you are trying to track versioned releases on websites.

## Building Binaries
```bash
./scripts/build
```

## Building docker image
```bash
sudo docker build -t your-dockerhub-username/new_version_resource .
sudo docker push your-dockerhub-username/new_version_resource
```

For deploying as the Buildpacks team:

```bash
./scripts/publish
```

## Usage

You will need to setup the new resource type, and specify a `url`, and `csspath` to the versions strings on that page.
An example is below:

```yaml
---
resource_types:
- name: new_version_resource
  type: docker-image
  source:
    repository: cfbuildpacks/new_version_resource

resources:
  - name: pecl-igbinary
    type: new_version_resource
    source:
      url: https://pecl.php.net/package/igbinary
      csspath: table[cellpadding="2"][cellspacing="1"] tr:has(td:contains("stable")) th
      use_semver: true
```
