# semver
Semantic Versioning as a Service

### Create Semver
```
$ curl https://semver.co/v1/new
e84e9872-fbf7-4d76-b222-68ba1f3e72b3
```

### Get Current Version
```
$ curl https://semver.co/v1/e84e9872-fbf7-4d76-b222-68ba1f3e72b3
0.0.1
```

### Bump Version
```
$ curl https://semver.co/v1/e84e9872-fbf7-4d76-b222-68ba1f3e72b3/bump
0.0.2
$ curl https://semver.co/v1/e84e9872-fbf7-4d76-b222-68ba1f3e72b3/bump?level=major
1.0.0
$ curl https://semver.co/v1/e84e9872-fbf7-4d76-b222-68ba1f3e72b3/bump?level=minor
1.1.0
```
