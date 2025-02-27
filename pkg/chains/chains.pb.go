// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: zetachain/zetacore/pkg/chains/chains.proto

package chains

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// ReceiveStatus represents the status of an outbound
// TODO: Rename and move
// https://github.com/zeta-chain/node/issues/2257
type ReceiveStatus int32

const (
	// Created is used for inbounds
	ReceiveStatus_created ReceiveStatus = 0
	ReceiveStatus_success ReceiveStatus = 1
	ReceiveStatus_failed  ReceiveStatus = 2
)

var ReceiveStatus_name = map[int32]string{
	0: "created",
	1: "success",
	2: "failed",
}

var ReceiveStatus_value = map[string]int32{
	"created": 0,
	"success": 1,
	"failed":  2,
}

func (x ReceiveStatus) String() string {
	return proto.EnumName(ReceiveStatus_name, int32(x))
}

func (ReceiveStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{0}
}

// ChainName represents the name of the chain
// Deprecated(v19): replaced with Chain.Name as string
type ChainName int32

const (
	ChainName_empty            ChainName = 0
	ChainName_eth_mainnet      ChainName = 1
	ChainName_zeta_mainnet     ChainName = 2
	ChainName_btc_mainnet      ChainName = 3
	ChainName_polygon_mainnet  ChainName = 4
	ChainName_bsc_mainnet      ChainName = 5
	ChainName_goerli_testnet   ChainName = 6
	ChainName_mumbai_testnet   ChainName = 7
	ChainName_bsc_testnet      ChainName = 10
	ChainName_zeta_testnet     ChainName = 11
	ChainName_btc_testnet      ChainName = 12
	ChainName_sepolia_testnet  ChainName = 13
	ChainName_goerli_localnet  ChainName = 14
	ChainName_btc_regtest      ChainName = 15
	ChainName_amoy_testnet     ChainName = 16
	ChainName_optimism_mainnet ChainName = 17
	ChainName_optimism_sepolia ChainName = 18
	ChainName_base_mainnet     ChainName = 19
	ChainName_base_sepolia     ChainName = 20
	ChainName_solana_mainnet   ChainName = 21
	ChainName_solana_devnet    ChainName = 22
	ChainName_solana_localnet  ChainName = 23
)

var ChainName_name = map[int32]string{
	0:  "empty",
	1:  "eth_mainnet",
	2:  "zeta_mainnet",
	3:  "btc_mainnet",
	4:  "polygon_mainnet",
	5:  "bsc_mainnet",
	6:  "goerli_testnet",
	7:  "mumbai_testnet",
	10: "bsc_testnet",
	11: "zeta_testnet",
	12: "btc_testnet",
	13: "sepolia_testnet",
	14: "goerli_localnet",
	15: "btc_regtest",
	16: "amoy_testnet",
	17: "optimism_mainnet",
	18: "optimism_sepolia",
	19: "base_mainnet",
	20: "base_sepolia",
	21: "solana_mainnet",
	22: "solana_devnet",
	23: "solana_localnet",
}

var ChainName_value = map[string]int32{
	"empty":            0,
	"eth_mainnet":      1,
	"zeta_mainnet":     2,
	"btc_mainnet":      3,
	"polygon_mainnet":  4,
	"bsc_mainnet":      5,
	"goerli_testnet":   6,
	"mumbai_testnet":   7,
	"bsc_testnet":      10,
	"zeta_testnet":     11,
	"btc_testnet":      12,
	"sepolia_testnet":  13,
	"goerli_localnet":  14,
	"btc_regtest":      15,
	"amoy_testnet":     16,
	"optimism_mainnet": 17,
	"optimism_sepolia": 18,
	"base_mainnet":     19,
	"base_sepolia":     20,
	"solana_mainnet":   21,
	"solana_devnet":    22,
	"solana_localnet":  23,
}

func (x ChainName) String() string {
	return proto.EnumName(ChainName_name, int32(x))
}

func (ChainName) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{1}
}

// Network represents the network of the chain
// there is a single instance of the network on mainnet
// then the network can have eventual testnets or devnets
type Network int32

const (
	Network_eth      Network = 0
	Network_zeta     Network = 1
	Network_btc      Network = 2
	Network_polygon  Network = 3
	Network_bsc      Network = 4
	Network_optimism Network = 5
	Network_base     Network = 6
	Network_solana   Network = 7
)

