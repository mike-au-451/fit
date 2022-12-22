package main

import (
	"fmt"
	"io"
	"os"

	"main/erf"
	"main/fit"
)

var argF string

func args() bool {
	argF = "-"
	return true
}

func main() {
	var (
		fh *os.File
		err error
	)

	if !args() {
		err =  erf.Errorf("bad args")
		fmt.Printf("%s\n", err)
		return
	}

	if argF == "-" {
		fh = os.Stdin
	} else {
		fh, err = os.Open(argF)
		if err != nil {
			err =  erf.Errorf("failed to open %s", argF)
			fmt.Printf("%s\n", err)
			return
		}
		defer fh.Close()
	}

	err = ReadFit(fh)
	if err != nil {
		err =  erf.Errorf("failed to read %s: %s", argF, err)
		fmt.Printf("%s\n", err)
	}
}

func ReadFit(fh *os.File) error {
	const (
		FIT_DEFINITION_MESSAGE = 0b01000000
		FIT_DATA_MESSAGE       = 0b00000000
	)

	var (
		// TODO: supporting type "Message"
		at, nn int
		hdr fit.FileHeader
		def fit.DefinitionMsg
		dat fit.DataMsg
		b1 [1]byte
		err error
	)

	hdr = fit.NewHeader()
	nn, err = hdr.ReadHeader(fh)
	if err != nil {
		return erf.Errorf("failed to read header: %s", err)
	}
	// TODO:
	// 1.  Check protocol and profile versions.
	// 2. Check crc
	hdr.Dump()

	for at < hdr.DataSize - 2 {
		// fmt.Printf("...top of loop at %d\n", at)

		nn, err = io.ReadFull(fh, b1[:])
		if err != nil {
			if err != io.EOF {
				return erf.Errorf("failed to read record header: %s", err)
			}
			break
		}
		// fmt.Printf("...record header: %d\n", b1[0])
		at += nn

		if (b1[0] & FIT_DEFINITION_MESSAGE) != 0 {
			// fmt.Printf("...definition message\n")
			def = fit.NewDefinitionMsg(int(b1[0]))
			nn, err = def.ReadMsg(fh)
			if err != nil {
				return erf.Errorf("failed to read definition record: %s", err)
			}
			at += nn
		} else {
			// fmt.Printf("...data message\n")
			dat = fit.NewDataMsg(int(b1[0]), &def)
			nn, err = dat.ReadMsg(fh)
			if err != nil {
				return erf.Errorf("failed to read definition record: %s", err)
			}
			at += nn
		}
		if err != nil {
			return erf.Errorf("failed to read record: %s", err)
		}

		switch b1[0] & 0b01000000 {
		case FIT_DEFINITION_MESSAGE:
			// fmt.Printf("...dump defn message\n")
			def.Dump()
		case FIT_DATA_MESSAGE:
			// fmt.Printf("...dump data message\n")
			dat.Dump()
		}
	}

	return nil
}