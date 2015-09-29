package common

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"

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

func Decode(filename string) (Gamesinfo, error) {
	// read the files
	data, err := ioutil.ReadFile(filename)

	var result Gamesinfo
	err = xml.Unmarshal(data, &result)
	return result, err
}
