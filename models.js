//@flow

// DO NOT EDIT -- automatically generated by goflow

// Accepted is used to describe commands accepted by the service.
// Note that Interrogate is always accepted.
export type Accepted = number

export type BpfHdr = &{C struct_bpf_hdr}

export type BpfInsn = &{C struct_bpf_insn}

export type BpfProgram = &{C struct_bpf_program}

export type BpfStat = &{C struct_bpf_stat}

export type BpfTimeval = &{C struct_bpf_timeval}

export type BpfVersion = &{C struct_bpf_version}

export type BpfZbuf = &{C struct_bpf_zbuf}

export type BpfZbufHeader = &{C struct_bpf_zbuf_header}

// Cmd represents service state change request. It is sent to a service
// by the service manager, and should be actioned upon by the service.
export type Cmd = number

export type Cmsghdr = &{C struct_cmsghdr}

export type Dirent = &{C struct_dirent}

export type EpollEvent = &{C struct_my_epoll_event}

// Errors should be an array of strings
export type Errors = Array<string>

export type Fbootstraptransfer_t = &{C struct_fbootstraptransfer}

export type FdSet = &{C fd_set}

// FieldMap allows users to customize the names of keys for various fields.
// As an example:
// formatter := &JSONFormatter{
//   	FieldMap: FieldMap{
// 		 FieldKeyTime: "@timestamp",
// 		 FieldKeyLevel: "@level",
// 		 FieldKeyLevel: "@message",
//    },
// }
export type FieldMap = { [key: fieldKey]: string }

// Fields type, used to pass to `WithFields`.
export type Fields = { [key: string]:  }

export type Flock_t = &{C struct_flock}

export type Fsid = &{C struct_fsid}

export type Fstore_t = &{C struct_fstore}

// Handle returns d's module handle.
export type Handle = numberptr

export type ICMPv6Filter = &{C struct_icmp6_filter}

export type IPMreq = &{C struct_ip_mreq}

export type IPMreqn = &{C struct_ip_mreqn}

export type IPv6MTUInfo = &{C struct_ip6_mtuinfo}

export type IPv6Mreq = &{C struct_ipv6_mreq}

export type IfAddrmsg = &{C struct_ifaddrmsg}

export type IfAnnounceMsghdr = &{C struct_if_announcemsghdr}

export type IfData = &{C struct_if_data}

export type IfInfomsg = &{C struct_ifinfomsg}

export type IfMsghdr = &{C struct_if_msghdr}

export type IfaMsghdr = &{C struct_ifa_msghdr}

export type IfmaMsghdr = &{C struct_ifma_msghdr}

export type IfmaMsghdr2 = &{C struct_ifma_msghdr2}

export type Inet4Pktinfo = &{C struct_in_pktinfo}

export type Inet6Pktinfo = &{C struct_in6_pktinfo}

export type InotifyEvent = &{C struct_inotify_event}

export type Iovec = &{C struct_iovec}

export type IpMaskString = IpAddressString

export type Kevent_t = &{C struct_kevent}

// Key is a handle to an open Windows registry key.
// Keys can be obtained by calling OpenKey; there are
// also some predefined root keys such as CURRENT_USER.
// Keys can be used directly in the Windows API.
export type Key = &{syscall Handle}

// Level type
export type Level = number

export type LevelHooks = { [key: Level]: Hook }

export type Linger = &{C struct_linger}

export type Log2phys_t = &{C struct_log2phys}

// MapKeyPtr is a string pointer key
export type MapKeyPtr = { [key: ?string]: Animal }

// MapKeyValPtr is a string pointer key
export type MapKeyValPtr = { [key: ?string]: ?Animal }

// MapNoPtr is a map of string to Animal, no pointer
export type MapNoPtr = { [key: string]: Animal }

// MapNumPtr should transform int64 to number
export type MapNumPtr = { [key: number]: Animal }

// MapValPtr is a string pointer value
export type MapValPtr = { [key: string]: ?Animal }

export type Mclpool = Array<number>

export type Msghdr = &{C struct_msghdr}

export type NlAttr = &{C struct_nlattr}

export type NlMsgerr = &{C struct_nlmsgerr}

export type NlMsghdr = &{C struct_nlmsghdr}

export type Note = string

// Payrate should be a number
export type Payrate = number

// People should be an array of Person
export type People = Array<Person>

export type PtraceRegs = &{C PtraceRegs}

export type Radvisory_t = &{C struct_radvisory}

export type RawSockaddr = &{C struct_sockaddr}

export type RawSockaddrAny = &{C struct_sockaddr_any}

export type RawSockaddrDatalink = &{C struct_sockaddr_dl}

export type RawSockaddrHCI = &{C struct_sockaddr_hci}

export type RawSockaddrInet4 = &{C struct_sockaddr_in}

export type RawSockaddrInet6 = &{C struct_sockaddr_in6}

