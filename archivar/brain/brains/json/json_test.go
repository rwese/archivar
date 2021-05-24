package json

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/sirupsen/logrus"
)

func TestJsonBrain_Load(t *testing.T) {
	type fields struct {
		backendArchiver archivers.Archiver
		backendGatherer archivers.Gatherer
		logger          *logrus.Logger
		Memory          Memory
		BrainFile       io.Reader
	}
	tests := []struct {
		name       string
		fields     fields
		wantMemory Memory
		wantErr    bool
	}{
		{
			name: "Simple loading",
			fields: fields{
				BrainFile: bytes.NewReader([]byte(`{"Memory":{"myMemory":0,"myOtherMemory":1}}`)),
			},
			wantMemory: Memory{
				map[string]int64{
					"myMemory":      0,
					"myOtherMemory": 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &JsonBrain{
				backendArchiver: tt.fields.backendArchiver,
				backendGatherer: tt.fields.backendGatherer,
				logger:          tt.fields.logger,
			}
			if err := i.Load(tt.fields.BrainFile); (err != nil) != tt.wantErr {
				t.Errorf("JsonBrain.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(i.Memory, tt.wantMemory) && !tt.wantErr {
				t.Fatalf("Memory is failing me\nWant: %+v\nHave: %+v", tt.wantMemory, i.Memory)
			}
		})
	}
}

func TestJsonBrain_Save(t *testing.T) {
	type fields struct {
		backendArchiver archivers.Archiver
		backendGatherer archivers.Gatherer
		logger          *logrus.Logger
		Memory          Memory
	}
	tests := []struct {
		name    string
		fields  fields
		wantFh  string
		wantErr bool
	}{
		{
			name: "Simple saving",
			fields: fields{
				Memory: Memory{
					map[string]int64{
						"myMemory":      0,
						"myOtherMemory": 1,
					},
				},
			},
			wantFh: `{"Memory":{"myMemory":0,"myOtherMemory":1}}
`,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &JsonBrain{
				backendArchiver: tt.fields.backendArchiver,
				backendGatherer: tt.fields.backendGatherer,
				logger:          tt.fields.logger,
				Memory:          tt.fields.Memory,
			}
			fh := &bytes.Buffer{}
			if err := i.Save(fh); (err != nil) != tt.wantErr {
				t.Errorf("JsonBrain.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFh := fh.String(); gotFh != tt.wantFh {
				t.Errorf("JsonBrain.Save() = %v, want '%v'", gotFh, tt.wantFh)
			}
		})
	}
}

func Test_serialize(t *testing.T) {
	type args struct {
		in []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name:    "simple serialize",
			args:    args{[]interface{}{"Some", "thing", "remembered", 1983}},
			wantOut: "Some;thing;remembered;1983",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := serialize(tt.args.in); gotOut != tt.wantOut {
				t.Errorf("serialize() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestJsonBrain_Remember(t *testing.T) {
	type fields struct {
		file            io.ReadWriter
		backendArchiver archivers.Archiver
		backendGatherer archivers.Gatherer
		logger          *logrus.Logger
		Memory          Memory
	}
	type args struct {
		data []interface{}
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		argsRemembered args
		wantRemembered bool
	}{
		{
			name: "",
			fields: fields{
				file: bytes.NewBuffer([]byte(``)),
			},
			args: args{
				[]interface{}{"Some", "thing", "remembered", 1983},
			},
			argsRemembered: args{
				[]interface{}{"Some", "thing", "remembered", 1983},
			},
			wantRemembered: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &JsonBrain{
				file:            tt.fields.file,
				backendArchiver: tt.fields.backendArchiver,
				backendGatherer: tt.fields.backendGatherer,
				logger:          tt.fields.logger,
				Memory:          tt.fields.Memory,
			}
			err := i.Connect()
			if err != nil && err != io.EOF {
				t.Error(err)
			}
			i.Remember(tt.args.data...)
			if i.DoYouKnow(tt.argsRemembered.data...) != tt.wantRemembered {
				t.Errorf("failed to remember")
			}
		})
	}
}
