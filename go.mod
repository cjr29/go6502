module github.com/cjr29/go6502

go 1.21.6

require (
	github.com/cjr29/go6502/asm v0.0.0-20240519135954-448b93fea12a
	github.com/cjr29/go6502/host v0.0.0-20240519135954-448b93fea12a
)

require (
	github.com/beevik/cmd v0.2.0 // indirect
	github.com/beevik/prefixtree v0.3.0 // indirect
	github.com/cjr29/go6502/cpu v0.0.0-20240519135954-448b93fea12a // indirect
	github.com/cjr29/go6502/disasm v0.0.0-20240519135954-448b93fea12a // indirect
	github.com/cjr29/go6502/term v0.0.0-20240519135954-448b93fea12a // indirect
	golang.org/x/sys v0.20.0 // indirect
)

replace github.com/cjr29/go6502/cpu v0.0.0 => ./cpu

replace github.com/cjr29/go6502/dashboard v0.0.0 => ./dashboard

replace github.com/cjr29/go6502/host v0.0.0 => ./host

replace github.com/cjr29/go6502/asm v0.0.0 => ./asm

replace github.com/cjr29/go6502/term v0.0.0 => ./term
