---
maintainer:
    - ssst0n3
---
# env/suid

Find and list SUID files visible to the current process. This is useful when
checking whether local privilege-escalation flows can use a root-owned SUID
helper, and when auditing container images or lab VMs for unexpected SUID
programs.

## Usage

```shell
$ ctrsploit env suid
===========SUID Files===========
total:                  6
[/usr/bin/chfn]
path:                   /usr/bin/chfn
mode:                   -rwsr-xr-x
owner:                  0:0
size:                   72792
usability:              root-owned executable
```

By default, the command scans common executable directories such as `/bin`,
`/sbin`, `/usr/bin`, `/usr/sbin`, `/usr/lib`, and `/usr/libexec`.

Scan the whole filesystem:

```shell
$ ctrsploit env suid --all
```

Scan explicit paths:

```shell
$ ctrsploit env suid --path /bin,/usr/bin,/usr/local/bin
```

Skip additional directories while scanning:

```shell
$ ctrsploit env suid --all --skip /mnt,/media
```

The mode field uses `ls -l` style permissions, so SUID files show `s` or `S`
in the owner execute position, for example `-rwsr-xr-x`.
