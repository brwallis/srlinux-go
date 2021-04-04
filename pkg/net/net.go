package net

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	log "k8s.io/klog"
)

const fileMode0755 = os.FileMode(0755)

// IPInfo contains a v4/v6 IP address and mask
type IPInfo struct {
	IPCIDR    string
	IPAddress net.IP
	IPMask    net.IPMask
	IPNetwork *net.IPNet
}

// Interface contains a single interface
type Interface struct {
	Name       string
	UdevName   string
	PCIAddress string
	MTU        int
	HWAddr     string
	Flags      string
	Driver     string
	IPInfo     []IPInfo
}

// Namespace contains a single Linux namespace
type Namespace struct {
	NS        int    `json:"ns,omitempty"`
	Type      string `json:"type,omitempty"`
	Processes int    `json:"nprocs,omitempty"`
	PID       int    `json:"pid,omitempty"`
	User      string `json:"user,omitempty"`
	NetNSID   string `json:"netnsid,omitempty"`
	NSFS      string `json:"nsfs,omitempty"`
	Command   string `json:"command,omitempty"`
}

type Namespaces struct {
	Namespaces []Namespace `json:"namespaces,omitempty"`
}

// ShellOut executes the provides command in a bash shell, and returns the stdout and stderr, along with rc of the command
func ShellOut(command string) (error, string, string) {
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	return err, stdoutBuf.String(), stderrBuf.String()
}

// ListNamespaces lists all namespaces in Linux and returns them in a list
func ListNamespaces() map[int]Namespace {
	var namespaces Namespaces
	namespaceList := make(map[int]Namespace)

	err, out, errOut := ShellOut("lsns -t net -J")
	if err != nil {
		log.Fatalf("Unable to list namespaces in host: %s, stderr: %s", err, errOut)
	}
	jsonOut := strings.Replace(out, "\n", "", -1)
	log.Infof("Output from lsns: %s", jsonOut)

	err = json.Unmarshal([]byte(jsonOut), &namespaces)
	if err != nil {
		log.Fatalf("Unable to unmarshal lsns output to Namespaces struct: output %s, error: %s", jsonOut, err)
	}

	for _, ns := range namespaces.Namespaces {
		namespaceList[ns.NS] = ns
		// namespaceList = append(namespaceList, ns)
	}
	return namespaceList
}

// GetNamespace returns a Namespace object from an inode
func GetNamespace(inode int) Namespace {
	var ns Namespace

	namespaces := ListNamespaces()
	if val, ok := namespaces[inode]; ok {
		ns = val
	} else {
		log.Infof("Unable to find namespace with inode: %s", inode)
	}

	return ns
}

// WriteToFile opens a file in write only mode and writes bytes to it
func WriteToFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY, fileMode0755)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

// BindVFIO takes a PCI address, and a PCI path, and then unbinds the current driver (if needed), before binding vfio-pci
func BindVFIO(pciAddress string) error {
	var err error

	// Get the path to the PCI address
	pciPath := GetPCIPath(pciAddress)
	// Unbind if needed
	unbindPath := fmt.Sprintf("%s/driver/unbind", pciPath)
	if err := WriteToFile(unbindPath, []byte(pciAddress)); err != nil {
		return err
	}
	//file, err := os.OpenFile(unbindPath, os.O_APPEND, 0200)
	//if err != nil {
	//	datawriter := bufio.NewWriter(file)
	//	_, _ = datawriter.WriteString(pciAddress + "\n")
	//	datawriter.Flush()
	//	file.Close()
	//} else {
	//	log.Infof("Unable to open driver unbind, probably already unbound. Error: %e", err)
	//}

	// Bind
	overridePath := fmt.Sprintf("%s/driver_override", pciPath)
	if err := WriteToFile(overridePath, []byte("vfio-pci")); err != nil {
		return err
	}
	//overridePath := fmt.Sprintf("%s/driver_override", pciPath)
	//file, err = os.OpenFile(overridePath, os.O_APPEND, 0200)
	//if err != nil {
	//	datawriter := bufio.NewWriter(file)
	//	_, _ = datawriter.WriteString("vfio-pci" + "\n")
	//	datawriter.Flush()
	//	file.Close()
	//} else {
	//	log.Fatalf("Unable to open driver override path: %s. Error: %e", overridePath, err)
	//}

	// Probe driver
	probePath := "/sys/bus/pci/drivers_probe"
	if err := WriteToFile(probePath, []byte(pciAddress)); err != nil {
		return err
	}
	//probePath := "/sys/bus/pci/drivers_probe"
	//file, err = os.OpenFile(probePath, os.O_APPEND, 0200)
	//if err != nil {
	//	datawriter := bufio.NewWriter(file)
	//	_, _ = datawriter.WriteString(pciAddress + "\n")
	//	datawriter.Flush()
	//	file.Close()
	//} else {
	//	log.Fatalf("Unable to open PCI probe interface: %s. Error: %e", probePath, err)
	//}

	return err
}

