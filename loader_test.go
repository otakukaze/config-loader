package confloader

import (
	"os"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	type ObjKey struct {
		Name string `json:"name" yaml:"name" toml:"name"`
	}
	type ArrObj struct {
		Key string `json:"key" yaml:"key" toml:"key"`
	}
	type Config struct {
		StrKey   string   `json:"strKey" yaml:"strKey" toml:"strKey" default:"def value" env:"TEST_STR"`
		IntKey   int      `json:"intKey" yaml:"intKey" toml:"intKey" env:"TEST_INT"`
		BoolKey  bool     `json:"boolKey" yaml:"boolKey" toml:"boolKey"`
		FloatKey float64  `json:"floatKey" yaml:"floatKey" toml:"floatKey"`
		StrArr   []string `json:"strArr" yaml:"strArr" toml:"strArr" default:"arrval" length:"1"`
		StrArr2  []string `json:"strArr2" yaml:"strArr2" toml:"strArr2" default:"arrval2" length:"2"`
		ObjKey   ObjKey   `json:"objKey" yaml:"objKey" toml:"objKey"`
		ArrObj   []ArrObj `json:"arrObj" yaml:"arrObj" toml:"arrObj"`
	}

	os.Setenv("TEST_STR", "string value2")
	os.Setenv("TEST_INT", "2000")

	expected := Config{
		StrKey:   "string value",
		IntKey:   1000,
		BoolKey:  true,
		FloatKey: 1.2345,
		StrArr:   []string{"arr1"},
		StrArr2:  []string{"arrval2", "arrval2"},
		ObjKey: ObjKey{
			Name: "name",
		},
		ArrObj: []ArrObj{{Key: "val"}},
	}
	expectedWithEnv := Config{
		StrKey:   "string value2",
		IntKey:   2000,
		BoolKey:  true,
		FloatKey: 1.2345,
		StrArr:   []string{"arr1"},
		StrArr2:  []string{"arrval2", "arrval2"},
		ObjKey: ObjKey{
			Name: "name",
		},
		ArrObj: []ArrObj{{Key: "val"}},
	}

	// src := Config{}

	type args struct {
		i    interface{}
		opts *LoadOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test load json file with default",
			args: args{
				i: &Config{},
				opts: &LoadOptions{
					ConfigFile: &ConfigFile{
						Type: ConfigFileTypeJSON,
						Path: "./test/config.json",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test load yaml file with default",
			args: args{
				i: &Config{},
				opts: &LoadOptions{
					ConfigFile: &ConfigFile{
						Type: ConfigFileTypeYAML,
						Path: "./test/config.yaml",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test load toml file with default",
			args: args{
				i: &Config{},
				opts: &LoadOptions{
					ConfigFile: &ConfigFile{
						Type: ConfigFileTypeTOML,
						Path: "./test/config.toml",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test load json file with default and env",
			args: args{
				i: &Config{},
				opts: &LoadOptions{
					ConfigFile: &ConfigFile{
						Type: ConfigFileTypeJSON,
						Path: "./test/config.json",
					},
					FromEnv: true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(tt.args.i, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.args.opts != nil && tt.args.opts.FromEnv {
				if !reflect.DeepEqual(tt.args.i, expectedWithEnv) {
					t.Errorf("Load and expected not match")
				}
			} else {
				if !reflect.DeepEqual(tt.args.i, expected) {
					t.Errorf("Load and expected not match")
				}
			}
		})
	}
}

func Test_loadDefaultIntoStruct(t *testing.T) {
	type SubStruct struct {
		S string `default:"sss"`
	}
	type ArrStruct struct {
		S string `default:"inslice"`
	}

	type SS struct {
		S   string `default:"str"`
		I   int    `default:"123"`
		B   bool   `default:"true"`
		SS  SubStruct
		SL  []string    `default:"z" length:"1"`
		SL2 []ArrStruct `length:"1"`
		SL3 []int       `default:"100" length:"3"`
	}

	s := SS{}
	expected := SS{
		S:   "str",
		I:   123,
		B:   true,
		SS:  SubStruct{S: "sss"},
		SL:  []string{"z"},
		SL2: []ArrStruct{{S: "inslice"}},
		SL3: []int{100, 100, 100},
	}

	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test 1",
			args: args{
				i: &s,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadDefaultIntoStruct(tt.args.i)
			reflect.DeepEqual(s, expected)
		})
	}
}
