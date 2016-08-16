package adapter

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

func om_cluster_docs_replica_set_json() ([]byte, error) {
	return bindata_read([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0xcc, 0x57,
		0x4d, 0x73, 0xa3, 0x38, 0x10, 0xbd, 0xfb, 0x57, 0xb8, 0xd8, 0xab, 0x3f,
		0x88, 0xe3, 0x64, 0x93, 0x9c, 0xc6, 0x93, 0x1c, 0xe6, 0xb0, 0xd9, 0x4c,
		0x25, 0x33, 0x93, 0xc3, 0x4c, 0xca, 0x25, 0x90, 0x30, 0xda, 0x11, 0x12,
		0x25, 0x89, 0xc4, 0x0e, 0xc5, 0x7f, 0x5f, 0x49, 0x60, 0x8c, 0x40, 0x38,
		0x9e, 0xaa, 0xf1, 0x66, 0x39, 0x99, 0x56, 0xab, 0xbb, 0xf5, 0xf4, 0x5e,
		0xd3, 0xce, 0x07, 0x43, 0xf5, 0x78, 0x2c, 0x95, 0x98, 0x51, 0xe1, 0x5d,
		0x0d, 0x4b, 0x83, 0x31, 0x42, 0xf6, 0x42, 0x09, 0x03, 0xf0, 0x23, 0x10,
		0x48, 0xad, 0x78, 0xd3, 0x67, 0xc0, 0xa7, 0x04, 0x07, 0xd3, 0x84, 0xd1,
		0x15, 0x83, 0xc1, 0x38, 0x49, 0xc4, 0x18, 0x64, 0x92, 0x25, 0x40, 0xef,
		0xf6, 0x46, 0xee, 0xad, 0x8f, 0x98, 0xaa, 0x57, 0x1d, 0xdb, 0xbb, 0xbe,
		0xfa, 0xf1, 0xa3, 0x67, 0xb3, 0xd9, 0x5b, 0x94, 0x21, 0x3c, 0xe3, 0x73,
		0x13, 0x7c, 0x43, 0x5c, 0x54, 0x65, 0x7d, 0x6f, 0xd4, 0x15, 0x64, 0x98,
		0x40, 0x63, 0xac, 0x6d, 0xdb, 0x27, 0xef, 0x58, 0xca, 0x1d, 0x58, 0x6a,
		0xff, 0xf3, 0xf9, 0xc8, 0xbd, 0x1e, 0x11, 0xf0, 0xcc, 0xb8, 0xae, 0xd0,
		0xeb, 0xf1, 0x58, 0x61, 0x59, 0x55, 0xa3, 0xbd, 0xe6, 0xb3, 0xf9, 0x65,
		0x78, 0x02, 0x67, 0xc1, 0xd9, 0xe5, 0xe5, 0x25, 0x0a, 0x82, 0xe8, 0x24,
		0x82, 0xd1, 0x69, 0x10, 0xfa, 0xc8, 0x47, 0xa7, 0xc1, 0x69, 0x14, 0x9d,
		0x85, 0x3e, 0x00, 0xd1, 0xac, 0x2f, 0x5a, 0x02, 0xd6, 0x77, 0xa2, 0x11,
		0xaf, 0xd7, 0x0f, 0xd3, 0xc3, 0xfc, 0x18, 0xcc, 0x08, 0x32, 0x90, 0x3c,
		0xf5, 0xb8, 0xa4, 0x04, 0xc8, 0x88, 0xf1, 0x44, 0x87, 0x61, 0x62, 0xdd,
		0x17, 0x29, 0xe3, 0x44, 0x7b, 0xc4, 0x52, 0xa6, 0xe2, 0x6a, 0x3a, 0x8d,
		0x80, 0x90, 0x90, 0x4c, 0xaa, 0x3b, 0x9b, 0x30, 0xbe, 0x52, 0x04, 0xa0,
		0xd9, 0xba, 0xa6, 0x80, 0x79, 0x1b, 0xaf, 0x2f, 0xce, 0x97, 0xe7, 0xf3,
		0x31, 0x48, 0xc0, 0x2b, 0xa3, 0xe3, 0xd3, 0xc9, 0x6c, 0xf2, 0xe7, 0x44,
		0xae, 0x5e, 0xfb, 0x92, 0xbc, 0x60, 0x3a, 0xf3, 0xfd, 0x8b, 0x94, 0x64,
		0xba, 0xe4, 0x08, 0x10, 0x81, 0xfa, 0x3d, 0xbf, 0x5d, 0xdf, 0x23, 0x88,
		0x85, 0xbc, 0x21, 0x64, 0x1f, 0x04, 0x0d, 0xcf, 0xbb, 0x9a, 0xcc, 0xbd,
		0x68, 0x34, 0xbc, 0xbf, 0xf2, 0x43, 0xe3, 0x36, 0x2f, 0xa2, 0xe3, 0x5c,
		0x58, 0x96, 0x46, 0x5e, 0x8f, 0x82, 0xc4, 0xc8, 0xc7, 0xc0, 0x52, 0xf1,
		0xbc, 0x5a, 0xf7, 0x02, 0x10, 0xfe, 0xcc, 0xd2, 0x1e, 0x9e, 0xc7, 0x4c,
		0xc8, 0xed, 0xe6, 0x3c, 0xa7, 0x0c, 0x22, 0x31, 0xf9, 0xee, 0x3f, 0x15,
		0x45, 0x53, 0x69, 0x84, 0xad, 0x3e, 0x03, 0x19, 0xef, 0xf4, 0xc9, 0x56,
		0x3d, 0xfa, 0x9c, 0x96, 0xc9, 0xc6, 0x60, 0x85, 0xa8, 0x9c, 0x28, 0xbf,
		0x56, 0x98, 0x7b, 0x26, 0x81, 0x44, 0x56, 0x0b, 0x30, 0x4b, 0x02, 0xbf,
		0xa2, 0x2f, 0x31, 0x47, 0x22, 0x66, 0x04, 0xde, 0x7e, 0x54, 0x0e, 0x27,
		0xbe, 0xef, 0xdb, 0x70, 0x79, 0x12, 0x27, 0x3b, 0xa7, 0x4f, 0x5c, 0x1f,
		0x65, 0x36, 0x1f, 0xd8, 0xe0, 0xe8, 0x53, 0xd7, 0xfa, 0xc6, 0x92, 0x71,
		0x4c, 0x57, 0xff, 0xd1, 0xd1, 0x77, 0x09, 0xdf, 0xf7, 0xf8, 0xc6, 0x39,
		0xe5, 0x2c, 0x44, 0x42, 0x94, 0x7a, 0xcd, 0xf3, 0x3f, 0x10, 0x08, 0xe3,
		0xa1, 0x39, 0x63, 0xb1, 0x65, 0x51, 0x03, 0x0a, 0xc0, 0x57, 0x62, 0xb6,
		0x3c, 0xef, 0x16, 0x46, 0x91, 0xec, 0x18, 0xcb, 0xf8, 0x8c, 0xeb, 0x95,
		0xd9, 0x85, 0xaa, 0xd3, 0x5a, 0x2d, 0x5a, 0x55, 0x73, 0x94, 0x12, 0x1c,
		0x96, 0xfd, 0xd7, 0x19, 0x49, 0x3b, 0x3c, 0x20, 0xf9, 0x77, 0x75, 0x13,
		0x69, 0x18, 0x2d, 0xb5, 0xc9, 0xdb, 0x1b, 0x55, 0x28, 0xa0, 0x15, 0xc8,
		0xee, 0x88, 0x30, 0xb0, 0xae, 0xec, 0x39, 0x04, 0xe9, 0x54, 0x6f, 0x40,
		0xf5, 0xcd, 0x41, 0x20, 0xc1, 0x1b, 0x09, 0x36, 0x42, 0xa2, 0xe4, 0x2f,
		0x75, 0x83, 0xee, 0x14, 0x48, 0x48, 0x4c, 0xb7, 0xa7, 0xf2, 0x22, 0x4c,
		0x90, 0x43, 0xdc, 0x5e, 0xda, 0xad, 0x63, 0x23, 0x1a, 0x1c, 0x5a, 0xea,
		0x1b, 0xd9, 0x56, 0x65, 0xf8, 0x62, 0x17, 0x35, 0x70, 0x94, 0xd7, 0xe2,
		0xad, 0x8c, 0xb1, 0xe8, 0x50, 0xf6, 0x98, 0x3c, 0x73, 0x74, 0x1d, 0x57,
		0x15, 0x15, 0x01, 0xbf, 0x6c, 0x52, 0xe3, 0x52, 0x9e, 0xb1, 0xe9, 0xf0,
		0xbc, 0xeb, 0x74, 0x65, 0xd7, 0x6a, 0xac, 0x29, 0x55, 0xc5, 0x0f, 0x61,
		0x8c, 0x12, 0xb0, 0xeb, 0x87, 0x67, 0x25, 0xc1, 0x15, 0x97, 0x71, 0x34,
		0xfc, 0x40, 0xd4, 0x47, 0xa3, 0x50, 0x2f, 0x48, 0xf5, 0xf5, 0xa2, 0x18,
		0xe5, 0xf9, 0x14, 0x47, 0x15, 0xb3, 0xd5, 0x6f, 0xcd, 0x76, 0xf3, 0xb6,
		0x95, 0x43, 0xc5, 0x42, 0xc5, 0xb3, 0x76, 0x17, 0x58, 0x62, 0x68, 0xd1,
		0xae, 0x51, 0x45, 0x82, 0x92, 0x00, 0xf1, 0xf6, 0x10, 0xe0, 0x16, 0x93,
		0x59, 0xb1, 0x71, 0x2c, 0x23, 0xe7, 0xf9, 0x07, 0x35, 0x9b, 0xa0, 0x75,
		0xd1, 0x66, 0x98, 0xbe, 0x46, 0x0b, 0x3c, 0x3b, 0xc7, 0xee, 0x90, 0x23,
		0x25, 0x4f, 0x35, 0x5c, 0x20, 0x7e, 0x47, 0xc9, 0x46, 0xed, 0x90, 0x3c,
		0x43, 0x23, 0x85, 0x2f, 0x56, 0xbd, 0x46, 0x6a, 0x83, 0x6f, 0x15, 0x51,
		0x41, 0x52, 0x63, 0xd2, 0x84, 0xa3, 0x7c, 0x9e, 0xec, 0x56, 0xc1, 0x99,
		0xfd, 0x59, 0xf7, 0x44, 0x0c, 0x38, 0x54, 0x6d, 0xac, 0xb2, 0x0d, 0xea,
		0x2b, 0xb1, 0x47, 0x37, 0xdd, 0xfa, 0xbe, 0x0a, 0x64, 0x46, 0x9a, 0xfe,
		0x29, 0x4d, 0x5b, 0x3f, 0xbf, 0xc0, 0xf2, 0xa4, 0xc3, 0x14, 0x08, 0xf1,
		0xc2, 0x38, 0x1c, 0xda, 0x64, 0x81, 0x0a, 0x79, 0xb6, 0x49, 0x54, 0xd3,
		0x5c, 0xa8, 0x34, 0xb7, 0x28, 0x8c, 0x01, 0xc5, 0x22, 0xe9, 0xce, 0x5f,
		0xde, 0xc3, 0xf5, 0xfd, 0xe2, 0x76, 0xfc, 0xf0, 0x69, 0x31, 0x3e, 0xd9,
		0x01, 0xd6, 0xfc, 0x12, 0xfe, 0x44, 0x9b, 0x2a, 0x97, 0xfa, 0xd5, 0x4a,
		0xa3, 0x2c, 0x46, 0xa9, 0x96, 0x22, 0xff, 0x61, 0x81, 0xb0, 0xe4, 0x18,
		0x32, 0x1a, 0xe1, 0x4a, 0xa1, 0x4b, 0x96, 0x4c, 0x74, 0xc4, 0x66, 0xad,
		0x58, 0x80, 0x80, 0x20, 0xd8, 0x9d, 0x2a, 0xbc, 0x4c, 0xa1, 0x21, 0x6e,
		0x10, 0x41, 0xd2, 0x2c, 0x37, 0x27, 0x83, 0x72, 0xed, 0x11, 0xd0, 0x6a,
		0xc9, 0x3a, 0x95, 0xb3, 0x87, 0xe9, 0x2a, 0x01, 0x54, 0xd3, 0x99, 0xab,
		0xb1, 0xd4, 0x77, 0xe6, 0x1c, 0x28, 0xdc, 0x13, 0xea, 0x01, 0x71, 0xad,
		0xf8, 0xda, 0x31, 0x54, 0xe3, 0x93, 0x62, 0xde, 0x6d, 0xf9, 0x69, 0xeb,
		0x4e, 0x24, 0xfa, 0x29, 0x3a, 0x56, 0xc7, 0x44, 0x64, 0xce, 0xbf, 0x65,
		0x4a, 0xfb, 0x4b, 0xe9, 0x3a, 0x21, 0x56, 0x2e, 0x6e, 0xde, 0xec, 0xeb,
		0xdb, 0xff, 0x7b, 0x24, 0x17, 0xc6, 0xdf, 0x8d, 0xa3, 0x3b, 0xca, 0xef,
		0x2b, 0x81, 0x23, 0x00, 0x17, 0x74, 0x73, 0xa3, 0x3e, 0x7e, 0x81, 0xfe,
		0xbb, 0xf5, 0x3e, 0x55, 0x68, 0x22, 0x18, 0x14, 0x8e, 0x56, 0x0a, 0x61,
		0x21, 0x20, 0x87, 0x02, 0xf2, 0xa8, 0x9a, 0xe8, 0x7b, 0x41, 0xf1, 0x56,
		0xfe, 0x5f, 0x16, 0x56, 0x73, 0xfa, 0x3e, 0xa2, 0xa8, 0xac, 0x73, 0x0e,
		0x3b, 0x79, 0xea, 0x92, 0x7a, 0x1d, 0xf6, 0x88, 0xae, 0x07, 0xdd, 0x43,
		0xa0, 0x3d, 0x54, 0x6b, 0xae, 0x7b, 0xfd, 0x1d, 0x69, 0xdf, 0xd6, 0xd7,
		0xb1, 0x32, 0x1f, 0xa8, 0xa9, 0x5f, 0x4f, 0xbf, 0x57, 0x4a, 0x07, 0xf1,
		0xf8, 0x98, 0x60, 0xf7, 0xe6, 0x6c, 0x9b, 0xba, 0xba, 0x69, 0x49, 0xc1,
		0xe4, 0x5c, 0xf6, 0x0b, 0xc2, 0x39, 0x6a, 0xe8, 0xd1, 0xc6, 0x9a, 0x57,
		0x74, 0xb0, 0xce, 0x80, 0x52, 0x0c, 0x8a, 0xc1, 0xbf, 0x01, 0x00, 0x00,
		0xff, 0xff, 0xf4, 0x31, 0xf6, 0xbb, 0x02, 0x13, 0x00, 0x00,
	},
		"om_cluster_docs/replica_set.json",
	)
}

