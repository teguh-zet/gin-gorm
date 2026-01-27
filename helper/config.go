package helper

import (
	"github.com/spf13/viper"
)

type Config struct {
	// Server Configuration
	AppPort string `mapstructure:"APP_PORT"`
	GinMode string `mapstructure:"GIN_MODE"`
	Schema  string `mapstructure:"SCHEMA"`
	// Database & Security Configuration
	DB        string `mapstructure:"DB"`
	JWTSecret string `mapstructure:"JWT_SECRET"`

	// Admin Seeding Configuration
	ADMIN_EMAIL   string `mapstructure:"ADMIN_EMAIL"`
	AdminPassword string `mapstructure:"ADMIN_PASSWORD"`
	AdminName     string `mapstructure:"ADMIN_NAME"`

	// Cloudinary Configuration
	CloudinaryAPISecret string `mapstructure:"CLOUDINARY_API_SECRET"`
	CloudinaryAPIKey    string `mapstructure:"CLOUDINARY_API_KEY"`
	CloudinaryCloudName string `mapstructure:"CLOUDINARY_CLOUD_NAME"`

	ALLOW_ORIGIN string `mapstructure:"ALLOW_ORIGIN"`
	LOG_FILE     string `mapstructure:"LOG_FILE"`
	AUTO_MIGRATE string `mapstructure:"AUTO_MIGRATE"`

	//nats 
	NatsUrl   string `mapstructure:"NATS_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	// Menambahkan path lokasi file config
	viper.AddConfigPath(path)
	viper.AddConfigPath(".") // Cek juga di root folder

	// Memberitahu Viper untuk mencari file spesifik bernama ".env"
	viper.SetConfigFile(".env")

	// Membaca Environment Variables dari sistem (jika ada yang di-override)
	viper.AutomaticEnv()

	// Mulai membaca file
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// Memasukkan nilai dari file ke dalam Struct Config
	err = viper.Unmarshal(&config)
	return
}
