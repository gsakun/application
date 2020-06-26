# 接口列表

接口文档包含如下接口：

- login
- clusters
- projects
- namespace
- configmap
  - 创建configmap
  - 更新configmap
  - 删除configmap
- workloads
  - 获取workload
  - 查询workload负载
- pods
  - 获取pod
  - 查询pod负载
  - 删除pod
- applicationConfigurationTemplate
- application
  - 对接开发功能
  - 数据定义
  - 创建application
  - 更新application
  - 删除application

# login(样例)

## 接口

```shell
https://10.10.111.54:444/v3-public/localProviders/local?action=login
```

## verb

```
"Methods": [
        "POST"
  ]
```

## headers

| key          | value            |
| ------------ | ---------------- |
| content-type | application/json |

## body

```json
{"username":"admin","password":"111111"}
```

## respons

```json
{
    "authProvider": "local",
    "baseType": "token",
    "created": "2019-11-28T07:43:52Z",
    "id": "token-4b8v4",
    "labels": {
        "authn.management.cattle.io/token-userId": "user-5js2x",
        "cattle.io/creator": "norman"
    },
    "name": "token-4b8v4",
    "token": "token-4b8v4:g2dl8r98jzhk5wslwqfw87hks4xt9s927xqj86r5clwd76bk8q6rpb",
    "ttl": 0,
    "type": "token",
    "userId": "user-5js2x",
    "uuid": "d28c077f-11b2-11ea-9e36-0242ac110003"，
    ...
}
```

# 获取Clusters

## 接口

```
/v3/clusters?limit=-1&sort=name
```

## verb

```
"Methods": [
        "GET"
  ]
```

## headers

| key           | value                                                        |
| ------------- | ------------------------------------------------------------ |
| content-type  | application/json                                             |
| Authorization | Bearer token-5tgq4:xnq57jc7vgmjlcwzhtl8qdfbftt7tns7jf4bk6p4tstrh6bjx9whsd |

## response

```json
{
    "type": "collection",
    "links": {
        "self": "https://10.10.111.54:444/v3/clusters"
    },
    "createTypes": {
        "cluster": "https://10.10.111.54:444/v3/clusters"
    },
    "resourceType": "cluster",
    "data": [
        {
            "id": "c-x67ps",
            "links": {
                "namespaces": "https://10.10.111.54:444/v3/clusters/c-x67ps/namespaces",
                "projects": "https://10.10.111.54:444/v3/clusters/c-x67ps/projects",
                ...
            },
            "name": "54",
            "state": "active",
            "type": "cluster",
            "uuid": "ebcbcc89-04f9-11ea-bdf4-0242ac110003",
            ...
        }
    ]
}
```

解析json获取Cluster名称以及ID  {data[].name}  {data[].id} 

# 获取Project列表

## 接口

```
/v3/clusters/$(clusterid)/projects
```

## verb

```
"Methods": [
        "GET"
  ]
```

## headers

| key           | value                                                        |
| ------------- | ------------------------------------------------------------ |
| content-type  | application/json                                             |
| Authorization | Bearer token-5tgq4:xnq57jc7vgmjlcwzhtl8qdfbftt7tns7jf4bk6p4tstrh6bjx9whsd |

## response

```json
{
    "type": "collection",
    "links": {
        "self": "https://10.10.111.54:444/v3/clusters/c-x67ps/projects"
    },
    "createTypes": {
        "project": "https://10.10.111.54:444/v3/projects"
    },

    "resourceType": "project",
    "data": [
        {
            "created": "2019-11-12T03:10:03Z",
            "createdTS": 1573528203000,
            "creatorId": "user-5js2x",
            "description": "System project created for the cluster",
            "id": "c-x67ps:p-gxbks",
            "name": "System",
            "state": "active",
            "type": "project",
            "uuid": "ebcfb498-04f9-11ea-bdf4-0242ac110003",
            ...
        },
        ...
    ]
}
```

解析json获取Project名称以及ID  {data[].name}  {data[].id} 

# Namespace

## 获取namespaces

### 接口

```shell
/v3/clusters/$(clusterid)/namespaces?limit=-1&sort=name
```

### verb

```
"Methods": [
        "GET"
  ]
```

### headers

| key           | value                                                        |
| ------------- | ------------------------------------------------------------ |
| content-type  | application/json                                             |
| Authorization | Bearer token-5tgq4:xnq57jc7vgmjlcwzhtl8qdfbftt7tns7jf4bk6p4tstrh6bjx9whsd |

### response

```json
{
    "type": "collection",
    "links": {
        "self": "https://10.10.111.54:444/v3/cluster/c-x67ps/namespaces"
    },
    "createTypes": {
        "namespace": "https://10.10.111.54:444/v3/cluster/c-x67ps/namespaces"
    },
    "resourceType": "namespace",
    "data": [
        {
            "baseType": "namespace",
            "created": "2019-10-28T01:56:51Z",
            "createdTS": 1572227811000,
            "creatorId": null,
            "id": "kube-system",
            "labels": {
                "field.cattle.io/projectId": "p-5kb4h"
            },
            "name": "kube-system",
            "projectId": "c-86x26:p-5kb4h",
            "state": "active",
            "transitioning": "no",
            "transitioningMessage": "",
            "type": "namespace",
            "uuid": "35a9205d-f926-11e9-ba4d-fa163ecee4c9",
            ...
        },
        ...
    ]
}
```

## 创建namespace

### 接口

```shell
/v3/clusters/c-x67ps/namespace
```

### verb

```json
"Methods": [
        "POST"
  ]
```

### payload

```json
{
	"type": "namespace", # 必选
	"name": string, # 必选
	"clusterId": string, # 必选 集群id "c-x67ps"  
	"projectId": string # Projectid # 必选 示例"c-x67ps:p-sjjrk",
	"resourceQuota": {
		pods: string, # 可选
		services: string, # 可选
		replicationControllers: string, # 可选
		secrets: string, # 可选
		configMaps: string, # 可选
		persistentVolumeClaims: string, # 可选
		servicesNodePorts: string, # 可选
		servicesLoadBalancers: string, # 可选
		requestsCpu: string, # 可选
		requestsMemory: string, # 可选
		requestsStorage: string, # 可选
		limitsCpu: string, # 可选
		limitsMemory: string # 可选
	} # 可选
}
```



### headers

| key           | value                                                        |
| ------------- | ------------------------------------------------------------ |
| content-type  | application/json                                             |
| Authorization | Bearer token-5tgq4:xnq57jc7vgmjlcwzhtl8qdfbftt7tns7jf4bk6p4tstrh6bjx9whsd |

### response