var Network_name = map[int32]string{
	0: "eth",
	1: "zeta",
	2: "btc",
	3: "polygon",
	4: "bsc",
	5: "optimism",
	6: "base",
	7: "solana",
}

var Network_value = map[string]int32{
	"eth":      0,
	"zeta":     1,
	"btc":      2,
	"polygon":  3,
	"bsc":      4,
	"optimism": 5,
	"base":     6,
	"solana":   7,
}

func (x Network) String() string {
	return proto.EnumName(Network_name, int32(x))
}

func (Network) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{2}
}

// NetworkType represents the network type of the chain
// Mainnet, Testnet, Privnet, Devnet
type NetworkType int32

const (
	NetworkType_mainnet NetworkType = 0
	NetworkType_testnet NetworkType = 1
	NetworkType_privnet NetworkType = 2
	NetworkType_devnet  NetworkType = 3
)

var NetworkType_name = map[int32]string{
	0: "mainnet",
	1: "testnet",
	2: "privnet",
	3: "devnet",
}

var NetworkType_value = map[string]int32{
	"mainnet": 0,
	"testnet": 1,
	"privnet": 2,
	"devnet":  3,
}

func (x NetworkType) String() string {
	return proto.EnumName(NetworkType_name, int32(x))
}

func (NetworkType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{3}
}

// Vm represents the virtual machine type of the chain to support smart
// contracts
type Vm int32

const (
	Vm_no_vm Vm = 0
	Vm_evm   Vm = 1
	Vm_svm   Vm = 2
)

var Vm_name = map[int32]string{
	0: "no_vm",
	1: "evm",
	2: "svm",
}

var Vm_value = map[string]int32{
	"no_vm": 0,
	"evm":   1,
	"svm":   2,
}

func (x Vm) String() string {
	return proto.EnumName(Vm_name, int32(x))
}

func (Vm) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{4}
}

// Consensus represents the consensus algorithm used by the chain
// this can represent the consensus of a L1
// this can also represent the solution of a L2
type Consensus int32

const (
	Consensus_ethereum         Consensus = 0
	Consensus_tendermint       Consensus = 1
	Consensus_bitcoin          Consensus = 2
	Consensus_op_stack         Consensus = 3
	Consensus_solana_consensus Consensus = 4
)

var Consensus_name = map[int32]string{
	0: "ethereum",
	1: "tendermint",
	2: "bitcoin",
	3: "op_stack",
	4: "solana_consensus",
}

var Consensus_value = map[string]int32{
	"ethereum":         0,
	"tendermint":       1,
	"bitcoin":          2,
	"op_stack":         3,
	"solana_consensus": 4,
}

func (x Consensus) String() string {
	return proto.EnumName(Consensus_name, int32(x))
}

func (Consensus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{5}
}

// CCTXGateway describes for the chain the gateway used to handle CCTX outbounds
type CCTXGateway int32

const (
	// zevm is the internal CCTX gateway to process outbound on the ZEVM and read
	// inbound events from the ZEVM only used for ZetaChain chains
	CCTXGateway_zevm CCTXGateway = 0
	// observers is the CCTX gateway for chains relying on the observer set to
	// observe inbounds and TSS for outbounds
	CCTXGateway_observers CCTXGateway = 1
)

var CCTXGateway_name = map[int32]string{
	0: "zevm",
	1: "observers",
}

var CCTXGateway_value = map[string]int32{
	"zevm":      0,
	"observers": 1,
}

func (x CCTXGateway) String() string {
	return proto.EnumName(CCTXGateway_name, int32(x))
}

func (CCTXGateway) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{6}
}

