package goicns

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"

	"github.com/nfnt/resize"
)

const (
	icp4 = 16  // 16x16
	icp5 = 32  // 32x32
	icp6 = 64  // 64x64
	ic07 = 128 // 128x128
	ic08 = 256 // 256x256, 10.5
	ic09 = 512 // 512x512, 10.5

	ic11 = 32   // 16x16@2x retina
	ic12 = 64   // 32x32@2x retina
	ic13 = 256  // 128x128@2x retina
	ic14 = 512  // 256x256@2x retina
	ic10 = 1024 // 1024x1024/512x512@2x
)

var (
	IcnsHeader = []byte{0x69, 0x63, 0x6e, 0x73}
	// All icon sizes
	sizes = []int{16, 32, 64, 128, 256, 512, 1024}

	// All icon sizes mapped to their respective possible OSTypes.
	// This includes old OSTypes such as ic08 recognized on 10.5.
	// TODO: Would only one OSType suffice?
	sizeToType = map[int][]string{
		16:   {"icp4"},
		32:   {"icp5", "ic11"},
		64:   {"icp6", "ic12"},
		128:  {"ic07"},
		256:  {"ic08", "ic13"},
		512:  {"ic09", "ic14"},
		1024: {"ic10"},
	}
)

func NewICNS(img image.Image) *ICNS {
	return &ICNS{
		BaseImage: img,
	}
}

type ICNS struct {
	BaseImage image.Image
	Data      *[]byte
}

func (i ICNS) WriteTo(w io.Writer) (n int, err error) {
	return w.Write(*i.Data)
}

func (i ICNS) WriteToFile(path string, perm os.FileMode) error {
	return ioutil.WriteFile(path, *i.Data, perm)
}

func (i *ICNS) Construct() (err error) {
	// Create a new buffer to hold the series of icons generated via resizing
	icns := new(bytes.Buffer)

	for _, s := range sizes {
		imgBuf := new(bytes.Buffer)
		// Resize the BaseImage to each specific size
		resized := resize.Resize(uint(s), uint(s), i.BaseImage, resize.MitchellNetravali)
		if err = png.Encode(imgBuf, resized); err != nil {
			return
		}

		// Each icon type is prefixed with a 4-byte OSType marker and a 4-byte size
		// header (which includes the ostype/size header). Add the size of the total
		// icon to lenByt in big-endian format.
		lenByt := make([]byte, 4, 4)
		binary.BigEndian.PutUint32(lenByt, uint32(imgBuf.Len()+8))

		// Iterate through every OSType and append the icon to icns
		for _, ostype := range sizeToType[s] {
			if _, err = icns.Write([]byte(ostype)); err != nil {
				return
			}
			if _, err = icns.Write(lenByt); err != nil {
				return
			}
			if _, err = icns.Write(imgBuf.Bytes()); err != nil {
				return
			}
		}
	}

	// Each ICNS file is prefixed with a 4 byte header and 4 bytes marking
	// the length of the file, MSB first.
	lenByt := make([]byte, 4, 4)
	binary.BigEndian.PutUint32(lenByt, uint32(icns.Len()+8))

	data := IcnsHeader
	data = append(data, lenByt...)
	data = append(data, icns.Bytes()...)

	i.Data = &data

	return
}
