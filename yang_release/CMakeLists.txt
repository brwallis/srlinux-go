###########################################################################
# Description:
#
# Copyright (c) 2018 Nokia
###########################################################################

# Target name
set(target yang_modules)
set(target_zip yang_modules_zip)

message(STATUS "Building Lib ${target}")
message(STATUS "Building Lib ${target_zip}")

###############################################################################
#
# To generate an extra .proto file:
#    - Add your .yang files to PROTO_YANG_FILES_DEPS
#       *note* your module name must match the filename!
#    - Add their directories to PROTO_YANG_DIRECTORIES
#    - Add their modules to PROTO_YANG_MODULES
#    - Add the corresponding .proto output files to PROTO_FILES
#    --> the .proto files will be generated in the <build>/protos/yang directory
#    - The .proto files can then be used in srlinux_proto target to
#    -   generate c++/python/grpc classes
#    - See protos/CMakeLists.txt
#
###############################################################################

# Sources
set(cc_output_dir    "${PROJECT_BINARY_DIR}/yang-modules-proto")
set(proto_output_dir "${PROJECT_BINARY_DIR}/protos/yang")
set(protoc           "${PROJECT_BINARY_DIR}/third_party/install/protobuf/bin/protoc")
set(source_dir       "${CMAKE_CURRENT_SOURCE_DIR}")
set(yang2proto       "${PROJECT_BINARY_DIR}/bin/yang2proto")
set(yang_output_dir  "${PROJECT_BINARY_DIR}/bin")
set(header_dir       "${PROJECT_BINARY_DIR}/include/yang")

# Make YANG files available in the bin-directory
srlinux_symlink_dir(${source_dir}/models ${yang_output_dir} "")

