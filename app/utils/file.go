/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Tue Apr 2 17:35:07 2024 +0800
 */
package utils

import (
	"io"
	"os"
)

func MakeDir(storage string) error {
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

// 读取文件所有数据，返回字符串
func FileReadString(filePath string) (string, error) {
	f, err := os.Open(filePath)
	defer func(file *os.File) {
		// ignore file close error
		file.Close()
	}(f)
	if err != nil {
		return "", err
	}

	var result []byte
	readBuff := make([]byte, 1024*4)
	for {
		n, err := f.Read(readBuff)
		if err != nil {
			if err == io.EOF {
				if n != 0 {
					result = append(result, readBuff[:n]...)
				}
				break
			}
			return "", err
		}
		result = append(result, readBuff[:n]...)
	}
	return string(result), nil
}
