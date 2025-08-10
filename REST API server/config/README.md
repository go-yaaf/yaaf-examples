# Config
This folder includes the service configuration utility which provides application-wise accessor to configuration variables.
The service configuration (based on environment variables) is provided by the orchestration service (e.g. Kubernetes or Docker manifests)
and via the configuration one can select the concrete implementation of each middleware facade

## Configuration Variables
The configuration service is derived from the `yaaf-common/config/BaseConfig` utility which includes some common variables
like LOG LEVEL.
On top of the common variables, the developer should add additional variables and accessor methods to this configuration utility.
All the configuration variables must be listed also in the main application `README.md` file for visibility.


