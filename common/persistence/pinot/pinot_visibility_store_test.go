// The MIT License (MIT)

// Copyright (c) 2017-2020 Uber Technologies Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package pinotvisibility

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/uber/cadence/.gen/go/indexer"
	"github.com/uber/cadence/common"
	"github.com/uber/cadence/common/definition"
	"github.com/uber/cadence/common/dynamicconfig"
	"github.com/uber/cadence/common/log"
	"github.com/uber/cadence/common/log/testlogger"
	"github.com/uber/cadence/common/mocks"
	p "github.com/uber/cadence/common/persistence"
	pnt "github.com/uber/cadence/common/pinot"
	"github.com/uber/cadence/common/service"
	"github.com/uber/cadence/common/types"
)

var (
	testIndex        = "test-index"
	testDomain       = "test-domain"
	testDomainID     = "bfd5c907-f899-4baf-a7b2-2ab85e623ebd"
	testPageSize     = 10
	testEarliestTime = int64(1547596872371000000)
	testLatestTime   = int64(2547596872371000000)
	testWorkflowType = "test-wf-type"
	testWorkflowID   = "test-wid"
	testCloseStatus  = int32(1)
	testTableName    = "test-table-name"

	testContextTimeout = 5 * time.Second

	validSearchAttr = definition.GetDefaultIndexedKeys()
)

