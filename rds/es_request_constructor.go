//   Copyright 2017 MSolution.IO
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

// Package costs gets billing information from an ElasticSearch.
package rds

import (
	"time"

	"gopkg.in/olivere/elastic.v5"
)

// createQueryAccountFilterRds creates and return a new *elastic.TermsQuery on the accountList array
func createQueryAccountFilterRds(accountList []string) *elastic.TermsQuery {
	accountListFormatted := make([]interface{}, len(accountList))
	for i, v := range accountList {
		accountListFormatted[i] = v
	}
	return elastic.NewTermsQuery("account", accountListFormatted...)
}

// GetElasticSearchRdsParams is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
// It takes as paramters :
// 	- accountList []string : A slice of strings representing aws account number, in the format of the field
//	'awsdetailedlineitem.linked_account_id'
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on wich to execute the query. In this context the default value
//	should be "rds-reports"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func GetElasticSearchRdsParams(accountList []string, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	if len(accountList) > 0 {
		query = query.Filter(createQueryAccountFilterRds(accountList))
	}
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("top_reports", elastic.NewTermsAggregation().Field("account").
		SubAggregation("top_reports_hits", elastic.NewTopHitsAggregation().Sort("reportDate", false).Size(1)))
	return search
}

// createQueryAccountFilterBill creates and return a new *elastic.TermsQuery on the accountList array
func createQueryAccountFilterBill(accountList []string) *elastic.TermsQuery {
	accountListFormatted := make([]interface{}, len(accountList))
	for i, v := range accountList {
		accountListFormatted[i] = v
	}
	return elastic.NewTermsQuery("usageAccountId", accountListFormatted...)
}

// GetElasticSearchCostParams is used to construct an ElasticSearch *elastic.SearchService used to perform a request on ES
// It takes as paramters :
// 	- accountList []string : A slice of strings representing aws account number, in the format of the field
//	'awsdetailedlineitem.linked_account_id'
//	- client *elastic.Client : an instance of *elastic.Client that represent an Elastic Search client.
//	It needs to be fully configured and ready to execute a client.Search()
//	- index string : The Elastic Search index on wich to execute the query. In this context the default value
//	should be "rds-reports"
// This function excepts arguments passed to it to be sanitize. If they are not, the following cases will make
// it crash :
//	- If the client is nil or malconfigured, it will crash
//	- If the index is not an index present in the ES, it will crash
func GetElasticSearchCostParams(accountList []string, client *elastic.Client, index string) *elastic.SearchService {
	query := elastic.NewBoolQuery()
	if len(accountList) > 0 {
		query = query.Filter(createQueryAccountFilterBill(accountList))
	}
	query = query.Filter(elastic.NewTermQuery("productCode", "AmazonRDS"))
	dateEnd := time.Now().UTC()
	dateStart := time.Date(dateEnd.Year(), dateEnd.Month(), 1, 0, 0, 0, 0, dateEnd.Location()).UTC()
	query = query.Filter(elastic.NewRangeQuery("usageStartDate").
		From(dateStart).To(dateEnd))
	search := client.Search().Index(index).Size(0).Query(query)
	search.Aggregation("accounts",  elastic.NewTermsAggregation().Field("usageAccountId").Size(len(accountList)).
		SubAggregation("instances", elastic.NewTermsAggregation().Field("resourceId").Size(0x7FFFFFFF).
		SubAggregation("cost",  elastic.NewSumAggregation().Field("unblendedCost"))))
	return search
}