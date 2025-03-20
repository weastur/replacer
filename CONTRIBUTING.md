# Contribution Guide

Thank you for your interest in contributing to this project! Your contributions help improve and maintain the project
for everyone. Please follow the guidelines below to ensure a smooth contribution process.

## 1. Getting Started

1. **Fork the Repository** – Click the fork button on GitHub and clone your fork locally:

   ```sh
   git clone https://github.com/weastur/replacer.git
   cd replacer
   ```

1. **Install Dependencies**

    - [Go](https://golang.org/dl/)
    - [direnv](https://direnv.net/) (optional but recommended)
    - [pre-commit](https://pre-commit.com/) (optional but recommended)

1. **Install pre-commit hooks** – This will ensure your code is formatted and linted before committing:

   ```sh
   pre-commit install
   ```

1. **Allow direnv** – If you have direnv installed, allow the `.envrc` file:

   ```sh
   direnv allow
   ```

1. **Use make** – The project is using self-describing makefile:

   ```sh
   make
   ```

## 2. Making Changes

- Follow the existing code style and structure.
- Write clear, self-documented code and add comments where necessary.
- Ensure your changes do not introduce new errors or warnings.

## 3. Commit and Push

1. **Follow commit style (gitlint will do all checks, just ensure git hooks are installed)**:

   ```sh
   git commit -m "Add new awesome feature"
   ```

2. **Push Your Changes**:

   ```sh
   git push origin your-branch
   ```

## 4. Submitting a Pull Request

1. Open a pull request with a clear description of your changes.
2. Ensure your PR follows best practices and passes all checks.
3. Request a review from the maintainers.

## 5. Code Review Process

- Be open to feedback and make necessary revisions.
- Address requested changes promptly.
- Once approved, your PR will be merged!

## 6. Reporting Issues

If you find bugs or have feature requests, please
[open an issue](https://github.com/weastur/replacer/issues) with relevant details.

---

Thank you for your contributions! 🚀
