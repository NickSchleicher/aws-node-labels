# aws-node-labels

Easily label your kubernetes cluster's nodes based on their EC2 instance tags. Cluster deployments such as EKS does not allow for automatically attaching labels to nodes, a feature that was included with [kops](https://github.com/kubernetes/kops).

## Prerequisites
Your AWS CLI should be setup with the same account your cluster is using
Your kubectl should be setup with your cluster
Your nodes on the cluster should be tagged if you want a label attached to the node

## Example
Running the following
```
aws-node-labels
```
gives the below results:  
  
The default tag delimiter is "k8s-label_"
The following instance(s) with the tag
```"Key: k8s-label_disktype"```
```"Value: ssd"```
Will be tagged with the kubernetes node label "disktype=ssd"

This will allow you to specify which node a pod is deployed on by:
```
  nodeSelector:
    disktype: ssd
```




