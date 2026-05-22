---
maintainer:
    - ssst0n3
---
# Crash

`pkg/crash` provides restart-trigger strategies for code paths that need the
current container to exit and rely on an external restart policy to start it
again.

Available methods:

- `auto`: expand to the default chain.
- `all`: expand to every built-in method.
- `cgroup-kill`: write `1` to cgroup v2 `cgroup.kill` for the target process.
- `sigkill`: send `SIGKILL` to the target process, usually PID 1.
- `kill-all`: send `SIGKILL` to other visible processes so PID 1 can exit when
  it depends on child process status.
- `oom`: bias OOM selection toward the target process and allocate memory until
  the configured or detected limit is exceeded.

The default chain is `cgroup-kill,sigkill,kill-all`.
The `all` chain is `cgroup-kill,sigkill,kill-all,oom`.
