package json

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/brain/brains"
	"github.com/rwese/archivar/internal/utils/config"
	"github.com/sirupsen/logrus"
)

func init() {
	brains.RegisterBrain(NewBrain)
}

type JsonBrainConfig struct {
	File string
}

type JsonBrain struct {
	file            io.ReadWriter
	backendArchiver archivers.Archiver
	backendGatherer archivers.Gatherer
	logger          *logrus.Logger
	Memory
}

type Memory struct {
	Memory map[string]int64
}

func NewBrain(c interface{}, archiver archivers.Archiver, gatherer archivers.Gatherer, logger *logrus.Logger) (i brains.Brain) {
	igc := JsonBrainConfig{}

	config.ConfigFromStruct(c, &igc)

	fh, err := os.OpenFile(igc.File, os.O_EXCL, 0660)
	if err != nil {
		logger.Fatalf("Failed to open brain file: %v", err)
	}

	return &JsonBrain{
		file:            fh,
		logger:          logger,
		backendArchiver: archiver,
		backendGatherer: gatherer,
	}
}

func (i *JsonBrain) DoYouKnow(data ...interface{}) bool {
	memoryLine := serialize(data)
	if _, exists := i.Memory.Memory[memoryLine]; !exists {
		return false
	}
	return true
}

func (i *JsonBrain) Init(data ...interface{}) {
	i.Memory.Memory = make(map[string]int64)
}

func (i *JsonBrain) Remember(data ...interface{}) {
	memoryLine := serialize(data)
	i.Memory.Memory[memoryLine] = time.Now().Unix()
}

func (i *JsonBrain) Connect() (err error) {
	i.Init()
	return i.Load(i.file)
}

func (i *JsonBrain) Disconnect() (err error) {
	return i.Save(i.file)
}

func (i *JsonBrain) SaveFile(brainfile string) (err error) {
	brainfh, err := os.Open(brainfile)
	if err != nil {
		return err
	}
	defer func() {
		if err := brainfh.Close(); err != nil {
			i.logger.Warnf("failed to close brainfile: %v", err)
		}
	}()

	return i.Save(brainfh)
}

func (i *JsonBrain) Save(fh io.Writer) (err error) {
	brainDump := json.NewEncoder(fh)
	return brainDump.Encode(i.Memory)
}

func (i *JsonBrain) LoadFile(brainfile string) (err error) {
	brainfh, err := os.Open(brainfile)
	if err != nil {
		return err
	}
	defer func() {
		if err := brainfh.Close(); err != nil {
			i.logger.Warnf("failed to close brainfile: %v", err)
		}
	}()
	if err != nil {
		fmt.Println(err.Error())
	}
	return i.Load(brainfh)
}

func (i *JsonBrain) Load(fh io.Reader) (err error) {
	jsonParser := json.NewDecoder(fh)
	return jsonParser.Decode(&i.Memory)
}

const serialize_separator = ";"

// serialize is a poor implementation of a serializer
func serialize(in []interface{}) (out string) {
	var inStrings []string
	for _, i := range in {
		inStrings = append(inStrings, fmt.Sprintf("%v", i))
	}

	return strings.Join(inStrings, serialize_separator)
}
