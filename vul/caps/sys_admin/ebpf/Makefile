EBPF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

# vmlinux header file
BPFTOOL = $(shell which bpftool || /bin/false)
BTFFILE = /sys/kernel/btf/vmlinux
DBGVMLINUX = /usr/lib/debug/boot/vmlinux-$(shell uname -r)
VMLINUXH = $(EBPF_DIR)headers/vmlinux.h

.PHONY: vmlinuxh generate-ebpf
vmlinuxh: $(VMLINUXH)

$(VMLINUXH):
ifeq ($(wildcard $(BPFTOOL)),)
	@echo "ERROR: could not find bpftool"
	@exit 1
endif
	@if [ -f $(DBGVMLINUX) ]; then \
		echo "INFO: found dbg kernel, generating $(VMLINUXH) from $(DBGVMLINUX)"; \
		$(BPFTOOL) btf dump file $(DBGVMLINUX) format c > $(VMLINUXH); \
	fi
	@if [ ! -f $(BTFFILE) ] && [ ! -f $(DBGVMLINUX) ]; then \
		echo "ERROR: kernel does not seem to support BTF"; \
		exit 1; \
	fi
	@if [ ! -f $(VMLINUXH) ]; then \
		echo "INFO: generating $(VMLINUXH) from $(BTFFILE)"; \
		$(BPFTOOL) btf dump file $(BTFFILE) format c > $(VMLINUXH); \
	fi

BPF_CFLAGS = -target bpf -D__TARGET_ARCH_x86
generate-ebpf: vmlinuxh
	BPF_CFLAGS="$(BPF_CFLAGS)" go generate -x ./vul/caps/sys_admin/ebpf/...
