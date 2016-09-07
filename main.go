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
	"sync"
	"time"

	"github.com/gocraft/web"

	"github.com/trustedanalytics/tap-catalog/api"
	"github.com/trustedanalytics/tap-catalog/data"
	httpGoCommon "github.com/trustedanalytics/tap-go-common/http"
	"github.com/trustedanalytics/tap-go-common/util"
)

type appHandler func(web.ResponseWriter, *web.Request) error

var waitGroup = &sync.WaitGroup{}

func main() {
	rand.Seed(time.Now().UnixNano())

	err := (&data.RepositoryConnector{}).CreateDirs(os.Getenv("CORE_ORGANIZATION"))
	go util.TerminationObserver(waitGroup, "Catalog")
	if err != nil {
		log.Fatalln("Can't create directories in ETCD!", err)
	}

	context := api.Context{}

	r := web.New(context)
	r.Middleware(web.LoggerMiddleware)
	r.Get("/healthz", context.GetCatalogHealth)

	apiRouter := r.Subrouter(context, "/api")

	basicAuthRouter := apiRouter.Subrouter(context, "/v1")
	route(basicAuthRouter, &context)

	// for testing purpose, where v1 is current version
	v1AliasRouter := apiRouter.Subrouter(context, "/v1.0")
	route(v1AliasRouter, &context)

	r.Get("/", (*api.Context).Index)
	r.Error((*api.Context).Error)

	if os.Getenv("CATALOG_SSL_CERT_FILE_LOCATION") != "" {
		httpGoCommon.StartServerTLS(os.Getenv("CATALOG_SSL_CERT_FILE_LOCATION"),
			os.Getenv("CATALOG_SSL_KEY_FILE_LOCATION"), r)
	} else {
		httpGoCommon.StartServer(r)
	}

}

func route(router *web.Router, context *api.Context) {
	router.Middleware((*context).BasicAuthorizeMiddleware)
	router.Middleware((*context).OrganizationSetupMiddleware)

	router.Get("/services", (*context).Services)
	router.Get("/services/:serviceId", (*context).GetService)
	router.Post("/services", (*context).AddService)
	router.Patch("/services/:serviceId", (*context).PatchService)
	router.Delete("/services/:serviceId", (*context).DeleteService)

	router.Get("/services/:serviceId/plans", (*context).Plans)
	router.Get("/services/:serviceId/plans/:planId", (*context).GetPlan)
	router.Post("/services/:serviceId/plans", (*context).AddPlan)
	router.Patch("/services/:serviceId/plans/:planId", (*context).PatchPlan)
	router.Delete("/services/:serviceId/plans/:planId", (*context).DeletePlan)

	router.Get("/services/instances", (*context).ServicesInstances)
	router.Get("/services/:serviceId/instances", (*context).ServiceInstances)
	router.Get("/services/:serviceId/instances/:instanceId", (*context).GetInstance)
	router.Post("/services/:serviceId/instances", (*context).AddServiceInstance)
	router.Patch("/services/:serviceId/instances/:instanceId", (*context).PatchInstance)
	router.Delete("/services/:serviceId/instances/:instanceId", (*context).DeleteInstance)

	router.Get("/applications", (*context).Applications)
	router.Get("/applications/:applicationId", (*context).GetApplication)
	router.Post("/applications", (*context).AddApplication)
	router.Patch("/applications/:applicationId", (*context).PatchApplication)
	router.Delete("/applications/:applicationId", (*context).DeleteApplication)

	router.Get("/applications/instances", (*context).ApplicationsInstances)
	router.Get("/applications/:applicationId/instances", (*context).ApplicationInstances)
	router.Get("/applications/:applicationId/instances/:instanceId", (*context).GetInstance)
	router.Post("/applications/:applicationId/instances", (*context).AddApplicationInstance)
	router.Patch("/applications/:applicationId/instances/:instanceId", (*context).PatchInstance)
	router.Delete("/applications/:applicationId/instances/:instanceId", (*context).DeleteInstance)

	router.Get("/images", (*context).Images)
	router.Get("/images/:imageId", (*context).GetImage)
	router.Post("/images", (*context).AddImage)
	router.Patch("/images/:imageId", (*context).PatchImage)
	router.Delete("/images/:imageId", (*context).DeleteImage)

	router.Get("/instances", (*context).Instances)
	router.Get("/instances/:instanceId", (*context).GetInstance)
	router.Get("/instances/:instanceId/bindings", (*context).GetInstanceBindings)
	router.Delete("/instances/:instanceId", (*context).DeleteInstance)
	router.Patch("/instances/:instanceId", (*context).PatchInstance)

	router.Get("/templates", (*context).Templates)
	router.Post("/templates", (*context).AddTemplate)
	router.Get("/templates/:templateId", (*context).GetTemplate)
	router.Delete("/templates/:templateId", (*context).DeleteTemplate)
	router.Patch("/templates/:templateId", (*context).PatchTemplate)
}
