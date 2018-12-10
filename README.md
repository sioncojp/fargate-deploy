# fargate-deploy

- (optional): Log in to ECR
- (optional): Lifecycle Image in ECR Repository
    - Remove unused images. 
    - Ignore git CcommitHash, latest, prod stg ... etc from config.toml. 
- (optional): Build & push on Dockerfile
    - set tag
        - prod, stg .... etc from config.toml.
        - git commitHash from current.
        - If prod, latest is set.
- (required): Update and register task. And update Fargate service.

## Usage

sample [config.toml](examples/config.toml)

```shell
fargate-deploy -h

fargate-deploy -c config.toml -e prod -i "111111111" -p
```

## Development

```
make help

make dep/init
make dep

### tar.gz for darwin-amd64, linux-amd64
make dist
```

# License
The MIT License

Copyright Shohei Koyama / sioncojp 

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.