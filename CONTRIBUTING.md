# contributing

> ⌨️ Contributing guide for the cerbero project.

## Table of contents

<!--toc:start-->
- [contributing](#contributing)
  - [Table of contents](#table-of-contents)
  - [Before contributing](#before-contributing)
  - [Bootstrapping a development environment](#bootstrapping-a-development-environment)
    - [Development environment for cerbero-web](#development-environment-for-cerbero-web)
  - [Workflow](#workflow)
    - [Branching Strategy](#branching-strategy)
        - [Create a branch](#create-a-branch)
        - [Make changes](#make-changes)
        - [Create a pull request](#create-a-pull-request)
        - [Merge your pull request](#merge-your-pull-request)
        - [Delete your branch](#delete-your-branch)
          - [Checkout into main](#checkout-into-main)
          - [Pull changes in main](#pull-changes-in-main)
          - [Delete your branch locally](#delete-your-branch-locally)
          - [Prune the remote](#prune-the-remote)
    - [Commit Convention](#commit-convention)
      - [Commit Message Header](#commit-message-header)
        - [Type](#type)
        - [Scope](#scope)
      - [Commit Message Body](#commit-message-body)
      - [Commit Message Footer](#commit-message-footer)
      - [More information about conventional commits](#more-information-about-conventional-commits)
    - [Pull Requests](#pull-requests)
      - [When to pull request](#when-to-pull-request)
    - [Issues](#issues)
<!--toc:end-->

## Before contributing

No matter if you are a inside `K!nd4SUS` or a random dude that suddenly wants to contribute to this project, before making any change to the codebase, please read this document.

## Bootstrapping a development environment

### Development environment for cerbero-web

If you wish to contribute to cerbero-web, you can follow our [frontend development quick-start guide](/web/frontend/README.md#quick-start) and our [backend development quick-stark guide](/web/backend/README.md#quick-start). You will have a working development environment in no time.

## Workflow

### Branching Strategy

> This project follows the [GitHubFlow](https://docs.github.com/en/get-started/quickstart/github-flow) branching strategy.

Following is a summary of a typical GitHubFlow workflow:

##### Create a branch

Create a branch (from `main`) in your repository. The branch name should be short and descriptive, for example: `increase-test-timeout` or `add-code-of-conduct`.

```sh
git checkout -b <name_of_the_branch>
```

If a branch targets a specific issue **the name of the branch should begin with the issue_id** e.g. `123-fix-users-endpoint`.

##### Make changes

On your branch, make the desired changes to the repository, then commit and push your changes to your branch.

When committing your changes, make sure to follow the guidelines described in the <a href="#commits">commits section</a>.

##### Create a pull request

When you create a pull request, **include a summary of the changes** and what problem they solve.

##### Merge your pull request

Once your pull request has been approved, merge your pull request. This will automatically merge your branch so that your changes appear on the default branch.

##### Delete your branch

After you merge your pull request, **delete your branch**. Deleting the branch on GitHub can be done in the branches section of the repository.

To delete your branch locally you can run the following commands:

###### Checkout into main

```sh
git checkout main
```

###### Pull changes in main

```sh
git pull
```

###### Delete your branch locally

```sh
git branch -d <name_of_the_branch>
```

###### Prune the remote

This is **optional** and will remove references to remote branches that were deleted.

```sh
git remote prune origin
```

### Commit Convention

> This project follow the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) specification, this increases consistency and readability of commits.

> Following is a summary of the conventional commits specification, modified to fit the needs of this project.

Each commit message consists of a **header**, a **body**, and a **footer**.

```
<header>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

#### Commit Message Header

```
<type>(<scope>): <short summary>
  │       │             │
  │       │             └─⫸ Summary in present tense. Not capitalized. No period at the end.
  │       │
  │       └─⫸ Commit Scope: components|controllers|db|hooks|middlewares|models|pages|readme|services|utils|<workflow_name>|<test_name>
  │
  └─⫸ Commit Type: build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test
```

The `<type>` and `<summary>` fields are mandatory, the `(<scope>)` field is optional.

##### Type

MUST be one of the following:

* **build**: Changes that affect the build system
* **chore**: Maintain project or external dependencies
* **ci**: Changes to CI configuration files and scripts
* **docs**: Documentation only changes
* **enhance**: A code change that improves an already existing feature
* **feat**: A code change that adds a new functionality to the application
* **fix**: A code change that fixes incorrect code (bugs)
* **perf**: A code change that improves performance
* **refactor**: A code change that doesn't alter the behaviour of the application
* **revert**: Revert a previous change
* **style**: Fix code style, trailing spaces, semi-colons, tab size...
* **test**: Add missing tests or correct existing ones

##### Scope

The scope is the part of the codebase where the changes happened e.g. `feat(components): create IncredibleButton component`.

- If the change doesn't target any particular scope then the commit scope can be empty e.g. `chore: update beautiful stuff`.

- If a commit changes multiple parts of the codebase then an `*` sign can be used as the scope specifier.

#### Commit Message Body

- Use imperative, present tense: “change” not “changed” nor “changes”.

- Include motivation for the change and contrasts with previous behavior.

#### Commit Message Footer

All breaking changes have to be mentioned in footer with the description of the change, justification and migration notes (e.g. `BREAKING CHANGE: desc...`).

- If a commit targets a specific issue, the issue_id must be specified in the footer e.g. `Closes #123`, in case of multiple issues `Closes #123, #124, #125`.

#### More information about conventional commits

Please take a look at the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) specification before committing changes.

### Pull Requests

> Pull requests are the final part of this workflow and they allow contributors to **review and share opinions on code** with each other. Furthermore such mechanism opens the doors to **automated workflow runs**.

#### When to pull request

- A pull request should only be opened when the work is *done* and ready for production.

- If a pull request doesn't pass automated testing/linting/building **it shouldn't be merged**: fix the problems and then push your fixes again until it passes.

### Issues

> Issues can be opened for everything that has to do with the project: from asking questions to requesting new fetures or bug-fixes.

Issues should be labeled as follows:

- A `priority` label
    - `priority: 0` &larr; **Highest**
    - `priority: 1`
    - `priority: 2`
    - `priority: 3`
    - `priority: 4` &larr; **Lowest**
- A `scope` label
    - One of the scopes described in the [commit scope section](#scope).
- A `type` label
    - One or more of the types described in the [commit type section](#type).

