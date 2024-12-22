package chr

const hostnamTemplate = `Value Hostname (.*)

Start
 ^\s*name:\s${Hostname}
`

const versionTemplate = `Value Uptime (.*)
Value Version (.*)
Value BuildTime (\d+-\d+-\d+ \d+:\d+:\d+)
Value FactorySoftware (.*)
Value FreeMemory (\d*\.?\d*)
Value TotalMemory (\d*\.?\d*)
Value Cpu (\S*)
Value CpuCount (\d+)
Value CpuFrequency (\d+)
Value CpuLoad (\d+)
Value FreeHddSpace (\d*\.?\d*)
Value TotalHddSpace (\d*\.?\d*)
Value WriteSectSinceReboot (\d+)
Value WriteSecTotal (\d+)
Value ArchitectureName (.*)
Value BoardName (.*)

Start
	^\s*uptime: ${Uptime}
	^\s*version: ${Version}
	^\s*build-time: ${BuildTime}
	^\s*factory-software: ${FactorySoftware}
	^\s*free-memory: ${FreeMemory}MiB
	^\s*total-memory: ${TotalMemory}MiB
	^\s*cpu: ${Cpu}
	^\s*cpu-count: ${CpuCount}
	^\s*cpu-frequency: ${CpuFrequency}MHz
	^\s*cpu-load: ${CpuLoad}%
	^\s*free-hdd-space: ${FreeHddSpace}MiB
	^\s*total-hdd-space: ${TotalHddSpace}MiB
	^\s*write-sect-since-reboot: ${WriteSectSinceReboot}
	^\s*write-sect-total: ${WriteSecTotal}
	^\s*architecture-name: ${ArchitectureName}
	^\s*board-name: ${BoardName}`

const interfaceAboutTemplate = `Value Name (\S*)
Value Type (\S*)
Value NameOriginal (\S*)
Value MTU (\d+)
Value MAC (\w\w:\w\w:\w\w:\w\w:\w\w:\w\w)
Value Disabled (\w+)
Value Running (\w+)
Value Comment (\S*)

Start
 ^.*;;;\s${Comment} -> Continue
 ^\s*\d?\s+\w?\s+name="${Name}" -> Continue
 ^.*default-name="${NameOriginal}" -> Continue
 ^.*type="${Type}" -> Continue
 ^.*mtu=${MTU} -> Continue
 ^.*mac-address=${MAC} -> Continue
 ^.*disabled=${Disabled} -> Continue
 ^.*running=${Running} -> Next.Record
`

const interfaceCounterTemplate = `Value InBytes (\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d+)
Value OutBytes (\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d+)
Value InPkts (\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d+)
Value OutPkts (\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d+)
Value InDrops (\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d+)
Value OutDrops (\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d+)
Value ReadError (\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d+)
Value OutError (\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d+)

Start
 ^.*rx-byte=${InBytes} -> Continue
 ^.*tx-byte=${OutBytes} -> Continue
 ^.*rx-packet=${InPkts} -> Continue
 ^.*tx-packet=${OutPkts} -> Continue
 ^.*rx-drop=${InDrops} -> Continue
 ^.*tx-drop=${OutDrops} -> Continue
 ^.*rx-error=${ReadError} -> Continue
 ^.*tx-error=${OutError} -> Next.Record`
