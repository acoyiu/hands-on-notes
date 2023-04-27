use futures::prelude::*;
use k8s_openapi::api::{
    apps::v1::{Deployment, ReplicaSet},
    core::v1::Pod,
};
use serde_json::json;
use std::env;
use tracing::*;

use kube::{
    api::{Api, DeleteParams, ListParams, Patch, PatchParams, PostParams, ResourceExt},
    runtime::watcher,
    runtime::WatchStreamExt,
    Client,
};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::fmt::init();

    let cli_args: Vec<String> = env::args().collect();

    let action = cli_args.get(1).expect("No argument entered");

    //
    // --- Get K8s Config -------------------------------------------------
    //

    /*
      Get default config on "/root/.kube/config"
      or when failed, get sa token from "/var/run/secrets/kubernetes.io/serviceaccount"
      ::try_from() for Rust Trait of Config Obj
    */

    // set KUBECONFIG if needed
    if let Some(arg2) = env::args().nth(2) {
        env::set_var("KUBECONFIG", &arg2);
        match env::var("KUBECONFIG") {
            Ok(v) => {
                if v != "" {
                    println!("KUBECONFIG specified as {}", v)
                }
            }
            Err(e) => println!("{}", e),
        }
    }

    // will try load 1:KUBECONFIG - 2:~/.kube/config - 3:/var/run/secrets/kubernetes.io/serviceaccount
    let client = Client::try_default().await?;

    //
    // --- CRUD Resources ------------------------------------------------------
    //

    //
    // ===== Create Deployment ====================
    //
    if action == "create-deploy" {
        /* target to current or specific Namespace */
        let deployment_in_target_ns: Api<Deployment> = Api::default_namespaced(client.clone());
        //  deployment_in_target_ns: Api<Pod> = Api::namespaced(client, "kube-system");
        // Api::all 不能用於 client.get, 只能用於 List 等 !!

        info!("Creating Deployment instance blog");
        let deployment_details: Deployment = serde_json::from_value(json!({
            "kind": "Deployment",
            "apiVersion": "apps/v1",
            "metadata": {
                "name": "blog",
            },
            "spec": {
                "replicas": 1,
                "selector": { "matchLabels": { "app": "blog" } },
                "template": {
                    "metadata": { "labels": { "app": "blog" }},
                    "spec": {"containers": [ { "name": "nginx", "image": "nginx:alpine" } ]}
                }
            }
        }))?;

        // Post Param is used for something like "--dry-run"
        let post_param = PostParams::default();

        // Handle cases for successful or failed pod creation
        match deployment_in_target_ns
            .create(&post_param, &deployment_details)
            .await
        {
            Ok(created_deploy_instance) => {
                // name_any => meta.name || meta.generatedName
                let name = created_deploy_instance.name_any();
                info!("Created Deployment with Name: {}", name);
            }
            // if case of create failed by k8s
            Err(kube::Error::Api(err_res)) => {
                info!(err_res.code);
                info!(err_res.message);
            }
            // any other case is probably bad
            Err(e) => {
                print!("error = {}", e);
                return Err(e.into());
            }
        }
    }
    //
    // ===== Read Pod ====================
    //
    else if action == "read-pod" {
        // get pod in namespace
        let pod_in_target_ns: Api<Pod> = Api::default_namespaced(client.clone());

        // add label filter
        let mut list_param = ListParams::default();
        list_param.label_selector = Option::Some(String::from("app=blog"));

        // Start Query
        let found_pod = pod_in_target_ns.list(&list_param).await?;
        for pod_instance in found_pod {
            info!("Found Pod: {}", pod_instance.name_any());
        }
    }
    //
    // ===== Read Deployment ====================
    //
    else if action == "read-deploy" {
        let deployment: Api<Deployment> = Api::default_namespaced(client.clone());
        let asd = deployment.get("blog").await.unwrap();
        info!(
            "deploy name: {}, uid: {} ",
            asd.name_any(),
            asd.metadata.uid.unwrap()
        );
    }
    // ===== Update Deploy ====================
    else if action == "update-deploy" {
        let deployment: Api<Deployment> = Api::default_namespaced(client.clone());

        let replica_num = env::args().nth(3).expect("3rd argr must be provided");

        let replica_num = replica_num
            .parse::<i32>()
            .expect("3rd argr must be digit(int)");

        // Create an updated version of the resource
        let patch = serde_json::json!({
            "metadata": { "name": "blog" },
            "spec": { "replicas": replica_num }
        });

        let params = PatchParams::default();
        let patch = Patch::Merge(&patch);
        match deployment.patch("blog", &params, &patch).await {
            Ok(updated_deployment) => {
                info!(
                    "updated Deployment: {} to replicas {}",
                    updated_deployment.name_any(),
                    replica_num
                );
            }
            Err(e) => {
                print!("update failed = {}", e);
            }
        }
    }
    //
    // ===== Delete Deploy ====================
    //
    else if action == "delete-deploy" {
        let deployment: Api<Deployment> = Api::default_namespaced(client.clone());
        deployment.delete("blog", &DeleteParams::default()).await?;
        info!("Deployment {} deleted", "blog");
    }
    //
    // ===== pod Watch ========================
    //
    else if action == "watch-pod" {
        let pod_watcher: Api<Pod> = Api::default_namespaced(client.clone());
        let target_pod_label = ListParams::default();
        watcher(pod_watcher, target_pod_label)
            .applied_objects()
            .try_for_each(|pod_instance| async move {
                info!("saw {}", pod_instance.name_any());
                Ok(())
            })
            .await?;
    }
    //
    // ===== rs Watch ========================
    //
    else if action == "watch-rs" {
        let rs_watcher: Api<ReplicaSet> = Api::default_namespaced(client.clone());
        let target_rs_label = ListParams::default();
        watcher(rs_watcher, target_rs_label)
            .applied_objects()
            .try_for_each(|rs_instance| async move {
                info!("saw rs {}", rs_instance.name_any());
                Ok(())
            })
            .await?;
    }
    //
    // ===== deploy Watch ========================
    //
    else if action == "watch-deploy" {
        let dp_watcher: Api<Deployment> = Api::default_namespaced(client.clone());
        let target_dp_label = ListParams::default();
        watcher(dp_watcher, target_dp_label)
            .applied_objects()
            .try_for_each(|deployment_instance| async move {
                info!("saw deployment {}", deployment_instance.name_any());
                Ok(())
            })
            .await?;
    }
    //
    // ===== No action found ==================
    //
    else {
        panic!("No correct action specified");
    }

    Ok(())
}
