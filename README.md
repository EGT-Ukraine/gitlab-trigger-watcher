gitlab-trigger-watcher
---

Application to check Gitlab trigger's pipeline status

### Build and run
```bash
make build
./gitlab-trigger-watcher -h
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
  * `--variables` - (optional) custom variables for the project (format example: `variable1:value,variable2:value`)
  
### Run:

```bash
./gitlab-trigger-watcher --privateToken ${PRIVATE_TOKEN} --token ${TOKEN} --host gitlab.egt.com --projectID 123 run
```