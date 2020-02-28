package lib

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/ljt000/openvpn-manage/models"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

//Cert
//https://groups.google.com/d/msg/mailing.openssl.users/gMRbePiuwV0/wTASgPhuPzkJ
type Cert struct {
	EntryType   string
	Expiration  string
	ExpirationT time.Time
	Revocation  string
	RevocationT time.Time
	Serial      string
	FileName    string
	Details     *Details
}

type Details struct {
	Name         string
	CN           string
	Country      string
	Organisation string
	Email        string
}
type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ReadCerts(path string) ([]*Cert, error) {
	certs := make([]*Cert, 0, 0)
	text, err := ioutil.ReadFile(path)
	if err != nil {
		return certs, err
	}
	lines := strings.Split(trim(string(text)), "\n")
	for _, line := range lines {
		fields := strings.Split(trim(line), "\t")
		if len(fields) != 6 {
			return certs,
				fmt.Errorf("Incorrect number of lines in line: \n%s\n. Expected %d, found %d",
					line, 6, len(fields))
		}
		if fields[0] == "R" {
			continue
		}
		expT, _ := time.Parse("060102150405Z", fields[1])
		//revT, _ := time.Parse("060102150405Z", fields[2])
		c := &Cert{
			EntryType:   fields[0],
			Expiration:  fields[1],
			ExpirationT: expT,
			Serial:      fields[3],
			FileName:    fields[4],
			Details:     parseDetails(fields[5]),
		}
		certs = append(certs, c)
	}

	return certs, nil
}

func parseDetails(d string) *Details {
	details := &Details{}
	lines := strings.Split(trim(string(d)), "/")
	for _, line := range lines {
		if strings.Contains(line, "") {
			fields := strings.Split(trim(line), "=")
			switch fields[0] {
			case "name":
				details.Name = fields[1]
			case "CN":
				details.Name = fields[1]
				details.CN = fields[1]
			case "C":
				details.Country = fields[1]
			case "O":
				details.Organisation = fields[1]
			case "emailAddress":
				details.Email = fields[1]
			default:
				beego.Warn(fmt.Sprintf("Undefined entry: %s", line))
			}
		}
	}
	return details
}

func trim(s string) string {
	return strings.Trim(strings.Trim(s, "\r\n"), "\n")
}

func CreateCertificate(name string) bool {
	_, err := os.Stat(models.GlobalCfg.OVConfigPath + "easy-rsa/pki/issued/" + name + ".crt")
	if !os.IsNotExist(err) {
		return true
	}
	return runCmd("EASYRSA_CERT_EXPIRE=3650 ./easyrsa build-client-full " + name + " nopass")
}
func runCmd(command string) bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("/bin/bash", "-c", command)
	}
	cmd.Dir = models.GlobalCfg.OVConfigPath + "easy-rsa/"
	output, err := cmd.CombinedOutput()
	if err != nil {
		beego.Debug(string(output))
		beego.Error(err)
		return false
	}
	Dump(ConvertByte2String(output, GB18030))
	return false
}
func DelCertificate(name string) bool {
	runCmd("./easyrsa --batch revoke " + name)
	runCmd("EASYRSA_CRL_DAYS=3650 ./easyrsa gen-crl")

	path := models.GlobalCfg.OVConfigPath + "easy-rsa/pki/"
	os.Remove(path + "issued/" + name + ".crt")
	os.Remove(path + "private/" + name + ".key")
	os.Remove(path + "reqs/" + name + ".req")
	CopyFile(path+"crl.pem", models.GlobalCfg.OVConfigPath+"crl.pem")
	runCmd("chown nobody:nobody " + models.GlobalCfg.OVConfigPath + "crl.pem")
	return false
}

func CopyFile(srcFileName string, dstFileName string) (int64, error) {
	srcFile, err := os.Open(srcFileName)

	if err != nil {
		fmt.Printf("open file err = %v\n", err)
		return 0, nil
	}

	defer srcFile.Close()

	//打开dstFileName

	dstFile, err := os.OpenFile(dstFileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Printf("open file err = %v\n", err)
		return 0, nil
	}

	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}
func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
