// Package ipvs provides IP Virtual Server (IPVS) management functionality
// for load balancing traffic across backend servers using Linux kernel features.
package ipvs

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"kube-proxy-ipvs/pkg/config"
)

// Handler manages IPVS configuration through ipvsadm commands
type Handler struct{}

// New creates a new IPVS handler and validates that ipvsadm is available
func New() (*Handler, error) {
	// Check if ipvsadm is available in PATH
	if _, err := exec.LookPath("ipvsadm"); err != nil {
		return nil, fmt.Errorf("ipvsadm not found: %v\nNote: ipvsadm is required and only available on Linux systems", err)
	}
	return &Handler{}, nil
}

// Close performs cleanup (no-op for command-based approach)
func (h *Handler) Close() {
	// No cleanup needed for command-based approach
}

// Apply configures IPVS with the provided configuration
func (h *Handler) Apply(cfg *config.Config) error {
	// Clear existing IPVS rules first to ensure clean state
	if err := h.clearRules(); err != nil {
		log.Printf("Warning: failed to clear existing rules: %v", err)
	}

	// Create virtual service with appropriate protocol (-t for TCP, -u for UDP)
	proto := strings.ToLower(cfg.Service.Protocol)
	cmd := exec.Command("ipvsadm", "-A", "-t", fmt.Sprintf("%s:%d", cfg.Service.VIP, cfg.Service.Port), "-s", cfg.Service.Scheduler)
	if proto == "udp" {
		cmd = exec.Command("ipvsadm", "-A", "-u", fmt.Sprintf("%s:%d", cfg.Service.VIP, cfg.Service.Port), "-s", cfg.Service.Scheduler)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add virtual service: %v", err)
	}

	// Add backend servers with weight=1 (equal distribution)
	for _, b := range cfg.Backends {
		cmd := exec.Command("ipvsadm", "-a", "-t", fmt.Sprintf("%s:%d", cfg.Service.VIP, cfg.Service.Port), "-r", fmt.Sprintf("%s:%d", b.IP, b.Port), "-w", "1")
		if proto == "udp" {
			cmd = exec.Command("ipvsadm", "-a", "-u", fmt.Sprintf("%s:%d", cfg.Service.VIP, cfg.Service.Port), "-r", fmt.Sprintf("%s:%d", b.IP, b.Port), "-w", "1")
		}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add backend %s:%d: %v", b.IP, b.Port, err)
		}
		log.Printf("Added backend %s:%d", b.IP, b.Port)
	}

	log.Printf("Created virtual server %s:%d [%s]", cfg.Service.VIP, cfg.Service.Port, cfg.Service.Scheduler)
	return nil
}

// clearRules removes all existing IPVS rules using ipvsadm -C
func (h *Handler) clearRules() error {
	cmd := exec.Command("ipvsadm", "-C")
	return cmd.Run()
}

// ShowStatus displays the current IPVS configuration using ipvsadm -Ln
func (h *Handler) ShowStatus() error {
	log.Println("Current IPVS configuration:")
	cmd := exec.Command("ipvsadm", "-Ln") // -L list, -n numeric output
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get IPVS status: %v", err)
	}
	log.Printf("\n%s", string(output))
	return nil
}
