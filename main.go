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
	"time"
	"os"

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
	basicAuthRouter.Get("/services/:service_id", (*api.Context).GetService)
	basicAuthRouter.Post("/services", (*api.Context).AddService)
	basicAuthRouter.Put("/services/:service_id", (*api.Context).UpdateService)
	basicAuthRouter.Delete("/services/:service_id", (*api.Context).DeleteService)

	basicAuthRouter.Get("/services/:service_id/plans", (*api.Context).Plans)
	basicAuthRouter.Get("/services/:service_id/plans/:plan_id", (*api.Context).GetPlan)
	basicAuthRouter.Post("/services/:service_id/plans", (*api.Context).AddPlan)
	basicAuthRouter.Put("/services/:service_id/plans/:plan_id", (*api.Context).UpdatePlan)
	basicAuthRouter.Delete("/services/:service_id/plans/:plan_id", (*api.Context).DeletePlan)

	basicAuthRouter.Get("/services/:service_id/instances", (*api.Context).Instances)
	basicAuthRouter.Get("/services/:service_id/instances/:instances_id", (*api.Context).GetInstance)
	basicAuthRouter.Post("/services/:service_id/instances", (*api.Context).AddInstance)
	basicAuthRouter.Put("/services/:service_id/instances/:instances_id", (*api.Context).UpdateInstance)
	basicAuthRouter.Delete("/services/:service_id/instances/:instances_id", (*api.Context).DeleteInstance)

	basicAuthRouter.Get("/applications", (*api.Context).Applications)
	basicAuthRouter.Get("/applications/:application_id", (*api.Context).GetApplication)
	basicAuthRouter.Post("/applications", (*api.Context).AddApplication)
	basicAuthRouter.Put("/applications/:application_id", (*api.Context).UpdateApplication)
	basicAuthRouter.Delete("/applications/:application_id", (*api.Context).DeleteApplication)

	basicAuthRouter.Get("/applications/:application_id/instances", (*api.Context).Instances)
	basicAuthRouter.Get("/applications/:application_id/instances/:instances_id", (*api.Context).GetInstance)
	basicAuthRouter.Post("/applications/:application_id/instances", (*api.Context).AddInstance)
	basicAuthRouter.Put("/applications/:application_id/instances/:instances_id", (*api.Context).UpdateInstance)
	basicAuthRouter.Delete("/applications/:application_id/instances/:instances_id", (*api.Context).DeleteInstance)

	basicAuthRouter.Get("/instances", (*api.Context).Instances)
	basicAuthRouter.Get("/instances/:instances_id", (*api.Context).GetInstance)
	basicAuthRouter.Get("/instances/:instances_id", (*api.Context).DeleteInstance)
	basicAuthRouter.Put("/instances/:instances_id", (*api.Context).UpdateInstance)
	basicAuthRouter.Put("/instances/:instances_id/states/state_name", (*api.Context).UpdateInstanceState)
	basicAuthRouter.Put("/instances/:instances_id/bindings/:binding_id", (*api.Context).AddInstanceBinding)
	basicAuthRouter.Delete("/instances/:instances_id/bindings/:binding_id", (*api.Context).DeleteInstanceBinding)
	basicAuthRouter.Put("/instances/:instances_id/meta/:key", (*api.Context).AddInstanceMetadata)
	basicAuthRouter.Delete("/instances/:instances_id/meta/:key", (*api.Context).DeleteInstanceMetadata)

	basicAuthRouter.Get("/templates", (*api.Context).Templates)
	basicAuthRouter.Get("/templates/:template_id", (*api.Context).GetTemplate)
	basicAuthRouter.Delete("/templates/:template_id", (*api.Context).DeleteTemplate)
	basicAuthRouter.Put("/templates/:templates_id/states/state_name", (*api.Context).UpdateTemplateState)

	port := os.Getenv("PORT")
	log.Println("Will listen on:", "0.0.0.0:"+port)
	err := http.ListenAndServe("0.0.0.0:"+port, r)
	if err != nil {
		log.Println("Couldn't serve app on port ", port, " Application will be closed now.")
	}
}
