package makeenvfile

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var configSeparator = "="
var reConfig = regexp.MustCompile(`(.*)=(.*)`)
var reComment = regexp.MustCompile(`^#|^\s*$`)

// Config to map from env
type Config struct {
	Name  string
	Value string
}

// Report of configs
type Report struct {
	DestExists     bool
	Equals         []Config
	NotEquals      []Config
	SourceNotFound []Config
	DestNotFound   []Config
}

func init() {
	log.SetPrefix("make-envfile ")
}

// GetFlags parse command line options
func GetFlags() (string, string, bool, int) {
	sourceFile := flag.String("source", ".env.example", "config file source")
	destFile := flag.String("dest", ".env", "config file destination")
	overwrite := flag.Bool("overwrite", false, "overwrite if config are not equal")
	sleep := flag.Int("sleep", 0, "seconds to sleep at the end if differences are found")

	flag.Parse()

	return *sourceFile, *destFile, *overwrite, *sleep
}

// ParseConfig and populate Config struct
func ParseConfig(filename string) ([]Config, error) {
	config := []Config{}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lineNb := 0
	for {
		line, _, err := reader.ReadLine()
		lineNb++
		configData := reConfig.FindString(string(line))
		commentData := reComment.MatchString(string(line))

		if err == io.EOF {
			break
		}

		if configData != "" {
			c := strings.Split(configData, configSeparator)
			newConfig := Config{
				Name:  c[0],
				Value: c[1],
			}
			config = append(config, newConfig)
		} else if !commentData {
			var errMsg strings.Builder
			errMsg.WriteString("Unknown property on line: ")
			errMsg.WriteString(strconv.Itoa(lineNb))
			err = errors.New(errMsg.String())
			return config, err
		}
	}
	return config, nil
}

// GetConfigDiff returns diff between two configs
func GetConfigDiff(sourceConfig []Config, destConfig []Config) ([]Config, []Config, []Config, []Config) {
	notEquals := []Config{}
	notFoundScs := []Config{}
	foundScs := []Config{}
	for _, sc := range sourceConfig {
		foundSc := false
		for j, dc := range destConfig {
			if sc.Name == dc.Name {
				foundSc = true
				foundScs = append(foundScs, sc)
				if sc.Value == dc.Value {
					// Found same
					destConfig = append(destConfig[:j], destConfig[j+1:]...)
					break
				} else {
					// Found but not same
					destConfig = append(destConfig[:j], destConfig[j+1:]...)
					notEquals = append(notEquals, sc)
					break
				}
			}
		}
		// Not found
		if !foundSc {
			notFoundScs = append(notFoundScs, sc)
		}
	}
	// Found or not
	return foundScs, notEquals, notFoundScs, destConfig
}

// CopyFile from source to dest
func CopyFile(source string, dest string) error {
	log.Printf("copy %s to %s", source, dest)

	// Read all content of src to data
	data, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}
	// Write data to dst
	err = ioutil.WriteFile(dest, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// CheckEnv to check config files
func CheckEnv(sourceFile string, destFile string) (Report, error) {
	report := Report{
		DestExists:     false,
		Equals:         []Config{},
		NotEquals:      []Config{},
		SourceNotFound: []Config{},
		DestNotFound:   []Config{},
	}

	if !fileExists(sourceFile) {
		errorMsg := fmt.Sprintf("source file: %v, not found", sourceFile)
		return report, errors.New(errorMsg)
	}

	sourceConfig, err := ParseConfig(sourceFile)
	if err != nil {
		return report, err
	}

	// Dest file not existing, just copy
	if !fileExists(destFile) {
		report.DestExists = false
		report.NotEquals = sourceConfig
		return report, nil
	} else {
		report.DestExists = true
	}

	destConfig, err := ParseConfig(destFile)
	if err != nil {
		return report, err
	}

	equals, notEquals, notFoundScs, notFoundDcs := GetConfigDiff(sourceConfig, destConfig)
	report.Equals = equals
	report.NotEquals = notEquals
	report.SourceNotFound = notFoundScs
	report.DestNotFound = notFoundDcs

	return report, nil
}

// Sleep for n seconds
func Sleep(sleep int) error {
	if sleep < 0 {
		return errors.New("sleep must be a positive value")
	} else if sleep > 0 {
		log.Printf("sleep for %d seconds", sleep)
	}

	time.Sleep(time.Duration(sleep) * time.Second)
	return nil
}

// MakeEnv to generate env file
func MakeEnv(sourceFile string, destFile string, overwrite bool, report Report, sleep int) error {
	if !report.DestExists {
		log.Printf("dest file does not exists")
		err := CopyFile(sourceFile, destFile)
		if err != nil {
			return err
		}
	} else if len(report.NotEquals) > 0 || len(report.SourceNotFound) > 0 || len(report.DestNotFound) > 0 {
		total := len(report.NotEquals) + len(report.SourceNotFound) + len(report.DestNotFound)
		log.Printf("warning, %d config files values are not equal:", total)
		if len(report.NotEquals) > 0 {
			for _, diff := range report.NotEquals {
				log.Printf("\t- %s are different", diff.Name)
			}
		}
		if len(report.SourceNotFound) > 0 {
			for _, diff := range report.SourceNotFound {
				log.Printf("\t- %s is not in dest config `%s`", diff.Name, sourceFile)
			}
		}
		if len(report.DestNotFound) > 0 {
			for _, diff := range report.DestNotFound {
				log.Printf("\t- %s is not in source config `%s`", diff.Name, destFile)
			}
		}
		if overwrite {
			log.Printf("overwrite enable:")
			err := CopyFile(sourceFile, destFile)
			if err != nil {
				return err
			}
		} else {
			log.Printf("overwrite disable, do nothing")
		}
		Sleep(sleep)
	} else {
		log.Printf("config files values are equal, do nothing")
	}

	return nil
}
