package aristoteles

import (
	"fmt"
)

type BuilderImpl struct {
}

func NewBuilderImpl() *BuilderImpl {
	return &BuilderImpl{}
}

func (b *BuilderImpl) MatchQuery(term, queryWord string) map[string]interface{} {
	return map[string]interface{}{
		"query": map[string]interface{}{
			"match_phrase": map[string]string{
				term: queryWord,
			},
		},
	}
}

func (b *BuilderImpl) MatchAll() map[string]interface{} {
	return map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
}

func (b *BuilderImpl) MultipleMatch(mappedFields []map[string]string) map[string]interface{} {
	var must []map[string]interface{}

	for _, mappedField := range mappedFields {
		for key, value := range mappedField {
			matchItem := map[string]interface{}{
				"match": map[string]interface{}{
					key: value,
				},
			}
			must = append(must, matchItem)
		}
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
	}

	return query
}

func (b *BuilderImpl) MultiMatchWithGram(queryWord, field string) map[string]interface{} {
	return map[string]interface{}{
		"size": 15,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": queryWord,
				"type":  "bool_prefix",
				"fields": [3]string{
					field, fmt.Sprintf("%s._2gram", field), fmt.Sprintf("%s._3gram", field),
				},
			},
		},
	}
}

func (b *BuilderImpl) MatchPhrasePrefixed(queryWord, field string) map[string]interface{} {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_phrase_prefix": map[string]interface{}{
				field: queryWord,
			},
		},
	}

	return query
}

func (b *BuilderImpl) Aggregate(aggregate, field string) map[string]interface{} {
	return map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			aggregate: map[string]interface{}{
				"terms": map[string]interface{}{
					"field": field,
					"size":  500,
				},
			},
		},
	}
}

func (b *BuilderImpl) FilteredAggregate(term, queryWord, aggregate, field string) map[string]interface{} {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_phrase": map[string]interface{}{
				term: queryWord,
			},
		},
		"size": 0,
		"aggs": map[string]interface{}{
			aggregate: map[string]interface{}{
				"terms": map[string]interface{}{
					"field": field,
					"size":  500,
				},
			},
		},
	}

	return query
}

func (b *BuilderImpl) SearchAsYouTypeIndex(searchWord string) map[string]interface{} {
	return map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				searchWord: map[string]interface{}{
					"type": "search_as_you_type",
				},
			},
		},
	}
}

func (b *BuilderImpl) TextIndex(policyName string) map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 1,
				"lifecycle.name":     policyName,
			},
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"greek_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"greek_stop",
							"greek_stemmer",
						},
					},
				},
				"filter": map[string]interface{}{
					"greek_stop": map[string]interface{}{
						"type":      "stop",
						"stopwords": "_greek_",
					},
					"greek_stemmer": map[string]interface{}{
						"type":     "stemmer",
						"language": "greek",
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"author": map[string]interface{}{
					"type": "keyword",
				},
				"greek": map[string]interface{}{
					"type":     "text",
					"analyzer": "greek_analyzer",
				},
				"translations": map[string]interface{}{
					"type": "text",
				},
				"book": map[string]interface{}{
					"type": "integer",
				},
				"chapter": map[string]interface{}{
					"type": "integer",
				},
				"section": map[string]interface{}{
					"type": "integer",
				},
				"perseusTextLink": map[string]interface{}{
					"type": "keyword",
				},
			},
		},
	}
}

func (b *BuilderImpl) QuizIndex(policyName string) map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   2, // Adjust as needed based on cluster size
				"number_of_replicas": 1,
				"lifecycle.name":     policyName, // Add this line to associate the policy
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"method": map[string]interface{}{
					"type": "keyword",
				},
				"category": map[string]interface{}{
					"type": "keyword",
				},
				"greek": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"translation": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"chapter": map[string]interface{}{
					"type": "integer",
				},
				// Add additional fields here if needed
			},
		},
	}
}

func (b *BuilderImpl) GrammarIndex(policyName string) map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 1,
				"lifecycle.name":     policyName, // Add this line to associate the policy
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"declension": map[string]interface{}{
					"type": "keyword",
				},
				"ruleName": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"searchTerm": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
			},
		},
	}
}

func (b *BuilderImpl) DictionaryIndex(min, max int, policyName string) map[string]interface{} {
	nGramDiff := max - min
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   2,
				"number_of_replicas": 2,
				"max_ngram_diff":     nGramDiff,
				"lifecycle.name":     policyName, // Add this line to associate the policy
			},
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"greek_analyzer": map[string]interface{}{
						"tokenizer": "greek_tokenizer",
					},
				},
				"tokenizer": map[string]interface{}{
					"greek_tokenizer": map[string]interface{}{
						"type":        "ngram",
						"min_gram":    min,
						"max_gram":    max,
						"token_chars": []string{"letter"},
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"greek": map[string]interface{}{
					"type":     "text",
					"analyzer": "greek_analyzer",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"english": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"dutch": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
			},
		},
	}
}

func (b *BuilderImpl) CreateTraceIndexMapping(policyName string) map[string]interface{} {
	return map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"items": map[string]interface{}{
					"type": "nested", // Use "nested" type for arrays of complex objects
				},
				"isActive": map[string]interface{}{
					"type": "boolean",
				},
				"timeStarted": map[string]interface{}{
					"type":   "date",
					"format": "yyyy-MM-dd'T'HH:mm:ss.SSS",
				},
				"timeEnded": map[string]interface{}{
					"type":       "date",
					"format":     "yyyy-MM-dd'T'HH:mm:ss.SSS",
					"null_value": "1970-01-01T00:00:00.000",
				},
				"totalTime": map[string]interface{}{
					"type": "long",
				},
				"responseCode": map[string]interface{}{
					"type": "short",
				},
				// Add additional fields here if needed
			},
		},
		"settings": map[string]interface{}{
			"index.lifecycle.name":                   policyName,
			"index.lifecycle.rollover_alias":         "trace",
			"index.lifecycle.parse_origination_date": true,
		},
	}
}

func (b *BuilderImpl) Index() map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 1,
			},
		},
	}
}
