# Puffin Transfer
A reverese engineered iCloud solution for Unix systems written in Golang. Goal is to automatically rotate unused, unwanted and unnecessary files out of your system into a cloud of your choice. Current state supports rotation of identified files into a GCP firestore instance upon command. WIP - use with caution.

## Dependencies
- [Golang](https://go.dev/doc/install): ^1.20.4
- Access to cloud system of your choice, i.e. [GCP - Firestore](https://cloud.google.com)

## Build instructions
1. Retrieve the repository
  ```shell
  git clone https://github.com/OliverKlukas/puffin-transfer.git
  cd puffin-transfer
  ```
3. Build the binary with
  ```shell 
  go build -o puffin-transfer cmd/cli/main.go
  ```
3. Run the binary with 
  ```shell
  ./puffin-transfer
  ```

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
