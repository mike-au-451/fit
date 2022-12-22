package fit

import (
	"fmt"
	"io"
	"os"
	"strings"

	"main/erf"
	"main/util"
)

type File struct {
	Hdr FileHeader
	rcs []DataMsg
}

type FileHeader struct {
	size, protocolver, profilever, DataSize, crc int
	datatype string
}

type DataMsg struct {
	header int
	defn *DefinitionMsg
	data []interface{}	// this is kinda complicated...,
}

type DefinitionMsg struct {
	header, reserved, architecture, globalmsgno, fitfieldcnt, devfieldcnt int
	FitDefns, DevDefns []FieldDef
}

type FieldDef struct {
	defno, size, basetype int
}

func NewHeader() (FileHeader) {
	// fmt.Printf(">>>NewHeader\n")
	return FileHeader{}
}

func NewDefinitionMsg(hd int) (DefinitionMsg) {
	// fmt.Printf(">>>NewDefinitionMsg\n")
	msg := DefinitionMsg{}
	msg.header = hd
	return msg
}

func NewDataMsg(hd int, defn *DefinitionMsg) (DataMsg) {
	// fmt.Printf(">>>NewDataMsg\n")
	msg := DataMsg{}
	msg.header = hd
	msg.defn = defn

	return msg
}

// TODO: crap naming
func (hdr *FileHeader) ReadHeader(fh *os.File) (int, error) {

	// fmt.Printf(">>>ReadHeader\n")

	var (
		at, nn int
		b1 [1]byte
		b2 [2]byte
		b4 [4]byte

		err error
	)

	nn, err = io.ReadFull(fh, b1[:])
	if err != nil {
		return at, erf.Errorf("failed to read header byte: %s", err)
	}
	at += nn
	hdr.size = int(b1[0])

	switch b1[0] {
	case 12, 14:
	default:
		return at, erf.Errorf("unexpected header size: %d", b1[0])
	}

	nn, err = io.ReadFull(fh, b1[:])
	if err != nil {
		return at, erf.Errorf("failed to read protocol version: %s", err)
	}
	at += nn
	hdr.protocolver = int(b1[0])

	nn, err = io.ReadFull(fh, b2[:])
	if err != nil {
		return at, erf.Errorf("failed to read profile version: %s", err)
	}
	at += nn
	hdr.profilever = int(b2[1]) << 8 + int(b2[0])

	nn, err = io.ReadFull(fh, b4[:])
	if err != nil {
		return at, erf.Errorf("failed to read data size: %s", err)
	}
	at += nn
	hdr.DataSize = int(b4[3]) << 24 + int(b4[2]) << 16 + int(b4[1]) << 8 + int(b4[0])

	nn, err = io.ReadFull(fh, b4[:])
	if err != nil {
		return at, erf.Errorf("failed to read file type: %s", err)
	}
	at += nn
	hdr.datatype = string(b4[:])

	nn, err = io.ReadFull(fh, b2[:])
	if err != nil {
		return at, erf.Errorf("failed to read file crc: %s", err)
	}
	at += nn
	hdr.crc = int(b2[1]) << 8 + int(b2[0])

	return at, nil
}

func (msg *DefinitionMsg) ReadMsg(fh *os.File) (int, error) {

	// fmt.Printf(">>>DefinitionMsg.ReadMsg\n")

	var (
		at, nn int
		b1 [1]byte
		b2 [2]byte
		err error
	)

	nn, err = io.ReadFull(fh, b1[:])
	if err != nil {
		return at, erf.Errorf("failed to read defn reserved byte: %s", err)
	}
	at += nn
	msg.reserved = int(b1[0])

	nn, err = io.ReadFull(fh, b1[:])
	if err != nil {
		return at, erf.Errorf("failed to read defn architecture byte: %s", err)
	}
	at += nn
	msg.architecture = int(b1[0])

	nn, err = io.ReadFull(fh, b2[:])
	if err != nil {
		return at, erf.Errorf("failed to read defn global msg no: %s", err)
	}
	at += nn
	msg.globalmsgno = util.Int(b2[:], int(msg.architecture))

	nn, err = io.ReadFull(fh, b1[:])
	if err != nil {
		return at, erf.Errorf("failed to read data field count: %s", err)
	}
	at += nn
	msg.fitfieldcnt = int(b1[0])

	msg.FitDefns, nn, err = ReadDataDefFields(fh, msg.fitfieldcnt)
	if err != nil {
		return at, erf.Errorf("failed to read field definitions: %s", err)
	}
	at += nn

	const FIT_EXTENDED_DATA = 0b00100000
	if msg.header & FIT_EXTENDED_DATA != 0 {
		// extended data
		nn, err = io.ReadFull(fh, b1[:])
		if err != nil {
			return at, erf.Errorf("failed to read extended field count: %s", err)
		}
		at += nn
		msg.devfieldcnt = int(b1[0])

		msg.DevDefns, nn, err = ReadDataDefFields(fh, msg.devfieldcnt)
		if err != nil {
			return at, erf.Errorf("failed to read extended definitions: %s", err)
		}
		at += nn
	}

	return at, nil
}