```json
{
    "baseType": "namespace",
    "created": "2019-12-06T12:53:29Z",
    "createdTS": 1575636809000,
    "creatorId": "user-5js2x",
    "id": "test-1111",
    "labels": {
        "cattle.io/creator": "norman"
    },
    "links": {
        "remove": "https://10.10.111.54:444/v3/cluster/c-x67ps/namespaces/test-1111",
        "self": "https://10.10.111.54:444/v3/cluster/c-x67ps/namespaces/test-1111",
        "update": "https://10.10.111.54:444/v3/cluster/c-x67ps/namespaces/test-1111",
        "yaml": "https://10.10.111.54:444/v3/cluster/c-x67ps/namespaces/test-1111/yaml"
    },
    "name": "test-1111",
    "projectId": "c-x67ps:p-sjjrk",
    "state": "activating",
    "type": "namespace",
    "uuid": "6684967d-1827-11ea-96e2-fa163ecee4c9,
    ...
}
```

# ConfigMap

## 创建configmap

### 接口

```
/v3/projects/${projectid}/configmap
```

### headers

| key           | value                                                        |
| ------------- | ------------------------------------------------------------ |
| content-type  | application/json                                             |
| Authorization | Bearer token-5tgq4:xnq57jc7vgmjlcwzhtl8qdfbftt7tns7jf4bk6p4tstrh6bjx9whsd |

### verb

```json
"Methods": [
        "POST"
  ]
```

### payload

```json
{
	"type": "configMap", //固定参数 
	"data": "map[string]string", // key:value 
	"name": "string", // configmap name
	"namespaceId": "string" //namespace name
}

示例
{
    "type":"configMap",
    "data":{
        "test":"test"
    },
    "name":"test",
    "namespaceId":"service"
}
```

### response

```json
{
    "annotations":{
    },
    "baseType":"configMap",
    "data":{
        "test":"test"
    },
    "id":"service:test",
    "labels":{
        "cattle.io/creator":"norman"
    },
	...
    "name":"test",
    "namespaceId":"service",
    "projectId":"local:p-rnqx6",
    "type":"configMap",
	...
}
```

## 更新configmap

### 接口

```
/v3/project/${projectid}/configMaps/${namespace}:${configmapname}
```

### headers

| key           | value                                                        |
| ------------- | ------------------------------------------------------------ |
| content-type  | application/json                                             |
| Authorization | Bearer token-5tgq4:xnq57jc7vgmjlcwzhtl8qdfbftt7tns7jf4bk6p4tstrh6bjx9whsd |

### verb

```json
"Methods": [
        "PUT"
  ]
```

### payload

```json
{
	"data": "map[string]string" // key:value 
}
```

### response

```json
{
    "annotations":{
    },
    "baseType":"configMap",
    "data":{
        "test":"test"
    },
    "id":"service:test",
    "labels":{
        "cattle.io/creator":"norman"
    },
	...
    "name":"test",
    "namespaceId":"service",
    "projectId":"local:p-rnqx6",
    "type":"configMap",
	...
}
```

## 删除configmap

### 接口

```
/v3/project/${projectid}/configMaps/${namespace}:${configmapname}
```

### headers

| key           | value                                                        |
| ------------- | ------------------------------------------------------------ |
| content-type  | application/json                                             |
| Authorization | Bearer token-5tgq4:xnq57jc7vgmjlcwzhtl8qdfbftt7tns7jf4bk6p4tstrh6bjx9whsd |

### verb

```json
"Methods": [
        "DELETE"
  ]
```

# Workloads

## 获取workloads

### 接口

```
/v3/project/c-x67ps:p-sjjrk/workloads?limit=-1&sort=name
```

### 扩展

根据namespace 和 名字 过滤出指定workload

```shell
/v3/project/c-x67ps:p-sjjrk/workloads?limit=-1&sort=name&namespaceId=${namespace}&name=${applicationname}-${compomentname}-workload-${version}
## 示例
/v3/project/c-x67ps:p-sjjrk/workloads?limit=-1&sort=name&namespaceId=service&name=fyhtest-fyhtest-workload-v1
```

### verb

```
"Methods": [
        "GET"
  ]
```

### headers

| key           | value                                                        |
| ------------- | ------------------------------------------------------------ |
| content-type  | application/json                                             |
| Authorization | Bearer token-5tgq4:xnq57jc7vgmjlcwzhtl8qdfbftt7tns7jf4bk6p4tstrh6bjx9whsd |

### response

```json
{
    "type": "collection",
    "links": {
        "self": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/workloads"
    },
    "createTypes": {
        "workload": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/workloads"
    },
    "resourceType": "workload",
    "data": [
        {
            "baseType": "workload",
            "containers": [
                {
                    "image": "docker.io/citizenstig/httpbin@sha256:b81c818ccb8668575eb3771de2f72f8a5530b515365842ad374db76ad8bcf875",
                    "imagePullPolicy": "IfNotPresent",
                    "initContainer": false,
                    "name": "httpbin",
                    "ports": [
                        {
                            "containerPort": 8000,
                            "protocol": "TCP",
                            "sourcePort": 0,
                            "type": "/v3/project/schemas/containerPort"
                        }
                    ],
                    "resources": {
                        "type": "/v3/project/schemas/resourceRequirements"
                    },
                    "restartCount": 0,
                    "type": "/v3/project/schemas/container"
                }
            ],
            "created": "2019-12-03T02:39:41Z",
            "createdTS": 1575340781000,
            "creatorId": null,
            "deploymentConfig": {
                "maxSurge": 1,
                "maxUnavailable": 1,
                "minReadySeconds": 0,
                "progressDeadlineSeconds": 2147483647,
                "revisionHistoryLimit": 2147483647,
                "strategy": "RollingUpdate"
            },
            "deploymentStatus": {
                "availableReplicas": 1,
                "conditions": [
                    {
                        "lastTransitionTime": "2019-12-03T02:39:41Z",
                        "lastTransitionTimeTS": 1575340781000,
                        "lastUpdateTime": "2019-12-03T02:39:41Z",
                        "lastUpdateTimeTS": 1575340781000,
                        "message": "Deployment has minimum availability.",
                        "reason": "MinimumReplicasAvailable",
                        "status": "True",
                        "type": "Available"
                    }
                ],
                "observedGeneration": 1,
                "readyReplicas": 1,
                "replicas": 1,
                "type": "/v3/project/schemas/deploymentStatus",
                "unavailableReplicas": 0,
                "updatedReplicas": 1
            },
            "dnsPolicy": "ClusterFirst",
            "hostIPC": false,
            "hostNetwork": false,
            "hostPID": false,
            "id": "deployment:istio-test:httpbin",
            "labels": {
                "app": "httpbin",
                "version": "v1"
            },
            "links": {
                "update": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/workloads/deployment:istio-test:httpbin",
                ...
            },
            "name": "httpbin",
            "namespaceId": "istio-test",
            "paused": false,
            "projectId": "c-x67ps:p-sjjrk",
            "restartPolicy": "Always",
            "scale": 1,
            "schedulerName": "default-scheduler",
            "selector": {
                "matchLabels": {
                    "app": "httpbin",
                    "version": "v1"
                },
                "type": "/v3/project/schemas/labelSelector"
            },
            "state": "active",
            "terminationGracePeriodSeconds": 30,
            "type": "deployment",
            "workloadAnnotations": {
                "deployment.kubernetes.io/revision": "1"，
            },
            "workloadLabels": {
                "app": "httpbin",
                "version": "v1"
            }
        }，
        ...
    ],
    ...
}
```

