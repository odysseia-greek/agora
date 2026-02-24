package aristoteles

import (
	"bytes"
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/odysseia-greek/agora/aristoteles/models"
)

type Client interface {
	Query() Query
	Document() Document
	Index() Index
	Builder() Builder
	Health() Health
	Access() Access
	Policy() Policy
}

type Query interface {
	Match(index string, request map[string]interface{}) (*models.Response, error)
	MatchWithContext(ctx context.Context, index string, request map[string]interface{}) (*models.Response, error)
	MatchWithSort(index, mode, sort string, size int, request map[string]interface{}) (*models.Response, error)
	MatchWithSortWithContext(ctx context.Context, index, mode, sort string, size int, request map[string]interface{}) (*models.Response, error)
	MatchWithScroll(index string, request map[string]interface{}) (*models.Response, error)
	MatchWithScrollWithContext(ctx context.Context, index string, request map[string]interface{}) (*models.Response, error)
	MatchAggregate(index string, request map[string]interface{}) (*models.Aggregations, error)
	MatchAggregateWithContext(ctx context.Context, index string, request map[string]interface{}) (*models.Aggregations, error)
	MatchRaw(index string, request map[string]interface{}) ([]byte, error)
	MatchRawWithContext(ctx context.Context, index string, request map[string]interface{}) ([]byte, error)
	CountRaw(ctx context.Context, index string, request map[string]interface{}) (*models.CountResponse, error)
	GetById(ctx context.Context, index, id string) (*models.DirectResponse, error)
}

type Document interface {
	Create(index string, body []byte) (*models.CreateResult, error)
	CreateWithContext(ctx context.Context, index string, body []byte) (*models.CreateResult, error)
	Update(index, id string, body []byte) (*models.CreateResult, error)
	UpdateWithContext(ctx context.Context, index, id string, body []byte) (*models.CreateResult, error)
	AddItemToDocument(index, id, body, paramName string) (*models.CreateResult, error)
	AddItemToDocumentWithContext(ctx context.Context, index, id, body, paramName string) (*models.CreateResult, error)
	CreateWithId(index, documentId string, body []byte) (*models.CreateResult, error)
	CreateWithIdWithContext(ctx context.Context, index, documentId string, body []byte) (*models.CreateResult, error)
	CreateWithIdAndFirstItem(index, documentId, body, paramName string) (*models.CreateResult, error)
	CreateWithIdAndFirstItemWithContext(ctx context.Context, index, documentId, body, paramName string) (*models.CreateResult, error)
	Bulk(buf bytes.Buffer, index string) (*BulkResponse, error)
	BulkWithContext(ctx context.Context, buf bytes.Buffer, index string) (*BulkResponse, error)
}

type Index interface {
	CreateDocument(index string, body []byte) (*models.CreateResult, error)
	CreateDocumentWithContext(ctx context.Context, index string, body []byte) (*models.CreateResult, error)
	Create(index string, request map[string]interface{}) (*models.IndexCreateResult, error)
	CreateWithContext(ctx context.Context, index string, request map[string]interface{}) (*models.IndexCreateResult, error)
	CreateWithAlias(indexName string, request map[string]interface{}) (*models.IndexCreateResult, error)
	CreateWithAliasWithContext(ctx context.Context, indexName string, request map[string]interface{}) (*models.IndexCreateResult, error)
	Delete(index string) (bool, error)
	DeleteWithContext(ctx context.Context, index string) (bool, error)
	IndexExists(index string) (bool, *models.IndexInfo, error)
	IndexExistsWithContext(ctx context.Context, index string) (bool, *models.IndexInfo, error)
}

type Policy interface {
	CreatePolicyWithRollOver(name, maxAge, phase string) (*models.IndexCreateResult, error)
	CreatePolicyWithRollOverWithContext(ctx context.Context, name, maxAge, phase string) (*models.IndexCreateResult, error)
	CreatePolicy(name, phase string) (*models.IndexCreateResult, error)
	CreatePolicyWithContext(ctx context.Context, name, phase string) (*models.IndexCreateResult, error)
}

type Builder interface {
	MatchQuery(term, queryWord string) map[string]interface{}
	MatchAll() map[string]interface{}
	MultipleMatch(mappedFields []map[string]string) map[string]interface{}
	MultiMatchWithGram(queryWord, field string) map[string]interface{}
	MatchPhrasePrefixed(queryWord, field string) map[string]interface{}
	Aggregate(aggregate, field string) map[string]interface{}
	FilteredAggregate(term, queryWord, aggregate, field string) map[string]interface{}
	SearchAsYouTypeIndex(searchWord string) map[string]interface{}
	Index() map[string]interface{}
	TextIndex(policyName string) map[string]interface{}
	GrammarIndex(policyName string) map[string]interface{}
}

type Health interface {
	Check(ticks, tick time.Duration) bool
	Info() (elasticHealth models.DatabaseHealth)
}

