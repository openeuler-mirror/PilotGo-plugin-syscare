package service

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/client"
	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/config"
)

func PathInit(storage string) error {
	if _, err := os.Stat(storage); os.IsNotExist(err) {
		err := os.MkdirAll(storage, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

type rpms struct {
	Version string `json:"version"`
	Rpms    Rpms   `json:"rpm"`
}
type Rpms struct {
	SrcRpm    string `json:"srcRpm"`
	DebugInfo string `json:"debugInfo"`
}

func searchKernelSrcRpm() (map[string]string, error) {
	srcRpm := make(map[string]string)
	files, err := os.ReadDir(config.Config().Storage.Path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "kernel") && strings.HasSuffix(file.Name(), "src.rpm") {
			matches := regexp.MustCompile(`kernel-(.+)\.src\.rpm`).FindStringSubmatch(file.Name())
			if len(matches) >= 2 {
				version := matches[1]
				srcRpm[version] = file.Name()
			}
			continue
		}
	}
	return srcRpm, nil
}
func searchKernelDebugRpm() (map[string]string, error) {
	debuginfoRpm := make(map[string]string)
	files, err := os.ReadDir(config.Config().Storage.Path)
	if err != nil {
		return nil, err
	}

	arch := client.GetClient().KernelArch
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "kernel-debuginfo") && strings.HasSuffix(file.Name(), arch+".rpm") {
			reString := fmt.Sprintf(`kernel-debuginfo-(.+)\.%s\.rpm`, arch)
			matches := regexp.MustCompile(reString).FindStringSubmatch(file.Name())
			if len(matches) >= 2 {
				version := matches[1]
				debuginfoRpm[version] = file.Name()
			}
			continue
		}
	}
	return debuginfoRpm, nil
}
