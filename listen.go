package main

import (
	"fmt"
	"net"
	"net/http"
)

type KeepAliveListener struct {
	net.Listener
}

func (l *KeepAliveListener) Accept() (c net.Conn, err error) {
	c, err = l.Listener.Accept()
	if err != nil {
		Log.Printf("Listener.Accept() failed (%s)", err.Error())
		return
	}

	// set keepalive
	if tc, ok := c.(*net.TCPConn); !ok {
		panic("Assection type failed c.(net.TCPConn)")
	} else {
		err = tc.SetKeepAlive(true)
		if err != nil {
			Log.Printf("tc.SetKeepAlive(true) failed (%s)", err.Error())
			return
		}
	}

	return
}

func Listen() error {
	addr := fmt.Sprintf("%s:%d", Conf.Addr, Conf.Port)
	if Conf.TCPKeepAlive == 1 {
		server := &http.Server{}
		l, err := net.Listen("tcp", addr)
		if err != nil {
			Log.Printf("net.Listen(\"tcp\", \"%s\") failed (%s)", addr, err.Error())
			return err
		}

		return server.Serve(&KeepAliveListener{Listener: l})
	} else {
		if err := http.ListenAndServe(addr, nil); err != nil {
			Log.Printf("http.ListenAdServe(\"%s\") failed (%s)", addr, err.Error())
			return err
		}
	}
	// nerve here
	return nil
}
