package ctxfmt

import (
	"io"
	"unicode/utf8"
)

type printer struct {
	To      io.Writer
	written int
	err     error
}

func (p *printer) Write(buf []byte) (int, error) {
	if p.err != nil {
		return 0, p.err
	}

	return p.doWrite(buf)
}

func (p *printer) doWrite(buf []byte) (int, error) {
	return p.upd(p.To.Write(buf))
}

func (p *printer) upd(n int, err error) (int, error) {
	p.written += n
	p.err = err
	return n, err
}

func (p *printer) WriteByte(b byte) error {
	if p.err != nil {
		return p.err
	}

	if bw, ok := p.To.(io.ByteWriter); ok {
		err := bw.WriteByte(b)
		if err != nil {
			p.err = err
		} else {
			p.written++
		}
		return err
	}

	_, err := p.doWrite([]byte{b})
	return err
}

func (p *printer) WriteString(s string) (int, error) {
	if p.err != nil {
		return 0, p.err
	}

	if sw, ok := p.To.(io.StringWriter); ok {
		return p.upd(sw.WriteString(s))
	}
	return p.doWrite(unsafeBytes(s))
}

func (p *printer) WriteRune(r rune) error {
	if p.err != nil {
		return p.err
	}

	if rw, ok := p.To.(interface{ WriteRune(rune) error }); ok {
		p.err = rw.WriteRune(r)
		return p.err
	}

	if r < utf8.RuneSelf {
		return p.WriteByte(byte(r))
	}

	var runeBuf [utf8.UTFMax]byte
	n := utf8.EncodeRune(runeBuf[:], r)
	_, err := p.doWrite(runeBuf[:n])
	return err
}

func (p *printer) onString(s string) {
	p.WriteString(s)
}
