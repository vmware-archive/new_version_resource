# new_version_resource

You will need to setup the new resource type, and specify a `url`, and `csspath` to the versions strings on that page.
An example is below

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
```
