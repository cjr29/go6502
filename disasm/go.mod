module github.com/cjr29/go6502/disasm

go 1.22.3

require github.com/cjr29/go6502/cpu v0.0.0-20240519131451-46a8867f26b0

require github.com/cjr29/go6502/asm v0.0.0-20240519131451-46a8867f26b0 // indirect

replace github.com/cjr29/go6502/cpu v0.0.0 => ../cpu

replace github.com/cjr29/go6502/asm v0.0.0 => ../asm
