package dpx

// #cgo LDFLAGS: -ldpx
// #include <dpx.h>
// #include <stdlib.h>
import "C"

import (
	"net"
	"runtime"
	"strconv"
)

type Peer struct {
	peer *C.dpx_peer
}

func newPeer() *Peer {
	context := C.dpx_init()
	peer := C.dpx_peer_new(context)

	p := &Peer{
		peer: peer,
	}

	runtime.SetFinalizer(p, func(p *Peer) {
		C.dpx_peer_close(p)
		C.dpx_peer_free(p)
		C.dpx_cleanup(context)
	})

	return p
}

func (p *Peer) Open(method string) *Channel {
	cMethod := C.CString(method)
	defer C.free(cMethod)
	cChan := C.dpx_peer_open(p.peer, cMethod)
	return fromCChannel(cChan)
}

func (p *Peer) Accept() *Channel {
	cChan := C.dpx_peer_accept(p.peer)
	return fromCChannel(cChan)
}

func (p *Peer) Close() error {
	return ParseError(uint64(C.dpx_peer_close(p.peer)))
}

func (p *Peer) Connect(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	iport, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	chost := C.CString(host)
	defer C.free(host)

	return ParseError(C.dpx_peer_connect(p.peer, chost, C.int(iport)))
}

func (p *Peer) Bind(addr string) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	iport, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	chost := C.CString(host)
	defer C.free(host)

	return ParseError(C.dpx_peer_bind(p.peer, chost, C.int(iport)))
}
