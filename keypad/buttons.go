package keypad

const (
	// Each button on the keypad is mapped to a single bit in a 64bit integer.
	// This allows for all possible combinations of keypresses and releases to be
	// tracked. The base buttons are labeled with their default label, though
	// some convnience combos exist, such as the arrow buttons.
	BtnBacktick = 1 << iota
	Btn1
	Btn2
	Btn3
	Btn4
	Btn5
	Btn6
	Btn7
	Btn8
	Btn9
	Btn0
	BtnUnderscore
	BtnEqual
	BtnBackspace
	// Row 2
	BtnTab
	BtnQ
	BtnW
	BtnE
	BtnR
	BtnT
	BtnY
	BtnU
	BtnI
	BtnO
	BtnP
	BtnBraceLeft
	BtnBraceRight
	BtnBackslash
	// Row 3
	BtnFn
	BtnShift
	BtnA
	BtnS
	BtnD
	BtnF
	BtnG
	BtnH
	BtnJ
	BtnK
	BtnL
	BtnSemicolon
	BtnQuote
	BtnEnter
	// Row 4
	BtnCtrl
	BtnOpt
	BtnAlt
	BtnZ
	BtnX
	BtnC
	BtnV
	BtnB
	BtnN
	BtnM
	BtnComma
	BtnPeriod
	BtnSlash
	BtnSpace

	// Convenience Aliases
	BtnEsc   = BtnFn | BtnBacktick
	BtnDel   = BtnFn | BtnBackspace
	BtnUp    = BtnFn | BtnSemicolon
	BtnDown  = BtnFn | BtnPeriod
	BtnRight = BtnFn | BtnSlash
	BtnLeft  = BtnFn | BtnComma

	// Implement F1-F12 keys as Fn+Numbers
	BtnF1  = BtnFn | Btn1
	BtnF2  = BtnFn | Btn2
	BtnF3  = BtnFn | Btn3
	BtnF4  = BtnFn | Btn4
	BtnF5  = BtnFn | Btn5
	BtnF6  = BtnFn | Btn6
	BtnF7  = BtnFn | Btn7
	BtnF8  = BtnFn | Btn8
	BtnF9  = BtnFn | Btn9
	BtnF10 = BtnFn | Btn0
	BtnF11 = BtnFn | BtnUnderscore
	BtnF12 = BtnFn | BtnEqual
)

const (
	// BtnSpecialMask is a bitmask for all the buttons that should not produce
	// a character (if using (*Device).) when pressed on their own.
	BtnSpecialMask = (BtnCtrl | BtnFn | BtnOpt | BtnShift | BtnAlt)
)