# List of YANG files dependencies which are used as source files for yang->proto generation
# cmake will make sure this target is rebuilt is any of them changes
set(PROTO_YANG_FILES_DEPS
    # these are built-in libyang (just copied here for reference purposes, not really sourced)
    ${source_dir}/models/ietf/ietf-datastores@2017-08-17.yang
    ${source_dir}/models/ietf/ietf-datastores@2018-02-14.yang
    ${source_dir}/models/ietf/ietf-inet-types@2013-07-15.yang
    ${source_dir}/models/ietf/ietf-yang-library@2016-06-21.yang
    ${source_dir}/models/ietf/ietf-yang-library@2017-08-17.yang
    ${source_dir}/models/ietf/ietf-yang-library@2019-01-04.yang
    ${source_dir}/models/ietf/ietf-yang-metadata@2016-08-05.yang
    ${source_dir}/models/ietf/ietf-yang-types@2013-07-15.yang
    ${source_dir}/models/ietf/yang@2017-02-20.yang
    # regular yang files
    ${source_dir}/models/iana/iana-crypt-hash.yang
    ${source_dir}/models/iana/iana-if-type.yang
    ${source_dir}/models/iana/iana-timezones.yang
    ${source_dir}/models/ietf/ietf-interfaces.yang
    ${source_dir}/models/ietf/ietf-ip.yang
    ${source_dir}/models/ietf/ietf-ipfix-psamp.yang
    ${source_dir}/models/ietf/ietf-netconf-acm.yang
    ${source_dir}/models/ietf/ietf-netconf-monitoring.yang
    ${source_dir}/models/ietf/ietf-netconf-notifications.yang
    ${source_dir}/models/ietf/ietf-netconf-partial-lock.yang
    ${source_dir}/models/ietf/ietf-netconf-with-defaults.yang
    ${source_dir}/models/ietf/ietf-netconf.yang
    ${source_dir}/models/ietf/ietf-snmp-common.yang
    ${source_dir}/models/ietf/ietf-snmp-community.yang
    ${source_dir}/models/ietf/ietf-snmp-engine.yang
    ${source_dir}/models/ietf/ietf-snmp-notification.yang
    ${source_dir}/models/ietf/ietf-snmp-proxy.yang
    ${source_dir}/models/ietf/ietf-snmp-ssh.yang
    ${source_dir}/models/ietf/ietf-snmp-target.yang
    ${source_dir}/models/ietf/ietf-snmp-tls.yang
    ${source_dir}/models/ietf/ietf-snmp-tsm.yang
    ${source_dir}/models/ietf/ietf-snmp-usm.yang
    ${source_dir}/models/ietf/ietf-snmp-vacm.yang
    ${source_dir}/models/ietf/ietf-snmp.yang
    ${source_dir}/models/ietf/ietf-system.yang
    ${source_dir}/models/ietf/ietf-x509-cert-to-name.yang
    ${source_dir}/models/ietf/ietf-yang-smiv2.yang
    ${source_dir}/models/srl_nokia/models/acl/srl_nokia-acl.yang
    ${source_dir}/models/srl_nokia/models/acl/srl_nokia-packet-match-types.yang
    ${source_dir}/models/srl_nokia/models/acl/srl_nokia-tools-acl.yang
    ${source_dir}/models/srl_nokia/models/bfd/srl_nokia-bfd.yang
    ${source_dir}/models/srl_nokia/models/bfd/srl_nokia-micro-bfd.yang
    ${source_dir}/models/srl_nokia/models/bfd/srl_nokia-tools-bfd.yang
    ${source_dir}/models/srl_nokia/models/common/srl_nokia-common.yang
    ${source_dir}/models/srl_nokia/models/common/srl_nokia-extensions.yang
    ${source_dir}/models/srl_nokia/models/common/srl_nokia-features.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-if-ip.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces-bridge-table-mac-duplication-entries.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces-bridge-table-mac-learning-entries.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces-bridge-table-mac-table.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces-bridge-table-statistics.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces-bridge-table.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces-ethernet-segment-association.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces-ip-dhcp-relay.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces-ip-dhcp.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces-ip-vrrp.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces-lag.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces-nbr.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces-nbr-evpn.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces-router-adv.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces-vlans.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-lacp.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-tools-interfaces-bridge-table-mac-duplication.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-tools-interfaces-bridge-table-mac-learning.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-tools-interfaces-ip-dhcp-relay.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-tools-interfaces-ip-vrrp.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-tools-interfaces-nbr.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-tools-interfaces.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-aggregate-routes.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-bgp-evpn.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-bgp-vpn.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-bgp.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-bridge-table-mac-duplication-entries.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-bridge-table-mac-duplication.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-bridge-table-mac-learning-entries.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-bridge-table-mac-learning.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-bridge-table-mac-limit.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-bridge-table-mac-table.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-bridge-table-static-mac.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-bridge-table.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-icmp.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-ip-route-tables.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-ip-route-tables.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-isis.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-ldp.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-linux.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-mpls-route-tables.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-mpls.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-network-instance-mtu.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-network-instance.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-next-hop-groups.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-ospf-lsdb.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-ospf-types.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-ospf.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-ospfv3-lsas.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-rib-bgp.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-static-routes.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-tcp-udp.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-tools-bgp.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-tools-bridge-table-mac-duplication.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-tools-bridge-table-mac-learning.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-tools-icmp.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-tools-isis.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-tools-network-instance.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-tools-ospf.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-tunnel-tables.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-traffic-engineering.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-acl.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-cgroup.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-chassis.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-control.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-cpu.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-datapath-resources.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-disk.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-fabric.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-fan.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-lc.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-memory.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-mtu.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-psu.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-qos.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-redundancy.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-resource-mgmt.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform-resource-monitoring.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-platform.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-tools-platform-redundancy.yang
    ${source_dir}/models/srl_nokia/models/platform/srl_nokia-tools-platform.yang
    ${source_dir}/models/srl_nokia/models/qos/srl_nokia-qos.yang
    ${source_dir}/models/srl_nokia/models/qos/srl_nokia-tools-qos.yang
    ${source_dir}/models/srl_nokia/models/routing-policy/srl_nokia-policy-types.yang
    ${source_dir}/models/srl_nokia/models/routing-policy/srl_nokia-routing-policy.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-aaa-tacacs.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-aaa-types.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-aaa.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-app-mgmt.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-boot.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-configuration.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-configuration-role.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-dns.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-ftp.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-gnmi-server.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-json-rpc.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-keychains.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-lldp-types.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-lldp.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-load-balancing.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-logging.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-maintenance-mode.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-mpls-label-management.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-mtu.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-ntp.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-sflow.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-snmp-trace.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-snmp.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-ssh.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-system-banner.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-system-bridge-table.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-system-info.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-system-name.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-system-network-instance-bgp-evpn-ethernet-segments.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-system-network-instance-bgp-vpn.yang  
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-system-network-instance.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-system.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-timezone.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-tls.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-tools-app-mgmt.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-tools-boot.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-tools-cgroup.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-tools-configuration.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-tools-packet-trace.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-tools-system.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-tools-tls.yang
    ${source_dir}/models/srl_nokia/models/tunnel/srl_nokia-tools-tunnel.yang
    ${source_dir}/models/srl_nokia/models/tunnel/srl_nokia-tools-vxlan-tunnel.yang
    ${source_dir}/models/srl_nokia/models/tunnel/srl_nokia-tunnel-interfaces-vxlan-interface-bridge-table-multicast-destinations.yang
    ${source_dir}/models/srl_nokia/models/tunnel/srl_nokia-tunnel-interfaces-vxlan-interface-bridge-table-unicast-destinations.yang
    ${source_dir}/models/srl_nokia/models/tunnel/srl_nokia-tunnel-interfaces-vxlan-interface-bridge-table-unicast-es-destination-vteps.yang
    ${source_dir}/models/srl_nokia/models/tunnel/srl_nokia-tunnel-interfaces-vxlan-interface-bridge-table.yang
    ${source_dir}/models/srl_nokia/models/tunnel/srl_nokia-tunnel-interfaces.yang
    ${source_dir}/models/srl_nokia/models/tunnel/srl_nokia-tunnel.yang
    ${source_dir}/models/srl_nokia/models/tunnel/srl_nokia-vxlan-tunnel-vtep.yang
    ${source_dir}/models/srl_test/srl_test-mgmt-server-module-1.yang
    ${source_dir}/models/srl_test/srl_test-mgmt-server-module-4.yang
    ${source_dir}/models/srl_test/srl_test-mgmt-server-module-5.yang
    ${source_dir}/models/srl_test/srl_test-mgmt-server-module-6.yang
    ${source_dir}/models/srl_test/srl_test-mgmt-server-module-misc.yang
    ${source_dir}/models/srl_test/srl_test-sdk-agent.yang
    ${source_dir}/models/srl_test/srl_test-state-generator.yang
)

