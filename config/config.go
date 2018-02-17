package config

import (
	"log"
	"time"

	ini "gopkg.in/ini.v1"
)

// Config bundles all configuration-options
// available in the config.ini
var Config struct {
	Server     Server           `ini:"server"`
	DB         DBConfig         `ini:"db"`
	Grades     Grades           `ini:"grades"`
	Auth       Authentification `ini:"auth"`
	Mail       Mail             `ini:"mail"`
	Sessioner  Sessioner        `ini:"sessioner"`
	FileSystem FileSystem       `ini:"filesystem"`
}

// FileSystem information
type FileSystem struct {
	// MDStoragePath holds the directory in which this server stores document-files.
	// Make sure is exists before executing the web command
	MDStoragePath string `ini:"mdstoragepath"`
	// AssetsStoragePath holds the directory in which this server stores
	// document-assets-directories. Make sure is exists before executing the web command.
	AssetsStoragePath string `ini:"assetsstoragepath"`
}

// Sessioner information
type Sessioner struct {
	// StoragePath holds the directory in which this server stores files
	// for session-management. Make sure is exists before executing the web command
	StoragePath string `ini:"storagepath"`
}

// Server information
type Server struct {
	// Address holds the domain, which is used to host this service. It is needed
	// to create the confirm-links.
	Address string `ini:"address"`
	// IP holds the ip used to host this service. Use "localhost" or "0.0.0.0"
	IP string `ini:"ip"`
	// Port holds the port used to host this service
	Port int `ini:"port"`
	// DevMode is used to toggle between development and production use.
	// In devmode more log-data is produced and templates are re-rendered
	// every single time. Also panics will be shown in the client's browser
	DevMode bool `ini:"devmode"`
	// ResetMonth holds the month (1 = January; 12 = December), when the server
	// resets itself. It will keep the admin-user and the tos and about documents
	ResetMonth int `ini:"resetmonth"`
	// Advent must not be defined in the config-file. It is set to the current
	// month in development mode or to time.December in production mode
	Advent time.Month
}

// Mail information
type Mail struct {
	// Sender holds the e-mail-address displayed to the customer, when this service
	// sends a e-mail
	Sender string `ini:"sender"`
	// Address holds the url/ip of your e-mail-server
	Address string `ini:"address"`
	// Port holds the port of your e-mail-server for sending e-mails via smtp
	Port int `ini:"port"`
	// Username holds the username of a valid account on your e-mail-server
	Username string `ini:"username"`
	// Password holds the password of a valid account on your e-mail-server
	Password string `ini:"password"`
}

// Authentification information
type Authentification struct {
	// MinPwLength holds the minimum length a user's password may have
	MinPwLength uint `ini:"minpwlength"`
	// MaxPwLength holds the maximum length a user's password may have
	MaxPwLength uint `ini:"maxpwlength"`
	// AdminMail holds the server-admin's mail-address. It will be available
	// right when the server starts. The admin can upload new documents
	AdminMail string `ini:"adminmail"`
	// AdminHash is the admin's hash converted to a bcrypt-hash. You can use this
	// (https://bcrypt-generator.com/) website to generate such a hash
	AdminHash string `ini:"adminhash"`
}

// Grades boundaries
type Grades struct {
	// Min holds the minimum grade a user of this service may be in
	Min uint `ini:"min"`
	// Max holds the maximum grade a user of this service may be in
	Max uint `ini:"max"`
}

// DBConfig stores database-related information
// for connection and authentication
type DBConfig struct {
	// Host holds the url/ip to your database-server
	Host string `ini:"host"`
	// Port holds the port on which your database-server is hosted
	Port uint `ini:"port"`
	// User hods the database-server's user. The user must have access to
	// the database you specify in DBConfig.Name
	User string `ini:"user"`
	// Password holds the password for the user stored in DBConfig.User.
	Password string `ini:"password"`
	// Name name holds the name of your database
	Name string `ini:"name"`
}

// Load loads the provided config file into memory.
func Load(path string) {
	f, err := ini.Load(path)
	if err != nil {
		log.Fatal(err)
	}
	err = f.MapTo(&Config)
	if err != nil {
		log.Fatal(err)
	}
}
