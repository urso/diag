package ctxfmt

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
)

type (
	renamedBool       bool
	renamedInt        int
	renamedInt8       int8
	renamedInt16      int16
	renamedInt32      int32
	renamedInt64      int64
	renamedUint       uint
	renamedUint8      uint8
	renamedUint16     uint16
	renamedUint32     uint32
	renamedUint64     uint64
	renamedUintptr    uintptr
	renamedString     string
	renamedBytes      []byte
	renamedFloat32    float32
	renamedFloat64    float64
	renamedComplex64  complex64
	renamedComplex128 complex128
)

func TestFmtInterface(t *testing.T) {
	var i1 interface{}
	i1 = "abc"
	s := fmt.Sprintf("%s", i1)
	if s != "abc" {
		t.Errorf(`Sprintf("%%s", empty("abc")) = %q want %q`, s, "abc")
	}
}

var (
	NaN    = math.NaN()
	posInf = math.Inf(1)
	negInf = math.Inf(-1)

	intVar = 0

	array  = [5]int{1, 2, 3, 4, 5}
	iarray = [4]interface{}{1, "hello", 2.5, nil}
	slice  = array[:]
	islice = iarray[:]
)

type A struct {
	i int
	j uint
	s string
	x []int
}

type I int

func (i I) String() string { return fmt.Sprintf("<%d>", int(i)) }

type B struct {
	I I
	j int
}

type C struct {
	i int
	B
}

type F int

func (f F) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "<%c=F(%d)>", c, int(f))
}

type G int

func (g G) GoString() string {
	return fmt.Sprintf("GoString(%d)", int(g))
}

type S struct {
	F F // a struct field that Formats
	G G // a struct field that GoStrings
}

type SI struct {
	I interface{}
}

// P is a type with a String method with pointer receiver for testing %p.
type P int

var pValue P

func (p *P) String() string {
	return "String(p)"
}

var barray = [5]renamedUint8{1, 2, 3, 4, 5}
var bslice = barray[:]

type byteStringer byte

func (byteStringer) String() string {
	return "X"
}

var byteStringerSlice = []byteStringer{'h', 'e', 'l', 'l', 'o'}

type byteFormatter byte

func (byteFormatter) Format(f fmt.State, _ rune) {
	fmt.Fprint(f, "X")
}

var byteFormatterSlice = []byteFormatter{'h', 'e', 'l', 'l', 'o'}

type writeStringFormatter string

func (sf writeStringFormatter) Format(f fmt.State, c rune) {
	if sw, ok := f.(io.StringWriter); ok {
		sw.WriteString("***" + string(sf) + "***")
	}
}

