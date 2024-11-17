# kro Application example

This example creates a ResourceGroup called `App` and then instaciates it with
the default nginx container image.

### Create ResourceGroup called App

Apply the RG to your cluster:

```
kubectl apply -f rg.yaml
```

Validate the RG status is Active:

```
kubectl get rg app.kro.run
```

Expected result:

```
NAME          APIVERSION   KIND   STATE    AGE
app.kro.run   v1alpha1     App    Active    6m
```

### Create an Instance of kind App

Apply the provided instance.yaml

```
kubectl apply -f instance.yaml
```

Validate instance status:

```
kubectl get apps test-app
```

Expected result:

```
NAME       STATE    SYNCED   AGE
test-app   ACTIVE   True     16m
```

### Validate the app is working

Get the ingress url:

```
kubectl get ingress test-app -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
```

Either navigate in the browser or curl it:

```
curl -s $(kubectl get ingress test-app -o jsonpath='{.status.loadBalancer.ingress[0].hostname}') | sed -n '/<body>/,/<\/body>/p' | sed -e 's/<[^>]*>//g'
```

Expected result:

```
Welcome to nginx!
If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.

For online documentation and support please refer to
nginx.org.
Commercial support is available at
nginx.com.

Thank you for using nginx.
```

### Clean up

Remove the instance:

```
kubectl delete apps test-app
```

Remove the resourcegroup:

```
kubectl delete rg app.kro.run
```