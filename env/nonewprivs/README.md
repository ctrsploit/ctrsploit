# env/no-new-privs

Show the current process `NoNewPrivs` status from `/proc/self/status`.

When `NoNewPrivs` is enabled, `execve` cannot grant additional privileges.
That means SUID programs and file capabilities cannot elevate the process, so
SUID-based local privilege-escalation checks and exploit fallbacks should treat
the environment as blocked.

## Usage

```shell
$ ctrsploit env no-new-privs
===========NoNewPrivs===========
[Y]  NoNewPrivs enabled        # current process cannot gain new privileges across execve
[Y]  SUID privilege gain blocked       # SUID and file capabilities cannot raise privileges when NoNewPrivs is enabled
source:                 /proc/self/status
```

Aliases:

```shell
$ ctrsploit env nnp
$ ctrsploit env no-new-privilege
$ ctrsploit env no-new-privileges
```

The machine-readable output contains:

```json
{"status_path":"/proc/self/status","enabled":true}
```
