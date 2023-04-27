use anyhow::Result;
use futures::StreamExt;
use k8s_openapi::api::core::v1::{Container, PodSpec};
use k8s_openapi::apimachinery::pkg::apis::meta::v1::OwnerReference;
use kube::api::DeleteParams;
use kube::core::ObjectMeta;
use kube::runtime::watcher::Error;
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use serde_json::json;

use std::env;
use std::io::BufRead;
use std::time::Duration;
use std::{collections::BTreeMap, sync::Arc};
use tracing::*;
use validator::Validate;

use k8s_openapi::{
    api::core::v1::Pod,
    apiextensions_apiserver::pkg::apis::apiextensions::v1::CustomResourceDefinition,
};
use kube::{
    api::{Api, ListParams, PostParams, ResourceExt},
    core::crd::CustomResourceExt,
    runtime::controller::{Action, Controller},
    Client, CustomResource,
};

//
// ---- crd -----------------------------------------------------------------------------------------------------
//

/*
binding   APIVersion       Kind
acoapps   aco.app.dev/v1   AcoApp
*/
// "k get acoapps.aco.app.dev"

// #[kube(scale = r#"{"specReplicasPath":".spec.replicas", "statusReplicasPath":".status.replicas"}"#)]
#[derive(CustomResource, Deserialize, Serialize, Clone, Debug, Validate, JsonSchema)]
#[kube(group = "aco.app.dev", version = "v1", kind = "AcoApp", namespaced)]
#[kube(status = "AcoAppStatus")]
#[kube(printcolumn = r#"{"name":"Target", "jsonPath": ".spec.target", "type": "string"}"#)]
#[kube(printcolumn = r#"{"name":"Ready", "jsonPath": ".status.is_ready", "type": "string"}"#)]
#[kube(printcolumn = r#"{"name":"Status", "jsonPath": ".status.state", "type": "string"}"#)]
pub struct AcoAppSpec {
    #[validate(length(min = 3))]
    target: String,
}

#[derive(Deserialize, Serialize, Clone, Debug, Default, JsonSchema)]
pub struct AcoAppStatus {
    is_ready: bool,
    state: String,
    child_pod_uid: String,
}

//
// ---------------------------------------------------------------------------------------------------------
//

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt::init();

    /*
      Get default config on "/root/.kube/config"
      or when failed, get sa token from "/var/run/secrets/kubernetes.io/serviceaccount"
      ::try_from() for Rust Trait of Config Obj
    */
    // set KUBECONFIG if needed
    if let Some(arg1) = env::args().nth(2) {
        env::set_var("KUBECONFIG", &arg1);
        match env::var("KUBECONFIG") {
            Ok(v) => {
                if v != "" {
                    println!("KUBECONFIG specified as {}", v)
                }
            }
            Err(e) => println!("{}", e),
        }
    }

    //
    // --- Get K8s Config -------------------------------------------------
    //

    // will try load 1:KUBECONFIG - 2:~/.kube/config - 3:/var/run/secrets/kubernetes.io/serviceaccount
    let client = Client::try_default().await?;

    //
    // --- Enter Loop -------------------------------------------------
    //

    dbg!("Enter Signals:");
    for input_str in std::io::BufReader::new(std::io::stdin()).lines() {
        let input_str = input_str.unwrap_or_default();
        dbg!(&input_str);
        match input_str.as_str() {
            "create crd" => {
                create_crd(&client).await;
            }
            "create instance" => {
                create_instance(&client).await;
            }
            "watch" => {
                watch(&client).await;
            }
            "update instance" => {
                update_status(&client).await;
            }
            "delete instance" => {
                delete_instance(&client).await;
            }
            "exit" => {
                std::process::exit(0);
            }
            &_ => {
                dbg!("Error: no valid action set");
                ()
            }
        };
        dbg!("Enter Signals:");
    }

    Ok(())
}

