package worker

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	POSConfName = "./pos.yml"
)

type POSCfg struct {
	TempDir    string `yaml:"tmp_dir"`
	Temp2Dir   string `yaml:"tmp2_dir"`
	FinalDir   string `yaml:"final_dir"`
	FileName   string `yaml:"filename"`
	Size       string `yaml:"size"`
	PlotMemo   string `yaml:"plot_memo"`
	PlotID     string `yaml:"plot_id"`
	Buffer     string `yaml:"buffer"`
	StripeSize string `yaml:"stripe_size"`
	NumThreads string `yaml:"num_threads"`
	NumBuckets string `yaml:"num_buckets"`
	NoBitField string `yaml:"nobitfield"`
}

func (pc *POSCfg) GetPOSCfg() (*POSCfg, error) {
	yamlFile, err := ioutil.ReadFile(POSConfName)
	if err != nil {
		return &POSCfg{}, fmt.Errorf("yamlFile.Get err: %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, pc)
	if err != nil {
		return &POSCfg{}, fmt.Errorf("unmarshal: %v", err)
	}

	return pc, nil
}

func GetPOSArgs() ([]string, error) {
	var args []string
	var posCfg POSCfg
	_, err := posCfg.GetPOSCfg()
	if err != nil {
		return []string{}, err
	}

	args = append(args, "create")
	args = append(args, "--tempdir", posCfg.TempDir)
	args = append(args, "--tempdir2", posCfg.Temp2Dir)
	args = append(args, "--finaldir", posCfg.FinalDir)
	args = append(args, "--file", posCfg.FileName)
	args = append(args, "--size", posCfg.Size)
	args = append(args, "--memo", posCfg.PlotMemo)
	args = append(args, "--id", posCfg.PlotID)
	args = append(args, "--buffer", posCfg.Buffer)
	args = append(args, "--stripes", posCfg.StripeSize)
	args = append(args, "--threads", posCfg.NumThreads)
	args = append(args, "--buckets", posCfg.NumBuckets)
	args = append(args, "--nobitfield", posCfg.NoBitField)

	return args, nil
}
