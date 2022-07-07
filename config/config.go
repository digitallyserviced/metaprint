package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/oxodao/metaprint/modules"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Battery     map[string]modules.Battery     `yaml:"battery"`
	Backlight   map[string]modules.Backlight   `yaml:"backlight"`
	Custom      map[string]modules.Custom      `yaml:"custom"`
	DateTime    map[string]modules.Date        `yaml:"datetime"`
	HackSpeed   map[string]modules.HackSpeed   `yaml:"hackspeed"`
	CpuInfo     map[string]modules.CpuInfo     `yaml:"cpu_info"`
	CpuUsage    map[string]modules.CpuUsage    `yaml:"cpu_usage"`
	NetUsage    map[string]modules.NetUsage    `yaml:"net_usage"`
	LoadAvg     map[string]modules.LoadAvg     `yaml:"load_avg"`
	Ip          map[string]modules.IP          `yaml:"ip"`
	Music       map[string]modules.Music       `yaml:"music"`
	PulseAudio  map[string]modules.PulseAudio  `yaml:"pulseaudio"`
	Ram         map[string]modules.Ram         `yaml:"ram"`
	Storage     map[string]modules.Storage     `yaml:"storage"`
	Temperature map[string]modules.Temperature `yaml:"temperature"`
	Uptime      map[string]modules.Uptime      `yaml:"uptime"`
}


func GetModulesAvailable() []string {
	modulesAvailable := []string{}

	t := reflect.TypeOf(Config{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		modulesAvailable = append(modulesAvailable, field.Tag.Get("yaml"))
	}

	return modulesAvailable
}

func getFieldNameFromModuleName(moduleType string) string {
	t := reflect.TypeOf(Config{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("yaml") == moduleType {
			return field.Name
		}
	}

	return ""
}

func (c Config) FindModule(moduleType, name string) (modules.Module, error) {
	fieldName := getFieldNameFromModuleName(moduleType)
	if len(fieldName) == 0 {
		return nil, nil
	}

	t := reflect.ValueOf(c)
	value := t.FieldByName(fieldName).MapIndex(reflect.ValueOf(name))

	if value.Kind() == reflect.Invalid {
		return nil, errors.New("could not find the " + fieldName + " module named " + name)
	}

	module := value.Interface().(modules.Module)

	return module, nil
}

func (c Config) GetModuleFormatsAvailable() (map[string]interface{}, error) {
    mods := GetModulesAvailable()
    out := make(map[string]interface{})
    for _, v := range mods {
        formats, ok := c.GetModuleFormats(v)

        if ok != nil {
            return nil, fmt.Errorf("%V", ok)
        }
        fmts := make(map[string]interface{})
        for formats.Next() {
            k := formats.Key()
            v := formats.Value()
            fields := reflect.VisibleFields(v.Type())

            fmtKeys := make(map[string]string)
            for _, field := range fields {
                fvkey := field.Name
                fvval := v.FieldByIndex(field.Index)
                if fvval.Kind() != reflect.Invalid && fvval.Kind().String() == "string" {
                   fmtKeys[fvkey] = fvval.String() 
                   continue
                }
            }
            fmts[k.String()]=fmtKeys
        }
        if len(fmts) == 0 {
            fmts["NA"] = "None Available"
        }
        // for _, f := range formats {
        //     fmtConfig, fok := c.FindModule(v, f.String())
        //     if fok != nil {
        //         break
        //     }
        //     for fmtConfig.Next
        //     rfmt := reflect.ValueOf(fmtConfig)
        //     fmts[f.String()] = rfmt.MapIndex()
        // }
        out[v] = fmts
    }
    return out, nil
}

func (c Config) GetModuleFormats(moduleType string) (*reflect.MapIter, error) {
	fieldName := getFieldNameFromModuleName(moduleType)
	if len(fieldName) == 0 {
		return nil, nil
	}

	t := reflect.ValueOf(c)
    moduleFormatKeys := t.FieldByName(fieldName).MapRange()

	// if len(moduleFormatKeys) == 0 {
	// 	return nil, errors.New("could not find any formats for module " + fieldName)
	// }

	return moduleFormatKeys, nil
}

// Load the configuration struct from a json file
func Load() *Config {
	var config Config

	globalPath := getPath() + "config.yml"

	hostname, err := os.Hostname()
	hasHostname := err == nil
	hostPath := getPath() + hostname + ".yml"

	anyFound := false
    // cfg, err := findAndLoadConfig()
	if _, err := os.Stat(globalPath); !os.IsNotExist(err) {
        if err != nil {
            fmt.Println(fmt.Sprintf("%s", err))
        }
        // actually check err, we need a working config otherwise complain
		err = loadConfig(&config, globalPath)
		if err == nil {
			anyFound = true
		}
	}

	if !hasHostname {
		return &config
	}

	// Loading it after will override the previous values if those exists
	if _, err := os.Stat(hostPath); !os.IsNotExist(err) {
		err = loadConfig(&config, hostPath)
		if err == nil {
			anyFound = true
		}
	}

	if !anyFound {
		fmt.Println(fmt.Sprintf("Could not load any config !\n %s", err))
		os.Exit(1)
	}

	return &config
}

func loadConfig(cfg *Config, path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, cfg)

	return err
}

func getPath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "./"
	}

	configPath := dirname + "/.config/metaprint"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return "./"
	}

	return configPath + "/"
}
