import * as k8s from '@kubernetes/client-node';

// load the default path to the token
const kc = new k8s.KubeConfig();
kc.loadFromDefault();

// k8s client with such token and session pod, cm
const k8sApi_corev1 = kc.makeApiClient(k8s.CoreV1Api);

// k8s client with such token and session for deployment and others
const k8sApi_appv1 = kc.makeApiClient(k8s.AppsV1Api);

k8sApi_corev1.listNamespacedPod('default').then((res: any) => {
    console.log(res.body);
});