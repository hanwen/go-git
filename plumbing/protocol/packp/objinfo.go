package packp

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/pktline"
)

// TODO
type ObjectInfoRequest struct {
	Want string
	OID  []plumbing.Hash
}

func (r *ObjectInfoRequest) Encode(w io.Writer) error {
	e := pktline.NewEncoder(w)
	e.EncodeString("command=object-info\n")
	e.Delimiter()
	e.EncodeString("size\n")
	for _, id := range r.OID {
		e.Encodef("oid %s\n", id)
	}
	e.Flush()
	return nil
}

type ObjectInfoSize struct {
	Hash plumbing.Hash
	Size uint64
}

type ObjectInfoResponse struct {
	Sizes []ObjectInfoSize
}

func (rep *ObjectInfoResponse) Decode(r io.Reader) error {
	d := pktline.NewScanner(r)
	d.Scan()

	if line := d.Bytes(); string(line) != "size" {
		return fmt.Errorf("got %q, want 'size'", line)
	}

	for d.Scan() {
		line := d.Bytes()
		if len(line) == 0 {
			continue
		}
		idx := bytes.IndexByte(line, ' ')
		if idx == -1 {
			continue
		}
		h := plumbing.NewHash(string(line[:idx]))
		line = line[idx+1:]
		if len(line) == 0 {
			continue
		}
		sz, err := strconv.ParseUint(string(line), 10, 64)
		if err != nil {
			return err
		}

		rep.Sizes = append(rep.Sizes, ObjectInfoSize{
			Hash: h,
			Size: sz,
		})
	}
	return nil
}
