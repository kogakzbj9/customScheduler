# customScheduler

## Description

This project implements a custom scheduler plugin for Kubernetes. The custom scheduler plugin allows you to define custom scheduling logic for your Kubernetes cluster.
This project is under construction.

## Building the Custom Scheduler Binary

To build the custom scheduler binary, follow these steps:

1. Ensure you have Go installed on your machine.
2. Navigate to the project directory:
   ```sh
   cd customScheduler
   ```
3. Build the custom scheduler binary:
   ```sh
   go build -o custom-scheduler .
   ```

## Building the Docker Image

To build the Docker image for the custom scheduler, follow these steps:

1. Ensure you have Docker installed on your machine.
2. If you are using Minikube, run the following command to set up the Docker environment:
   ```sh
   eval $(minikube docker-env)
   ```
   This command configures Docker to use the Minikube Docker daemon, allowing you to build Docker images directly inside the Minikube environment.
3. Navigate to the project directory:
   ```sh
   cd customScheduler
   ```
4. Build the Docker image:
   ```sh
   docker build -t custom-scheduler:latest .
   ```

## Pushing the Docker Image to a Registry

To push the Docker image to a registry, follow these steps:

1. Tag the Docker image with the registry URL:
   ```sh
   docker tag custom-scheduler:latest <your-registry-url>/custom-scheduler:latest
   ```
2. Push the Docker image to the registry:
   ```sh
   docker push <your-registry-url>/custom-scheduler:latest
   ```

## Deploying the Custom Scheduler Plugin

To deploy the custom scheduler plugin to your Kubernetes cluster, follow these steps:

1. Create a ConfigMap for the custom scheduler configuration:
   ```sh
   kubectl create configmap custom-scheduler-config --from-file=config.yaml -n kube-system
   ```
   The ConfigMap should contain the following keys with their default values:
   - `cpuThreshold`: 50
   - `waitTime`: 10
2. Verify the creation of the ConfigMap:
   ```sh
   kubectl get configmap custom-scheduler-config -o yaml
   ```
3. Create a service account for the custom scheduler:
   ```sh
   kubectl create serviceaccount custom-scheduler -n kube-system
   ```
4. Create a cluster role binding for the custom scheduler:
   ```sh
   kubectl create clusterrolebinding custom-scheduler --clusterrole=system:kube-scheduler --serviceaccount=kube-system:custom-scheduler
   ```
5. Create a role binding for `extension-apiserver-authentication-reader`:
   ```sh
   kubectl create rolebinding -n kube-system extension-apiserver-authentication-reader --role=extension-apiserver-authentication-reader --serviceaccount=kube-system:custom-scheduler
   ```
6. Create a role binding for `storage-admin`:
   ```sh
   kubectl create rolebinding -n kube-system storage-admin --role=storage-admin --serviceaccount=kube-system:custom-scheduler
   ```
7. Verify the creation of the custom-scheduler deployment:
   ```sh
   kubectl get deployments -n kube-system -l app=custom-scheduler
   ```
8. Deploy the custom scheduler plugin as a Kubernetes Deployment:
   ```sh
   kubectl apply -f custom-scheduler-deployment.yaml -n kube-system
   ```
9. Verify that the custom scheduler plugin is running:
   ```sh
   kubectl get pods -n kube-system -l app=custom-scheduler
   ```
10. Verify the correctness of the command execution results:
   ```sh
   kubectl logs <pod-name> -n kube-system
   ```

## Deploying the Custom Scheduler Using the Docker Image

To deploy the custom scheduler using the Docker image, follow these steps:

1. Update the `custom-scheduler-deployment.yaml` file to use the Docker image from your registry:
   ```yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: custom-scheduler
     namespace: kube-system
     labels:
       app: custom-scheduler
   spec:
     replicas: 1
     selector:
       matchLabels:
         app: custom-scheduler
     template:
       metadata:
         labels:
           app: custom-scheduler
       spec:
         containers:
         - name: custom-scheduler
           image: <your-registry-url>/custom-scheduler:latest
           imagePullPolicy: IfNotPresent
           ports:
           - containerPort: 10259
         serviceAccountName: custom-scheduler
   ```
2. Apply the updated deployment:
   ```sh
   kubectl apply -f custom-scheduler-deployment.yaml
   ```
3. Verify that the custom scheduler is running using the Docker image:
   ```sh
   kubectl get pods -n kube-system -l app=custom-scheduler
   ```

## Testing the Custom Scheduler Plugin

To test the custom scheduler plugin, follow these steps:

1. Create a test pod with the custom scheduler name:
   ```sh
   kubectl apply -f test-pod.yaml
   ```
2. Verify that the test pod is scheduled by the custom scheduler:
   ```sh
   kubectl get pod test-pod -o jsonpath='{.spec.schedulerName}'
   ```

## Using the cpuSpike Annotation

The custom scheduler plugin can also use the `cpuSpike` annotation to override the default CPU threshold for a specific Pod. To use the `cpuSpike` annotation, add it to the Pod's metadata with the desired CPU threshold value. For example:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  annotations:
    cpuSpike: "75"
spec:
  containers:
  - name: test-container
    image: nginx
```

In this example, the `cpuSpike` annotation is set to `75`, which means the custom scheduler will use a CPU threshold of 75% for this Pod instead of the default value.

## Contributing

If you would like to contribute to this project, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.
