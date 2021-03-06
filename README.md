gitlab-trigger-watcher [![Build Status](https://travis-ci.org/EGT-Ukraine/gitlab-trigger-watcher.svg?branch=master)](https://travis-ci.org/EGT-Ukraine/gitlab-trigger-watcher)
---

Application to run and check until the end of Gitlab trigger's job

### Build and run
```bash
make build
./gtw -h
```

options:
  * `--privateToken` - user's personal private token;
  * `--token` - project token;
  * `--projectID` - integer ID of your project;
  * `--skipVerifyTLS` - (optional) skip to verify tls certificate;
  * `--host` - (optional) set your custom Gitlab's host;
  * `--schema` - (optional. default: https) set http or https connection type;
  * `--ref` - (optional) branch for the project. (default: master);
  * `--urlPrefix` - (optional) if you are use some prefix for your Gitlab (final URL will be looked like: /prefix/api/v4/...);
  * `--variables` - (optional) custom variables for the project
  
### Run:

```bash
./gtw --privateToken ${PRIVATE_TOKEN} --token ${TOKEN} --host gitlab.egt.com --projectID 123 --variables KEY1:VALUE1 --variables KEY2:VALUE2 run
```