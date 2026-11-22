# ADR 0061: A documentation website built by MkDocs and deployed to GitHub Pages

- Status: Accepted
- Effective date: 2026-11-22
- Deciders: eduardo.ferro

## Context

Issue #39 asks to publish the user documentation as a website on GitHub Pages, built with MkDocs from the same Markdown files the repository already has, with search; a CI job must build it on every push, fail on broken links and deploy from `main`. The user documents live in three places (`README.md` and `CHANGELOG.md` at the root, the guides in `docs/`, the catalogue index in `scenarios/README.md`) and link to files a website would not publish: the specification, the ADRs, `AGENTS.md`, the `.ubc` models, `examples/` and the licence. MkDocs only renders what is inside its `docs_dir`, which cannot be the directory that holds `mkdocs.yml`, and treats a link to a file outside it as broken.

## Decision

- **MkDocs 1.6.1 with its built-in `mkdocs` theme and `search` plugin.** Seven pages need no more; a third-party theme would be another dependency to pin and follow. MkDocs and its dependencies are pinned in `tools/docsite/requirements.txt` and installed by `make docs` into a virtual environment in `.tools/`, like the Go tools.
- **What is published:** the pages listed in the `nav` of `mkdocs.yml` — the README as the home page, the tutorial, the modelling guide, the how-to guides, the reference, the scenario catalogue index and the changelog. The specification, the vision, the ADRs and `AGENTS.md` are for contributors and stay on GitHub.
- **The site is assembled, not moved.** `tools/docsite` (a Go program, like `tools/releasenotes`, ADR 0060) reads `repo_url`, `docs_dir` and the `nav` pages from `mkdocs.yml`, and copies each page into `build/docs` at its repository path, so links between published pages keep working unchanged. Every other relative link is pointed at the file on GitHub (`<repo_url>/blob/main/<path>`, or `tree` for a directory), keeping its anchor; a link to a file that does not exist, or outside the repository, fails the tool. The Markdown sources stay where they are and keep working on GitHub.
- **Broken links fail twice.** `mkdocs.yml` sets `strict: true` and every `validation` check to `warn`, so a link to a missing page or a missing anchor fails `mkdocs build`. The Go test suite also assembles the real site, so a dead link to a repository file fails `make validate` without Python.
- **`make docs`** runs the tool and `mkdocs build --strict` into `build/site`; `build/` is ignored by git. It is not part of `make validate`, which keeps needing only Go.
- **Deployment:** the `Docs` workflow (`.github/workflows/docs.yml`) runs `make docs` on every push to `main` and every pull request; on a push it uploads `build/site` with `actions/upload-pages-artifact` and a second job publishes it with `actions/deploy-pages` to `https://eferro.github.io/uncertain-bc/`. Pages must be enabled once in the repository settings with "GitHub Actions" as the source, which only the maintainer can do.

## Consequences

- The guides, the README and the changelog have one source: editing them updates both GitHub and the website on the next push to `main`.
- A new page is published by adding it to the `nav`; the tool then copies it and checks its links.
- The tool reads `mkdocs.yml` line by line instead of parsing YAML, so the keys it reads must keep their simple one-line form.
- The website shows released and unreleased documentation alike: it follows `main`, not the latest tag.

## Alternatives considered

- **Material for MkDocs:** nicer navigation and search, but a larger dependency set for seven pages.
- **Moving the documents under one `docs_dir`, or symlinking them:** moving would break the paths users and the test suite know (the README must stay at the root); symlinks are fragile on Windows, and links to unpublished files would still be broken.
- **`exclude_docs` with the repository root as `docs_dir`:** needs `mkdocs.yml` in a subdirectory and a long exclusion list, and still breaks the links to models and ADRs.
- **The `gh-pages` branch with `mkdocs gh-deploy`:** pushes generated files into the repository; the official Pages actions deploy an artifact instead.
