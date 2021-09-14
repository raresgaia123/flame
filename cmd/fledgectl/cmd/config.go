// Copyright (c) 2021 Cisco Systems, Inc. and its affiliates
// All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ApiServer ApiServer `yaml:"apiserver"`
	User      string    `yaml:"user"`
}

type ApiServer struct {
	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
}

var fs afero.Fs

func init() {
	fs = afero.NewOsFs()
}

func loadConfig(configPath string) (*Config, error) {
	data, err := afero.ReadFile(fs, configPath)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
