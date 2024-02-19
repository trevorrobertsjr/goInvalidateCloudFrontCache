# goInvalidateCloudFrontCache

Compile with the following parameters:

```bash
GOARCH=arm64 GOOS=linux go build -tags lambda.norpc -o bootstrap main.go
```

The Amazon Linux 2023 and Amazon Linux 2 OS-only runtimes require the executable to be called **bootstrap**. Also, they do not feature an RPC client like the go1.x runtime did. So, the **lambda.norpc** tag specifices that the compiler not include the RPC component from the ```aws-lambda-go``` package. This decreases the compiled binary size and improves performance.

You can read more on James Beswick's AWS blog for [Migrating AWS Lambda functions from the Go1.x runtime to the custom runtime on Amazon Linux](https://aws.amazon.com/blogs/compute/migrating-aws-lambda-functions-from-the-go1-x-runtime-to-the-custom-runtime-on-amazon-linux-2/).