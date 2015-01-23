package common

import (
	"bytes"
	"encoding/xml"

	. "model"
)

func Encode(gameinfo *Gamesinfo) ([]byte, error) {
	result := bytes.NewBuffer(nil)
	data, err := xml.MarshalIndent(gameinfo, "  ", "  ")

	if err != nil {
		return data, err
	}

	result.WriteString(xml.Header)
	result.Write(data)

	return result.Bytes(), nil

}
