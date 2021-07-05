// Simplified STUN or SSTUN for short is a simple protocol that enables a subet of the original
// stun protocol. In particular it enables discovery of the public facing ip and port (if the peer is behind a nat)
// and NAT hole punching at the same time. This is implemented as a udp client and server.
package sstun

import(
   "net"
   "fmt"
   "encoding/binary"
)

const (
   MAX_TOKEN_SIZE = 4294967296
)

type Server struct {
   ListenAddr net.UDPAddr
   conn *net.UDPConn
};

// Create a new server that listens on the given address.
func NewServer(addr net.UDPAddr) Server {
   return Server {
      ListenAddr: addr,
      conn: nil,
   }
}

// Starts listening on the specified address, blocking.
func (s *Server) Listen() error {
   // a 20 byte buffer for reading messages,
   // for most use-cases we don't actually need this 
   // because we are only interested in the connection info 
   var buff [20]byte

   conn, err := net.ListenUDP(s.ListenAddr.Network(), &s.ListenAddr )

   s.conn = conn;

   if err != nil {
      return err
   }

   for {
      // we read the message, this should not exceed four bytes, as the four bytes are simpy a tracking id set by the the client
      n, addr, err := conn.ReadFromUDP(buff[:])

      // TODO if n < 4 we should keep reading
      if err != nil || n != 4 {
         // ignore error
         continue
      }

      fmt.Printf("Received a message from %s\n", addr.IP)

      // Let us now send the received information back
      s.reply(buff[:], addr)
   }

   return nil
}

// Write the given uint32 in little endian to the given byte slice
func writeUint32(n uint32, buff []byte) {
   binary.LittleEndian.PutUint32(buff, n)
}


// Formulates and sends a reply to the server
func (s *Server) reply(token []byte, addr *net.UDPAddr) {
   // a packet consists of the received token (4 bytes) followed by the public facing IP 
   // and finally followed by the public port of the peer
   packet := make([]byte, 12)
   copy(packet[0:4], token)

   // TODO: support ipv6
   ip4 := addr.IP.To4()
   if ip4 == nil {
      return
   } else {
      copy(packet[4:8], ip4)
      writeUint32(uint32(addr.Port), packet[8:12])
   }

   // finally send it
   s.conn.WriteTo(packet, addr)
}
