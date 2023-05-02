# go-tui-file-project aka *puffin transfer*

## Project structure
### cmd/cli
- main tui for user input

### config
- configuration of static properties

#### internal/fileservice
- connects to firestore in GCP and uploads/downloads files
- requires service account [json file](https://console.cloud.google.com/iam-admin/serviceaccounts/details/114598002818126335278/keys?authuser=1&project=puffin-transfer&supportedpurview=project)
- stored files are saved in the [files collection](https://console.cloud.google.com/firestore/databases/-default-/data/panel/files/vbsEePbrUSzBaZpqshOP?referrer=search&authuser=1&project=puffin-transfer&supportedpurview=project) for now
