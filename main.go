package main

import (
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/format/rtmp"
)

func main() {
	upstreamKey := "foo"
	upstreamURL := "rtmp://live-cdg.twitch.tv/app/" + upstreamKey

	server := &rtmp.Server{}

	server.HandlePublish = func(conn *rtmp.Conn) {
		proxyConn, _ := rtmp.Dial(upstreamURL)
		defer proxyConn.Close()

		avutil.CopyFile(proxyConn, conn)
	}

	server.ListenAndServe()
}
