// Package webserver
// Created by RTT.
// Author: teocci@yandex.com on 2022-Apr-26
package webserver

import (
	"embed"
	"fmt"
	"log"
	"mime"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/teocci/go-basic-web/src/config"
)

const (
	defaultProtocol = "http"
	formatAddress   = "%s:%d"
	formatURL       = "%s://%s/%s"
	defaultPage     = "page.html"
)

var (
	f       embed.FS
	address string
)

func Start() {
	address = fmt.Sprintf(formatAddress, GetLocalIp(), config.Data.Web.Port)
	gin.SetMode(gin.ReleaseMode)
	_ = mime.AddExtensionType(".js", "application/javascript")

	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*")

	router.StaticFS("/css", http.Dir("web/static/css"))
	router.StaticFS("/js", http.Dir("web/static/js"))
	router.StaticFS("/img", http.Dir("web/static/img"))
	router.StaticFile("/"+defaultPage, "web/static/"+defaultPage)
	router.StaticFile("/favicon.ico", "web/static/favicon.ico")

	router.Use(CORSMiddleware())

	fmt.Println("[url]", urlFormat(address))

	err := router.Run(address)
	if err != nil {
		log.Fatalln("Start HTTP Server error", err)
	}
}

func GetLocalIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "localhost"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func addressFormat(a string) string {
	s := strings.Split(a, ":")
	if s[0] == "" {
		s[0] = GetLocalIp()
	}
	return strings.Join(s[:], ":")
}

func urlFormat(a string) string {
	s := fmt.Sprintf(formatURL, defaultProtocol, a, defaultPage)

	return s
}
