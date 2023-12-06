# git-releaser

![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/git-releaser/git-releaser)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/git-releaser/git-releaser/release.yaml)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/git-releaser/git-releaser)

`git-releaser` is a tool for creating Git releases based on [Semantic Versioning](https://semver.org/) heavily inspired by [release-please](https://github.com/release-please/release-please).

It is designed to be used in CI/CD pipelines to automate the creation of Git releases using PRs.

## Providers
Currently, `git-releaser` works best with GitLab. GitHub support is planned and more providers can be added in the future.

![](https://img.shields.io/badge/gitlab-ready-blue)
![](https://img.shields.io/badge/github-work_in_progress-orange)

# CLI Installation


### Linux/Mac via brew

```
brew tap git-releaser/git-releaser
brew install git-releaser
```

<details>
  <summary>Failing Installation on WSL or Linux (missing gcc)</summary>
  When installing Homebrew on WSL or Linux, you may encounter the following error:

  ```
  ==> Installing git-releaser from git-releaser/git-releaser Error: The following formula cannot be installed from a bottle and must be
  built from the source. git-releaser Install Clang or run brew install gcc.
  ```

If you install gcc as suggested, the problem will persist. Therefore, you need to install the build-essential package.
  ```
     sudo apt-get update
     sudo apt-get install build-essential
  ```
</details>

## Quick Start

### GitLab
* Create a new [GitLab Personal Access Token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) with `api` scope.
* Create a new `.git-releaser-manifest.yaml` in your repository with the following content:
```json
{"version": "0.1.0"}
```

* Create a new `.gitlab-ci.yml` file in your repository with the following content:
```yaml
stages:
  - release

variables:
  DOCKER_IMAGE: "ghcr.io/git-releaser/git-releaser:dev-202312060656"
  GIT_RELEASER_PROVIDER: "gitlab"
  GIT_RELEASER_USER_ID: "<your user-id>"
  GIT_RELEASER_TOKEN: $PAT
  GIT_RELEASER_PROJECT_URL: "<your-project-url>"
  GIT_RELEASER_PROJECT_ID: $CI_PROJECT_ID

run-release:
  stage: release
  script:
    - git clone $GIT_RELEASER_PROJECT_URL
    - cd git-releaser-demo
    - /git-releaser update
  image:
    name: $DOCKER_IMAGE
    entrypoint: [""]
  rules:
    - if: $CI_COMMIT_BRANCH == "main"
      when: always
```

* If you run the SaaS version of GitLab, the API URL is automatically detected. If you run a self-hosted version of GitLab, you need to set the `GIT_RELEASER_API_URL` variable to the URL of your GitLab instance.

* Afterwards, commit and push the changes to your repository. This will trigger a new pipeline which will create a new release based on the latest PRs merged into the `main` branch.

* When a new change is merged into the `main` branch, a new release merge will be created automatically.

* When this merge request is merged, a new release will be created automatically.

### Updating the version in config files
`git-releaser` can also update the version in config files. To do so, you need to specify the extra files in a `.git-releaser-config.yaml` file:

```yaml
extra_files:
- path: test2.txt
```

If the test2.txt file contains the following content:
```
my_version: 0.0.1 # x-git-releaser-version
other_version: 0.0.2
```

git-releaser will update the version specified n my_version during the release.


## Contributing
Please read our [contributing guide](./CONTRIBUTING.md).

## Community
<a href="https://github.com/git-releaser/git-releaser/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=git-releaser/git-releaser" />
</a>
