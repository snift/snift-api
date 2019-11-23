# Contributing to Snift API

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Setting Up the project locally](#setting-up-the-project-locally)
- [Running Tests](#running-tests)

## Code of Conduct

We have a code of conduct you can find [here](./CODE-OF-CONDUCT.md) and every
contributor is expected to obey the rules therein. Any issues or PRs that don't
abide by the code of conduct may be closed.

## Setting Up the project locally

To install the project you need to have `go`. Follow the instructions [here](https://golang.org/doc/install) to setup go.

1.  [Fork](https://help.github.com/articles/fork-a-repo/) the project, clone
    your fork:

    ```sh
    # Clone your fork
    git clone https://github.com/<your-username>/snift-api.git

    # Navigate to the newly cloned directory
    cd snift-api
    ```

2.  Your environment needs to be running `go v1.11.4`.
3.  from the root of the project: `go install` to install all dependencies
4.  from the root of the project: `go run main.go`
    - this starts the api server on port `9700` by default

> Tip: Keep your `master` branch pointing at the original repository and make
> pull requests from branches on your fork. To do this, run:
>
> ```sh
> git remote add upstream https://github.com/snift/snift-api.git
> git fetch upstream
> git branch --set-upstream-to=upstream/master master
> ```
>
> This will add the original repository as a "remote" called "upstream," then
> fetch the git information from that remote, then set your local `master`
> branch to use the upstream master branch whenever you run `git pull`. Then you
> can make all of your pull request branches based on this `master` branch.
> Whenever you want to update your version of `master`, do a regular `git pull`.

## Running Tests

Also, make sure to run the tests and lint the code before you commit your
changes.

```sh
go test ./.../
```

Thank you for taking the time to contribute! 👍
