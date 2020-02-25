package mi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//ParsePid gets pid from string
func ParsePid(input string) (int64, error) {
	a := strings.Split(trim(input), "\n")
	if len(a) != 1 {
		return int64(0), fmt.Errorf("Wrong number of lines, expected %d, got %d", 1, len(a))
	}
	if !isSuccess(a[0]) {
		return int64(0), fmt.Errorf("Bad response: %s", a[0])
	}
	return strconv.ParseInt(stripPrefix(a[0], "SUCCESS: pid="), 10, 64)
}

//ParseVersion gets version information from string
func ParseVersion(input string) (*Version, error) {
	v := Version{}
	a := strings.Split(trim(input), "\n")
	if len(a) != 3 {
		return nil, fmt.Errorf("Wrong number of lines, expected %d, got %d", 3, len(a))
	}
	v.OpenVPN = stripPrefix(a[0], "OpenVPN Version: ")
	v.Management = stripPrefix(a[1], "Management Version: ")

	return &v, nil
}

//ParseStats gets stats from string
func ParseStats(input string) (*LoadStats, error) {
	ls := LoadStats{}
	a := strings.Split(trim(input), "\n")

	if len(a) != 1 {
		return nil, fmt.Errorf("预期的行数错误 %d, got %d", 1, len(a))
	}
	line := a[0]
	if !isSuccess(line) {
		return nil, fmt.Errorf("反应不佳: %s", line)
	}

	dString := stripPrefix(line, "SUCCESS: ")
	dElements := strings.Split(dString, ",")
	var err error
	ls.NClients, err = getLStatsValue(dElements[0])
	if err != nil {
		return nil, err
	}
	ls.BytesIn, err = getLStatsValue(dElements[1])
	if err != nil {
		return nil, err
	}
	ls.BytesOut, err = getLStatsValue(dElements[2])
	if err != nil {
		return nil, err
	}
	return &ls, nil
}

//ParseStatus gets status information from string
func ParseStatus(input string) (*Status, error) {
	region, err := New("conf/ip2region.db")
	s := Status{}
	defer region.Close()
	if err != nil {
		fmt.Println(err)
		return &s, err
	}
	s.ClientList = make([]*OVClient, 0, 0)
	s.RoutingTable = make([]*RoutingPath, 0, 0)
	a := strings.Split(trim(input), "\n")
	for _, line := range a {
		fields := strings.Split(trim(line), ",")
		c := fields[0]
		switch {
		case c == "TITLE":
			s.Title = fields[1]
		case c == "TIME":
			s.Time = fields[1]
			s.TimeT = fields[2]
		case c == "ROUTING_TABLE":
			item := &RoutingPath{
				VirtualAddress: fields[1],
				CommonName:     fields[2],
				RealAddress:    fields[3],
				LastRef:        fields[4],
				LastRefT:       fields[5],
			}
			if fields[2] != "UNDEF" {
				s.RoutingTable = append(s.RoutingTable, item)
			}
		case c == "CLIENT_LIST":
			bytesR, _ := strconv.ParseInt(fields[5], 10, 64)
			bytesS, _ := strconv.ParseInt(fields[6], 10, 64)
			ip := strings.Split(fields[2], ":")[0]
			ipAddress, _ := region.BtreeSearch(ip)
			dates := strings.Split(fields[7], " ")
			connectedSince := dates[4] + "-" + month(dates[1]) + "-" + dates[2] + " " + dates[3]
			item := &OVClient{
				CommonName:      fields[1],
				RealAddress:     ip,
				VirtualAddress:  fields[3],
				BytesReceived:   bytesR,
				BytesSent:       bytesS,
				ConnectedSince:  connectedSince,
				TimeOnline:      timeOnline(connectedSince),
				ConnectedSinceT: fields[8],
				ClientAddress:   ipAddress.Country + "-" + ipAddress.Province + "-" + ipAddress.City + "-" + ipAddress.ISP,
			}
			if fields[1] != "UNDEF" {
				s.ClientList = append(s.ClientList, item)
			}
		}
	}
	return &s, nil
}
func timeOnline(startTime string) string {
	var day, hour, min, sec int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	t2 := time.Now()

	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix()
		day = diff / (24 * 60 * 60)
		hour = diff/(60*60) - day*24
		min = (diff / 60) - day*24*60 - hour*60
		sec = diff - day*24*60*60 - hour*60*60 - min*60
		return strconv.FormatInt(day, 10) + "天" + strconv.FormatInt(hour, 10) + "小时" + strconv.FormatInt(min, 10) + "分" + strconv.FormatInt(sec, 10) + "秒"
	}
	return "0"
}
func month(month string) string {
	if strings.EqualFold(month, "Jan") {
		return "01"
	}
	if strings.EqualFold(month, "Feb") {
		return "02"
	}
	if strings.EqualFold(month, "Mar") {
		return "03"
	}
	if strings.EqualFold(month, "Apr") {
		return "04"
	}
	if strings.EqualFold(month, "May") {
		return "05"
	}
	if strings.EqualFold(month, "Jun") {
		return "06"
	}
	if strings.EqualFold(month, "Jul") {
		return "07"
	}
	if strings.EqualFold(month, "Aug") {
		return "08"
	}
	if strings.EqualFold(month, "Sep") {
		return "09"
	}
	if strings.EqualFold(month, "Oct") {
		return "10"
	}
	if strings.EqualFold(month, "Nov") {
		return "11"
	}
	if strings.EqualFold(month, "Dec") {
		return "12"
	}
	return ""
}

//ParseSignal checks for error in response string
func ParseSignal(input string) error {
	a := strings.Split(trim(input), "\n")
	if len(a) != 1 {
		return fmt.Errorf("预期的行数错误 %d, got %d", 1, len(a))
	}
	if !isSuccess(a[0]) {
		return fmt.Errorf("错误回应： %s", a[0])
	}
	return nil
}

//ParseKillSession gets kill command result from string
func ParseKillSession(input string) (string, error) {
	a := strings.Split(trim(input), "\n")

	if len(a) != 1 {
		return "", fmt.Errorf("预期的行数错误 %d, got %d", 1, len(a))
	}
	line := a[0]
	if !isSuccess(line) {
		return "", errors.New(line)
	}

	return stripPrefix(line, "SUCCESS: "), nil
}

func getLStatsValue(s string) (int64, error) {
	a := strings.Split(s, "=")
	if len(a) != 2 {
		return int64(-1), errors.New("解析错误")
	}
	return strconv.ParseInt(a[1], 10, 64)
}

func trim(s string) string {
	return strings.Trim(strings.Trim(s, "\r\n"), "\n")
}

func stripPrefix(s, prefix string) string {
	return trim(strings.Replace(s, prefix, "", 1))
}

func isSuccess(s string) bool {
	if strings.HasPrefix(s, "SUCCESS: ") {
		return true
	}
	return false
}
