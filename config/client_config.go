/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"fmt"
)

import (
	"github.com/glory-go/glory/tools"
)

type NetworkConfig struct {
	Timeout int `yaml:"timeout" default:"3000"`
}

type ClientConfig struct {
	ConfigSource  string `yaml:"config_source"`
	ServerAddress string `yaml:"server_address"`
	ServiceID     string `yaml:"service_id"`
	Protocol      string `yaml:"protocol"`
	RegistryKey   string `yaml:"registry_key"`
	// NetworkConfig now store timeout for client
	FiltersKey    []string       `yaml:"filters_key"`
	NetworkConfig *NetworkConfig `yaml:"network"`
}

func (g *ClientConfig) checkAndFix() {
	if err := tools.ReadFromEnvIfNeed(g); err != nil {
		fmt.Println("warn: GloryClientConfig checkAndFix failed with err = ", err)
	}
}
