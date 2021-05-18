package worker

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	workerConfName = "./worker.yml"
)

type WorkerCfg struct {
	ExecDir   string `yaml:"exec_dir"`
	OutputDir string `yaml:"output_dir"`
}

func (wc *WorkerCfg) GetWorkerCfg() (*WorkerCfg, error) {
	yamlFile, err := ioutil.ReadFile(workerConfName)
	if err != nil {
		return &WorkerCfg{}, fmt.Errorf("yamlFile.Get err: %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, wc)
	if err != nil {
		return &WorkerCfg{}, fmt.Errorf("unmarshal: %v", err)
	}

	return wc, nil
}
