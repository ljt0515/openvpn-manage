package mi

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	IndexBlockLength  = 12
	TotalHeaderLength = 8192
)

var err error
var ipInfo IpInfo

type Ip2Region struct {
	// db file handler
	dbFileHandler *os.File

	//header block info

	headerSip []int64
	headerPtr []int64
	headerLen int64

	// super block index info
	firstIndexPtr int64
	lastIndexPtr  int64
	totalBlocks   int64

	// for memory mode only
	// the original db binary string

	dbBinStr []byte
	dbFile   string
}

type IpInfo struct {
	CityId   int64
	Country  string
	Region   string
	Province string
	City     string
	ISP      string
}

func (ip IpInfo) String() string {
	return strconv.FormatInt(ip.CityId, 10) + "|" + ip.Country + "|" + ip.Region + "|" + ip.Province + "|" + ip.City + "|" + ip.ISP
}

func getIpInfo(cityId int64, line []byte) IpInfo {

	lineSlice := strings.Split(string(line), "|")
	ipInfo := IpInfo{}
	length := len(lineSlice)
	ipInfo.CityId = cityId
	if length < 5 {
		for i := 0; i <= 5-length; i++ {
			lineSlice = append(lineSlice, "")
		}
	}

	ipInfo.Country = lineSlice[0]
	ipInfo.Region = lineSlice[1]
	ipInfo.Province = lineSlice[2]
	ipInfo.City = lineSlice[3]
	ipInfo.ISP = lineSlice[4]
	return ipInfo
}

func New(path string) (*Ip2Region, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &Ip2Region{
		dbFile:        path,
		dbFileHandler: file,
	}, nil
}

func (region *Ip2Region) Close() {
	region.dbFileHandler.Close()
}

func (region *Ip2Region) MemorySearch(ipStr string) (ipInfo IpInfo, err error) {
	ipInfo = IpInfo{}

	if region.totalBlocks == 0 {
		region.dbBinStr, err = ioutil.ReadFile(region.dbFile)

		if err != nil {

			return ipInfo, err
		}

		region.firstIndexPtr = getLong(region.dbBinStr, 0)
		region.lastIndexPtr = getLong(region.dbBinStr, 4)
		region.totalBlocks = (region.lastIndexPtr-region.firstIndexPtr)/IndexBlockLength + 1
	}

	ip, err := ip2long(ipStr)
	if err != nil {
		return ipInfo, err
	}

	h := region.totalBlocks
	var dataPtr, l int64
	for l <= h {

		m := (l + h) >> 1
		p := region.firstIndexPtr + m*IndexBlockLength
		sip := getLong(region.dbBinStr, p)
		if ip < sip {
			h = m - 1
		} else {
			eip := getLong(region.dbBinStr, p+4)
			if ip > eip {
				l = m + 1
			} else {
				dataPtr = getLong(region.dbBinStr, p+8)
				break
			}
		}
	}
	if dataPtr == 0 {
		return ipInfo, errors.New("not found")
	}

	dataLen := ((dataPtr >> 24) & 0xFF)
	dataPtr = (dataPtr & 0x00FFFFFF)
	ipInfo = getIpInfo(getLong(region.dbBinStr, dataPtr), region.dbBinStr[(dataPtr)+4:dataPtr+dataLen])
	return ipInfo, nil

}

func (region *Ip2Region) BinarySearch(ipStr string) (ipInfo IpInfo, err error) {
	ipInfo = IpInfo{}
	if region.totalBlocks == 0 {
		region.dbFileHandler.Seek(0, 0)
		superBlock := make([]byte, 8)
		region.dbFileHandler.Read(superBlock)
		region.firstIndexPtr = getLong(superBlock, 0)
		region.lastIndexPtr = getLong(superBlock, 4)
		region.totalBlocks = (region.lastIndexPtr-region.firstIndexPtr)/IndexBlockLength + 1
	}

	var l, dataPtr, p int64

	h := region.totalBlocks

	ip, err := ip2long(ipStr)

	if err != nil {
		return
	}

	for l <= h {
		m := (l + h) >> 1

		p = m * IndexBlockLength

		_, err = region.dbFileHandler.Seek(region.firstIndexPtr+p, 0)
		if err != nil {
			return
		}

		buffer := make([]byte, IndexBlockLength)
		_, err = region.dbFileHandler.Read(buffer)

		if err != nil {

		}
		sip := getLong(buffer, 0)
		if ip < sip {
			h = m - 1
		} else {
			eip := getLong(buffer, 4)
			if ip > eip {
				l = m + 1
			} else {
				dataPtr = getLong(buffer, 8)
				break
			}
		}

	}

	if dataPtr == 0 {
		err = errors.New("not found")
		return
	}

	dataLen := (dataPtr >> 24) & 0xFF
	dataPtr = dataPtr & 0x00FFFFFF

	region.dbFileHandler.Seek(dataPtr, 0)
	data := make([]byte, dataLen)
	region.dbFileHandler.Read(data)
	ipInfo = getIpInfo(getLong(data, 0), data[4:dataLen])
	err = nil
	return
}

