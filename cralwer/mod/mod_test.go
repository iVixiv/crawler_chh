package mod

import (
	"testing"
	"encoding/json"
)

func Test_Mod(t *testing.T) {
	result, err := new(Mod).Crawler("forum.php?mod=forumdisplay\u0026fid=80\u0026mobile=2")
	if err != nil {
		log.Error("", err)
	} else {
		jsonStr, _ := json.Marshal(result)
		log.Notice("Crawler output: %s", string(jsonStr))
	}
}
