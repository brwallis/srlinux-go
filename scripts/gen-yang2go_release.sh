#!/bin/bash
YANG_SRC_PATH="../yang_release/models"
SRL_SRC_PATH="${YANG_SRC_PATH}/srl_nokia/models"
GO_OUT_PATH="../pkg/yang_release"

echo "Generate Go bindings for SR Linux YANG modules"
go mod download

YGOT_DIR=`go list -f '{{ .Dir }}' -m github.com/openconfig/ygot`

mkdir -p ${GO_OUT_PATH}
go run $YGOT_DIR/generator/generator.go \
   -path=${YANG_SRC_PATH}/ -output_file=${GO_OUT_PATH}/model.go -package_name=srl_yang_release -generate_fakeroot \
   ${SRL_SRC_PATH}/acl/srl_nokia-acl.yang \
   ${SRL_SRC_PATH}/acl/srl_nokia-packet-match-types.yang \
   ${SRL_SRC_PATH}/bfd/srl_nokia-bfd.yang \
   ${SRL_SRC_PATH}/bfd/srl_nokia-micro-bfd.yang \
   ${SRL_SRC_PATH}/common/srl_nokia-extensions.yang \
   ${SRL_SRC_PATH}/common/srl_nokia-common.yang \
   ${SRL_SRC_PATH}/common/srl_nokia-features.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-if-ip.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces-bridge-table-mac-duplication-entries.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces-bridge-table-mac-learning-entries.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces-bridge-table-mac-table.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces-bridge-table-statistics.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces-bridge-table.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces-ethernet-segment-association.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces-ip-dhcp-relay.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces-ip-dhcp.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces-ip-vrrp.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces-lag.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces-local-mirror-destination.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces-nbr-evpn.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces-nbr.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces-router-adv.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces-vlans.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-interfaces.yang \
   ${SRL_SRC_PATH}/interfaces/srl_nokia-lacp.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-aggregate-routes.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-bgp-evpn.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-bgp-vpn.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-bgp.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-bridge-table-mac-duplication-entries.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-bridge-table-mac-duplication.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-bridge-table-mac-learning-entries.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-bridge-table-mac-learning.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-bridge-table-mac-limit.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-bridge-table-mac-table.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-bridge-table-static-mac.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-bridge-table.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-icmp.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-ip-route-tables.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-isis.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-ldp.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-linux.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-mpls-route-tables.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-mpls.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-network-instance-mtu.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-network-instance.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-next-hop-groups.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-ospf-lsdb.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-ospf-types.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-ospf.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-ospfv3-lsas.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-rib-bgp.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-static-routes.yang \
   ${SRL_SRC_PATH}/network-instance/srl_nokia-tcp-udp.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-acl.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-cgroup.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-chassis.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-control.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-cpu.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-datapath-resources.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-disk.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-fabric.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-fan.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-lc.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-memory.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-mtu.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-psu.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-qos.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-redundancy.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-resource-mgmt.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform-resource-monitoring.yang \
   ${SRL_SRC_PATH}/platform/srl_nokia-platform.yang \
   ${SRL_SRC_PATH}/qos/srl_nokia-qos.yang \
   ${SRL_SRC_PATH}/routing-policy/srl_nokia-policy-types.yang \
   ${SRL_SRC_PATH}/routing-policy/srl_nokia-routing-policy.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-aaa-types.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-aaa.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-app-mgmt.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-boot.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-configuration-role.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-configuration.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-dns.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-ftp.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-gnmi-server.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-json-rpc.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-keychains.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-lldp-types.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-lldp.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-load-balancing.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-logging.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-maintenance-mode.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-mirroring.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-mpls-label-management.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-mtu.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-ntp.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-sflow.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-snmp-trace.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-snmp.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-ssh.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-system-banner.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-system-bridge-table.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-system-info.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-system-name.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-system-network-instance-bgp-evpn-ethernet-segments.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-system-network-instance-bgp-vpn.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-system-network-instance.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-system.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-timezone.yang \
   ${SRL_SRC_PATH}/system/srl_nokia-tls.yang \
   ${SRL_SRC_PATH}/tunnel/srl_nokia-tunnel-interfaces-vxlan-interface-bridge-table-multicast-destinations.yang \
   ${SRL_SRC_PATH}/tunnel/srl_nokia-tunnel-interfaces-vxlan-interface-bridge-table-unicast-destinations.yang \
   ${SRL_SRC_PATH}/tunnel/srl_nokia-tunnel-interfaces-vxlan-interface-bridge-table-unicast-es-destination-vteps.yang \
   ${SRL_SRC_PATH}/tunnel/srl_nokia-tunnel-interfaces-vxlan-interface-bridge-table.yang \
   ${SRL_SRC_PATH}/tunnel/srl_nokia-tunnel-interfaces.yang \
   ${SRL_SRC_PATH}/tunnel/srl_nokia-tunnel.yang \
   ${SRL_SRC_PATH}/tunnel/srl_nokia-vxlan-tunnel-vtep.yang 

go mod tidy