set(PROTO_YANG_DIRECTORIES
    ${source_dir}/models
    ${source_dir}/models/iana
    ${source_dir}/models/ietf
    ${source_dir}/models/srl_nokia/models/acl
    ${source_dir}/models/srl_nokia/models/bfd
    ${source_dir}/models/srl_nokia/models/common
    ${source_dir}/models/srl_nokia/models/interfaces
    ${source_dir}/models/srl_nokia/models/network-instance
    ${source_dir}/models/srl_nokia/models/platform
    ${source_dir}/models/srl_nokia/models/qos
    ${source_dir}/models/srl_nokia/models/routing-policy
    ${source_dir}/models/srl_nokia/models/system
    ${source_dir}/models/srl_nokia/models/tunnel
    ${source_dir}/models/srl_test
)

set(PROTO_YANG_MODULES
    srl_nokia-aaa
    srl_nokia-aaa-types
    srl_nokia-acl
    srl_nokia-aggregate-routes
    srl_nokia-app-mgmt
    srl_nokia-bfd
    srl_nokia-bgp
    srl_nokia-bgp-evpn
    srl_nokia-bgp-vpn
    srl_nokia-boot
    srl_nokia-bridge-table
    srl_nokia-bridge-table-mac-duplication
    srl_nokia-bridge-table-mac-duplication-entries
    srl_nokia-bridge-table-mac-learning
    srl_nokia-bridge-table-mac-learning-entries
    srl_nokia-bridge-table-mac-limit
    srl_nokia-bridge-table-mac-table
    srl_nokia-bridge-table-static-mac
    srl_nokia-configuration
    srl_nokia-configuration-role
    srl_nokia-dns
    srl_nokia-ftp
    srl_nokia-gnmi-server
    srl_nokia-icmp
    srl_nokia-if-ip
    srl_nokia-interfaces
    srl_nokia-interfaces-bridge-table
    srl_nokia-interfaces-bridge-table-mac-duplication-entries
    srl_nokia-interfaces-bridge-table-mac-learning-entries
    srl_nokia-interfaces-bridge-table-mac-table
    srl_nokia-interfaces-bridge-table-statistics
    srl_nokia-interfaces-ethernet-segment-association
    srl_nokia-interfaces-ip-dhcp
    srl_nokia-interfaces-ip-dhcp-relay
    srl_nokia-interfaces-ip-vrrp
    srl_nokia-interfaces-lag
    srl_nokia-interfaces-nbr
    srl_nokia-interfaces-nbr-evpn
    srl_nokia-interfaces-router-adv
    srl_nokia-interfaces-vlans
    srl_nokia-ip-route-tables
    srl_nokia-isis
    srl_nokia-ldp
    srl_nokia-json-rpc
    srl_nokia-keychains
    srl_nokia-lacp
    srl_nokia-linux
    srl_nokia-lldp
    srl_nokia-lldp-types
    srl_nokia-load-balancing
    srl_nokia-logging
    srl_nokia-maintenance-mode
    srl_nokia-micro-bfd
    srl_nokia-mpls
    srl_nokia-mpls-label-management
    srl_nokia-mpls-route-tables
    srl_nokia-mtu
    srl_nokia-network-instance
    srl_nokia-network-instance-mtu
    srl_nokia-next-hop-groups
    srl_nokia-ntp
    srl_nokia-ospf
    srl_nokia-packet-match-types
    srl_nokia-platform
    srl_nokia-platform-acl
    srl_nokia-platform-cgroup
    srl_nokia-platform-chassis
    srl_nokia-platform-control
    srl_nokia-platform-cpu
    srl_nokia-platform-datapath-resources
    srl_nokia-platform-disk
    srl_nokia-platform-fabric
    srl_nokia-platform-fan
    srl_nokia-platform-lc
    srl_nokia-platform-memory
    srl_nokia-platform-mtu
    srl_nokia-platform-psu
    srl_nokia-platform-qos
    srl_nokia-platform-redundancy
    srl_nokia-platform-resource-mgmt
    srl_nokia-platform-resource-monitoring
    srl_nokia-qos
    srl_nokia-rib-bgp
    srl_nokia-routing-policy
    srl_nokia-sflow
    srl_nokia-snmp
    srl_nokia-ssh
    srl_nokia-static-routes
    srl_nokia-system
    srl_nokia-system-banner
    srl_nokia-system-bridge-table
    srl_nokia-system-info
    srl_nokia-system-name
    srl_nokia-system-network-instance
    srl_nokia-system-network-instance-bgp-evpn-ethernet-segments
    srl_nokia-system-network-instance-bgp-vpn
    srl_nokia-tcp-udp
    srl_nokia-tls
    srl_nokia-tools-acl
    srl_nokia-tools-app-mgmt
    srl_nokia-tools-bfd
    srl_nokia-tools-bgp
    srl_nokia-tools-boot
    srl_nokia-tools-cgroup
    srl_nokia-tools-bridge-table-mac-duplication
    srl_nokia-tools-bridge-table-mac-learning
    srl_nokia-tools-configuration
    srl_nokia-tools-icmp
    srl_nokia-tools-interfaces
    srl_nokia-tools-interfaces-bridge-table-mac-duplication
    srl_nokia-tools-interfaces-bridge-table-mac-learning
    srl_nokia-tools-interfaces-ip-dhcp-relay
    srl_nokia-tools-interfaces-ip-vrrp
    srl_nokia-tools-interfaces-nbr
    srl_nokia-tools-isis
    srl_nokia-tools-network-instance
    srl_nokia-tools-ospf
    srl_nokia-tools-packet-trace
    srl_nokia-tools-platform
    srl_nokia-tools-platform-redundancy
    srl_nokia-tools-qos
    srl_nokia-tools-system
    srl_nokia-tools-tls
    srl_nokia-tools-tunnel
    srl_nokia-tools-vxlan-tunnel
    srl_nokia-traffic-engineering 
    srl_nokia-tunnel
    srl_nokia-tunnel-interfaces
    srl_nokia-tunnel-interfaces-vxlan-interface-bridge-table
    srl_nokia-tunnel-interfaces-vxlan-interface-bridge-table-multicast-destinations
    srl_nokia-tunnel-interfaces-vxlan-interface-bridge-table-unicast-destinations
    srl_nokia-tunnel-interfaces-vxlan-interface-bridge-table-unicast-es-destination-vteps
    srl_nokia-tunnel-tables
    srl_nokia-vxlan-tunnel-vtep
    srl_test-mgmt-server-module-1
    srl_test-mgmt-server-module-4
    srl_test-mgmt-server-module-5
    srl_test-mgmt-server-module-6
    srl_test-mgmt-server-module-misc
    srl_test-sdk-agent
    srl_test-state-generator
)

