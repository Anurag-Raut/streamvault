package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool

func Connect() {
	var err error
	pool, err = pgxpool.Connect(context.Background(), "host=localhost user=postgres password=postgres dbname=streamvault sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}

func Disconnect() {
	pool.Close()
}

func sendError(w http.ResponseWriter, err string) {
	type Error struct {
		Error string `json:"error"`
	}
	errorResponse := Error{
		Error: err,
	}
	result, _ := json.MarshalIndent(errorResponse, "", "  ")
	w.Write(result)
}

func AddStream(title, description, category, thumbnail, userId string) (string, error) {
	ctx := context.Background()
	var id string
	err := pool.QueryRow(ctx,
		`INSERT INTO "Video" (title, description, category, thumbnail, "isStreaming", "userId")
		 VALUES ($1, $2, $3, $4, false, $5)
		 RETURNING id`,
		title, description, category, thumbnail, userId).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func GetStreams(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	rows, err := pool.Query(ctx,
		`SELECT v.id, v.title, v.description, v.category, v.thumbnail, v."isStreaming", u.username ,u.id AS "userId"
		 FROM "Video" v
		 JOIN "User" u ON v."userId" = u.id`)
	if err != nil {
		sendError(w, err.Error())
		return
	}
	defer rows.Close()

	streams := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, title, description, category, thumbnail string
		var isStreaming bool
		var user struct {
			Username string `json:"username"`
			UserId   string `json:"userId"`
		}

		err = rows.Scan(&id, &title, &description, &category, &thumbnail, &isStreaming, &user.Username, &user.UserId)
		if err != nil {
			sendError(w, err.Error())
			return
		}
		fmt.Println(user, "user")
		stream := map[string]interface{}{
			"id":          id,
			"title":       title,
			"description": description,
			"category":    category,
			"thumbnail":   thumbnail,
			"isStreaming": isStreaming,
			"user":        user,
		}
		streams = append(streams, stream)
	}
	if err = rows.Err(); err != nil {
		sendError(w, err.Error())
		return
	}

	result, _ := json.MarshalIndent(streams, "", "  ")
	w.Write(result)
}

func UserExists(userId string) (bool, error) {
	ctx := context.Background()
	var count int
	err := pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM "User" WHERE id = $1`,
		userId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CreateUser(username string) (string, error) {
	ctx := context.Background()
	var id string
	err := pool.QueryRow(ctx,
		`INSERT INTO "User" (username)
		 VALUES ($1)
		 RETURNING "id"`,
		username).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func UpdateStatus(streamId string, status bool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx,
		`UPDATE videos
		 SET is_streaming = $1
		 WHERE id = $2`,
		status, streamId)
	return err
}

func GetUserId(w http.ResponseWriter, r *http.Request) {
	var userIdd *string
	var userIdResponse struct {
		UserId *string `json:"userId"`
	}

	cookie, err := r.Cookie("jwt")
	if err != nil {
		userIdResponse.UserId = userIdd
		response, _ := json.MarshalIndent(userIdResponse, "", "  ")
		w.Write(response)
		// http.Error(w, "No token found", http.StatusUnauthorized)
		return
	}
	tokenString := cookie.Value

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// You should provide the secret key or the key used for signing the token here
		return []byte("eat shit"), nil
	})

	if err != nil {
		userIdResponse.UserId = userIdd
		response, _ := json.MarshalIndent(userIdResponse, "", "  ")
		w.Write(response)
		return
	}

	if !token.Valid {
		userIdResponse.UserId = userIdd
		response, _ := json.MarshalIndent(userIdResponse, "", "  ")
		w.Write(response)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Access the username claim
		if userId, exists := claims["userId"].(string); exists {
			// Now you have the username
			fmt.Println("userId:", userId)
			userExists, _ := UserExists(userId)

			if !userExists {
				userIdResponse.UserId = userIdd
				response, _ := json.MarshalIndent(userIdResponse, "", "  ")
				w.Write(response)
				return
			}

			userIdd = &userId
			userIdResponse.UserId = userIdd
			response, _ := json.MarshalIndent(userIdResponse, "", "  ")
			w.Write(response)

		} else {
			userIdResponse.UserId = userIdd
			response, _ := json.MarshalIndent(userIdResponse, "", "  ")
			w.Write(response)
			return
		}
	} else {
		userIdResponse.UserId = userIdd
		response, _ := json.MarshalIndent(userIdResponse, "", "  ")
		w.Write(response)
		return
	}

}

func GetVideoData(w http.ResponseWriter, r *http.Request) {
	var videoId string
	err := json.NewDecoder(r.Body).Decode(&videoId)
	if err != nil {
		sendError(w, "Error decoding videoId")
		return
	}
	ctx := context.Background()
	userId := r.Context().Value("userId")

	type VideoData struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Thumbnail   string `json:"thumbnail"`
		IsStreaming bool   `json:"isStreaming"`
		User        struct {
			Username string `json:"username"`
			ID       string `json:"id"`
		} `json:"user"`
	}
	var videoData VideoData

	err = pool.QueryRow(ctx,
		`SELECT v.id, v.title, v.description, v.category, v.thumbnail, v."isStreaming", u.username, u.id
		 FROM "Video" v
		 JOIN "User" u ON v."userId" = u."id"
		 WHERE v.id = $1`,
		videoId).Scan(&videoData.ID, &videoData.Title, &videoData.Description, &videoData.Category, &videoData.Thumbnail, &videoData.IsStreaming, &videoData.User.Username, &videoData.User.ID)
	if err != nil {
		fmt.Println(err)
		sendError(w, "Error fetching video data")
		return
	}
	// fmt.Println(videoData, "videoData")
	var res struct {
		Likes    int `json:"likes"`
		Dislikes int `json:"dislikes"`
	}

	err = pool.QueryRow(ctx,
		`SELECT
		COALESCE(SUM(CASE WHEN "isLike" = true THEN 1 ELSE 0 END), 0) AS likes,
		COALESCE(SUM(CASE WHEN "isLike" = false THEN 1 ELSE 0 END), 0) AS dislikes
		 FROM "Like"
		 WHERE "videoId" = $1`,
		videoId).Scan(&res.Likes, &res.Dislikes)
	if err != nil {
		fmt.Println(err)
		sendError(w, "Error fetching likes")
		return
	}

	var subscriberRes struct {
		Subscribers int `json:"subscribers"`
	}

	err = pool.QueryRow(ctx,
		`SELECT 
			COALESCE(
				(SELECT COUNT(*) AS subscribers
				 FROM "Subscription"
				 WHERE "creatorId" = (SELECT "userId" FROM "Video" WHERE id = $1)
				),
				0
			) AS subscribers;
		`,
		videoId).Scan(&subscriberRes.Subscribers)
	if err != nil {
		fmt.Println(err)
		sendError(w, "Error fetching subscribers")
		return
	}

	const (
		Neutral = 0
		Like    = 1
		Dislike = 2
	)

	var response struct {
		VideoData
		Likes       int  `json:"likes"`
		Dislikes    int  `json:"dislikes"`
		Likestate   int  `json:"likeState"`
		Subscribers int  `json:"subscribers"`
		Subscribed  bool `json:"isSubscribed"`
	}

	response.VideoData = videoData
	response.Likes = res.Likes
	response.Dislikes = res.Dislikes
	response.Subscribers = subscriberRes.Subscribers

	if userId != nil {
		err = pool.QueryRow(ctx,
			`SELECT CASE
					WHEN EXISTS (
						SELECT 1
						FROM "Like"
						WHERE "videoId" = $1 AND "userId" = $2
					) THEN
						CASE
							WHEN (
								SELECT "isLike"
								FROM "Like"
								WHERE "videoId" = $1 AND "userId" = $2
							) THEN 1  
							ELSE 2  
						END
					ELSE 0  
				END AS "likeState";
	`, videoId, userId).Scan(&response.Likestate)
		fmt.Println(videoData.User.ID, "creatorId", userId, "userId")

		if err != nil {
			fmt.Println(err)

			sendError(w, "Error fetching like state")
			return
		}

		err = pool.QueryRow(ctx,
			`SELECT CASE
			WHEN EXISTS (
				SELECT 1
				FROM "Subscription"
				WHERE "creatorId" = $1 AND "subscriberId" = $2
			) THEN true
			ELSE false
		END AS "isSubscribed";
	`, videoData.User.ID, userId).Scan(&response.Subscribed)
		fmt.Println(response.Subscribed, "subscribed")

		if err != nil {
			fmt.Println(err)
			sendError(w, "Error fetching subscription state")
			return
		}
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	w.Write(result)
}

func Like(w http.ResponseWriter, r *http.Request) {
	var videoId string
	err := json.NewDecoder(r.Body).Decode(&videoId)
	if err != nil {
		sendError(w, "error decoding videoId")
		return
	}
	userId := r.Context().Value("userId").(string)

	ctx := context.Background()

	_, err = pool.Exec(ctx,
		`INSERT INTO "Like" ("videoId", "userId", "isLike")
		 VALUES ($1, $2, true)
		 ON CONFLICT ("videoId", "userId") DO UPDATE SET "isLike" = true`,
		videoId, userId)

	if err != nil {
		fmt.Println(err)
		sendError(w, "error liking video")
		return
	}
	var response struct {
		Message string `json:"message"`
	}
	response.Message = "Liked"
	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		sendError(w, "error marshalling response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func Dislike(w http.ResponseWriter, r *http.Request) {
	var videoId string
	err := json.NewDecoder(r.Body).Decode(&videoId)
	if err != nil {
		sendError(w, "error decoding videoId")
		return
	}
	ctx := context.Background()
	userId := r.Context().Value("userId").(string)

	_, err = pool.Exec(ctx,
		`INSERT INTO "Like" ("videoId", "userId", "isLike")
		 VALUES ($1, $2, false)
		 ON CONFLICT ("videoId", "userId") DO UPDATE SET "isLike" = false`,
		videoId, userId)

	if err != nil {
		sendError(w, "error disliking video")
		return
	}
	var response struct {
		Message string `json:"message"`
	}
	response.Message = "Disliked"
	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		sendError(w, "error marshalling response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func RemoveLike(w http.ResponseWriter, r *http.Request) {
	var videoId string
	err := json.NewDecoder(r.Body).Decode(&videoId)
	if err != nil {
		sendError(w, "error decoding videoId")
		return
	}
	ctx := context.Background()
	userId := r.Context().Value("userId").(string)

	_, err = pool.Exec(ctx,
		`DELETE FROM "Like"
		 WHERE "videoId" = $1 AND "userId" = $2`,
		videoId, userId)

	if err != nil {
		sendError(w, "error removing like")
		return
	}
	var response struct {
		Message string `json:"message"`
	}
	response.Message = "Like removed"
	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		sendError(w, "error marshalling response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	var channelId string
	err := json.NewDecoder(r.Body).Decode(&channelId)
	if err != nil {
		sendError(w, "error decoding channelId")
		return
	}
	ctx := context.Background()
	userId := r.Context().Value("userId").(string)

	_, err = pool.Exec(ctx,
		`INSERT INTO "Subscription" ("creatorId", "subscriberId")
		 VALUES ($1, $2)
		 ON CONFLICT ("creatorId", "subscriberId") DO NOTHING`,
		channelId, userId)

	if err != nil {
		sendError(w, "error subscribing to channel")
		return
	}
	var response struct {
		Message string `json:"message"`
	}
	response.Message = "Subscribed"
	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		sendError(w, "error marshalling response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func Unsubscribe(w http.ResponseWriter, r *http.Request) {
	var channelId string
	err := json.NewDecoder(r.Body).Decode(&channelId)
	if err != nil {
		sendError(w, "error decoding channelId")
		return
	}
	ctx := context.Background()
	userId := r.Context().Value("userId").(string)

	_, err = pool.Exec(ctx,
		`DELETE FROM "Subscription"
		 WHERE "creatorId" = $1 AND "subscriberId" = $2`,
		channelId, userId)

	if err != nil {
		sendError(w, "error unsubscribing to channel")
		return
	}
	var response struct {
		Message string `json:"message"`
	}
	response.Message = "Unsubscribed"
	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		sendError(w, "error marshalling response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func GetContent(w http.ResponseWriter, r *http.Request) {
	var userId = r.Context().Value("userId").(string)
	var pageSize = 10
	var pageNumber = 1
	var err error = nil
	pageSize, err = strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		fmt.Println(err)
		pageSize = 10
	}

	pageNumber, err = strconv.Atoi(r.URL.Query().Get("pageNumber"))
	if err != nil {
		fmt.Println(err)
		pageNumber = 1
	}

	fmt.Println(pageSize, pageNumber, "pageSize,pageNumber")

	type Content struct {
		Id          string    `json:"id"`
		Thumbnail   string    `json:"thumbnail"`
		Title       string    `json:"title"`
		CreatedAt   time.Time `json:"createdAt"`
		Likes       int       `json:"likes"`
		Comments    int       `json:"comments"`
		Description string    `json:"description"`
		Category    string    `json:"category"`
		Dislikes    int       `json:"dislikes"`
		Subscribers int       `json:"subscribers"`
	}

	ctx := context.Background()
	rows, err := pool.Query(ctx, `
    SELECT
        v.id, v.title, v.description, v.category, v.thumbnail, v."createdAt",
        (SELECT COUNT(*) FROM "Like" WHERE "videoId" = v.id AND "isLike" = true) AS likes,
        (SELECT COUNT(*) FROM "Like" WHERE "videoId" = v.id AND "isLike" = false) AS dislikes,
        (SELECT COUNT(*) FROM "Subscription" WHERE "creatorId" = $1) AS subscribers,
        (SELECT COUNT(*) FROM "Comment" WHERE "videoId" = v.id) AS comments
    FROM "Video" v
	ORDER BY v."createdAt" DESC
    LIMIT $2 OFFSET $3

	`, userId, pageSize, (pageNumber-1)*pageSize)

	if err != nil {
		sendError(w, err.Error())
		return
	}
	defer rows.Close()

	var response []Content
	for rows.Next() {
		// var id, title, description, category, thumbnail string
		// var likes, comments, dislikes, subscribers int
		// var createdAt time.Time
		var videoData Content
		err = rows.Scan(&videoData.Id, &videoData.Title, &videoData.Description, &videoData.Category, &videoData.Thumbnail, &videoData.CreatedAt, &videoData.Likes, &videoData.Dislikes, &videoData.Subscribers, &videoData.Comments)
		if err != nil {
			sendError(w, err.Error())
			return
		}
		response = append(response, videoData)

	}
	if err = rows.Err(); err != nil {
		sendError(w, err.Error())
		return
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	w.Write(result)

}

func GetChats(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var request struct {
		VideoId string `json:"videoId"`
	}
	err:=json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println("erorr decoding",err)
		sendError(w, err.Error())
		return
	}

	rows, err := pool.Query(ctx,
		`SELECT c."userId",c.text,c."createdAt"
		 FROM "Comment" c
		 Where c."videoId" = $1
		 ORDER BY c."createdAt" DESC
		`,request.VideoId)
	if err != nil {
		sendError(w, err.Error())
		return
	}
	defer rows.Close()

	type Chat struct {
		UserId    string    `json:"userId"`
		Message   string    `json:"message"`
		CreatedAt time.Time `json:"createdAt"`

	}
	chats := make([]Chat, 0)
	for rows.Next() {
		var chat Chat
		err = rows.Scan( &chat.UserId, &chat.Message, &chat.CreatedAt)
		if err != nil {
			sendError(w, err.Error())
			return
		}
		chats = append(chats, chat)
	}
	if err = rows.Err(); err != nil {
		sendError(w, err.Error())
		return
	}

	result, _ := json.MarshalIndent(chats, "", "  ")
	w.Write(result)
}


func PostChat( videoId string,  userId string , message string ) error {
	ctx := context.Background()
	_,err := pool.Exec(ctx,
		`INSERT INTO "Comment" ("videoId", "userId", "text")
		 VALUES ($1, $2, $3)`,
		videoId, userId, message)
	if err != nil {
		fmt.Println(err,"reallly")
		return  err
	}
	return nil


}

type UserDetails struct {
	Username string `json:"username"`
	ProfileImage *string `json:"profileImage"`
	UserId string `json:"userId"`
}
func GetUserDetailsFromDatabase(userId string) (UserDetails,error){

	var userDetails UserDetails

	ctx := context.Background()
	err := pool.QueryRow(ctx,
		`SELECT username, "profileImage", id
		 FROM "User"
		 WHERE id = $1`,
		userId).Scan(&userDetails.Username, &userDetails.ProfileImage, &userDetails.UserId)
	if err != nil {
		fmt.Println(err)
		return userDetails,err
	}
	return userDetails,nil


}