var fmtTests = []struct {
	fmt string
	val interface{}
	out string
}{
	{"%d", 12345, "12345"},
	{"%v", 12345, "12345"},
	{"%t", true, "true"},

	// basic string
	{"%s", "abc", "abc"},
	{"%q", "abc", `"abc"`},
	{"%x", "abc", "616263"},
	{"%x", "\xff\xf0\x0f\xff", "fff00fff"},
	{"%X", "\xff\xf0\x0f\xff", "FFF00FFF"},
	{"%x", "", ""},
	{"% x", "", ""},
	{"%#x", "", ""},
	{"%# x", "", ""},
	{"%x", "xyz", "78797a"},
	{"%X", "xyz", "78797A"},
	{"% x", "xyz", "78 79 7a"},
	{"% X", "xyz", "78 79 7A"},
	{"%#x", "xyz", "0x78797a"},
	{"%#X", "xyz", "0X78797A"},
	{"%# x", "xyz", "0x78 0x79 0x7a"},
	{"%# X", "xyz", "0X78 0X79 0X7A"},

	// basic bytes
	{"%s", []byte("abc"), "abc"},
	{"%s", [3]byte{'a', 'b', 'c'}, "abc"},
	{"%s", &[3]byte{'a', 'b', 'c'}, "&abc"},
	{"%q", []byte("abc"), `"abc"`},
	{"%x", []byte("abc"), "616263"},
	{"%x", []byte("\xff\xf0\x0f\xff"), "fff00fff"},
	{"%X", []byte("\xff\xf0\x0f\xff"), "FFF00FFF"},
	{"%x", []byte(""), ""},
	{"% x", []byte(""), ""},
	{"%#x", []byte(""), ""},
	{"%# x", []byte(""), ""},
	{"%x", []byte("xyz"), "78797a"},
	{"%X", []byte("xyz"), "78797A"},
	{"% x", []byte("xyz"), "78 79 7a"},
	{"% X", []byte("xyz"), "78 79 7A"},
	{"%#x", []byte("xyz"), "0x78797a"},
	{"%#X", []byte("xyz"), "0X78797A"},
	{"%# x", []byte("xyz"), "0x78 0x79 0x7a"},
	{"%# X", []byte("xyz"), "0X78 0X79 0X7A"},

	// escaped strings
	{"%q", "", `""`},
	{"%#q", "", "``"},
	{"%q", "\"", `"\""`},
	{"%#q", "\"", "`\"`"},
	{"%q", "`", `"` + "`" + `"`},
	{"%#q", "`", `"` + "`" + `"`},
	{"%q", "\n", `"\n"`},
	{"%#q", "\n", `"\n"`},
	{"%q", `\n`, `"\\n"`},
	{"%#q", `\n`, "`\\n`"},
	{"%q", "abc", `"abc"`},
	{"%#q", "abc", "`abc`"},
	{"%q", "日本語", `"日本語"`},
	{"%+q", "日本語", `"\u65e5\u672c\u8a9e"`},
	{"%#q", "日本語", "`日本語`"},
	{"%#+q", "日本語", "`日本語`"},
	{"%q", "\a\b\f\n\r\t\v\"\\", `"\a\b\f\n\r\t\v\"\\"`},
	{"%+q", "\a\b\f\n\r\t\v\"\\", `"\a\b\f\n\r\t\v\"\\"`},
	{"%#q", "\a\b\f\n\r\t\v\"\\", `"\a\b\f\n\r\t\v\"\\"`},
	{"%#+q", "\a\b\f\n\r\t\v\"\\", `"\a\b\f\n\r\t\v\"\\"`},
	{"%q", "☺", `"☺"`},
	{"% q", "☺", `"☺"`}, // The space modifier should have no effect.
	{"%+q", "☺", `"\u263a"`},
	{"%#q", "☺", "`☺`"},
	{"%#+q", "☺", "`☺`"},
	{"%10q", "⌘", `       "⌘"`},
	{"%+10q", "⌘", `  "\u2318"`},
	{"%-10q", "⌘", `"⌘"       `},
	{"%+-10q", "⌘", `"\u2318"  `},
	{"%010q", "⌘", `0000000"⌘"`},
	{"%+010q", "⌘", `00"\u2318"`},
	{"%-010q", "⌘", `"⌘"       `}, // 0 has no effect when - is present.
	{"%+-010q", "⌘", `"\u2318"  `},
	{"%#8q", "\n", `    "\n"`},
	{"%#+8q", "\r", `    "\r"`},
	{"%#-8q", "\t", "`	`     "},
	{"%#+-8q", "\b", `"\b"    `},
	{"%q", "abc\xffdef", `"abc\xffdef"`},
	{"%+q", "abc\xffdef", `"abc\xffdef"`},
	{"%#q", "abc\xffdef", `"abc\xffdef"`},
	{"%#+q", "abc\xffdef", `"abc\xffdef"`},
	// Runes that are not printable.
	{"%q", "\U0010ffff", `"\U0010ffff"`},
	{"%+q", "\U0010ffff", `"\U0010ffff"`},
	{"%#q", "\U0010ffff", "`􏿿`"},
	{"%#+q", "\U0010ffff", "`􏿿`"},
	// Runes that are not valid.
	{"%q", string(0x110000), `"�"`},
	{"%+q", string(0x110000), `"\ufffd"`},
	{"%#q", string(0x110000), "`�`"},
	{"%#+q", string(0x110000), "`�`"},

	// characters
	{"%c", uint('x'), "x"},
	{"%c", 0xe4, "ä"},
	{"%c", 0x672c, "本"},
	{"%c", '日', "日"},
	{"%.0c", '⌘', "⌘"}, // Specifying precision should have no effect.
	{"%3c", '⌘', "  ⌘"},
	{"%-3c", '⌘', "⌘  "},
	// Runes that are not printable.
	{"%c", '\U00000e00', "\u0e00"},
	{"%c", '\U0010ffff', "\U0010ffff"},
	// Runes that are not valid.
	{"%c", -1, "�"},
	{"%c", 0xDC80, "�"},
	{"%c", rune(0x110000), "�"},
	{"%c", int64(0xFFFFFFFFF), "�"},
	{"%c", uint64(0xFFFFFFFFF), "�"},

	// escaped characters
	{"%q", uint(0), `'\x00'`},
	{"%+q", uint(0), `'\x00'`},
	{"%q", '"', `'"'`},
	{"%+q", '"', `'"'`},
	{"%q", '\'', `'\''`},
	{"%+q", '\'', `'\''`},
	{"%q", '`', "'`'"},
	{"%+q", '`', "'`'"},
	{"%q", 'x', `'x'`},
	{"%+q", 'x', `'x'`},
	{"%q", 'ÿ', `'ÿ'`},
	{"%+q", 'ÿ', `'\u00ff'`},
	{"%q", '\n', `'\n'`},
	{"%+q", '\n', `'\n'`},
	{"%q", '☺', `'☺'`},
	{"%+q", '☺', `'\u263a'`},
	{"% q", '☺', `'☺'`},  // The space modifier should have no effect.
	{"%.0q", '☺', `'☺'`}, // Specifying precision should have no effect.
	{"%10q", '⌘', `       '⌘'`},
	{"%+10q", '⌘', `  '\u2318'`},
	{"%-10q", '⌘', `'⌘'       `},
	{"%+-10q", '⌘', `'\u2318'  `},
	{"%010q", '⌘', `0000000'⌘'`},
	{"%+010q", '⌘', `00'\u2318'`},
	{"%-010q", '⌘', `'⌘'       `}, // 0 has no effect when - is present.
	{"%+-010q", '⌘', `'\u2318'  `},
	// Runes that are not printable.
	{"%q", '\U00000e00', `'\u0e00'`},
	{"%q", '\U0010ffff', `'\U0010ffff'`},
	// Runes that are not valid.
	{"%q", int32(-1), "%!q(int32=-1)"},
	{"%q", 0xDC80, `'�'`},
	{"%q", rune(0x110000), "%!q(int32=1114112)"},
	{"%q", int64(0xFFFFFFFFF), "%!q(int64=68719476735)"},
	{"%q", uint64(0xFFFFFFFFF), "%!q(uint64=68719476735)"},

	// width
	{"%5s", "abc", "  abc"},
	{"%5s", []byte("abc"), "  abc"},
	{"%2s", "\u263a", " ☺"},
	{"%2s", []byte("\u263a"), " ☺"},
	{"%-5s", "abc", "abc  "},
	{"%-5s", []byte("abc"), "abc  "},
	{"%05s", "abc", "00abc"},
	{"%05s", []byte("abc"), "00abc"},
	{"%5s", "abcdefghijklmnopqrstuvwxyz", "abcdefghijklmnopqrstuvwxyz"},
	{"%5s", []byte("abcdefghijklmnopqrstuvwxyz"), "abcdefghijklmnopqrstuvwxyz"},
	{"%.5s", "abcdefghijklmnopqrstuvwxyz", "abcde"},
	{"%.5s", []byte("abcdefghijklmnopqrstuvwxyz"), "abcde"},
	{"%.0s", "日本語日本語", ""},
	{"%.0s", []byte("日本語日本語"), ""},
	{"%.5s", "日本語日本語", "日本語日本"},
	{"%.5s", []byte("日本語日本語"), "日本語日本"},
	{"%.10s", "日本語日本語", "日本語日本語"},
	{"%.10s", []byte("日本語日本語"), "日本語日本語"},
	{"%08q", "abc", `000"abc"`},
	{"%08q", []byte("abc"), `000"abc"`},
	{"%-8q", "abc", `"abc"   `},
	{"%-8q", []byte("abc"), `"abc"   `},
	{"%.5q", "abcdefghijklmnopqrstuvwxyz", `"abcde"`},
	{"%.5q", []byte("abcdefghijklmnopqrstuvwxyz"), `"abcde"`},
	{"%.5x", "abcdefghijklmnopqrstuvwxyz", "6162636465"},
	{"%.5x", []byte("abcdefghijklmnopqrstuvwxyz"), "6162636465"},
	{"%.3q", "日本語日本語", `"日本語"`},
	{"%.3q", []byte("日本語日本語"), `"日本語"`},
	{"%.1q", "日本語", `"日"`},
	{"%.1q", []byte("日本語"), `"日"`},
	{"%.1x", "日本語", "e6"},
	{"%.1X", []byte("日本語"), "E6"},
	{"%10.1q", "日本語日本語", `       "日"`},
	{"%10.1q", []byte("日本語日本語"), `       "日"`},
	{"%10v", nil, "     <nil>"},
	{"%-10v", nil, "<nil>     "},

	// integers
	{"%d", uint(12345), "12345"},
	{"%d", int(-12345), "-12345"},
	{"%d", ^uint8(0), "255"},
	{"%d", ^uint16(0), "65535"},
	{"%d", ^uint32(0), "4294967295"},
	{"%d", ^uint64(0), "18446744073709551615"},
	{"%d", int8(-1 << 7), "-128"},
	{"%d", int16(-1 << 15), "-32768"},
	{"%d", int32(-1 << 31), "-2147483648"},
	{"%d", int64(-1 << 63), "-9223372036854775808"},
	{"%.d", 0, ""},
	{"%.0d", 0, ""},
	{"%6.0d", 0, "      "},
	{"%06.0d", 0, "      "},
	{"% d", 12345, " 12345"},
	{"%+d", 12345, "+12345"},
	{"%+d", -12345, "-12345"},
	{"%b", 7, "111"},
	{"%b", -6, "-110"},
	{"%#b", 7, "0b111"},
	{"%#b", -6, "-0b110"},
	{"%b", ^uint32(0), "11111111111111111111111111111111"},
	{"%b", ^uint64(0), "1111111111111111111111111111111111111111111111111111111111111111"},
	{"%b", int64(-1 << 63), zeroFill("-1", 63, "")},
	{"%o", 01234, "1234"},
	{"%o", -01234, "-1234"},
	{"%#o", 01234, "01234"},
	{"%#o", -01234, "-01234"},
	{"%O", 01234, "0o1234"},
	{"%O", -01234, "-0o1234"},
	{"%o", ^uint32(0), "37777777777"},
	{"%o", ^uint64(0), "1777777777777777777777"},
	{"%#X", 0, "0X0"},
	{"%x", 0x12abcdef, "12abcdef"},
	{"%X", 0x12abcdef, "12ABCDEF"},
	{"%x", ^uint32(0), "ffffffff"},
	{"%X", ^uint64(0), "FFFFFFFFFFFFFFFF"},
	{"%.20b", 7, "00000000000000000111"},
	{"%10d", 12345, "     12345"},
	{"%10d", -12345, "    -12345"},
	{"%+10d", 12345, "    +12345"},
	{"%010d", 12345, "0000012345"},
	{"%010d", -12345, "-000012345"},
	{"%20.8d", 1234, "            00001234"},
	{"%20.8d", -1234, "           -00001234"},
	{"%020.8d", 1234, "            00001234"},
	{"%020.8d", -1234, "           -00001234"},
	{"%-20.8d", 1234, "00001234            "},
	{"%-20.8d", -1234, "-00001234           "},
	{"%-#20.8x", 0x1234abc, "0x01234abc          "},
	{"%-#20.8X", 0x1234abc, "0X01234ABC          "},
	{"%-#20.8o", 01234, "00001234            "},

	// Test correct f.intbuf overflow checks.
	{"%068d", 1, zeroFill("", 68, "1")},
	{"%068d", -1, zeroFill("-", 67, "1")},
	{"%#.68x", 42, zeroFill("0x", 68, "2a")},
	{"%.68d", -42, zeroFill("-", 68, "42")},
	{"%+.68d", 42, zeroFill("+", 68, "42")},
	{"% .68d", 42, zeroFill(" ", 68, "42")},
	{"% +.68d", 42, zeroFill("+", 68, "42")},

	// unicode format
	{"%U", 0, "U+0000"},
	{"%U", -1, "U+FFFFFFFFFFFFFFFF"},
	{"%U", '\n', `U+000A`},
	{"%#U", '\n', `U+000A`},
	{"%+U", 'x', `U+0078`},       // Plus flag should have no effect.
	{"%# U", 'x', `U+0078 'x'`},  // Space flag should have no effect.
	{"%#.2U", 'x', `U+0078 'x'`}, // Precisions below 4 should print 4 digits.
	{"%U", '\u263a', `U+263A`},
	{"%#U", '\u263a', `U+263A '☺'`},
	{"%U", '\U0001D6C2', `U+1D6C2`},
	{"%#U", '\U0001D6C2', `U+1D6C2 '𝛂'`},
	{"%#14.6U", '⌘', "  U+002318 '⌘'"},
	{"%#-14.6U", '⌘', "U+002318 '⌘'  "},
	{"%#014.6U", '⌘', "  U+002318 '⌘'"},
	{"%#-014.6U", '⌘', "U+002318 '⌘'  "},
	{"%.68U", uint(42), zeroFill("U+", 68, "2A")},
	{"%#.68U", '日', zeroFill("U+", 68, "65E5") + " '日'"},

	// floats
	{"%+.3e", 0.0, "+0.000e+00"},
	{"%+.3e", 1.0, "+1.000e+00"},
	{"%+.3x", 0.0, "+0x0.000p+00"},
	{"%+.3x", 1.0, "+0x1.000p+00"},
	{"%+.3f", -1.0, "-1.000"},
	{"%+.3F", -1.0, "-1.000"},
	{"%+.3F", float32(-1.0), "-1.000"},
	{"%+07.2f", 1.0, "+001.00"},
	{"%+07.2f", -1.0, "-001.00"},
	{"%-07.2f", 1.0, "1.00   "},
	{"%-07.2f", -1.0, "-1.00  "},
	{"%+-07.2f", 1.0, "+1.00  "},
	{"%+-07.2f", -1.0, "-1.00  "},
	{"%-+07.2f", 1.0, "+1.00  "},
	{"%-+07.2f", -1.0, "-1.00  "},
	{"%+10.2f", +1.0, "     +1.00"},
	{"%+10.2f", -1.0, "     -1.00"},
	{"% .3E", -1.0, "-1.000E+00"},
	{"% .3e", 1.0, " 1.000e+00"},
	{"% .3X", -1.0, "-0X1.000P+00"},
	{"% .3x", 1.0, " 0x1.000p+00"},
	{"%+.3g", 0.0, "+0"},
	{"%+.3g", 1.0, "+1"},
	{"%+.3g", -1.0, "-1"},
	{"% .3g", -1.0, "-1"},
	{"% .3g", 1.0, " 1"},
	{"%b", float32(1.0), "8388608p-23"},
	{"%b", 1.0, "4503599627370496p-52"},
	// Test sharp flag used with floats.
	{"%#g", 1e-323, "1.00000e-323"},
	{"%#g", -1.0, "-1.00000"},
	{"%#g", 1.1, "1.10000"},
	{"%#g", 123456.0, "123456."},
	{"%#g", 1234567.0, "1.234567e+06"},
	{"%#g", 1230000.0, "1.23000e+06"},
	{"%#g", 1000000.0, "1.00000e+06"},
	{"%#.0f", 1.0, "1."},
	{"%#.0e", 1.0, "1.e+00"},
	{"%#.0x", 1.0, "0x1.p+00"},
	{"%#.0g", 1.0, "1."},
	{"%#.0g", 1100000.0, "1.e+06"},
	{"%#.4f", 1.0, "1.0000"},
	{"%#.4e", 1.0, "1.0000e+00"},
	{"%#.4x", 1.0, "0x1.0000p+00"},
	{"%#.4g", 1.0, "1.000"},
	{"%#.4g", 100000.0, "1.000e+05"},
	{"%#.0f", 123.0, "123."},
	{"%#.0e", 123.0, "1.e+02"},
	{"%#.0x", 123.0, "0x1.p+07"},
	{"%#.0g", 123.0, "1.e+02"},
	{"%#.4f", 123.0, "123.0000"},
	{"%#.4e", 123.0, "1.2300e+02"},
	{"%#.4x", 123.0, "0x1.ec00p+06"},
	{"%#.4g", 123.0, "123.0"},
	{"%#.4g", 123000.0, "1.230e+05"},
	{"%#9.4g", 1.0, "    1.000"},
	// The sharp flag has no effect for binary float format.
	{"%#b", 1.0, "4503599627370496p-52"},
	// Precision has no effect for binary float format.
	{"%.4b", float32(1.0), "8388608p-23"},
	{"%.4b", -1.0, "-4503599627370496p-52"},
	// Test correct f.intbuf boundary checks.
	{"%.68f", 1.0, zeroFill("1.", 68, "")},
	{"%.68f", -1.0, zeroFill("-1.", 68, "")},
	// float infinites and NaNs
	{"%f", posInf, "+Inf"},
	{"%.1f", negInf, "-Inf"},
	{"% f", NaN, " NaN"},
	{"%20f", posInf, "                +Inf"},
	{"% 20F", posInf, "                 Inf"},
	{"% 20e", negInf, "                -Inf"},
	{"% 20x", negInf, "                -Inf"},
	{"%+20E", negInf, "                -Inf"},
	{"%+20X", negInf, "                -Inf"},
	{"% +20g", negInf, "                -Inf"},
	{"%+-20G", posInf, "+Inf                "},
	{"%20e", NaN, "                 NaN"},
	{"%20x", NaN, "                 NaN"},
	{"% +20E", NaN, "                +NaN"},
	{"% +20X", NaN, "                +NaN"},
	{"% -20g", NaN, " NaN                "},
	{"%+-20G", NaN, "+NaN                "},
	// Zero padding does not apply to infinities and NaN.
	{"%+020e", posInf, "                +Inf"},
	{"%+020x", posInf, "                +Inf"},
	{"%-020f", negInf, "-Inf                "},
	{"%-020E", NaN, "NaN                 "},
	{"%-020X", NaN, "NaN                 "},

	// complex values
	{"%.f", 0i, "(0+0i)"},
	{"% .f", 0i, "( 0+0i)"},
	{"%+.f", 0i, "(+0+0i)"},
	{"% +.f", 0i, "(+0+0i)"},
	{"%+.3e", 0i, "(+0.000e+00+0.000e+00i)"},
	{"%+.3x", 0i, "(+0x0.000p+00+0x0.000p+00i)"},
	{"%+.3f", 0i, "(+0.000+0.000i)"},
	{"%+.3g", 0i, "(+0+0i)"},
	{"%+.3e", 1 + 2i, "(+1.000e+00+2.000e+00i)"},
	{"%+.3x", 1 + 2i, "(+0x1.000p+00+0x1.000p+01i)"},
	{"%+.3f", 1 + 2i, "(+1.000+2.000i)"},
	{"%+.3g", 1 + 2i, "(+1+2i)"},
	{"%.3e", 0i, "(0.000e+00+0.000e+00i)"},
	{"%.3x", 0i, "(0x0.000p+00+0x0.000p+00i)"},
	{"%.3f", 0i, "(0.000+0.000i)"},
	{"%.3F", 0i, "(0.000+0.000i)"},
	{"%.3F", complex64(0i), "(0.000+0.000i)"},
	{"%.3g", 0i, "(0+0i)"},
	{"%.3e", 1 + 2i, "(1.000e+00+2.000e+00i)"},
	{"%.3x", 1 + 2i, "(0x1.000p+00+0x1.000p+01i)"},
	{"%.3f", 1 + 2i, "(1.000+2.000i)"},
	{"%.3g", 1 + 2i, "(1+2i)"},
	{"%.3e", -1 - 2i, "(-1.000e+00-2.000e+00i)"},
	{"%.3x", -1 - 2i, "(-0x1.000p+00-0x1.000p+01i)"},
	{"%.3f", -1 - 2i, "(-1.000-2.000i)"},
	{"%.3g", -1 - 2i, "(-1-2i)"},
	{"% .3E", -1 - 2i, "(-1.000E+00-2.000E+00i)"},
	{"% .3X", -1 - 2i, "(-0X1.000P+00-0X1.000P+01i)"},
	{"%+.3g", 1 + 2i, "(+1+2i)"},
	{"%+.3g", complex64(1 + 2i), "(+1+2i)"},
	{"%#g", 1 + 2i, "(1.00000+2.00000i)"},
	{"%#g", 123456 + 789012i, "(123456.+789012.i)"},
	{"%#g", 1e-10i, "(0.00000+1.00000e-10i)"},
	{"%#g", -1e10 - 1.11e100i, "(-1.00000e+10-1.11000e+100i)"},
	{"%#.0f", 1.23 + 1.0i, "(1.+1.i)"},
	{"%#.0e", 1.23 + 1.0i, "(1.e+00+1.e+00i)"},
	{"%#.0x", 1.23 + 1.0i, "(0x1.p+00+0x1.p+00i)"},
	{"%#.0g", 1.23 + 1.0i, "(1.+1.i)"},
	{"%#.0g", 0 + 100000i, "(0.+1.e+05i)"},
	{"%#.0g", 1230000 + 0i, "(1.e+06+0.i)"},
	{"%#.4f", 1 + 1.23i, "(1.0000+1.2300i)"},
	{"%#.4e", 123 + 1i, "(1.2300e+02+1.0000e+00i)"},
	{"%#.4x", 123 + 1i, "(0x1.ec00p+06+0x1.0000p+00i)"},
	{"%#.4g", 123 + 1.23i, "(123.0+1.230i)"},
	{"%#12.5g", 0 + 100000i, "(      0.0000 +1.0000e+05i)"},
	{"%#12.5g", 1230000 - 0i, "(  1.2300e+06     +0.0000i)"},
	{"%b", 1 + 2i, "(4503599627370496p-52+4503599627370496p-51i)"},
	{"%b", complex64(1 + 2i), "(8388608p-23+8388608p-22i)"},
	// The sharp flag has no effect for binary complex format.
	{"%#b", 1 + 2i, "(4503599627370496p-52+4503599627370496p-51i)"},
	// Precision has no effect for binary complex format.
	{"%.4b", 1 + 2i, "(4503599627370496p-52+4503599627370496p-51i)"},
	{"%.4b", complex64(1 + 2i), "(8388608p-23+8388608p-22i)"},
	// complex infinites and NaNs
	{"%f", complex(posInf, posInf), "(+Inf+Infi)"},
	{"%f", complex(negInf, negInf), "(-Inf-Infi)"},
	{"%f", complex(NaN, NaN), "(NaN+NaNi)"},
	{"%.1f", complex(posInf, posInf), "(+Inf+Infi)"},
	{"% f", complex(posInf, posInf), "( Inf+Infi)"},
	{"% f", complex(negInf, negInf), "(-Inf-Infi)"},
	{"% f", complex(NaN, NaN), "( NaN+NaNi)"},
	{"%8e", complex(posInf, posInf), "(    +Inf    +Infi)"},
	{"%8x", complex(posInf, posInf), "(    +Inf    +Infi)"},
	{"% 8E", complex(posInf, posInf), "(     Inf    +Infi)"},
	{"% 8X", complex(posInf, posInf), "(     Inf    +Infi)"},
	{"%+8f", complex(negInf, negInf), "(    -Inf    -Infi)"},
	{"% +8g", complex(negInf, negInf), "(    -Inf    -Infi)"},
	{"% -8G", complex(NaN, NaN), "( NaN    +NaN    i)"},
	{"%+-8b", complex(NaN, NaN), "(+NaN    +NaN    i)"},
	// Zero padding does not apply to infinities and NaN.
	{"%08f", complex(posInf, posInf), "(    +Inf    +Infi)"},
	{"%-08g", complex(negInf, negInf), "(-Inf    -Inf    i)"},
	{"%-08G", complex(NaN, NaN), "(NaN     +NaN    i)"},

	// old test/fmt_test.go
	{"%e", 1.0, "1.000000e+00"},
	{"%e", 1234.5678e3, "1.234568e+06"},
	{"%e", 1234.5678e-8, "1.234568e-05"},
	{"%e", -7.0, "-7.000000e+00"},
	{"%e", -1e-9, "-1.000000e-09"},
	{"%f", 1234.5678e3, "1234567.800000"},
	{"%f", 1234.5678e-8, "0.000012"},
	{"%f", -7.0, "-7.000000"},
	{"%f", -1e-9, "-0.000000"},
	{"%g", 1234.5678e3, "1.2345678e+06"},
	{"%g", float32(1234.5678e3), "1.2345678e+06"},
	{"%g", 1234.5678e-8, "1.2345678e-05"},
	{"%g", -7.0, "-7"},
	{"%g", -1e-9, "-1e-09"},
	{"%g", float32(-1e-9), "-1e-09"},
	{"%E", 1.0, "1.000000E+00"},
	{"%E", 1234.5678e3, "1.234568E+06"},
	{"%E", 1234.5678e-8, "1.234568E-05"},
	{"%E", -7.0, "-7.000000E+00"},
	{"%E", -1e-9, "-1.000000E-09"},
	{"%G", 1234.5678e3, "1.2345678E+06"},
	{"%G", float32(1234.5678e3), "1.2345678E+06"},
	{"%G", 1234.5678e-8, "1.2345678E-05"},
	{"%G", -7.0, "-7"},
	{"%G", -1e-9, "-1E-09"},
	{"%G", float32(-1e-9), "-1E-09"},
	{"%20.5s", "qwertyuiop", "               qwert"},
	{"%.5s", "qwertyuiop", "qwert"},
	{"%-20.5s", "qwertyuiop", "qwert               "},
	{"%20c", 'x', "                   x"},
	{"%-20c", 'x', "x                   "},
	{"%20.6e", 1.2345e3, "        1.234500e+03"},
	{"%20.6e", 1.2345e-3, "        1.234500e-03"},
	{"%20e", 1.2345e3, "        1.234500e+03"},
	{"%20e", 1.2345e-3, "        1.234500e-03"},
	{"%20.8e", 1.2345e3, "      1.23450000e+03"},
	{"%20f", 1.23456789e3, "         1234.567890"},
	{"%20f", 1.23456789e-3, "            0.001235"},
	{"%20f", 12345678901.23456789, "  12345678901.234568"},
	{"%-20f", 1.23456789e3, "1234.567890         "},
	{"%20.8f", 1.23456789e3, "       1234.56789000"},
	{"%20.8f", 1.23456789e-3, "          0.00123457"},
	{"%g", 1.23456789e3, "1234.56789"},
	{"%g", 1.23456789e-3, "0.00123456789"},
	{"%g", 1.23456789e20, "1.23456789e+20"},

	// arrays
	{"%v", array, "[1 2 3 4 5]"},
	{"%v", iarray, "[1 hello 2.5 <nil>]"},
	{"%v", barray, "[1 2 3 4 5]"},
	{"%v", &array, "&[1 2 3 4 5]"},
	{"%v", &iarray, "&[1 hello 2.5 <nil>]"},
	{"%v", &barray, "&[1 2 3 4 5]"},

	// slices
	{"%v", slice, "[1 2 3 4 5]"},
	{"%v", islice, "[1 hello 2.5 <nil>]"},
	{"%v", bslice, "[1 2 3 4 5]"},
	{"%v", &slice, "&[1 2 3 4 5]"},
	{"%v", &islice, "&[1 hello 2.5 <nil>]"},
	{"%v", &bslice, "&[1 2 3 4 5]"},

	// byte arrays and slices with %b,%c,%d,%o,%U and %v
	{"%b", [3]byte{65, 66, 67}, "[1000001 1000010 1000011]"},
	{"%c", [3]byte{65, 66, 67}, "[A B C]"},
	{"%d", [3]byte{65, 66, 67}, "[65 66 67]"},
	{"%o", [3]byte{65, 66, 67}, "[101 102 103]"},
	{"%U", [3]byte{65, 66, 67}, "[U+0041 U+0042 U+0043]"},
	{"%v", [3]byte{65, 66, 67}, "[65 66 67]"},
	{"%v", [1]byte{123}, "[123]"},
	{"%012v", []byte{}, "[]"},
	{"%#012v", []byte{}, "[]byte{}"},
	{"%6v", []byte{1, 11, 111}, "[     1     11    111]"},
	{"%06v", []byte{1, 11, 111}, "[000001 000011 000111]"},
	{"%-6v", []byte{1, 11, 111}, "[1      11     111   ]"},
	{"%-06v", []byte{1, 11, 111}, "[1      11     111   ]"},
	{"%#v", []byte{1, 11, 111}, "[]byte{0x1, 0xb, 0x6f}"},
	{"%#6v", []byte{1, 11, 111}, "[]byte{   0x1,    0xb,   0x6f}"},
	{"%#06v", []byte{1, 11, 111}, "[]byte{0x000001, 0x00000b, 0x00006f}"},
	{"%#-6v", []byte{1, 11, 111}, "[]byte{0x1   , 0xb   , 0x6f  }"},
	{"%#-06v", []byte{1, 11, 111}, "[]byte{0x1   , 0xb   , 0x6f  }"},
	// f.space should and f.plus should not have an effect with %v.
	{"% v", []byte{1, 11, 111}, "[ 1  11  111]"},
	{"%+v", [3]byte{1, 11, 111}, "[1 11 111]"},
	{"%# -6v", []byte{1, 11, 111}, "[]byte{ 0x1  ,  0xb  ,  0x6f }"},
	{"%#+-6v", [3]byte{1, 11, 111}, "[3]uint8{0x1   , 0xb   , 0x6f  }"},
	// f.space and f.plus should have an effect with %d.
	{"% d", []byte{1, 11, 111}, "[ 1  11  111]"},
	{"%+d", [3]byte{1, 11, 111}, "[+1 +11 +111]"},
	{"%# -6d", []byte{1, 11, 111}, "[ 1      11     111  ]"},
	{"%#+-6d", [3]byte{1, 11, 111}, "[+1     +11    +111  ]"},

	// floates with %v
	{"%v", 1.2345678, "1.2345678"},
	{"%v", float32(1.2345678), "1.2345678"},

	// complexes with %v
	{"%v", 1 + 2i, "(1+2i)"},
	{"%v", complex64(1 + 2i), "(1+2i)"},

	// structs
	{"%v", A{1, 2, "a", []int{1, 2}}, `{1 2 a [1 2]}`},
	{"%+v", A{1, 2, "a", []int{1, 2}}, `{i:1 j:2 s:a x:[1 2]}`},

	// +v on structs with Stringable items
	{"%+v", B{1, 2}, `{I:<1> j:2}`},
	{"%+v", C{1, B{2, 3}}, `{i:1 B:{I:<2> j:3}}`},

	// other formats on Stringable items
	{"%s", I(23), `<23>`},
	{"%q", I(23), `"<23>"`},
	{"%x", I(23), `3c32333e`},
	{"%#x", I(23), `0x3c32333e`},
	{"%# x", I(23), `0x3c 0x32 0x33 0x3e`},
	// Stringer applies only to string formats.
	{"%d", I(23), `23`},
	// Stringer applies to the extracted value.
	{"%s", reflect.ValueOf(I(23)), `<23>`},

	// go syntax
	{"%#v", A{1, 2, "a", []int{1, 2}}, `ctxfmt.A{i:1, j:0x2, s:"a", x:[]int{1, 2}}`},
	{"%#v", new(byte), "(*uint8)(0xPTR)"},
	{"%#v", TestFmtInterface, "(func(*testing.T))(0xPTR)"},
	{"%#v", make(chan int), "(chan int)(0xPTR)"},
	{"%#v", uint64(1<<64 - 1), "0xffffffffffffffff"},
	{"%#v", 1000000000, "1000000000"},
	{"%#v", map[string]int{"a": 1}, `map[string]int{"a":1}`},
	{"%#v", map[string]B{"a": {1, 2}}, `map[string]ctxfmt.B{"a":ctxfmt.B{I:1, j:2}}`},
	{"%#v", []string{"a", "b"}, `[]string{"a", "b"}`},
	{"%#v", SI{}, `ctxfmt.SI{I:interface {}(nil)}`},
	{"%#v", []int(nil), `[]int(nil)`},
	{"%#v", []int{}, `[]int{}`},
	{"%#v", array, `[5]int{1, 2, 3, 4, 5}`},
	{"%#v", &array, `&[5]int{1, 2, 3, 4, 5}`},
	{"%#v", iarray, `[4]interface {}{1, "hello", 2.5, interface {}(nil)}`},
	{"%#v", &iarray, `&[4]interface {}{1, "hello", 2.5, interface {}(nil)}`},
	{"%#v", map[int]byte(nil), `map[int]uint8(nil)`},
	{"%#v", map[int]byte{}, `map[int]uint8{}`},
	{"%#v", "foo", `"foo"`},
	{"%#v", barray, `[5]ctxfmt.renamedUint8{0x1, 0x2, 0x3, 0x4, 0x5}`},
	{"%#v", bslice, `[]ctxfmt.renamedUint8{0x1, 0x2, 0x3, 0x4, 0x5}`},
	{"%#v", []int32(nil), "[]int32(nil)"},
	{"%#v", 1.2345678, "1.2345678"},
	{"%#v", float32(1.2345678), "1.2345678"},

	// Whole number floats are printed without decimals. See Issue 27634.
	{"%#v", 1.0, "1"},
	{"%#v", 1000000.0, "1e+06"},
	{"%#v", float32(1.0), "1"},
	{"%#v", float32(1000000.0), "1e+06"},

	// Only print []byte and []uint8 as type []byte if they appear at the top level.
	{"%#v", []byte(nil), "[]byte(nil)"},
	{"%#v", []uint8(nil), "[]byte(nil)"},
	{"%#v", []byte{}, "[]byte{}"},
	{"%#v", []uint8{}, "[]byte{}"},
	{"%#v", reflect.ValueOf([]byte{}), "[]uint8{}"},
	{"%#v", reflect.ValueOf([]uint8{}), "[]uint8{}"},
	{"%#v", &[]byte{}, "&[]uint8{}"},
	{"%#v", &[]byte{}, "&[]uint8{}"},
	{"%#v", [3]byte{}, "[3]uint8{0x0, 0x0, 0x0}"},
	{"%#v", [3]uint8{}, "[3]uint8{0x0, 0x0, 0x0}"},

	// slices with other formats
	{"%#x", []int{1, 2, 15}, `[0x1 0x2 0xf]`},
	{"%x", []int{1, 2, 15}, `[1 2 f]`},
	{"%d", []int{1, 2, 15}, `[1 2 15]`},
	{"%d", []byte{1, 2, 15}, `[1 2 15]`},
	{"%q", []string{"a", "b"}, `["a" "b"]`},
	{"% 02x", []byte{1}, "01"},
	{"% 02x", []byte{1, 2, 3}, "01 02 03"},

	// Padding with byte slices.
	{"%2x", []byte{}, "  "},
	{"%#2x", []byte{}, "  "},
	{"% 02x", []byte{}, "00"},
	{"%# 02x", []byte{}, "00"},
	{"%-2x", []byte{}, "  "},
	{"%-02x", []byte{}, "  "},
	{"%8x", []byte{0xab}, "      ab"},
	{"% 8x", []byte{0xab}, "      ab"},
	{"%#8x", []byte{0xab}, "    0xab"},
	{"%# 8x", []byte{0xab}, "    0xab"},
	{"%08x", []byte{0xab}, "000000ab"},
	{"% 08x", []byte{0xab}, "000000ab"},
	{"%#08x", []byte{0xab}, "00000xab"},
	{"%# 08x", []byte{0xab}, "00000xab"},
	{"%10x", []byte{0xab, 0xcd}, "      abcd"},
	{"% 10x", []byte{0xab, 0xcd}, "     ab cd"},
	{"%#10x", []byte{0xab, 0xcd}, "    0xabcd"},
	{"%# 10x", []byte{0xab, 0xcd}, " 0xab 0xcd"},
	{"%010x", []byte{0xab, 0xcd}, "000000abcd"},
	{"% 010x", []byte{0xab, 0xcd}, "00000ab cd"},
	{"%#010x", []byte{0xab, 0xcd}, "00000xabcd"},
	{"%# 010x", []byte{0xab, 0xcd}, "00xab 0xcd"},
	{"%-10X", []byte{0xab}, "AB        "},
	{"% -010X", []byte{0xab}, "AB        "},
	{"%#-10X", []byte{0xab, 0xcd}, "0XABCD    "},
	{"%# -010X", []byte{0xab, 0xcd}, "0XAB 0XCD "},
	// Same for strings
	{"%2x", "", "  "},
	{"%#2x", "", "  "},
	{"% 02x", "", "00"},
	{"%# 02x", "", "00"},
	{"%-2x", "", "  "},
	{"%-02x", "", "  "},
	{"%8x", "\xab", "      ab"},
	{"% 8x", "\xab", "      ab"},
	{"%#8x", "\xab", "    0xab"},
	{"%# 8x", "\xab", "    0xab"},
	{"%08x", "\xab", "000000ab"},
	{"% 08x", "\xab", "000000ab"},
	{"%#08x", "\xab", "00000xab"},
	{"%# 08x", "\xab", "00000xab"},
	{"%10x", "\xab\xcd", "      abcd"},
	{"% 10x", "\xab\xcd", "     ab cd"},
	{"%#10x", "\xab\xcd", "    0xabcd"},
	{"%# 10x", "\xab\xcd", " 0xab 0xcd"},
	{"%010x", "\xab\xcd", "000000abcd"},
	{"% 010x", "\xab\xcd", "00000ab cd"},
	{"%#010x", "\xab\xcd", "00000xabcd"},
	{"%# 010x", "\xab\xcd", "00xab 0xcd"},
	{"%-10X", "\xab", "AB        "},
	{"% -010X", "\xab", "AB        "},
	{"%#-10X", "\xab\xcd", "0XABCD    "},
	{"%# -010X", "\xab\xcd", "0XAB 0XCD "},

	// renamings
	{"%v", renamedBool(true), "true"},
	{"%d", renamedBool(true), "%!d(ctxfmt.renamedBool=true)"},
	{"%o", renamedInt(8), "10"},
	{"%d", renamedInt8(-9), "-9"},
	{"%v", renamedInt16(10), "10"},
	{"%v", renamedInt32(-11), "-11"},
	{"%X", renamedInt64(255), "FF"},
	{"%v", renamedUint(13), "13"},
	{"%o", renamedUint8(14), "16"},
	{"%X", renamedUint16(15), "F"},
	{"%d", renamedUint32(16), "16"},
	{"%X", renamedUint64(17), "11"},
	{"%o", renamedUintptr(18), "22"},
	{"%x", renamedString("thing"), "7468696e67"},
	{"%d", renamedBytes([]byte{1, 2, 15}), `[1 2 15]`},
	{"%q", renamedBytes([]byte("hello")), `"hello"`},
	{"%x", []renamedUint8{'h', 'e', 'l', 'l', 'o'}, "68656c6c6f"},
	{"%X", []renamedUint8{'h', 'e', 'l', 'l', 'o'}, "68656C6C6F"},
	{"%s", []renamedUint8{'h', 'e', 'l', 'l', 'o'}, "hello"},
	{"%q", []renamedUint8{'h', 'e', 'l', 'l', 'o'}, `"hello"`},
	{"%v", renamedFloat32(22), "22"},
	{"%v", renamedFloat64(33), "33"},
	{"%v", renamedComplex64(3 + 4i), "(3+4i)"},
	{"%v", renamedComplex128(4 - 3i), "(4-3i)"},

	// Formatter
	{"%x", F(1), "<x=F(1)>"},
	{"%x", G(2), "2"},
	{"%+v", S{F(4), G(5)}, "{F:<v=F(4)> G:5}"},

	// GoStringer
	{"%#v", G(6), "GoString(6)"},
	{"%#v", S{F(7), G(8)}, "ctxfmt.S{F:<v=F(7)>, G:GoString(8)}"},

	// %T
	{"%T", byte(0), "uint8"},
	{"%T", reflect.ValueOf(nil), "reflect.Value"},
	{"%T", (4 - 3i), "complex128"},
	{"%T", renamedComplex128(4 - 3i), "ctxfmt.renamedComplex128"},
	{"%T", intVar, "int"},
	{"%6T", &intVar, "  *int"},
	{"%10T", nil, "     <nil>"},
	{"%-10T", nil, "<nil>     "},

	// %p with pointers
	{"%p", (*int)(nil), "0x0"},
	{"%#p", (*int)(nil), "0"},
	{"%p", &intVar, "0xPTR"},
	{"%#p", &intVar, "PTR"},
	{"%p", &array, "0xPTR"},
	{"%p", &slice, "0xPTR"},
	{"%8.2p", (*int)(nil), "    0x00"},
	{"%-20.16p", &intVar, "0xPTR  "},
	// %p on non-pointers
	{"%p", make(chan int), "0xPTR"},
	{"%p", make(map[int]int), "0xPTR"},
	{"%p", func() {}, "0xPTR"},
	{"%p", 27, "%!p(int=27)"},  // not a pointer at all
	{"%p", nil, "%!p(<nil>)"},  // nil on its own has no type ...
	{"%#p", nil, "%!p(<nil>)"}, // ... and hence is not a pointer type.
	// pointers with specified base
	{"%b", &intVar, "PTR_b"},
	{"%d", &intVar, "PTR_d"},
	{"%o", &intVar, "PTR_o"},
	{"%x", &intVar, "PTR_x"},
	{"%X", &intVar, "PTR_X"},
	// %v on pointers
	{"%v", nil, "<nil>"},
	{"%#v", nil, "<nil>"},
	{"%v", (*int)(nil), "<nil>"},
	{"%#v", (*int)(nil), "(*int)(nil)"},
	{"%v", &intVar, "0xPTR"},
	{"%#v", &intVar, "(*int)(0xPTR)"},
	{"%8.2v", (*int)(nil), "   <nil>"},
	{"%-20.16v", &intVar, "0xPTR  "},
	// string method on pointer
	{"%s", &pValue, "String(p)"}, // String method...
	{"%p", &pValue, "0xPTR"},     // ... is not called with %p.

	// %d on Stringer should give integer if possible
	{"%s", time.Time{}.Month(), "January"},
	{"%d", time.Time{}.Month(), "1"},

	// erroneous things
	// ctxfmt does not handle 'EXTRA', but returns the list of extra arguments
	// for the caller to handle.
	//
	// {"", nil, "%!(EXTRA <nil>)"},
	// {"", 2, "%!(EXTRA int=2)"},
	// {"no args", "hello", "no args%!(EXTRA string=hello)"},
	{"%s %", "hello", "hello %!(NOVERB)"},
	{"%s %.2", "hello", "hello %!(NOVERB)"},
	// {"%017091901790959340919092959340919017929593813360", 0, "%!(NOVERB)%!(EXTRA int=0)"},
	// {"%184467440737095516170v", 0, "%!(NOVERB)%!(EXTRA int=0)"},
	// Extra argument errors should format without flags set.
	// {"%010.2", "12345", "%!(NOVERB)%!(EXTRA string=12345)"},

	// Test that maps with non-reflexive keys print all keys and values.
	{"%v", map[float64]int{NaN: 1, NaN: 1}, "map[NaN:1 NaN:1]"},

	// Comparison of padding rules with C printf.
	/*
		C program:
		#include <stdio.h>

		char *format[] = {
			"[%.2f]",
			"[% .2f]",
			"[%+.2f]",
			"[%7.2f]",
			"[% 7.2f]",
			"[%+7.2f]",
			"[% +7.2f]",
			"[%07.2f]",
			"[% 07.2f]",
			"[%+07.2f]",
			"[% +07.2f]"
		};

		int main(void) {
			int i;
			for(i = 0; i < 11; i++) {
				printf("%s: ", format[i]);
				printf(format[i], 1.0);
				printf(" ");
				printf(format[i], -1.0);
				printf("\n");
			}
		}

		Output:
			[%.2f]: [1.00] [-1.00]
			[% .2f]: [ 1.00] [-1.00]
			[%+.2f]: [+1.00] [-1.00]
			[%7.2f]: [   1.00] [  -1.00]
			[% 7.2f]: [   1.00] [  -1.00]
			[%+7.2f]: [  +1.00] [  -1.00]
			[% +7.2f]: [  +1.00] [  -1.00]
			[%07.2f]: [0001.00] [-001.00]
			[% 07.2f]: [ 001.00] [-001.00]
			[%+07.2f]: [+001.00] [-001.00]
			[% +07.2f]: [+001.00] [-001.00]

	*/
	{"%.2f", 1.0, "1.00"},
	{"%.2f", -1.0, "-1.00"},
	{"% .2f", 1.0, " 1.00"},
	{"% .2f", -1.0, "-1.00"},
	{"%+.2f", 1.0, "+1.00"},
	{"%+.2f", -1.0, "-1.00"},
	{"%7.2f", 1.0, "   1.00"},
	{"%7.2f", -1.0, "  -1.00"},
	{"% 7.2f", 1.0, "   1.00"},
	{"% 7.2f", -1.0, "  -1.00"},
	{"%+7.2f", 1.0, "  +1.00"},
	{"%+7.2f", -1.0, "  -1.00"},
	{"% +7.2f", 1.0, "  +1.00"},
	{"% +7.2f", -1.0, "  -1.00"},
	{"%07.2f", 1.0, "0001.00"},
	{"%07.2f", -1.0, "-001.00"},
	{"% 07.2f", 1.0, " 001.00"},
	{"% 07.2f", -1.0, "-001.00"},
	{"%+07.2f", 1.0, "+001.00"},
	{"%+07.2f", -1.0, "-001.00"},
	{"% +07.2f", 1.0, "+001.00"},
	{"% +07.2f", -1.0, "-001.00"},

	// Complex numbers: exhaustively tested in TestComplexFormatting.
	{"%7.2f", 1 + 2i, "(   1.00  +2.00i)"},
	{"%+07.2f", -1 - 2i, "(-001.00-002.00i)"},

	// Use spaces instead of zero if padding to the right.
	{"%0-5s", "abc", "abc  "},
	{"%-05.1f", 1.0, "1.0  "},

	// float and complex formatting should not change the padding width
	// for other elements. See issue 14642.
	{"%06v", []interface{}{+10.0, 10}, "[000010 000010]"},
	{"%06v", []interface{}{-10.0, 10}, "[-00010 000010]"},
	{"%06v", []interface{}{+10.0 + 10i, 10}, "[(000010+00010i) 000010]"},
	{"%06v", []interface{}{-10.0 + 10i, 10}, "[(-00010+00010i) 000010]"},

	// integer formatting should not alter padding for other elements.
	{"%03.6v", []interface{}{1, 2.0, "x"}, "[000001 002 00x]"},
	{"%03.0v", []interface{}{0, 2.0, "x"}, "[    002 000]"},

	// Complex fmt used to leave the plus flag set for future entries in the array
	// causing +2+0i and +3+0i instead of 2+0i and 3+0i.
	{"%v", []complex64{1, 2, 3}, "[(1+0i) (2+0i) (3+0i)]"},
	{"%v", []complex128{1, 2, 3}, "[(1+0i) (2+0i) (3+0i)]"},

	// Incomplete format specification caused crash.
	{"%.", 3, "%!(NOVERB)(int=3)"},

	// Padding for complex numbers. Has been bad, then fixed, then bad again.
	{"%+10.2f", +104.66 + 440.51i, "(   +104.66   +440.51i)"},
	{"%+10.2f", -104.66 + 440.51i, "(   -104.66   +440.51i)"},
	{"%+10.2f", +104.66 - 440.51i, "(   +104.66   -440.51i)"},
	{"%+10.2f", -104.66 - 440.51i, "(   -104.66   -440.51i)"},
	{"%+010.2f", +104.66 + 440.51i, "(+000104.66+000440.51i)"},
	{"%+010.2f", -104.66 + 440.51i, "(-000104.66+000440.51i)"},
	{"%+010.2f", +104.66 - 440.51i, "(+000104.66-000440.51i)"},
	{"%+010.2f", -104.66 - 440.51i, "(-000104.66-000440.51i)"},

	// []T where type T is a byte with a Stringer method.
	{"%v", byteStringerSlice, "[X X X X X]"},
	{"%s", byteStringerSlice, "hello"},
	{"%q", byteStringerSlice, "\"hello\""},
	{"%x", byteStringerSlice, "68656c6c6f"},
	{"%X", byteStringerSlice, "68656C6C6F"},
	{"%#v", byteStringerSlice, "[]ctxfmt.byteStringer{0x68, 0x65, 0x6c, 0x6c, 0x6f}"},

	// And the same for Formatter.
	{"%v", byteFormatterSlice, "[X X X X X]"},
	{"%s", byteFormatterSlice, "hello"},
	{"%q", byteFormatterSlice, "\"hello\""},
	{"%x", byteFormatterSlice, "68656c6c6f"},
	{"%X", byteFormatterSlice, "68656C6C6F"},
	// This next case seems wrong, but the docs say the Formatter wins here.
	{"%#v", byteFormatterSlice, "[]ctxfmt.byteFormatter{X, X, X, X, X}"},

	// pp.WriteString
	{"%s", writeStringFormatter(""), "******"},
	{"%s", writeStringFormatter("xyz"), "***xyz***"},
	{"%s", writeStringFormatter("⌘/⌘"), "***⌘/⌘***"},

	// reflect.Value handled specially in Go 1.5, making it possible to
	// see inside non-exported fields (which cannot be accessed with Interface()).
	// Issue 8965.
	{"%v", reflect.ValueOf(A{}).Field(0).String(), "<int Value>"}, // Equivalent to the old way.
	{"%v", reflect.ValueOf(A{}).Field(0), "0"},                    // Sees inside the field.

	// verbs apply to the extracted value too.
	{"%s", reflect.ValueOf("hello"), "hello"},
	{"%q", reflect.ValueOf("hello"), `"hello"`},
	{"%#04x", reflect.ValueOf(256), "0x0100"},

	// invalid reflect.Value doesn't crash.
	{"%v", reflect.Value{}, "<invalid reflect.Value>"},
	{"%v", &reflect.Value{}, "<invalid Value>"},
	{"%v", SI{reflect.Value{}}, "{<invalid Value>}"},

	/*
		// Tests to check that not supported verbs generate an error string.
		{"%☠", nil, "%!☠(INVALID)(<nil>)"},
		{"%☠", interface{}(nil), "%!☠(INVALID)(<nil>)"},
		{"%☠", int(0), "%!☠(INVALID)(int=0)"},
		{"%☠", uint(0), "%!☠(INVALID)(uint=0)"},
		{"%☠", []byte{0, 1}, "[%!☠(INVALID)(uint8=0) %!☠(INVALID)(uint8=1)]"},
		{"%☠", []uint8{0, 1}, "[%!☠(uint8=0) %!☠(uint8=1)]"},
		{"%☠", [1]byte{0}, "[%!☠(INVALID)(uint8=0)]"},
		{"%☠", [1]uint8{0}, "[%!☠(uint8=0)]"},
		{"%☠", "hello", "%!☠(string=hello)"},
		{"%☠", 1.2345678, "%!☠(float64=1.2345678)"},
		{"%☠", float32(1.2345678), "%!☠(float32=1.2345678)"},
		{"%☠", 1.2345678 + 1.2345678i, "%!☠(complex128=(1.2345678+1.2345678i))"},
		{"%☠", complex64(1.2345678 + 1.2345678i), "%!☠(complex64=(1.2345678+1.2345678i))"},
		{"%☠", &intVar, "%!☠(*int=0xPTR)"},
		{"%☠", make(chan int), "%!☠(chan int=0xPTR)"},
		{"%☠", func() {}, "%!☠(func()=0xPTR)"},
		{"%☠", reflect.ValueOf(renamedInt(0)), "%!☠(ctxfmt.renamedInt=0)"},
		{"%☠", SI{renamedInt(0)}, "{%!☠(ctxfmt.renamedInt=0)}"},
		{"%☠", &[]interface{}{I(1), G(2)}, "&[%!☠(ctxfmt.I=1) %!☠(ctxfmt.G=2)]"},
		{"%☠", SI{&[]interface{}{I(1), G(2)}}, "{%!☠(*[]interface {}=&[1 2])}"},
		{"%☠", reflect.Value{}, "<invalid reflect.Value>"},
		{"%☠", map[float64]int{NaN: 1}, "map[%!☠(float64=NaN):%!☠(int=1)]"},
	*/
}