## 判断workload 状态

```
{
    "data":[
        {
            "deploymentStatus":{
                "availableReplicas":1,
                "observedGeneration":1,
                "readyReplicas":1,
                "replicas":1,
                "type":"/v3/project/schemas/deploymentStatus",
                "unavailableReplicas":0,
                "updatedReplicas":1
            }
    ]
}
```



根据response **data.[].deploymentStatus.unavailableReplicas** 值判断

- 若该字段值为0 则woload状态正常
- 若该字段值不为0 则woload状态异常

## 查询workload负载

### 接口

```shell
/v3/projectmonitorgraphs?action=query
```

### verb

```json
"Methods": [
        "POST"
  ]
```

### payload

```json
{
	"filters": {
		"resourceType": "workload", ## 资源类型
		"projectId": "c-gwnnw:p-5ggc4" ## projectid
	},
	"metricParams": {
		"workloadName": "${workload-type}:${namespace}:${workload-name}"  # 示例deployment:test:httpbin
	},
	"interval": "5s",
	"isDetails": true,
	"from": "now-5m", # 查询时间范围起点
	"to": "now" # 查询时间范围终点
}
```



### headers

| key           | value                                                        |
| ------------- | ------------------------------------------------------------ |
| content-type  | application/json                                             |
| Authorization | Bearer token-5tgq4:xnq57jc7vgmjlcwzhtl8qdfbftt7tns7jf4bk6p4tstrh6bjx9whsd |

### response

```json
{
	"type": "collection",
	"data": [{
		"graphID": "p-5ggc4:workload-network-packet",
		"series": [{
			"name": "Receive errors(grafana-project-monitoring-6b69d54985-bktmt)",
			"points": [
				[0, 1576056851000],
				...
				[0, 1576057151000]
			]
		}, {
			"name": "Receive dropped(grafana-project-monitoring-6b69d54985-bktmt)",
			"points": [
				[0, 1576056851000],
				...
				[0, 1576057151000]
			]
		}, {
			"name": "Transmit errors(grafana-project-monitoring-6b69d54985-bktmt)",
			"points": [
				[0, 1576056851000],
				...
				[0, 1576057151000]
			]
		}, {
			"name": "Transmit dropped(grafana-project-monitoring-6b69d54985-bktmt)",
			"points": [
				[0, 1576056851000],
				...
				[0, 1576057151000]
			]
		}, {
			"name": "Receive packets(grafana-project-monitoring-6b69d54985-bktmt)",
			"points": [
				[5.100504919199245, 1576056851000],
				....
				[4.57875, 1576057146000],
				[4.57875, 1576057151000]
			]
		}, {
			"name": "Transmit packets(grafana-project-monitoring-6b69d54985-bktmt)",
			"points": [
				[5.0355085453127355, 1576056851000],
				...
				[4.57125, 1576057151000]
			]
		}]
	}, {
		"graphID": "p-5ggc4:workload-disk-io",
		"series": [{
			"name": "Write(grafana-project-monitoring-6b69d54985-bktmt)",
			"points": [
				[0, 1576056851000],
				...
				[0, 1576057151000]
			]
		}, {
			"name": "Read(grafana-project-monitoring-6b69d54985-bktmt)",
			"points": [
				[0, 1576056851000],
				...
				[0, 1576057151000]
			]
		}]
	}, {
		"graphID": "p-5ggc4:workload-network-io",
		"series": [{
			"name": "Receive(grafana-project-monitoring-6b69d54985-bktmt)",
			"points": [
				[477.5146754970512, 1576056851000],
				...
				[476.95875, 1576057151000]
			]
		}, {
			"name": "Transmit(grafana-project-monitoring-6b69d54985-bktmt)",
			"points": [
				[633.2357247437775, 1576056851000],
				...
				[567.34125, 1576057151000]
			]
		}]
	}, {
		"graphID": "p-5ggc4:workload-memory-usage-bytes-sum",
		"series": [{
			"name": "grafana-project-monitoring-6b69d54985-bktmt",
			"points": [
				[56700928, 1576056851000],
				...
				[57503744, 1576057151000]
			]
		}]
	}, {
		"graphID": "p-5ggc4:workload-cpu-usage",
		"series": [{
			"name": "CPU system seconds(grafana-project-monitoring-6b69d54985-bktmt)",
			"points": [
				[0.00013666523168160972, 1576056851000],
				...
				[0.0001195616451144745, 1576057151000]
			]
		}, {
			"name": "CPU cfs throttled(grafana-project-monitoring-6b69d54985-bktmt)",
			"points": [
				[0, 1576056851000],
				...
				[0, 1576057151000]
			]
		}, {
			"name": "CPU usage(grafana-project-monitoring-6b69d54985-bktmt)",
			"points": [
				[0.001228605559192756, 1576056851000],
				...
				[0.001172037997012214, 1576057151000]
			]
		}, {
			"name": "CPU user seconds(grafana-project-monitoring-6b69d54985-bktmt)",
			"points": [
				[0.0005124946188060365, 1576056851000],
				...
				[0.0005181004621626473, 1576057151000]
			]
		}]
	}]
}
```

# Pods

## 获取pods

### 接口

```
/v3/project/${projectid}/pods?limit=-1&sort=name
```

### 扩展

```shell
/v3/project/${projectid}/pods?limit=-1&sort=name&namespaceId=${namespace}&workloadId=${workloadtype}:${namespaceid}:${applicationname}-${compomentname}-workload-${version}
# 示例 /v3/project/c-x67ps:p-sjjrk/pods?limit=-1&sort=name&namespaceId=service&workloadId=deployment:service:fyhservice-fyhcompoment-workload-v1
```

根据namespace 和 名字 过滤出指pod

### verb

```
"Methods": [
        "GET"
  ]
```

### headers

| key           | value                                                        |
| ------------- | ------------------------------------------------------------ |
| content-type  | application/json                                             |
| Authorization | Bearer token-5tgq4:xnq57jc7vgmjlcwzhtl8qdfbftt7tns7jf4bk6p4tstrh6bjx9whsd |

### response

