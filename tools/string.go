package tools

func TextParser(text string) []string {
	const (
		FLAG_NONE   = 0
		FLAG_STRING = 1 // start and end with "
		FLAG_RAW    = 2
	)

	var result []string
	var perMessage []byte
	var curFlag int = FLAG_NONE
	var prevFlag int = FLAG_NONE

	for i := 0; i < len(text); i++ {
		curByte := text[i]
		perDone := false

		switch curFlag {
		case FLAG_NONE:
			switch curByte {
			case '\\':
				prevFlag = curFlag
				curFlag = FLAG_RAW
			case '"':
				curFlag = FLAG_STRING
			case ' ':
				if len(perMessage) != 0 {
					perDone = true
				}
			default:
				perMessage = append(perMessage, curByte)
			}
		case FLAG_STRING:
			switch curByte {
			case '\\':
				prevFlag = curFlag
				curFlag = FLAG_RAW
			case '"':
				curFlag = FLAG_RAW
				perDone = true
			default:
				perMessage = append(perMessage, curByte)
			}
		case FLAG_RAW:
			switch curByte {
			case 't':
				curByte = '\t'
			case 'n':
				curByte = '\n'
			case 'b':
				curByte = '\b'
			case 'a':
				curByte = '\a'
			}
			perMessage = append(perMessage, curByte)
			curFlag = prevFlag // to deal like \" in a "string"
		}

		if perDone {
			result = append(result, string(perMessage))
			perMessage = []byte{}
		}
	}

	// clear what we left
	// for a complete string left, the flag should be FLAG_RAW
	// FLAG_STRING only occur in half of process of parsing a string
	if len(perMessage) != 0 && curFlag != FLAG_STRING {
		result = append(result, string(perMessage))
	}

	return result
}