// NameToPCI takes a NIC name in udevadm path format, and returns the PCI address that it references
func NameToPCI(name string) (string, error) {
	var pciAddress string
	var err error

	if strings.Contains(name, "enp") {
		re := regexp.MustCompile(`enp(?P<bus>[0-9]+)s(?P<device>[0-9]+)f(?P<function>[0-9]+)`)
		matches := re.FindStringSubmatch(name)
		busInt, err := strconv.Atoi(matches[1])
		if err != nil {
			return "", err
		}
		deviceInt, err := strconv.Atoi(matches[2])
		if err != nil {
			return "", err
		}
		functionInt, err := strconv.Atoi(matches[3])
		if err != nil {
			return "", err
		}
		pciAddress = fmt.Sprintf(
			"0000:%02s:%02s.%s",
			strconv.FormatInt(int64(busInt), 16),
			strconv.FormatInt(int64(deviceInt), 16),
			strconv.FormatInt(int64(functionInt), 16),
		)
	} else {
		err = errors.New((fmt.Sprintf("Interface in incorrect format; received %s, expected enp*s*f*", name)))
	}
	return pciAddress, err
}

// PCIToName returns the NIC name of a PCI address in the udevadm PATH format
func PCIToName(pciAddress string) (string, error) {
	var nicName string
	var err error
	// Break up the NIC name into its domain:bus:device.function
	// Values are initially provided in hex, need to convert to base10
	// 0000:12:00.1
	// ens2f1
	// enp18s0f1

	parts := strings.Split(pciAddress, ":")
	bus := parts[1]
	deviceParts := strings.Split(parts[2], ".")
	device := deviceParts[0]
	function := deviceParts[1]

	// Process bus
	busInt, err := strconv.ParseInt(bus, 16, 64)
	if err != nil {
		log.Errorf("Unable to parse bus for PCI address: %s, bus %s", pciAddress, bus)
	}

	// Process device
	deviceInt, err := strconv.ParseInt(device, 16, 64)
	if err != nil {
		log.Errorf("Unable to parse device for PCI address: %s, bus %s", pciAddress, device)
	}

	// Process function
	functionInt, err := strconv.ParseInt(function, 16, 64)
	if err != nil {
		log.Errorf("Unable to parse device for PCI address: %s, bus %s", pciAddress, device)
	}

	// Build NIC name
	nicName = fmt.Sprintf("enp%ds%df%d", busInt, deviceInt, functionInt)

	return nicName, err
}

// GetIPInfo returns the IPv4 addresses for an interface
func GetIPInfo(netInfo *net.Interface) ([]IPInfo, error) {
	var ipInfo []IPInfo
	var err error
	addresses, err := netInfo.Addrs()
	log.Infof("Found addresses for device: %s; %v", netInfo.Name, addresses)
	if err != nil {
		log.Errorf("Unable to parse IP info for device %s: %e", addresses, err)
	}
	// Parse addresses and remove ones we don't care about
	for _, prefix := range addresses {
		log.Infof("Processing address: %v", prefix)
		ip, network, err := net.ParseCIDR(prefix.String())
		if err != nil {
			log.Errorf("Unable to ParseCIDR on prefix: %s; %e", prefix, err)
		}
		if isGlobalUnicast := ip.IsGlobalUnicast(); isGlobalUnicast {
			ipInfo = append(ipInfo, IPInfo{
				IPCIDR:    prefix.String(),
				IPAddress: ip,
				IPMask:    network.Mask,
				IPNetwork: network,
			})
		} else {
			log.Infof("Not processing non-Global Unicast address: %s", ip)
		}
	}
	return ipInfo, err
}

// GetPCIPath returns the full real path to a PCI device
func GetPCIPath(pciAddress string) string {
	var pciPath string
	var pciPathReal string
	pciPath = fmt.Sprintf("/sys/bus/pci/devices/%s", pciAddress)
	pciPathReal, err := os.Readlink(pciPath)
	if err != nil {
		log.Errorf("Unable to parse PCI path for PCI address: %s: %e", pciPath, err)
	}
	return pciPathReal
}

// GetPCIAddress returns the PCI address for an interface
func GetPCIAddress(nicName string) (string, error) {
	var pciAddress string
	var err error
	devicePath := fmt.Sprintf("/sys/class/net/%s/device", nicName)
	devicePathReal, err := os.Readlink(devicePath)
	if err != nil {
		return "", err
		//log.Errorf("Unable to parse PCI address for device %s: %e", nicName, err)
	}
	// Parse the address
	words := strings.Split(devicePathReal, "/")
	pciAddress = words[len(words)-1]
	log.Infof("Found address: %s for device: %s", pciAddress, nicName)
	return pciAddress, err
}