# List of generated '.proto' files
set(PROTO_FILES
    ${proto_output_dir}/nokia_extensions.proto
    ${proto_output_dir}/srl_nokia-boot.proto
    ${proto_output_dir}/srl_nokia_acl.proto
    ${proto_output_dir}/srl_nokia_bfd.proto
    ${proto_output_dir}/srl_nokia_interfaces.proto
    ${proto_output_dir}/srl_nokia_network_instance.proto
    ${proto_output_dir}/srl_nokia_platform.proto
    ${proto_output_dir}/srl_nokia_qos.proto
    ${proto_output_dir}/srl_nokia_routing_policy.proto
    ${proto_output_dir}/srl_nokia_system.proto
    ${proto_output_dir}/srl_nokia_tools_acl.proto
    ${proto_output_dir}/srl_nokia_tools_bfd.proto
    ${proto_output_dir}/srl_nokia_tools_boot.proto
    ${proto_output_dir}/srl_nokia_tools_cgroup.proto
    ${proto_output_dir}/srl_nokia_tools_interfaces.proto
    ${proto_output_dir}/srl_nokia_tools_network_instance.proto
    ${proto_output_dir}/srl_nokia_tools_platform.proto
    ${proto_output_dir}/srl_nokia_tools_system.proto
    ${proto_output_dir}/srl_nokia_tools_tunnel.proto
    ${proto_output_dir}/srl_nokia_tunnel.proto
    ${proto_output_dir}/srl_nokia_tunnel_interfaces.proto
    ${proto_output_dir}/srl_test_mgmt_server_module_1.proto
    ${proto_output_dir}/srl_test_mgmt_server_module_4.proto
    ${proto_output_dir}/srl_test_mgmt_server_module_misc.proto
    ${proto_output_dir}/srl_test_sdk_agent.proto
    ${proto_output_dir}/srl_test_state_generator.proto
)

