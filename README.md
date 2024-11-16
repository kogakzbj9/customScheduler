# customScheduler

## Description

This project implements a custom scheduler plugin for Kubernetes. The custom scheduler plugin allows you to define custom scheduling logic for your Kubernetes cluster.
This project is under construction.


## Building the Custom Scheduler Plugin

To build the custom scheduler plugin, follow these steps:

1. Ensure you have a Kubernetes cluster and the necessary tools for development, such as Go and kubectl.
2. Clone the repository:
   ```sh
   git clone https://github.com/kogakzbj9/customScheduler.git
   ```
3. Navigate to the project directory:
   ```sh
   cd customScheduler
   ```
4. Build the custom scheduler plugin:
   ```sh
   go build -o custom-scheduler .
   ```

## Deploying the Custom Scheduler Plugin

To deploy the custom scheduler plugin to your Kubernetes cluster, follow these steps:

1. Create a ConfigMap for the custom scheduler configuration:
   ```sh
   kubectl create configmap custom-scheduler-config --from-file=config.yaml
   ```
   The ConfigMap should contain the following keys with their default values:
   - `cpuThreshold`: 50
   - `waitTime`: 10
2. Deploy the custom scheduler plugin as a Kubernetes Deployment:
   ```sh
   kubectl apply -f custom-scheduler-deployment.yaml
   ```
3. Verify that the custom scheduler plugin is running:
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
