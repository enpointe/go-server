package server

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	tests := []struct {
		filename string
		want     *Config
		wantErr  bool
	}{
		{"testdata/testConfig1.json", &Config{"TestKey"}, false},
		{"testdata/nonExistent.json", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got, err := ReadConfig(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleReadConfig() {
	config, err := ReadConfig("Configuration.json")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(config)

	//Output:
	//&{secretKey}
}