set(proto_stamp ${proto_output_dir}/generated)

# Generate .proto files
add_custom_command(
    OUTPUT ${proto_stamp}
    COMMAND ${CMAKE_COMMAND} -E make_directory ${proto_output_dir} &&
            ${yang2proto} --output-directory ${proto_output_dir} --yangs ${PROTO_YANG_DIRECTORIES} ${PROTO_YANG_MODULES}
    COMMAND touch ${proto_stamp}
    DEPENDS ${PROTO_YANG_FILES_DEPS} yang2proto
    COMMENT "Generating .proto files into folder ${proto_output_dir}"
    BYPRODUCTS ${PROTO_FILES}
)

add_custom_target(${target} DEPENDS ${proto_stamp} ${target_zip})

set(yangzip_output_dir ${PROJECT_BINARY_DIR}/test)
set(YANGZIP_FILES
    ${yangzip_output_dir}/YANG.zip
    ${yangzip_output_dir}/TOOLSYANG.zip
)
set(TOOLSDIR ${SRLINUX_SOURCE_DIR}/test_yang_tools)
set(yangzip ${TOOLSDIR}/build_yang_zip.sh)
# Generate YANG.zip and TOOLSYANG.zip files when any .yang file has been modified
add_custom_command(
    OUTPUT ${YANGZIP_FILES}
    COMMAND ${CMAKE_COMMAND} -E make_directory ${yangzip_output_dir} &&
            ${yangzip} ${yangzip_output_dir} ${PROTO_YANG_FILES_DEPS}  ${source_dir}/test_yang_tools/fib_agent.yang
    DEPENDS ${PROTO_YANG_FILES_DEPS}
    COMMENT "Generating YANG.zip and TOOLSYANG.zip files into folder ${yangzip_output_dir}"
)
add_custom_target(${target_zip} DEPENDS ${proto_stamp} ${YANGZIP_FILES})

