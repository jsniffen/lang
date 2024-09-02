nasm -f win64 out.asm
link out.obj kernel32.lib legacy_stdio_definitions.lib ucrt.lib /entry:main /subsystem:console