// ScancodeToBytes maps the pressed buttons to a character.
// This is used by the (*Device).WriteByteCallback() method.
// If you want a key-combo to result in a character sequence or
// want to override a character sequence, you can modify this
// map to do so.
//
// E.G., To have Ctrl-m to send a newline:
// keypad.ScancodeToBytes[keypad.BtnCtrl|keypad.BtnM] = []byte{'\n'}
var (
	ScancodeToBytes = map[int64][]byte{
		// Row 1
		BtnBacktick:   []byte{'`'},
		Btn1:          []byte{'1'},
		Btn2:          []byte{'2'},
		Btn3:          []byte{'3'},
		Btn4:          []byte{'4'},
		Btn5:          []byte{'5'},
		Btn6:          []byte{'6'},
		Btn7:          []byte{'7'},
		Btn8:          []byte{'8'},
		Btn9:          []byte{'9'},
		Btn0:          []byte{'0'},
		BtnUnderscore: []byte{'_'},
		BtnEqual:      []byte{'='},
		BtnBackspace:  []byte{'\b'},
		// Row 2
		BtnTab:        []byte{'\t'},
		BtnQ:          []byte{'q'},
		BtnW:          []byte{'w'},
		BtnE:          []byte{'e'},
		BtnR:          []byte{'r'},
		BtnT:          []byte{'t'},
		BtnY:          []byte{'y'},
		BtnU:          []byte{'u'},
		BtnI:          []byte{'i'},
		BtnO:          []byte{'o'},
		BtnP:          []byte{'p'},
		BtnBraceLeft:  []byte{'['},
		BtnBraceRight: []byte{']'},
		BtnBackslash:  []byte{'\\'},
		// Row 3
		BtnA:         []byte{'a'},
		BtnS:         []byte{'s'},
		BtnD:         []byte{'d'},
		BtnF:         []byte{'f'},
		BtnG:         []byte{'g'},
		BtnH:         []byte{'h'},
		BtnJ:         []byte{'j'},
		BtnK:         []byte{'k'},
		BtnL:         []byte{'l'},
		BtnSemicolon: []byte{';'},
		BtnQuote:     []byte{'\''},
		BtnEnter:     []byte{'\n'},
		// Row 4
		BtnZ:     []byte{'z'},
		BtnX:     []byte{'x'},
		BtnC:     []byte{'c'},
		BtnV:     []byte{'v'},
		BtnB:     []byte{'b'},
		BtnN:     []byte{'n'},
		BtnM:     []byte{'m'},
		BtnUp:    []byte{0x1b, '[', 'A'},
		BtnDown:  []byte{0x1b, '[', 'B'},
		BtnRight: []byte{0x1b, '[', 'C'},
		BtnLeft:  []byte{0x1b, '[', 'D'},
		BtnSpace: []byte{' '},

		// With Shift
		// Row 1
		BtnShift | BtnBacktick:   []byte{'~'},
		BtnShift | Btn1:          []byte{'!'},
		BtnShift | Btn2:          []byte{'@'},
		BtnShift | Btn3:          []byte{'#'},
		BtnShift | Btn4:          []byte{'$'},
		BtnShift | Btn5:          []byte{'%'},
		BtnShift | Btn6:          []byte{'^'},
		BtnShift | Btn7:          []byte{'&'},
		BtnShift | Btn8:          []byte{'*'},
		BtnShift | Btn9:          []byte{'('},
		BtnShift | Btn0:          []byte{')'},
		BtnShift | BtnUnderscore: []byte{'-'},
		BtnShift | BtnEqual:      []byte{'+'},
		BtnShift | BtnBackspace:  []byte{'\b'},
		// Row 2
		BtnShift | BtnTab:        []byte{0x1b, '[', 'Z'},
		BtnShift | BtnQ:          []byte{'Q'},
		BtnShift | BtnW:          []byte{'W'},
		BtnShift | BtnE:          []byte{'E'},
		BtnShift | BtnR:          []byte{'R'},
		BtnShift | BtnT:          []byte{'T'},
		BtnShift | BtnY:          []byte{'Y'},
		BtnShift | BtnU:          []byte{'U'},
		BtnShift | BtnI:          []byte{'I'},
		BtnShift | BtnO:          []byte{'O'},
		BtnShift | BtnP:          []byte{'P'},
		BtnShift | BtnBraceLeft:  []byte{'{'},
		BtnShift | BtnBraceRight: []byte{'}'},
		BtnShift | BtnBackslash:  []byte{'|'},
		// Row 3
		BtnShift | BtnA:         []byte{'A'},
		BtnShift | BtnS:         []byte{'S'},
		BtnShift | BtnD:         []byte{'D'},
		BtnShift | BtnF:         []byte{'F'},
		BtnShift | BtnG:         []byte{'G'},
		BtnShift | BtnH:         []byte{'H'},
		BtnShift | BtnJ:         []byte{'J'},
		BtnShift | BtnK:         []byte{'K'},
		BtnShift | BtnL:         []byte{'L'},
		BtnShift | BtnSemicolon: []byte{':'},
		BtnShift | BtnQuote:     []byte{'"'},
		BtnShift | BtnEnter:     []byte{'\n'},
		// Row 4
		BtnShift | BtnZ:     []byte{'Z'},
		BtnShift | BtnX:     []byte{'X'},
		BtnShift | BtnC:     []byte{'C'},
		BtnShift | BtnV:     []byte{'V'},
		BtnShift | BtnB:     []byte{'B'},
		BtnShift | BtnN:     []byte{'N'},
		BtnShift | BtnM:     []byte{'M'},
		BtnShift | BtnUp:    []byte{0x1b, '[', '1', ';', '2', 'A'}, // Up button
		BtnShift | BtnDown:  []byte{0x1b, '[', '1', ';', '2', 'B'}, // Down button
		BtnShift | BtnRight: []byte{0x1b, '[', '1', ';', '2', 'C'}, // Right button
		BtnShift | BtnLeft:  []byte{0x1b, '[', '1', ';', '2', 'D'}, // Left button
		BtnShift | BtnSpace: []byte{' '},

		// Other Fn buttons
		BtnEsc: []byte{'\x1b'},

		// Ctrl Combos
		// NUL Null
		BtnCtrl | BtnShift | Btn2: []byte{0x00},
		// SOH Start of Heading
		BtnCtrl | BtnA: []byte{0x01},
		// SOT Start of Text
		BtnCtrl | BtnB: []byte{0x02},
		// ETX End of Text (interrupt)
		BtnCtrl | BtnC: []byte{0x03},
		// EOT End of Transmission (EOF)
		BtnCtrl | BtnD: []byte{0x04},
		// ENQ Eqnuiry
		BtnCtrl | BtnE: []byte{0x05},
		// ACK Acknowledge
		BtnCtrl | BtnF: []byte{0x06},
		// BEL Bell
		BtnCtrl | BtnG: []byte{0x07},
		// BS Backspace
		BtnCtrl | BtnH: []byte{0x08},
		// HT Horizontal Tab
		BtnCtrl | BtnI: []byte{0x09},
		// LF Line Feed
		BtnCtrl | BtnJ: []byte{0x0A},
		// VT Vertical Tab
		BtnCtrl | BtnK: []byte{0x0B},
		// FF Form Feed
		BtnCtrl | BtnL: []byte{0x0C},
		// CR Carriage Feed (Carriage Return)
		BtnCtrl | BtnM: []byte{0x0D},
		// SO Shift Out
		BtnCtrl | BtnN: []byte{0x0E},
		// SI Shift In
		BtnCtrl | BtnO: []byte{0x0F},
		// DLE Data Line Escape
		BtnCtrl | BtnP: []byte{0x10},
		// DC1 Device Control 1
		BtnCtrl | BtnQ: []byte{0x11},
		// DC2 Device Control 2
		BtnCtrl | BtnR: []byte{0x12},
		// DC3 Device Control 3
		BtnCtrl | BtnS: []byte{0x13},
		// DC4 Device Control 4
		BtnCtrl | BtnT: []byte{0x14},
		// NAK Negative Acknowledge
		BtnCtrl | BtnU: []byte{0x15},
		// SYN Synchronous Idle
		BtnCtrl | BtnV: []byte{0x16},
		// ETB End of Transmission Block
		BtnCtrl | BtnW: []byte{0x17},
		// CAN Cancel
		BtnCtrl | BtnX: []byte{0x18},
		// EM End of Medium
		BtnCtrl | BtnY: []byte{0x19},
		// SUB Substitute
		BtnCtrl | BtnZ: []byte{0x1A},
		// ESC Escape
		BtnCtrl | BtnBraceLeft: []byte{0x1B},
		// FS File Separator
		BtnCtrl | BtnBackslash: []byte{0x1C},
		// GS Group Separator
		BtnCtrl | BtnBraceRight: []byte{0x1D},
		// RS Record Separator
		BtnCtrl | BtnShift | Btn6: []byte{0x1E},
		// US Unit Separator
		BtnCtrl | BtnUnderscore: []byte{0x1F},

		//TODO: add Alt+<> combos
		// Seems like this just sends ESC + button pressed
		// Do this programmatically? Mask off Alt bit and send keys
		// pressed prefixed with ESC?
	}
)
