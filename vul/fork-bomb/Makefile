all:
	nasm -f elf64 -o bomb_amd64.o bomb_amd64.asm
	ld -o bomb_amd64 bomb_amd64.o
	rm bomb_amd64.o
	# TODO: arm64