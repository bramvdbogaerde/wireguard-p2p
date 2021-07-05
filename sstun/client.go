package sstun

import(
   "net"
   "errors"
   "encoding/binary"
)

type Client struct {
   // The address of the server to connect to
   ServerAddr *net.UDPAddr
   // Keeps a reference to the "established" connection with the server.
   conn *net.UDPConn
   // Keeps track of the current token
   token int
};

type ClientInfo struct {
   IP net.IP
   Port int
}

// Creates a new client based on the address of the server to connect to
func NewClient(addr *net.UDPAddr) Client {
   return Client {
      ServerAddr: addr,
      conn:       nil,
      token:      0,
   }
}

// Generates a unique token. As the token are generated rather infrequently, this is simply a rolling 32-bit counter, that gets incremented at each call to newToken.
func (c *Client) newToken() int {
   c.token = (c.token + 1) % MAX_TOKEN_SIZE
   return c.token
}

// Creates a UDP socket for sending messages to the server
func (c *Client) establishConn() error {
   conn, err := net.DialUDP(c.ServerAddr.Network(), nil, c.ServerAddr)
   c.conn = conn
   return err
}

// Send the 4 byte token to the server (encoded as little endian)
func (c *Client) sendToken(token int) error {
   packet := make([]byte, 4)
   writeUint32(uint32(token), packet)
   // TODO: keep writing if n < 4
   _, err := c.conn.Write(packet)
   if err != nil {
      return err
   }

   return nil
}


// Send the query to the server
func (c *Client) sendQuery() (int, error) {
   token := c.newToken()
   err := c.establishConn()

   if err != nil {
      return 0, err
   }

   err = c.sendToken(token)

   if err != nil {
      return 0, err
   }

   return token, nil
}

func (c *Client) waitForAnswer(token int) (*ClientInfo, error) {
   // We are expecting only twelve bytes
   buff := make([]byte, 12)
   // TODO: if n < 12 then keep reading
   _, err := c.conn.Read(buff)
   if err != nil {
      return nil, err
   }

   // the received token must match the sent one
   rtoken := binary.LittleEndian.Uint32(buff[0:4])
   if rtoken != uint32(token) {
      return nil, errors.New("Received token was not the sent token")
   }

   port  := binary.LittleEndian.Uint32(buff[8:12])

   return &ClientInfo {
      IP: buff[4:8],
      Port: int(port),
   }, nil
}

// Ask the public ip of the client and public facing port
func (c *Client) Ask() (*ClientInfo, error) {
   // The protocol is as follows:
   // 1. the client sends a packet containing a four byte token to the server
   // 2. the server replies with the IP and port that it observes
   // 3. the server sends a reply to that port and IP with that same token, ip and port (little endian)
   // 4. the client waits until it receives the packet corresponding to its token and parses it

   token, err := c.sendQuery()
   if err != nil {
      return nil, err
   }

   return c.waitForAnswer(token)
}