#RPM
if(${BUILD_RPM})
    include("${SRLINUX_SOURCE_DIR}/cmake/Srlinuxpackage.cmake")
    configure_file("${CMAKE_CURRENT_SOURCE_DIR}/rpm.spec.in" "${CMAKE_CURRENT_BINARY_DIR}/rpm.spec" @ONLY IMMEDIATE)
    srlinux_add_package(${target})

    set(target ${target_zip})
    configure_file("${CMAKE_CURRENT_SOURCE_DIR}/yang_zip.rpm.spec.in" "${CMAKE_CURRENT_BINARY_DIR}/yang_zip.rpm.spec" @ONLY IMMEDIATE)
    srlinux_add_package(${target} RPMSPEC ${CMAKE_CURRENT_BINARY_DIR}/yang_zip.rpm.spec)
    add_dependencies(${target_zip}_rpm yang_modules_rpm)

endif(${BUILD_RPM})


set(target_telemetry_mgr_impl telemetry_mgr_impl)
message(STATUS "Building Lib ${target_telemetry_mgr_impl}")

#file(GLOB telemetry_mgr_generated_srcs ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_*_mgr_impl.cc)
set(telemetry_mgr_sources
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_aaa_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_acl_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_app_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_arp_nd_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_bfd_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_bgp_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_chassis_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_common_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_dhcp_relay_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_evpn_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_fib_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_gnmi_server_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_isis_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_ldp_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_json_rpc_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_l2_mac_learn_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_l2_mac_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_l2_static_mac_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_label_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_lag_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_linux_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_lldp_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_log_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_maint_mode_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_mcid_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_mgmt_server_test_app_1_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_mgmt_server_test_app_4_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_mgmt_server_test_app_5_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_mgmt_server_test_app_6_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_mgmt_server_test_app_misc_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_mpls_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_net_inst_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_ospf_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_plcy_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_qos_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_sdk_apple_agent_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_state_generator_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_static_route_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_te_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_vrrp_mgr_mgr_impl.cc
    ${PROJECT_BINARY_DIR}/protos/telemetry/telemetry_vxlan_mgr_mgr_impl.cc
)

set(MGR_IMPL_DEPS
    srlinux_proto
    protobuf
    telemetry
    srlinux_telemetry
    srlinux_telemetry_test
)
set (TARGET_DEPS
    srlinux_telemetry_test
)
set(sources ${telemetry_mgr_sources} ${HDRS})
set_source_files_properties(${sources} PROPERTIES GENERATED TRUE)
cc_library(${target_telemetry_mgr_impl} SHARED SRCS ${sources} DEPS ${MGR_IMPL_DEPS} TARGET_DEPS ${TARGET_DEPS})
set_target_properties(${target_telemetry_mgr_impl} PROPERTIES COMPILE_FLAGS -s)

add_subdirectory(planned)


if(${BUILD_RPM})
    set(target ${target_telemetry_mgr_impl})
    include("${SRLINUX_SOURCE_DIR}/cmake/Srlinuxpackage.cmake")
    configure_file("${SRLINUX_SOURCE_DIR}/maketools/rpm.spec.in" "${CMAKE_CURRENT_BINARY_DIR}/telemetry_impl_rpm.spec" @ONLY IMMEDIATE)
    srlinux_add_package(${target} RPMSPEC "${CMAKE_CURRENT_BINARY_DIR}/telemetry_impl_rpm.spec")
    add_dependencies(${target}_rpm yang_modules_rpm ${target_zip}_rpm)
endif(${BUILD_RPM})
