package native

/*
#cgo CFLAGS: -I${SRCDIR}/libs/
#cgo windows,amd64 LDFLAGS: ${SRCDIR}/libs/libdeflate_windows_amd64.a
#cgo linux,amd64 LDFLAGS: ${SRCDIR}/libs/libdeflate_linux_amd64.a
#cgo darwin,amd64 LDFLAGS: ${SRCDIR}/libs/libdeflate_darwin_amd64.a
*/
import "C"
