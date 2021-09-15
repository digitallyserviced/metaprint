package modules

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/oxodao/metaprint/utils"
	"github.com/oliveagle/jsonpath"
)

type Temperature struct {
	Prefix string
	Suffix string

	Rounding int
	JsonPath string `yaml:"json_path"`
}

func (t Temperature) Print(args []string) string {
	data, err := exec.Command("sensors", "-j").Output()
	if err != nil {
		return "Could not get sensors data. Do you have lm_sensors installed & working?"
	}

	var sensorsData interface{}
	err = json.Unmarshal(data, &sensorsData)
	if err != nil {
		return "Could not get sensors data. Do you have lm_sensors installed & working?"
	}

	res, err := jsonpath.JsonPathLookup(sensorsData, t.JsonPath)
	if err != nil {
		return err.Error()
	}

	if val, ok := res.(float64); ok {
		return fmt.Sprintf("%v", utils.GetRoundedFloat(val, t.Rounding))
	}

	return fmt.Sprintf("%v", res)
}

func (t Temperature) GetPrefix() string {
	return t.Prefix
}

func (t Temperature) GetSuffix() string {
	return t.Suffix
}
