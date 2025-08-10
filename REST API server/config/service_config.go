package config

import (
	bc "github.com/go-yaaf/yaaf-common/config"
	"sync"
)

const (
	CfgRunAsJob       = "RUN_AS_JOB"       // Run this service as a scheduled job to execute maintenance tasks
	CfgLogJsonFormat  = "LOG_JSON_FORMAT"  // Enable Json log format
	CfgDatabaseUri    = "DATABASE_URI"     // Configuration database URI
	CfgDataCacheUri   = "DATACACHE_URI"    // Distributed cache middleware URI
	CfgFileStorageUri = "FILE_STORAGE_URI" // File storage location URI
	CfgExposeHttpPort = "EXPOSE_HTTP_PORT" // Port number to expose HTTP REST API endpoint
	CfgInitialAdmin   = "INIT_ADMIN_EMAIL" // On system startup, set the initial administrator email if not exists
	CfgMailRelayUri   = "MAIL_RELAY_URI"   // Mail Relay URI
	CfgMailRelayUsr   = "MAIL_RELAY_USR"   // Mail Relay User
	CfgMailRelayPwd   = "MAIL_RELAY_PWD"   // Mail Relay Password
	CfgMailRelayTls   = "MAIL_RELAY_TLS"   // Mail Relay TLS flag

)

type ServiceConfig struct {
	bc.BaseConfig
}

var cfg *ServiceConfig
var initOnce sync.Once

func NewConfig() *ServiceConfig {
	c := &ServiceConfig{
		BaseConfig: *bc.Get(),
	}
	c.AddConfigVar(CfgRunAsJob, "false")
	c.AddConfigVar(CfgLogJsonFormat, "false")
	c.AddConfigVar(CfgDatabaseUri, "")
	c.AddConfigVar(CfgDataCacheUri, "")
	c.AddConfigVar(CfgFileStorageUri, "")
	c.AddConfigVar(CfgExposeHttpPort, "8080")
	c.AddConfigVar(CfgInitialAdmin, "")
	c.AddConfigVar(CfgMailRelayUri, "")
	c.AddConfigVar(CfgMailRelayUsr, "")
	c.AddConfigVar(CfgMailRelayPwd, "")
	c.AddConfigVar(CfgMailRelayTls, "false")
	return c
}

// GetConfig creates or gets the configuration singleton.
func GetConfig() *ServiceConfig {
	initOnce.Do(func() {
		cfg = NewConfig()
		cfg.ScanEnvVariables()
	})
	return cfg
}

// ServerPort returns the http port for the REST API
func (c *ServiceConfig) ServerPort() (result int) {
	return c.GetIntParamValueOrDefault(CfgExposeHttpPort, 0)
}

// RunAsJob returns json job scheduler flag
func (c *ServiceConfig) RunAsJob() bool {
	return c.GetBoolParamValueOrDefault(CfgRunAsJob, false)
}

// EnableLogJsonFormat returns json log format flag
func (c *ServiceConfig) EnableLogJsonFormat() bool {
	return c.GetBoolParamValueOrDefault(CfgLogJsonFormat, false)
}

// DatabaseUri returns the database URI
func (c *ServiceConfig) DatabaseUri() string {
	return c.GetStringParamValueOrDefault(CfgDatabaseUri, "")
}

// DataCacheUri returns the distributed cache middleware URI
func (c *ServiceConfig) DataCacheUri() string {
	return c.GetStringParamValueOrDefault(CfgDataCacheUri, "")
}

// FileStorageUri returns the file storage location URI
func (c *ServiceConfig) FileStorageUri() string {
	return c.GetStringParamValueOrDefault(CfgFileStorageUri, "")
}

// InitialAdminEmail returns the initial administrator email if not exists
func (c *ServiceConfig) InitialAdminEmail() string {
	return c.GetStringParamValueOrDefault(CfgInitialAdmin, "admin@org.io")
}

// MailRelayUri returns the mail relay URI
func (c *ServiceConfig) MailRelayUri() string {
	return c.GetStringParamValueOrDefault(CfgMailRelayUri, "smtp://smtp.gmail.com:587")
}

// MailRelayUsr returns the mail relay user
func (c *ServiceConfig) MailRelayUsr() string {
	return c.GetStringParamValueOrDefault(CfgMailRelayUsr, "")
}

// MailRelayPwd returns the mail relay password
func (c *ServiceConfig) MailRelayPwd() string {
	return c.GetStringParamValueOrDefault(CfgMailRelayPwd, "")
}

// MailRelayTls returns the mail relay TLS flag
func (c *ServiceConfig) MailRelayTls() bool {
	return c.GetBoolParamValueOrDefault(CfgMailRelayTls, true)
}
