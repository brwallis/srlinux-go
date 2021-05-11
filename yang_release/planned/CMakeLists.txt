###########################################################################
# Description:
#
# Copyright (c) 2018 Nokia
###########################################################################

# Target name
set(target yang_planned_modules)

message(STATUS "Building Lib ${target}")

###############################################################################
#
# To generate an extra .proto file:
#    - Add your .yang files to PROTO_YANG_FILES_DEPS
#       *note* your module name must match the filename!
#    - Add their directories to PROTO_YANG_DIRECTORIES
#    - Add their modules to PROTO_YANG_MODULES
#
###############################################################################

# Sources
set(cc_output_dir    "${PROJECT_BINARY_DIR}/yang-planned-modules-proto")
set(proto_output_dir "${PROJECT_BINARY_DIR}/protos/yang-planned")
set(protoc           "${PROJECT_BINARY_DIR}/third_party/install/protobuf/bin/protoc")
set(source_dir       "${CMAKE_CURRENT_SOURCE_DIR}")
set(yang2proto       "${PROJECT_BINARY_DIR}/bin/yang2proto")
set(yang_output_dir  "${PROJECT_BINARY_DIR}/bin/yang-planned")
set(header_dir       "${PROJECT_BINARY_DIR}/include/yang-planned")

# Make YANG files available in the bin-directory
srlinux_symlink_dir(${source_dir}/models ${yang_output_dir} "")

# List of YANG files dependencies which are used as source files for yang->proto generation
# cmake will make sure this target is rebuilt is any of them changes
set(PROTO_YANG_FILES_DEPS
    # these are built-in libyang (just copied here for reference purposes, not really sourced)
    ${source_dir}/models/ietf/ietf-datastores@2017-08-17.yang
    ${source_dir}/models/ietf/ietf-inet-types@2013-07-15.yang
    ${source_dir}/models/ietf/ietf-yang-library@2016-06-21.yang
    ${source_dir}/models/ietf/ietf-yang-library@2017-08-17.yang
    ${source_dir}/models/ietf/ietf-yang-metadata@2016-08-05.yang
    ${source_dir}/models/ietf/ietf-yang-types@2013-07-15.yang
    ${source_dir}/models/ietf/yang@2017-02-20.yang
    # regular yang files
    ${source_dir}/models/iana/iana-crypt-hash.yang
    ${source_dir}/models/iana/iana-if-type.yang
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
    ${source_dir}/models/srl_nokia/models/bfd/srl_nokia-bfd.yang
    ${source_dir}/models/srl_nokia/models/common/srl_nokia-extensions.yang
    ${source_dir}/models/srl_nokia/models/common/srl_nokia-common.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-if-ip.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces-nbr.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-interfaces.yang
    ${source_dir}/models/srl_nokia/models/interfaces/srl_nokia-tools-interfaces.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-aggregate-routes.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-bgp.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-linux.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-mpls.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-network-instance.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-next-hop-groups.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-static-routes.yang
    ${source_dir}/models/srl_nokia/models/network-instance/srl_nokia-rib-bgp.yang
    ${source_dir}/models/srl_nokia/models/routing-policy/srl_nokia-routing-policy.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-aaa-tacacs.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-aaa-types.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-aaa.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-app-mgmt.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-gnmi-server.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-json-rpc.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-lldp-types.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-lldp.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-logging.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-maintenance-mode.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-system.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-dns.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-ntp.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-ssh.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-ftp.yang
    ${source_dir}/models/srl_nokia/models/system/srl_nokia-snmp.yang
    ${source_dir}/models/srl_nokia/models/acl/srl_nokia-acl.yang
    ${source_dir}/models/srl_nokia/models/acl/srl_nokia-packet-match-types.yang
    ${source_dir}/models/srl_nokia/models/acl/srl_nokia-tools-acl.yang
    ${source_dir}/models/srl_nokia/models/qos/srl_nokia-qos.yang
)

set(PROTO_YANG_DIRECTORIES
    ${source_dir}/models
)

set(PROTO_YANG_MODULES
    srl_nokia-aaa
    srl_nokia-aaa-types
    srl_nokia-acl
    srl_nokia-tools-acl
    srl_nokia-app-mgmt
    srl_nokia-packet-match-types
    srl_nokia-aggregate-routes
    srl_nokia-bfd
    srl_nokia-bgp
    srl_nokia-gnmi-server
    srl_nokia-if-ip
    srl_nokia-interfaces
    srl_nokia-interfaces-nbr
    srl_nokia-tools-interfaces
    srl_nokia-json-rpc
    srl_nokia-linux
    srl_nokia-lldp
    srl_nokia-lldp-types
    srl_nokia-logging
    srl_nokia-mpls
    srl_nokia-maintenance-mode
    srl_nokia-network-instance
    srl_nokia-next-hop-groups
    srl_nokia-rib-bgp
    srl_nokia-routing-policy
    srl_nokia-qos
    srl_nokia-static-routes
    srl_nokia-system
    srl_nokia-dns
    srl_nokia-ntp
    srl_nokia-ssh
    srl_nokia-ftp
    srl_nokia-snmp
)

# List of generated '.proto' files
#  set(PROTO_FILES
#      ${proto_output_dir}/nokia_extensions.proto
#      ${proto_output_dir}/srl_nokia_system.proto
#      ${proto_output_dir}/srl_nokia_interfaces.proto
#      ${proto_output_dir}/srl_nokia_tools_interfaces.proto
#      ${proto_output_dir}/srl_nokia_network_instance.proto
#      ${proto_output_dir}/srl_nokia_routing_policy.proto
#      ${proto_output_dir}/srl_nokia_bfd.proto
#      ${proto_output_dir}/srl_nokia_acl.proto
#      ${proto_output_dir}/srl_nokia_tools_acl.proto
#      ${proto_output_dir}/srl_nokia_qos.proto
#  )

set(proto_stamp ${proto_output_dir}/generated)

# Generate .proto files
add_custom_command(
    OUTPUT ${proto_stamp}
    COMMAND ${CMAKE_COMMAND} -E make_directory ${proto_output_dir} &&
            ${yang2proto} --output-directory ${proto_output_dir} --yangs ${PROTO_YANG_DIRECTORIES} ${PROTO_YANG_MODULES}
    COMMAND touch ${proto_stamp}
    DEPENDS ${PROTO_YANG_FILES_DEPS} yang2proto
    COMMENT "Generating .proto files into folder ${proto_output_dir}"
)

# add_custom_target(yang_modules DEPENDS ${PROTO_YANG_FILES_DEPS} yang2proto ${PROTO_FILES} ${proto_stamp})
add_custom_target(yang_planned_modules DEPENDS ${proto_stamp})