async fn create_crd(client: &Client) {
    let crds: Api<CustomResourceDefinition> = Api::all(client.clone());

    // Create the CRD so we can create AcoApp in kube
    let aco_app_ard = AcoApp::crd();
    let post_params = PostParams::default();

    // create crd
    match crds.create(&post_params, &aco_app_ard).await {
        Ok(created_crd) => info!(
            "Created CRD: {}",
            created_crd.metadata.name.unwrap_or_default()
        ),
        Err(kube::Error::Api(ae)) => info!("api error = {}, {}", ae.code, ae.message),
        Err(e) => {
            dbg!(e);
        }
    };
}

async fn create_instance(client: &Client) {
    info!("Creating AcoApp instance the-aco-app");

    let mut app_instance = AcoApp::new(
        "the-aco-app",
        AcoAppSpec {
            target: String::from("google.com"),
        },
    );

    // add annotation
    let mut map = BTreeMap::new();
    map.insert(String::from("usage"), String::from("aco"));
    app_instance.metadata.annotations = Some(map);

    // create client
    let aco_app: Api<AcoApp> = Api::default_namespaced(client.clone());
    let post_params = PostParams::default();

    // create app resource instance
    let aco_app_instance = aco_app.create(&post_params, &app_instance).await.unwrap();
    info!("Created aco.app instance: {}", aco_app_instance.name_any());
}

async fn watch(client: &Client) {
    info!("Now watching crd resource...");

    // the list of watching distance
    let crd_instance_watcher: Api<AcoApp> = Api::all(client.clone());
    let pod_instance_watcher: Api<Pod> = Api::all(client.clone());

    info!("starting controller");

    // Controller::run 是用於啟動控制器並在事件發生時處理事件的主要函數，
    // 而 Controller::reconcile_all_on 是用於對特定類型的所有資源執行一次性協調的實用函數。

    Controller::new(crd_instance_watcher, ListParams::default())
        .owns(pod_instance_watcher, ListParams::default()) // 必須要有 OwnerReference 才有效
        .shutdown_on_signal() // use graceful shutdown
        .run(
            reconcile,
            error_policy,
            Arc::new(Data {
                client: client.clone(),
            }),
        )
        .for_each(|res| async move {
            match res {
                Ok(o) => info!("reconciled {:?}", o),
                Err(e) => warn!("reconcile failed: {}", e),
            }
        })
        .await;
}

//
// =--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--
//

async fn update_status(client: &Client) {
    // 將 instance Status 更新
    let name_of_res = "the-aco-app";

    // get the resource
    let aco_app: Api<AcoApp> = Api::default_namespaced(client.clone());
    let found_res = aco_app.get(&name_of_res).await.unwrap();

    // Update status on instance (cannot be done through replace/create/patch direct)
    let json_spec = json!({
        "apiVersion": "aco.app.dev/v1",
        "kind": "AcoApp",
        "metadata": {
            "name": &name_of_res,
            "resourceVersion": found_res.resource_version(), // Updates need to provide our last observed version
        },
        "status": AcoAppStatus { is_ready: false, child_pod_uid:String::from(""), state: String::from("Pending") }
    });

    // Post param for update
    let post_params = PostParams::default();

    // update target resource's status
    let aco_app_instance = aco_app
        .replace_status(
            &name_of_res,
            &post_params,
            serde_json::to_vec(&json_spec).expect("must able to parse"),
        )
        .await;

    let status = aco_app_instance
        .expect("must get status here")
        .status
        .unwrap_or_default();
    info!("Status Replaced -- is_ready: {}", status.is_ready);
}

//
// =--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--=--
//

use either::{Left, Right};

async fn delete_instance(client: &Client) {
    info!("Deleting instance...");

    let name_of_res = "the-aco-app";

    let aco_app: Api<AcoApp> = Api::default_namespaced(client.clone());
    let dp = DeleteParams::default();

    match aco_app.delete(name_of_res, &dp).await {
        Ok(ei) => match ei {
            Left(l) => {
                print!("Left: {}\n", l.metadata.name.expect("must ok"));
            }
            Right(r) => {
                print!(
                    "Right: {}, {}, {}, {}\n",
                    r.code,
                    r.details.clone().unwrap().name,
                    r.details.clone().unwrap().kind,
                    r.reason
                );
            }
        },
        Err(e) => {
            dbg!(e);
        }
    }
}

