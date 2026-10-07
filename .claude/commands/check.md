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

Prefix IDs by source so they don't collide: `V-` for verify (`V-W1`, `V-S1`) and `FE-` for the FE review (`FE-B1`, `FE-D2`). Every finding appears exactly once in the whole report.

1. **Scope**: one line each. Which OpenSpec change was verified (and whether the subagent had to guess), and whether the FE review ran or was skipped.
2. **Must fix**: every verify Critical or Warning, every FE Bug, and any finding that both reports flagged independently (list it once with both IDs, e.g. `V-S1 + FE-U1`). Each entry keeps its ID, `file:line`, tags, and consequence.
3. **Your call**: every other finding, in the order the reports gave them, with the same detail. Do not add your own ranking or opinions.
4. **Clean checks**: one line for each check that found nothing (e.g. "OpenSpec verify: no issues"). If verify found only suggestions, say so in this line and put the suggestions under Your call.

Do not write a separate overlaps section. Do not paraphrase away file references, IDs, or severities. Do not fix anything. If a subagent fails or returns nothing useful, say so plainly instead of filling the gap.