func TestSprintf(t *testing.T) {
	for i, tt := range fmtTests {
		t.Run(fmt.Sprintf("%d: %v -> %v", i, tt.fmt, tt.out), func(t *testing.T) {
			s, _ := Sprintf(nil, tt.fmt, tt.val)
			i := strings.Index(tt.out, "PTR")
			if i >= 0 && i < len(s) {
				var pattern, chars string
				switch {
				case strings.HasPrefix(tt.out[i:], "PTR_b"):
					pattern = "PTR_b"
					chars = "01"
				case strings.HasPrefix(tt.out[i:], "PTR_o"):
					pattern = "PTR_o"
					chars = "01234567"
				case strings.HasPrefix(tt.out[i:], "PTR_d"):
					pattern = "PTR_d"
					chars = "0123456789"
				case strings.HasPrefix(tt.out[i:], "PTR_x"):
					pattern = "PTR_x"
					chars = "0123456789abcdef"
				case strings.HasPrefix(tt.out[i:], "PTR_X"):
					pattern = "PTR_X"
					chars = "0123456789ABCDEF"
				default:
					pattern = "PTR"
					chars = "0123456789abcdefABCDEF"
				}
				p := s[:i] + pattern
				for j := i; j < len(s); j++ {
					if !strings.ContainsRune(chars, rune(s[j])) {
						p += s[j:]
						break
					}
				}
				s = p
			}
			if s != tt.out {
				if _, ok := tt.val.(string); ok {
					// Don't requote the already-quoted strings.
					// It's too confusing to read the errors.
					t.Errorf("Sprintf(%q, %q) = <%s> want <%s>", tt.fmt, tt.val, s, tt.out)
				} else {
					t.Errorf("Sprintf(%q, %v) = %q want %q", tt.fmt, tt.val, s, tt.out)
				}
			}
		})
	}
}

// zeroFill generates zero-filled strings of the specified width. The length
// of the suffix (but not the prefix) is compensated for in the width calculation.
func zeroFill(prefix string, width int, suffix string) string {
	return prefix + strings.Repeat("0", width-len(suffix)) + suffix
}

func TestFunny(t *testing.T) {
	var str = "%."
	t.Log(fmt.Sprintf(str, 3))
}
