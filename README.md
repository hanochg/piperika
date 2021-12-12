![](assets/logo.png)


## About this plugin
This is a JFrog Pipelines CLI utility with an emphasis on the consumer experience.</br>
As a consumer of CI/CD system, you don't always want to understand the complexity but just "use" it.</br></br>
Piperika provides a single command that does all the complex work for you.</br>
All you need to do is to commit and push your code, and launch Piperika. It will start the CI on your branch, follow the build process and show you all the relevant data including your progress. 

## Installation with JFrog CLI
Since this plugin is currently not included in [JFrog CLI Plugins Registry](https://github.com/jfrog/jfrog-cli-plugins-reg), it needs to be built and installed manually. Follow these steps to install and use this plugin with JFrog CLI.
1. Make sure JFrog CLI is installed on you machine by running ```jfrog```. If it is not installed, [install](https://jfrog.com/getcli/) it.
2. Create a directory named ```plugins``` under ```~/.jfrog/``` if it does not exist already.
3. Clone this repository.
4. `cd` into the root directory of the cloned project.
5. Run ```make build``` to create the binary in the current directory.
6. Copy the binary into the ```~/.jfrog/plugins``` directory.

**Or**, do `make install` in the sources root folder.

## Building from source
To build the **piperika** binary
```shell
make build
```
To build the **piperika** binary for multiple operating systems and architectures (Mac, Linux and Windows)
```shell
make build-multi-os
```

## Usage
### What will Piperika do?
* Validate that your local commit is on the remote git server (so the CI server could reach it)
* Sync Pipelines with your branch and latest commit SHA (if needed).
* Check if there is a CI pipe that is already running with your commit, if not, it will trigger it.
* Follow up your CI run, providing the current state and progress, information about your run steps (in progress, succeed and failed steps), and displays information about tests failure.

### Commands
* `piperika run` (or just `piperika r`): </br>
It will start the Piperika magic.</br>


## Release Notes
The release notes are available [here](RELEASE.md).

## Owners
Hanoch Giner</br>
Itai Raz</br>
Omer Karjevsky