module github.com/cjr29/go6502

go 1.22.3

require (
	github.com/beevik/cmd v0.2.0
	github.com/beevik/prefixtree v0.3.0
	github.com/cjr29/go6502/asm v0.0.0-20240520005320-9fb32dbc95b2
	github.com/cjr29/go6502/cpu v0.0.0-20240520005320-9fb32dbc95b2
	github.com/cjr29/go6502/disasm v0.0.0-20240520005320-9fb32dbc95b2
	github.com/cjr29/go6502/term v0.0.0-20240525125723-5dc44534dbc7
)

require golang.org/x/sys v0.20.0 // indirect

replace github.com/cjr29/go6502/host v0.0.0 => ./host

replace github.com/cjr29/go6502/asm v0.0.0 => ./asm

replace github.com/cjr29/go6502/cpu v0.0.0 => ../cpu

replace github.com/cjr29/go6502/disasm v0.0.0 => ../disasm
