package raft

import "encoding/json"

func encodeProposeMsg(from uint64, data []byte) ([]byte, error) {
	msg := &ProposeMsg{
		Body: data,
		From: from,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return data, nil

}


func decodeProposeMsg(data []byte) (*ProposeMsg, error) {
	msg := &ProposeMsg{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func decodeResponse(data []byte) (*Response, error) {
	res := &Response{}
	err := json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func decodeSNRequest(data []byte) (*SNRequest, error) {
	snRequest := new(SNRequest)
	err := json.Unmarshal(data, snRequest)
	if err != nil {
		return nil, err
	}
	return snRequest, nil
}
