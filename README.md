# git-releaser

![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/git-releaser/git-releaser)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/git-releaser/git-releaser/release.yaml)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/git-releaser/git-releaser)

`git-releaser` is a tool for creating Git releases based on [Semantic Versioning](https://semver.org/) heavily inspired by [release-please](https://github.com/release-please/release-please).

It is designed to be used in CI/CD pipelines to automate the creation of Git releases using PRs.

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

### GitHub
TBD

### GitLab
TBD

## Contributing
Please read our [contributing guide](./CONTRIBUTING.md).

## Community
<a href="https://github.com/git-releaser/git-releaser/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=git-releaser/git-releaser" />
</a>
