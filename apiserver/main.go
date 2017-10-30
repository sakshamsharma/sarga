package apiserver

import (
	"fmt"
	"math/rand"
	"time"

	arg "github.com/alexflint/go-arg"
	"github.com/sakshamsharma/sarga/common/iface"
	"github.com/sakshamsharma/sarga/impl/httpnet"
	"github.com/sakshamsharma/sarga/impl/sdht"
)

type ServerArgs struct {
	iface.CommonArgs

	Seeds          []string
	RandomDHTCount int
}

func Init() error {
	var args ServerArgs
	arg.MustParse(&args)

	if args.Port == 0 {
		return fmt.Errorf("port not provided. Please provide a port using --port=<integer>")
	}

	if args.IP == "" {
		args.IP = "127.0.0.1"
	}

	for _, seed := range args.Seeds {
		fmt.Println(seed)
	}

	seeds, err := iface.ParseAddresses(args.Seeds)
	if err != nil {
		return err
	}

	dhtInst := &sdht.SDHT{}
	if err = dhtInst.Init(iface.Address{"0.0.0.0", 8080},
		seeds, &httpnet.HTTPNet{}); err != nil {
		return err
	}

	if args.RandomDHTCount > 0 {
		time.Sleep(2 * time.Second)
		for i := 0; i < args.RandomDHTCount; i++ {
			nodeDHT := &sdht.SDHT{}
			addr := iface.Address{"0.0.0.0", rand.Intn(3000) + 4000}
			nodeDHT.Init(addr, []iface.Address{{"0.0.0.0", 8080}}, &httpnet.HTTPNet{})
		}
		time.Sleep(2 * time.Second)
	}

	StartAPIServer(args.CommonArgs, dhtInst)

	return nil
}
