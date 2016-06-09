# semver
Semantic Versioning as a Service

### Create Semver
```
$ curl https://semver.co/new
e84e9872-fbf7-4d76-b222-68ba1f3e72b3
```

### Get Current Version
```
$ curl "https://semver.co/v1/e84e9872-fbf7-4d76-b222-68ba1f3e72b3"
0.0.1
```

### Bump Version
```
$ curl "https://semver.co/v1/e84e9872-fbf7-4d76-b222-68ba1f3e72b3/bump"
0.0.2
$ curl "https://semver.co/v1/e84e9872-fbf7-4d76-b222-68ba1f3e72b3/bump?type=major"
1.0.0
$ curl "https://semver.co/v1/e84e9872-fbf7-4d76-b222-68ba1f3e72b3/bump?type=minor"
1.1.0
```

### Set Version
```
$ curl "https://semver.co/v1/e84e9872-fbf7-4d76-b222-68ba1f3e72b3" -d "version=3.1.0"
3.1.0
```

### List Versions (History)
```
$ curl "https://semver.co/v1/e84e9872-fbf7-4d76-b222-68ba1f3e72b3/history"
0.0.1
0.0.2
1.0.0
1.1.0
3.1.0
```

### Delete Project
```
$ curl "https://semver.co/v1/e84e9872-fbf7-4d76-b222-68ba1f3e72b3" -XDELETE
ok
```

### XML, JSON, and Plain-Text Response
```
$ curl "https://semver.co/v1/e84e9872-fbf7-4d76-b222-68ba1f3e72b3?output=json"
{
  "version": "0.0.1",
  "major": 0,
  "minor": 0,
  "patch": 1
}
$ curl "https://semver.co/v1/e84e9872-fbf7-4d76-b222-68ba1f3e72b3?output=xml"
<Versioning>
   <version>0.0.1</version>
   <major>0</major>
   <minor>0</minor>
   <patch>1</patch>
</Versioning>
```

## Contributing

Everyone is encouraged to help improve this project. Here are a few ways you can help:

- [Report bugs](https://github.com/samuelngs/semver/issues)
- Fix bugs and [submit pull requests](https://github.com/samuelngs/semver/pulls)
- Write, clarify, or fix documentation
- Suggest or add new features

## License

This project is distributed under the MIT license found in the [LICENSE](./LICENSE) file.

```
The MIT License (MIT)

Copyright (c) 2016 Samuel

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

