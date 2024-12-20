package main

import (
	"fmt"
	"github.com/sirikothe/gotextfsm"
)

const template = `Value InBytes (\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d*\s?\d+)
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

const input = `Flags: D - dynamic; X - disabled; I - inactive, R - running; S - slave; P - passthrough
 0   R   ;;; 123
         rx-byte=650 439 tx-byte=7 775 866 rx-packet=6 679 tx-packet=28 428 rx-drop=500 tx-drop=300 rx-error=0 tx-error=1000

 1       rx-byte=0 tx-byte=0 rx-packet=0 tx-packet=0 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

 2       rx-byte=0 tx-byte=0 rx-packet=0 tx-packet=0 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

 3       rx-byte=0 tx-byte=0 rx-packet=0 tx-packet=0 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

 4       rx-byte=0 tx-byte=0 rx-packet=0 tx-packet=0 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

 5       rx-byte=0 tx-byte=0 rx-packet=0 tx-packet=0 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

 6       rx-byte=0 tx-byte=0 rx-packet=0 tx-packet=0 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

 7       ;;; 1234
         rx-byte=0 tx-byte=0 rx-packet=0 tx-packet=0 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

 8   R   rx-byte=541 728 tx-byte=541 728 rx-packet=3 344 tx-packet=3 344 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

`


func main() {
	fms := gotextfsm.TextFSM{}
	err := fms.ParseString(template)

	if err != nil {
		fmt.Println(err)
	}

	parser := gotextfsm.ParserOutput{}
	err = parser.ParseTextString(input, fms, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(parser.Dict)
}
