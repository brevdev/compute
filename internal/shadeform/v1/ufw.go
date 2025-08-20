package v1

import (
	"encoding/base64"
	"fmt"
	v1 "github.com/brevdev/cloud/pkg/v1"
)

const (
	ufwForceReset           = "ufw --force reset"
	ufwDefaultDropIncoming  = "ufw default deny incoming"
	ufwDefaultAllowOutgoing = "ufw default allow outgoing"
	ufwDefaultAllowPort22   = "ufw allow 22/tcp"
	ufwDefaultAllowPort2222 = "ufw allow 2222/tcp"
	ufwForceEnable          = "ufw --force enable"
)

func (c *ShadeformClient) generateFirewallScript(firewallRules v1.FirewallRules) (string, error) {
	commands := []string{ufwForceReset, ufwDefaultDropIncoming, ufwDefaultAllowOutgoing, ufwDefaultAllowPort22, ufwDefaultAllowPort2222}

	for _, firewallRule := range firewallRules.IngressRules {
		commands = append(commands, c.convertIngressFirewallRuleToUfwCommand(firewallRule)...)
	}

	for _, firewallRule := range firewallRules.EgressRules {
		commands = append(commands, c.convertEgressFirewallRuleToUfwCommand(firewallRule)...)
	}

	// Add the enable command
	commands = append(commands, ufwForceEnable)

	script := ""
	for _, command := range commands {
		script = script + fmt.Sprintf("%v\n", command)
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(script))
	return encoded, nil
}

func (c *ShadeformClient) convertIngressFirewallRuleToUfwCommand(firewallRule v1.FirewallRule) []string {
	cmds := []string{}
	portSpec := ""
	if firewallRule.FromPort == firewallRule.ToPort {
		portSpec = fmt.Sprintf("port %d", firewallRule.FromPort)
	} else {
		portSpec = fmt.Sprintf("port %d:%d", firewallRule.FromPort, firewallRule.ToPort)
	}

	if len(firewallRule.IPRanges) == 0 {
		cmds = append(cmds, fmt.Sprintf("ufw allow in from any to any port %s", portSpec))
	}

	for _, ipRange := range firewallRule.IPRanges {
		cmds = append(cmds, fmt.Sprintf("ufw allow in from %s to any port %s", ipRange, portSpec))
	}
	return cmds
}

func (c *ShadeformClient) convertEgressFirewallRuleToUfwCommand(firewallRule v1.FirewallRule) []string {
	cmds := []string{}
	portSpec := ""
	if firewallRule.FromPort == firewallRule.ToPort {
		portSpec = fmt.Sprintf("port %d", firewallRule.FromPort)
	} else {
		portSpec = fmt.Sprintf("port %d:%d", firewallRule.FromPort, firewallRule.ToPort)
	}

	if len(firewallRule.IPRanges) == 0 {
		cmds = append(cmds, fmt.Sprintf("ufw allow out to any port %s", portSpec))
	}

	for _, ipRange := range firewallRule.IPRanges {
		cmds = append(cmds, fmt.Sprintf("ufw allow out to %s port %s", ipRange, portSpec))
	}
	return cmds
}
