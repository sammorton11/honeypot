package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func convertAddressToCountry(address string) string {
	return ""
}

func hashMessage(buf []byte, n int) string {
	hash := sha512.Sum512(buf[:n])
	stringHash := hex.EncodeToString(hash[:])

	return stringHash
}
func trimIP(ip string) string {

	//temp := strings.split(ip, ".")
	//sb := strings.builder{}
	//sb.grow(len(ip))
	// create a slice up until half of the ip address or something

	//sb.writestring(temp[:])

	temp := []rune(ip)
	half := (len(temp) - 1) / 2
	_temp := temp[:half]
	newIP := string(_temp)
	hash := md5.Sum([]byte(newIP))
	hash2 := sha256.Sum256(hash[:])

	return string(hash2[:])
}
