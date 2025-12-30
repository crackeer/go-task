package sshrun

import (
	"encoding/json"
	"fmt"
	"go-task/util"
)

type Input struct {
	Config  []string `json:"config"`
	Command []string `json:"command"`
}

// Run
//
//	@param data
//	@param sendFunc
//	@return string
//	@return error
func Run(data string, sendFunc func(string)) (string, error) {
	var input Input
	if err := json.Unmarshal([]byte(data), &input); err != nil {
		return "", fmt.Errorf("unmarshal input error: %s", err.Error())
	}

	retData := [][]string{}
	for _, item := range input.Config {
		sendFunc(fmt.Sprintf("new ssh client %s", item))
		sshClient, err := util.NewSSHClientByConfig(item)
		if err != nil {
			sendFunc(fmt.Sprintf("new ssh client error: %s", err.Error()))
			continue
		}

		for _, cmd := range input.Command {
			sendFunc(fmt.Sprintf("run command %s", cmd))
			output, err := sshClient.Run(cmd)
			if err != nil {
				sendFunc(fmt.Sprintf("run command error: %s", err.Error()))
				continue
			}
			retData = append(retData, []string{output})
		}
		sshClient.Close()
	}

	byteData, err := json.Marshal(retData)
	if err != nil {
		sendFunc(fmt.Sprintf("marshal error: %s", err.Error()))
		return "", fmt.Errorf("marshal error: %s", err.Error())
	}

	return string(byteData), nil
}
