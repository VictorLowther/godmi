/*
* godmi.go
* godmi command
*
* Chapman Ou <ochapman.cn@gmail.com>
*
* Thu Jul 31 22:44:14 CST 2014
 */

package main

import (
	"fmt"

	"github.com/VictorLowther/godmi"
)

func main() {
	godmi.Init()
	fmt.Printf("%s\n", godmi.GetBIOSInformation())
	fmt.Printf("%s\n", godmi.GetSystemInformation())
	fmt.Printf("%s\n", godmi.GetBaseboardInformation())
	fmt.Printf("%s\n", godmi.GetChassisInformation())
	fmt.Printf("%s\n", godmi.GetProcessorInformation())
	fmt.Printf("%s\n", godmi.GetCacheInformation())
	fmt.Printf("%s\n", godmi.GetPortInformation())
	fmt.Printf("%s\n", godmi.GetSystemSlot())
	fmt.Printf("%s\n", godmi.GetOEMStrings())
	fmt.Printf("%s\n", godmi.GetSystemConfigurationOptions())
	fmt.Printf("%s\n", godmi.GetBIOSLanguageInformation())
	fmt.Printf("%s\n", godmi.GetGroupAssociations())
	fmt.Printf("%s\n", godmi.GetPhysicalMemoryArray())
	fmt.Printf("%s\n", godmi.GetMemoryDevice())
	fmt.Printf("%s\n", godmi.Get_32BitMemoryErrorInformation())
	fmt.Printf("%s\n", godmi.GetBuiltinPointingDevice())
	fmt.Printf("%s\n", godmi.GetPortableBattery())
	fmt.Printf("%s\n", godmi.GetSystemPowerSupplyInformation())
}
