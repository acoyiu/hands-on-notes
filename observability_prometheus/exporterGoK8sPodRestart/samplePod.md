# sample pod

```json
{
    "kind": "PodList",
    "apiVersion": "v1",
    "metadata": {
        "selfLink": "/api/v1/namespaces/jaeger/pods",
        "resourceVersion": "14759701"
    },
    "items": [
        {
            "metadata": {
                "name": "jaeger-deployment-5574878fdc-7cbfv",
                "generateName": "jaeger-deployment-5574878fdc-",
                "namespace": "jaeger",
                "selfLink": "/api/v1/namespaces/jaeger/pods/jaeger-deployment-5574878fdc-7cbfv",
                "uid": "02206347-5386-4f7f-8e75-6b78dc538d5e",
                "resourceVersion": "14606394",
                "creationTimestamp": "2022-07-27T13:36:57Z",
                "labels": {
                    "pod-template-hash": "5574878fdc",
                    "usage": "jaeger"
                },
                "annotations": {
                    "cni.projectcalico.org/podIP": "10.1.43.74/32",
                    "cni.projectcalico.org/podIPs": "10.1.43.74/32",
                    "kubernetes.io/limit-ranger": "LimitRanger plugin set: memory request for container jaeger-container; memory limit for container jaeger-container"
                },
                "ownerReferences": [
                    {
                        "apiVersion": "apps/v1",
                        "kind": "ReplicaSet",
                        "name": "jaeger-deployment-5574878fdc",
                        "uid": "63d0fb79-7588-4c4e-a6b8-4d67412831cd",
                        "controller": true,
                        "blockOwnerDeletion": true
                    }
                ],
                "managedFields": [
                    {
                        "manager": "calico",
                        "operation": "Update",
                        "apiVersion": "v1",
                        "time": "2022-07-27T13:36:58Z",
                        "fieldsType": "FieldsV1",
                        "fieldsV1": {
                            "f:metadata": {
                                "f:annotations": {
                                    "f:cni.projectcalico.org/podIP": {},
                                    "f:cni.projectcalico.org/podIPs": {}
                                }
                            }
                        }
                    },
                    {
                        "manager": "kubelite",
                        "operation": "Update",
                        "apiVersion": "v1",
                        "time": "2022-08-11T08:29:11Z",
                        "fieldsType": "FieldsV1",
                        "fieldsV1": {
                            "f:metadata": {
                                "f:generateName": {},
                                "f:labels": {
                                    ".": {},
                                    "f:pod-template-hash": {},
                                    "f:usage": {}
                                },
                                "f:ownerReferences": {
                                    ".": {},
                                    "k:{\"uid\":\"63d0fb79-7588-4c4e-a6b8-4d67412831cd\"}": {
                                        ".": {},
                                        "f:apiVersion": {},
                                        "f:blockOwnerDeletion": {},
                                        "f:controller": {},
                                        "f:kind": {},
                                        "f:name": {},
                                        "f:uid": {}
                                    }
                                }
                            },
                            "f:spec": {
                                "f:containers": {
                                    "k:{\"name\":\"jaeger-container\"}": {
                                        ".": {},
                                        "f:env": {
                                            ".": {},
                                            "k:{\"name\":\"COLLECTOR_OTLP_ENABLED\"}": {
                                                ".": {},
                                                "f:name": {},
                                                "f:value": {}
                                            }
                                        },
                                        "f:image": {},
                                        "f:imagePullPolicy": {},
                                        "f:name": {},
                                        "f:ports": {
                                            ".": {},
                                            "k:{\"containerPort\":14268,\"protocol\":\"TCP\"}": {
                                                ".": {},
                                                "f:containerPort": {},
                                                "f:name": {},
                                                "f:protocol": {}
                                            },
                                            "k:{\"containerPort\":16686,\"protocol\":\"TCP\"}": {
                                                ".": {},
                                                "f:containerPort": {},
                                                "f:name": {},
                                                "f:protocol": {}
                                            }
                                        },
                                        "f:resources": {},
                                        "f:terminationMessagePath": {},
                                        "f:terminationMessagePolicy": {}
                                    }
                                },
                                "f:dnsPolicy": {},
                                "f:enableServiceLinks": {},
                                "f:restartPolicy": {},
                                "f:schedulerName": {},
                                "f:securityContext": {},
                                "f:terminationGracePeriodSeconds": {}
                            },
                            "f:status": {
                                "f:conditions": {
                                    "k:{\"type\":\"ContainersReady\"}": {
                                        ".": {},
                                        "f:lastProbeTime": {},
                                        "f:lastTransitionTime": {},
                                        "f:status": {},
                                        "f:type": {}
                                    },
                                    "k:{\"type\":\"Initialized\"}": {
                                        ".": {},
                                        "f:lastProbeTime": {},
                                        "f:lastTransitionTime": {},
                                        "f:status": {},
                                        "f:type": {}
                                    },
                                    "k:{\"type\":\"Ready\"}": {
                                        ".": {},
                                        "f:lastProbeTime": {},
                                        "f:lastTransitionTime": {},
                                        "f:status": {},
                                        "f:type": {}
                                    }
                                },
                                "f:containerStatuses": {},
                                "f:hostIP": {},
                                "f:phase": {},
                                "f:podIP": {},
                                "f:podIPs": {
                                    ".": {},
                                    "k:{\"ip\":\"10.1.43.74\"}": {
                                        ".": {},
                                        "f:ip": {}
                                    }
                                },
                                "f:startTime": {}
                            }
                        }
                    }
                ]
            },
            "spec": {
                "volumes": [
                    {
                        "name": "kube-api-access-msgjz",
                        "projected": {
                            "sources": [
                                {
                                    "serviceAccountToken": {
                                        "expirationSeconds": 3607,
                                        "path": "token"
                                    }
                                },
                                {
                                    "configMap": {
                                        "name": "kube-root-ca.crt",
                                        "items": [
                                            {
                                                "key": "ca.crt",
                                                "path": "ca.crt"
                                            }
                                        ]
                                    }
                                },
                                {
                                    "downwardAPI": {
                                        "items": [
                                            {
                                                "path": "namespace",
                                                "fieldRef": {
                                                    "apiVersion": "v1",
                                                    "fieldPath": "metadata.namespace"
                                                }
                                            }
                                        ]
                                    }
                                }
                            ],
                            "defaultMode": 420
                        }
                    }
                ],
                "containers": [
                    {
                        "name": "jaeger-container",
                        "image": "jaegertracing/all-in-one:1.35",
                        "ports": [
                            {
                                "name": "dashboard",
                                "containerPort": 16686,
                                "protocol": "TCP"
                            },
                            {
                                "name": "trace",
                                "containerPort": 14268,
                                "protocol": "TCP"
                            }
                        ],
                        "env": [
                            {
                                "name": "COLLECTOR_OTLP_ENABLED",
                                "value": "true"
                            }
                        ],
                        "resources": {
                            "limits": {
                                "memory": "1Gi"
                            },
                            "requests": {
                                "memory": "1Gi"
                            }
                        },
                        "volumeMounts": [
                            {
                                "name": "kube-api-access-msgjz",
                                "readOnly": true,
                                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
                            }
                        ],
                        "terminationMessagePath": "/dev/termination-log",
                        "terminationMessagePolicy": "File",
                        "imagePullPolicy": "IfNotPresent"
                    }
                ],
                "restartPolicy": "Always",
                "terminationGracePeriodSeconds": 30,
                "dnsPolicy": "ClusterFirst",
                "serviceAccountName": "default",
                "serviceAccount": "default",
                "nodeName": "ppwikay-all-series",
                "securityContext": {},
                "schedulerName": "default-scheduler",
                "tolerations": [
                    {
                        "key": "node.kubernetes.io/not-ready",
                        "operator": "Exists",
                        "effect": "NoExecute",
                        "tolerationSeconds": 300
                    },
                    {
                        "key": "node.kubernetes.io/unreachable",
                        "operator": "Exists",
                        "effect": "NoExecute",
                        "tolerationSeconds": 300
                    }
                ],
                "priority": 0,
                "enableServiceLinks": true,
                "preemptionPolicy": "PreemptLowerPriority"
            },
            "status": {
                "phase": "Running",
                "conditions": [
                    {
                        "type": "Initialized",
                        "status": "True",
                        "lastProbeTime": null,
                        "lastTransitionTime": "2022-07-27T13:36:58Z"
                    },
                    {
                        "type": "Ready",
                        "status": "True",
                        "lastProbeTime": null,
                        "lastTransitionTime": "2022-08-11T08:29:11Z"
                    },
                    {
                        "type": "ContainersReady",
                        "status": "True",
                        "lastProbeTime": null,
                        "lastTransitionTime": "2022-08-11T08:29:11Z"
                    },
                    {
                        "type": "PodScheduled",
                        "status": "True",
                        "lastProbeTime": null,
                        "lastTransitionTime": "2022-07-27T13:36:57Z"
                    }
                ],
                "hostIP": "192.168.0.35",
                "podIP": "10.1.43.74",
                "podIPs": [
                    {
                        "ip": "10.1.43.74"
                    }
                ],
                "startTime": "2022-07-27T13:36:58Z",
                "containerStatuses": [
                    {
                        "name": "jaeger-container",
                        "state": {
                            "running": {
                                "startedAt": "2022-08-11T08:29:10Z"
                            }
                        },
                        "lastState": {
                            "terminated": {
                                "exitCode": 137,
                                "reason": "Error",
                                "startedAt": "2022-08-09T09:21:38Z",
                                "finishedAt": "2022-08-11T08:29:08Z",
                                "containerID": "containerd://a2276ed254282f775192a3e4aa3174c45fcdc6f28f0fc327bc870bd2d7a37355"
                            }
                        },
                        "ready": true,
                        "restartCount": 2,
                        "image": "docker.io/jaegertracing/all-in-one:1.35",
                        "imageID": "docker.io/jaegertracing/all-in-one@sha256:0e0c081fd0cf40e9028533d386177fc2492e4e839e1151237e847b2c9486715b",
                        "containerID": "containerd://8fb5faea1ae195495013a415385dde0bfd0e3de991ae56d7d03f8b4199e678a9",
                        "started": true
                    }
                ],
                "qosClass": "Burstable"
            }
        }
    ]
}
```