func (region *Ip2Region) BtreeSearch(ipStr string) (ipInfo IpInfo, err error) {
	ipInfo = IpInfo{}
	ip, err := ip2long(ipStr)

	if region.headerLen == 0 {
		region.dbFileHandler.Seek(8, 0)

		buffer := make([]byte, TotalHeaderLength)
		region.dbFileHandler.Read(buffer)
		var idx int64
		for i := 0; i < TotalHeaderLength; i += 8 {
			startIp := getLong(buffer, int64(i))
			dataPar := getLong(buffer, int64(i+4))
			if dataPar == 0 {
				break
			}

			region.headerSip = append(region.headerSip, startIp)
			region.headerPtr = append(region.headerPtr, dataPar)
			idx++
		}

		region.headerLen = idx
	}

	var l, sptr, eptr int64
	h := region.headerLen

	for l <= h {
		m := int64(l+h) >> 1
		if m < region.headerLen {
			if ip == region.headerSip[m] {
				if m > 0 {
					sptr = region.headerPtr[m-1]
					eptr = region.headerPtr[m]
				} else {
					sptr = region.headerPtr[m]
					eptr = region.headerPtr[m+1]
				}
				break
			}
			if ip < region.headerSip[m] {
				if m == 0 {
					sptr = region.headerPtr[m]
					eptr = region.headerPtr[m+1]
					break
				} else if ip > region.headerSip[m-1] {
					sptr = region.headerPtr[m-1]
					eptr = region.headerPtr[m]
					break
				}
				h = m - 1
			} else {
				if m == region.headerLen-1 {
					sptr = region.headerPtr[m-1]
					eptr = region.headerPtr[m]
					break
				} else if ip <= region.headerSip[m+1] {
					sptr = region.headerPtr[m]
					eptr = region.headerPtr[m+1]
					break
				}
				l = m + 1
			}
		}

	}

	if sptr == 0 {
		err = errors.New("not found")
		return
	}

	blockLen := eptr - sptr
	region.dbFileHandler.Seek(sptr, 0)
	index := make([]byte, blockLen+IndexBlockLength)
	region.dbFileHandler.Read(index)
	var dataptr int64
	h = blockLen / IndexBlockLength
	l = 0

	for l <= h {
		m := int64(l+h) >> 1
		p := m * IndexBlockLength
		sip := getLong(index, p)
		if ip < sip {
			h = m - 1
		} else {
			eip := getLong(index, p+4)
			if ip > eip {
				l = m + 1
			} else {
				dataptr = getLong(index, p+8)
				break
			}
		}
	}

	if dataptr == 0 {
		err = errors.New("not found")
		return
	}

	dataLen := (dataptr >> 24) & 0xFF
	dataPtr := dataptr & 0x00FFFFFF

	region.dbFileHandler.Seek(dataPtr, 0)
	data := make([]byte, dataLen)
	region.dbFileHandler.Read(data)
	ipInfo = getIpInfo(getLong(data, 0), data[4:])
	return
}

func getLong(b []byte, offset int64) int64 {

	val := int64(b[offset]) |
		int64(b[offset+1])<<8 |
		int64(b[offset+2])<<16 |
		int64(b[offset+3])<<24

	return val

}

func ip2long(IpStr string) (int64, error) {
	bits := strings.Split(IpStr, ".")
	if len(bits) != 4 {
		return 0, errors.New("ip format error")
	}

	var sum int64
	for i, n := range bits {
		bit, _ := strconv.ParseInt(n, 10, 64)
		sum += bit << uint(24-8*i)
	}

	return sum, nil
}
