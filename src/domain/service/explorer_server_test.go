package service

import (
	"errors"
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

	nowTime := time.Now()

	testCases := []testCaseData{
		// User 1 likers: user id 2
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
							UnixTimestamp: uint64(nowTime.Unix()),
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
						CreatedAt:   nowTime,
						UpdatedAt:   nowTime,
					},
				},
				dbError: nil,
			},
		},
		// User 1 likers returns error
		{
			inputData: inputData{
				ctx: context.Background(),
				request: &explore.ListLikedYouRequest{
					RecipientUserId: "1",
					PaginationToken: nil,
				},
			},
			expectations: expectation{
				response: nil,
				err:      errors.New("error getting liked decisions for recipient id: Error executing query"),
			},
			mocksData: mockData{
				dbDecisions: nil,
				dbError:     errors.New("Error executing query"),
			},
		},
	}
	// Add more tests and improve testing functions

	for _, testCase := range testCases {
		listLikedYou_tester_func(t, testCase)
	}
}

func listLikedYou_tester_func(t *testing.T, testCase testCaseData) {
	repositoryMock := &repository_mock.MockExplorerRepository{}

	repositoryMock.
		On("GetDecisionsForRecipientId", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("*bool")).
		Once().Return(testCase.mocksData.dbDecisions, testCase.mocksData.dbError)

	explorerService := NewExplorerServer(repositoryMock)

	response, err := explorerService.ListLikedYou(testCase.inputData.ctx, testCase.inputData.request)

	if testCase.expectations.err != nil && err == nil {
		t.FailNow()
	} else if testCase.expectations.err == nil && err != nil {
		t.FailNow()
	}

	if testCase.expectations.err != nil && err != nil {
		assert.Equal(t, err.Error(), testCase.expectations.err.Error())
		assert.Equal(t, response, testCase.expectations.response)
	} else {
		assert.Equal(t, len(response.Likers), len(testCase.expectations.response.Likers))

		if len(response.Likers) == len(testCase.expectations.response.Likers) {
			for i, like := range response.Likers {
				assert.Equal(t, like.ActorId, testCase.expectations.response.Likers[i].ActorId)
				assert.Equal(t, like.UnixTimestamp, uint64(testCase.mocksData.dbDecisions[i].UpdatedAt.Unix()))
			}
		}

		assert.Equal(t, response.NextPaginationToken, testCase.expectations.response.NextPaginationToken)
	}

	repositoryMock.AssertExpectations(t)
}
