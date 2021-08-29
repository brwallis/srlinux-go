package srlinux

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/containernetworking/cni/pkg/skel"

	"github.com/intel/userspace-cni-network-plugin/logging"
	"github.com/intel/userspace-cni-network-plugin/pkg/types"

	gpb "github.com/openconfig/gnmi/proto/gnmi"

	"github.com/brwallis/srlinux-go/pkg/gnmi"

	log "k8s.io/klog"
)

// IPPrefix is a prefix
type IPPrefix struct {
	Prefix string `json:"ip-prefix"`
}

// SubInterfaceIP is an instance of a subinterface
type SubInterfaceIP struct {
	Addresses []IPPrefix `json:"address"`
}

// SubInterfaceShort is an instance of a subinterface
type SubInterfaceShort struct {
	Name string `json:"name"`
}

// SubInterface is an instance of a subinterface
type SubInterface struct {
	Index      int    `json:"index"`
	Type       string `json:"type,omitempty"`
	AdminState string `json:"admin-state,omitempty"`
}

// Interface is an instance of an interface
type Interface struct {
	Name          string         `json:"name,omitempty"`
	SubInterfaces []SubInterface `json:"subinterface,omitempty"`
	AdminState    string         `json:"admin-state,omitempty"`
}

// NetworkInstance is a instaince of a network instance (heh)
type NetworkInstance struct {
	Name       string              `json:"name,omitempty"`
	Type       string              `json:"type,omitempty"`
	Interfaces []SubInterfaceShort `json:"interface,omitempty"`
	AdminState string              `json:"admin-state,omitempty"`
	RouterID   string              `json:"router-id,omitempty"`
}

// YangNetworkInstance is the top level structure of network instance YANG
type YangNetworkInstance struct {
	NetworkInstance []NetworkInstance `json:"srl_nokia-network-instance:network-instance"`
}

// YangInterface is the top level structure of interface YANG
type YangInterface struct {
	Interface []Interface `json:"srl_nokia-interfaces:interface"`
}

// GetNetworkInstance gets a network instance
func GetNetworkInstance(netInstName string) bool {
	var found bool
	// Run a get
	gnmi.Get(fmt.Sprintf("/network-instance[name=%s]", netInstName))
	return found
}

// AddNetworkInstance creates a network instance
func AddNetworkInstance(name string, netInstType string) error {
	var err error
	var netInstanceList []NetworkInstance

	newNetworkInstance := NetworkInstance{
		Name:       name,
		Type:       netInstType,
		AdminState: "enable",
	}
	netInstanceList = append(netInstanceList, newNetworkInstance)
	// Create the network instance
	JSONNetworkInstance := &YangNetworkInstance{
		NetworkInstance: netInstanceList,
	}
	out, err := json.MarshalIndent(JSONNetworkInstance, "", "  ")
	if err != nil {
		log.Infof("Unable to Marshal network instance: %e", err)
	}
	gNMINetworkInstance := &gpb.TypedValue{
		Value: &gpb.TypedValue_JsonIetfVal{
			JsonIetfVal: out,
		},
	}
	gNMINetworkInstancePath := "/"
	log.Infof("AddNetworkInstance: Adding instance; name: %s, type: %s, path: %s", name, netInstType, gNMINetworkInstancePath)
	_, err = gnmi.Update(context.Background(), gNMINetworkInstancePath, gNMINetworkInstance)
	if err != nil {
		log.Errorf("gNMI set failed to create network-instance: %s, err: %e", name, err)
	}
	return err
}

// DeleteNetworkInstance deletes a network instance
func DeleteNetworkInstance(name string) error {
	var err error
	return err
}

// DeleteIfEmptyNetworkInstance deletes a network instance if it is empty of subinterfaces
func DeleteIfEmptyNetworkInstance(netInstName string) error {
	var err error
	// Check if the network instance has any subinterfaces present
	//path := fmt.Sprintf("/network-instance[name=%s]/interface[name=*]", netInstName)
	//resp, err := gnmi.Get(path)

	// If it has subinterfaces, return no error

	// If it doesn't have subinterfaces, call DeleteNetworkInstance

	// Return any errors
	return err
}

