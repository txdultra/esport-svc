// Copyright 2013 Julien Schmidt. All rights reserved.
// http://www.julienschmidt.com
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package sphinxql

import (
	"crypto/sha1"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

// Logger
var (
	errLog *log.Logger
)

func init() {
	errLog = log.New(os.Stderr, "[SphinxQL] ", log.Ldate|log.Ltime|log.Lshortfile)

	dsnPattern = regexp.MustCompile(
		`^(?:(?P<user>.*?)(?::(?P<passwd>.*))?@)?` + // [user[:password]@]
			`(?:(?P<net>[^\(]*)(?:\((?P<addr>[^\)]*)\))?)?` + // [net[(addr)]]
			`\/(?P<dbname>.*?)` + // /dbname
			`(?:\?(?P<params>[^\?]*))?$`) // [?param1=value1&paramN=valueN]
}

// Data Source Name Parser
var dsnPattern *regexp.Regexp

func parseDSN(dsn string) *config {
	cfg := new(config)
	cfg.params = make(map[string]string)

	matches := dsnPattern.FindStringSubmatch(dsn)
	names := dsnPattern.SubexpNames()

	for i, match := range matches {
		switch names[i] {
		case "user":
			cfg.user = match
		case "passwd":
			cfg.passwd = match
		case "net":
			cfg.net = match
		case "addr":
			cfg.addr = match
		case "dbname":
			cfg.dbname = match
		case "params":
			for _, v := range strings.Split(match, "&") {
				param := strings.SplitN(v, "=", 2)
				if len(param) != 2 {
					continue
				}
				cfg.params[param[0]] = param[1]
			}
		}
	}

	// Set default network if empty
	if cfg.net == "" {
		cfg.net = "tcp"
	}

	// Set default adress if empty
	if cfg.addr == "" {
		cfg.addr = "127.0.0.1:9306"
	}

	return cfg
}

// Encrypt password using 4.1+ method
// http://forge.mysql.com/wiki/MySQL_Internals_ClientServer_Protocol#4.1_and_later
func scramblePassword(scramble, password []byte) []byte {
	if len(password) == 0 {
		return nil
	}

	// stage1Hash = SHA1(password)
	crypt := sha1.New()
	crypt.Write(password)
	stage1Hash := crypt.Sum(nil)

	// scrambleHash = SHA1(scramble + SHA1(stage1Hash))
	// inner Hash
	crypt.Reset()
	crypt.Write(stage1Hash)
	scrambleHash := crypt.Sum(nil)

	// outer Hash
	crypt.Reset()
	crypt.Write(scramble)
	crypt.Write(scrambleHash)
	scrambleHash = crypt.Sum(nil)

	// token = scrambleHash XOR stage1Hash
	result := make([]byte, 20)
	for i := range result {
		result[i] = scrambleHash[i] ^ stage1Hash[i]
	}
	return result
}

/******************************************************************************
*                       Convert from and to bytes                             *
******************************************************************************/

func readLengthEnodedString(b []byte) ([]byte, bool, int, error) {
	// Get length
	num, isNull, n := readLengthEncodedInteger(b)
	if num < 1 {
		return nil, isNull, n, nil
	}

	n += int(num)

	// Check data length
	if len(b) >= n {
		return b[n-int(num) : n], false, n, nil
	}
	return nil, false, n, io.EOF
}

func skipLengthEnodedString(b []byte) (int, error) {
	// Get length
	num, _, n := readLengthEncodedInteger(b)
	if num < 1 {
		return n, nil
	}

	n += int(num)

	// Check data length
	if len(b) >= n {
		return n, nil
	}
	return n, io.EOF
}

func readLengthEncodedInteger(b []byte) (num uint64, isNull bool, n int) {
	switch (b)[0] {

	// 251: NULL
	case 0xfb:
		n = 1
		isNull = true
		return

	// 252: value of following 2
	case 0xfc:
		n = 3

	// 253: value of following 3
	case 0xfd:
		n = 4

	// 254: value of following 8
	case 0xfe:
		n = 9

	// 0-250: value of first byte
	default:
		num = uint64(b[0])
		n = 1
		return
	}

	switch n - 1 {
	case 2:
		num = uint64(b[0]) | uint64(b[1])<<8
		return
	case 3:
		num = uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16
		return
	default:
		num = uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 |
			uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 |
			uint64(b[6])<<48 | uint64(b[7])<<54
	}
	return
}

func lengthEncodedIntegerToBytes(n uint64) []byte {
	switch {
	case n <= 250:
		return []byte{byte(n)}

	case n <= 0xffff:
		return []byte{0xfc, byte(n), byte(n >> 8)}

	case n <= 0xffffff:
		return []byte{0xfd, byte(n), byte(n >> 8), byte(n >> 16)}
	}
	return nil
}
