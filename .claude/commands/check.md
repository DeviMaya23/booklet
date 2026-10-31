---
description: Run OpenSpec verify and the FE review on parallel subagents and report what each found. Flags only, no edits.
argument-hint: [optional: OpenSpec change name]
allowed-tools: Agent, Bash(git status:*), Bash(git diff:*), Read, Glob, Grep
---

Arguments: $ARGUMENTS (an optional OpenSpec change name; if empty, verify the active change).

Run openspec verify and fe-reviewer on subagents, in parallel, and come back with their findings. Nothing gets edited, created, or deleted.

**Verify subagent** (general-purpose). Tell it to read `.claude/skills/openspec-verify-change/SKILL.md` and follow it for the change "$ARGUMENTS" (the active change if empty). It cannot ask questions: if the skill would ask something, it picks the most reasonable reading, says which at the top of its report, and carries on.

**FE subagent** (`fe-reviewer`). Tell it to review the current uncommitted changes under `frontend/` in diff mode, following `.claude/review/fe-review-rubric.md`. First check `git diff --name-only` and `git diff --staged --name-only`: if no changed file is under `frontend/`, skip it and say so.

When both are back, report what each found, short and plain, one section per report. Keep each finding's own ID (verify's C/W/S, the FE review's B/D/S) and file reference, and say which report it came from. Do not add your own ranking, must-fix list, or overlap section. If a subagent fails or returns nothing useful, say so.
