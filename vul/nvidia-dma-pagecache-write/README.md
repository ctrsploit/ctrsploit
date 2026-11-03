---
maintainer:
    - Hpd0ger
---
# NVIDIA GPU DMA page-cache write

CVE-2026-47591, CVE-2026-47595 and CVE-2026-47599 share one root cause in the
NVIDIA driver memory-descriptor and DMA layer: a caller can pin a read-only
file page, but the read-only attribute is not propagated to the device mapping
and the page is mapped for bidirectional DMA. The GPU can therefore write the
host page cache without ever gaining CPU write access to the file.

An attacker uses this to poison the host `libnvidia-ml.so.<version>` page cache
through a writable GPU mapping. The next root process on the host that maps the
poisoned library executes the attacker's payload, which turns the primitive into
container escape.

| CVE | CWE | Root cause |
| --- | --- | --- |
| CVE-2026-47591 | CWE-863 | `RmcCreateOsDescriptor` pins a read-only page when the caller sets `USER_READ_ONLY` |
| CVE-2026-47595 | CWE-281 | `MEMDESC_FLAGS_USER_READ_ONLY` does not become `DEVICE_READ_ONLY` |
| CVE-2026-47599 | CWE-281 | `nv-dma.c` unconditionally uses `DMA_BIDIRECTIONAL` for OS-descriptor pages |

## 1. Affected versions

| Branch | Vulnerable | Fixed | Notes |
| --- | --- | --- | --- |
| R535 | all versions | not stated | 535.183.01 is confirmed exploitable; no fixed release is stated, so the whole branch is treated as vulnerable |
| R580 | `< 580.178.04` | `>= 580.178.04` | |
| R595 | `< 595.91.07` | `>= 595.91.07` | |
| R610 | `< 610.57.04` | `>= 610.57.04` | |
| R615 | `< 615.71.09` | not asserted | the fixed boundary is only stated in the analysis README; at or above it the probe reports undetermined |
| vGPU 19 | affected | not stated | no reliable read-only vGPU marker |
| MIG | not assessed | not assessed | not discussed in the analysis |

Other branches are reported as undetermined, never as fixed.

## 2. Detection scope

The probe is read-only. It never opens an NVIDIA device node, never issues an
ioctl, and never runs `nvidia-smi`, because `nvidia-smi` opens `/dev/nvidiactl`
and issues NVML ioctls.

Version sources, in priority order:

1. `/proc/driver/nvidia/version` and `/sys/module/nvidia/version`, which reflect
   the loaded host kernel module;
2. `libnvidia-ml.so.<driver-version>` filenames in `LD_LIBRARY_PATH` and the
   common NVIDIA/CUDA library directories.

Kernel-module evidence wins over library filenames. When only library filenames
are available, a vulnerable version is still reported as vulnerable, but an
all-fixed set is downgraded to undetermined: the library may come from the
container image rather than from the host driver.

The prerequisite for reachability is that `/dev/nvidiactl`, `/dev/nvidia0` and
`/dev/nvidia-uvm` all exist, are character devices, and pass
`access(R_OK|W_OK)`. The probe does not check `nvidia1` or other GPU numbers, so
a container that exposes only a non-zero GPU index is reported as unreachable.

The overall check is `[Y]` only when the driver version is vulnerable and the
three device nodes are reachable. A `[N]` result with a `Response` block
distinguishes:

- the version is fixed but the device nodes are not all visible;
- the version cannot be determined from read-only sources (undetermined, not
  fixed);
- the version is vulnerable but the escape interface is unreachable.

## 3. Usage

```console
ctrsploit checksec nvidia-dma-pagecache-write
ctrsploit checksec 47591
ctrsploit vul nvidia-dma-pagecache-write checksec
ctrsploit module nvidia-driver nvidia-dma-pagecache-write checksec
```

Add `--debug` to see every version source, the parsed version, its verdict and
its trust tier.

`ctrsploit checksec auto` can only report a boolean, so an undetermined version
appears as `[N]` there. Use the dedicated command when the difference between
"fixed" and "undetermined" matters.

## 4. Limitations

- Different builds of the same driver version exist; the probe classifies by
  version only and does not use file size, MD5 or byte signatures.
- A library shipped by the container image does not prove anything about the
  host kernel module, so library-only fixed evidence is reported as
  undetermined.
- The driver version alone does not prove exploitability: a downstream
  backport with the same version string is not distinguishable.
- Whether a root process on the host maps the poisoned library cannot be
  determined from inside the container.
- vGPU 19 has no stated fixed boundary and no reliable read-only marker; MIG
  is not covered by the analysis.
