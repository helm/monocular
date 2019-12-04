/*
Copyright (c) 2019 The Helm Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

//BodyAPIListResponse is an API body response in list format including the number of results pages
type BodyAPIListResponse struct {
	Data *ApiListResponse `json:"data"`
	Meta Meta             `json:"meta,omitempty"`
}

//BodyAPIResponse is an API body response in non-list format
type BodyAPIResponse struct {
	Data ApiResponse `json:"data"`
}

//ApiResponse is an API response in non-list format
type ApiResponse struct {
	ID            string      `json:"id"`
	Type          string      `json:"type"`
	Attributes    interface{} `json:"attributes"`
	Links         interface{} `json:"links"`
	Relationships RelMap      `json:"relationships"`
}

//ApiListResponse is an API response in list format
type ApiListResponse []*ApiResponse

//SelfLink the self-referencing URL to a chart in a response
type SelfLink struct {
	Self string `json:"self"`
}

//RelMap maps elements e.g. Charts to other elements of a response e.g. Chart Versions
type RelMap map[string]Rel

//Rel describes a relationship between element(s) in a response
type Rel struct {
	Data  interface{} `json:"data"`
	Links SelfLink    `json:"links"`
}

//Meta the number of pages in the response
type Meta struct {
	TotalPages int `json:"totalPages"`
}
