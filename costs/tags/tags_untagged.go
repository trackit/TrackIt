//   Copyright 2017-2018 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package tags

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/errors"
	"github.com/trackit/trackit/es"
)

type (
	esUntaggedValueResults struct {
		Buckets []struct {
			AccountID    string `json:"key"`
			ResourceType struct {
				Buckets []struct {
					ResourceType string `json:"key"`
					ResourceID   struct {
						Buckets []struct {
							ResourceID string `json:"key"`
						} `json:"buckets"`
					} `json:"resourceId"`
				} `json:"buckets"`
			} `json:"resourceType"`
		} `json:"buckets"`
	}

	// UntaggedResourceID contains the ResourceID
	UntaggedResourceID struct {
		ResourceID string `json:"resource_id"`
	}

	// UntaggedResourceType contains the ResourceType and a list of UntaggedResourceID
	UntaggedResourceType struct {
		ResourceType string               `json:"resource_type"`
		ResourceID   []UntaggedResourceID `json:"resource_ids"`
	}

	// UntaggedAccountID contains the AccountID and a list of UntaggedResourceType
	UntaggedAccountID struct {
		AccountID    string                 `json:"account_id"`
		ResourceType []UntaggedResourceType `json:"resource_types"`
	}

	// UntaggedResourcesResponse is the format for the endpoint response
	UntaggedResourcesResponse map[string][]UntaggedAccountID
)

// getUntaggedResourcesWithParsedParams will parse the data from teh ElasticSearch and return it
func getUntaggedResourcesWithParsedParams(ctx context.Context, params untaggedQueryParams) (int, interface{}) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	var typedDocument esUntaggedValueResults
	response := UntaggedResourcesResponse{}
	res, returnCode, err := makeElasticSearchRequestForUntaggedResources(ctx, params, es.Client)
	if err != nil {
		if returnCode == http.StatusOK {
			return returnCode, nil
		}
		return returnCode, errors.GetErrorMessage(ctx, err)
	}
	err = json.Unmarshal(*res.Aggregations["accountId"], &typedDocument)
	if err != nil {
		l.Error("Error while unmarshalling", err)
		return http.StatusInternalServerError, errors.GetErrorMessage(ctx, err)
	}
	var resultAccountID []UntaggedAccountID
	for _, accountID := range typedDocument.Buckets {
		var resultResourceType []UntaggedResourceType
		for _, resourceType := range accountID.ResourceType.Buckets {
			var resultResourceID []UntaggedResourceID
			for _, resourceID := range resourceType.ResourceID.Buckets {
				resultResourceID = append(resultResourceID, UntaggedResourceID{
					ResourceID: resourceID.ResourceID,
				})
			}
			resultResourceType = append(resultResourceType, UntaggedResourceType{
				ResourceType: resourceType.ResourceType,
				ResourceID:   resultResourceID,
			})
		}
		resultAccountID = append(resultAccountID, UntaggedAccountID{
			AccountID:    accountID.AccountID,
			ResourceType: resultResourceType,
		})
	}
	response[params.TagKey] = resultAccountID
	return http.StatusOK, response
}

// makeElasticSearchRequestForTagsValues will make the actual request to the ElasticSearch
// It will return the data, an http status code (as int) and an error.
// Because an error can be generated, but is not critical and is not needed to be known by
// the user (e.g if the index does not exists because it was not yet indexed ) the error will
// be returned, but instead of having a 500 status code, it will return the provided status code
// with empty data
func makeElasticSearchRequestForUntaggedResources(ctx context.Context, params untaggedQueryParams, client *elastic.Client) (*elastic.SearchResult, int, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	query := getUntaggedQuery(params)
	index := strings.Join(params.IndexList, ",")
	aggregation := getUntaggedAggregation()
	search := client.Search().Index(index).Size(0).Query(query).Pretty(true)
	search.Aggregation("accountId", aggregation)
	res, err := search.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists", map[string]interface{}{
				"index": index,
				"error": err.Error(),
			})
			return nil, http.StatusOK, err
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details.Type == "search_phase_execution_exception" {
			l.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{
				"err": err.Error(),
			})
		}
		return nil, http.StatusInternalServerError, err
	}
	return res, http.StatusOK, nil
}

// getUntaggedAggregation will generate the Aggregation for the query
func getUntaggedAggregation() *elastic.TermsAggregation {
	aggregation := elastic.NewTermsAggregation().Field("usageAccountId").Size(maxAggregationSize).SubAggregation("resourceType",
		elastic.NewTermsAggregation().Field("productCode").Size(maxAggregationSize).SubAggregation("resourceId",
			elastic.NewTermsAggregation().Field("resourceId").Size(maxAggregationSize)))
	return aggregation
}

// getUntaggedQuery will generate a query based on the params
func getUntaggedQuery(params untaggedQueryParams) *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	if len(params.AccountList) > 0 {
		query = query.Filter(createQueryAccountFilter(params.AccountList))
	}
	query = query.Filter(elastic.NewRangeQuery("usageStartDate").From(params.DateBegin).To(params.DateEnd))
	query = query.MustNot(elastic.NewBoolQuery().Filter(elastic.NewTermQuery("resourceId", "")))
	query = query.Must(elastic.NewNestedQuery("tags", getUntaggedNestedFilterQuery(params)))
	return query
}

// getUntaggedNestedFilterQuery will generate the nested Filter query based on the param
func getUntaggedNestedFilterQuery(params untaggedQueryParams) *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("tags.key", params.TagKey), elastic.NewTermQuery("tags.tag", ""))
	return query
}
