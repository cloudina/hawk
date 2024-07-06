
# HAWK

## Introduction
Multi Cloud antivirus scanning API based on CLAMAV and YARA for AWS S3, AZURE Blob Storage, GCP Cloud Storage.

## Features
-   Microservice for scanning stream with YARA and CLAMAV
-   Scans S3 Bucket Object
-   Moves Clean S3 Objects to another S3 Bucket
-   Quarantines Infected S3 Objects to another S3 Bucket
-   CLAMAV DB auto is updated to latest
-   [TODO] AZURE and GCP support
-   [TODO] Merge Various YARA rules to one set
-   [TODO] Auto Update YARA rules
-   [TODO] Support Yextend
-   [TODO] Improve Logging using logrus [https://github.com/antonfisher/nested-logrus-formatter]
-   [TODO] Harden Image


## API
Available API are
```
POST /scanstream - scan stream

POST -d '{"bucketname": $S3_SCANNING_BUCKET, "key": $S3_OBJECT_TO_SCAN, "clean_files_bucket": $S3_CLEAN_FILES_BUCKET, "qurantine_files_bucket": $S3_QUARNTINE_FILES_BUCKET}' /s3/scanfile - scan a file which is in s3 ( in scanning bucket )

GET /ruleset/ - list all loaded ruleset

GET /ruleset/{ruleset} - list all rules from a loaded rule

GET /metrics - get metrics
GET /health - get health info 
GET / - get index

```

## Installation

Automated builds of the image are available on [Registry](https://hub.docker.com/r/cloudina/hawk) and is the recommended method of installation.

```bash
docker pull hub.docker.com/cloudina/hawk:(imagetag)
```

The following image tags are available:
* `latest` - Most recent release of ClamAV with REST API

# Quick Start

Run hawk docker image:
```bash
docker run -p 9000:9999 -itd --name hawk cloudina/hawk
docker run -p 9000:9999 -v $HOME/.aws/credentials:/go/src/app/.aws/credentials:ro -itd --name hawk cloudina/hawk
```

Test that service detects common test virus signature:

**EXAMPLES**
```bash
# Request - Scanning a file from S3 , ./testsamples/request/s3filescan has config for s3
curl --data "@./testsamples/request/s3filescan" http://0.0.0.0:9000/s3/scanfile -H 'Content-Type: application/json'

# Response
{"filename":"stream","matches":[{"Rule":"Win.Test.EICAR_HDB-1","namespace":"","tags":null}],"status":"INFECTED"}%                                 

# Request - Uploading sample virus file to API
curl --data "@./testsamples/scanfiles/eicar" http://0.0.0.0:9000/scanstream -H 'Content-Type: application/json'

# Response
{"filename":"stream","matches":[{"Rule":"Win.Test.EICAR_HDB-1","namespace":"","tags":null}],"status":"INFECTED"}                           

# Request - Uploading sample clean file to API
curl --data "@./testsamples/scanfiles/hello.txt" http://0.0.0.0:9000/scanstream -H 'Content-Type: application/json'

# Response
{"filename":"stream","matches":[],"status":"CLEAN"} 
                                                                                         
```
## Networking

| Port | Description |
|-----------|-------------|
| `3310`    | ClamD Listening Port |
| `9999`    | HAWK Container Port |

## Debug
For debugging the running container
```bash
docker exec -it (whatever your container name is e.g. hawk) /bin/ash
```

## Build
For building
```bash
docker build -t (whatever your image name is e.g. hawk) .
```

## Prebuild Image
```bash
docker pull cloudina/hawk

```

## Acknowledgements

* [yarascanner](https://github.com/jheise/yarascanner)
* [clamscanner](https://github.com/ifad/clammit)

## References

* https://www.clamav.net
* https://virustotal.github.io/yara/
