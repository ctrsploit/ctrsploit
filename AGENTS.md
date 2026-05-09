# Agent Instructions

## Project Skills

This repository keeps shared AI-agent skills under `.opencode/skills` so they
can be used by Opencode and referenced by Codex-compatible agents.

When a user request matches a skill below, open the referenced `SKILL.md` and
follow its instructions:

- `commit-message`: Generate semantic commit messages and commit safely.
  File: `.opencode/skills/commit-message/SKILL.md`
- `record-terminal-session`: Record ctrsploit terminal demonstrations and
  convert asciinema `.cast` captures into project-style `video.svg` or
  `video.gif` artifacts.
  File: `.opencode/skills/record-terminal-session/SKILL.md`

If a skill references files such as `references/`, resolve them relative to that
skill directory.
