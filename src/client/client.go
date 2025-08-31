package main

import (
	"context"
	"fmt"
	"log"
	"time"

	ep "github.com/lokker96/grpc_project/infrastructure/proto/explore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create new connection to localhost with insecure credentials for development purposes
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to grpc server")
	}

	// defer closing connection for later
	defer conn.Close()

	// make new client for explore gRPC service calls
	client := ep.NewExploreServiceClient(conn)

	// create new context with 1 second timeout which should be plenty for this exercise
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	// Check how many users like user id 1
	count, err := client.CountLikedYou(ctx, &ep.CountLikedYouRequest{
		RecipientUserId: "1",
	})
	if err != nil {
		log.Fatal("error calling function CountLikedYou: %w", err)
	}

	fmt.Println("CountLikedYou - RecipientUserId: 1, Count: ", count.Count)

	// Check how many new users like user id 1
	listNewLikedYouResponse, err := client.ListNewLikedYou(ctx, &ep.ListLikedYouRequest{
		RecipientUserId: "1",
	})
	if err != nil {
		log.Fatal("error calling function ListNewLikedYou: %w", err)
	}

	fmt.Println("\nListNewLikedYou - RecipientUserId: 1, ids: ")
	for _, like := range listNewLikedYouResponse.Likers {
		fmt.Print(like.ActorId, " - ", like.UnixTimestamp) // Not formatting timestamp, just checking that it returns
		fmt.Printf("\n")
	}
	fmt.Printf("\n")

	// Check how many users like user id 1
	listLikedYouResponse, err := client.ListLikedYou(ctx, &ep.ListLikedYouRequest{
		RecipientUserId: "1",
	})
	if err != nil {
		log.Fatal("error calling function ListLikedYou: %w", err)
	}

	fmt.Println("ListLikedYou - RecipientUserId: 1, ids: ")
	for _, like := range listLikedYouResponse.Likers {
		fmt.Print(like.ActorId, " - ", like.UnixTimestamp)
		fmt.Printf("\n")
	}
	fmt.Printf("\n")

	// Put decisions user id 3 likes user id 1
	putDecisionResponse, err := client.PutDecision(ctx, &ep.PutDecisionRequest{
		ActorUserId:     "3",
		RecipientUserId: "1",
		LikedRecipient:  true,
	})
	if err != nil {
		log.Fatal("error calling function ListLikedYou: %w", err)
	}

	fmt.Println("\nActorId: 3 likes RecipientUserId: 1")
	fmt.Println("Mutual Like - ActorId: 3 and RecipientUserId: 1, response: ", putDecisionResponse.MutualLikes)

	// Test that we can alter the decisions for user id 3 when he does not like user id 1
	putDecisionResponse, err = client.PutDecision(ctx, &ep.PutDecisionRequest{
		ActorUserId:     "3",
		RecipientUserId: "1",
		LikedRecipient:  false,
	})
	if err != nil {
		log.Fatal("error calling function ListLikedYou: %w", err)
	}

	fmt.Println("\nActorId: 3 does not likes RecipientUserId: 1")
	fmt.Println("Mutual Like - ActorId: 3 and RecipientUserId: 1, response: ", putDecisionResponse.MutualLikes)

	// Test that we can alter the decisions for user id 1 when he does not like user id 2
	// and that we get the mutual like equal to true
	putDecisionResponse, err = client.PutDecision(ctx, &ep.PutDecisionRequest{
		ActorUserId:     "1",
		RecipientUserId: "2",
		LikedRecipient:  true,
	})
	if err != nil {
		log.Fatal("error calling function ListLikedYou: %w", err)
	}

	fmt.Println("\nActorId: 1 likes RecipientUserId: 2")
	fmt.Println("Mutual Like - ActorId: 1 and RecipientUserId: 2, response: ", putDecisionResponse.MutualLikes)

	// Check how many new users like user id 4
	listNewLikedYouResponse, err = client.ListNewLikedYou(ctx, &ep.ListLikedYouRequest{
		RecipientUserId: "4",
	})
	if err != nil {
		log.Fatal("error calling function ListNewLikedYou: %w", err)
	}

	fmt.Println("\nListNewLikedYou - RecipientUserId: 4, ids: ")
	for _, like := range listNewLikedYouResponse.Likers {
		fmt.Print(like.ActorId, " - ", like.UnixTimestamp)
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}
