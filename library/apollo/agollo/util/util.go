package util

import (
	"bytes"
	"path"

	"github.com/pkg/errors"
	"github.com/why444216978/codec/json"
	"github.com/why444216978/codec/toml"
	"github.com/why444216978/codec/xml"
	"github.com/why444216978/codec/yaml"
)

var (
	jsonCodec = json.JSONCodec{}
	tomlCodec = toml.TomlCodec{}
	yamlCodec = yaml.YamlCodec{}
	xmlCodec  = xml.XMLCodec{}
)

func ExtractConf(namespace, content string, conf interface{}) error {
	switch path.Ext(namespace) {
	case "":
		return errors.New("ext is empty, maybe it is not support properties type")
	case ".json":
		return jsonCodec.Decode(bytes.NewReader([]byte(content)), conf)
	case ".txt":
		return tomlCodec.Decode(bytes.NewReader([]byte(content)), conf)
	case ".yaml", ".yml":
		return yamlCodec.Decode(bytes.NewReader([]byte(content)), conf)
	case ".xml":
		return xmlCodec.Decode(bytes.NewReader([]byte(content)), conf)
	}

	return errors.New("namespace ext error")
}