// GetCurrentDriver returns the current driver for an interface
func GetCurrentDriver(nicName string) (string, error) {
	var driver string
	var err error
	pciAddress, err := NameToPCI(nicName)
	if err != nil {
		log.Errorf("Unable to parse PCI address for device %s: %e", nicName, err)
		return "", err
	}
	devicePath := fmt.Sprintf("/sys/bus/pci/devices/%s", pciAddress)
	driverPath := fmt.Sprintf("%s/driver", devicePath)
	driverPathReal, err := os.Readlink(driverPath)
	if err != nil {
		log.Errorf("Unable to parse driver for device %s: %e", nicName, err)
		return "", err
	}
	// Parse the driver
	words := strings.Split(driverPathReal, "/")
	driver = words[len(words)-1]
	log.Infof("Found driver: %s for device: %s", driver, nicName)
	return driver, err
}

// ProcessInterface takes in a PCI address and returns an Interface{} struct populated with interface information
func ProcessInterface(pciAddress string) (Interface, error) {
	var currentInterface Interface
	var netInfoFound bool
	var ipInfo []IPInfo
	log.Infof("Processing PCI device: %v", pciAddress)
	// Retrieve NIC udev name
	name, err := PCIToName(pciAddress)
	if err != nil {
		log.Infof("Unable to parse name for PCI device: %s", pciAddress)
		return currentInterface, err
	}
	// Retrieve netInfo
	netInfo, err := net.InterfaceByName(name)
	if err != nil {
		log.Infof("Unable to parse netInfo for PCI device")
		netInfoFound = false
		//return currentInterface, err
	} else {
		netInfoFound = true
	}
	// Retrieve current driver
	driver, err := GetCurrentDriver(name)
	if err != nil {
		log.Infof("Unable to parse driver for interface: %s", name)
		// Continue processing this interface
	}
	currentInterface = Interface{
		Name:       name,
		UdevName:   name,
		PCIAddress: pciAddress,
		Driver:     driver,
	}
	if netInfoFound {
		// Retrieve IPv4 info
		ipInfo, err = GetIPInfo(netInfo)
		if err != nil {
			log.Infof("Unable to parse IPv4 addresses for interface: %s", name)
			// Continue processing this interface
		} else {
			currentInterface.IPInfo = ipInfo
		}
		currentInterface.MTU = netInfo.MTU
		currentInterface.HWAddr = netInfo.HardwareAddr.String()
		currentInterface.Flags = netInfo.Flags.String()
	}
	err = nil
	return currentInterface, err
}

// ListInterfaces returns a list of interfaces in the system
func ListInterfaces() []Interface {
	//func main() {
	var interfaces []Interface
	var currentInterface Interface
	var err error

	log.Info("Listing interfaces...")
	items, _ := ioutil.ReadDir("/sys/class/net")
	// Process local interfaces bound to net drivers
	for _, item := range items {
		name := item.Name()
		if name == "eth0" || name == "lo" {
			// We don't want to process these interfaces
			continue
		} else {
			// Retrieve PCI address
			pciAddress, err := GetPCIAddress(name)
			if err != nil {
				// If it doesn't have a PCI address, we can ignore it
				continue
			}
			log.Infof("Processing interface: %v", name)
			currentInterface, err = ProcessInterface(pciAddress)
			interfaces = append(interfaces, currentInterface)
		}
	}

	// Pull env vars from file written to disk
	envVars := make(map[string]string)
	f, _ := os.Open("/etc/opt/srlinux/env")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		env := strings.SplitN(scanner.Text(), "=", 2)
		envKey := env[0]
		envVal := env[1]
		envVars[envKey] = envVal
		if strings.Contains(envKey, "PCIDEVICE") {
			// Process the key
			log.Infof("Found passed through interface: %s, val: %s", envKey, envVal)
			currentInterface, err = ProcessInterface(envVal)
			if err != nil {
				log.Infof("Unable to parse Kubernetes-passed device type: %s, address: %e, error: %e", envKey, envVal, err)
				continue
			}
			interfaces = append(interfaces, currentInterface)
		}
	}

	// Process local interfaces that are passed through to us (if we are in the vxdp container)
	for _, env := range os.Environ() {
		envPair := strings.SplitN(env, "=", 2)
		key := envPair[0]
		value := envPair[1]
		log.Infof("Processing env var: %s, val: %s", key, value)

		if strings.Contains(key, "PCIDEVICE") {
			// Process the key
			log.Infof("Found passed through interface: %s, val: %s", key, value)
			currentInterface, err = ProcessInterface(value)
			if err != nil {
				log.Infof("Unable to parse Kubernetes-passed device type: %s, address: %e, error: %e", key, value, err)
				continue
			}
			interfaces = append(interfaces, currentInterface)
		}
	}

	log.Infof("Interfaces discovered: %#v", interfaces)
	return interfaces
}
