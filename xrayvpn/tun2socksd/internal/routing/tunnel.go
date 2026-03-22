package routing

import (
	"fmt"
	"log"
	"time"

	"github.com/realglebivanov/hstd/hstdlib"
	"github.com/realglebivanov/hstd/tun2socksd/internal/routing/state"
	"github.com/vishvananda/netlink"
)

type Tunnel struct {
	Gw      *state.DefaultGateway
	TunLink netlink.Link
	TunAddr *netlink.Addr
}

var (
	TunDev = hstdlib.EnvOr("TUN_DEV", "xray0")
	ApdDev = hstdlib.EnvOr("APD_DEV", "wlp4s0")
	WanDev = hstdlib.EnvOr("WAN_DEV", "eno1")
	TunMTU = hstdlib.EnvOrInt("TUN_MTU", 1500)
)

func TearDownTunnel(tun *Tunnel) error {
	if err := tearDownNftables(); err != nil {
		return fmt.Errorf("tear down nftables: %w", err)
	}
	if err := tearDownRoutes(tun); err != nil {
		return fmt.Errorf("tear down route table: %w", err)
	}
	if err := tearDownRules(); err != nil {
		return fmt.Errorf("clean tear down route rules: %w", err)
	}
	return nil
}

func SetUpTunnel() (*Tunnel, error) {
	tunnel, err := setUpLink()
	if err != nil {
		return tunnel, err
	}

	if err := setUpRules(); err != nil {
		return tunnel, fmt.Errorf("set up route rules: %w", err)
	}
	if err := setUpRoutes(tunnel); err != nil {
		return tunnel, fmt.Errorf("set up route table: %w", err)
	}
	if err := setUpNftables(); err != nil {
		return tunnel, fmt.Errorf("set up nftables: %w", err)
	}
	return tunnel, nil
}

func setUpLink() (*Tunnel, error) {
	gw, err := preserveDefaultGateway()
	if err != nil {
		return nil, err
	}

	tunLink, err := awaitTunLink(10 * time.Second)
	if err != nil {
		return nil, err
	}

	tunAddr, err := configureTun(tunLink)
	if err != nil {
		return nil, err
	}

	return &Tunnel{Gw: gw, TunLink: tunLink, TunAddr: tunAddr}, nil
}

func awaitTunLink(timeout time.Duration) (netlink.Link, error) {
	log.Printf("waiting for %s ...", TunDev)

	deadline := time.Now().Add(timeout)
	for {
		if link, err := netlink.LinkByName(TunDev); err == nil {
			return link, nil
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for %s", TunDev)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func configureTun(link netlink.Link) (*netlink.Addr, error) {
	addr, err := netlink.ParseAddr("198.18.0.1/16")
	if err != nil {
		return nil, err
	}
	if err := netlink.AddrReplace(link, addr); err != nil {
		return nil, fmt.Errorf("addr add on %s: %w", TunDev, err)
	}
	if err := netlink.LinkSetUp(link); err != nil {
		return nil, fmt.Errorf("link set up %s: %w", TunDev, err)
	}
	return addr, nil
}
