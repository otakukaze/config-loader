package confloader

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
	"strconv"

	"git.trj.tw/golang/utils"
	"github.com/BurntSushi/toml"
	"github.com/otakukaze/envconfig"
	"gopkg.in/yaml.v2"
)

type ConfigFileType int

const (
	ConfigFileTypeJSON ConfigFileType = iota
	ConfigFileTypeYAML
	ConfigFileTypeTOML
)

type ConfigFile struct {
	Type ConfigFileType
	Path string
}

type LoadOptions struct {
	ConfigFile *ConfigFile
	FromEnv    bool
}

func Load(i interface{}, opts *LoadOptions) error {
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Ptr {
		return errors.New("input arg not ptr")
	}

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return errors.New("input not a struct")
	}

	// load default value
	LoadDefaultIntoStruct(i)

	// not config file opts, return
	if opts == nil {
		return nil
	}

	// load config file
	if opts.ConfigFile != nil {
		if opts.ConfigFile.Path == "" {
			return errors.New("config file path empty")
		}

		// resolve file path
		opts.ConfigFile.Path = utils.ParsePath(opts.ConfigFile.Path)
		// check file exists
		if !utils.CheckExists(opts.ConfigFile.Path, false) {
			return errors.New("config file not found")
		}

		filebyte, err := ioutil.ReadFile(opts.ConfigFile.Path)
		if err != nil {
			return err
		}

		switch opts.ConfigFile.Type {
		case ConfigFileTypeJSON:
			err := json.Unmarshal(filebyte, i)
			if err != nil {
				return err
			}
			break
		case ConfigFileTypeTOML:
			err := toml.Unmarshal(filebyte, i)
			if err != nil {
				return err
			}
			break
		case ConfigFileTypeYAML:
			err := yaml.Unmarshal(filebyte, i)
			if err != nil {
				return err
			}
			break
		default:
			return errors.New("file type not impl")
		}
	}

	// load config from env
	if opts.FromEnv {
		envconfig.Parse(i)
	}

	return nil
}

func LoadDefaultIntoStruct(i interface{}) {
	t := reflect.ValueOf(i)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// not struct skip
	if t.Kind() != reflect.Struct {
		return
	}

	fieldLen := t.NumField()
	for idx := 0; idx < fieldLen; idx++ {
		v := t.Field(idx)
		f := t.Type().Field(idx)

		val, tagExists := f.Tag.Lookup("default")

		if v.Type().Kind() == reflect.Slice {
			minLen := 0
			if defLen := f.Tag.Get("length"); defLen != "" {
				if convInt, err := strconv.ParseInt(defLen, 10, 64); err == nil {
					minLen = int(convInt)
				}
			}
			if minLen < 1 {
				return
			}

			val, tagExists := f.Tag.Lookup("default")

			slice := reflect.MakeSlice(f.Type, minLen, minLen)

			item := reflect.Indirect(slice.Index(0))

			if item.Type().Kind() == reflect.Slice {
				//slice in slice  skip proc
			} else if item.Type().Kind() == reflect.Struct {
				LoadDefaultIntoStruct(item.Addr().Interface())
			} else {
				if tagExists {
					procValue(item, val)
				}
			}

			for i := 0; i < slice.Len(); i++ {
				slice.Index(i).Set(item)
			}
			v.Set(slice)
		} else if v.Type().Kind() == reflect.Struct {
			LoadDefaultIntoStruct(v.Addr().Interface())
		} else {
			if tagExists {
				procValue(v, val)
			}
		}
	}
}

func procValue(v reflect.Value, val string) {
	if !v.IsValid() || !v.CanSet() {
		return
	}
	switch v.Type().Kind() {
	case reflect.String:
		v.SetString(val)
		break
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		if convInt, err := strconv.ParseInt(val, 10, 64); err == nil {
			if !v.OverflowInt(convInt) {
				v.SetInt(convInt)
			}
		}
		break
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		if convUint, err := strconv.ParseUint(val, 10, 64); err == nil {
			if !v.OverflowUint(convUint) {
				v.SetUint(convUint)
			}
		}
		break
	case reflect.Float32:
	case reflect.Float64:
		if convFloat, err := strconv.ParseFloat(val, 64); err == nil {
			if !v.OverflowFloat(convFloat) {
				v.SetFloat(convFloat)
			}
		}
		break
	case reflect.Bool:
		if convBool, err := strconv.ParseBool(val); err == nil {
			v.SetBool(convBool)
		}
		break
	}
}

func procSlice(field *reflect.StructField) {
	minLen := 0
	if defLen := field.Tag.Get("length"); defLen != "" {
		if convInt, err := strconv.ParseInt(defLen, 10, 64); err == nil {
			minLen = int(convInt)
		}
	}
	if minLen < 1 {
		return
	}

	val := field.Tag.Get("default")

	slice := reflect.MakeSlice(field.Type, minLen, minLen)

	item := reflect.Indirect(slice.Index(0))

	switch item.Kind() {
	case reflect.String:
		for i := 0; i < slice.Len(); i++ {
			slice.Index(i).Set(reflect.ValueOf(val))
		}
		break
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		if convInt, err := strconv.ParseInt(val, 10, 64); err == nil {
			if !slice.Index(0).OverflowInt(convInt) {
				for i := 0; i < slice.Len(); i++ {
					slice.Index(i).Set(reflect.ValueOf(convInt))
				}
			}
		}
		break
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		if convUint, err := strconv.ParseUint(val, 10, 64); err == nil {
			if !slice.Index(0).OverflowUint(convUint) {
				for i := 0; i < slice.Len(); i++ {
					slice.Index(i).Set(reflect.ValueOf(convUint))
				}
			}
		}
		break
	case reflect.Float32,
		reflect.Float64:
		if convFloat, err := strconv.ParseFloat(val, 64); err == nil {
			if !slice.Index(0).OverflowFloat(convFloat) {
				for i := 0; i < slice.Len(); i++ {
					slice.Index(i).Set(reflect.ValueOf(convFloat))
				}
			}
		}
		break
	case reflect.Bool:
		if conv, err := strconv.ParseBool(val); err == nil {
			for i := 0; i < slice.Len(); i++ {
				slice.Index(i).Set(reflect.ValueOf(conv))
			}
		}
		break
	case reflect.Struct:
		break
	}

	v := reflect.ValueOf(field)
	v.Set(slice)
}
