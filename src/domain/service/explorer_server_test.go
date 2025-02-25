package service

import (
	"testing"
	"time"

	"context"

	"github.com/lokker96/grpc_project/domain/entity"
	"github.com/lokker96/grpc_project/infrastructure/proto/explore"
	repository_mock "github.com/lokker96/grpc_project/mocks/github.com/lokker96/grpc_project/domain/repository"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/mock"
)

type inputData struct {
	ctx     context.Context
	request *explore.ListLikedYouRequest
}

type expectation struct {
	response *explore.ListLikedYouResponse
	err      error
}

type mockData struct {
	dbDecisions []entity.Decision
	dbError     error
}

type testCaseData struct {
	inputData    inputData
	expectations expectation
	mocksData    mockData
}

func Test_ListLikedYou(t *testing.T) {
	testCases := []testCaseData{
		{
			inputData: inputData{
				ctx: context.Background(),
				request: &explore.ListLikedYouRequest{
					RecipientUserId: "1",
					PaginationToken: nil,
				},
			},
			expectations: expectation{
				response: &explore.ListLikedYouResponse{
					Likers: []*explore.ListLikedYouResponse_Liker{
						{
							ActorId:       "2",
							UnixTimestamp: 0,
						},
					},
					NextPaginationToken: nil,
				},
				err: nil,
			},
			mocksData: mockData{
				dbDecisions: []entity.Decision{
					{
						ID:          1,
						AuthorID:    2,
						RecipientID: 1,
						Liked:       true,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				},
				dbError: nil,
			},
		},
		{
			// Add more tests and improve testing functions
		},
	}

	for _, testCase := range testCases {
		listLikedYou_tester_func(t, testCase)
	}
}

func listLikedYou_tester_func(t *testing.T, testCase testCaseData) {
	repositoryMock := &repository_mock.MockExplorerRepository{}

	repositoryMock.
		On("GetDecisionsForRecipientId", mock.AnythingOfType("int"), mock.AnythingOfType("*bool")).
		Once().Return(testCase.mocksData.dbDecisions, testCase.mocksData.dbError)

	explorerService := NewExplorerServer(repositoryMock)

	response, err := explorerService.ListLikedYou(testCase.inputData.ctx, testCase.inputData.request)

	for i, like := range response.Likers {
		assert.Equal(t, like.ActorId, testCase.expectations.response.Likers[i].ActorId)

		// Unnecessary but would be nice to mock
		// assert.Equal(like.UnixTimestamp, testCase.expectations.response.Likers[i].UnixTimestamp)
	}

	assert.Equal(t, response.NextPaginationToken, testCase.expectations.response.NextPaginationToken)

	assert.Equal(t, err, testCase.expectations.err)

	repositoryMock.AssertExpectations(t)
}
