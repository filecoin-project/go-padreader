module github.com/filecoin-project/go-padreader

go 1.13

require (
	github.com/filecoin-project/filecoin-ffi v0.0.0-00010101000000-000000000000
	github.com/google/go-cmp v0.4.0 // indirect
	gotest.tools v2.2.0+incompatible
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
