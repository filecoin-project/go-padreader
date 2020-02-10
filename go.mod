module github.com/filecoin-project/go-padreader

go 1.13

require (
	github.com/filecoin-project/filecoin-ffi v0.0.0-00010101000000-000000000000
	github.com/filecoin-project/specs-actors v0.0.0-20200210130641-2d1fbd8672cf
	github.com/stretchr/testify v1.4.0
)

replace github.com/filecoin-project/filecoin-ffi => ./extern/filecoin-ffi
