package main

import (
	"os"

	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/format/rtmp"
)

type TwitchService struct{}

func (TwitchService) Stream(key string) string {
	return "rtmp://live-cdg.twitch.tv/app/" + key
}

func main() {
	twitch := TwitchService{}
	proxyServer := &RTMPProxy{
		Dsts: []string{
			twitch.Stream(os.Getenv("TWITCH_KEY")),
			twitch.Stream(os.Getenv("TWITCH_KEY2")),
		},
	}

	server := &rtmp.Server{
		Addr: ":" + os.Getenv("PORT"),
	}
	server.HandlePublish = func(conn *rtmp.Conn) {
		proxyServer.ServeRTMP(conn)
	}

	server.ListenAndServe()
}

type RTMPProxy struct {
	Dsts []string
}

func (s RTMPProxy) ServeRTMP(conn *rtmp.Conn) {
	m := &muxer{}

	for _, dst := range s.Dsts {
		m.Open(dst)
	}
	defer m.Close()

	avutil.CopyFile(m, conn)
}

type muxer struct {
	conns []*rtmp.Conn
}

func (m *muxer) Open(url string) error {
	c, err := rtmp.Dial(url)
	if err != nil {
		return err
	}

	m.conns = append(m.conns, c)
	return nil
}

func (m *muxer) Close() {
	for _, c := range m.conns {
		c.Close()
	}
}

func (m *muxer) WritePacket(p av.Packet) error {
	for _, c := range m.conns {
		c.WritePacket(p)
	}

	return nil
}

func (m *muxer) WriteHeader(streams []av.CodecData) error {
	for _, c := range m.conns {
		c.WriteHeader(streams)
	}

	return nil
}

func (m *muxer) WriteTrailer() error {
	for _, c := range m.conns {
		c.WriteTrailer()
	}

	return nil
}
