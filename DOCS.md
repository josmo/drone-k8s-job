Use this drone plugin to creates and watches a K8s Job.  The idea is to be able to deploy
a service and then deploy a job that can test the service in k8s and clean it up.

The following parameters are used to configure this plugin:

- `url` - url to your cluster server
- `token` - Token used to connect to the cluster
- `ca` - certificate auth to connect to cluster
- `insecure` - allow for insecure cluster connection 
- `namespace` - namespace (will use default if not set)  
- `template` - file location for the job template [sample](job.yml)
- `cleanup` - default true: will remove the job upon success or failure
- `timeout` -  default 120: will timeout watching the job after 120 seconds and "try" to clean up the job
- `debug` - default false: will add debug information
The following is a sample k8s deployment configuration in your `.drone.yml` file:

```yaml
  - name: k8s-job
    image: pelotech/drone-k8s-job
    settings:
      url: https://k8s.server/
      token: asldkfj
      insecure: false
      template: job.yml
```

Or with no cleanup, different timeout, and different namespace

```yaml
  - name: k8s-job
    image: pelotech/drone-k8s-job
    settings: 
      url: https://k8s.server/
      token: asldkfj
      namespace: testing
      insecure: false
      template: job.yml
      cleanup: false
      debug: true
      timeout: 200
```

if you want to add secrets for the token it's KUBERNETES_TOKEN, KUBERNETES_URL, KUBERNETES_CA, JOB_TEMPLATE
