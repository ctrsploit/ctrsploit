.section .text
.global _start

_start:
loop:
    mov     x8, #220
    mov     x0, #17
    mov     x1, #0
    svc     #0
    b       loop