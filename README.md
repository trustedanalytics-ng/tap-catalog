# TAP Catalog

Catalog is a microservice developed to be a part of TAP platform.
The Catalog acts as the central registry and coordination point for the entire TAP NG instance.
It provides an integrated, logical, view of the platform offerings including their deployment status, state and dependencies.
It is implemented on top of a distributed configuration system (etcd) that provides key primitives such as replicated state and atomic updates.

The specific details stored include:
* The list of available applications and services including details such as image references and desired replication level
* The list of service instances
* application and service instance bindings
* defined users and organizations

## REQUIREMENTS

### Binary
etcd instance is required to run Catalog.
Firstly, you need to download it from github:
```
curl -L  https://github.com/coreos/etcd/releases/download/v2.3.6/etcd-v2.3.6-linux-amd64.tar.gz -o etcd-v2.3.6-linux-amd64.tar.gz
tar xzvf etcd-v2.3.6-linux-amd64.tar.gz
```
And finally run it:
```
cd etcd-v2.3.6-linux-amd64
./etcd --listen-peer-urls=http://localhost:2382 --advertise-client-urls http://localhost:2379 -listen-client-urls=http://localhost:2379
```

### Compilation
* git (for pulling repository)
* go >= 1.6

## Compilation
To build project:
```
  git clone https://github.com/intel-data/tapng-catalog
  cd tapng-catalog
  make build_anywhere
```
Binaries are available in ./application/

## USAGE

To provide IP and port for the application, you have to setup system environment variables
```
export BIND_ADDRESS=127.0.0.1
export PORT=80
```

Catalog endpoints are documented in swagger.yaml file.
Below you can find sample Catalog usage.

#### Creating template
```
curl -XPOST -H 'Content-type: application/json' http://127.0.0.1/api/v1/templates -d "{}" --user admin:password
{"templateId":"4fcee3a1-201a-4db2-782e-7eae3e654535","state":"IN_PROGRESS","auditTrail":{"createdOn":1472563175,"createdBy":"admin","lastUpdatedOn":1472563175,"lastUpdateBy":"admin"}}
```

#### Creating service
```
curl -XPOST -H 'Content-type: application/json' http://127.0.0.1/api/v1/services -d '{"state":"DEPLOYING", "templateId":"4fcee3a1-201a-4db2-782e-7eae3e654535", "name": "logstash", "description": "logstash description", "bindable": true, "tags": ["logstash", "logstash14"]}' --user admin:password
{"id":"0f506b25-9cb2-4d87-4b26-b6d702714b5f","name":"logstash","description":"logstash description","bindable":true,"templateId":"4fcee3a1-201a-4db2-782e-7eae3e654535","state":"DEPLOYING","plans":null,"auditTrail":{"createdOn":1472568204,"createdBy":"admin","lastUpdatedOn":1472568204,"lastUpdateBy":"admin"},"metadata":null}%
```

#### Creating service instance
```
curl -XPOST -H 'Content-type: application/json' "http://127.0.0.1/api/v1/services/0f506b25-9cb2-4d87-4b26-b6d702714b5f/instances?isServiceBroker=false" -d '{"name":"logstash", "type":"SERVICE", "classId":"0f506b25-9cb2-4d87-4b26-b6d702714b5f"}' --user admin:password
{"id":"b1e18756-fc55-486b-5c7b-9a7b7ef30d10","name":"logstash","type":"SERVICE","classId":"0f506b25-9cb2-4d87-4b26-b6d702714b5f","bindings":null,"metadata":null,"state":"REQUESTED","auditTrail":{"createdOn":1472568909,"createdBy":"admin","lastUpdatedOn":1472568909,"lastUpdateBy":"admin"}}
```

#### Listing service instances
```
curl http://127.0.0.1/api/v1/services/instances --user admin:password
{"id":"b1e18756-fc55-486b-5c7b-9a7b7ef30d10","name":"logstash","type":"SERVICE","classId":"0f506b25-9cb2-4d87-4b26-b6d702714b5f","bindings":null,"metadata":null,"state":"REQUESTED","auditTrail":{"createdOn":1472568909,"createdBy":"admin","lastUpdatedOn":1472568909,"lastUpdateBy":"admin"}}
```

#### Removing service instance
```
curl -XDELETE "http://127.0.0.1/api/v1/instances/b1e18756-fc55-486b-5c7b-9a7b7ef30d10" --user admin:password
```