//
// ================================================================
//

// Data we want access to in error/reconcile calls
struct Data {
    client: Client,
}

// The controller triggers this on reconcile errors
fn error_policy(_object: Arc<AcoApp>, _error: &Error, _ctx: Arc<Data>) -> Action {
    Action::requeue(Duration::from_secs(1))
}

/// Controller triggers this whenever our main object or our children changed
async fn reconcile(_generator: Arc<AcoApp>, context: Arc<Data>) -> Result<Action, Error> {
    /*
    Reconcile 有兩種結局
        - requeue
        - await_change
     */

    /*
    Action::requeue
        - 設置若干時間後 reconcile 一次
        - 這是最佳實踐操作，即使在錯過更改（可能發生）的情況下也能確保控制器的最終一致性。
        - 通常不會錯過觀察事件，因此每小時運行一次（“默認”）作為後備是合理的。
        - Ok(Action::requeue(Duration::from_secs(300)))
     */

    /*
    Action::await_change
        - 在檢測到更改之前不執行任何操作
        - 這會停止 Controller 定期協調此對象，直到檢測到相關的監視事件。
        - 如果您的 watch 不同步，則可能會完全錯過 changes。
        - 因此，不建議以這種方式禁用重新排隊
        - 除非您經常更改底層對象，或使用其他一些掛鉤來保持最終一致性。
        - Ok(Action::await_change())
     */

    custom_reconciler(_generator, context).await
}

//
// Reconciler ---------------------------------------------------------
//

async fn custom_reconciler(
    current_resource: Arc<AcoApp>, // the _generator is the resource itself
    data: Arc<Data>,
) -> Result<Action, Error> {
    //
    // The Arc<Data> here is the context passed back from above "Data" instance
    let api_version = data.client.apiserver_version().await.unwrap_or_default();
    info!("api_version: {:?}", api_version);

    let name = current_resource.metadata.name.clone().unwrap_or_default();
    info!("name: {}", &name);

    // check current resource status
    let status = current_resource.status.clone().unwrap_or_default();
    info!("status: {:?}", status);

    let context: Client = data.client.clone();
    let current_resource = current_resource.as_ref().clone();

    // // cast current resource from Arc<AcoApp> to AcoApp
    // let current_resource: AcoApp = current_resource.clone().as_ref().clone();

    // have a match on status.state,
    // if status.state is empty, means non-initialized resource,
    //     create child pod with owner reference to current crd resource,
    //         and update current resource "status.state" as "Pending"
    //         and update current resource "status.child_pod_uid" as child pod uid
    //         and update current resource "status.is_ready" as false
    //             then finish reconcile and requeue
    // else if status.state is not empty, means initializing resource or in wanted status,
    //     which has state as "Pending", "Running"
    //     if status.state is "Pending", check child pod status,
    //         if child pod is existed and running, update current resource status.state as "Running"
    //         if child pod is existed and not running, finish reconcile and requeue
    //         if child pod is not existed, create child pod with owner reference to current crd resource,
    //             and update current resource "status.state" as "Pending"
    //             and update current resource "status.child_pod_uid" as child pod uid
    //             and update current resource "status.is_ready" as false
    //     if status.state is "Running", do nothing

    info!("status.state: {:?}", status.state);
    match status.state.as_str() {
        "" => {
            info!("current resource is non-initialized, create child pod...");

            // create child pod
            match create_child_pod(&context, &current_resource).await {
                Ok(pod) => {
                    info!("pod: {:?}", pod);

                    // update current resource status
                    update_current_resource_status(&context, &current_resource, pod).await;
                }
                Err(e) => {
                    panic!("error: {:?}", e);
                }
            }

            Ok(Action::requeue(Duration::from_secs(300)))
        }
        _ => {
            info!("current resource is initialized, check child pod status...");

            match status.state.as_str() {
                "Pending" => {
                    info!("current resource is in pending status, check child pod status...");

                    // check child pod status
                    match check_child_pod_status(&context, &current_resource).await {
                        Ok(pod) => {
                            info!("pod: {:?}", pod);

                            // update current resource status
                            update_current_resource_status(&context, &current_resource, pod).await;
                        }
                        Err(e) => {
                            panic!("error: {:?}", e);
                        }
                    }

                    Ok(Action::requeue(Duration::from_secs(300)))
                }
                "Running" => {
                    info!("current resource is in running status, do nothing...");
                    Ok(Action::requeue(Duration::from_secs(300)))
                }
                _ => Ok(Action::requeue(Duration::from_secs(300))),
            }
        }
    }
}