func (msg *DataMsg) ReadMsg(fh *os.File) (int, error) {
	// fmt.Printf(">>>DataMsg.ReadMsg\n")

	var (
		at, nn int
		fielddefs []FieldDef
		buf []byte
		err error
	)

	fielddefs = append(fielddefs, msg.defn.FitDefns...)
	fielddefs = append(fielddefs, msg.defn.DevDefns...)

	for ii := 0; ii < len(fielddefs); ii++ {
		buf = make([]byte, fielddefs[ii].size)
		nn, err = io.ReadFull(fh, buf)
		if err != nil {
			return at, erf.Errorf("failed to read actual data: %s", err)
		}
		at += nn

		// TODO: rationialize variant data
		switch fielddefs[ii].basetype {
		case 0x00, 0x01, 0x02, 0x0d:/// enum, sint8, uint8, byte
			msg.data = append(msg.data, util.Int8(buf, msg.defn.architecture))
		case 0x07:// string
			msg.data = append(msg.data, string(buf))
		case 0x83, 0x84:// sint16, uint16
			msg.data = append(msg.data, util.Int16(buf, msg.defn.architecture))
		case 0x85, 0x86:// sint32, uint32
			msg.data = append(msg.data, util.Int32(buf, msg.defn.architecture))
		case 0x89:// float64
			panic("unhandled float conversion")
		case 0x8c:// uint32z
		default:
			panic(fmt.Sprintf("bug: mssing basetype: %d", fielddefs[ii].basetype))
		}

	}

	return at, nil
}

func ReadDataDefFields(fh *os.File, cnt int) ([]FieldDef, int, error) {
	// fmt.Printf(">>>DefinitionMsg.ReadDataDefFields(%d)\n", cnt)

	var (
		at, nn int
		defs []FieldDef
		b3 [3]byte
		err error
	)

	for ii := 0; ii < cnt; ii++ {
		nn, err = io.ReadFull(fh, b3[:])
		if err != nil {
			return nil, at, erf.Errorf("failed to read fields: %s", err)
		}
		at += nn

		defs = append(defs, FieldDef{ int(b3[0]), int(b3[1]), int(b3[2]) })
	}

	return defs, at, nil
}

func (hd *FileHeader) Dump() {
	fmt.Printf("FileHeader:\n%+v\n", *hd)
}

func (msg *DefinitionMsg) Dump() {
	archs := map[int]string{
		0: "little endian",
		1: "little endian",
	}
	// from garmin xls
	gmns := map[int]string{
		0x00: "file_id",
		0x08: "hr_zone",
		0x09: "power_zone",
		0x0c: "sport",
		0x12: "session",
		0x13: "lap",
		0x14: "record",
		0x15: "event",
		0x17: "device_info",
		0x1a: "workout",
		0x22: "activity",
		0xce: "field_description",
		0xcf: "developer_data_id",

		0xff00: "reserved - manufacturer specific message",
		0xff01: "reserved - manufacturer specific message",
		0xff04: "reserved - manufacturer specific message",
	}

	bts := map[int]string{
		0x00: "enum",
		0x01: "sint8",
		0x02: "uint8",
		0x07: "string",
		0x0d: "byte",
		0x83: "sint16",
		0x84: "uint16",
		0x85: "sint32",
		0x86: "uint32",
		0x89: "float64",
		0x8c: "uint32z",
	}

	fmt.Printf("DefinitionMsg @ %p\n", msg)
	fmt.Printf("\theader: %d\n", msg.header)
	fmt.Printf("\treserved: %d\n", msg.reserved)
	fmt.Printf("\tarchitecture: %d (%s)\n", msg.architecture, archs[msg.architecture])
	gmn, ok := gmns[msg.globalmsgno]
	if !ok {
		// panic(fmt.Sprintf("missing global message number: %d", msg.globalmsgno))
		fmt.Printf("missing global message number: %d", msg.globalmsgno)
	}
	fmt.Printf("\tglobalmsgno: %d (%s)\n", msg.globalmsgno, gmn)
	fmt.Printf("\t(fitfieldcnt, len): (%d, %d)\n", msg.fitfieldcnt, len(msg.FitDefns))
	if msg.fitfieldcnt != len(msg.FitDefns) {
		panic("bug")
	}
	fmt.Printf("\t(devfieldcnt, len): (%d, %d)\n", msg.devfieldcnt, len(msg.DevDefns))
	if msg.devfieldcnt != len(msg.DevDefns) {
		panic("bug")
	}

	defns := []string{}
	for ii := 0; ii < len(msg.FitDefns); ii++ {
		defns = append(defns, fmt.Sprintf("{ %d, %d, %s }", msg.FitDefns[ii].defno, msg.FitDefns[ii].size, bts[msg.FitDefns[ii].basetype]))
	}
	for ii := 0; ii < len(msg.DevDefns); ii++ {
		defns = append(defns, fmt.Sprintf("{ %d, %d, %s }", msg.DevDefns[ii].defno, msg.DevDefns[ii].size, bts[msg.DevDefns[ii].basetype]))
	}
	fmt.Printf("\tFieldDefs: %s\n", strings.Join(defns, " "))
}

func (msg *DataMsg) Dump() {
	fmt.Printf("DataMsg:\n%+v\n", *msg)
}

func (fd *FieldDef) Dump() {
	bts := map[int]string{
		0x00: "enum",
		0x01: "sint8",
		0x02: "uint8",
		0x07: "string",
		0x0d: "byte",
		0x83: "sint16",
		0x84: "uint16",
		0x85: "sint32",
		0x86: "uint32",
		0x89: "float64",
		0x8c: "uint32z",
	}
	fmt.Printf("FieldDef @ %p\n", fd)
	bt, ok := bts[fd.basetype]
	if !ok {
		panic(fmt.Sprintf("missing bt: %x", fd.basetype))
	}
	fmt.Printf("\t{ %d, %d, %s }\n", fd.defno, fd.size, bt)
}