```json
{
    "type": "collection",
    "links": {
        "self": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/pods"
    },
    "createTypes": {
        "pod": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/pods"
    },
    "actions": {},
    "pagination": {
        "limit": 1000,
        "total": 1
    },
    "filters": {
        "name": null,
        "namespaceId": [
            {
                "modifier": "eq",
                "value": "service"
            }
        ],
        "workloadId": [
            {
                "modifier": "eq",
                "value": "deployment:service:fyhservice-fyhcompoment-workload"
            }
        ]
    },
    "resourceType": "pod",
    "data": [
        {
            "annotations": {
                "cni.projectcalico.org/podIP": "10.244.0.207/32",
            },
            "baseType": "pod",
            "containers": [
                {
                    "exitCode": null,
                    "image": "socp.io/library/httpbin:1031",
                    "imagePullPolicy": "IfNotPresent",
                    "initContainer": false,
                    "name": "fyhcontainer",
                    "ports": [
                        {
                            "containerPort": 8000,
                            "protocol": "UDP",
                            "sourcePort": 0,
                            "type": "/v3/project/schemas/containerPort"
                        }
                    ],
                    "resources": {
                        "limits": {
                            "cpu": "500m",
                            "memory": "200Mi"
                        },
                        "requests": {
                            "cpu": "500m",
                            "memory": "200Mi"
                        },
                        "type": "/v3/project/schemas/resourceRequirements"
                    },
                    "restartCount": 0,
                    "state": "running",
                    "stdin": false,
                    "stdinOnce": false,
                    "terminationMessagePath": "/dev/termination-log",
                    "terminationMessagePolicy": "File",
                    "transitioning": "no",
                    "transitioningMessage": "",
                    "tty": false,
                    "type": "/v3/project/schemas/container",
                    "volumeMounts": [
                        {
                            "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                            "name": "default-token-pmsdc",
                            "readOnly": true,
                            "type": "/v3/project/schemas/volumeMount"
                        }
                    ]
                },
               ...
            ],
            "created": "2019-12-20T06:59:51Z",
            "createdTS": 1576825191000,
            "creatorId": null,
            "dnsPolicy": "ClusterFirst",
            "hostIPC": false,
            "hostNetwork": false,
            "hostPID": false,
            "id": "service:fyhservice-fyhcompoment-workload-58cd864985-vvnfv",
            "labels": {
                "app": "fyhservice-fyhcompoment-workload",
                "pod-template-hash": "58cd864985",
                "version": ""
            },
            "links": {
                "remove": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/pods/service:fyhservice-fyhcompoment-workload-58cd864985-vvnfv",
                "self": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/pods/service:fyhservice-fyhcompoment-workload-58cd864985-vvnfv",
                "update": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/pods/service:fyhservice-fyhcompoment-workload-58cd864985-vvnfv",
                "yaml": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/pods/service:fyhservice-fyhcompoment-workload-58cd864985-vvnfv/yaml"
            },
            "name": "fyhservice-fyhcompoment-workload-58cd864985-vvnfv",
            "namespaceId": "service",
            "nodeId": "c-x67ps:machine-prxns",
            "ownerReferences": [
                {
                    "apiVersion": "apps/v1",
                    "blockOwnerDeletion": true,
                    "controller": true,
                    "kind": "ReplicaSet",
                    "name": "fyhservice-fyhcompoment-workload-58cd864985",
                    "type": "/v3/project/schemas/ownerReference",
                    "uid": "51b25eef-22f6-11ea-8e5d-fa163ecee4c9"
                }
            ],
            "priority": 0,
            "projectId": "c-x67ps:p-sjjrk",
            "restartPolicy": "Always",
            "schedulerName": "default-scheduler",
            "scheduling": {
                "tolerate": [
                    {
                        "effect": "NoExecute",
                        "key": "node.kubernetes.io/not-ready",
                        "operator": "Exists",
                        "tolerationSeconds": 300,
                        "type": "/v3/project/schemas/toleration"
                    },
                    {
                        "effect": "NoExecute",
                        "key": "node.kubernetes.io/unreachable",
                        "operator": "Exists",
                        "tolerationSeconds": 300,
                        "type": "/v3/project/schemas/toleration"
                    }
                ]
            },
            "serviceAccountName": "default",
            "state": "running",
            "status": {
   				...
                "nodeIp": "10.10.111.54",
                "phase": "Running",
                "podIp": "10.244.0.207",
                "qosClass": "Burstable",
                "startTime": "2019-12-20T06:59:52Z",
                "startTimeTS": 1576825192000,
                "type": "/v3/project/schemas/podStatus"
            },
            "terminationGracePeriodSeconds": 30,
            "transitioning": "no",
            "transitioningMessage": "",
            "type": "pod",
            "uuid": "51d20a5c-22f6-11ea-8e5d-fa163ecee4c9",
            "volumes": [
                {
                    "name": "default-token-pmsdc",
                    "secret": {
                        "defaultMode": 420,
                        "secretName": "default-token-pmsdc",
                        "type": "/v3/project/schemas/secretVolumeSource"
                    },
                    "type": "/v3/project/schemas/volume"
                },
                ...
            ],
            "workloadId": "deployment:service:fyhservice-fyhcompoment-workload"
        }
    ]
}
```

## 删除pod

### 接口

```
/v3/project/${projectid}/pods/${namespace}:${podname}
```

### verb

```
"Methods": [
        "DELETE"
  ]
```



## 查询Pod负载

### 接口

```
/v3/projectmonitorgraphs?action=query
```

### verb

```html
"Methods": [
        "POST"
  ]
```

### headers

| key           | value                                                        |
| ------------- | ------------------------------------------------------------ |
| content-type  | application/json                                             |
| Authorization | Bearer token-5tgq4:xnq57jc7vgmjlcwzhtl8qdfbftt7tns7jf4bk6p4tstrh6bjx9whsd |

### payload

```json
{
	"filters": {
		"resourceType": "pod",
		"projectId": "${projectid}"
	},
	"metricParams": {
		"podName": "${namespace}:${podname}"
	},
	"interval": "5s",
	"isDetails": true,
	"from": "now-5m",
	"to": "now"
}


### 示例
{
	"filters": {
		"resourceType": "pod",
		"projectId": "c-gwnnw:p-hmmh9"
	},
	"metricParams": {
		"podName": "ingress-nginx:nginx-ingress-controller-qmnf9"
	},
	"interval": "5s",
	"isDetails": true,
	"from": "now-5m",
	"to": "now"
}
```

### responese

