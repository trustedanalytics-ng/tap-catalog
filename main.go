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

var waitGroup = &sync.WaitGroup{}
var logger, _ = commonLogger.InitLogger("main")

func main() {
	rand.Seed(time.Now().UnixNano())

	go util.TerminationObserver(waitGroup, "Catalog")

	repository := setupRepository()
	context := setupContext(repository)

	startMetrics(repository)

	r := setupRouter(context)

	httpGoCommon.StartServer(r)
}

func setupRouter(context api.Context) *web.Router {
	r := web.New(context)
	r.Middleware(web.LoggerMiddleware)
	r.Get("/healthz", context.GetCatalogHealth)
	r.Get("/metrics", metricsHandler())

	apiRouter := r.Subrouter(context, "/api")

	basicAuthRouter := apiRouter.Subrouter(context, "/v1")
	route(basicAuthRouter, &context)

	v1AliasRouter := apiRouter.Subrouter(context, "/v1.0")
	route(v1AliasRouter, &context)

	r.Get("/", context.Index)
	r.Error(context.Error)
	return r
}

func setupContext(repository data.RepositoryApi) api.Context {
	context, err := api.NewContext(repository, getDefaultOrganization())
	if err != nil {
		logger.Fatalf("Cannot create new Context: %v", err)
	}
	return context
}

func setupRepository() data.RepositoryApi {
	etcdAddress, etcdPort, err := util.GetConnectionHostAndPortFromEnvs("ETCD")
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

func route(router *web.Router, context *api.Context) {
	router.Middleware(context.BasicAuthorizeMiddleware)
	router.Middleware(context.OrganizationSetupMiddleware)

	router.Get("/services", context.Services)
	router.Get("/services/:serviceId", context.GetService)
	router.Post("/services", context.AddService)
	router.Patch("/services/:serviceId", context.PatchService)
	router.Delete("/services/:serviceId", context.DeleteService)

	router.Get("/services/:serviceId/plans", context.Plans)
	router.Get("/services/:serviceId/plans/:planId", context.GetPlan)
	router.Post("/services/:serviceId/plans", context.AddPlan)
	router.Patch("/services/:serviceId/plans/:planId", context.PatchPlan)
	router.Delete("/services/:serviceId/plans/:planId", context.DeletePlan)

	router.Get("/services/instances", context.ServicesInstances)
	router.Get("/services/:serviceId/instances", context.ServiceInstances)
	router.Get("/services/:serviceId/instances/:instanceId", context.GetServiceInstance)
	router.Post("/services/:serviceId/instances", context.AddServiceInstance)
	router.Patch("/services/:serviceId/instances/:instanceId", context.PatchServiceInstance)
	router.Delete("/services/:serviceId/instances/:instanceId", context.DeleteServiceInstance)

	router.Get("/applications", context.Applications)
	router.Get("/applications/:applicationId", context.GetApplication)
	router.Post("/applications", context.AddApplication)
	router.Patch("/applications/:applicationId", context.PatchApplication)
	router.Delete("/applications/:applicationId", context.DeleteApplication)

	router.Get("/applications/instances", context.ApplicationsInstances)
	router.Get("/applications/:applicationId/instances", context.ApplicationInstances)
	router.Get("/applications/:applicationId/instances/:instanceId", context.GetApplicationInstance)
	router.Post("/applications/:applicationId/instances", context.AddApplicationInstance)
	router.Patch("/applications/:applicationId/instances/:instanceId", context.PatchApplicationInstance)
	router.Delete("/applications/:applicationId/instances/:instanceId", context.DeleteApplicationInstance)

	router.Get("/images", context.Images)
	router.Get("/images/nextState", context.MonitorImagesStates)
	router.Get("/images/:imageId", context.GetImage)
	router.Get("/images/:imageId/nextState", context.MonitorSpecificImageState)
	router.Post("/images", context.AddImage)
	router.Patch("/images/:imageId", context.PatchImage)
	router.Delete("/images/:imageId", context.DeleteImage)

	router.Get("/instances", context.Instances)
	router.Get("/instances/nextState", context.MonitorInstancesStates)
	router.Get("/instances/:instanceId", context.GetInstance)
	router.Get("/instances/:instanceId/nextState", context.MonitorSpecificInstanceState)
	router.Get("/instances/:instanceId/bindings", context.GetInstanceBindings)
	router.Delete("/instances/:instanceId", context.DeleteInstance)
	router.Patch("/instances/:instanceId", context.PatchInstance)

	router.Get("/templates", context.Templates)
	router.Post("/templates", context.AddTemplate)
	router.Get("/templates/:templateId", context.GetTemplate)
	router.Delete("/templates/:templateId", context.DeleteTemplate)
	router.Patch("/templates/:templateId", context.PatchTemplate)
}
