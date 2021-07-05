package main


import(
   "github.com/urfave/cli/v2"
   "vdb.space/wireguard-p2p/sstun"
   "fmt"
)

func RunServer(context *cli.Context) error {
   addr := sstun.NewUDPAddr("0.0.0.0", 9090)
   server := sstun.NewServer(addr)
   server.Listen()
   return nil
}

func RunClient(context *cli.Context) error {
   // TODO: read remote address from config or from flags
   addr := sstun.NewUDPAddr("10.0.0.75", 9090)
   client := sstun.NewClient(&addr)
   info, err := client.Ask()
   if err != nil {
      return err
   }

   fmt.Printf("Got client info %s", info);
   return nil
}

func CreateApp() *cli.App {
   app := &cli.App {
      Commands: []*cli.Command {
         {
            Name: "server",
            Usage: "run the sstun server in the foreground",
            Action: RunServer,
         },
         {
            Name: "client",
            Usage: "run the sstun demo client",
            Action: RunClient,
         },
      },
   }

   return app
}
