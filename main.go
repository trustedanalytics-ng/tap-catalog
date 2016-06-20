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
	"net/http"
	"os"
	"time"

	"github.com/gocraft/web"

	"github.com/trustedanalytics/tap-catalog/api"
)

type appHandler func(web.ResponseWriter, *web.Request) error

func main() {
	rand.Seed(time.Now().UnixNano())

	r := web.New(api.Context{})
	r.Middleware(web.LoggerMiddleware)
	basicAuthRouter := r.Subrouter(api.Context{}, "/v1")
	basicAuthRouter.Middleware((*api.Context).BasicAuthorizeMiddleware)

	r.Get("/", (*api.Context).Index)

	r.Error((*api.Context).Error)
	basicAuthRouter.Get("/services", (*api.Context).Services)
	basicAuthRouter.Get("/services/:serviceId", (*api.Context).GetService)
	basicAuthRouter.Post("/services", (*api.Context).AddService)
	basicAuthRouter.Put("/services/:serviceId", (*api.Context).UpdateService)
	basicAuthRouter.Patch("/services/:serviceId", (*api.Context).UpdateService)
	basicAuthRouter.Delete("/services/:serviceId", (*api.Context).DeleteService)

	basicAuthRouter.Get("/services/:serviceId/plans", (*api.Context).Plans)
	basicAuthRouter.Get("/services/:serviceId/plans/:planId", (*api.Context).GetPlan)
	basicAuthRouter.Post("/services/:serviceId/plans", (*api.Context).AddPlan)
	basicAuthRouter.Patch("/services/:serviceId/plans/:planId", (*api.Context).UpdatePlan)
	basicAuthRouter.Put("/services/:serviceId/plans/:planId", (*api.Context).UpdatePlan)
	basicAuthRouter.Delete("/services/:serviceId/plans/:planId", (*api.Context).DeletePlan)

	basicAuthRouter.Get("/services/:serviceId/instances", (*api.Context).Instances)
	basicAuthRouter.Get("/services/:serviceId/instances/:instanceId", (*api.Context).GetInstance)
	basicAuthRouter.Post("/services/:serviceId/instances", (*api.Context).AddInstance)
	basicAuthRouter.Patch("/services/:serviceId/instances/:instanceId", (*api.Context).UpdateInstance)
	basicAuthRouter.Put("/services/:serviceId/instances/:instanceId", (*api.Context).UpdateInstance)
	basicAuthRouter.Delete("/services/:serviceId/instances/:instanceId", (*api.Context).DeleteInstance)

	basicAuthRouter.Get("/applications", (*api.Context).Applications)
	basicAuthRouter.Get("/applications/:applicationId", (*api.Context).GetApplication)
	basicAuthRouter.Post("/applications", (*api.Context).AddApplication)
	basicAuthRouter.Put("/applications/:applicationId", (*api.Context).UpdateApplication)
	basicAuthRouter.Patch("/applications/:applicationId", (*api.Context).UpdateApplication)
	basicAuthRouter.Delete("/applications/:applicationId", (*api.Context).DeleteApplication)

	basicAuthRouter.Get("/applications/:applicationId/instances", (*api.Context).Instances)
	basicAuthRouter.Get("/applications/:applicationId/instances/:instanceId", (*api.Context).GetInstance)
	basicAuthRouter.Post("/applications/:applicationId/instances", (*api.Context).AddInstance)
	basicAuthRouter.Put("/applications/:applicationId/instances/:instanceId", (*api.Context).UpdateInstance)
	basicAuthRouter.Patch("/applications/:applicationId/instances/:instanceId", (*api.Context).UpdateInstance)
	basicAuthRouter.Delete("/applications/:applicationId/instances/:instanceId", (*api.Context).DeleteInstance)

	basicAuthRouter.Get("/instances", (*api.Context).Instances)
	basicAuthRouter.Get("/instances/:instanceId", (*api.Context).GetInstance)
	basicAuthRouter.Delete("/instances/:instanceId", (*api.Context).DeleteInstance)
	basicAuthRouter.Patch("/instances/:instanceId", (*api.Context).UpdateInstance)

	basicAuthRouter.Post("/instances/:instanceId/bindings", (*api.Context).AddInstanceBinding)
	basicAuthRouter.Delete("/instances/:instanceId/bindings/:bindingId", (*api.Context).DeleteInstanceBinding)
	basicAuthRouter.Post("/instances/:instanceId/meta", (*api.Context).AddInstanceMetadata)
	basicAuthRouter.Delete("/instances/:instanceId/meta/:key", (*api.Context).DeleteInstanceMetadata)

	basicAuthRouter.Get("/templates", (*api.Context).Templates)
	basicAuthRouter.Post("/templates", (*api.Context).AddTemplate)
	basicAuthRouter.Get("/templates/:templateId", (*api.Context).GetTemplate)
	basicAuthRouter.Delete("/templates/:templateId", (*api.Context).DeleteTemplate)
	basicAuthRouter.Put("/templates/:templateId", (*api.Context).UpdateTemplate)

	port := os.Getenv("PORT")
	log.Println("Will listen on:", "0.0.0.0:"+port)
	err := http.ListenAndServe("0.0.0.0:"+port, r)
	if err != nil {
		log.Println("Couldn't serve app on port ", port, " Application will be closed now.")
	}
}
