package sstun

import "net"


func NewUDPAddr(host string, port int) net.UDPAddr {
   return net.UDPAddr {
      IP: net.ParseIP(host),
      Port: port,
      Zone: "",
   }
}
