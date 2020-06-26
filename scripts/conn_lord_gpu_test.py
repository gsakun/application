import requests
import json
import urllib3
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

'''
{
                  "name": "lbj-3",
                  "namespaceId": "istio-test",
                  #"version": "v1",
                  "components": [{
                      "name": "lbj-3",
                      "workloadType": "Server",
                      "containers": [{
                          "ports": [
                              {
                                  "containerPort": 8000
                              }
                          ],
                          "name": "lbj-3",
                          "imagePullPolicy": "IfNotPresent",
                          "image": "socp.io/library/hjw-ner:v0"
                      }],
                      "optTraits": {
                              "ingress": {
                                  "host": "www.lbj-3.com",
                                  "path": "/",
                                  "serverPort": 8000
                              },
                              "manualScaler": {
                                  "replicas": 1
                              },
                              "whiteList":{
                                  "users": [""]
                              }
                      },
                      "devTraits": {
                          "ingressLB": {
                              "consistentType": "sourceIP"
                          },
                          "staticIP": False
                      }
                    }]
              }
'''
url = "https://10.10.111.45:31888"
#url = "https://10.10.111.54:444"
#url = "https://183.131.12.19:32067"
headers = {
    'Content-Type': "application/json"
}

class User:
    token = ""
    def __init__(self,username,password):
        self.username = username
        self.password = password
        loginurl = url + "/v3-public/localProviders/local?action=login"
        body = {
            "username": self.username,
            "password": self.password
        }
        response = requests.post(loginurl, json=body, headers=headers, verify=False)
        reponseinfo = json.loads(response.text)
        self.token = reponseinfo["token"]
    def getclusters(self):
        clustersurl = url + "/v3/clusters?limit=-1&sort=name"
        #token = self.login()
        getheaders = {
            'Content-Type': "application/json",
            'Authorization': "Bearer" + " " + self.token
        }
        clusterlist = []
        response = requests.get(clustersurl, headers=getheaders, verify=False)
        reponseinfo = json.loads(response.text)
        clusterinfo = reponseinfo["data"]
        for cluster in clusterinfo:
            print ('{} --- {}'.format(cluster["name"],cluster["id"]))
            clusterlist.append(cluster["id"])
        return clusterlist
    def getproject(self,clusterid):
        projecturl = url + "/v3/clusters/" + clusterid + "/" + "projects"
        getheaders = {
            'Content-Type': "application/json",
            'Authorization': "Bearer" + " " + self.token
        }
        response = requests.get(projecturl, headers=getheaders, verify=False)
        reponseinfo = json.loads(response.text)
        projectinfo = reponseinfo["data"]
        for project in projectinfo:
            print ('{} --- {}'.format(project["name"],project["id"]))
    def createapplication(self):
        createapplicationurl = url + "/v3/projects/local:p-rnqx6/application?_replace=true"
        getheaders = {
            'Content-Type': "application/json",
            'Authorization': "Bearer" + " " + self.token
        }

        body = {
	"name": "gpu-test",
	"namespaceId": "service",
	"labels": {
		"projectId": "local_p-rnqx6",
		"applicationTemplateId": "gpu-test"
	},
	"components": [{
		"name": "gpu-test",
		"version": "v1",
		"workloadType": "Server",
		"containers": [{
			"ports": [{
				"containerPort": 8888
			}],
			"env": [{
				"name": "PATH",
				"value": "/root/miniconda3/bin:/root/miniconda3/condabin:/usr/local/nvidia/bin:/usr/local/cuda/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
			}],
			"name": "gpu-test",
			"imagePullPolicy": "IfNotPresent",
			"image": "socp.io/library/copysearch-gpu:v1.1",
			"command": [
				"python",
				"/service/server.py"
			],
			"resources": {
				"cpu": "2000m",
				"memory": "2048Mi",
				"volumes": [{
					"name": "vol1",
					"mountPath": "/service",
					"disk": {
						"ephemeral": False,
						"required": "/home/imageSearch"
					}
				}]
			}
		}],
		"optTraits": {
			"custommetric": {
				"enable": True,
				"uri": "http://127.0.0.1:8888/jaminfo"
			},
			"autoscaling": {
				"metric": "avg---data---all---1m",
				"threshold": 100,
				"maxreplicas": 3,
				"minreplicas": 1
			},
			"httpretry": {
				"attempts": 3,
				"perTryTimeout": "5s"
			},
			"schedulePolicy": {
				"nodeAffinity": {
					"hardAffinity": True,
					"labelSelectorRequirement": {
						"key": "kubernetes.io/hostname",
						"operator": "In",
						"values": [
							"k8s-5.novalocal"
						]
					}
				}
			},
			"ingress": {
				"host": "gputest.socp.io",
				"path": "/",
				"serverPort": 8888
			},
			"manualScaler": {
				"replicas": 1
			},
			"whiteList": {
				"users": ["zk-1@qq.com"]
			}
		}
	}]
}
        response = requests.post(createapplicationurl, json=body, headers=getheaders,verify=False)
        print (response.text)
        print (response.status_code)
    def updateapplication(self,namespaceid,name):
        createapplicationurl = url + "/v3/projects/c-bnjtk:p-f89vh/applications/" + namespaceid + ":" + name
        getheaders = {
            'Content-Type': "application/json",
            'Authorization': "Bearer" + " " + self.token
        }

        body = {
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
                            "type": "/v3/project/schemas/appPort"
                        }
                    ],
                    "type": "/v3/project/schemas/componentContainer"
                }
            ],
            "devTraits": {
                "ingressLB": {
                    "consistentType": "sourceIP",
                    "type": "/v3/project/schemas/ingressLB"
                },
                "staticIP": False,
                "type": "/v3/project/schemas/componentTraitsForDev"
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
                        "dwda@wqeqw.com"
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
    "links": {
        "remove": "https://10.10.111.54:444/v3/project/c-bnjtk:p-f89vh/applications/service:fyhtest",
        "self": "https://10.10.111.54:444/v3/project/c-bnjtk:p-f89vh/applications/service:fyhtest",
        "update": "https://10.10.111.54:444/v3/project/c-bnjtk:p-f89vh/applications/service:fyhtest",
        "yaml": "https://10.10.111.54:444/v3/project/c-bnjtk:p-f89vh/applications/service:fyhtest/yaml"
    },
    "name": "fyhtest",
    "namespaceId": "service",
    "projectId": "c-bnjtk:p-f89vh",
    "state": "active",
    "transitioning": "no",
    "transitioningMessage": "",
    "type": "application",
    "uuid": "7bb975e2-2176-11ea-a471-fa163ecee4c9"
}
        response = requests.put(createapplicationurl, json=body, headers=getheaders,verify=False)
        print (response.text)
        print (response.status_code)
    def createtemplate(self):
        createurl = url + "/v3/applicationConfigurationTemplate?_replace=true"
        getheaders = {
            'Content-Type': "application/json",
            'Authorization': "Bearer" + " " + self.token
        }
        body = {
            "components": [{
                "name": "TEST",
                "workloadType": "Server",
                "containers": [{
                    "name": "test1204",
                    "image": "busybox:1.28.3"
                }],
            }],
            "labels": {},
            "name": "test-application",
            "namespaceId": "istio-test"
        }
        response = requests.post(createurl, json=body, headers=getheaders,verify=False)
        print (response.text)
        print (response.status_code)
    def deletetemplate(self,name):
        deleteurl = url + "/v3/applicationConfigurationTemplate/" + name
        getheaders = {
            'Content-Type': "application/json",
            'Authorization': "Bearer" + " " + self.token
        }
        response = requests.delete(deleteurl,headers=getheaders,verify=False)
        print (response.text)
        print (response.status_code)
    def deleteapplication(self,name,id,namespace):
        deleteurl = url + "/v3/projects/" + id + "/application/" + namespace + ":" + name 
        getheaders = {
            'Content-Type': "application/json",
            'Authorization': "Bearer" + " " + self.token
        }
        response = requests.delete(deleteurl,headers=getheaders,verify=False)
        print (response.text)
        print (response.status_code)
if __name__ == '__main__':
    user = User("admin","socpcloud")
    print("---Clusters---")
    clusterlist = user.getclusters()
    print("---Projects---")
    for id in clusterlist:
        user.getproject(id)
    #user.createtemplate()
    #user.deletetemplate("fyh-test-service")
    user.createapplication()
    #user.updateapplication("service","fyhtest")
    #user.deleteapplication("gpu-test","local:p-rnqx6","service")