```json
{
	"type": "collection",
	"data": [{
		"graphID": "p-hmmh9:pod-network-packet",
		"series": [{
			"name": "Receive errors",
			"points": [
				[0, 1576834551000],
				...
				[0, 1576834851000]
			]
		}, {
			"name": "Receive dropped",
			"points": [
				[0, 1576834551000],
				...
				[0, 1576834851000]
			]
		}, {
			"name": "Transmit errors",
			"points": [
				[0, 1576834551000],
				...
				[0, 1576834851000]
			]
		}, {
			"name": "Transmit packets",
			"points": [
				[3436.049359075862, 1576834551000],
				...
				[3210.96378671229, 1576834851000]
			]
		}, {
			"name": "Transmit dropped",
			"points": [
				[0, 1576834551000],
				...
				[0, 1576834851000]
			]
		}, {
			"name": "Receive packets",
			"points": [
				[1471.4889104105523, 1576834551000],
				...
				[1367.739520958084, 1576834851000]
			]
		}]
	}, {
		"graphID": "p-hmmh9:pod-memory-usage-bytes-sum",
		"series": [{
			"name": "Memory usage(nginx-ingress-controller)",
			"points": [
				[234831872, 1576834551000],
				...
				[234913792, 1576834851000]
			]
		}]
	}, {
		"graphID": "p-hmmh9:pod-disk-io",
		"series": [{
			"name": "Write(nginx-ingress-controller)",
			"points": [
				[0, 1576834551000],
				...
				[0, 1576834851000]
			]
		}, {
			"name": "Read(nginx-ingress-controller)",
			"points": [
				[0, 1576834551000],
				...
				[0, 1576834851000]
			]
		}]
	}, {
		"graphID": "p-hmmh9:pod-cpu-usage",
		"series": [{
			"name": "CPU system seconds(nginx-ingress-controller)",
			"points": [
				[0.012992010485047202, 1576834551000],
				...
				[0.010964132692197434, 1576834851000]
			]
		}, {
			"name": "CPU usage(nginx-ingress-controller)",
			"points": [
				[0.03289497521611722, 1576834551000],
				...
				[0.0320943907946117, 1576834851000]
			]
		}, {
			"name": "CPU user seconds(nginx-ingress-controller)",
			"points": [
				[0.018249774259058677, 1576834551000],
				...
				[0.017213314124274005, 1576834851000]
			]
		}]
	}, {
		"graphID": "p-hmmh9:pod-network-io",
		"series": [{
			"name": "Transmit",
			"points": [
				[8049937.402646183, 1576834551000],
				...
				[7569497.943398917, 1576834851000]
			]
		}, {
			"name": "Receive",
			"points": [
				[4295550.995908838, 1576834551000],
				...
				[4055065.00570288, 1576834851000]
			]
		}]
	}]
}
```

# Application

## 对接开发功能

1. 服务容器配置
   1. 容器配置
      - 容器基本配置（components.[].containers[](name image imagePullPolicy Command args)）√
      - 容器副本数配置 （**components.[].componentTraits.replicas** ~~components.[]optTraits.manualScaler.replicas~~）√ 
      - 镜像拉取秘钥配置（**optTraits.imagePullConfig** ~~components.[].devTraits.imagePullConfig~~）√ (暂时不需要开发)
      - 计算资源配置以及本地目录挂载(components.[].containers[].resources) √
      - configmap 配置 (components.[].containers[].config) √
      - 环境变量配置 （components.[].containers[].env）√
      - 容器健康状态检查（components.[].containers[].livenessProbe）√
      - 容器服务状态检查(components.[].containers[].readinessProbe) √
      - 容器退出预处理（components.[].containers[].lifecycle.preStop）√
      - 容器自定义退出时间（**components.[].componentTraits.terminationGracePeriodSeconds** ~~components.[].optTraits.terminationGracePeriodSeconds~~）
      - 容器调度策略（**components.[].componentTraits.schedulePolicy**  ~~components.[].optTraits.schedulePolicy~~）√ 
      - 容器自定义metric配置（**components.[].componentTraits.custommetric**）√ 
      - 容器自定义日志采集配置(**components.[].componentTraits.logcollect**) √ 
2. 服务治理配置
   - 负载均衡策略（**optTraits.loadBalancer**  ~~components.[].devTraits.ingressLB~~）√ 
   - 路由规则配置 （**optTraits.ingress** ~~components.[].optTraits.ingress~~）√
   - 手动熔断配置 （**optTraits.fusing ** ~~components.[].optTraits.fusing~~）√
   - 自动熔断配置 （**optTraits.circuitbreaking ** ~~components.[].optTraits.circuitbreaking~~）√
   - 限流配置 （**optTraits.rateLimit** ~~components.[].optTraits.rateLimit~~）√
   - 访问控制白名单配置 **optTraits.whiteList** ~~components.[].optTraits.whiteList~~）√
   - httpRetry 连接重试逻辑（**optTraits.httpretry** ~~components.[].optTraits.httpretry~~）√
   - 自动扩缩逻辑 （**components.[].componentTraits.autoscaling** ~~components.[].optTraits.autoscaling~~）√
   - 灰度发布（**optTraits.grayRelease**）√ 
3. changelog
   - 上述字段 容器自定义metric配置 容器自定义日志采集配置 灰度发布为新增字段 其余加粗字段为修改字段

## 数据定义