type Access interface {
	CreateRole(name string, roleRequest models.CreateRoleRequest) (bool, error)
	CreateRoleWithContext(ctx context.Context, name string, roleRequest models.CreateRoleRequest) (bool, error)
	CreateUser(name string, userCreation models.CreateUserRequest) (bool, error)
	CreateUserWithContext(ctx context.Context, name string, userCreation models.CreateUserRequest) (bool, error)
	ListUsers() ([]string, error)
	ListUsersWithContext(ctx context.Context) ([]string, error)
	DeleteUser(name string) (bool, error)
	DeleteUserWithContext(ctx context.Context, name string) (bool, error)
}

type Elastic struct {
	document *DocumentImpl
	query    *QueryImpl
	index    *IndexImpl
	builder  *BuilderImpl
	health   *HealthImpl
	access   *AccessImpl
	policy   *PolicyImpl
}

func NewClient(config models.Config) (Client, error) {
	var err error
	var esClient *elasticsearch.Client
	if config.ElasticCERT != "" {
		esClient, err = createWithTLS(config)
		if err != nil {
			return nil, err
		}
	} else {
		esClient, err = create(config)
		if err != nil {
			return nil, err
		}
	}

	query, err := NewQueryImpl(esClient)
	if err != nil {
		return nil, err
	}

	index, err := NewIndexImpl(esClient)
	if err != nil {
		return nil, err
	}

	health, err := NewHealthImpl(esClient)
	if err != nil {
		return nil, err
	}

	access, err := NewAccessImpl(esClient)
	if err != nil {
		return nil, err
	}

	document, err := NewDocumentImpl(esClient)
	if err != nil {
		return nil, err
	}

	policy, err := NewPolicyImpl(esClient)
	if err != nil {
		return nil, err
	}

	builder := NewBuilderImpl()

	es := &Elastic{query: query, index: index, builder: builder, health: health, access: access, document: document, policy: policy}

	return es, nil
}

func NewMockClient(fixtureFiles interface{}, statusCode int) (Client, error) {
	var files []string
	switch t := fixtureFiles.(type) {
	case string:
		files = []string{t}
	case []string:
		files = t
	case [][]byte:
		for _, data := range t {

			tempFile, err := os.CreateTemp("", "mockclient-*.json")
			if err != nil {
				return nil, err
			}
			defer tempFile.Close()

			if _, err := tempFile.Write(data); err != nil {
				return nil, err
			}
			AddRawFixture(tempFile.Name(), data)
			files = append(files, tempFile.Name())
		}
	default:
		return nil, errors.New("unsupported fixtureFiles type")
	}

	esClient, err := CreateMockClient(files, statusCode)
	if err != nil {
		return nil, err
	}

	query, err := NewQueryImpl(esClient)
	if err != nil {
		return nil, err
	}

	index, err := NewIndexImpl(esClient)
	if err != nil {
		return nil, err
	}

	health, err := NewHealthImpl(esClient)
	if err != nil {
		return nil, err
	}

	access, err := NewAccessImpl(esClient)
	if err != nil {
		return nil, err
	}

	document, err := NewDocumentImpl(esClient)
	if err != nil {
		return nil, err
	}

	policy, err := NewPolicyImpl(esClient)
	if err != nil {
		return nil, err
	}

	builder := NewBuilderImpl()

	es := &Elastic{query: query, index: index, builder: builder, health: health, access: access, document: document, policy: policy}

	return es, nil
}

func create(config models.Config) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Username:  config.Username,
		Password:  config.Password,
		Addresses: []string{config.Service},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return es, nil
}

func createWithTLS(config models.Config) (*elasticsearch.Client, error) {
	caCert := []byte(config.ElasticCERT)

	// --> Clone the default HTTP transport

	tp := http.DefaultTransport.(*http.Transport).Clone()

	// --> Initialize the set of root certificate authorities
	//
	var err error

	if tp.TLSClientConfig.RootCAs, err = x509.SystemCertPool(); err != nil {
		return nil, fmt.Errorf("problem adding system CA: %w", err)
	}

	// --> Add the custom certificate authority
	//
	if ok := tp.TLSClientConfig.RootCAs.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("problem adding CA from certificate data")
	}

	cfg := elasticsearch.Config{
		Username:  config.Username,
		Password:  config.Password,
		Addresses: []string{config.Service},
		Transport: tp,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return es, nil
}

func (e *Elastic) Query() Query {
	if e == nil {
		return nil
	}
	return e.query
}

func (e *Elastic) Document() Document {
	if e == nil {
		return nil
	}
	return e.document
}

func (e *Elastic) Index() Index {
	if e == nil {
		return nil
	}
	return e.index
}

func (e *Elastic) Health() Health {
	if e == nil {
		return nil
	}
	return e.health
}

func (e *Elastic) Builder() Builder {
	if e == nil {
		return nil
	}
	return e.builder
}

func (e *Elastic) Access() Access {
	if e == nil {
		return nil
	}
	return e.access
}

func (e *Elastic) Policy() Policy {
	if e == nil {
		return nil
	}
	return e.policy
}
