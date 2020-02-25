package lib

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

//CreateValidationMap ranslates validation structure to map
// that can be easly presented in template
func CreateValidationMap(valid validation.Validation) map[string]map[string]string {
	v := make(map[string]map[string]string)
	/*
			{
				"email": {
					"Requrired" : "Can not be empty"
				},
				"password" :{

			  }
		  }
	*/
	for _, err := range valid.Errors {
		beego.Notice(err.Key, err.Message)
		k := strings.Split(err.Key, ".")
		var field, errorType string
		if len(k) > 1 {
			field = k[0]
			errorType = k[1]
		} else {
			field = err.Key
			errorType = " "
		}
		beego.Error(field)
		if _, ok := v[field]; !ok {
			v[field] = make(map[string]string)
		}
		v[field][errorType] = err.Message
	}
	return v

}

//Dump any structure as json string
func Dump(obj interface{}) {
	result, _ := json.MarshalIndent(obj, "", "\t")
	beego.Debug(string(result))
}

//CopyStruct serializes src and tries to deserialize it to dst
func CopyStruct(src interface{}, dst interface{}) error {
	jsonString, err := json.Marshal(&src)
	if err != nil {
		beego.Error("Unable to marshal object")
		return err
	}

	if err := json.Unmarshal([]byte(jsonString), &dst); err != nil {
		beego.Error("Unable to unmarshal object")
		return err
	}

	return nil
}

func ReadLine(filePth string) (string, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return "", nil
	}
	defer f.Close()
	var str []string
	bfRd := bufio.NewReader(f)
	flag := true
	for {
		line, err := bfRd.ReadBytes('\n')
		strLine := string(line)
		if flag {
			flag = !strings.Contains(strLine, "BEGIN")
		} else {
			if strLine != "" {
				str = append(str, strLine)
			}
			if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
				if err == io.EOF {
					return strings.Join(str, ""), nil
				}
				return "", err
			}
		}
	}
	return "", nil
}
