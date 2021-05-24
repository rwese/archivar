package json

import (
	"bytes"
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
	File   string // File is the location within the archiver where to store the brain
	MaxAge string // MaxAge of entries remembered Format, default: "7d"
}

type JsonBrain struct {
	file            io.ReadWriter
	maxMemoryAge    int64
	backendArchiver archivers.Archiver
	backendGatherer archivers.Gatherer
	logger          *logrus.Logger
	Brain           Brain
}

type Brain struct {
	Memory map[string]int64
}

func NewBrain(c interface{}, archiver archivers.Archiver, gatherer archivers.Gatherer, logger *logrus.Logger) brains.BrainInterface {
	igc := JsonBrainConfig{}

	config.ConfigFromStruct(c, &igc)

	jb := &JsonBrain{
		logger:          logger,
		backendArchiver: archiver,
		backendGatherer: gatherer,
	}
	if igc.File != "" {
		fh, err := os.OpenFile(igc.File, os.O_EXCL, 0660)
		if err != nil {
			logger.Fatalf("Failed to open brain file: %v", err)
		}

		jb.file = fh
	} else {
		logger.Warn("no file for the brain defined, no state will be saved")
		jb.file = bytes.NewBuffer([]byte(`{}`))
	}

	if igc.MaxAge == "" {
		igc.MaxAge = "168h"
	}

	maxMemoryAge, err := time.ParseDuration(igc.MaxAge)
	if err != nil {
		logger.Fatalf("failed to parse given MaxAge:'%s'", igc.MaxAge)
	}

	jb.maxMemoryAge = int64(maxMemoryAge.Seconds())
	return jb
}

func (i *JsonBrain) DoYouKnow(data ...interface{}) bool {
	memoryLine := serialize(data)
	if _, exists := i.Brain.Memory[memoryLine]; !exists {
		return false
	}

	i.Brain.Memory[memoryLine] = time.Now().Unix()
	return true
}

func (i *JsonBrain) Init(data ...interface{}) {
	i.Brain.Memory = make(map[string]int64)
}

func (i *JsonBrain) Remember(data ...interface{}) {
	memoryLine := serialize(data)
	i.Brain.Memory[memoryLine] = time.Now().Unix()
}

func (i *JsonBrain) Connect() (err error) {
	i.Init()
	if i.file == nil {
		return nil
	}

	return i.Load(i.file)
}

func (i *JsonBrain) Disconnect() (err error) {
	if i.file == nil {
		return nil
	}

	return i.Save(i.file)
}

func (i *JsonBrain) SaveToFile(brainfile string) (err error) {
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
	return brainDump.Encode(i.Brain)
}

func (i *JsonBrain) LoadFromFile(brainfile string) (err error) {
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
	return jsonParser.Decode(&i.Brain)
}

// CleanMemory will meditate on the old (>maxAge) and remove those memories
func (i *JsonBrain) CleanMemory() {
	oldestMemoryAge := time.Now().Unix() - i.maxMemoryAge
	for memory, lastSeen := range i.Brain.Memory {
		if lastSeen > oldestMemoryAge {
			continue
		}

		i.logger.Debugf("Brain has forgotten about: %s", memory)
		delete(i.Brain.Memory, memory)
	}
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
