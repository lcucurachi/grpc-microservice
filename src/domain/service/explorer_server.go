package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/lokker96/grpc_project/domain/entity"
	"github.com/lokker96/grpc_project/domain/repository"
	ep "github.com/lokker96/grpc_project/infrastructure/proto/explore"
)

// Embeds the gRPC server that provides the endpoints and implements them
type ExploreServer struct {
	ep.UnimplementedExploreServiceServer
	explorerRepository repository.ExplorerRepository // Explorer repository which implements a postgres DB and method to access the data
}

func NewExplorerServer(explorerRepository repository.ExplorerRepository) *ExploreServer {
	return &ExploreServer{explorerRepository: explorerRepository}
}

// Helper function for making testing easier
// Dataset:
// User IDs: [1, 2, 3, 4]
// Like: 1 -> 2
// Like: 2 -> 1
// Like: 4 -> 1
func (s *ExploreServer) BuildDummyDataset() {
	s.explorerRepository.CreateUser(&entity.User{})
	s.explorerRepository.CreateUser(&entity.User{})
	s.explorerRepository.CreateUser(&entity.User{})
	s.explorerRepository.CreateUser(&entity.User{})

	s.explorerRepository.CreateDecision(&entity.Decision{
		AuthorID:    1,
		RecipientID: 2,
		Liked:       true,
	})

	s.explorerRepository.CreateDecision(&entity.Decision{
		AuthorID:    2,
		RecipientID: 1,
		Liked:       true,
	})

	s.explorerRepository.CreateDecision(&entity.Decision{
		AuthorID:    4,
		RecipientID: 1,
		Liked:       true,
	})
}

func (s *ExploreServer) ListLikedYou(ctx context.Context, request *ep.ListLikedYouRequest) (*ep.ListLikedYouResponse, error) {
	likers := make([]*ep.ListLikedYouResponse_Liker, 0)

	recipientUserID, err := strconv.Atoi(request.GetRecipientUserId())
	if err != nil {
		return nil, fmt.Errorf("error converting recipient user id string: %w", err)
	}

	// Ideally we should check that the recipient user id exists first by calling a method to check

	liked := true

	decisions, err := s.explorerRepository.GetDecisionsForRecipientId(recipientUserID, &liked)
	if err != nil {
		return nil, fmt.Errorf("error getting liked decisions for recipient id: %w", err)
	}

	for _, dec := range decisions {
		likers = append(likers, &ep.ListLikedYouResponse_Liker{
			ActorId:       strconv.Itoa(int(dec.AuthorID)),
			UnixTimestamp: uint64(time.Now().Unix()),
			// Is this the timestamp when the profile liked the user?
			// If yes, I would get it from the DB using the updated_at field from the decision table
		})
	}

	return &ep.ListLikedYouResponse{
		Likers: likers,
	}, nil
}

func (s *ExploreServer) ListNewLikedYou(ctx context.Context, request *ep.ListLikedYouRequest) (*ep.ListLikedYouResponse, error) {
	likers := make([]*ep.ListLikedYouResponse_Liker, 0)

	recipientUserID, err := strconv.Atoi(request.GetRecipientUserId())
	if err != nil {
		return nil, fmt.Errorf("error converting recipient user id string: %w", err)
	}

	// Ideally we should check that the recipient user id exists first by calling a method to check

	liked := true

	userDecisionsLiked, err := s.explorerRepository.GetDecisionsForUserId(recipientUserID, &liked)
	if err != nil {
		return nil, fmt.Errorf("error getting user decisions list: %w", err)
	}

	recipientDecisionsLiked, err := s.explorerRepository.GetDecisionsForRecipientId(recipientUserID, &liked)
	if err != nil {
		return nil, fmt.Errorf("error getting liked decisions for user id: %w", err)
	}

	usersLiked := map[int]bool{}
	for _, userLike := range userDecisionsLiked {
		usersLiked[int(userLike.RecipientID)] = true
	}

	// This is bad and I would not reccomend to use it in production.
	// We need to use a more specific SQL query with JOINS to simplify this process.
	// Basically, here we are checking that there isn't a mutual like becasue
	// we want new users who liked the recipient.
	for _, recipientLike := range recipientDecisionsLiked {
		_, ok := usersLiked[int(recipientLike.AuthorID)]
		if !ok {
			likers = append(likers, &ep.ListLikedYouResponse_Liker{
				ActorId: strconv.Itoa(int(recipientLike.AuthorID)),
			})
		}
	}

	return &ep.ListLikedYouResponse{
		Likers: likers,
	}, nil
}

func (s *ExploreServer) CountLikedYou(ctx context.Context, request *ep.CountLikedYouRequest) (*ep.CountLikedYouResponse, error) {
	recipientUserID, err := strconv.Atoi(request.GetRecipientUserId())
	if err != nil {
		return nil, fmt.Errorf("error converting recipient user id string: %w", err)
	}

	// Ideally we should check that the recipient user id exists first by calling a method to check
	result := s.explorerRepository.GetLikesCountByProfileId(recipientUserID)

	return &ep.CountLikedYouResponse{
		Count: uint64(result),
	}, nil
}

func (s *ExploreServer) PutDecision(ctx context.Context, request *ep.PutDecisionRequest) (*ep.PutDecisionResponse, error) {
	actorUserId, err := strconv.Atoi(request.GetActorUserId())
	if err != nil {
		return nil, fmt.Errorf("error converting actor user id string: %w", err)
	}

	recipientUserId, err := strconv.Atoi(request.GetRecipientUserId())
	if err != nil {
		return nil, fmt.Errorf("error converting recipient user id string: %w", err)
	}

	// Ideally we should check that both the user ids exists before calling this.
	// This is especially true if this routine is called from other microservices and not only by the user app.
	err = s.explorerRepository.UpdateDecision(actorUserId, recipientUserId, request.GetLikedRecipient())
	if err != nil {
		return nil, fmt.Errorf("error putting decision: %w", err)
	}

	mutualLikes := s.explorerRepository.FindMutualLike(actorUserId, recipientUserId)

	return &ep.PutDecisionResponse{
		MutualLikes: mutualLikes,
	}, nil
}
