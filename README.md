# `folkctl`

> [!caution]
> **Status:** Under development


## ℹ️ About

Folkctl is a CLI tool for finding and copying national identity numbers (fødselsnummer) for BankID test users that also exist in the Norwegian National Population Register (DSF) test database.

This project is the CLI tool for the [Folkomaten](https://github.com/olefredrik/Folkomaten) project.

## ✨ Features

- Search for users in the Folkctl test database.
- Copy the fødselsnummer to the clipboard.
- Copy the name to the clipboard.
- Copy the all information to the clipboard.

## 🛠️ Technologies

The project is built with:

- [Cobra](https://cobra.dev)
- [fzf](https://github.com/junegunn/fzf)

## 📋 Requirements

Before starting, make sure `Go` is installed on your system.

```bash
go help
```

## 📦 Installation

### Manual installation

```bash
git clone https://github.com/antoniorodr/folkctl
cd folkctl
go build && sudo mv folkctl /usr/local/bin/
```

<!-- ### Homebrew installation -->
<!---->
<!-- ```bash -->
<!-- brew install your-project -->
<!-- ``` -->

## 🚀 Getting Started

Once installed, run:

```bash
folkctl
```

## ❤️ Do you like my work?

If you find the project useful, you can support the author here:

[![GitHub Sponsor](https://img.shields.io/badge/Sponsor_on_GitHub-30363D?logo=github&style=for-the-badge)](https://github.com/sponsors/antoniorodr)