// Chain represents static data about a blockchain network
// it is identified by a unique chain ID
type Chain struct {
	// ChainId is the unique identifier of the chain
	ChainId int64 `protobuf:"varint,2,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	// ChainName is the name of the chain
	// Deprecated(v19): replaced with Name
	ChainName ChainName `protobuf:"varint,1,opt,name=chain_name,json=chainName,proto3,enum=zetachain.zetacore.pkg.chains.ChainName" json:"chain_name,omitempty"` // Deprecated: Do not use.
	// Network is the network of the chain
	Network Network `protobuf:"varint,3,opt,name=network,proto3,enum=zetachain.zetacore.pkg.chains.Network" json:"network,omitempty"`
	// NetworkType is the network type of the chain: mainnet, testnet, etc..
	NetworkType NetworkType `protobuf:"varint,4,opt,name=network_type,json=networkType,proto3,enum=zetachain.zetacore.pkg.chains.NetworkType" json:"network_type,omitempty"`
	// Vm is the virtual machine used in the chain
	Vm Vm `protobuf:"varint,5,opt,name=vm,proto3,enum=zetachain.zetacore.pkg.chains.Vm" json:"vm,omitempty"`
	// Consensus is the underlying consensus algorithm used by the chain
	Consensus Consensus `protobuf:"varint,6,opt,name=consensus,proto3,enum=zetachain.zetacore.pkg.chains.Consensus" json:"consensus,omitempty"`
	// IsExternal describe if the chain is ZetaChain or external
	IsExternal bool `protobuf:"varint,7,opt,name=is_external,json=isExternal,proto3" json:"is_external,omitempty"`
	// CCTXGateway is the gateway used to handle CCTX outbounds
	CctxGateway CCTXGateway `protobuf:"varint,8,opt,name=cctx_gateway,json=cctxGateway,proto3,enum=zetachain.zetacore.pkg.chains.CCTXGateway" json:"cctx_gateway,omitempty"`
	// Name is the name of the chain
	Name string `protobuf:"bytes,9,opt,name=name,proto3" json:"name,omitempty"`
}

func (m *Chain) Reset()         { *m = Chain{} }
func (m *Chain) String() string { return proto.CompactTextString(m) }
func (*Chain) ProtoMessage()    {}
func (*Chain) Descriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{0}
}
func (m *Chain) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Chain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Chain.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Chain) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Chain.Merge(m, src)
}
func (m *Chain) XXX_Size() int {
	return m.Size()
}
func (m *Chain) XXX_DiscardUnknown() {
	xxx_messageInfo_Chain.DiscardUnknown(m)
}

var xxx_messageInfo_Chain proto.InternalMessageInfo

func (m *Chain) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

// Deprecated: Do not use.
func (m *Chain) GetChainName() ChainName {
	if m != nil {
		return m.ChainName
	}
	return ChainName_empty
}

func (m *Chain) GetNetwork() Network {
	if m != nil {
		return m.Network
	}
	return Network_eth
}

func (m *Chain) GetNetworkType() NetworkType {
	if m != nil {
		return m.NetworkType
	}
	return NetworkType_mainnet
}

func (m *Chain) GetVm() Vm {
	if m != nil {
		return m.Vm
	}
	return Vm_no_vm
}

func (m *Chain) GetConsensus() Consensus {
	if m != nil {
		return m.Consensus
	}
	return Consensus_ethereum
}

func (m *Chain) GetIsExternal() bool {
	if m != nil {
		return m.IsExternal
	}
	return false
}

func (m *Chain) GetCctxGateway() CCTXGateway {
	if m != nil {
		return m.CctxGateway
	}
	return CCTXGateway_zevm
}

func (m *Chain) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterEnum("zetachain.zetacore.pkg.chains.ReceiveStatus", ReceiveStatus_name, ReceiveStatus_value)
	proto.RegisterEnum("zetachain.zetacore.pkg.chains.ChainName", ChainName_name, ChainName_value)
	proto.RegisterEnum("zetachain.zetacore.pkg.chains.Network", Network_name, Network_value)
	proto.RegisterEnum("zetachain.zetacore.pkg.chains.NetworkType", NetworkType_name, NetworkType_value)
	proto.RegisterEnum("zetachain.zetacore.pkg.chains.Vm", Vm_name, Vm_value)
	proto.RegisterEnum("zetachain.zetacore.pkg.chains.Consensus", Consensus_name, Consensus_value)
	proto.RegisterEnum("zetachain.zetacore.pkg.chains.CCTXGateway", CCTXGateway_name, CCTXGateway_value)
	proto.RegisterType((*Chain)(nil), "zetachain.zetacore.pkg.chains.Chain")
}

func init() {
	proto.RegisterFile("zetachain/zetacore/pkg/chains/chains.proto", fileDescriptor_236b85e7bff6130d)
}

var fileDescriptor_236b85e7bff6130d = []byte{
	// 765 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xcd, 0x8e, 0xe4, 0x34,
	0x10, 0xee, 0x24, 0xfd, 0x5b, 0x99, 0x1f, 0xaf, 0x77, 0x80, 0xb0, 0x12, 0xcd, 0xc0, 0x01, 0x9a,
	0x11, 0xf4, 0x08, 0x38, 0x72, 0x41, 0x3b, 0x62, 0x11, 0x42, 0xec, 0x21, 0xac, 0x56, 0x88, 0x4b,
	0x70, 0xdc, 0x45, 0xda, 0x9a, 0x38, 0x8e, 0x62, 0x77, 0x76, 0x9b, 0xa7, 0xe0, 0x21, 0x38, 0x20,
	0xf1, 0x22, 0x1c, 0xf7, 0xc8, 0x11, 0xcd, 0x3c, 0x08, 0xc8, 0x8e, 0x93, 0x5e, 0x0e, 0xec, 0xce,
	0xa9, 0xed, 0xaf, 0xbf, 0xaf, 0xea, 0xab, 0x2a, 0x57, 0xe0, 0xe2, 0x17, 0x34, 0x8c, 0x6f, 0x99,
	0xa8, 0x2e, 0xdd, 0x49, 0x35, 0x78, 0x59, 0x5f, 0x17, 0x97, 0x0e, 0xd2, 0xfe, 0x67, 0x5d, 0x37,
	0xca, 0x28, 0xfa, 0xce, 0xc0, 0x5d, 0xf7, 0xdc, 0x75, 0x7d, 0x5d, 0xac, 0x3b, 0xd2, 0x83, 0xb3,
	0x42, 0x15, 0xca, 0x31, 0x2f, 0xed, 0xa9, 0x13, 0xbd, 0xff, 0x4f, 0x04, 0x93, 0x2b, 0x4b, 0xa0,
	0x6f, 0xc3, 0xdc, 0x31, 0x33, 0xb1, 0x49, 0xc2, 0xf3, 0x60, 0x15, 0xa5, 0x33, 0x77, 0xff, 0x66,
	0x43, 0xbf, 0x05, 0xe8, 0xfe, 0xaa, 0x98, 0xc4, 0x24, 0x38, 0x0f, 0x56, 0x27, 0x9f, 0xad, 0xd6,
	0xaf, 0x4c, 0xb7, 0x76, 0x41, 0x1f, 0x33, 0x89, 0x0f, 0xc3, 0x24, 0x48, 0x17, 0xbc, 0xbf, 0xd2,
	0x2f, 0x61, 0x56, 0xa1, 0x79, 0xa6, 0x9a, 0xeb, 0x24, 0x72, 0x91, 0x3e, 0x78, 0x4d, 0xa4, 0xc7,
	0x1d, 0x3b, 0xed, 0x65, 0xf4, 0x3b, 0x38, 0xf2, 0xc7, 0xcc, 0xec, 0x6b, 0x4c, 0xc6, 0x2e, 0xcc,
	0xc5, 0xdd, 0xc2, 0x3c, 0xd9, 0xd7, 0x98, 0xc6, 0xd5, 0xe1, 0x42, 0x3f, 0x85, 0xb0, 0x95, 0xc9,
	0xc4, 0x05, 0x79, 0xef, 0x35, 0x41, 0x9e, 0xca, 0x34, 0x6c, 0x25, 0x7d, 0x04, 0x0b, 0xae, 0x2a,
	0x8d, 0x95, 0xde, 0xe9, 0x64, 0x7a, 0xb7, 0x7e, 0xf4, 0xfc, 0xf4, 0x20, 0xa5, 0xef, 0x42, 0x2c,
	0x74, 0x86, 0xcf, 0x0d, 0x36, 0x15, 0x2b, 0x93, 0xd9, 0x79, 0xb0, 0x9a, 0xa7, 0x20, 0xf4, 0x57,
	0x1e, 0xb1, 0xa5, 0x72, 0x6e, 0x9e, 0x67, 0x05, 0x33, 0xf8, 0x8c, 0xed, 0x93, 0xf9, 0x9d, 0x4a,
	0xbd, 0xba, 0x7a, 0xf2, 0xc3, 0xd7, 0x9d, 0x22, 0x8d, 0xad, 0xde, 0x5f, 0x28, 0x85, 0xb1, 0x1b,
	0xe1, 0xe2, 0x3c, 0x58, 0x2d, 0x52, 0x77, 0xbe, 0xf8, 0x02, 0x8e, 0x53, 0xe4, 0x28, 0x5a, 0xfc,
	0xde, 0x30, 0xb3, 0xd3, 0x34, 0x86, 0x19, 0x6f, 0x90, 0x19, 0xdc, 0x90, 0x91, 0xbd, 0xe8, 0x1d,
	0xe7, 0xa8, 0x35, 0x09, 0x28, 0xc0, 0xf4, 0x67, 0x26, 0x4a, 0xdc, 0x90, 0xf0, 0xc1, 0xf8, 0xf7,
	0xdf, 0x96, 0xc1, 0xc5, 0x1f, 0x11, 0x2c, 0x86, 0x49, 0xd3, 0x05, 0x4c, 0x50, 0xd6, 0x66, 0x4f,
	0x46, 0xf4, 0x14, 0x62, 0x34, 0xdb, 0x4c, 0x32, 0x51, 0x55, 0x68, 0x48, 0x40, 0x09, 0x1c, 0x59,
	0xab, 0x03, 0x12, 0x5a, 0x4a, 0x6e, 0xf8, 0x00, 0x44, 0xf4, 0x3e, 0x9c, 0xd6, 0xaa, 0xdc, 0x17,
	0xaa, 0x1a, 0xc0, 0xb1, 0x63, 0xe9, 0x03, 0x6b, 0x42, 0x29, 0x9c, 0x14, 0x0a, 0x9b, 0x52, 0x64,
	0x06, 0xb5, 0xb1, 0xd8, 0xd4, 0x62, 0x72, 0x27, 0x73, 0x76, 0xc0, 0x66, 0xbd, 0xb0, 0x07, 0x60,
	0x70, 0xd0, 0x23, 0x71, 0xef, 0xa0, 0x07, 0x8e, 0xac, 0x03, 0x8d, 0xb5, 0x2a, 0xc5, 0x81, 0x75,
	0x6c, 0x41, 0x9f, 0xb0, 0x54, 0x9c, 0x95, 0x16, 0x3c, 0xe9, 0xa5, 0x0d, 0x16, 0x96, 0x48, 0x4e,
	0x6d, 0x74, 0x26, 0xd5, 0x7e, 0xd0, 0x11, 0x7a, 0x06, 0x44, 0xd5, 0x46, 0x48, 0xa1, 0xe5, 0x60,
	0xff, 0xde, 0x7f, 0x50, 0x9f, 0x8b, 0x50, 0xab, 0xce, 0x99, 0xc6, 0x81, 0x77, 0x7f, 0x40, 0x7a,
	0xce, 0x99, 0x2d, 0x52, 0xab, 0x92, 0x55, 0x87, 0x1e, 0xbe, 0x41, 0xef, 0xc1, 0xb1, 0xc7, 0x36,
	0xd8, 0x5a, 0xe8, 0x4d, 0x57, 0x43, 0x07, 0x0d, 0x76, 0xdf, 0xf2, 0xd3, 0x42, 0x98, 0xf9, 0x2d,
	0xa0, 0x33, 0x88, 0xd0, 0x6c, 0xc9, 0x88, 0xce, 0x61, 0x6c, 0xbb, 0x42, 0x02, 0x0b, 0xe5, 0x86,
	0x93, 0xd0, 0xce, 0xdc, 0xcf, 0x81, 0x44, 0x0e, 0xd5, 0x9c, 0x8c, 0xe9, 0x11, 0xcc, 0x7b, 0xe3,
	0x64, 0x62, 0x65, 0xd6, 0x1e, 0x99, 0xda, 0x47, 0xd1, 0xe5, 0x23, 0x33, 0x9f, 0xe6, 0x11, 0xc4,
	0x2f, 0x2d, 0x9b, 0x0d, 0xd7, 0x1b, 0x76, 0xef, 0xa9, 0xef, 0x50, 0xe0, 0x12, 0x35, 0xa2, 0xed,
	0x9e, 0x03, 0xc0, 0xd4, 0xd7, 0x10, 0xf9, 0x38, 0x1f, 0x42, 0xf8, 0x54, 0xda, 0x47, 0x55, 0xa9,
	0xac, 0x95, 0x64, 0xe4, 0x4c, 0xb7, 0xb2, 0xb3, 0xaa, 0x5b, 0x39, 0xbc, 0xc2, 0x9f, 0x60, 0x31,
	0xac, 0x97, 0xf5, 0x89, 0x66, 0x8b, 0x0d, 0xee, 0xac, 0xe4, 0x04, 0xc0, 0x60, 0xb5, 0xc1, 0x46,
	0x8a, 0xca, 0xa7, 0xcc, 0x85, 0xe1, 0x4a, 0x54, 0x24, 0xec, 0x4a, 0xca, 0xb4, 0x61, 0xfc, 0x9a,
	0x44, 0x76, 0x32, 0xbe, 0x71, 0xc3, 0x82, 0x92, 0xb1, 0xcf, 0xf0, 0x31, 0xc4, 0x2f, 0x2d, 0x55,
	0xd7, 0x34, 0x67, 0xe9, 0x18, 0x16, 0x2a, 0xd7, 0xd8, 0xb4, 0xd8, 0x68, 0x12, 0x74, 0xec, 0x87,
	0x57, 0x7f, 0xde, 0x2c, 0x83, 0x17, 0x37, 0xcb, 0xe0, 0xef, 0x9b, 0x65, 0xf0, 0xeb, 0xed, 0x72,
	0xf4, 0xe2, 0x76, 0x39, 0xfa, 0xeb, 0x76, 0x39, 0xfa, 0xf1, 0xa3, 0x42, 0x98, 0xed, 0x2e, 0x5f,
	0x73, 0x25, 0xdd, 0x07, 0xfd, 0x93, 0xff, 0xfd, 0xb6, 0xe7, 0x53, 0xf7, 0x81, 0xfe, 0xfc, 0xdf,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x34, 0x40, 0x75, 0xa2, 0x03, 0x06, 0x00, 0x00,
}

func (m *Chain) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Chain) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Chain) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintChains(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x4a
	}
	if m.CctxGateway != 0 {
		i = encodeVarintChains(dAtA, i, uint64(m.CctxGateway))
		i--
		dAtA[i] = 0x40
	}
	if m.IsExternal {
		i--
		if m.IsExternal {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x38
	}
	if m.Consensus != 0 {
		i = encodeVarintChains(dAtA, i, uint64(m.Consensus))
		i--
		dAtA[i] = 0x30
	}
	if m.Vm != 0 {
		i = encodeVarintChains(dAtA, i, uint64(m.Vm))
		i--
		dAtA[i] = 0x28
	}
	if m.NetworkType != 0 {
		i = encodeVarintChains(dAtA, i, uint64(m.NetworkType))
		i--
		dAtA[i] = 0x20
	}
	if m.Network != 0 {
		i = encodeVarintChains(dAtA, i, uint64(m.Network))
		i--
		dAtA[i] = 0x18
	}
	if m.ChainId != 0 {
		i = encodeVarintChains(dAtA, i, uint64(m.ChainId))
		i--
		dAtA[i] = 0x10
	}
	if m.ChainName != 0 {
		i = encodeVarintChains(dAtA, i, uint64(m.ChainName))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintChains(dAtA []byte, offset int, v uint64) int {
	offset -= sovChains(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Chain) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ChainName != 0 {
		n += 1 + sovChains(uint64(m.ChainName))
	}
	if m.ChainId != 0 {
		n += 1 + sovChains(uint64(m.ChainId))
	}
	if m.Network != 0 {
		n += 1 + sovChains(uint64(m.Network))
	}
	if m.NetworkType != 0 {
		n += 1 + sovChains(uint64(m.NetworkType))
	}
	if m.Vm != 0 {
		n += 1 + sovChains(uint64(m.Vm))
	}
	if m.Consensus != 0 {
		n += 1 + sovChains(uint64(m.Consensus))
	}
	if m.IsExternal {
		n += 2
	}
	if m.CctxGateway != 0 {
		n += 1 + sovChains(uint64(m.CctxGateway))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovChains(uint64(l))
	}
	return n
}

func sovChains(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozChains(x uint64) (n int) {
	return sovChains(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Chain) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowChains
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Chain: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Chain: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainName", wireType)
			}
			m.ChainName = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChainName |= ChainName(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			m.ChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChainId |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Network", wireType)
			}
			m.Network = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Network |= Network(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NetworkType", wireType)
			}
			m.NetworkType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NetworkType |= NetworkType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vm", wireType)
			}
			m.Vm = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Vm |= Vm(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Consensus", wireType)
			}
			m.Consensus = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Consensus |= Consensus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsExternal", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsExternal = bool(v != 0)
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CctxGateway", wireType)
			}
			m.CctxGateway = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CctxGateway |= CCTXGateway(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthChains
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthChains
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipChains(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthChains
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipChains(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowChains
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowChains
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowChains
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthChains
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupChains
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthChains
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthChains        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowChains          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupChains = fmt.Errorf("proto: unexpected end of group")
)