// AddInterface creates an interface
func AddInterface(name string, intType string, subIntType string, netInstName string) error {
	var err error
	var subInterfaceList []SubInterface
	var newSubInterface SubInterface
	var interfaceList []Interface
	var newInterface Interface
	var subInterfaceShortList []SubInterfaceShort
	var newSubInterfaceShort SubInterfaceShort

	// Create the subinterface
	newSubInterface = SubInterface{
		Index:      0,
		AdminState: "enable",
		Type:       fmt.Sprintf("srl_nokia-interfaces:%s", subIntType),
	}
	subInterfaceList = append(subInterfaceList, newSubInterface)
	// Create the interface
	newInterface = Interface{
		Name:          name,
		AdminState:    "enable",
		SubInterfaces: subInterfaceList,
	}
	interfaceList = append(interfaceList, newInterface)

	JSONInterface := &YangInterface{
		Interface: interfaceList,
	}
	out, err := json.MarshalIndent(JSONInterface, "", "  ")
	if err != nil {
		log.Errorf("Unable to Marshal interface: %e", err)
		return err
	}
	gNMIInterface := &gpb.TypedValue{
		Value: &gpb.TypedValue_JsonIetfVal{
			JsonIetfVal: out,
		},
	}
	gNMIInterfacePath := "/"
	log.Infof("AddInterface: Adding interface; name: %s, type: %s, network-instance: %s", name, intType, netInstName)
	_, err = gnmi.Update(context.Background(), gNMIInterfacePath, gNMIInterface)
	if err != nil {
		log.Errorf("gNMI set failed for interface; name: %s, type: %s: %e", name, intType, err)
		return err
	}
	// Add subinterface to the network instance
	shortSubIntName := newInterface.Name + "." + strconv.Itoa(newSubInterface.Index)
	newSubInterfaceShort = SubInterfaceShort{
		Name: shortSubIntName,
	}
	subInterfaceShortList = append(subInterfaceShortList, newSubInterfaceShort)
	JSONNetInstance := NetworkInstance{
		//Name:       netInstName,
		Interfaces: subInterfaceShortList,
	}
	out, err = json.MarshalIndent(JSONNetInstance, "", "  ")
	if err != nil {
		log.Errorf("Unable to Marshal network instance: %e", err)
		return err
	}
	gNMINetInstance := &gpb.TypedValue{
		Value: &gpb.TypedValue_JsonIetfVal{
			JsonIetfVal: out,
		},
	}
	gNMINetInstancePath := fmt.Sprintf("/network-instance[name=%s]", netInstName)
	log.Infof("AddSubInterface: Adding subinterfaces to network-instance %s", netInstName)
	//gnmi.Set(gNMINetInstancePath, gNMINetInstance)
	_, err = gnmi.Update(context.Background(), gNMINetInstancePath, gNMINetInstance)
	return err
}

// DeleteInterface deletes an interface
func DeleteInterface(name string, containerNameShort string, netInstName string) error {
	var err error
	log.Infof("DeleteInterface: Deleting interface; name: %s, network-instance: %s", name, netInstName)
	// Delete the subinterface
	path := fmt.Sprintf("/network-instance[name=%s]/interface[name=%s]", netInstName, name+".0")
	_, err = gnmi.Delete(context.Background(), path)
	if err != nil {
		log.Errorf("Unable to delete interface %s from network-instance %s: %e", name, netInstName, err)
		return err
	}
	// Delete the interface
	path = fmt.Sprintf("/interface[name=%s]", name)
	_, err = gnmi.Delete(context.Background(), path)
	if err != nil {
		log.Errorf("Unable to delete interface %s: %e", name, err)
		return err
	}

	return err
}

// AddOrExistsNetworkInstance will either add or return true for a given network instance
func AddOrExistsNetworkInstance(conf *types.NetConf, args *skel.CmdArgs) error {
	var err error

	if found := GetNetworkInstance(conf.HostConf.BridgeConf.BridgeName); !found {
		logging.Debugf("AddOrExistsNetworkInstance(): Bridge %s not found, creating", conf.HostConf.BridgeConf.BridgeName)
		err = AddNetworkInstance(conf.HostConf.BridgeConf.BridgeName, "mac-vrf")

		//if err == nil {
		//	// Bridge is always created because it is required for interface.
		//	// If bridge type was actually called out, then set the
		//	// bridge up as L2 bridge. Otherwise, a controller is
		//	// responsible for writing flows to OvS.
		//	if conf.HostConf.NetType == "bridge" {
		//		err = configL2Bridge(conf.HostConf.BridgeConf.BridgeName)
		//	}
		//}
	} else {
		logging.Debugf("addLocalNetworkBridge(): Bridge %s exists, skip creating", conf.HostConf.BridgeConf.BridgeName)
	}

	return err
}

//func DeleteNetworkInstanceIfEmpty(conf *types.NetConf, args *skel.CmdArgs, data *OvsSavedData) error {
//	var err error
//
//	if containInterfaces := doesBridgeContainInterfaces(conf.HostConf.BridgeConf.BridgeName); containInterfaces == false {
//		logging.Debugf("delSRLinuxMACVRF(): No interfaces found, deleting Bridge %s", conf.HostConf.BridgeConf.BridgeName)
//		err = DeleteNetworkInstance(conf.HostConf.BridgeConf.BridgeName)
//	} else {
//		logging.Debugf("delSRLinuxMACVRF(): Interfaces found, skip deleting Bridge %s", conf.HostConf.BridgeConf.BridgeName)
//	}
//	return err
//}
