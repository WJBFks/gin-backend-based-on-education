package service

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *ServiceSetup) ServiceTestAdd(td Test) (string, error) {

	eventID := "eventTestAdd"
	reg, _ := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)

	// 将td对象序列化成为字节数组
	b, err := json.Marshal(td)
	if err != nil {
		return "", fmt.Errorf("指定的td对象序列化时发生错误")
	}

	respone, err := t.Client.Execute(channel.Request{
		ChaincodeID: t.ChaincodeID,
		Fcn:         "TestAdd",
		Args:        [][]byte{b},
	})
	if err != nil {
		return "", err
	}

	// err = eventResult(notifier, eventID)
	// if err != nil {
	// 	return "", err
	// }

	return string(respone.TransactionID), nil
}

func (t *ServiceSetup) ServiceQueryTest(name string) (string, error) {
	respone, err := t.Client.Query(channel.Request{
		ChaincodeID: t.ChaincodeID,
		Fcn:         "TestQueryByName",
		Args:        [][]byte{[]byte(name)},
	})
	if err != nil {
		return "", err
	}

	return string(respone.Payload), nil
}
