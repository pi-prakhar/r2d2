# 🤝 Contributing to R2D2

Thank you for your interest in contributing to R2D2! This document outlines the standards and practices we follow in this project.

## 🚀 Getting Started

1. **Fork the Repository**: Start by forking the repository to your GitHub account.
2. **Clone Your Fork**: Clone your fork to your local machine.
   ```bash
   git clone https://github.com/YOUR_USERNAME/r2d2.git
   cd r2d2
   ```
3. **Add Upstream Remote**: Set up the original repository as your "upstream" to keep your fork in sync.
   ```bash
   git remote add upstream https://github.com/pi-prakhar/r2d2.git
   ```

## 🌿 Branch Naming Convention

We follow a standard branch naming convention to keep our work organized:

```
<type>/<description>
```

### Branch Types

| Type | Description |
|------|-------------|
| **feature/** | New features or enhancements |
| **fix/** | Bug fixes |
| **docs/** | Documentation updates |
| **refactor/** | Code refactoring without functional changes |
| **test/** | Adding or updating tests |
| **chore/** | Maintenance tasks, dependency updates, etc. |
| **style/** | Code style or formatting changes |

### Examples

```
feature/pod-level-monitoring
fix/deployment-restart-issue
docs/improve-installation-guide
refactor/kubernetes-client
test/watch-tags-unit-tests
chore/update-dependencies
style/code-formatting
```

## 💬 Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification for our commit messages.

### Format

```
<type>: <description>

[optional body]

[optional footer]
```

### Commit Types

| Type | Description |
|------|-------------|
| **feat** | A new feature |
| **fix** | A bug fix |
| **docs** | Documentation changes |
| **style** | Changes that do not affect the meaning of the code (formatting, etc.) |
| **refactor** | Code changes that neither fix a bug nor add a feature |
| **perf** | Performance improvements |
| **test** | Adding or fixing tests |
| **chore** | Changes to the build process, tools, or dependencies |

### Examples

```
feat: add pod-level monitoring to watch-tags command
fix: prevent deployment restart when namespace is empty
docs: update installation guide for Windows users
style: apply consistent code formatting
refactor: simplify Kubernetes client initialization
perf: optimize deployment info retrieval
test: add unit tests for watch-tags command
chore: update Go modules to latest versions
```

### Additional Guidelines

- Use the imperative mood ("add" not "added" or "adds")
- Don't capitalize the first letter of the description
- No period at the end of the description
- Reference issues in the footer with "Fixes #123" or "Closes #123"

## 🔄 Pull Request Process

1. **Update Your Fork**: Make sure your fork is up-to-date with the upstream repository.
   ```bash
   git fetch upstream
   git checkout master
   git merge upstream/master
   ```

2. **Create a Branch**: Create a new branch following our naming convention.
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make Changes**: Make your changes, following our coding standards.

4. **Test Your Changes**: Ensure your changes don't break existing functionality.

5. **Commit Your Changes**: Use meaningful commit messages following our commit guidelines.

6. **Push to Your Fork**: Push your branch to your forked repository.
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request**: Open a PR from your fork to the main repository.
   - Provide a clear description of the changes
   - Reference any related issues
   - Add reviewers if appropriate

8. **Address Review Feedback**: Make any requested changes and push them to your branch.

9. **Merge**: Once approved, your PR will be merged into the main branch.

## 📏 Code Style and Standards

- Follow Go best practices and idiomatic Go code style
- Use `gofmt` to format your code
- Add comments for public functions and packages
- Write unit tests for new functionality

## 🐛 Reporting Issues

When reporting issues, please include:

- A clear description of the issue
- Steps to reproduce the problem
- Expected behavior
- Actual behavior
- Environment details (OS, Go version, etc.)
- Screenshots or logs if applicable

## 📜 License

By contributing to R2D2, you agree that your contributions will be licensed under the project's [LICENSE](LICENSE).

Thank you for contributing to R2D2! 