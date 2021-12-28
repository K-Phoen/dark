# Contributing Guidelines

`dark` is [MIT licensed](LICENSE) and accepts contributions via
GitHub pull requests. This document outlines some of the conventions on
development workflow, commit message formatting, contact points, and other
resources to make it easier to get your contribution accepted.

## Support Channels

The official support channels, for both users and contributors, are:

- GitHub [issues](https://github.com/K-Phoen/dark/issues)

## How to Contribute

Pull Requests (PRs) are the main and exclusive way to contribute to the project.

### Setup

[Fork][fork], then clone the repository:

```
git clone git@github.com:your_github_username/dark.git
cd dark
git remote add upstream https://github.com/K-Phoen/dark.git
git fetch upstream
```

Make sure you have the required tools:

```
make dev-env-check-binaries
```

Start a development environment:

```
make dev-env-start
```

This command will start a lightweight Kubernetes cluster using [k3d](https://k3d.io/v5.2.2/),
provisioned with Grafana, Prometheus and Loki.

The following command will run `dark` against this cluster.

```
make run
```

### Making Changes

Start by creating a new branch for your changes:

```
git checkout master
git fetch upstream
git rebase upstream/master
git checkout -b new-feature
```

Make your changes, then ensure that `make lint` and `make test` still pass. If
you're satisfied with your changes, push them to your fork.

```
git push origin new-feature
```

Then use the GitHub UI to open a pull request.

At this point, you're waiting on us to review your changes. We *try* to respond
to issues and pull requests within a few business days, and we may suggest some
improvements or alternatives. Once your changes are approved, one of the
project maintainers will merge them.

We're much more likely to approve your changes if you:

* Add tests for new functionality.
* Write a [good commit message][commit-message].
* Maintain backward compatibility.

[fork]: https://github.com/K-Phoen/dark/fork
[commit-message]: http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html
