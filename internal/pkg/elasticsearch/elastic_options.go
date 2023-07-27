package elasticsearch

import (
	"github.com/iancoleman/strcase"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetTypeNameByT[ElasticOptions]())

type ElasticOptions struct {
	URL string `mapstructure:"url"`
}

func provideConfig(environment environemnt.Environment) (*ElasticOptions, error) {
	return config.BindConfigKey[*ElasticOptions](optionName, environment)
}