export type RawSockaddrLinklayer = &{C struct_sockaddr_ll}

export type RawSockaddrNetlink = &{C struct_sockaddr_nl}

export type RawSockaddrUnix = &{C struct_sockaddr_un}

export type Rlimit = &{C struct_rlimit}

export type RtAttr = &{C struct_rtattr}

export type RtGenmsg = &{C struct_rtgenmsg}

export type RtMetrics = &{C struct_rt_metrics}

export type RtMsg = &{C struct_rtmsg}

export type RtMsghdr = &{C struct_rt_msghdr}

export type RtNexthop = &{C struct_rtnexthop}

export type Rusage = &{C struct_rusage}

// Signal table
export type Signal = number

export type SockFilter = &{C struct_sock_filter}

export type SockFprog = &{C struct_sock_fprog}

export type SockaddrGen = Array<number>

export type Stat_t = &{C struct_stat}

// State describes service execution state (Stopped, Running and so on).
export type State = number

export type Statfs_t = &{C struct_statfs}

// Strings should be an array of strings
export type Strings = Array<string>

export type Sysctlnode = &{C struct_sysctlnode}

export type Sysinfo_t = &{C struct_sysinfo}

export type TCPInfo = &{C struct_tcp_info}

export type Termio = &{C struct_termio}

export type Termios = &{C struct_termios}

export type Time_t = number

// Timespec is an invented structure on Windows, but here for
// consistency with the corresponding package for other operating systems.
export type Timespec = &{C struct_timespec}

export type Timeval = &{C struct_timeval}

export type Timeval32 = &{C struct_timeval32}

export type Timex = &{C struct_timex}

export type Tms = &{C struct_tms}

export type Token = Handle

export type Ucred = &{C struct_ucred}

export type Ustat_t = &{C struct_ustat}

export type Utimbuf = &{C struct_utimbuf}

export type Utsname = &{C struct_utsname}

export type WaitStatus = number

export type Winsize = &{C struct_winsize}

export type _C_int = number

export type _C_long = number

export type _C_long_long = number

export type _C_short = number

export type _Gid_t = number

export type _Socklen = number

export type errString = string

export type fieldKey = string

export type ifData = &{C struct_if_data}

export type ifMsghdr = &{C struct_if_msghdr}

export type ptraceFpregs = &{C ptraceFpregs}

export type ptracePer = &{C ptracePer}

export type ptracePsw = &{C ptracePsw}

export type syscallFunc = numberptr

// Animal is anything, but should probably have a master
// @strict
export type Animal = {|
	breed: string,
	name: string,
|}

export type EmbeddedAnimal = {
	breed: string,
	name: string,
	some_horse_attrib: string,
	doohickey: string,
	doohickey2: string,	// doohickey two
}

export type EmbeddedAnimal2 = {
	breed: string,
	name: string,
	birthday: string,	// birthday comment
	date: string,
	duration: string,	// a duration
	age: number,
}

export type Horse = {
	some_horse_attrib: string,
	doohickey: string,
	doohickey2: string,	// doohickey two
}

// Maps is for testing maps. These are the hardest part.
// The maps were not fun.
export type Maps = {
	base_map: { [key: string]: Person },
	base_map_ptr_key: { [key: ?string]: Person },
	base_map_ptr_val: { [key: string]: ?Person },
	map_of_slice: { [key: string]: Person },
	slice_of_map_of_slices: Array<Person>,
}

// NoIgnoredComment should NOT be ignored since flowignore is not the only
// thing there
// flowignore will not ignore here
export type NoIgnoredComment = {
	something: string,
}

// Parse tags after looking for embedded types.
// Ignore any json-ignore tags
export type Parse = {
	: ,
}

// Person has many types and should all convert correctly
export type Person = {
	name: string,	// This is a name comment
	age: number,
	StringOverride: String,	// Override `string` with `String`
	age64: number,
	flow_is_awesome: boolean,
	nullable:  ?string,
	animals_array: Array<Animal>,	// I have no pointer
	animals_array_ptr:  ?Array<Animal>,	// I am a pointer
	animals_array_ptr_2: Array<Animal>,	// I hold pointers
	payrate: Payrate,
	hascomma: string,
	some_generator: Generator,
	has_lots_of_tags: string,
	inner_struct: Object,	// I have a comment in a nested struct
	map_data: { [key: string]: number },
}

// TestFlowTags is to test all the possible flow flags
export type TestFlowTags = {
	person: Person,
	persona: Person,
	override_name_b: Person,	// should have new name
	personc: OverrideTypeA,	// should have original name but overriding type
	override_name_d: OverrideTypeB,
	override_name_f: Person,	// should have new name
}

// Time at which the log entry was created
export type Time = {
	the_time: string,
}

export type Whatever = {
	doohickey: string,
	doohickey2: string,	// doohickey two
}

export type Whatever2 = {
	doohickey2: string,	// doohickey two
}