func om_cluster_docs_sharded_set_json() ([]byte, error) {
	return bindata_read([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0xd4, 0x58,
		0x4b, 0x6f, 0xdb, 0x38, 0x10, 0xbe, 0x2f, 0xb0, 0xff, 0xc1, 0xd0, 0x5e,
		0x76, 0x81, 0xf8, 0x99, 0x34, 0xdb, 0xe4, 0xd4, 0x34, 0x39, 0x74, 0x81,
		0x4d, 0x5b, 0xd8, 0x7d, 0x1c, 0xda, 0xc0, 0xa0, 0x44, 0xca, 0xe2, 0x96,
		0x22, 0x05, 0x92, 0x76, 0xe2, 0x1a, 0xfa, 0xef, 0x3b, 0xa4, 0x64, 0x8b,
		0x7a, 0xd0, 0x76, 0xdb, 0xdd, 0x20, 0x15, 0x72, 0x88, 0xa9, 0x6f, 0x1e,
		0x9c, 0xf9, 0x66, 0x38, 0xd4, 0xe6, 0xd7, 0x5f, 0x7a, 0xbd, 0x5e, 0x20,
		0x32, 0x4d, 0x05, 0x57, 0xc1, 0x65, 0xaf, 0x58, 0xb0, 0x4f, 0x80, 0xc5,
		0x3d, 0x67, 0x02, 0xe1, 0x97, 0x48, 0x11, 0x78, 0x15, 0x0c, 0x57, 0x48,
		0x0e, 0x19, 0x0d, 0x87, 0xa9, 0xe0, 0x0b, 0x81, 0xc3, 0x7e, 0x9a, 0xaa,
		0x3e, 0x5a, 0x6a, 0x91, 0x22, 0x23, 0x1e, 0x9c, 0x78, 0x64, 0x3f, 0x52,
		0x0e, 0x3f, 0x8d, 0xf6, 0xe0, 0xfa, 0xf2, 0xf3, 0x67, 0x8f, 0x74, 0x21,
		0x9c, 0x97, 0x4a, 0x02, 0x8b, 0xba, 0x09, 0x3f, 0x10, 0xa9, 0x4a, 0xd7,
		0x3e, 0xb9, 0xbe, 0x85, 0x4b, 0xca, 0xb0, 0x5d, 0xad, 0x16, 0xb7, 0xcf,
		0xa6, 0xbd, 0x54, 0xc8, 0x50, 0x6d, 0x24, 0xce, 0xcf, 0x4e, 0x3c, 0x80,
		0x98, 0xa1, 0x95, 0x90, 0xc6, 0xcf, 0xc0, 0x07, 0x59, 0x50, 0x5d, 0xba,
		0x64, 0x60, 0x67, 0x93, 0xb3, 0x8b, 0x68, 0x8c, 0x27, 0xe1, 0xb3, 0x8b,
		0x8b, 0x0b, 0x12, 0x86, 0xf1, 0x38, 0xc6, 0xf1, 0x69, 0x18, 0x8d, 0xc8,
		0x88, 0x9c, 0x86, 0xa7, 0x71, 0xfc, 0x2c, 0x1a, 0x21, 0x14, 0x4f, 0xbc,
		0xea, 0x52, 0xf4, 0xf0, 0x46, 0x39, 0x0a, 0xfd, 0x40, 0xca, 0x8f, 0x04,
		0x0a, 0xbc, 0x64, 0xc4, 0x46, 0xe6, 0xce, 0x87, 0xc9, 0x18, 0xd2, 0xb1,
		0x90, 0xa9, 0x51, 0x24, 0xd4, 0x83, 0x57, 0xd7, 0x52, 0x32, 0x03, 0x49,
		0xb4, 0xce, 0xd4, 0xe5, 0x70, 0x18, 0x23, 0xa5, 0x31, 0x1b, 0x94, 0xf9,
		0x1b, 0x08, 0xb9, 0x00, 0x36, 0xf0, 0xe5, 0xc3, 0x8e, 0x0f, 0xf6, 0x57,
		0xff, 0xe1, 0xf9, 0xf9, 0xfc, 0xfc, 0xac, 0x8f, 0x52, 0xf4, 0x55, 0xf0,
		0xfe, 0xe9, 0x60, 0x32, 0xf8, 0x73, 0xa0, 0x17, 0x5f, 0xbd, 0x56, 0xee,
		0x29, 0x9f, 0x8c, 0x46, 0xcf, 0x33, 0xb6, 0x34, 0x5e, 0xc7, 0x88, 0x29,
		0xb2, 0x07, 0xfa, 0xe1, 0x7a, 0x4a, 0x30, 0x55, 0xfa, 0x86, 0xb1, 0xbd,
		0x71, 0x70, 0xa0, 0x6f, 0x76, 0xec, 0xf6, 0x87, 0xc4, 0x81, 0xbf, 0x97,
		0x47, 0x6b, 0x76, 0x13, 0xd2, 0x46, 0xe7, 0xf5, 0x25, 0xd7, 0x76, 0xc0,
		0x51, 0x6a, 0x6b, 0xca, 0x86, 0x67, 0xcb, 0xfd, 0x2d, 0x22, 0x08, 0x51,
		0xf4, 0x65, 0x99, 0xb9, 0xdc, 0x2f, 0x5e, 0x18, 0x40, 0x55, 0x1f, 0x54,
		0x0b, 0x49, 0xf9, 0xa2, 0x0d, 0x73, 0x0b, 0x25, 0x11, 0x4a, 0x6f, 0x8d,
		0x6d, 0x36, 0x5c, 0x60, 0xa2, 0x06, 0x9f, 0x46, 0x77, 0x79, 0x5e, 0x2b,
		0x57, 0x26, 0x16, 0x6f, 0x91, 0x4e, 0xaa, 0x2a, 0x17, 0x0b, 0x4f, 0x95,
		0x0f, 0x2b, 0xc3, 0x7d, 0xb4, 0x20, 0x5c, 0x0f, 0x00, 0xdb, 0xd4, 0x35,
		0x15, 0x1a, 0x69, 0x52, 0x6f, 0x27, 0xf6, 0x9d, 0xa2, 0x5f, 0xc9, 0xbb,
		0x44, 0x12, 0x95, 0x08, 0x86, 0x6f, 0x5f, 0x02, 0x62, 0x3c, 0x1a, 0x8d,
		0x1a, 0x91, 0x0e, 0x34, 0x4d, 0x2b, 0xd4, 0x2b, 0x69, 0x76, 0x36, 0x39,
		0xab, 0x30, 0x65, 0x5c, 0xf3, 0x2a, 0x26, 0x56, 0x2a, 0x93, 0x22, 0x22,
		0x4a, 0x15, 0xdc, 0xdf, 0x6c, 0x7e, 0x23, 0x28, 0x4a, 0x7a, 0x76, 0xc3,
		0xf9, 0x2e, 0x13, 0x6e, 0x64, 0x90, 0x5c, 0xa8, 0xc9, 0xfc, 0xbc, 0xc3,
		0x4b, 0x4e, 0x74, 0x7b, 0xb5, 0xb0, 0x21, 0xa4, 0x79, 0x35, 0x79, 0x0e,
		0x5e, 0xd7, 0x5f, 0xe7, 0x8d, 0x4d, 0x80, 0x03, 0x34, 0xee, 0xfd, 0x4e,
		0xd5, 0x5f, 0x7c, 0x96, 0x20, 0x89, 0x7b, 0x2f, 0xa0, 0x0b, 0x92, 0x87,
		0x3f, 0xf2, 0x06, 0x2b, 0x02, 0x49, 0x32, 0x46, 0xa3, 0xa2, 0x07, 0x76,
		0x1b, 0x35, 0x88, 0x19, 0xd1, 0xaf, 0xcb, 0x2c, 0x2a, 0xa3, 0x6e, 0xbe,
		0xd9, 0x60, 0xba, 0x2a, 0x95, 0xf6, 0x4e, 0x7b, 0x90, 0xcf, 0x43, 0xfe,
		0x0c, 0x69, 0xdc, 0x34, 0xbe, 0x73, 0xf2, 0x5a, 0xf0, 0x98, 0x2e, 0xbc,
		0x3e, 0x5a, 0x93, 0x90, 0x71, 0x8f, 0x83, 0x11, 0x94, 0xae, 0x26, 0x72,
		0x2a, 0x98, 0x75, 0x30, 0xb2, 0xca, 0xd4, 0x4a, 0xfe, 0x80, 0x4b, 0x09,
		0x52, 0x33, 0x60, 0x19, 0x30, 0xcc, 0xef, 0x54, 0xf1, 0xde, 0xe3, 0x13,
		0x0e, 0x6b, 0x84, 0x5e, 0x45, 0x28, 0x1b, 0x1a, 0x09, 0xb2, 0xe3, 0x35,
		0x46, 0x1a, 0x7d, 0x8f, 0x87, 0x81, 0x5a, 0xc3, 0x6e, 0xd3, 0xbf, 0x85,
		0x2f, 0x1c, 0xc0, 0x38, 0x4d, 0xf9, 0x36, 0xa3, 0x41, 0x4c, 0x21, 0x2c,
		0x1d, 0xcd, 0x24, 0xc8, 0xda, 0x0e, 0xae, 0x95, 0x53, 0x7a, 0x73, 0x43,
		0xde, 0xad, 0xbb, 0xb6, 0xcc, 0x1a, 0xde, 0x3a, 0x25, 0x71, 0xe2, 0x2d,
		0x79, 0x9d, 0x50, 0xd5, 0xae, 0xf6, 0xff, 0xb9, 0x42, 0xbb, 0x5a, 0x5d,
		0xa7, 0x2b, 0x65, 0xd1, 0xbe, 0x5b, 0x67, 0x25, 0xc6, 0x59, 0x28, 0x33,
		0xdf, 0x90, 0x58, 0x55, 0x3d, 0xb7, 0x68, 0x9e, 0xee, 0x4b, 0xe8, 0x53,
		0xc9, 0x2c, 0x4a, 0x48, 0x8a, 0xaa, 0xd6, 0xfc, 0xac, 0x6c, 0x17, 0x05,
		0xb3, 0x5e, 0x30, 0x38, 0xc5, 0x72, 0xf8, 0x41, 0xe0, 0x9c, 0xc9, 0xf3,
		0x93, 0x5a, 0x8a, 0xe1, 0x87, 0xe9, 0x1a, 0x79, 0xb3, 0xb5, 0x94, 0x35,
		0x0a, 0x45, 0x58, 0x6b, 0x2e, 0x19, 0x92, 0x9a, 0x9a, 0x34, 0x13, 0xfc,
		0xba, 0xde, 0x67, 0x8a, 0xa7, 0x8b, 0x1c, 0x73, 0x8a, 0xdd, 0x22, 0xee,
		0xda, 0xe2, 0x0e, 0x9b, 0x92, 0x34, 0x24, 0xb2, 0x66, 0xb2, 0x08, 0x61,
		0xf7, 0xc9, 0xe4, 0x99, 0x7b, 0x1c, 0xb3, 0x95, 0x3d, 0xcf, 0xe1, 0x56,
		0x04, 0x51, 0xc2, 0xa4, 0x44, 0xe4, 0x1b, 0xce, 0xd6, 0x07, 0x0e, 0x64,
		0x0b, 0x4f, 0x28, 0xc6, 0x84, 0x1f, 0x85, 0x04, 0x62, 0xfa, 0x98, 0xd0,
		0x02, 0x67, 0x92, 0xc2, 0x29, 0xa3, 0x8d, 0x0b, 0xe3, 0x7d, 0x38, 0x05,
		0x43, 0x1b, 0xb9, 0x21, 0x0c, 0x19, 0x64, 0x93, 0xa8, 0x35, 0xe4, 0x4a,
		0x68, 0x7b, 0x38, 0x8c, 0xbb, 0x31, 0xc7, 0x10, 0xa4, 0xf9, 0x34, 0x08,
		0xe3, 0x3e, 0x77, 0x07, 0x5b, 0x4b, 0x97, 0xe4, 0x41, 0xca, 0x9c, 0x1e,
		0x62, 0xca, 0xf7, 0x72, 0x63, 0x6f, 0xe8, 0x1e, 0x81, 0x12, 0xe5, 0x68,
		0x72, 0x71, 0xf7, 0xe4, 0x78, 0xe1, 0x91, 0x3d, 0x1c, 0xd2, 0xbd, 0xfe,
		0x3d, 0x5e, 0x48, 0xc7, 0xcd, 0x71, 0xaf, 0x25, 0xf2, 0x13, 0xc5, 0x74,
		0xf2, 0x44, 0x62, 0x3a, 0x7e, 0x7a, 0x31, 0x3d, 0xdc, 0x83, 0x1a, 0xc7,
		0x9a, 0x33, 0xd6, 0x39, 0x6d, 0xa3, 0xab, 0x05, 0x59, 0xa4, 0xb7, 0xbd,
		0x1c, 0x7b, 0x1e, 0xee, 0xb1, 0xb0, 0xb3, 0xa4, 0xd1, 0x62, 0xef, 0x4d,
		0xcd, 0x82, 0xbe, 0xe1, 0x10, 0xdd, 0xc9, 0xd8, 0xee, 0xd8, 0x16, 0x79,
		0x9c, 0xb3, 0xa0, 0xab, 0x67, 0x73, 0x77, 0x9e, 0x27, 0xb8, 0x5f, 0x4e,
		0xd1, 0x9d, 0xfd, 0xbd, 0x18, 0xa9, 0x67, 0x44, 0xc2, 0xf4, 0xe3, 0x8b,
		0x4e, 0x0d, 0x34, 0x2d, 0x26, 0x96, 0x6a, 0x1c, 0x9f, 0x2a, 0x8f, 0x62,
		0xc6, 0x48, 0x54, 0xdd, 0x8f, 0xbd, 0x84, 0x29, 0x29, 0x63, 0x06, 0xac,
		0xe0, 0xb2, 0x4a, 0x60, 0x00, 0xd7, 0x60, 0x14, 0x32, 0x02, 0xf9, 0x68,
		0x14, 0x91, 0x81, 0x8a, 0xb7, 0xf7, 0xb8, 0x28, 0x1b, 0xe0, 0x86, 0x52,
		0xf7, 0x02, 0xee, 0x41, 0xb5, 0x34, 0x59, 0xd0, 0x7b, 0x65, 0x36, 0x15,
		0x78, 0xbf, 0x26, 0xc1, 0x40, 0x9d, 0x31, 0xb1, 0x4e, 0xe1, 0xb2, 0x79,
		0x05, 0xe6, 0x6f, 0x49, 0x94, 0x20, 0x4e, 0x55, 0xda, 0xe6, 0x63, 0x30,
		0xbb, 0x9e, 0x5e, 0xdd, 0xf6, 0x67, 0xaf, 0xae, 0xfa, 0xe3, 0x2a, 0xb3,
		0x4e, 0xb0, 0x82, 0x2f, 0x04, 0xca, 0xcc, 0xf8, 0x03, 0xff, 0x34, 0x5c,
		0x81, 0x15, 0x3b, 0xaa, 0x5f, 0x3a, 0x03, 0xf9, 0x3f, 0x22, 0x54, 0xb5,
		0x69, 0xbc, 0x08, 0x66, 0xb1, 0x34, 0x17, 0xe9, 0xc0, 0xe8, 0x73, 0x74,
		0x2c, 0x61, 0x2b, 0xea, 0x23, 0xe2, 0xda, 0x04, 0xa4, 0xe6, 0x5b, 0x93,
		0xf3, 0x70, 0x41, 0x01, 0x4b, 0x08, 0xa7, 0x94, 0xb7, 0x12, 0x13, 0x50,
		0xb8, 0x62, 0xdb, 0xd0, 0xf9, 0x23, 0x57, 0x22, 0xa5, 0xb0, 0x1f, 0x7b,
		0x8e, 0xfe, 0x0a, 0xb6, 0xcf, 0xac, 0xa3, 0x12, 0x30, 0x25, 0x1d, 0x6f,
		0x8b, 0xeb, 0xfe, 0xe1, 0xcf, 0x1b, 0x6d, 0x92, 0xdb, 0x68, 0x94, 0x89,
		0x6d, 0x7e, 0x35, 0xa8, 0xe9, 0xab, 0x9f, 0x06, 0x3f, 0x6b, 0xa4, 0xae,
		0x2c, 0xb4, 0x23, 0x4e, 0x1d, 0xd2, 0x3f, 0x6a, 0x52, 0x12, 0x84, 0xaf,
		0xf8, 0xfa, 0x06, 0xae, 0xb0, 0xa1, 0xf9, 0x3a, 0xfb, 0x38, 0x56, 0x4d,
		0x3e, 0xed, 0x2e, 0xff, 0x5b, 0xd3, 0x4c, 0x44, 0x88, 0x1d, 0xb1, 0xe1,
		0x8f, 0x70, 0x9c, 0x3e, 0xd6, 0x56, 0xf7, 0xda, 0xfb, 0x16, 0xe2, 0x17,
		0x9f, 0xf2, 0xba, 0x48, 0xdf, 0xd9, 0x9f, 0x6c, 0x03, 0x81, 0x69, 0x80,
		0x14, 0x1d, 0xe4, 0xae, 0xd1, 0x2a, 0x6b, 0xfd, 0xcf, 0x74, 0xd6, 0x76,
		0xc3, 0x03, 0xbd, 0xf0, 0xf7, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xbb,
		0x53, 0x04, 0x1d, 0xd9, 0x17, 0x00, 0x00,
	},
		"om_cluster_docs/sharded_set.json",
	)
}

