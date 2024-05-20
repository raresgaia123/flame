// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMongoService_DeleteDesign(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success_force_delete_true", func(mt *mtest.T) {
		db := &MongoService{
			designCollection: mt.Coll,
			jobCollection:    mt.Coll,
		}
		mt.AddMockResponses(
			bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 5}},
			bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 1}},
		)
		err := db.DeleteDesign("userid", "designid", true)
		assert.Nil(t, err)
	})

	mt.Run("success_force_delete_false", func(mt *mtest.T) {
		db := &MongoService{
			designCollection: mt.Coll,
			jobCollection:    mt.Coll,
		}
		mt.AddMockResponses(
			bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 5}},
			bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 1}},
		)
		err := db.DeleteDesign("userid", "designid", false)
		assert.Nil(t, err)
	})

	mt.Run("no_document_deleted", func(mt *mtest.T) {
		db := &MongoService{
			designCollection: mt.Coll,
			jobCollection:    mt.Coll,
		}
		mt.AddMockResponses(
			bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 5}},
			bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 0}},
		)
		err := db.DeleteDesign("userid", "designid", true)
		assert.NotNil(t, err)
		assert.Equal(t, "design id designid not found", err.Error(), "Unexpected error message")
	})
}
