# Piperika - spice up your CI

## About this plugin
The plugin allows you to get all Artifactory logs with a single click.<br>
You have the ability to `cat` and `tail -f` any log on an Artifactory node.<br>

## Installation with JFrog CLI
Since this plugin is currently not included in [JFrog CLI Plugins Registry](https://github.com/jfrog/jfrog-cli-plugins-reg), it needs to be built and installed manually. Follow these steps to install and use this plugin with JFrog CLI.
1. Make sure JFrog CLI is installed on you machine by running ```jfrog```. If it is not installed, [install](https://jfrog.com/getcli/) it.
2. Create a directory named ```plugins``` under ```~/.jfrog/``` if it does not exist already.
3. Clone this repository.
4. `cd` into the root directory of the cloned project.
5. Run ```make build``` to create the binary in the current directory.
6. Copy the binary into the ```~/.jfrog/plugins``` directory.

## Usage
### Commands
* logs
    - Arguments:
        - server_id - JFrog CLI Artifactory server id.
        - node_id - Selected Artifactory node id.
        - log_name - Selected Artifactory log name.
    - Flags:
        - i: Open interactive menu **[Default: false]**
        - f: Show the log and keep following for changes **[Default: false]**
    - Example:
    ```
  $ jfrog forest logs local-arti 2368364e2c78 console.log -f | grep INFO
  2020-12-06T19:21:52.549Z [jfac ] [INFO ] [6469d8c8e2ece130] [a.s.b.AccessServerRegistrar:73] [pool-26-thread-1    ] - [ACCESS BOOTSTRAP] JFrog Access registrar finished.
  2020-12-06T19:21:52.612Z [jfac ] [INFO ] [7ccdb881f0258729] [s.r.NodeRegistryServiceImpl:68] [27.0.0.1-8040-exec-8] - Cluster join: Successfully joined jfevt@01eqtrgsxaztsq1yq0a9s60289 with node id a15e67cc9bed
  2020-12-06T19:21:52.622Z [jfevt] [INFO ] [152a442b8f87bacc] [access_join.go:58             ] [main                ] - Cluster join: Successfully joined the cluster [application]
  2020-12-06T19:21:52.624Z [jfevt] [INFO ] [152a442b8f87bacc] [access_join.go:58             ] [main                ] - Executing Router register at: localhost:8046 [application]
  ```
    ```
  $ jfrog forest logs -i
  Select JFrog CLI server id
  ✔ local-arti
  Select node id
  ✔ 2368364e2c78
  Select log name
  ✔ console.log
  2020-12-06T19:21:52.549Z [jfac ] [INFO ] [6469d8c8e2ece130] [a.s.b.AccessServerRegistrar:73] [pool-26-thread-1    ] - [ACCESS BOOTSTRAP] JFrog Access registrar finished.
  2020-12-06T19:21:52.612Z [jfac ] [INFO ] [7ccdb881f0258729] [s.r.NodeRegistryServiceImpl:68] [27.0.0.1-8040-exec-8] - Cluster join: Successfully joined jfevt@01eqtrgsxaztsq1yq0a9s60289 with node id a15e67cc9bed
  2020-12-06T19:21:52.622Z [jfevt] [INFO ] [152a442b8f87bacc] [access_join.go:58             ] [main                ] - Cluster join: Successfully joined the cluster [application]
  2020-12-06T19:21:52.624Z [jfevt] [INFO ] [152a442b8f87bacc] [access_join.go:58             ] [main                ] - Executing Router register at: localhost:8046 [application]
  ```

## Additional info
- Admin permissions are required.
- If you get an argument wrong, the CLI will suggest the correct value.
  <br>For example:
```
$ jfrog forest logs local-artii 2368364e2c78 console.log
[Error] server id not found [local-artii], consider using one of the following server id values [remote-arti,local-arti]
```
## Release Notes
The release notes are available [here](RELEASE.md).
