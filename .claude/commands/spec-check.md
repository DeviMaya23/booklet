---
description: Check whether changes made outside the OpenSpec flow leave any spec in openspec/specs outdated. Flags only, no edits.
argument-hint: [optional: commit range or sha, default main..HEAD plus uncommitted changes]
allowed-tools: Bash(git log:*), Bash(git show:*), Bash(git diff:*), Bash(git status:*), Read, Glob, Grep
---

Arguments: $ARGUMENTS (an optional commit range or sha; if empty, use `main..HEAD` plus any uncommitted changes).

Decide whether the code changes warrant a change to the specs in `openspec/specs/`. Nothing gets edited, created, or deleted. This is a heuristic check, so show your work.

## 1. Scope

- List the commits (`git log <range> --reverse --format='%h %s'`) and, for each, its files (`git show --name-only --format= <sha>`). Include uncommitted changes (`git diff HEAD --name-only`) when no range was given.
- Skip any commit that touches files under `openspec/`: it came through the proposal flow and is covered by verify, not this check. Say which commits you skipped.

## 2. Classify each remaining commit

Read the diff and put it in one bucket:

- **Observable change**: user-visible behavior, an API request/response or status code, validation or limits, error handling, what is stored or deleted, ordering or filtering, permissions.
- **Spec-conforming fix**: restores behavior a spec already describes. Confirm by finding that requirement in step 3. If you cannot find it, treat it as an observable change.
- **Internal only**: refactor, tests, styling, logging, dependency bump, tooling, with no observable difference.

For **internal only**, stop for that commit and say why in one line. Do not search specs for it.

## 3. Find and read candidate specs

For each commit in the first two buckets:

1. Pull specific search terms from the diff: endpoint paths, route names, component and hook names, error messages and toast text, field names, constants and limits. Do not search generic words that appear in most specs (`PUT`, `input`, `fail`, `name`, `upload` on its own). Combine a domain word with a specific one (e.g. `upload_url`, `Failed to upload`, `thumbnail_upload_url`).
2. `Grep` `openspec/specs/*/spec.md` for them, one narrow search at a time. If a term matches more than about five specs, it is too broad: narrow it instead of reading the hits. Take the few best matches, not every hit. Also look at `openspec/specs/` directory names for obvious capability matches the grep missed.
3. Open each candidate with `Read` and read all of its requirements (`SHALL`) and scenarios (`#### Scenario`). For a very large spec you may read only the requirements that relate to the diff, but then say so in the report. Grep lines alone never count as reading a spec, and a spec you did not open cannot get a Silent verdict.
4. Also list the non-archived folders under `openspec/changes/`. If an active change already covers this behavior, say so: the spec update is already planned.

## 4. Verdict per candidate spec

- **Contradicts**: a requirement or scenario now says something the code no longer does. Quote it.
- **Uncovered**: the code has new observable behavior that no scenario describes. Say which capability it belongs to.
- **Silent**: the spec does not speak to this behavior, and no change is needed.

## 5. Report

Short and plain. Per commit: sha, subject, bucket and the one-line reason. Then:

- **Specs read**: every spec path you opened, each marked "in full" or "sections only: <which requirements>". Add a separate "Not opened" line for specs that matched a search but that you did not open, so the reader can see exactly what was and was not checked.
- **Findings**: for each Contradicts or Uncovered, the spec path, the requirement or scenario, and a one-line suggested delta (what it should say now).
- **Verdict**: either "no spec change warranted" or "spec change warranted: <spec names>".

When a change is warranted, mention that a small `/opsx:propose` change is how this repo updates specs. Do not run it, and do not edit `openspec/specs/` directly.

If nothing matched and you are unsure, say that instead of defaulting to "no spec change". A missing match is not proof that no spec applies.
