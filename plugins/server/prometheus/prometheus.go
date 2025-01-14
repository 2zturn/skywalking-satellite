// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/pkg/log"
	"github.com/apache/skywalking-satellite/internal/satellite/telemetry"
)

const (
	Name     = "prometheus-server"
	ShowName = "Prometheus Server"
)

type Server struct {
	config.CommonFields
	Address  string `mapstructure:"address"`  // The prometheus server address.
	Endpoint string `mapstructure:"endpoint"` // The prometheus server metrics endpoint.

	server *http.ServeMux // The prometheus server.
}

func (s *Server) Name() string {
	return Name
}

func (s *Server) ShowName() string {
	return ShowName
}

func (s *Server) Description() string {
	return "This is a prometheus server to export the metrics in Satellite."
}

func (s *Server) DefaultConfig() string {
	return `
# The prometheus server address.
address: ":1234"
# The prometheus server metrics endpoint.
endpoint: "/metrics"
`
}

func (s *Server) Prepare() error {
	s.server = http.NewServeMux()
	return nil
}

func (s *Server) Start() error {
	// add go info metrics.
	telemetry.Register(telemetry.WithMeta("processor_collector", prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{})),
		telemetry.WithMeta("go_collector", prometheus.NewGoCollector()))
	// register prometheus metrics exporter handler.
	s.server.Handle(s.Endpoint, promhttp.HandlerFor(telemetry.Gatherer, promhttp.HandlerOpts{ErrorLog: log.Logger}))
	go func() {
		log.Logger.WithField("address", s.Address).Info("prometheus server is starting...")
		err := http.ListenAndServe(s.Address, s.server)
		if err != nil {
			log.Logger.WithField("address", s.Address).Infof("prometheus server has failure when starting: %v", err)
		}
	}()
	return nil
}

func (s *Server) Close() error {
	log.Logger.Info("prometheus server is closed")
	return nil
}

func (s *Server) GetServer() interface{} {
	return s.server
}
