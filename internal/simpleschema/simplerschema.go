// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package simpleschema

import (
	"fmt"

	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// ToOpenAPISpec converts a SimpleSchema object to an OpenAPI schema.
//
// The input object is a map[string]interface{} where the key is the field name
// and the value is the field type.
func ToOpenAPISpec(obj map[string]interface{}) (*extv1.JSONSchemaProps, error) {
	tf := newTransformer()
	return tf.buildOpenAPISchema(obj)
}

// FromOpenAPISpec converts an OpenAPI schema to a SimpleSchema object.
func FromOpenAPISpec(schema *extv1.JSONSchemaProps) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}