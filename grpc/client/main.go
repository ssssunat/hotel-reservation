package main

import (
	// ...
	"context"
	"fmt"
	pb "grpc/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)


func main() {
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewUsersClient(conn)

	TestUsers(c)
}

func TestUsers(c pb.UsersClient) {
    users := []*pb.User{
        {Name: "Сергей", Email: "serge@example.com", Sex: pb.User_MALE},
        {Name: "Света", Email: "sveta@example.com", Sex: pb.User_FEMALE},
        {Name: "Денис", Email: "den@example.com", Sex: pb.User_MALE},
        // при добавлении этой записи должна вернуться ошибка:
        // пользователь с email sveta@example.com уже существует
        {Name: "Sveta", Email: "sveta@example.com", Sex: pb.User_FEMALE},
    }

	for _, user := range users {
		resp, err := c.AddUser(context.Background(), &pb.AddUserRequest{
			User: user,
		})
		if err != nil {
			log.Fatal(err)
		}
		if resp.Error != "" {
			fmt.Println(resp.Error)
		}
	}

	resp, err := c.DelUser(context.Background(), &pb.DelUserRequest{
		Email: "serge@example.com",
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp.Error != "" {
		fmt.Println(resp.Error)
	}

	for _, userEmail := range []string{"sveta@example.com", "serge@example.com"} {
        resp, err := c.GetUser(context.Background(), &pb.GetUserRequest{
            Email: userEmail,
        })
        if err != nil {
            if e, ok := status.FromError(err); ok {
				if e.Code() == codes.NotFound {
					fmt.Println(`NOT FOUND`, e.Message())
				} else {
					fmt.Println(e.Code(), e.Message())
				}
			}
        } else {
			fmt.Printf("Не получилось распарсить ошибку %v", err)
		} 
        if resp.Error == "" {
            fmt.Println(resp.User)
        } else {
            fmt.Println(resp.Error)
        }
    }

    // получаем список email пользователей
    emails, err := c.ListUsers(context.Background(), &pb.ListUsersRequest{
        Offset: 0,
        Limit:  100,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(emails.Count, emails.Emails)
}

