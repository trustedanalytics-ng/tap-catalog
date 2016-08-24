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
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gocraft/web"

	"github.com/trustedanalytics/tapng-catalog/api"
	"github.com/trustedanalytics/tapng-catalog/data"
	httpGoCommon "github.com/trustedanalytics/tapng-go-common/http"
)

type appHandler func(web.ResponseWriter, *web.Request) error

func main() {
	rand.Seed(time.Now().UnixNano())

	err := (&data.RepositoryConnector{}).CreateDirs()
	if err != nil {
		log.Fatalln("Can't create directories oin ETCD!", err)
	}

	context := api.Context{}

	r := web.New(context)
	r.Middleware(web.LoggerMiddleware)
	r.Get("/healthz", context.GetCatalogHealth)

	apiRouter := r.Subrouter(context, "/api")

	basicAuthRouter := apiRouter.Subrouter(context, "/v1")
	route(basicAuthRouter)

	// for testing purpose, where v1 is current version
	v1AliasRouter := apiRouter.Subrouter(context, "/v1.0")
	route(v1AliasRouter)

	r.Get("/", (*api.Context).Index)
	r.Error((*api.Context).Error)

	if os.Getenv("CATALOG_SSL_CERT_FILE_LOCATION") != "" {
		httpGoCommon.StartServerTLS(os.Getenv("CATALOG_SSL_CERT_FILE_LOCATION"),
			os.Getenv("CATALOG_SSL_KEY_FILE_LOCATION"), r)
	} else {
		httpGoCommon.StartServer(r)
	}

}

func route(router *web.Router) {
	router.Middleware((*api.Context).BasicAuthorizeMiddleware)

	router.Get("/services", (*api.Context).Services)
	router.Get("/services/:serviceId", (*api.Context).GetService)
	router.Post("/services", (*api.Context).AddService)
	router.Patch("/services/:serviceId", (*api.Context).PatchService)
	router.Delete("/services/:serviceId", (*api.Context).DeleteService)

	router.Get("/services/:serviceId/plans", (*api.Context).Plans)
	router.Get("/services/:serviceId/plans/:planId", (*api.Context).GetPlan)
	router.Post("/services/:serviceId/plans", (*api.Context).AddPlan)
	router.Patch("/services/:serviceId/plans/:planId", (*api.Context).PatchPlan)
	router.Delete("/services/:serviceId/plans/:planId", (*api.Context).DeletePlan)

	router.Get("/services/instances", (*api.Context).ServicesInstances)
	router.Get("/services/:serviceId/instances", (*api.Context).ServiceInstances)
	router.Get("/services/:serviceId/instances/:instanceId", (*api.Context).GetInstance)
	router.Post("/services/:serviceId/instances", (*api.Context).AddServiceInstance)
	router.Patch("/services/:serviceId/instances/:instanceId", (*api.Context).PatchInstance)
	router.Delete("/services/:serviceId/instances/:instanceId", (*api.Context).DeleteInstance)

	router.Get("/applications", (*api.Context).Applications)
	router.Get("/applications/:applicationId", (*api.Context).GetApplication)
	router.Post("/applications", (*api.Context).AddApplication)
	router.Patch("/applications/:applicationId", (*api.Context).PatchApplication)
	router.Delete("/applications/:applicationId", (*api.Context).DeleteApplication)

	router.Get("/applications/instances", (*api.Context).ApplicationsInstances)
	router.Get("/applications/:applicationId/instances", (*api.Context).ApplicationInstances)
	router.Get("/applications/:applicationId/instances/:instanceId", (*api.Context).GetInstance)
	router.Post("/applications/:applicationId/instances", (*api.Context).AddApplicationInstance)
	router.Patch("/applications/:applicationId/instances/:instanceId", (*api.Context).PatchInstance)
	router.Delete("/applications/:applicationId/instances/:instanceId", (*api.Context).DeleteInstance)

	router.Get("/images", (*api.Context).Images)
	router.Get("/images/:imageId", (*api.Context).GetImage)
	router.Post("/images", (*api.Context).AddImage)
	router.Patch("/images/:imageId", (*api.Context).PatchImage)
	router.Delete("/images/:imageId", (*api.Context).DeleteImage)

	router.Get("/instances", (*api.Context).Instances)
	router.Get("/instances/:instanceId", (*api.Context).GetInstance)
	router.Get("/instances/:instanceId/bindings", (*api.Context).GetInstanceBindings)
	router.Delete("/instances/:instanceId", (*api.Context).DeleteInstance)
	router.Patch("/instances/:instanceId", (*api.Context).PatchInstance)

	router.Get("/templates", (*api.Context).Templates)
	router.Post("/templates", (*api.Context).AddTemplate)
	router.Get("/templates/:templateId", (*api.Context).GetTemplate)
	router.Delete("/templates/:templateId", (*api.Context).DeleteTemplate)
	router.Patch("/templates/:templateId", (*api.Context).PatchTemplate)

}
