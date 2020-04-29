/*
* godmi.go
* DMI SMBIOS information
*
* Chapman Ou <ochapman.cn@gmail.com>
*
 */
package godmi

import (
	"bytes"
	"fmt"
	"github.com/digitalocean/go-smbios/smbios"
	"io/ioutil"
	"strconv"
	"sync"
)

const OUT_OF_SPEC = "<OUT OF SPEC>"

type SMBIOSStructureType byte

const (
	SMBIOSStructureTypeBIOS SMBIOSStructureType = iota
	SMBIOSStructureTypeSystem
	SMBIOSStructureTypeBaseBoard
	SMBIOSStructureTypeChassis
	SMBIOSStructureTypeProcessor
	SMBIOSStructureTypeMemoryController
	SMBIOSStructureTypeMemoryModule
	SMBIOSStructureTypeCache
	SMBIOSStructureTypePortConnector
	SMBIOSStructureTypeSystemSlots
	SMBIOSStructureTypeOnBoardDevices
	SMBIOSStructureTypeOEMStrings
	SMBIOSStructureTypeSystemConfigurationOptions
	SMBIOSStructureTypeBIOSLanguage
	SMBIOSStructureTypeGroupAssociations
	SMBIOSStructureTypeSystemEventLog
	SMBIOSStructureTypePhysicalMemoryArray
	SMBIOSStructureTypeMemoryDevice
	SMBIOSStructureType32_bitMemoryError
	SMBIOSStructureTypeMemoryArrayMappedAddress
	SMBIOSStructureTypeMemoryDeviceMappedAddress
	SMBIOSStructureTypeBuilt_inPointingDevice
	SMBIOSStructureTypePortableBattery
	SMBIOSStructureTypeSystemReset
	SMBIOSStructureTypeHardwareSecurity
	SMBIOSStructureTypeSystemPowerControls
	SMBIOSStructureTypeVoltageProbe
	SMBIOSStructureTypeCoolingDevice
	SMBIOSStructureTypeTemperatureProbe
	SMBIOSStructureTypeElectricalCurrentProbe
	SMBIOSStructureTypeOut_of_bandRemoteAccess
	SMBIOSStructureTypeBootIntegrityServices
	SMBIOSStructureTypeSystemBoot
	SMBIOSStructureType64_bitMemoryError
	SMBIOSStructureTypeManagementDevice
	SMBIOSStructureTypeManagementDeviceComponent
	SMBIOSStructureTypeManagementDeviceThresholdData
	SMBIOSStructureTypeMemoryChannel
	SMBIOSStructureTypeIPMIDevice
	SMBIOSStructureTypePowerSupply
	SMBIOSStructureTypeAdditionalInformation
	SMBIOSStructureTypeOnBoardDevicesExtendedInformation
	SMBIOSStructureTypeManagementControllerHostInterface                     /*42*/
	SMBIOSStructureTypeInactive                          SMBIOSStructureType = 126
	SMBIOSStructureTypeEndOfTable                        SMBIOSStructureType = 127
)

func (b SMBIOSStructureType) String() string {
	types := [...]string{
		"BIOS", /* 0 */
		"System",
		"Base Board",
		"Chassis",
		"Processor",
		"Memory Controller",
		"Memory Module",
		"Cache",
		"Port Connector",
		"System Slots",
		"On Board Devices",
		"OEM Strings",
		"System Configuration Options",
		"BIOS Language",
		"Group Associations",
		"System Event Log",
		"Physical Memory Array",
		"Memory Device",
		"32-bit Memory Error",
		"Memory Array Mapped Address",
		"Memory Device Mapped Address",
		"Built-in Pointing Device",
		"Portable Battery",
		"System Reset",
		"Hardware Security",
		"System Power Controls",
		"Voltage Probe",
		"Cooling Device",
		"Temperature Probe",
		"Electrical Current Probe",
		"Out-of-band Remote Access",
		"Boot Integrity Services",
		"System Boot",
		"64-bit Memory Error",
		"Management Device",
		"Management Device Component",
		"Management Device Threshold Data",
		"Memory Channel",
		"IPMI Device",
		"Power Supply",
		"Additional Information",
		"Onboard Device",
		"Management Controller Host Interface", /* 42 */
	}

	if int(b) >= len(types) {
		return "unspported type:" + strconv.Itoa(int(b))
	}
	return types[b]
}

type SMBIOSStructureHandle uint16

type infoCommon struct {
	smType SMBIOSStructureType
	length byte
	handle SMBIOSStructureHandle
}

type dmiHeader []byte

func (h dmiHeader) smType() SMBIOSStructureType {
	return SMBIOSStructureType(h[0x00])
}

func (h dmiHeader) len() int {
	return int(h[0x01])
}

func (h dmiHeader) handle() SMBIOSStructureHandle {
	return SMBIOSStructureHandle(u16(h[0x02:0x04]))
}

func (h dmiHeader) end() int {
	start := h.len()
	end := bytes.Index(h[start:], []byte{0, 0})
	if end == -1 {
		return -1
	}
	return start + end
}

func (h dmiHeader) data() []byte {
	res := make([]byte, 256)
	copy(res, h[:h.len()])
	return res
}

func (h dmiHeader) FieldString(idx int) string {
	end := h.end()
	if end == -1 || idx == 0 {
		return ""
	}
	bs := bytes.Split(h[h.len():end], []byte{0})
	if idx > len(bs) {
		return fmt.Sprintf("FieldString ### ERROR:strFields Len:%d, strIndex:%d", len(bs), idx)
	}
	return string(bs[idx-1])
}

func newdmiHeader(d []byte) dmiHeader {
	if len(d) < 0x04 {
		return nil
	}
	return dmiHeader(d)
}

func (h dmiHeader) Next() dmiHeader {
	index := h.end()

	if index == -1 {
		return nil
	}
	return newdmiHeader(h[index+2:])
}

func (h dmiHeader) decode() error {
	newfn, err := getTypeFunc(h.smType())
	if err != nil {
		return err
	}
	newfn(h)
	return nil
}

type dmiTyper interface {
	String() string
}

type newFunction func(d dmiHeader) dmiTyper

type typeFunc map[SMBIOSStructureType]newFunction

var g_typeFunc = make(typeFunc)

var g_lock sync.Mutex

var smbiosVersion string

func addTypeFunc(t SMBIOSStructureType, f newFunction) {
	g_lock.Lock()
	defer g_lock.Unlock()
	g_typeFunc[t] = f
}

func getTypeFunc(t SMBIOSStructureType) (fn newFunction, err error) {
	fn, ok := g_typeFunc[t]
	if !ok {
		return fn, fmt.Errorf("type %d have no NewFunction", int(t))
	}
	return fn, nil
}

func Init() error {
	stream, ep, err := smbios.Stream()
	if err != nil {
		return err
	}
	defer stream.Close()
	maj, min, _ := ep.Version()
	smbiosVersion = fmt.Sprintf("%d.%d", maj, min)
	tmem, err := ioutil.ReadAll(stream)
	if err != nil {
		return err
	}
	for hd := newdmiHeader(tmem); hd != nil; hd = hd.Next() {
		err := hd.decode()
		if err != nil {
			//fmt.Println("info: ", err)
			continue
		}
	}
	return nil
}
