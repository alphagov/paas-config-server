package config_test

import (
	. "config_server/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

var _ = Describe("ParseConfig", func() {

	Context("Config file does not exist", func() {
		It("should return an error", func() {
			_, err := ParseConfig("non-existent-file.yml")
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Config file exists", func() {

		var configFile *os.File

		BeforeEach(func() {
			configFile, _ = ioutil.TempFile(os.TempDir(), "server-config-")
		})

		AfterEach(func() {
			os.Remove(configFile.Name())
		})

		Context("has invalid content", func() {
			It("should return an error", func() {
				configFile.WriteString("garbage")
				_, err := ParseConfig(configFile.Name())
				Expect(err).ToNot(BeNil())
			})
		})

		Context("has valid content", func() {
			It("should return ServerConfig", func() {
				configFile.WriteString(`
{
   "port":9000,
   "certificate_file_path":"/path/to/cert",
   "private_key_file_path":"/path/to/key",
   "ca_certificate_file_path" : "/path/to/ca/cert",
   "ca_private_key_file_path": "/path/to/ca/private/key",
   "database":{
      "adapter":"postgres",
      "user":"uword",
      "password":"pword",
      "host":"http://www.yahoo.com",
      "port":4300,
      "db_name":"db",
      "connection_options":{
         "max_open_connections":12,
         "max_idle_connections":25
      }
   }
}
`)
				serverConfig, err := ParseConfig(configFile.Name())
				Expect(err).To(BeNil())

				Expect(serverConfig).ToNot(BeNil())
				Expect(serverConfig.Port).To(Equal(9000))
				Expect(serverConfig.CertificateFilePath).To(Equal("/path/to/cert"))
				Expect(serverConfig.PrivateKeyFilePath).To(Equal("/path/to/key"))
				Expect(serverConfig.CACertificatePath).To(Equal("/path/to/ca/cert"))
				Expect(serverConfig.CAPrivateKey).To(Equal("/path/to/ca/private/key"))
				Expect(serverConfig.PrivateKeyFilePath).To(Equal("/path/to/key"))
				Expect(serverConfig.Database).ToNot(BeNil())
				Expect(serverConfig.Database.Adapter).To(Equal("postgres"))
				Expect(serverConfig.Database.User).To(Equal("uword"))
				Expect(serverConfig.Database.Password).To(Equal("pword"))
				Expect(serverConfig.Database.Host).To(Equal("http://www.yahoo.com"))
				Expect(serverConfig.Database.Port).To(Equal(4300))
				Expect(serverConfig.Database.Name).To(Equal("db"))
				Expect(serverConfig.Database.ConnectionOptions).ToNot(BeNil())
				Expect(serverConfig.Database.ConnectionOptions.MaxOpenConnections).To(Equal(12))
				Expect(serverConfig.Database.ConnectionOptions.MaxIdleConnections).To(Equal(25))
			})
		})

		Context("has missing keys", func() {
			It("should error hen certificate_file_path is missing", func() {
				configFile.WriteString(`
{
   "port":9000,
   "private_key_file_path":"/path/to/key",
   "database":{
      "adapter":"postgres",
      "user":"uword",
      "password":"pword",
      "host":"http://www.yahoo.com",
      "port":4300,
      "db_name":"db",
      "connection_options":{
         "max_open_connections":12,
         "max_idle_connections":25
      }
   }
}
`)
				_, err := ParseConfig(configFile.Name())
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Certificate file path and key file path should be defined"))
			})

			It("should error hen private_key_file_path is missing", func() {
				configFile.WriteString(`
{
   "port":9000,
   "certificate_file_path":"/path/to/cert",
   "database":{
      "adapter":"postgres",
      "user":"uword",
      "password":"pword",
      "host":"http://www.yahoo.com",
      "port":4300,
      "db_name":"db",
      "connection_options":{
         "max_open_connections":12,
         "max_idle_connections":25
      }
   }
}
`)
				_, err := ParseConfig(configFile.Name())
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Certificate file path and key file path should be defined"))
			})
		})
	})
})
