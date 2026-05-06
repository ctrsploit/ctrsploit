# e2e Test Script Developer Guide

This document provides an overview and usage instructions for the e2e test script. The script automates the process of deploying your codebase to specified remote environments, setting up Docker-based test environments, and executing test commands. Follow the details below to configure and run the script effectively.

---

## Overview

The e2e test script is designed to:

- **Verify Required Tools:** Ensure the necessary commands (such as `yq`, `git`, `docker`, `scp`, and `ssh`) are installed.
- **Package & Upload Codebase:** Bundle the project code (excluding certain directories) into a tarball and upload it to a remote host.
- **Set Up Test Environments:** Clone/update the dqd repository and launch specified containers using `docker-compose`.
- **Execute Test Commands:** Run predefined test commands on the remote host after deploying the code and loading Docker images.

The script iterates over all `e2e.yml` files found recursively in the project directory and performs actions defined within.

---

## Usage Instructions

1. **Prepare Your Configuration:**

    - Place one or more `e2e.yml` files within the project tree.
    - Update each file with the necessary `test_envs` details (i.e., remote host, command, and dqd directory).

2. **Run the Script:**

    - Execute the script from the project directory:

      ```bash
      ./e2e_test.sh
      ```

3. **Monitor the Process:**

e.g.

```shell
$ make e2e
==================================
Handling: ./pkg/hostpath/e2e.yml
----------------------------------
TEST_ENV = docker-v19.03.13-e2e
TEST_CMD = docker run --rm -v $(pwd):/root/app --env TEST_ENV=docker-v19.03.13-e2e ghcr.io/ctrsploit/ctrsploit-dev:latest go test -v -run '^TestE2E' github.com/ctrsploit/ctrsploit/pkg/hostpath
DQD_DIR  = docker/v19.03.13-e2e
Already up to date.
[+] Running 1/0
 ✔ Container docker-19-03-13-e2e-vm-1  Running                                                                                                                                                                                 0.0s 
Warning: Permanently added '[127.0.0.1]:19313' (ED25519) to the list of known hosts.
ctrsploit.tar.gz                                                                                                                                                                                  100%   12MB 207.0MB/s   00:00    
Warning: Permanently added '[127.0.0.1]:19313' (ED25519) to the list of known hosts.
Loaded image: ghcr.io/ctrsploit/ctrsploit-dev:latest
go: downloading github.com/davecgh/go-spew v1.1.1
...
=== RUN   TestE2E_WritableAccessible
=== RUN   TestE2E_WritableAccessible/docker-v19.03.13-e2e:/
=== RUN   TestE2E_WritableAccessible/docker-v19.03.13-e2e:/root/app
=== RUN   TestE2E_WritableAccessible/docker-v19.03.13-e2e:/etc/hosts
=== RUN   TestE2E_WritableAccessible/docker-v19.03.13-e2e:/etc/hostname
=== RUN   TestE2E_WritableAccessible/docker-v19.03.13-e2e:/etc/resolv.conf
--- PASS: TestE2E_WritableAccessible (0.00s)
    --- PASS: TestE2E_WritableAccessible/docker-v19.03.13-e2e:/ (0.00s)
    --- PASS: TestE2E_WritableAccessible/docker-v19.03.13-e2e:/root/app (0.00s)
    --- PASS: TestE2E_WritableAccessible/docker-v19.03.13-e2e:/etc/hosts (0.00s)
    --- PASS: TestE2E_WritableAccessible/docker-v19.03.13-e2e:/etc/hostname (0.00s)
    --- PASS: TestE2E_WritableAccessible/docker-v19.03.13-e2e:/etc/resolv.conf (0.00s)
PASS
ok      github.com/ctrsploit/ctrsploit/pkg/hostpath     0.002s
```

---

## Prerequisites

Before running the script, make sure you have the following installed and accessible on your system:

yq

```shell
go install github.com/mikefarah/yq/v4@latest
```

dqd

```shell
git clone https://github.com/ctrsploit/dqd.git /tmp/dqd
/tmp/dqd/script/install_ssh_config.sh
```


---

## Script Structure

The script consists of several key sections:

### Function Descriptions

#### 1. `upload_codebase()`

- **Purpose:**  
  Packages the project codebase (excluding certain directories such as `.idea` and `bin`) into a tarball (`ctrsploit.tar.gz`) and uploads it to the remote host using `scp`.

- **Details:**
    - Determines the location of the project directory (assumed to be two levels up from the script's directory).
    - Uses `tar` to create a compressed archive.
    - Transfers the archive to the remote host.

#### 2. `startup_testEnv()`

- **Purpose:**  
  Prepares the test environment on the remote host by managing Docker containers.

- **Details:**
    - Checks if a local directory (`/tmp/dqd`) exists; if not, clones the dqd repository.
    - If the directory exists, it pulls the latest changes.
    - Uses `docker compose` (by combining `docker-compose.yml` and `docker-compose.kvm.yml`) to start the containers in detached mode.

#### 3. `test_cmd()`

- **Purpose:**  
  Executes the specified test command on the remote host.

- **Details:**
    - Connects via SSH, extracts the uploaded codebase tarball into a directory named `ctrsploit`, and loads a Docker image from `/root/dev-image.tar`.
    - Runs the test command that is defined per environment.

### Processing YAML Configuration Files

The script searches for all files named `e2e.yml` recursively in the current directory:

- For each `e2e.yml` file found, it determines the number of test environments provided under the `test_envs` key.
- For each test environment, it extracts:
    - **Environment Name (`ENV_NAME`):** Expected to be the remote host's address or identifier.
    - **Test Command (`ENV_CMD`):** The command to run on the remote host for executing tests.
    - **dqd Directory (`DQD_DIR`):** Specifies the subdirectory in the dqd repository to configure the environment.
- It then sequentially calls:
    - `startup_testEnv` to set up the Docker environment.
    - `upload_codebase` to transfer the project.
    - `test_cmd` to execute the command on the remote host.

---

## Configuration

### `e2e.yml` File Format

Each `e2e.yml` should define one or more test environments. Below is an example of the expected format:

```yaml
name: hostpath
test_envs:
  - name: docker-v19.03.13-e2e
    kind: dqd
    dqd_dir: docker/v19.03.13-e2e
    cmd: docker run --rm -v $(pwd):/root/app --env TEST_ENV=docker-v19.03.13-e2e ghcr.io/ctrsploit/ctrsploit-dev:latest
      go test -v -run '^TestE2E' github.com/ctrsploit/ctrsploit/pkg/hostpath
    # optional
    stop_flag: ok*github.com/ctrsploit/ctrsploit/pkg/hostpath
```

---

## Troubleshooting

--- 

Happy testing!
