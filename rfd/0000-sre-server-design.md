---
authors: Diego Rodriguez (diego@rodilla.me)
state: draft
---

# RFD 0 - SRE Server Design

## What
Initial design for SRE server HTTP API server that provides interactions with a
Kubernetes cluster.

## Why

As part of the SRE's job, they are required to interact with Kubernetes clusters
on regular basis. However, existing tools such as `kubectl` or using the
Kubernetes API directly, can pose considerable friction with the overwhelming
amount options that are available in those tools. Also, onboarding new SREs with
those tools can be non-trivial, requring them to spend weeks learning the
terminology and how to actually use them. To reduce the friction of interaction
with Kubernetes clusters and ramp up time, we therefore propose a new service
`sre-server` to help ease these pains. The `sre-server` will provide convenience
(and safety) on top of the existing Kuberenets API to simplify interactions and
thus save the SRE's time while exponentially increasing their productivity.

## Details

### Scope
After collecting feedback from the SREs, it turns out that they spend a large
chunk of their time working with Kubernetes deployments. Hence, for the first
pass, we will focus primarily on providing niceities over the Kubernetes
deployments. Depending on the success of this project, we may decide to further
expand the scope to other Kubernetes resources as well.


### API
All responses and request bodies shall be formated in JSON.


In general, all non 2XX/3XX statuses shall return a response as follows:

```json
{
  "message" : "Message explaining the status."
}
```


#### **`GET`** `/deployments[/{namespace}]`
Lists all of the deployments of the Kubernetes cluster by namespace

*Parameters*
| Path Parameter | Description | Required |
| -------------- | ----------- | -------- |
| `namespace`    | The namespace to limit the scope of the deployments to list. Omitting it the namespace will retrieve all the deployments of cluster.  | `false` |


*Response*
```json
[
  {
    "namespace": "default",
    "deployments": [
      "foo",
      "yeet",
      ...
    ]
  },
  ...
]
```

| HTTP Status Code | Description |
| ---------------- | ----------- |
| 200 | Successfully retrieved the list of deployments |
| 404 | The provided namespace does not exist |




#### **`GET`** `/deployments/{namespace}/{deployment}/replicas`
Retrieves the deployment replica count by name in the provided namespace.

For this particular endpoint, to prevent overwhelming the Kubernetes cluster API
with too many of these requests (impacting its performance), the replica count
shall be cached. For the initial implementation, it shall use an in-memory cache
implemented as a map for simplicity. Due to Go maps being [unsafe
for concurrent use](https://go.dev/doc/faq#atomic_maps), we'll need to synchronize
the read/write access (which may impact performance). Eventually, for better
efficiency and performance across multiple replicas, we'd want to use something
like Redis to centralize the cache.

In addition to caching the replica count, when a new entry is added, we will add
a deployment watcher to check for any updates to the replica count to keep the
cache up to date. In the event that the deployment no longer exists, the
associated cache entry shall eventually be removed by the watcher.

*Parameters*
| Path Parameter | Description | Required |
| -------------- | ----------- | -------- |
| `namespace`    | The deployment namespace. | `true` |
| `deployment`   | The name of the deployment. | `true` |


*Response*
```json
{
  "namespace": "default",
  "deployment": "cool",
  "replicas": 7
}
```

| HTTP Status Code | Description |
| ---------------- | ----------- |
| 200 | Successfully retrieved the replica count for the deployment |
| 404 | The provided namespace or deployment does not exist. |




#### **`PUT`** `/deployments/{namespace}/{deployment}/replicas`
Update the replica count for the given deployment.

*Parameters*
| Path Parameter | Description | Required |
| -------------- | ----------- | -------- |
| `namespace`    | The deployment namespace. | `true` |
| `deployment`   | The name of the deployment. | `true` |

| Body Parameter | Description | Required |
| -------------- | ----------- | -------- |
| `replicas`    | The new number of replicas for the deployment. | `true` |

```json
{
  "replicas": 20
}
```

*Response*
```json
{
  "message": "Successfully updated the replica count for the deployment"
}
```

| HTTP Status Code | Description |
| ---------------- | ----------- |
| 200 | Successfully updated the replica count for the deployment |
| 400 | Incorrect value for the number of replicas was set. |
| 404 | The provided namespace or deployment does not exist. |


### Security
We make use of mTLS for secured communications between the client and the
`sre-server`. With mTLS, both the server and client verify each other's
certificates (authenticate). For this initial design, in order to simplify
usage and development, we opted to make use of the same public and private
keys for both the server and the clients. However, for production, we will
want to use an existing internal CA (or use something like Hashicorp Vault)
to generate the TLS certifcate key pairs for each client and server for a
secure setup. For generating the TLS certificate key pairs themselves, we
recommend using Ed25519 for improved performance (and sometimes security) over
RSA and ECDSA. A key pair that can be used by both clients and the server
can be done though:

```shell
# For `-days` parameter, you can change the default to increase how long the
# certificate is valid for. We default it to a low number to catch accidental
# usage in production.
#
# For the `-addext "subjectAltName` option, we add a subject alternative name
# to prevent Go based clients from complaining about
# "certificate relies on legacy Common Name field, use SANs instead". In
# general, the industry seems to be moving towards using SANs over the common
# name.
openssl req \
  -newkey ed25519 \
  -new \
  -nodes \
  -x509 \
  -days 1 \
  -out public-certificate.pem \
  -keyout private-key.pem \
  -subj "/C=US/ST=Cali/L=Somewhere/O=Your Organization/OU=Your Unit/CN=localhost" \
  -addext "subjectAltName = DNS:localhost"
```

The `sre-server` shall enforce TLS 1.3 as the only option to enforce the
current best security practices, which is supported by the majority of
clients today. Using TLS 1.3 in Go also means that we're not able to choose or
prefer any particular cipher suite (see https://pkg.go.dev/crypto/tls#Config),
but default to the following:

-  `TLS_AES_128_GCM_SHA256`
-  `TLS_AES_256_GCM_SHA384`
-  `TLS_CHACHA20_POLY1305_SHA256`

Should we require to use TLS 1.2 to support legacy clients, then we still
recommend to use the above cipher suites as the preferred for the current
security best practices.


### Development


### Build, Release, and Delivery