async fn check_child_pod_status(
    context: &Client,
    current_resource: &AcoApp,
) -> Result<Pod, kube::Error> {
    //
    // declare k8s api client
    let pod_name = format!(
        "{}-pod",
        current_resource.metadata.name.clone().unwrap_or_default()
    );

    // get child pod
    Api::default_namespaced(context.clone())
        .get(&pod_name)
        .await
}

async fn update_current_resource_status(context: &Client, current_resource: &AcoApp, pod: Pod) {
    //
    // declare k8s api client
    let aco_app_client: Api<AcoApp> = Api::default_namespaced(context.clone());

    info!("----------1----------");

    // // update current resource status
    // let current_resource = AcoApp {
    //     metadata: ObjectMeta {
    //         name: Some(current_resource.metadata.name.clone().unwrap_or_default()),
    //     },
    //     spec: AcoAppSpec {
    //         target: current_resource.spec.target.clone(),
    //     },
    //     status: Some(AcoAppStatus {
    //         state: "Running".to_string(),
    //         child_pod_uid: pod.metadata.uid.clone().unwrap_or_default(),
    //         is_ready: true,
    //     }),
    // };
    info!("----------2----------");

    // let qweasd = serde_json::to_vec(&current_resource).unwrap();
    // info!("qweasd: {:?}", qweasd);

    let json_str = json!({
        "apiVersion": "aco.app.dev/v1",
        "kind": "AcoApp",
        "metadata": {
            "name": "the-aco-app",
            // Updates need to provide our last observed version:
            "resourceVersion": current_resource.resource_version(),
        },
        "status":  AcoAppStatus {
            state: "Running".to_string(),
            child_pod_uid: pod.metadata.uid.clone().unwrap_or_default(),
            is_ready: true,
        }
    });

    info!("----------3----------");

    // update current resource status
    match aco_app_client
        .replace_status(
            &current_resource.metadata.name.clone().unwrap_or_default(),
            &PostParams::default(),
            serde_json::to_vec(&json_str).expect("must be able to parse to json"),
        )
        .await
    {
        Ok(acoapp) => info!("acoapp: {:?}", acoapp),
        Err(e) => error!("error: {:?}", e),
    };
}

async fn create_child_pod(context: &Client, current_resource: &AcoApp) -> Result<Pod, kube::Error> {
    //
    // declare k8s api client
    let pod_client: Api<Pod> = Api::default_namespaced(context.clone());

    // create pod with owner reference
    let pod = Pod {
        metadata: ObjectMeta {
            name: Some(format!(
                "{}-pod",
                current_resource.metadata.name.clone().unwrap_or_default()
            )),
            owner_references: Some(vec![OwnerReference {
                api_version: "aco.app.dev/v1".to_string(),
                kind: "AcoApp".to_string(),
                name: current_resource.metadata.name.clone().unwrap_or_default(),
                uid: current_resource.metadata.uid.clone().unwrap_or_default(),
                controller: Some(true),
                block_owner_deletion: Some(true),
            }]),
            ..ObjectMeta::default()
        },
        spec: Some(PodSpec {
            containers: vec![Container {
                name: "nginx".to_string(),
                image: Some("nginx".to_string()),
                ..Container::default()
            }],
            ..PodSpec::default()
        }),
        ..Pod::default()
    };

    // create pod
    let pod = pod_client.create(&PostParams::default(), &pod).await?;

    Ok(pod)
}
