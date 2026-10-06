---
description: Run OpenSpec verify and the FE review in parallel subagents, then report one merged summary. Flags only, no edits.
argument-hint: [optional: OpenSpec change name]
allowed-tools: Agent, Bash(git status:*), Bash(git diff:*), Read, Glob, Grep
---

Arguments: $ARGUMENTS (an optional OpenSpec change name; if empty, verify the active change).

Run two independent checks as parallel subagents, then merge their findings. Do not edit, create, or delete any file.

## Step 1: decide what to run

Run `git status` and `git diff --name-only` (and `git diff --staged --name-only`).

- The OpenSpec verify always runs.
- The FE review runs only if at least one changed file is under `frontend/`. If none are, skip it and say so in the report.

## Step 2: launch the subagents in the same turn

Launch both at once so they run in parallel. Each starts with a fresh context.

**OpenSpec verify subagent** (general-purpose). Prompt it with:

> Read `.claude/skills/openspec-verify-change/SKILL.md` and follow it exactly to verify an OpenSpec change. The change to verify: "$ARGUMENTS" (if empty, use the active change, i.e. the one the current branch or most recently modified change under `openspec/changes/` refers to). You cannot ask questions. If the skill would ask which change to use or ask anything else, pick the most reasonable reading, state which one you picked at the top of your report, and carry on. Do not edit, create, or delete any file. Return the full verification report, including every issue with its severity and any file references.

**FE review subagent** (`fe-reviewer`, diff mode). Prompt it with:

> Review the current uncommitted changes under `frontend/` in diff mode, following `.claude/review/fe-review-rubric.md`. Return the full report in the rubric's diff-mode output format.

## Step 3: report one merged summary

When both finish, write a single report:

1. **Scope**: which OpenSpec change was verified (and whether the subagent had to guess), and whether the FE review ran or was skipped.
2. **Look at these first**: the top three findings across both reports, whatever their source.
3. **OpenSpec verify**: the findings, condensed, keeping severities and file references.
4. **FE review**: the findings, condensed, keeping tiers, `file:line`, and the new/pre-existing tags.
5. **Overlaps**: anything both reports flagged, or where the two disagree.

Do not paraphrase away file references or severities. Do not fix anything. If a subagent fails or returns nothing useful, say so plainly instead of filling the gap.
