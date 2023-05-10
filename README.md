# go-tui-file-project aka _puffin transfer_

## Build instructions
1. Build the binary with `go build -o puffin-transfer cmd/cli/main.go`
2. Run the binary with `./puffin-transfer`

## Project structure

### cmd/cli
- main tui for user input

### config
- configuration of static properties

### internal/firestore
- connects to firestore in GCP and uploads/downloads files
- requires service account [json file](https://console.cloud.google.com/iam-admin/serviceaccounts/details/114598002818126335278/keys?authuser=1&project=puffin-transfer&supportedpurview=project)
- stored files are saved in the [files collection](https://console.cloud.google.com/firestore/databases/-default-/data/panel/files/vbsEePbrUSzBaZpqshOP?referrer=search&authuser=1&project=puffin-transfer&supportedpurview=project) for now

### internal/fileanalyzer
- Walks a file tree and analyzes all files on different criteria. If a file violates a rule, it is transferred to the store.
- Supported criteria currently are: file size and duplicate files
