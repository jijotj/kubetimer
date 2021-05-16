**kubetimer** is a Kubernetes watcher that currently publishes latency metrics for pod state transitions. Currently 
supported metrics are:
- pod_scheduled_time: Time taken to schedule pods
- pod_init_time: Time taken for pods to initialize
- pod_containers_ready_time: Time taken for containers in pods to be ready
- pod_ready_time: Time taken for pods to be ready

## Install

