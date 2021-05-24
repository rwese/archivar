package json

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.Level = logrus.DebugLevel
}

func TestJsonBrain_Load(t *testing.T) {
	type fields struct {
		backendArchiver archivers.Archiver
		backendGatherer archivers.Gatherer
		logger          *logrus.Logger
		Brain           Brain
		BrainFile       io.Reader
	}
	tests := []struct {
		name      string
		fields    fields
		wantBrain Brain
		wantErr   bool
	}{
		{
			name: "Simple loading",
			fields: fields{
				BrainFile: bytes.NewReader([]byte(`{"Memory":{"myMemory":0,"myOtherMemory":1}}`)),
			},
			wantBrain: Brain{
				Memory: map[string]int64{
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
			if !reflect.DeepEqual(i.Brain, tt.wantBrain) && !tt.wantErr {
				t.Fatalf("Memory is failing me\nWant: %+v\nHave: %+v", tt.wantBrain, i.Brain)
			}
		})
	}
}

func TestJsonBrain_Save(t *testing.T) {
	type fields struct {
		backendArchiver archivers.Archiver
		backendGatherer archivers.Gatherer
		logger          *logrus.Logger
		Brain           Brain
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
				Brain: Brain{
					Memory: map[string]int64{
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
				Brain:           tt.fields.Brain,
			}
			fh := &bytes.Buffer{}
			if err := i.Save(fh); (err != nil) != tt.wantErr {
				t.Errorf("JsonBrain.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFh := fh.String(); gotFh != tt.wantFh {
				t.Errorf("JsonBrain.Save() = \nWant: %+v\nHave: '%+v'", tt.wantFh, gotFh)
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
		Brain           Brain
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
				Brain:           tt.fields.Brain,
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

func TestJsonBrain_CleanMemory(t *testing.T) {
	type fields struct {
		file         io.ReadWriter
		maxMemoryAge int64
		Brain        Brain
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "New and old",
			fields: fields{
				maxMemoryAge: 0,
				file:         bytes.NewBuffer([]byte(`{"Memory":{"myMemory":0,"myOtherMemory":1}}`)),
			},
		},
		{
			name: "Keep it all",
			fields: fields{
				maxMemoryAge: -1,
				file:         bytes.NewBuffer([]byte(`{"Memory":{"myMemory":0,"myOtherMemory":1}}`)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &JsonBrain{
				maxMemoryAge: tt.fields.maxMemoryAge,
				Brain:        tt.fields.Brain,
			}
			i := NewBrain(c, nil, nil, logger)
			i.Connect()
			i.Load(tt.fields.file)

			if !i.DoYouKnow("myOtherMemory") {
				t.Fatal("I should know that")
			}

			i.CleanMemory()

			if i.DoYouKnow("myMemory") {
				t.Fatal("I should have forgotten that")
			}

			if !i.DoYouKnow("myOtherMemory") {
				t.Fatal("I should have remembered that")
			}
		})
	}
}
