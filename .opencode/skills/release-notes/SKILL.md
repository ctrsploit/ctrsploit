---
name: release-notes
description: Draft ctrsploit release notes for a target version by comparing git tags, mapping merge commits to GitHub PRs, and grouping user-visible changes by module and object.
---

# Release Notes Skill

Use this when generating or updating ctrsploit release notes, changelogs, or
version summaries.

## Workflow

1. Identify the target version and previous release.
   - Prefer explicit user input, for example `v0.26.0-beta.5` since
     `v0.25.1`.
   - If missing, inspect tags with `git tag --sort=version:refname`.
   - Treat an existing tag's release notes as published. Do not backfill an
     already-published `docs/release-notes/<version>.md` for later changes;
     create or update the next target version instead.
   - For beta/pre-release notes, use the previous stable release as the
     comparison baseline, not the previous beta tag. For example,
     `v0.26.0-beta.6` should say `Changes since v0.25.1`, not
     `Changes since v0.26.0-beta.5`.
   - Later beta notes are cumulative snapshots for users upgrading from the
     previous stable release. Carry forward relevant previous beta notes,
     update versioned `blob/<version>/...` links to the new target version,
     then add the new changes since the prior beta.
2. Inspect the release range.
   - `git log --oneline --no-merges <previous>..HEAD`
   - `git log --oneline --merges <previous>..HEAD`
   - `git diff --stat <previous>..HEAD`
3. Map changes to PRs.
   - Parse merge commits such as `Merge pull request #368 ...`.
   - Use PR links in the final notes: `[#368]: https://github.com/ctrsploit/ctrsploit/pull/368`.
   - Attach PR references to the changed object heading, not only a footer.
   - If a commit was pushed directly to the release branch and has no PR,
     attach a commit link instead of forcing it into a nearby PR.
   - If one object includes both PR changes and direct-push commits, list both
     in the heading. Example: `### suid ([#368], [58afd74])`.
4. Write notes under `docs/release-notes/<version>.md`.
   - Use one file per release, for example
     `docs/release-notes/v0.26.0-beta.5.md`.
   - Do not create root-level `changelog.md` unless the repo has adopted a
     persistent root `CHANGELOG.md` convention.
   - Link changed objects to existing repository docs when available.
     Use GitHub absolute `blob/<version>/...` links for release notes that will
     be pasted into GitHub Releases, because release bodies are not rendered
     relative to `docs/release-notes/<version>.md`.
   - Add a `Contributors` section before link reference definitions.
5. Validate.
   - Run `git diff --check -- docs/release-notes/<version>.md`.
   - Review that major commits in the range are represented.
   - Verify an already-published previous notes file has no diff.
   - Verify `Changes since ...` points to the previous stable release for beta
     notes.
   - Verify the new file does not retain stale versioned links from the prior
     beta.

## Structure

Group by module first. Recommended module headings:

- `Vul`
- `Pipe Primitive`
- `Env`
- `Crash`
- `E2E`
- `Skills And Agent Docs`

Inside each module, group by changed object so the object appears once. Put
change type labels inside the object section:

```markdown
## Vul

### cve-2026-43284 ([#367], [#368])

Aliases: `43284`, `dirty-frag`.

New:

- Adds the Dirty Frag Linux kernel vulnerability module.

Improve:

- Adds `escape restart` support and documentation.

Fix:

- Fixes Dirty Frag worker setup.

Test:

- Adds no-`su` privilege escalation e2e coverage.
```

Avoid repeating the same object in multiple headings like
`### New: cve-2026-43284` and `### Improve: cve-2026-43284`; that makes the
notes harder to scan.

## Content Rules

- Emphasize user-visible commands, modules, and behavior.
- Name new commands explicitly, including aliases.
- Include important flags for new commands.
- Write from the perspective of users upgrading from the previous release, not
  from the perspective of intermediate beta or PR history.
- If an object did not exist in the previous release, put its final shipped
  behavior under `New`. Do not add separate `Improve` or `Fix` sections for
  beta-period changes to that same new object unless users of the previous
  release could have observed the earlier behavior.
- Mention compatibility or safety constraints when they change user decisions.
- Keep implementation details only when they explain behavior users will see.
- Include tests and e2e coverage, but keep them after feature/fix content.
- Use concrete paths for important docs or generated assets when helpful.
- Prefer linking the object heading to its README, SKILL, or documented section
  when the repository already has one. For GitHub Release body compatibility,
  use versioned absolute links such as
  `https://github.com/ctrsploit/ctrsploit/blob/<version>/pkg/pipe-primitive/README.md#clean`.
  Do not invent links for commands or objects that do not have a clear
  documentation target.
- Include release contributors from git authors in the compared range. Normalize
  duplicate casing for the same email, exclude bot/merge committers such as
  `GitHub <noreply@github.com>`, and do not treat tool footers such as
  `Made-with: ...` as contributors.
- Prefer GitHub `@username` mentions in the `Contributors` section when the
  contributor account is known, so GitHub release pages can render linked
  contributor profiles and avatars.

## PR Mapping Heuristics

- New vulnerability modules usually map to their feature PR plus later improve
  PRs touching the same module.
- Shared packages such as `pkg/pipe-primitive` may span several PRs; attach all
  relevant PRs to the shared object heading.
- Direct-push commits in the release range must still be represented. Use a
  markdown reference link to the commit short hash:
  - Heading example: `### make e2e ([69ce2f8])`
  - Footer example: `[69ce2f8]: https://github.com/ctrsploit/ctrsploit/commit/69ce2f8`
- Do not assign a direct-push commit to a nearby PR unless the git/GitHub
  history clearly shows that commit belongs to that PR.
- If a PR title is misleading, rely on the commits and touched files rather than
  the title alone.
- If a branch was split or rebased during review, confirm the merge commit range
  from `main` before finalizing the notes.