```json
{
	"name": "String", // 必选 application Name
	"namespaceId": "string", // 必选 应用部署目标Namespace
	"labels": "map[string]string"， //必选 必须添加此格式Label "projectId": "c-x67ps_p-sjjrk"，applicationTemplateId(template名称)
	"annotations": "map[string]string", //可选 可作为后期功能扩展字段
	"components": [{
		"name": "string", // 必选 组件名
		"workloadType": "string", // 服务类型 目前仅支持"Server"
		"workloadSettings": [{
			"name": "string",
			"type": "string",
			"value": "string",
			"fromparam": "string"
		}], // 扩展字段 可选
		"version": "string", // 服务版本 必选
		"parameters": [{
			"name": "string",
			"description": "string",
			"type": "string",
			"required": bool,
			"default": "string"
		}], // 可选
		"containers": [{
				"name": "string", // 必选 容器名
				"command": "[]string"， //可选 命令
				"args": "[]string", //可选 参数
				"config": [{
					"path": "string", //(必选）挂载到容器内路径
					"fileName": "string" //（必选）挂载到容器内的文件名
					"value": "string" //（必选）文件内容
				}], //可选 通过该项为容器创建configmap资源并挂载 
				"env": [{
					"fromParam": "string", //可选 见示例2 目前只支持(spec.nodeName 对应所在主机名 metadata.name 获取容器名 metadata.namespace获取所在namespace status.podIP获取容器ip)
					"name": "string", //必选
					"value": "string" //（必选）说明 如果fromparam不为空 value不需要再填 
				}], //可选 用于配置环境变量
				"image": "string", // 必选项
				"imagePullPolicy": "string", // 可选项 镜像拉取策略 默认为Always 可填字段仅限于Always，IfNotPresent，Never。
				"livenessProbe": {
					"exec": {
						"command": “[] string ", // 可选项 在容器内执行指定命令。如果命令退出时返回码为 0 则表明容器健康
					},
					"httpGet": {
						"port": "int", // (必选)
						"path": "string", // (必选)
						"httpHeaders": [{
							"name": "string", // (必选)
							"value": "string" // (必选)
						}] // 可选 配置请求头
					}, // 可选 通过httpGet判断容器内服务健康状态
					"tcpSocket": {
						"port": "int"， // 可选 配置监听容器内端口健康状态
					}, // ！！！ exec httpGet tcpSocket 三项只能选其中一项
					"initialDelaySeconds": "int"， // 容器启动和探针启动之间的秒数
					"periodSeconds": "int", //检查的频率（以秒为单位）。默认为10秒。最小值为1。
					"timeoutSeconds": "int" // 配置检查超时时间
					"successThreshold": "int" // 查成功的最小连续成功次数。默认为1.活跃度必须为1。最小值为1
					"failureThreshold": "int" // 当Pod成功启动且检查失败时，Kubernetes将在放弃之前尝试failureThreshold次。放弃生存检查意味着重新启动Pod。而放弃就绪检查，Pod将被标记为未就绪。默认为3.最小值为1。
				}, // 可选项 判断容器是否存活策略配置 见示例3
				"ports": [{
					"containerPort": "int", //可选 容器内服务监听端口
					"name": "string", //可选
					"protocol": "string" //可选
				}], // 可选
				"readinessProbe": {
					//内容于livenessProbe一致
				}, // 可选项 用于判断容器服务状态如果异常则从service endpoint列表移除
				"lifecycle": {
					"postStart": {
						// 同Prestop相同 目前不需要
					},
					"preStop": {
						"exec": {
							"command": "[]string", // 可选项 在容器内执行指定命令或脚本，做容器退出前的清理工作。
						},
						"httpGet": {
							"port": "int", // (必选)
							"path": "string", // (必选)
							"httpHeaders": [{
								"name": "string", // (必选)
								"value": "string" // (必选)
							}] // 可选 配置请求头
						}, // 可选 通过httpGet判断容器内服务健康状态
						"tcpSocket": {
							"port": "int"， // 可选 配置监听容器内端口健康状态
						} //// ！！！ exec httpGet tcpSocket 三项只能选其中一项
					}
				},
				"resources": {
					"cpu": "string", // 可选 cpu资源配额 单位m 1000m等价于1核cpu
					"gpu": "int", // 可选 gpu资源配额
					"memory": "string" // 可选 内存资源配额 单位Mi,Gi
					"volumes": [{
						"name": "string", // (必选)
						"mountPath": "string", // (必选)
						"accessMode": "string", // 可选
						"sharingPolicy": "string", // 可选，
						"disk": {
							"required": "string", //（如果ephemeral 为false 会给容器创建hostpath挂载卷 则此项必选 需要填写物理机上对应挂载目录）
							"ephemeral": "bool" // 是否需要持久化卷 false 对应创建hostpath true 对应创建emptydir 
						} // 这个数据结构对应容器挂载本地卷 
					}] // 可选
				}, // 可选
				"securityContext": {
					"RunAsNonRoot": bool // TODO
				}, // 可选 容器权限配置 TODO
			} // (必选) 容器配置
		], //可选  (平台托管必选 非平台托管此配置为空）
		"componentTraits": {
			"replicas": "int",
			"custommetric": {
				"enable": "bool", // 服务是否上报自定义指标
				"uri": "string" // 服务指标查询接口 示例/jaminfo 如果enable 为true 则需要填写uri
			}, // 自定义指标配置
			"logcollect": "bool", //是否采集规定目录用户自定义日志
			"terminationGracePeriodSeconds": "int", // 可选项 配置容器内进程完全退出所需处理时间
			"schedulePolicy": {
				"nodeSelector": "map[string]string", //根据一定的标签调度Pod到指定node
				"nodeAffinity": {
					"hardAffinity": "bool", // 硬限制（true） or 软限制(false)
					"labelSelectorRequirement": {
						"key": "string", // key (必选)
						"operator": "string", // 操作符（必选） In：label 的值在某个列表中 NotIn：label 的值不在某个列表中 Exists：某个 label 存在 DoesNotExist：某个 label 不存在（只能填这四种操作符 如果操作符为Exists或DoesNotExist 则不需用填写values）
						"values": "[]string" // （可选）
					}
				}, // Node亲和性规则（可选）
				"podAffinity": {
					"hardAffinity": "bool", // 硬限制（true） or 软限制(false)
					"labelSelectorRequirement": {
						"key": "string", // key (必选)
						"operator": "string", // 操作符（必选） In：label 的值在某个列表中 NotIn：label 的值不在某个列表中 Exists：某个 label 存在 DoesNotExist：某个 label 不存在（只能填这四种操作符 如果操作符为Exists或DoesNotExist 则不需用填写values）
						"values": "[]string" // （可选）
					}
				}, // Pod亲和性规则
				"podAntiAffinity": {
					"hardAffinity": "bool", // 硬限制（true） or 软限制(false)
					"labelSelectorRequirement": {
						"key": "string", // key (必选)
						"operator": "string", // 操作符（必选） In：label 的值在某个列表中 NotIn：label 的值不在某个列表中 Exists：某个 label 存在 DoesNotExist：某个 label 不存在（只能填这四种操作符 如果操作符为Exists或DoesNotExist 则不需用填写values）
						"values": "[]string" // （可选）
					}
				}, // Pod反亲和性规则
			}, // 容器调度策略(可选)
			"autoscaling": {
				"metric": "string", // 伸缩容依赖的指标
				"threshold": "float64", // 阈值
				"maxreplicas": "int", // 最大扩容副本数
				"minreplicas": "int" // 最小缩容副本数
			} // 自动扩缩配置
		}
	}], // 组件列表 必选
	"optTraits": {
		"imagePullConfig": {        
			"registry": "string", //（必选）
			        "username": "string", // (必选)
			        "password": "string" // (必选)
			     
		}, // 可选 配置镜像库config 暂时不需要开发相关功能
		"httpretry": {
			"attempts": "int" // 重试次数，
			"pertrytimeout": "string" // 重试时间间隔 示例3s
		}, // 请求重试配置
		"ingress": {
			"host": "string", // (必选) 访问入口域名
			"path": "[]string", // 可选	访问路径 默认为 ["/"]
			"serverPort": "int" //(必选) 服务端口
		}, // 必选 对外提供访问配置 ！！校验
		"eject": "[]string", // 保留字段 暂未开发相关功能
		"fusing": {
			"podlist": "[]string", // podname列表
			"action": "string" // in 或 out
		}, // 可选 手动熔断配置
		"rateLimit": {
			"timeDuration": "string", //(必选)
			"requestAmount": "int" // (必选),
			"overrides": [{
				"user": "string", //(必选)
				"requestAmount": "int", // (必选)
			}], // 可选覆盖默认参数
		}, // 可选 限流配置
		"loadBalancer": {
			"simple": "string", //负载均衡策略 "ROUND_ROBIN" or "LEAST_CONN" or "RANDOM"
			"consistentHash": {
				"useSourceIp": "bool" // 根据请求源IP绑定
			},
		}, // 负载均衡策略
		"grayRelease": "map[string]int", //灰度发布
		"circuitbreaking": {
			"connectionPool": {
				"tcp": {
					"maxConnections": "int", //最大TCP连接数（必选）
					"connectTimeout": "string", // 超时时间（可选）
				},
				"http": {
					"http1MaxPendingRequests": "int", //最大允许HTTP1等待连接数
					"http2MaxRequests": "int", // 最大允许HTTP2请求数
					"maxRequestsPerConnection": "int", // 最大允许HTTP2每个请求中的连接数
					"maxRetries": "int", // 最大重试次数
				},
			}, // 服务配置连接池
			"outlierDetection": {
				"consecutiveErrors": "int", // 错误连接数阈值
				"interval": "string", // 检查时间间隔
				"baseEjectionTime": "string", // 拒绝访问时间
				"maxEjectionPercent": "int" // 拒绝流量百分比
			}
		} // 可选 自动熔断配置
		"whiteList": {
			"users": "[]string"
		} // 可选 服务访问者白名单
	}, // 必选 运维人员配置
	"status": {
		"map[string]object{
		"componentId": "string",
		"workload": "string",
		"service": "string",
		"configMaps": "[]string",
		"imagePullSecret": "string",
		"gateway": "string",
		"policy": "string",
		"clusterRbacConfig": "string",
		"virtualService": "string",
		"serviceRole": "string",
		"serviceRoleBinding": "string",
		"DestinationRule": "string"
	} //保留字段
}

---
示例1
      env:
        - name: APPLOGLEVEL
          valueFrom:
            configMapKeyRef:
              name: cm-app
              key: apploglevel

----
示例2 

      env:
        - name: MY_NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        - name: MY_POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: MY_POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: MY_POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
---
示例3 

    livenessProbe:
      exec:
        command:
        - cat
        - /tmp/health  
      tcpSocket:
        port: 80
      httpGet:
        path: /_status/healthz
        port: 80
      initialDelaySeconds: 15
      timeoutSeconds: 1
```



