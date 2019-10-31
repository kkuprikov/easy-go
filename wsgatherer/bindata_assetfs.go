package wsgatherer

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
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

var _static_index_html = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x93\x3f\x6f\xdb\x30\x10\xc5\xf7\x7c\x8a\x03\x97\xd8\x48\x21\x46\x4e\x81\x02\x36\xe5\xa1\x40\x86\x0e\x9d\x52\xa0\x33\x4d\x5e\x2c\xd6\x12\x4f\x10\x4f\x56\x8d\xc0\xdf\xbd\x90\x68\x31\xb2\x8b\xfe\x99\x44\x1e\x1e\xdf\xef\x91\xba\x53\x25\xd7\xd5\xf6\x0e\x40\x95\xa8\xed\x56\xc9\xf1\x33\xec\x77\x64\x4f\xc3\x02\x40\x05\xd3\xba\x86\x81\x4f\x0d\x16\x82\xf1\x27\xcb\x1f\xfa\xa8\x63\x55\x44\x0d\xc0\x51\xb7\x10\xc8\x1c\xa0\x00\xdf\x55\xd5\x66\x56\xee\x43\xd7\x3a\x28\x40\xf4\x61\x2d\x65\xbe\xfa\x94\x3d\x66\x8f\x59\xbe\xce\x57\x4f\x1f\x65\x1f\x64\xe8\x76\x83\xd9\x0e\x65\x68\xd0\xb0\x66\x6a\x83\xcc\x57\x4f\x62\x73\x77\x71\xe9\x9d\xb7\xd4\x67\xe4\x2b\xd2\x16\x0a\x78\xed\xbc\x61\x47\x7e\xb1\x84\xb7\x8b\x04\xc0\x90\x0f\x54\x61\x56\xd1\x7e\x21\xa2\x54\x2c\x93\x05\xa4\x74\xd8\xc3\x77\xdc\xbd\x90\x39\x20\x2f\xc6\x6c\xb7\xaa\x8c\x3c\x35\xe8\xff\x04\xba\x41\x19\xf2\x1e\x0d\xa3\x05\x26\x10\xf0\x00\x93\xe7\xa4\x3e\xff\x6e\x6f\x2a\x0a\x38\xf7\xc7\x7f\x03\x1c\x79\x18\xcf\x59\x58\x0c\x18\xcc\x0c\x59\x84\x07\x10\x4b\xf1\x77\x5a\x8d\x21\xe8\xfd\x7f\xf3\x26\x79\x8b\x06\xdd\x11\xed\x1a\x22\xce\x6a\xd6\xd7\xa0\x5b\xe2\x64\x0f\x01\xbd\xbd\x7a\xb3\xa1\x0f\xea\xb0\x87\x02\x2c\x99\xae\x46\xcf\xd9\x1e\xf9\xb9\xc2\x61\xf9\xf9\xf4\xc5\x26\xac\x58\x66\x47\x5d\x75\xb8\xb9\xbe\xc5\xe8\x58\x87\x7d\x0a\x70\x8e\xcd\x29\x63\x1f\x5e\x5a\xb5\xcc\xb7\xe9\xdf\xc2\xb3\x29\x09\xbe\x61\x60\x25\xcb\xfc\x22\x78\xa5\xb6\x9e\x5a\x56\x35\xdb\xc4\xf8\x1a\xe1\xeb\x54\x50\xce\x37\x1d\xcf\x9e\xc8\xd9\x22\x45\x9c\x95\xdf\xa7\x62\x5e\x1d\x6f\x50\xdc\xbf\x09\x67\xc5\x1a\xc4\xd0\xcc\x1f\x40\xe0\x11\x3d\x0f\xfb\x12\x35\xd7\xba\x11\xe7\xfb\x74\x46\xa6\x54\xf2\x12\x4b\xc9\xf7\xb0\x6a\xd7\x31\x93\x87\xa1\x73\x9c\x39\x14\x22\x3e\xf0\x46\x6c\x5f\xd0\xdb\x29\xbd\x92\x51\x36\x0e\xb0\x8c\x13\xac\x64\x1c\xf1\x5f\x01\x00\x00\xff\xff\x7a\x79\xb7\x0e\xea\x03\x00\x00")

func static_index_html() ([]byte, error) {
	return bindata_read(
		_static_index_html,
		"static/index.html",
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
	"static/index.html": static_index_html,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func     func() ([]byte, error)
	Children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"static": &_bintree_t{nil, map[string]*_bintree_t{
		"index.html": &_bintree_t{static_index_html, map[string]*_bintree_t{}},
	}},
}}

func assetFS() *assetfs.AssetFS {
	assetInfo := func(path string) (os.FileInfo, error) {
		return os.Stat(path)
	}
	for k := range _bintree.Children {
		return &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: assetInfo, Prefix: k}
	}
	panic("unreachable")
}
