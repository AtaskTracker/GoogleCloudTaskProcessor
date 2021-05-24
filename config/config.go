package config

type Config struct {
	Storage Storage
}

type Storage struct {
	Bucket string
	Url string
}
