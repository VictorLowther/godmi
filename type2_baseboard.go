/*
* File Name:	type2_baseboard.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-08-18 22:58:31
 */

package godmi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type BaseboardFeatureFlags byte

var baseboardFeatureFlags = []string{
	"Board is a hosting board", /* 0 */
	"Board requires at least one daughter board",
	"Board is removable",
	"Board is replaceable",
	"Board is hot swappable", /* 4 */
}

func (f BaseboardFeatureFlags) String() string {
	var s string
	for i := uint32(0); i < 5; i++ {
		if f&(1<<i) != 0 {
			s += "\n\t\t" + baseboardFeatureFlags[i]
		}
	}
	return s
}

func (b BaseboardFeatureFlags) toMap() map[string]bool {
	res := map[string]bool{}
	for i := range baseboardFeatureFlags {
		if b>>uint(i)&1 > 0 {
			res[baseboardFeatureFlags[i]] = true
		}
	}
	return res
}

func (b BaseboardFeatureFlags) MarshalJSON() ([]byte, error) {
	ref := b.toMap()
	return json.Marshal(&ref)
}

type BaseboardType byte

var baseboardType = []string{
	"Unspecified",
	"Unknown",
	"Other",
	"Server Blade",
	"Connectivity Switch",
	"System Management Module",
	"Processor Module",
	"I/O Module",
	"Memory Module",
	"Daughter Board",
	"Motherboard",
	"Processor+Memory Module",
	"Processor+I/O Module",
	"Interconnect Board", /* 0x0D */
}

func (b BaseboardType) String() string {
	if int(b) >= len(baseboardType) {
		b = 0
	}
	return baseboardType[b]
}

func (b BaseboardType) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

type BaseboardInformation struct {
	infoCommon
	Manufacturer                   string
	ProductName                    string
	Version                        string
	SerialNumber                   string
	AssetTag                       string
	FeatureFlags                   BaseboardFeatureFlags
	LocationInChassis              string
	ChassisHandle                  uint16
	BoardType                      BaseboardType
	NumberOfContainedObjectHandles byte
	ContainedObjectHandles         []byte
}

func (b BaseboardInformation) String() string {
	return fmt.Sprintf("Base Board Information\n"+
		"\tManufacturer: %s\n"+
		"\tProduct Name: %s\n"+
		"\tVersion: %s\n"+
		"\tSerial Number: %s\n"+
		"\tAsset Tag: %s\n"+
		"\tFeatures:%s\n"+
		"\tLocation In Chassis: %s\n"+
		"\tType: %s",
		b.Manufacturer,
		b.ProductName,
		b.Version,
		b.SerialNumber,
		b.AssetTag,
		b.FeatureFlags,
		b.LocationInChassis,
		b.BoardType)
}

var BaseboardInformations []*BaseboardInformation

func newBaseboardInformation(h dmiHeader) dmiTyper {
	data := h.data()
	length := int(data[0x01])
	bi := &BaseboardInformation{}
	for _, idx := range []int{0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0d} {
		if idx >= length {
			break
		}
		switch idx {
		case 0x04:
			bi.Manufacturer = h.FieldString(int(data[0x04]))
		case 0x05:
			bi.ProductName = h.FieldString(int(data[0x05]))
		case 0x06:
			bi.Version = h.FieldString(int(data[0x06]))
		case 0x07:
			bi.SerialNumber = h.FieldString(int(data[0x07]))
		case 0x08:
			bi.AssetTag = h.FieldString(int(data[0x08]))
		case 0x09:
			bi.FeatureFlags = BaseboardFeatureFlags(data[0x09])
		case 0x0a:
			bi.LocationInChassis = h.FieldString(int(data[0x0A]))
		case 0x0d:
			bi.BoardType = BaseboardType(data[0x0D])
		}
	}

	BaseboardInformations = append(BaseboardInformations, bi)
	return bi
}

func GetBaseboardInformation() string {
	var ret string
	for i, v := range BaseboardInformations {
		ret += "\n baseboard infomation index:" + strconv.Itoa(i) + "\n" + v.String()
	}
	return ret
}

func init() {
	addTypeFunc(SMBIOSStructureTypeBaseBoard, newBaseboardInformation)
}