## 接口

```
/v3/projects/clusterid:projectid/application
```

## verb

```
"Methods": [
        "GET",
        "POST",
        "DELETE"
  ]
```

## headers

| key           | value                                                        |
| ------------- | ------------------------------------------------------------ |
| content-type  | application/json                                             |
| Authorization | Bearer token-5tgq4:xnq57jc7vgmjlcwzhtl8qdfbftt7tns7jf4bk6p4tstrh6bjx9whsd |

## 样例

### GET

#### url

```
/v3/projects/${projectid}/applications/${namespace}:${applicationname}
```

#### verb

```
"Methods": [
        "GET"
  ]
```

#### headers

| key           | value                                                        |
| ------------- | ------------------------------------------------------------ |
| Authorization | Bearer token-5tgq4:xnq57jc7vgmjlcwzhtl8qdfbftt7tns7jf4bk6p4tstrh6bjx9whsd |

#### response

```json
{
    "annotations": {},
    "baseType": "application",
    "components": [
        {
            "containers": [
                {
                    "image": "socp.io/library/httpbin:1031",
                    "imagePullPolicy": "IfNotPresent",
                    "name": "fyhfyh",
                    "ports": [
                        {
                            "containerPort": 8000,
                        }
                    ],
                }
            ],
            "devTraits": {
                "ingressLB": {
                    "consistentType": "sourceIP",
                },
                "staticIP": false,
            },
            "name": "fyh",
            "optTraits": {
                "ingress": {
                    "host": "fyh.com",
                    "path": "/",
                    "serverPort": 8080,
                    "type": "/v3/project/schemas/appIngress"
                },
                "manualScaler": {
                    "replicas": 1,
                    "type": "/v3/project/schemas/manualScaler"
                },
                "type": "/v3/project/schemas/componentTraitsForOpt",
                "whiteList": {
                    "type": "/v3/project/schemas/whiteList",
                    "users": [
                        ""
                    ]
                }
            },
            "type": "/v3/project/schemas/component",
            "workloadType": "Server"
        }
    ],
    "created": "2019-12-18T09:12:15Z",
    "createdTS": 1576660335000,
    "creatorId": "user-5js2x",
    "id": "service:fyhtest",
    "labels": {
        "cattle.io/creator": "norman"
    },
    "name": "fyhtest",
    "namespaceId": "service",
    "projectId": "c-x67ps:p-sjjrk",
    "state": "active",
    "type": "application",
    ...
}
```



### Create

#### url

```
/v3/projects/${projectid}/application?_replace=true
```

#### verb

```
"Methods": [
        "POST"
  ]
```

#### headers

| key           | value                                                        |
| ------------- | ------------------------------------------------------------ |
| content-type  | application/json                                             |
| Authorization | Bearer token-5tgq4:xnq57jc7vgmjlcwzhtl8qdfbftt7tns7jf4bk6p4tstrh6bjx9whsd |

#### body

```json
{
	"name": "zk-test1",
	"namespaceId": "service",
	"version": "v1",
	"components": [{
		"name": "zkhttpbin",
		"workloadType": "Server",
		"containers": [{
			"name": "zkhttpbin",
			"image": "docker.io/citizenstig/httpbin@sha256:b81c818ccb8668575eb3771de2f72f8a5530b515365842ad374db76ad8bcf875"
		}],
		"optTraits": {
			"ingress": {
				"host": "www.zkhttpbin.com",
				"serverPort": 8000
			},
			"manualScaler": {
				"replicas": 1
			},
			"whiteList": {
				"users": ["testsystem@keycloak.com"]
			}
		}
	}]
}     #详细数据结构见数据定义
```

#### response

```json
{
    "annotations": {},
    "baseType": "application",
    "components": [
        {
            "containers": [
                {
                    "image": "nginx:latest",
                    "initContainer": false,
                    "name": "test",
                    "restartCount": 0,
                    "stdin": false,
                    "stdinOnce": false,
                    "tty": false,
                    "type": "/v3/project/schemas/container"
                }
            ],
            "type": "/v3/project/schemas/component",
            "workloadType": "Server"
        }
    ],
    "created": "2019-11-27T07:42:16Z",
    "createdTS": 1574840536000,
    "creatorId": "user-5js2x",
    "id": "istio-test:application-test",
    "labels": {
        "cattle.io/creator": "norman"
    },
    "links": {
        "remove": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/applications/istio-test:application-test",
        "self": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/applications/istio-test:application-test",
        "update": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/applications/istio-test:application-test",
        "yaml": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/applications/istio-test:application-test/yaml"
    },
    "name": "application-test",
    "namespaceId": "istio-test",
    "projectId": "c-x67ps:p-sjjrk",
    "state": "active",
    "transitioning": "no",
    "transitioningMessage": "",
    "type": "application",
    "uuid": "6ee7eb00-10e9-11ea-af5e-fa163ecee4c9"
}
```

