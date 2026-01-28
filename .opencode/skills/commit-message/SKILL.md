---
name: commit-message
description: Generate semantic commit messages and commit safely
---

# Commit Message Skill

## Goal
- Generate a clear, professional Git commit message based on staged changes.
- Follow Conventional Commits with a concise, intent-focused subject line.
- Keep commits atomic; suggest splitting when changes are unrelated.
- Provide clear guidance for both English and Chinese-speaking developers.

## Workflow
1. Analyze staged changes (files, logic, and config).
2. Check atomicity:
   - If unrelated changes exist, propose split groupings.
   - Otherwise, proceed with a single commit message.
3. Draft output using the template below.
4. If committing is requested, follow the commit execution protocol exactly.

## Commit Message Format (Conventional Commits)
```
<type>(<scope>): <subject>

<body>
```

### Title Line
- Type: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert
- Scope: optional module or area (e.g., "vul/cve-2019-5736", "env", "pkg/runtime")
- Subject: imperative mood, no period, 50-72 characters max
- Example: `feat(vul/cve-2019-5736): add exploit for runc escape`

### Body (Use Sparingly)
- Prefer title-only when it is sufficient.
- Use a short body only for complex or non-obvious changes.
- Focus on intent/impact (why), not exhaustive details (what).
- Wrap at 72 characters.

## Common Types for This Project
- `feat`: New vulnerability detection or exploit functionality
- `fix`: Bug fixes in vulnerability checks or exploits
- `docs`: Documentation updates for vulnerabilities or usage
- `test`: Adding or updating tests for security checks
- `chore`: Maintenance tasks, dependency updates
- `refactor`: Code restructuring without functional changes
- `perf`: Performance improvements in detection or exploitation

## Examples
Good titles:
- `feat(env): add namespace detection capabilities`
- `fix(vul/cve-2022-0492): resolve nil pointer dereference`
- `docs(vul): update cve-2024-0132 usage instructions`
- `test(pkg/syscall): add unit tests for seccomp detection`

## Output Template
```
Files staged. Generating commit message:
Chinese Summary: <A clear Chinese summary; may include key files/vars>
English Draft:
<type>(<scope>): <subject>

<optional body>
```

## Commit Execution Protocol
- Always use a file-based commit message.
- Create `.commit-message` in the repo root using file editing tools.
- The file content must match the English Draft exactly.
- Commit with `git commit -F .commit-message`.
- Delete `.commit-message` after commit completes.

## Safety
- Do not change git config.
- Do not use force push, hard reset, or amend unless explicitly requested.
- For this security-focused project, ensure commit messages accurately reflect security implications.
- When dealing with vulnerability-related changes, be precise about the CVE or vulnerability reference.

## Best Practices
- Be specific about the vulnerability or security aspect in the scope
- Use CVE identifiers in parentheses when applicable: `fix(cve-2022-0492): ...`
- For multiple CVEs in one change, use the most relevant one in the scope
- When adding new exploits, use `feat` type: `feat(vul): add cve-2024-xxxx exploit`
- When updating existing security checks, use `feat` for new functionality, `fix` for corrections
