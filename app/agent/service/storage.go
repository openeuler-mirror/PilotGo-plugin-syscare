/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Fri Mar 22 16:27:24 2024 +0800
 */
package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/client"
	"gitee.com/openeuler/PilotGo-plugin-syscare/agent/config"
	"gitee.com/openeuler/PilotGo-plugin-syscare/utils"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

func DownloadPatchsFromServer(patch, taskId, downloadPath string) error {
	url := "http://" + client.GlobalClient.AgentServer + "/plugin/syscare/download/" + patch + "?path=" + taskId

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	out, err := os.Create(downloadPath + "/" + patch) // 保存文件的路径
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
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

func SearchDependentEnv() ([]*rpms, error) {
	buildEnv := []*rpms{}

	srcRpms, err := searchKernelSrcRpm()
	if err != nil {
		return nil, err
	}
	for version, srcRpm := range srcRpms {
		src := &rpms{
			Version: version,
			Rpms:    Rpms{SrcRpm: srcRpm},
		}
		buildEnv = append(buildEnv, src)
	}

	debugInfoRpms, err := searchKernelDebugRpm()
	if err != nil {
		return nil, err
	}
	for version, debugInfoRpm := range debugInfoRpms {
		found := false
		for _, existMap := range buildEnv {
			if existMap.Version == version {
				existMap.Rpms.DebugInfo = debugInfoRpm
				found = true
				break
			}
		}
		if !found {
			newRpm := &rpms{
				Version: version,
				Rpms:    Rpms{DebugInfo: debugInfoRpm},
			}
			buildEnv = append(buildEnv, newRpm)
		}
	}

	return buildEnv, nil
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
func ScanHotPatchRpm(path string) (string, string, error) {
	var kernelSrcRpm string
	var Rpm string
	files, err := os.ReadDir(path)
	if err != nil {
		return "", "", err
	}
	arch := client.GetClient().KernelArch
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "src.rpm") {
			matches := regexp.MustCompile(`(.+).src.rpm`).FindStringSubmatch(file.Name())
			if len(matches) >= 2 {
				kernelSrcRpm = file.Name()
			}
		}
		if strings.HasSuffix(file.Name(), arch+".rpm") {
			reString := fmt.Sprintf(`(.+)\.%s\.rpm`, arch)
			matches := regexp.MustCompile(reString).FindStringSubmatch(file.Name())
			if len(matches) >= 2 {
				Rpm = file.Name()
			}
		}
	}
	return kernelSrcRpm, Rpm, nil
}
func GetBuildLog(rootDir string) (string, error) {
	logFiles, err := walkLogFiles(rootDir)
	logger.Info("读取build.log: %v", logFiles[0])
	if err != nil {
		return "", err
	}
	logString, err := utils.FileReadString(logFiles[0])
	return logString, err
}

func walkLogFiles(rootDir string) ([]string, error) {
	var logFiles []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "build.log") {
			logFiles = append(logFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return logFiles, nil
}
