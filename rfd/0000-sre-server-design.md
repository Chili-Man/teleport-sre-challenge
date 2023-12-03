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

All responses shall be formated in JSON for easier machine parsing.


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
| 200 | Successfully retrieved the list of deployments
| 404 | The provided namespace does not exist