func om_cluster_docs_standalone_json() ([]byte, error) {
	return bindata_read([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0xcc, 0x56,
		0xcf, 0x73, 0xda, 0x38, 0x14, 0xbe, 0xf3, 0x57, 0x30, 0x3e, 0x17, 0x70,
		0x08, 0xc9, 0x26, 0xb9, 0xd1, 0xe4, 0xd0, 0xc3, 0x66, 0xda, 0x49, 0xda,
		0xe6, 0xd0, 0x76, 0x98, 0x67, 0x4b, 0xb6, 0xb5, 0x95, 0x25, 0x8f, 0x24,
		0x92, 0x90, 0x0c, 0xff, 0xfb, 0x4a, 0xb2, 0x01, 0x0b, 0x4b, 0xd4, 0x9d,
		0x29, 0x9b, 0xf5, 0x09, 0x9e, 0x3e, 0xbd, 0x5f, 0xfa, 0xbe, 0x27, 0xbd,
		0x0e, 0x86, 0xfa, 0x8b, 0x78, 0xa5, 0x08, 0x67, 0x32, 0xba, 0x1a, 0xd6,
		0x06, 0x6b, 0x44, 0xfc, 0x89, 0x51, 0x0e, 0xe8, 0x3d, 0x48, 0xac, 0x57,
		0xa2, 0xc9, 0x23, 0x88, 0x09, 0x25, 0xc9, 0xa4, 0xe4, 0x2c, 0xe7, 0x28,
		0x19, 0x95, 0xa5, 0x1c, 0xc1, 0x52, 0xf1, 0x12, 0xcc, 0xee, 0xe8, 0x9d,
		0x7f, 0xeb, 0x03, 0x61, 0xfa, 0xaf, 0xf1, 0x1d, 0x5d, 0x5f, 0x7d, 0xff,
		0x1e, 0xd8, 0x6c, 0xf7, 0xae, 0x6b, 0x17, 0x91, 0xc5, 0xdc, 0x24, 0x5f,
		0xb1, 0x90, 0x4d, 0x5a, 0xdf, 0x5a, 0x79, 0x25, 0x4b, 0x42, 0x91, 0x35,
		0x6e, 0x6d, 0x9b, 0xef, 0xb5, 0x63, 0xa9, 0x77, 0x10, 0x65, 0xf0, 0xe7,
		0xb3, 0x77, 0xfe, 0xf5, 0x8c, 0xc2, 0x23, 0x17, 0x26, 0xc3, 0x28, 0x80,
		0xc8, 0x89, 0x6a, 0xb2, 0x31, 0xa8, 0xd9, 0x74, 0x76, 0x99, 0x9e, 0xa0,
		0x69, 0x72, 0x76, 0x79, 0x79, 0x89, 0x93, 0x24, 0x3b, 0xc9, 0x50, 0x76,
		0x9a, 0xa4, 0x31, 0x8e, 0xf1, 0x69, 0x72, 0x9a, 0x65, 0x67, 0x69, 0x0c,
		0x90, 0x4d, 0x43, 0xde, 0x4a, 0x78, 0xfe, 0x28, 0x5b, 0xfe, 0x82, 0x38,
		0xc2, 0xfa, 0xe1, 0x38, 0x5a, 0x52, 0x6c, 0x5b, 0xf2, 0x23, 0x00, 0xa9,
		0x28, 0xa8, 0x8c, 0x8b, 0xd2, 0xb8, 0xe1, 0xf2, 0x39, 0xe4, 0x69, 0x29,
		0xa8, 0x41, 0x14, 0x4a, 0x55, 0xf2, 0x6a, 0x32, 0xc9, 0x40, 0x2a, 0x44,
		0xc7, 0xcd, 0x99, 0x8d, 0xb9, 0xc8, 0x35, 0x01, 0xd8, 0xf2, 0x79, 0x4b,
		0x01, 0xfb, 0x6f, 0xf4, 0x7c, 0x71, 0xbe, 0x38, 0x9f, 0x8d, 0xa0, 0x84,
		0x17, 0xce, 0x46, 0xa7, 0xe3, 0xe9, 0xf8, 0xaf, 0xb1, 0xca, 0x5f, 0x42,
		0x41, 0x9e, 0x08, 0x9b, 0xc6, 0xf1, 0x45, 0x45, 0x97, 0x26, 0xe5, 0x0c,
		0xa8, 0xc4, 0x61, 0xe4, 0xd7, 0xeb, 0x3b, 0x8c, 0x88, 0x54, 0x37, 0x94,
		0x1e, 0x6a, 0x41, 0x0b, 0xf9, 0x71, 0x4b, 0xe6, 0x60, 0x37, 0x5a, 0xe8,
		0x2f, 0xa2, 0xaf, 0xdf, 0xf6, 0x41, 0x74, 0xc0, 0x6b, 0xc7, 0xd2, 0x8a,
		0x1b, 0x31, 0x28, 0xad, 0x7c, 0x6c, 0x5b, 0x1a, 0x9e, 0x37, 0xeb, 0x51,
		0x02, 0xe9, 0xcf, 0x65, 0x15, 0xe0, 0x79, 0xc1, 0xa5, 0xda, 0x6c, 0x7e,
		0x7d, 0x65, 0x1c, 0x61, 0x39, 0xfe, 0x16, 0xff, 0x58, 0xaf, 0xdb, 0x4a,
		0xa3, 0x3c, 0xff, 0x04, 0xaa, 0xd8, 0xe9, 0x93, 0xe7, 0x01, 0x7d, 0x4e,
		0xea, 0x60, 0x23, 0xc8, 0x31, 0x53, 0x63, 0x8d, 0xdb, 0x73, 0x73, 0xc7,
		0x15, 0x28, 0xec, 0x8c, 0x00, 0xbb, 0x24, 0xc9, 0x0b, 0xfe, 0x5c, 0x08,
		0x2c, 0x0b, 0x4e, 0xd1, 0xed, 0x7b, 0x0d, 0x38, 0x89, 0xe3, 0xd8, 0x6d,
		0x57, 0xa4, 0x48, 0xb9, 0x03, 0x7d, 0x10, 0xa6, 0x94, 0xe9, 0x6c, 0xe0,
		0x36, 0xc7, 0x54, 0xbd, 0xd5, 0x37, 0x51, 0x5c, 0x10, 0x96, 0xff, 0x47,
		0xa5, 0xef, 0x02, 0xbe, 0x6d, 0xf9, 0x16, 0x5c, 0x09, 0x9e, 0x62, 0x29,
		0xf1, 0x7e, 0xd1, 0x20, 0x72, 0x39, 0x5d, 0x9c, 0x77, 0x53, 0x60, 0x58,
		0x75, 0x8c, 0xb5, 0x27, 0x2e, 0xcc, 0xca, 0xf4, 0x42, 0x67, 0xe4, 0xac,
		0xae, 0xf7, 0xf2, 0x93, 0xba, 0x78, 0x5d, 0xb8, 0xdf, 0x0b, 0x4a, 0x9c,
		0x36, 0x3e, 0xa6, 0x50, 0x4d, 0xcc, 0x06, 0xbc, 0xed, 0x26, 0x02, 0x05,
		0xd1, 0xe1, 0x00, 0x2b, 0xa9, 0x70, 0xf9, 0xb7, 0xee, 0xaa, 0x3f, 0x04,
		0x96, 0x8a, 0xb0, 0x7a, 0xd2, 0xeb, 0x38, 0x19, 0xa1, 0xd8, 0x23, 0xb8,
		0xa8, 0xea, 0xe6, 0xb1, 0x92, 0xad, 0x73, 0x5d, 0x18, 0x26, 0x6c, 0xb2,
		0xb2, 0x67, 0xe8, 0x26, 0x35, 0xf0, 0xa4, 0xd7, 0x9b, 0x4b, 0xc7, 0x24,
		0x80, 0x67, 0x1c, 0x04, 0x53, 0x69, 0xe8, 0xf1, 0x79, 0x55, 0x59, 0x5c,
		0x5d, 0x6d, 0x1b, 0xf0, 0xb8, 0x9b, 0x43, 0xf5, 0x4c, 0x69, 0xad, 0x69,
		0xce, 0x17, 0xf7, 0x69, 0x81, 0x4b, 0xd8, 0x4d, 0xab, 0x33, 0x97, 0x7e,
		0x02, 0x57, 0x94, 0xa4, 0x70, 0x8f, 0x55, 0x7b, 0x44, 0x46, 0x82, 0xbb,
		0x37, 0x48, 0x24, 0x0b, 0x10, 0x48, 0x2b, 0xa6, 0xb1, 0x0d, 0xb6, 0xfe,
		0xdd, 0x57, 0x82, 0x51, 0xd9, 0x17, 0x89, 0xed, 0xed, 0x19, 0x7e, 0x10,
		0x18, 0xeb, 0xa7, 0x27, 0x54, 0x17, 0x3e, 0xac, 0x40, 0xca, 0x27, 0x2e,
		0xd0, 0xd0, 0xad, 0x1c, 0xe9, 0xcc, 0xf8, 0xaa, 0xd4, 0xfa, 0x9c, 0xeb,
		0x30, 0xb7, 0x38, 0x2d, 0x80, 0x11, 0x59, 0x76, 0xaf, 0xfa, 0xe8, 0xfe,
		0xfa, 0x6e, 0x7e, 0x3b, 0xba, 0xff, 0x30, 0x1f, 0x9d, 0xec, 0x28, 0xd0,
		0x1e, 0xba, 0x3f, 0xf1, 0xaa, 0x89, 0xa5, 0x7f, 0xed, 0x85, 0xd1, 0x16,
		0x4b, 0x40, 0x87, 0x68, 0xff, 0xf0, 0x44, 0x3a, 0x2c, 0x4b, 0x39, 0xcb,
		0x48, 0x43, 0xbc, 0x05, 0x2f, 0xc7, 0xc6, 0x63, 0x3b, 0x57, 0x22, 0x21,
		0xa1, 0x18, 0x75, 0x2f, 0xb0, 0x68, 0xa9, 0xbb, 0x21, 0x6f, 0x30, 0xc5,
		0xca, 0x2e, 0xb7, 0x2f, 0xa1, 0x7a, 0xed, 0x01, 0x58, 0xb3, 0xe4, 0x54,
		0xe5, 0x95, 0xa6, 0xc9, 0x12, 0x90, 0x7e, 0x08, 0xf8, 0xf4, 0xb2, 0x3d,
		0x33, 0xef, 0xdd, 0xe5, 0x7f, 0x0c, 0xf5, 0xf0, 0xeb, 0xf8, 0x37, 0xc0,
		0x54, 0xdf, 0xd4, 0x0a, 0x8b, 0xdb, 0x7a, 0x8a, 0x76, 0x2f, 0x3f, 0xf3,
		0xad, 0x3b, 0x56, 0xcf, 0xe5, 0x6b, 0xeb, 0xdf, 0x30, 0x65, 0x7f, 0x28,
		0xfb, 0x2a, 0x24, 0x1a, 0xe2, 0xe7, 0xcd, 0xa1, 0x71, 0xf4, 0xbf, 0xef,
		0xe4, 0xdc, 0xe2, 0xfd, 0x7d, 0xf4, 0x7b, 0xf9, 0x73, 0x29, 0x08, 0x0c,
		0x68, 0xce, 0x56, 0x37, 0x7a, 0xa6, 0x27, 0xe6, 0x65, 0xff, 0x36, 0x59,
		0x18, 0x22, 0xd8, 0x2e, 0x1c, 0x2d, 0x15, 0xca, 0x53, 0xa0, 0x7d, 0x1b,
		0xf2, 0x20, 0x88, 0x7a, 0xab, 0x56, 0xfc, 0x2a, 0xfe, 0x6f, 0x0b, 0xab,
		0xfd, 0xd0, 0x3b, 0xa2, 0xa8, 0x9c, 0x3a, 0x87, 0x9d, 0x38, 0xdb, 0x94,
		0x82, 0x80, 0x03, 0xa2, 0x0b, 0x74, 0xb7, 0x4f, 0x6b, 0xfb, 0x6a, 0xcd,
		0x77, 0xae, 0x7f, 0x22, 0xec, 0xaf, 0xf5, 0x75, 0xac, 0xc8, 0x3d, 0x35,
		0xf5, 0xfb, 0xe1, 0x0f, 0x4a, 0xa9, 0x17, 0x8f, 0x8f, 0xd9, 0xec, 0x60,
		0xcc, 0x7d, 0x53, 0x57, 0x37, 0x7b, 0x52, 0xb0, 0x31, 0x17, 0x61, 0x41,
		0x78, 0x9f, 0x1a, 0xe6, 0x69, 0xe3, 0xbc, 0x57, 0x8c, 0xb3, 0xce, 0x03,
		0x65, 0x3d, 0x58, 0x0f, 0xfe, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x44, 0xd4,
		0x65, 0x75, 0x6d, 0x11, 0x00, 0x00,
	},
		"om_cluster_docs/standalone.json",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"om_cluster_docs/replica_set.json": om_cluster_docs_replica_set_json,
	"om_cluster_docs/sharded_set.json": om_cluster_docs_sharded_set_json,
	"om_cluster_docs/standalone.json": om_cluster_docs_standalone_json,
}
