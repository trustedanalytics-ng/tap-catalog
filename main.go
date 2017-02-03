/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/gocraft/web"

	"github.com/trustedanalytics/tap-catalog/api"
	"github.com/trustedanalytics/tap-catalog/data"
	"github.com/trustedanalytics/tap-catalog/etcd"
	"github.com/trustedanalytics/tap-catalog/metrics"
	httpGoCommon "github.com/trustedanalytics/tap-go-common/http"
	commonLogger "github.com/trustedanalytics/tap-go-common/logger"
	"github.com/trustedanalytics/tap-go-common/util"
	mutils "github.com/trustedanalytics/tap-metrics/utils"
)

const EtcdComponentName = "ETCD_CATALOG"

var waitGroup = &sync.WaitGroup{}
var logger, _ = commonLogger.InitLogger("main")

func main() {
	rand.Seed(time.Now().UnixNano())

	go util.TerminationObserver(waitGroup, "Catalog")

	repository := setupRepository()
	context := setupContext(repository)
	r := setupRouter(context)

	startMetrics(repository)

	httpGoCommon.StartServer(r)
}

func setupContext(repository data.RepositoryApi) api.Context {
	context, err := api.NewContext(repository, getDefaultOrganization())
	if err != nil {
		logger.Fatalf("Cannot create new Context: %v", err)
	}
	return context
}

func setupRouter(context api.Context) *web.Router {
	r := api.SetupRouter(context)
	r.Get("/metrics", metricsHandler())
	return r
}

func setupRepository() data.RepositoryApi {
	etcdAddress, etcdPort, err := util.GetConnectionHostAndPortFromEnvs(EtcdComponentName)
	if err != nil {
		logger.Fatalf("Cannot get ETCD address and port: %v", err)
	}
	etcdKVStore, err := etcd.NewEtcdKVStore(etcdAddress, etcdPort)
	if err != nil {
		logger.Fatalf("Cannot connect to ETCD on %s:%d: %v", etcdAddress, etcdPort, err)
	}
	return data.NewRepositoryAPI(etcdKVStore, data.DataMapper{})
}

func getDefaultOrganization() string {
	return os.Getenv("CORE_ORGANIZATION")
}

func startMetrics(repository data.RepositoryApi) {
	mcfenv := os.Getenv("METRICS_COLLECTING_FREQUENCY")
	mcf, err := time.ParseDuration(mcfenv)
	if err != nil {
		logger.Warningf("Couldn't parse metrics frequency setting (got: %s), fallback to default.", mcfenv)
		mcf = 15 * time.Second
	}
	metrics.EnableCollection(repository, mcf)
}

func metricsHandler() func(rw web.ResponseWriter, req *web.Request) {
	mHandler := mutils.GetHandler()
	return func(rw web.ResponseWriter, req *web.Request) {
		mHandler.ServeHTTP(rw, req.Request)
	}
}
