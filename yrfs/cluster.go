package yrfs

import (
	"log"
	"strings"
	"yrfs-exporter/common"
)

//获取存储集群的副本数
func GetClustercopies() (metaCopies uint8, storageCopies uint8, err error) {
	err, stdout, _ := common.ShellExec("yrcli --getentry / --unmounted")
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}
	var metaType, storageType string
	metaCopies, storageCopies = 1, 1
	for _, s := range strings.Split(stdout, "\n") {
		if len(s) == 0 {
			continue
		}
		line := strings.ToLower(s)
		if strings.HasPrefix(line, "meta redundancy") {
			metaType = strings.Split(line, ":")[1]
			metaType = strings.TrimSpace(metaType)
		}
		if strings.HasPrefix(line, "data redundancy") {
			storageType = strings.Split(line, ":")[1]
			storageType = strings.TrimSpace(storageType)
		}
	}
	if metaType == "mirror" {
		metaCopies = 2
	}
	if storageType == "mirror" {
		storageCopies = 2
	}
	return
}