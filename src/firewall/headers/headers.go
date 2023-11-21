package headers

const ipHeaderLength = 20
const tcpDataOffset = 12

func GetUDPHeaderLength() int {
	return ipHeaderLength + 8
}

func GetTCPHeaderLength(payload []byte) int {
	// https://en.wikipedia.org/wiki/Transmission_Control_Protocol
	return ipHeaderLength + ((int(payload[ipHeaderLength+tcpDataOffset])>>4)*(ipHeaderLength+tcpDataOffset))/8
}