func TestRecordWorkflowExecutionStarted(t *testing.T) {

	// test non-empty request fields match
	errorRequest := &p.InternalRecordWorkflowExecutionStartedRequest{
		WorkflowID: "wid",
		Memo:       p.NewDataBlob([]byte(`test bytes`), common.EncodingTypeThriftRW),
		SearchAttributes: map[string][]byte{
			"CustomStringField": []byte("test string"),
			"CustomTimeField":   []byte("2020-01-01T00:00:00Z"),
		},
	}

	request := &p.InternalRecordWorkflowExecutionStartedRequest{
		WorkflowID: "wid",
		Memo:       p.NewDataBlob([]byte(`test bytes`), common.EncodingTypeThriftRW),
	}

	tests := map[string]struct {
		request       *p.InternalRecordWorkflowExecutionStartedRequest
		expectedError error
	}{
		"Case1: error case": {
			request:       errorRequest,
			expectedError: fmt.Errorf("error"),
		},
		"Case2: normal case": {
			request:       request,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockProducer.On("Publish", mock.Anything, mock.MatchedBy(func(input *indexer.PinotMessage) bool {
				assert.Equal(t, request.WorkflowID, input.GetWorkflowID())
				return true
			})).Return(nil).Once()

			err := visibilityStore.RecordWorkflowExecutionStarted(context.Background(), test.request)
			if test.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRecordWorkflowExecutionClosed(t *testing.T) {
	// test non-empty request fields match
	errorRequest := &p.InternalRecordWorkflowExecutionClosedRequest{
		WorkflowID: "wid",
		Memo:       p.NewDataBlob([]byte(`test bytes`), common.EncodingTypeThriftRW),
		SearchAttributes: map[string][]byte{
			"CustomStringField": []byte("test string"),
			"CustomTimeField":   []byte("2020-01-01T00:00:00Z"),
		},
	}
	request := &p.InternalRecordWorkflowExecutionClosedRequest{
		WorkflowID: "wid",
		Memo:       p.NewDataBlob([]byte(`test bytes`), common.EncodingTypeThriftRW),
	}

	tests := map[string]struct {
		request       *p.InternalRecordWorkflowExecutionClosedRequest
		expectedError error
	}{
		"Case1: error case": {
			request:       errorRequest,
			expectedError: fmt.Errorf("error"),
		},
		"Case2: normal case": {
			request:       request,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockProducer.On("Publish", mock.Anything, mock.MatchedBy(func(input *indexer.PinotMessage) bool {
				assert.Equal(t, request.WorkflowID, input.GetWorkflowID())
				return true
			})).Return(nil).Once()

			err := visibilityStore.RecordWorkflowExecutionClosed(context.Background(), test.request)
			if test.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRecordWorkflowExecutionUninitialized(t *testing.T) {
	// test non-empty request fields match
	request := &p.InternalRecordWorkflowExecutionUninitializedRequest{
		WorkflowID: "wid",
	}

	tests := map[string]struct {
		request       *p.InternalRecordWorkflowExecutionUninitializedRequest
		expectedError error
	}{
		"Case1: normal case": {
			request:       request,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockProducer.On("Publish", mock.Anything, mock.MatchedBy(func(input *indexer.PinotMessage) bool {
				assert.Equal(t, request.WorkflowID, input.GetWorkflowID())
				return true
			})).Return(nil).Once()

			err := visibilityStore.RecordWorkflowExecutionUninitialized(context.Background(), test.request)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestUpsertWorkflowExecution(t *testing.T) {
	// test non-empty request fields match
	request := &p.InternalUpsertWorkflowExecutionRequest{}
	request.WorkflowID = "wid"
	memoBytes := []byte(`test bytes`)
	request.Memo = p.NewDataBlob(memoBytes, common.EncodingTypeThriftRW)

	tests := map[string]struct {
		request       *p.InternalUpsertWorkflowExecutionRequest
		expectedError error
	}{
		"Case1: normal case": {
			request:       request,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockProducer.On("Publish", mock.Anything, mock.MatchedBy(func(input *indexer.PinotMessage) bool {
				assert.Equal(t, request.WorkflowID, input.GetWorkflowID())
				return true
			})).Return(nil).Once()

			err := visibilityStore.UpsertWorkflowExecution(context.Background(), test.request)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestDeleteWorkflowExecution(t *testing.T) {
	// test non-empty request fields match
	request := &p.VisibilityDeleteWorkflowExecutionRequest{}
	request.WorkflowID = "wid"

	tests := map[string]struct {
		request       *p.VisibilityDeleteWorkflowExecutionRequest
		expectedError error
	}{
		"Case1: normal case": {
			request:       request,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockProducer.On("Publish", mock.Anything, mock.MatchedBy(func(input *indexer.PinotMessage) bool {
				assert.Equal(t, request.WorkflowID, input.GetWorkflowID())
				return true
			})).Return(nil).Once()

			err := visibilityStore.DeleteWorkflowExecution(context.Background(), test.request)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestDeleteUninitializedWorkflowExecution(t *testing.T) {
	// test non-empty request fields match
	request := &p.VisibilityDeleteWorkflowExecutionRequest{
		Domain:     "domain",
		DomainID:   "domainID",
		WorkflowID: "wid",
		RunID:      "rid",
		TaskID:     int64(111),
	}

	tests := map[string]struct {
		request       *p.VisibilityDeleteWorkflowExecutionRequest
		expectedError error
	}{
		"Case1: normal case": {
			request:       request,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockPinotClient.EXPECT().GetTableName().Return(testTableName).Times(1)
			mockPinotClient.EXPECT().CountByQuery(gomock.Any()).Return(int64(1), nil).Times(1)

			mockProducer.On("Publish", mock.Anything, mock.MatchedBy(func(input *indexer.PinotMessage) bool {
				assert.Equal(t, request.WorkflowID, input.GetWorkflowID())
				return true
			})).Return(nil).Once()

			err := visibilityStore.DeleteUninitializedWorkflowExecution(context.Background(), test.request)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestListOpenWorkflowExecutions(t *testing.T) {
	request := &p.InternalListWorkflowExecutionsRequest{
		Domain: DomainID,
	}

	tests := map[string]struct {
		request       *p.InternalListWorkflowExecutionsRequest
		expectedResp  *p.InternalListWorkflowExecutionsResponse
		expectedError error
	}{
		"Case1: normal case with nil response": {
			request:       request,
			expectedResp:  nil,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockPinotClient.EXPECT().GetTableName().Return(testTableName).Times(1)
			mockPinotClient.EXPECT().Search(gomock.Any()).Return(nil, nil).Times(1)
			resp, err := visibilityStore.ListOpenWorkflowExecutions(context.Background(), test.request)
			assert.Equal(t, test.expectedResp, resp)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestListClosedWorkflowExecutions(t *testing.T) {
	request := &p.InternalListWorkflowExecutionsRequest{
		Domain: DomainID,
	}

	tests := map[string]struct {
		request       *p.InternalListWorkflowExecutionsRequest
		expectedResp  *p.InternalListWorkflowExecutionsResponse
		expectedError error
	}{
		"Case1: normal case with nil response": {
			request:       request,
			expectedResp:  nil,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockPinotClient.EXPECT().GetTableName().Return(testTableName).Times(1)
			mockPinotClient.EXPECT().Search(gomock.Any()).Return(nil, nil).Times(1)
			resp, err := visibilityStore.ListClosedWorkflowExecutions(context.Background(), test.request)
			assert.Equal(t, test.expectedResp, resp)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestListOpenWorkflowExecutionsByType(t *testing.T) {
	request := &p.InternalListWorkflowExecutionsByTypeRequest{}

	tests := map[string]struct {
		request       *p.InternalListWorkflowExecutionsByTypeRequest
		expectedResp  *p.InternalListWorkflowExecutionsResponse
		expectedError error
	}{
		"Case1: normal case with nil response": {
			request:       request,
			expectedResp:  nil,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockPinotClient.EXPECT().GetTableName().Return(testTableName).Times(1)
			mockPinotClient.EXPECT().Search(gomock.Any()).Return(nil, nil).Times(1)
			resp, err := visibilityStore.ListOpenWorkflowExecutionsByType(context.Background(), test.request)
			assert.Equal(t, test.expectedResp, resp)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestListClosedWorkflowExecutionsByType(t *testing.T) {
	request := &p.InternalListWorkflowExecutionsByTypeRequest{}

	tests := map[string]struct {
		request       *p.InternalListWorkflowExecutionsByTypeRequest
		expectedResp  *p.InternalListWorkflowExecutionsResponse
		expectedError error
	}{
		"Case1: normal case with nil response": {
			request:       request,
			expectedResp:  nil,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockPinotClient.EXPECT().GetTableName().Return(testTableName).Times(1)
			mockPinotClient.EXPECT().Search(gomock.Any()).Return(nil, nil).Times(1)
			resp, err := visibilityStore.ListClosedWorkflowExecutionsByType(context.Background(), test.request)
			assert.Equal(t, test.expectedResp, resp)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestListOpenWorkflowExecutionsByWorkflowID(t *testing.T) {
	request := &p.InternalListWorkflowExecutionsByWorkflowIDRequest{}

	tests := map[string]struct {
		request       *p.InternalListWorkflowExecutionsByWorkflowIDRequest
		expectedResp  *p.InternalListWorkflowExecutionsResponse
		expectedError error
	}{
		"Case1: normal case with nil response": {
			request:       request,
			expectedResp:  nil,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockPinotClient.EXPECT().GetTableName().Return(testTableName).Times(1)
			mockPinotClient.EXPECT().Search(gomock.Any()).Return(nil, nil).Times(1)
			resp, err := visibilityStore.ListOpenWorkflowExecutionsByWorkflowID(context.Background(), test.request)
			assert.Equal(t, test.expectedResp, resp)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestListClosedWorkflowExecutionsByWorkflowID(t *testing.T) {
	request := &p.InternalListWorkflowExecutionsByWorkflowIDRequest{}

	tests := map[string]struct {
		request       *p.InternalListWorkflowExecutionsByWorkflowIDRequest
		expectedResp  *p.InternalListWorkflowExecutionsResponse
		expectedError error
	}{
		"Case1: normal case with nil response": {
			request:       request,
			expectedResp:  nil,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockPinotClient.EXPECT().GetTableName().Return(testTableName).Times(1)
			mockPinotClient.EXPECT().Search(gomock.Any()).Return(nil, nil).Times(1)
			resp, err := visibilityStore.ListClosedWorkflowExecutionsByWorkflowID(context.Background(), test.request)
			assert.Equal(t, test.expectedResp, resp)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestListClosedWorkflowExecutionsByStatus(t *testing.T) {
	request := &p.InternalListClosedWorkflowExecutionsByStatusRequest{}

	tests := map[string]struct {
		request       *p.InternalListClosedWorkflowExecutionsByStatusRequest
		expectedResp  *p.InternalListWorkflowExecutionsResponse
		expectedError error
	}{
		"Case1: normal case with nil response": {
			request:       request,
			expectedResp:  nil,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockPinotClient.EXPECT().GetTableName().Return(testTableName).Times(1)
			mockPinotClient.EXPECT().Search(gomock.Any()).Return(nil, nil).Times(1)
			resp, err := visibilityStore.ListClosedWorkflowExecutionsByStatus(context.Background(), test.request)
			assert.Equal(t, test.expectedResp, resp)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestGetClosedWorkflowExecution(t *testing.T) {
	request := &p.InternalGetClosedWorkflowExecutionRequest{}

	tests := map[string]struct {
		request       *p.InternalGetClosedWorkflowExecutionRequest
		expectedResp  *p.InternalGetClosedWorkflowExecutionRequest
		expectedError error
	}{
		"Case1: normal case with nil response": {
			request:       request,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockPinotClient.EXPECT().GetTableName().Return(testTableName).Times(1)
			mockPinotClient.EXPECT().Search(gomock.Any()).Return(&pnt.SearchResponse{
				Executions: []*p.InternalVisibilityWorkflowExecutionInfo{
					{
						DomainID: DomainID,
					},
				},
			}, nil).Times(1)
			_, err := visibilityStore.GetClosedWorkflowExecution(context.Background(), test.request)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestListWorkflowExecutions(t *testing.T) {
	request := &p.ListWorkflowExecutionsByQueryRequest{}

	tests := map[string]struct {
		request       *p.ListWorkflowExecutionsByQueryRequest
		expectedResp  *p.InternalListWorkflowExecutionsResponse
		expectedError error
	}{
		"Case1: normal case with nil response": {
			request:       request,
			expectedResp:  nil,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockPinotClient.EXPECT().GetTableName().Return(testTableName).Times(1)
			mockPinotClient.EXPECT().Search(gomock.Any()).Return(nil, nil).Times(1)
			resp, err := visibilityStore.ListWorkflowExecutions(context.Background(), test.request)
			assert.Equal(t, test.expectedResp, resp)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestScanWorkflowExecutions(t *testing.T) {
	request := &p.ListWorkflowExecutionsByQueryRequest{}

	tests := map[string]struct {
		request       *p.ListWorkflowExecutionsByQueryRequest
		expectedResp  *p.InternalListWorkflowExecutionsResponse
		expectedError error
	}{
		"Case1: normal case with nil response": {
			request:       request,
			expectedResp:  nil,
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			mockPinotClient.EXPECT().GetTableName().Return(testTableName).Times(1)
			mockPinotClient.EXPECT().Search(gomock.Any()).Return(nil, nil).Times(1)
			resp, err := visibilityStore.ScanWorkflowExecutions(context.Background(), test.request)
			assert.Equal(t, test.expectedResp, resp)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestGetName(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPinotClient := pnt.NewMockGenericClient(ctrl)
	mockProducer := &mocks.KafkaProducer{}
	mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
		ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
		ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
	}, mockProducer, testlogger.New(t))
	visibilityStore := mgr.(*pinotVisibilityStore)
	assert.NotEmpty(t, visibilityStore.GetName())
}

func TestNewPinotVisibilityStore(t *testing.T) {
	mockPinotClient := &pnt.MockGenericClient{}
	assert.NotPanics(t, func() {
		NewPinotVisibilityStore(mockPinotClient, &service.Config{
			ValidSearchAttributes: dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
		}, nil, log.NewNoop())
	})
}

func TestGetCountWorkflowExecutionsQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPinotClient := pnt.NewMockGenericClient(ctrl)
	mockProducer := &mocks.KafkaProducer{}
	mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
		ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
		ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
	}, mockProducer, testlogger.New(t))
	visibilityStore := mgr.(*pinotVisibilityStore)

	request := &p.CountWorkflowExecutionsRequest{
		DomainUUID: testDomainID,
		Domain:     testDomain,
		Query:      "WorkflowID = 'wfid'",
	}

	result := visibilityStore.getCountWorkflowExecutionsQuery(testTableName, request)
	expectResult := fmt.Sprintf(`SELECT COUNT(*)
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND WorkflowID = 'wfid'
`, testTableName)

	assert.Equal(t, result, expectResult)

	nilResult := visibilityStore.getCountWorkflowExecutionsQuery(testTableName, nil)
	assert.Equal(t, nilResult, "")
}

func TestGetListWorkflowExecutionQuery(t *testing.T) {
	token := pnt.PinotVisibilityPageToken{
		From: 11,
	}

	serializedToken, err := json.Marshal(token)
	if err != nil {
		panic(fmt.Sprintf("Serialized error in PinotVisibilityStoreTest!!!, %s", err))
	}

	tests := map[string]struct {
		input          *p.ListWorkflowExecutionsByQueryRequest
		expectedOutput string
	}{
		"complete request with keyword query only": {
			input: &p.ListWorkflowExecutionsByQueryRequest{
				DomainUUID:    testDomainID,
				Domain:        testDomain,
				PageSize:      testPageSize,
				NextPageToken: nil,
				Query:         "`Attr.CustomKeywordField` = 'keywordCustomized'",
			},
			expectedOutput: fmt.Sprintf(
				`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND (JSON_MATCH(Attr, '"$.CustomKeywordField"=''keywordCustomized''') or JSON_MATCH(Attr, '"$.CustomKeywordField[*]"=''keywordCustomized'''))
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName),
		},

		"complete request from search attribute worker": {
			input: &p.ListWorkflowExecutionsByQueryRequest{
				DomainUUID:    testDomainID,
				Domain:        testDomain,
				PageSize:      testPageSize,
				NextPageToken: nil,
				Query:         "CustomIntField=2 and CustomKeywordField='Update2' order by `Attr.CustomDatetimeField` DESC",
			},
			expectedOutput: fmt.Sprintf(
				`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND JSON_MATCH(Attr, '"$.CustomIntField"=''2''') and (JSON_MATCH(Attr, '"$.CustomKeywordField"=''Update2''') or JSON_MATCH(Attr, '"$.CustomKeywordField[*]"=''Update2'''))
order by CustomDatetimeField DESC
LIMIT 0, 10
`, testTableName),
		},

		"complete request with keyword query and other customized query": {
			input: &p.ListWorkflowExecutionsByQueryRequest{
				DomainUUID:    testDomainID,
				Domain:        testDomain,
				PageSize:      testPageSize,
				NextPageToken: nil,
				Query:         "CustomKeywordField = 'keywordCustomized' and CustomStringField = 'String and or order by'",
			},
			expectedOutput: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND (JSON_MATCH(Attr, '"$.CustomKeywordField"=''keywordCustomized''') or JSON_MATCH(Attr, '"$.CustomKeywordField[*]"=''keywordCustomized''')) and (JSON_MATCH(Attr, '"$.CustomStringField" is not null') AND REGEXP_LIKE(JSON_EXTRACT_SCALAR(Attr, '$.CustomStringField', 'string'), 'String and or order by*'))
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName),
		},

		"complete request with or query & customized attributes": {
			input: &p.ListWorkflowExecutionsByQueryRequest{
				DomainUUID:    testDomainID,
				Domain:        testDomain,
				PageSize:      testPageSize,
				NextPageToken: nil,
				Query:         "CustomStringField = 'Or' or CustomStringField = 'and' Order by StartTime DESC",
			},
			expectedOutput: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND ((JSON_MATCH(Attr, '"$.CustomStringField" is not null') AND REGEXP_LIKE(JSON_EXTRACT_SCALAR(Attr, '$.CustomStringField', 'string'), 'Or*')) or (JSON_MATCH(Attr, '"$.CustomStringField" is not null') AND REGEXP_LIKE(JSON_EXTRACT_SCALAR(Attr, '$.CustomStringField', 'string'), 'and*')))
Order by StartTime DESC
LIMIT 0, 10
`, testTableName),
		},

		"complex query": {
			input: &p.ListWorkflowExecutionsByQueryRequest{
				DomainUUID:    testDomainID,
				Domain:        testDomain,
				PageSize:      testPageSize,
				NextPageToken: nil,
				Query:         "WorkflowID = 'wid' and ((CustomStringField = 'custom and custom2 or custom3 order by') or CustomIntField between 1 and 10)",
			},
			expectedOutput: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND WorkflowID = 'wid' and ((JSON_MATCH(Attr, '"$.CustomStringField" is not null') AND REGEXP_LIKE(JSON_EXTRACT_SCALAR(Attr, '$.CustomStringField', 'string'), 'custom and custom2 or custom3 order by*')) or (JSON_MATCH(Attr, '"$.CustomIntField" is not null') AND CAST(JSON_EXTRACT_SCALAR(Attr, '$.CustomIntField') AS INT) >= 1 AND CAST(JSON_EXTRACT_SCALAR(Attr, '$.CustomIntField') AS INT) <= 10))
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName),
		},

		"or clause with custom attributes": {
			input: &p.ListWorkflowExecutionsByQueryRequest{
				DomainUUID:    testDomainID,
				Domain:        testDomain,
				PageSize:      testPageSize,
				NextPageToken: nil,
				Query:         "CustomIntField = 1 or CustomIntField = 2",
			},
			expectedOutput: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND (JSON_MATCH(Attr, '"$.CustomIntField"=''1''') or JSON_MATCH(Attr, '"$.CustomIntField"=''2'''))
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName),
		},

		"complete request with customized query with missing": {
			input: &p.ListWorkflowExecutionsByQueryRequest{
				DomainUUID:    testDomainID,
				Domain:        testDomain,
				PageSize:      testPageSize,
				NextPageToken: nil,
				Query:         "CloseTime = missing anD WorkflowType = 'some-test-workflow'",
			},
			expectedOutput: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND CloseTime = -1 and WorkflowType = 'some-test-workflow'
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName),
		},

		"complete request with customized query with NextPageToken": {
			input: &p.ListWorkflowExecutionsByQueryRequest{
				DomainUUID:    testDomainID,
				Domain:        testDomain,
				PageSize:      testPageSize,
				NextPageToken: serializedToken,
				Query:         "CloseStatus < 0 and CustomKeywordField = 'keywordCustomized' AND CustomIntField<=10 and CustomStringField = 'String field is for text' Order by DomainID Desc",
			},
			expectedOutput: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND CloseStatus < 0 and (JSON_MATCH(Attr, '"$.CustomKeywordField"=''keywordCustomized''') or JSON_MATCH(Attr, '"$.CustomKeywordField[*]"=''keywordCustomized''')) and (JSON_MATCH(Attr, '"$.CustomIntField" is not null') AND CAST(JSON_EXTRACT_SCALAR(Attr, '$.CustomIntField') AS INT) <= 10) and (JSON_MATCH(Attr, '"$.CustomStringField" is not null') AND REGEXP_LIKE(JSON_EXTRACT_SCALAR(Attr, '$.CustomStringField', 'string'), 'String field is for text*'))
Order by DomainID Desc
LIMIT 11, 10
`, testTableName),
		},

		"complete request with order by query": {
			input: &p.ListWorkflowExecutionsByQueryRequest{
				DomainUUID:    testDomainID,
				Domain:        testDomain,
				PageSize:      testPageSize,
				NextPageToken: nil,
				Query:         "Order by DomainId Desc",
			},
			expectedOutput: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
Order by DomainId Desc
LIMIT 0, 10
`, testTableName),
		},

		"complete request with filter query": {
			input: &p.ListWorkflowExecutionsByQueryRequest{
				DomainUUID:    testDomainID,
				Domain:        testDomain,
				PageSize:      testPageSize,
				NextPageToken: nil,
				Query:         "CloseStatus < 0",
			},
			expectedOutput: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND CloseStatus < 0
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName),
		},

		"complete request with empty query": {
			input: &p.ListWorkflowExecutionsByQueryRequest{
				DomainUUID:    testDomainID,
				Domain:        testDomain,
				PageSize:      testPageSize,
				NextPageToken: nil,
				Query:         "",
			},
			expectedOutput: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
LIMIT 0, 10
`, testTableName),
		},

		"empty request": {
			input: &p.ListWorkflowExecutionsByQueryRequest{},
			expectedOutput: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = ''
AND IsDeleted = false
LIMIT 0, 0
`, testTableName),
		},

		"nil request": {
			input:          nil,
			expectedOutput: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockPinotClient := pnt.NewMockGenericClient(ctrl)
			mockProducer := &mocks.KafkaProducer{}
			mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
				ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
				ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
			}, mockProducer, testlogger.New(t))
			visibilityStore := mgr.(*pinotVisibilityStore)

			output, err := visibilityStore.getListWorkflowExecutionsByQueryQuery(testTableName, test.input)
			assert.Equal(t, test.expectedOutput, output)
			assert.NoError(t, err)
		})
	}
}

func TestGetListWorkflowExecutionsQuery(t *testing.T) {
	request := &p.InternalListWorkflowExecutionsRequest{
		DomainUUID:    testDomainID,
		Domain:        testDomain,
		EarliestTime:  time.Unix(0, testEarliestTime),
		LatestTime:    time.Unix(0, testLatestTime),
		PageSize:      testPageSize,
		NextPageToken: nil,
	}

	closeResult, err1 := getListWorkflowExecutionsQuery(testTableName, request, true)
	openResult, err2 := getListWorkflowExecutionsQuery(testTableName, request, false)
	nilResult, err3 := getListWorkflowExecutionsQuery(testTableName, nil, true)
	expectCloseResult := fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND CloseTime BETWEEN 1547596871371 AND 2547596873371
AND CloseStatus >= 0
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName)
	expectOpenResult := fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND StartTime BETWEEN 1547596871371 AND 2547596873371
AND CloseStatus < 0
AND CloseTime = -1
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName)
	expectNilResult := ""

	assert.Equal(t, closeResult, expectCloseResult)
	assert.Equal(t, openResult, expectOpenResult)
	assert.Equal(t, nilResult, expectNilResult)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)
}

func TestGetListWorkflowExecutionsByTypeQuery(t *testing.T) {
	request := &p.InternalListWorkflowExecutionsByTypeRequest{
		InternalListWorkflowExecutionsRequest: p.InternalListWorkflowExecutionsRequest{
			DomainUUID:    testDomainID,
			Domain:        testDomain,
			EarliestTime:  time.Unix(0, testEarliestTime),
			LatestTime:    time.Unix(0, testLatestTime),
			PageSize:      testPageSize,
			NextPageToken: nil,
		},
		WorkflowTypeName: testWorkflowType,
	}

	closeResult, err1 := getListWorkflowExecutionsByTypeQuery(testTableName, request, true)
	openResult, err2 := getListWorkflowExecutionsByTypeQuery(testTableName, request, false)
	nilResult, err3 := getListWorkflowExecutionsByTypeQuery(testTableName, nil, true)
	expectCloseResult := fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND WorkflowType = 'test-wf-type'
AND CloseTime BETWEEN 1547596871371 AND 2547596873371
AND CloseStatus >= 0
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName)
	expectOpenResult := fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND WorkflowType = 'test-wf-type'
AND StartTime BETWEEN 1547596871371 AND 2547596873371
AND CloseStatus < 0
AND CloseTime = -1
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName)
	expectNilResult := ""

	assert.Equal(t, closeResult, expectCloseResult)
	assert.Equal(t, openResult, expectOpenResult)
	assert.Equal(t, nilResult, expectNilResult)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)
}

func TestGetListWorkflowExecutionsByWorkflowIDQuery(t *testing.T) {
	request := &p.InternalListWorkflowExecutionsByWorkflowIDRequest{
		InternalListWorkflowExecutionsRequest: p.InternalListWorkflowExecutionsRequest{
			DomainUUID:    testDomainID,
			Domain:        testDomain,
			EarliestTime:  time.Unix(0, testEarliestTime),
			LatestTime:    time.Unix(0, testLatestTime),
			PageSize:      testPageSize,
			NextPageToken: nil,
		},
		WorkflowID: testWorkflowID,
	}

	closeResult, err1 := getListWorkflowExecutionsByWorkflowIDQuery(testTableName, request, true)
	openResult, err2 := getListWorkflowExecutionsByWorkflowIDQuery(testTableName, request, false)
	nilResult, err3 := getListWorkflowExecutionsByWorkflowIDQuery(testTableName, nil, true)
	expectCloseResult := fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND WorkflowID = 'test-wid'
AND CloseTime BETWEEN 1547596871371 AND 2547596873371
AND CloseStatus >= 0
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName)
	expectOpenResult := fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND WorkflowID = 'test-wid'
AND StartTime BETWEEN 1547596871371 AND 2547596873371
AND CloseStatus < 0
AND CloseTime = -1
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName)
	expectNilResult := ""

	assert.Equal(t, closeResult, expectCloseResult)
	assert.Equal(t, openResult, expectOpenResult)
	assert.Equal(t, nilResult, expectNilResult)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)
}

func TestGetListWorkflowExecutionsByStatusQuery(t *testing.T) {
	tests := map[string]struct {
		inputRequest *p.InternalListClosedWorkflowExecutionsByStatusRequest
		expectResult string
		expectError  error
	}{
		"Case1: normal case": {
			inputRequest: nil,
			expectResult: "",
			expectError:  nil,
		},
		"Case2-0: normal case with close status is 0": {
			inputRequest: &p.InternalListClosedWorkflowExecutionsByStatusRequest{
				InternalListWorkflowExecutionsRequest: p.InternalListWorkflowExecutionsRequest{
					DomainUUID:    testDomainID,
					Domain:        testDomain,
					EarliestTime:  time.Unix(0, testEarliestTime),
					LatestTime:    time.Unix(0, testLatestTime),
					PageSize:      testPageSize,
					NextPageToken: nil,
				},
				Status: types.WorkflowExecutionCloseStatus(0),
			},
			expectResult: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND CloseStatus = '0'
AND CloseTime BETWEEN 1547596872371 AND 2547596872371
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName),
			expectError: nil,
		},
		"Case2-1: normal case with close status is 1": {
			inputRequest: &p.InternalListClosedWorkflowExecutionsByStatusRequest{
				InternalListWorkflowExecutionsRequest: p.InternalListWorkflowExecutionsRequest{
					DomainUUID:    testDomainID,
					Domain:        testDomain,
					EarliestTime:  time.Unix(0, testEarliestTime),
					LatestTime:    time.Unix(0, testLatestTime),
					PageSize:      testPageSize,
					NextPageToken: nil,
				},
				Status: types.WorkflowExecutionCloseStatus(1),
			},
			expectResult: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND CloseStatus = '1'
AND CloseTime BETWEEN 1547596872371 AND 2547596872371
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName),
			expectError: nil,
		},
		"Case2-2: normal case with close status is 2": {
			inputRequest: &p.InternalListClosedWorkflowExecutionsByStatusRequest{
				InternalListWorkflowExecutionsRequest: p.InternalListWorkflowExecutionsRequest{
					DomainUUID:    testDomainID,
					Domain:        testDomain,
					EarliestTime:  time.Unix(0, testEarliestTime),
					LatestTime:    time.Unix(0, testLatestTime),
					PageSize:      testPageSize,
					NextPageToken: nil,
				},
				Status: types.WorkflowExecutionCloseStatus(2),
			},
			expectResult: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND CloseStatus = '2'
AND CloseTime BETWEEN 1547596872371 AND 2547596872371
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName),
			expectError: nil,
		},
		"Case2-3: normal case with close status is 3": {
			inputRequest: &p.InternalListClosedWorkflowExecutionsByStatusRequest{
				InternalListWorkflowExecutionsRequest: p.InternalListWorkflowExecutionsRequest{
					DomainUUID:    testDomainID,
					Domain:        testDomain,
					EarliestTime:  time.Unix(0, testEarliestTime),
					LatestTime:    time.Unix(0, testLatestTime),
					PageSize:      testPageSize,
					NextPageToken: nil,
				},
				Status: types.WorkflowExecutionCloseStatus(3),
			},
			expectResult: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND CloseStatus = '3'
AND CloseTime BETWEEN 1547596872371 AND 2547596872371
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName),
			expectError: nil,
		},
		"Case2-4: normal case with close status is 4": {
			inputRequest: &p.InternalListClosedWorkflowExecutionsByStatusRequest{
				InternalListWorkflowExecutionsRequest: p.InternalListWorkflowExecutionsRequest{
					DomainUUID:    testDomainID,
					Domain:        testDomain,
					EarliestTime:  time.Unix(0, testEarliestTime),
					LatestTime:    time.Unix(0, testLatestTime),
					PageSize:      testPageSize,
					NextPageToken: nil,
				},
				Status: types.WorkflowExecutionCloseStatus(4),
			},
			expectResult: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND CloseStatus = '4'
AND CloseTime BETWEEN 1547596872371 AND 2547596872371
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName),
			expectError: nil,
		},
		"Case2-5: normal case with close status is 5": {
			inputRequest: &p.InternalListClosedWorkflowExecutionsByStatusRequest{
				InternalListWorkflowExecutionsRequest: p.InternalListWorkflowExecutionsRequest{
					DomainUUID:    testDomainID,
					Domain:        testDomain,
					EarliestTime:  time.Unix(0, testEarliestTime),
					LatestTime:    time.Unix(0, testLatestTime),
					PageSize:      testPageSize,
					NextPageToken: nil,
				},
				Status: types.WorkflowExecutionCloseStatus(5),
			},
			expectResult: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND CloseStatus = '5'
AND CloseTime BETWEEN 1547596872371 AND 2547596872371
Order BY StartTime DESC
LIMIT 0, 10
`, testTableName),
			expectError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualResult, actualError := getListWorkflowExecutionsByStatusQuery(testTableName, test.inputRequest)
			assert.Equal(t, test.expectResult, actualResult)
			assert.NoError(t, actualError)
		})
	}
}

func TestGetGetClosedWorkflowExecutionQuery(t *testing.T) {
	tests := map[string]struct {
		input          *p.InternalGetClosedWorkflowExecutionRequest
		expectedOutput string
	}{
		"complete request with empty RunId": {
			input: &p.InternalGetClosedWorkflowExecutionRequest{
				DomainUUID: testDomainID,
				Domain:     testDomain,
				Execution: types.WorkflowExecution{
					WorkflowID: testWorkflowID,
					RunID:      "",
				},
			},
			expectedOutput: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND CloseStatus >= 0
AND WorkflowID = 'test-wid'
`, testTableName),
		},

		"complete request with runId": {
			input: &p.InternalGetClosedWorkflowExecutionRequest{
				DomainUUID: testDomainID,
				Domain:     testDomain,
				Execution: types.WorkflowExecution{
					WorkflowID: testWorkflowID,
					RunID:      "runid",
				},
			},
			expectedOutput: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = 'bfd5c907-f899-4baf-a7b2-2ab85e623ebd'
AND IsDeleted = false
AND CloseStatus >= 0
AND WorkflowID = 'test-wid'
AND RunID = 'runid'
`, testTableName),
		},

		"empty request": {
			input: &p.InternalGetClosedWorkflowExecutionRequest{},
			expectedOutput: fmt.Sprintf(`SELECT *
FROM %s
WHERE DomainID = ''
AND IsDeleted = false
AND CloseStatus >= 0
AND WorkflowID = ''
`, testTableName),
		},

		"nil request": {
			input:          nil,
			expectedOutput: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output := getGetClosedWorkflowExecutionQuery(testTableName, test.input)
			assert.Equal(t, test.expectedOutput, output)
		})
	}
}

func TestParseLastElement(t *testing.T) {
	tests := map[string]struct {
		input           string
		expectedElement string
		expectedOrderBy string
	}{
		"Case1: only contains order by": {
			input:           "Order by TestInt DESC",
			expectedElement: "",
			expectedOrderBy: "Order by TestInt DESC",
		},
		"Case2: only contains order by": {
			input:           "TestString = 'cannot be used in order by'",
			expectedElement: "TestString = 'cannot be used in order by'",
			expectedOrderBy: "",
		},
		"Case3: not contains any order by": {
			input:           "TestInt = 1",
			expectedElement: "TestInt = 1",
			expectedOrderBy: "",
		},
		"Case4-1: with order by in string & real order by": {
			input:           "TestString = 'cannot be used in order by' Order by TestInt DESC",
			expectedElement: "TestString = 'cannot be used in order by'",
			expectedOrderBy: "Order by TestInt DESC",
		},
		"Case4-2: with non-string attribute & real order by": {
			input:           "TestDouble = 1.0 Order by TestInt DESC",
			expectedElement: "TestDouble = 1.0",
			expectedOrderBy: "Order by TestInt DESC",
		},
		"Case5: with random case order by": {
			input:           "TestString = 'cannot be used in OrDer by' ORdeR by TestInt DESC",
			expectedElement: "TestString = 'cannot be used in OrDer by'",
			expectedOrderBy: "ORdeR by TestInt DESC",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			element, orderBy := parseOrderBy(test.input)
			assert.Equal(t, test.expectedElement, element)
			assert.Equal(t, test.expectedOrderBy, orderBy)
		})
	}
}

func TestSplitElement(t *testing.T) {
	tests := map[string]struct {
		input       string
		expectedKey string
		expectedVal string
		expectedOp  string
	}{
		"Case1-1: A=B": {
			input:       "CustomizedTestField=Test",
			expectedKey: "CustomizedTestField",
			expectedVal: "Test",
			expectedOp:  "=",
		},
		"Case1-2: A=\"B\"": {
			input:       "CustomizedTestField=\"Test\"",
			expectedKey: "CustomizedTestField",
			expectedVal: "\"Test\"",
			expectedOp:  "=",
		},
		"Case1-3: A='B'": {
			input:       "CustomizedTestField='Test'",
			expectedKey: "CustomizedTestField",
			expectedVal: "'Test'",
			expectedOp:  "=",
		},
		"Case2: A<=B": {
			input:       "CustomizedTestField<=Test",
			expectedKey: "CustomizedTestField",
			expectedVal: "Test",
			expectedOp:  "<=",
		},
		"Case3: A>=B": {
			input:       "CustomizedTestField>=Test",
			expectedKey: "CustomizedTestField",
			expectedVal: "Test",
			expectedOp:  ">=",
		},
		"Case4: A = B": {
			input:       "CustomizedTestField = Test",
			expectedKey: "CustomizedTestField",
			expectedVal: "Test",
			expectedOp:  "=",
		},
		"Case5: A <= B": {
			input:       "CustomizedTestField <= Test",
			expectedKey: "CustomizedTestField",
			expectedVal: "Test",
			expectedOp:  "<=",
		},
		"Case6: A >= B": {
			input:       "CustomizedTestField >= Test",
			expectedKey: "CustomizedTestField",
			expectedVal: "Test",
			expectedOp:  ">=",
		},
		"Case7: A > B": {
			input:       "CustomizedTestField > Test",
			expectedKey: "CustomizedTestField",
			expectedVal: "Test",
			expectedOp:  ">",
		},
		"Case8: A < B": {
			input:       "CustomizedTestField < Test",
			expectedKey: "CustomizedTestField",
			expectedVal: "Test",
			expectedOp:  "<",
		},
		"Case9: empty": {
			input:       "",
			expectedKey: "",
			expectedVal: "",
			expectedOp:  "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			key, val, op := splitElement(test.input)
			assert.Equal(t, test.expectedKey, key)
			assert.Equal(t, test.expectedVal, val)
			assert.Equal(t, test.expectedOp, op)
		})
	}
}

func TestIsTimeStruct(t *testing.T) {
	var emptyInput []byte
	numberInput := []byte("1709601210000000000")
	errorInput := []byte("Not a timeStamp")
	testTime := time.UnixMilli(1709601210000)
	var legitInput []byte
	legitInput, err := json.Marshal(testTime)
	assert.NoError(t, err)
	legitOutput := testTime.UnixMilli()
	legitOutputJSON, _ := json.Marshal(legitOutput)

	tests := map[string]struct {
		input          []byte
		expectedOutput []byte
		expectedError  error
	}{
		"Case1: empty input": {
			input:          emptyInput,
			expectedOutput: nil,
			expectedError:  nil,
		},
		"Case2: error input": {
			input:          errorInput,
			expectedOutput: errorInput,
			expectedError:  nil,
		},
		"Case3: number input": {
			input:          numberInput,
			expectedOutput: numberInput,
			expectedError:  nil,
		},
		"Case4: legit input": {
			input:          legitInput,
			expectedOutput: legitOutputJSON,
			expectedError:  nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualOutput, actualError := isTimeStruct(test.input)
			assert.Equal(t, test.expectedOutput, actualOutput)
			assert.Equal(t, test.expectedError, actualError)
		})
	}
}

func TestClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPinotClient := pnt.NewMockGenericClient(ctrl)
	mockProducer := &mocks.KafkaProducer{}
	mgr := NewPinotVisibilityStore(mockPinotClient, &service.Config{
		ValidSearchAttributes:  dynamicconfig.GetMapPropertyFn(definition.GetDefaultIndexedKeys()),
		ESIndexMaxResultWindow: dynamicconfig.GetIntPropertyFn(3),
	}, mockProducer, testlogger.New(t))
	visibilityStore := mgr.(*pinotVisibilityStore)

	assert.NotPanics(t, func() {
		visibilityStore.Close()
	})
}