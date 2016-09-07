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
	router.Middleware(*Context).BasicAuthorizeMiddleware)
	router.Middleware(*Context).OrganizationSetupMiddleware)

	router.Get("/services", *Context).Services)
	router.Get("/services/:serviceId", *Context).GetService)
	router.Post("/services", *Context).AddService)
	router.Patch("/services/:serviceId", *Context).PatchService)
	router.Delete("/services/:serviceId", *Context).DeleteService)

	router.Get("/services/:serviceId/plans", *Context).Plans)
	router.Get("/services/:serviceId/plans/:planId", *Context).GetPlan)
	router.Post("/services/:serviceId/plans", *Context).AddPlan)
	router.Patch("/services/:serviceId/plans/:planId", *Context).PatchPlan)
	router.Delete("/services/:serviceId/plans/:planId", *Context).DeletePlan)

	router.Get("/services/instances", *Context).ServicesInstances)
	router.Get("/services/:serviceId/instances", *Context).ServiceInstances)
	router.Get("/services/:serviceId/instances/:instanceId", *Context).GetInstance)
	router.Post("/services/:serviceId/instances", *Context).AddServiceInstance)
	router.Patch("/services/:serviceId/instances/:instanceId", *Context).PatchInstance)
	router.Delete("/services/:serviceId/instances/:instanceId", *Context).DeleteInstance)

	router.Get("/applications", *Context).Applications)
	router.Get("/applications/:applicationId", *Context).GetApplication)
	router.Post("/applications", *Context).AddApplication)
	router.Patch("/applications/:applicationId", *Context).PatchApplication)
	router.Delete("/applications/:applicationId", *Context).DeleteApplication)

	router.Get("/applications/instances", *Context).ApplicationsInstances)
	router.Get("/applications/:applicationId/instances", *Context).ApplicationInstances)
	router.Get("/applications/:applicationId/instances/:instanceId", *Context).GetInstance)
	router.Post("/applications/:applicationId/instances", *Context).AddApplicationInstance)
	router.Patch("/applications/:applicationId/instances/:instanceId", *Context).PatchInstance)
	router.Delete("/applications/:applicationId/instances/:instanceId", *Context).DeleteInstance)

	router.Get("/images", *Context).Images)
	router.Get("/images/:imageId", *Context).GetImage)
	router.Post("/images", *Context).AddImage)
	router.Patch("/images/:imageId", *Context).PatchImage)
	router.Delete("/images/:imageId", *Context).DeleteImage)

	router.Get("/instances", *Context).Instances)
	router.Get("/instances/:instanceId", *Context).GetInstance)
	router.Get("/instances/:instanceId/bindings", *Context).GetInstanceBindings)
	router.Delete("/instances/:instanceId", *Context).DeleteInstance)
	router.Patch("/instances/:instanceId", *Context).PatchInstance)

	router.Get("/templates", *Context).Templates)
	router.Post("/templates", *Context).AddTemplate)
	router.Get("/templates/:templateId", *Context).GetTemplate)
	router.Delete("/templates/:templateId", *Context).DeleteTemplate)
	router.Patch("/templates/:templateId", *Context).PatchTemplate)
}
