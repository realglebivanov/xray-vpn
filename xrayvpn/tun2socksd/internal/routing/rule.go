package routing

import (
	"errors"
	"fmt"
	"log"
	"net"
	"syscall"

	"github.com/realglebivanov/hstd/hstdlib"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

const (
	localRulePriority        = 50
	transmissionRulePriority = 90
	xrayOutMarkRulePriority  = 100
)

func setUpRules() error {
	if err := netlink.RuleAdd(buildXrayOutMarkRule()); err != nil {
		return fmt.Errorf("set up xray out mark rule: %w", err)
	}

	uid, err := hstdlib.TransmissionUID()
	if err != nil {
		log.Printf("skip tranmission rule set up: %v", err)
		return nil
	}

	if err := netlink.RuleAdd(buildTransmissionLocalRule(uid)); err != nil {
		return fmt.Errorf("set up transmission uid local rule: %w", err)
	}

	if err := netlink.RuleAdd(buildTransmissionRule(uid)); err != nil {
		return fmt.Errorf("set up transmission uid rule: %w", err)
	}

	return nil
}

func tearDownRules() error {
	if err := netlink.RuleDel(buildXrayOutMarkRule()); err != nil && !errors.Is(err, syscall.ENOENT) {
		return fmt.Errorf("tear down xray out mark rule: %w", err)
	}

	uid, err := hstdlib.TransmissionUID()
	if err != nil {
		log.Printf("skip tranmission rule tear down: %v", err)
		return nil
	}

	if err := netlink.RuleDel(buildTransmissionLocalRule(uid)); err != nil && !errors.Is(err, syscall.ENOENT) {
		return fmt.Errorf("tear down transmission uid local rule: %w", err)
	}

	if err := netlink.RuleDel(buildTransmissionRule(uid)); err != nil && !errors.Is(err, syscall.ENOENT) {
		return fmt.Errorf("tear down transmission uid rule: %w", err)
	}

	return nil
}

func buildXrayOutMarkRule() *netlink.Rule {
	rule := netlink.NewRule()
	rule.Mark = hstdlib.XrayOutMark
	rule.Table = directRouteTable
	rule.Priority = xrayOutMarkRulePriority
	return rule
}

func buildTransmissionLocalRule(uid uint32) *netlink.Rule {
	rule := netlink.NewRule()
	rule.UIDRange = netlink.NewRuleUIDRange(uid, uid)
	rule.Dst = &net.IPNet{
		IP:   net.IPv4(192, 168, 2, 0),
		Mask: net.CIDRMask(24, 32),
	}
	rule.Table = unix.RT_TABLE_MAIN
	rule.Priority = localRulePriority
	return rule
}

func buildTransmissionRule(uid uint32) *netlink.Rule {
	rule := netlink.NewRule()
	rule.UIDRange = netlink.NewRuleUIDRange(uid, uid)
	rule.Table = directRouteTable
	rule.Priority = transmissionRulePriority
	return rule
}
