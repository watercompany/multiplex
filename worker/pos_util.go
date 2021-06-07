package worker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"

	"gopkg.in/yaml.v2"
)

const (
	POSConfName = "./pos.yml"
)

type POSCfg struct {
	TempDir         string `yaml:"tmp_dir"`
	Temp2Dir        string `yaml:"tmp2_dir"`
	FinalDir        string `yaml:"final_dir"`
	FinalDestDir    string `yaml:"final_dest_dir"`
	FileName        string `yaml:"filename" json:"filename"`
	Size            string `yaml:"size"`
	PlotMemo        string `yaml:"plot_memo" json:"plot_memo"`
	PlotID          string `yaml:"plot_id" json:"plot_id"`
	Buffer          string `yaml:"buffer"`
	StripeSize      string `yaml:"stripe_size"`
	NumThreads      string `yaml:"num_threads"`
	NumBuckets      string `yaml:"num_buckets"`
	NoBitField      string `yaml:"nobitfield"`
	FarmerPublicKey string `yaml:"farmer_public_key"`
	PoolPublicKey   string `yaml:"pool_public_key"`

	SizeInt       int  `json:"size"`
	BufferInt     int  `json:"buffer"`
	StripeSizeInt int  `json:"stripe_size"`
	NumThreadsInt int  `json:"num_threads"`
	NoBitFieldInt bool `json:"nobitfield"`
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

	posCfgFromChia, err := GeneratePOSConfigFromChia(pc.FarmerPublicKey, pc.PoolPublicKey)
	if err != nil {
		return &POSCfg{}, fmt.Errorf("generate pos from chia: %v", err)
	}

	pc.FileName = posCfgFromChia.FileName
	pc.PlotID = posCfgFromChia.PlotID
	pc.PlotMemo = posCfgFromChia.PlotMemo

	return pc, nil
}

func GetPOSArgs(tempDir, temp2Dir, finalDir string) ([]string, error) {
	var args []string
	var posCfg POSCfg
	_, err := posCfg.GetPOSCfg()
	if err != nil {
		return []string{}, err
	}

	if tempDir != "" {
		posCfg.TempDir = tempDir
	}
	if temp2Dir != "" {
		posCfg.Temp2Dir = temp2Dir
	}
	if finalDir != "" {
		posCfg.FinalDir = finalDir
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

func GetTempAndFinalDir() (string, string, string, string, error) {
	var temp string
	var final string
	var finalDest string
	var plotName string
	var posCfg POSCfg

	_, err := posCfg.GetPOSCfg()
	if err != nil {
		return "", "", "", "", err
	}

	temp = posCfg.TempDir
	final = posCfg.FinalDir
	finalDest = posCfg.FinalDestDir
	plotName = posCfg.FileName

	return temp, final, finalDest, plotName, nil
}

func GeneratePOSConfigFromChia(farmerPubKey, poolPubkey string) (POSCfg, error) {
	POSCfgJson, err := ChiaBlockchainCmd(farmerPubKey, poolPubkey)
	if err != nil {
		return POSCfg{}, err
	}

	var posCfg POSCfg
	err = json.Unmarshal([]byte(POSCfgJson), &posCfg)
	if err != nil {
		return POSCfg{}, err
	}

	return posCfg, nil
}

func ChiaBlockchainCmd(farmerPubKey, poolPubkey string) (string, error) {
	var output string = ""
	cmd := exec.Command("./chia", "plots", "create", "-k", "32", "-f", farmerPubKey, "-p", poolPubkey, "--json")
	cmd.Dir = "/home/cx/repo/chia-blockchain/venv/bin/"

	// create a pipe for the output of the script
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			output = output + scanner.Text()
		}
	}()

	err = cmd.Start()
	if err != nil {
		return "", err
	}

	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	return output, nil
}
