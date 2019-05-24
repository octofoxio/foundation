/*
 * Copyright (c) 2019. Octofox.io
 */

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExecute(t *testing.T) {
	var handler = execute("Get", "/test",
		EndpointHandler(func(ctx context.Context, request interface{}) (i interface{}, e error) {
			var req = request.(map[string]interface{})
			return map[string]interface{}{
				"message": req["id"],
			}, nil
		}),
		RequestDecoderMiddleware(func(ctx context.Context, r *http.Request) (i interface{}, e error) {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			var request map[string]interface{}
			err = json.Unmarshal(b, &request)
			if err != nil {
				return nil, err
			}
			return request, nil
		}),
		ResponseEncoderMiddleware(func(ctx context.Context, response interface{}) (i int, bytes []byte, e error) {
			b, _ := json.Marshal(response)
			return http.StatusCreated, b, nil
		}),
	)

	w := httptest.NewRecorder()
	b := bytes.NewBuffer([]byte("{\"id\": 29}"))
	req := httptest.NewRequest("GET", "/test", b)
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]int
	var err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, 29, response["message"])
}
