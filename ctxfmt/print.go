package ctxfmt

import (
	"io"
	"unicode/utf8"
)

type printer struct {
	To io.Writer
}

func (p *printer) Write(buf []byte) (int, error) {
	return p.To.Write(buf)
}

func (p *printer) WriteByte(b byte) error {
	if bw, ok := p.To.(io.ByteWriter); ok {
		return bw.WriteByte(b)
	}

	_, err := p.To.Write([]byte{b})
	return err
}

func (p *printer) WriteString(s string) (int, error) {
	if sw, ok := p.To.(io.StringWriter); ok {
		return sw.WriteString(s)
	}
	return p.To.Write(unsafeBytes(s))
}

func (p *printer) WriteRune(r rune) error {
	if rw, ok := p.To.(interface{ WriteRune(rune) error }); ok {
		return rw.WriteRune(r)
	}

	if r < utf8.RuneSelf {
		return p.WriteByte(byte(r))
	}

	var runeBuf [utf8.UTFMax]byte
	n := utf8.EncodeRune(runeBuf[:], r)
	_, err := p.Write(runeBuf[:n])
	return err
}

func (p *printer) onString(s string) {
	p.WriteString(s)
}
