package service

import (
	"testing"
	"encoding/json"
	"strings"
	"iax-test/models"
)

func Test_JSON_Parse(t *testing.T) {
	text := "{\"ads\":[{\"bidId\":\"4f1466313968208\",\"d\":\"1\",\"eid\":\"730080183\",\"height\":150,\"id\":\"1\",\"impId\":\"1\",\"link\":\"<html><head></head><body><html>  \n\t<body>\n\t\t<div>\n\t\t\t<a target='_blank' href='http://127.0.0.1:9090/v1/click?adxclick=http%3A%2F%2Ftracking.iax.optimix.asia%2Fiax-tracking%2Fadck%3Fs%3D1%26bid%3D4f1466313968208%26pub%3D2%26d%3D1%26lr%3Dhttp%253A%252F%252Fwww.baidu.com'>\n\t\t\t\t<img src='http://a03.optimix.asia/2015/12/18/bdbdc03fc9c4096bf5c2c49041da87a9853961.jpg' />\n\t\t\t</a>\n\t\t\t<img src='http://127.0.0.1:9090/v1/showup?bid=4f1466313968208&imp=1' height='0' width='0' />\n\t\t</div>\n\t</body>\n</html></body><script type=text/javascript src=http://127.0.0.1:9090/v1/winner?p=AAABVWcfiqcAAAAAAAAD6LH1-qBR_X7Zx9LhOA==&bid=4f1466313968208&rid=response-4f1466313968208&imp=1&seat=&adid=970734413&cur=CNY></script></html>\",\"noclick\":true,\"p\":5510000,\"position\":\"1001\",\"pubid\":2,\"rt\":5,\"sts\":\"e8ea3d75c8d779f6948d81e6cdaa5a21\",\"width\":300}],\"success\":true}"
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\t", "", -1)
	t.Log(text)
	var result models.Result
	err := json.Unmarshal([]byte(text), &result)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", result)

}
