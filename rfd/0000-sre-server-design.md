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
The APIs are documented formatted in the following structure:

```
**`VERB`** `/{service}/{subpaths}...`

Description of the API action goes here.

*Parameters*
| [Query | Body | Path] Paramater | Description | Required |
| --------------- | ----------- | -------- |
| `foo`           | Parameter description for `foo` | `true` |
| `yeet`          | The distance to yeet it for. Defaults to 50 meters | `false` |


*Response*

`` ```json
{
  "successful": true
}
``` ``

| HTTP Status Code | Description |
| ---------------- | ----------- |
| 200 | Successful request
...
```

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
| Path Paramater | Description | Required |
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

*Parameters*
| Path Paramater | Description | Required |
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
| Path Paramater | Description | Required |
| -------------- | ----------- | -------- |
| `namespace`    | The deployment namespace. | `true` |
| `deployment`   | The name of the deployment. | `true` |

| Body Paramater | Description | Required |
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
| 204 | Successfully updated the replica count for the deployment |
| 400 | Incorrect value for the number of replicas was set. |
| 404 | The provided namespace or deployment does not exist. |