### Update

#### url

```
/v3/projects/${projectid}/applications/${namespace}:${applicationname}
# 示例 https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/applications/service:zk-test1
```

#### verb

```
"Methods": [
        "PUT"
  ]
```

#### headers

| key           | value                                                        |
| ------------- | ------------------------------------------------------------ |
| content-type  | application/json                                             |
| Authorization | Bearer token-5tgq4:xnq57jc7vgmjlcwzhtl8qdfbftt7tns7jf4bk6p4tstrh6bjx9whsd |

#### body

```json
{
	"name": "zk-test1",
	"namespaceId": "service",
	"version": "v1",
	"components": [{
		"name": "zkhttpbin",
		"workloadType": "Server",
		"containers": [{
			"name": "zkhttpbin",
			"image": "socp.io/library/httpbin:1031"
		}],
		"optTraits": {
			"ingress": {
				"host": "www.zkhttpbin.com",
				"serverPort": 8000
			},
			"manualScaler": {
				"replicas": 2
			},
			"whiteList": {
				"users": ["testsystem@keycloak.com", "test2@qq.com"]
			}
		}
	}]
}      #详细数据结构见数据定义
```

#### response

```json
{
	"annotations": {},
	"baseType": "application",
	"components": [{
		"containers": [{
			"image": "socp.io/library/httpbin:1031",
			"name": "zkhttpbin",
			"type": "/v3/project/schemas/componentContainer"
		}],
		"name": "zkhttpbin",
		"optTraits": {
			"ingress": {
				"host": "www.zkhttpbin.com",
				"serverPort": 8000,
				"type": "/v3/project/schemas/appIngress"
			},
			"manualScaler": {
				"replicas": 2,
				"type": "/v3/project/schemas/manualScaler"
			},
			"type": "/v3/project/schemas/componentTraitsForOpt",
			"whiteList": {
				"type": "/v3/project/schemas/whiteList",
				"users": ["testsystem@keycloak.com", "test2@qq.com"]
			}
		},
		"type": "/v3/project/schemas/component",
		"workloadType": "Server"
	}],
	"created": "2019-12-18T03:06:07Z",
	"createdTS": 1576638367000,
	"creatorId": "user-5js2x",
	"id": "service:zk-test1",
	"labels": {
		"cattle.io/creator": "norman"
	},
	"links": {
		"remove": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/applications/service:zk-test1",
		"self": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/applications/service:zk-test1",
		"update": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/applications/service:zk-test1",
		"yaml": "https://10.10.111.54:444/v3/project/c-x67ps:p-sjjrk/applications/service:zk-test1/yaml"
	},
	"name": "zk-test1",
	"namespaceId": "service",
	"projectId": "c-x67ps:p-sjjrk",
	"state": "active",
	"transitioning": "no",
	"transitioningMessage": "",
	"type": "application",
	"uuid": "561c05a8-2143-11ea-a471-fa163ecee4c9"
}
```

### Delete

#### url

```
/v3/projects/${projectid}/application/${namespace}:${applicationname}
示例 /v3/projects/c-x67ps:p-sjjrk/application/istio-
test:application-test
```

#### headers

同上

#### verb

```
"Methods": [
        "DELETE"
  ]
```

## 完整JSON样例(不断更新中)

```json
{
	"name": "zk-1",
	"namespaceId": "test-ns",
	"labels": {
		"projectId": "c-x67ps_p-sjjrk",
		"applicationTemplateId": "zk-1"
	},
	"components": [{
		"name": "zk-1",
		"version": "v1",
		"workloadType": "Server",
		"containers": [{
			"ports": [{
				"containerPort": 8000
			}],
			"env": [{
				"name": "TEST",
				"value": "test"
			}, {
				"name": "PODIP",
				"fromParam": "status.podIP"
			}],
			"name": "zk-1",
			"imagePullPolicy": "IfNotPresent",
			"image": "socp.io/library/httpbin:1031",
			"resources": {
				"cpu": "300m",
				"memory": "300Mi",
				"volumes": [{
					"name": "test",
					"mountPath": "/mnt/test",
					"disk": {
						"ephemeral": false,
						"required": "/home/zk/test"
					}
				}, {
					"name": "test2",
					"mountPath": "/mnt/test2",
					"disk": {
						"ephemeral": true
					}
				}]
			},
			"config": [{
					"path": "/etc/test",
					"fileName": "test.yaml",
					"value": "wfwdaawfwadwawafaefaefalakklwak"
				},
				{
					"path": "/etc/test",
					"fileName": "test1.yaml",
					"value": "rrrrrlakklwak"
				}
			],
			"readinessProbe": {
				"exec": {
					"command": ["cat", "/etc/test/test.yaml"]
				},
				"initialDelaySeconds": 5,
				"periodSeconds": 5
			},
			"lifecycle": {
				"postStart": {
					"exec": {
						"command": ["cat", "/etc/test/test.yaml"]
					}
				},
				"prestop": {
					"exec": {
						"command": ["cat", "/etc/test/test.yaml"]
					}
				}
			},
			"livenessProbe": {
				"exec": {
					"command": ["cat", "/etc/test/test.yaml"]
				},
				"initialDelaySeconds": 5,
				"periodSeconds": 5
			}
		}],
		"devTraits": {
			"imagePullConfig": {
				"registry": "socp.io",
				"username": "zk",
				"password": "zk123456"
			}
		},
		"optTraits": {
			"ingress": {
				"host": "ingress.socp.io",
				"path": "/headers",
				"serverPort": 8000
			},
			"rateLimit": {
				"timeDuration": "1m",
				"requestAmount": 3,
				"overrides": [{
					"user": "zk-1@qq.com",
					"requestAmount": 8
				}]
			},
			"fusing": {
				"action": "out",
				"podlist": []
			},
			"manualScaler": {
				"replicas": 1
			},
			"circuitbreaking": {
				"connectionPool": {
					"tcp": {
						"maxConnections": 1
					},
					"http": {
						"http1MaxPendingRequests": 1,
						"maxRequestsPerConnection": 1
					}
				},
				"outlierDetection": {
					"consecutiveErrors": 2,
					"interval": "2s",
					"baseEjectionTime": "3m",
					"maxEjectionPercent": 100
				}
			},
			"whiteList": {
				"users": ["zk-1@qq.com"]
			}
		}
	}]
}
```



