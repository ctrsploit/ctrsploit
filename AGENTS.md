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
- `dqd-lab`: Start, stop, inspect, and connect to ctrsploit dqd lab
  environments from github.com/ctrsploit/dqd.
  File: `.opencode/skills/dqd-lab/SKILL.md`
- `release-notes`: Draft ctrsploit release notes for target versions by
  comparing tags, grouping user-visible changes, and preserving published
  release notes.
  File: `.opencode/skills/release-notes/SKILL.md`

If a skill references files such as `references/`, resolve them relative to that
skill directory.
