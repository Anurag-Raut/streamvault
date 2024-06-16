package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"streamvault/rmq"
	"time"

	"github.com/gofor-little/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type VideoData struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Thumbnail   string `json:"thumbnail"`
	IsStreaming bool   `json:"isStreaming"`
	User        struct {
		Username     string  `json:"username"`
		ID           string  `json:"id"`
		ProfileImage *string `json:"profileImage"`
	} `json:"user"`
}

type Chat struct {
	Message   string      `json:"message"`
	CreatedAt time.Time   `json:"createdAt"`
	User      UserDetails `json:"user"`
}

var pool *pgxpool.Pool

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	fmt.Println(hash)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Println(err, "errrr", hash)
	return err == nil
}

func Connect() {
	var err error
	pool, err = pgxpool.Connect(context.Background(), env.Get("DATABASE_URL","host=database user=postgres password=postgres dbname=streamvault sslmode=disable"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}

func Disconnect() {
	pool.Close()
}

func sendError(w http.ResponseWriter, err string, codes ...int) {
	code := http.StatusInternalServerError // Default value

	// Check if codes slice is not empty
	if len(codes) > 0 {
		code = codes[0]
	}
	type Error struct {
		Error string `json:"error"`
	}
	errorResponse := Error{
		Error: err,
	}
	result, _ := json.MarshalIndent(errorResponse, "", "  ")
	w.WriteHeader(code)
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
		`SELECT v.id, v.title, v.description, v.category, v.thumbnail, v."isStreaming",v."createdAt", u.username ,u.id,v."views",u."profileImage"
		 FROM "Video" v
		 JOIN "User" u ON v."userId" = u.id
		 WHERE v."isProcessed" = True
		 `)
	if err != nil {
		sendError(w, err.Error())
		return
	}
	defer rows.Close()
	type Video struct {
		Id          string      `json:"id"`
		Title       string      `json:"title"`
		Description string      `json:"description"`
		Category    string      `json:"category"`
		Thumbnail   string      `json:"thumbnail"`
		IsStreaming bool        `json:"isStreaming"`
		User        UserDetails `json:"user"`
		CreatedAt   time.Time   `json:"createdAt"`
		Views       int         `json:"views"`
	}
	videos := make([]Video, 0)
	for rows.Next() {
		var video Video

		err = rows.Scan(&video.Id, &video.Title, &video.Description, &video.Category, &video.Thumbnail, &video.IsStreaming, &video.CreatedAt, &video.User.Username, &video.User.UserId, &video.Views, &video.User.ProfileImage)
		if err != nil {
			fmt.Println(err, "errrurrr")
			sendError(w, err.Error())
			return
		}

		videos = append(videos, video)
	}
	if err = rows.Err(); err != nil {
		sendError(w, err.Error())
		return
	}

	result, _ := json.MarshalIndent(videos, "", "  ")
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

func CreateUser(username string, profileImage *string) (string, error) {
	ctx := context.Background()
	var id string

	// Check if the user already exists
	err := pool.QueryRow(ctx,
		`SELECT "id" FROM "User" WHERE username = $1`,
		username).Scan(&id)

	if err == pgx.ErrNoRows { // User doesn't exist, perform insert
		fmt.Println("User doesn't exist")
		err2 := pool.QueryRow(ctx,
			`INSERT INTO "User" (username,"profileImage")
             VALUES ($1,$2)
             RETURNING "id"`,
			username, profileImage).Scan(&id)
		if err2 != nil {
			fmt.Println(err2, "Error inserting user")
			return "", err
		}
	} else if err != nil { // Other error occurred
		fmt.Println(err, "Other error")
		return "", err
	}

	return id, nil
}

func UpdateStatus(streamId string, status bool) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx,
		`UPDATE video
		 SET is_streaming = $1
		 WHERE id = $2`,
		status, streamId)
	return err
}

func GetUserId(w http.ResponseWriter, r *http.Request) {
	var userIdd *string
	var userIdResponse struct {
		UserId       *string `json:"userId"`
		Username     string  `json:"username"`
		ProfileImage *string `json:"profileImage"`
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
			var userDetails UserDetails
			userDetails, err := GetUserDetailsFromDatabase(*userIdResponse.UserId)

			if err != nil {
				fmt.Println(err)
				userIdResponse.UserId = userIdd

				response, _ := json.MarshalIndent(userIdResponse, "", "  ")
				w.Write(response)
				return
			}
			fmt.Println(userDetails, "userDetails bbabbee")
			userIdResponse.Username = userDetails.Username
			userIdResponse.ProfileImage = userDetails.ProfileImage
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
			Username     string  `json:"username"`
			ID           string  `json:"id"`
			ProfileImage *string `json:"profileImage"`
		} `json:"user"`
	}
	var videoData VideoData

	err = pool.QueryRow(ctx,
		`SELECT v.id, v.title, v.description, v.category, v.thumbnail, v."isStreaming", u.username, u.id,u."profileImage"
		 FROM "Video" v
		 JOIN "User" u ON v."userId" = u."id"
		 WHERE v.id = $1`,
		videoId).Scan(&videoData.ID, &videoData.Title, &videoData.Description, &videoData.Category, &videoData.Thumbnail, &videoData.IsStreaming, &videoData.User.Username, &videoData.User.ID, &videoData.User.ProfileImage)
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

	_, err = pool.Exec(ctx,
		`UPDATE "Video"
		SET "views" = "views" + 1
		WHERE id = $1`, videoId)

	if err != nil {
		fmt.Println(err)
		sendError(w, "Error updating views")
		return
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	w.Write(result)
}

func GetVideoDataFromDatabase(videoId string) (VideoData, error) {

	ctx := context.Background()

	var videoData VideoData

	err := pool.QueryRow(ctx,
		`SELECT v.id, v.title, v.description, v.category, v.thumbnail, v."isStreaming", u.username, u.id,u."profileImage"
		 FROM "Video" v
		 JOIN "User" u ON v."userId" = u."id"
		 WHERE v.id = $1`,
		videoId).Scan(&videoData.ID, &videoData.Title, &videoData.Description, &videoData.Category, &videoData.Thumbnail, &videoData.IsStreaming, &videoData.User.Username, &videoData.User.ID, &videoData.User.ProfileImage)
	if err != nil {
		fmt.Println(err)

		return VideoData{}, err
	}

	return videoData, nil
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

// func GetSubscriptions(creatorId string) (int, error) {
// 	ctx := context.Background()
// 	var count int
// 	err := pool.QueryRow(ctx,
// 		`SELECT COUNT(*) FROM "Subscription" WHERE "creatorId" = $1`,
// 		creatorId).Scan(&count)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

func GetContent(w http.ResponseWriter, r *http.Request) {
	var username string
	var pageSize = 10
	var pageNumber = 1
	var err error = nil
	var vodConditions = "false,true"

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
	username = r.URL.Query().Get("username")

	isVOD, err := strconv.ParseBool(r.URL.Query().Get("isVOD"))
	if err != nil {
		fmt.Println(err)
		vodConditions = "FALSE,TRUE"

	} else {
		if isVOD {
			vodConditions = "TRUE"
		} else {
			vodConditions = "FALSE"
		}
	}

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
		Views       int       `json:"views"`
	}

	ctx := context.Background()
	rows, err := pool.Query(ctx,
		fmt.Sprintf(`
		SELECT
		v.id, v.title, v.description, v.category, v.thumbnail, v."createdAt",v."views",
		(SELECT COUNT(*) FROM "Like" WHERE "videoId" = v.id AND "isLike" = true) AS likes,
		(SELECT COUNT(*) FROM "Like" WHERE "videoId" = v.id AND "isLike" = false) AS dislikes,
		(SELECT COUNT(*) FROM "Subscription" WHERE "creatorId" = u.id) AS subscribers,
		(SELECT COUNT(*) FROM "Comment" WHERE "videoId" = v.id) AS comments
		FROM "Video" v
		JOIN "User" u ON v."userId" = u.id
	WHERE u.username = $1 AND v."isVOD" IN (%s)
	ORDER BY v."createdAt" DESC
	LIMIT $2 OFFSET $3
		`, vodConditions), username, pageSize, (pageNumber-1)*pageSize)

	if err != nil {
		sendError(w, err.Error())
		return
	}
	defer rows.Close()

	var response = make([]Content, 0)
	for rows.Next() {

		var videoData Content
		err = rows.Scan(&videoData.Id, &videoData.Title, &videoData.Description, &videoData.Category, &videoData.Thumbnail, &videoData.CreatedAt, &videoData.Views, &videoData.Likes, &videoData.Dislikes, &videoData.Subscribers, &videoData.Comments)
		if err != nil {
			sendError(w, err.Error())
			return
		}
		response = append(response, videoData)

	}
	fmt.Println(response, "response")
	if err = rows.Err(); err != nil {
		sendError(w, err.Error())
		return
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(result))
	w.Write(result)

}

func GetDashboardContent(w http.ResponseWriter, r *http.Request) {
	var userId string
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
	userId = r.Context().Value("userId").(string)

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
		Views       int       `json:"views"`
	}

	ctx := context.Background()
	rows, err := pool.Query(ctx,
		`
		SELECT
		v.id, v.title, v.description, v.category, v.thumbnail, v."createdAt",v."views",
		(SELECT COUNT(*) FROM "Like" WHERE "videoId" = v.id AND "isLike" = true) AS likes,
		(SELECT COUNT(*) FROM "Like" WHERE "videoId" = v.id AND "isLike" = false) AS dislikes,
		(SELECT COUNT(*) FROM "Subscription" WHERE "creatorId" = u.id) AS subscribers,
		(SELECT COUNT(*) FROM "Comment" WHERE "videoId" = v.id) AS comments
		FROM "Video" v
		JOIN "User" u ON v."userId" = u.id
	WHERE u."id" = $1
	ORDER BY v."createdAt" DESC
	LIMIT $2 OFFSET $3
		`, userId, pageSize, (pageNumber-1)*pageSize)

	if err != nil {
		sendError(w, err.Error())
		return
	}
	defer rows.Close()

	var response = make([]Content, 0)
	for rows.Next() {

		var videoData Content
		err = rows.Scan(&videoData.Id, &videoData.Title, &videoData.Description, &videoData.Category, &videoData.Thumbnail, &videoData.CreatedAt, &videoData.Views, &videoData.Likes, &videoData.Dislikes, &videoData.Subscribers, &videoData.Comments)
		if err != nil {
			sendError(w, err.Error())
			return
		}
		response = append(response, videoData)

	}
	fmt.Println(response, "response")
	if err = rows.Err(); err != nil {
		sendError(w, err.Error())
		return
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(result))
	w.Write(result)

}

func GetChats(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var request struct {
		VideoId   string `json:"videoId"`
		NoOfChats *int   `json:"noOfChats"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println("erorr decoding", err)
		sendError(w, err.Error())
		return
	}
	if request.NoOfChats == nil {
		noOfChats := 0
		request.NoOfChats = &noOfChats
	}
	pageSize := 10
	offset := (*request.NoOfChats / pageSize) * pageSize
	fmt.Println("Offset ", offset, " No of chats", *request.NoOfChats)

	rows, err := pool.Query(ctx,
		`SELECT c."userId",c.text,c."createdAt",u.username,u."profileImage"
		FROM "Comment" c 
		JOIN "User" u ON c."userId" = u.id
		 Where c."videoId" = $1
		 ORDER BY c."createdAt" DESC
		 LIMIT $2 OFFSET $3
		`, request.VideoId, pageSize, offset)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Chat struct {
		Message   string      `json:"message"`
		CreatedAt time.Time   `json:"createdAt"`
		User      UserDetails `json:"user"`
	}
	chats := make([]Chat, 0)
	for rows.Next() {
		var chat Chat
		err = rows.Scan(&chat.User.UserId, &chat.Message, &chat.CreatedAt, &chat.User.Username, &chat.User.ProfileImage)
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

	type Response struct {
		Chats  []Chat `json:"chats"`
		Finish bool   `json:"finish"`
	}

	resp := Response{
		Chats: chats,
	}

	if len(chats) < pageSize {
		resp.Finish = true
	}

	result, _ := json.MarshalIndent(resp, "", "  ")
	w.Write(result)
}

func GetChatsFromDatabase(videoId string, numberOfChats int) []Chat {
	ctx := context.Background()

	rows, err := pool.Query(ctx,
		`SELECT c."userId",c.text,c."createdAt",u.username,u."profileImage"
		FROM "Comment" c 
		JOIN "User" u ON c."userId" = u.id
		 Where c."videoId" = $1
		 ORDER BY c."createdAt" DESC
		 LIMIT $2
		`, videoId, numberOfChats)
	if err != nil {

		return []Chat{}
	}
	defer rows.Close()

	chats := make([]Chat, 0)
	for rows.Next() {
		var chat Chat
		err = rows.Scan(&chat.User.UserId, &chat.Message, &chat.CreatedAt, &chat.User.Username, &chat.User.ProfileImage)
		if err != nil {
			return []Chat{}
		}
		chats = append(chats, chat)
	}
	if err = rows.Err(); err != nil {
		return []Chat{}
	}

	return chats
}

func PostChat(videoId string, userId string, message string) error {
    ctx := context.Background()

    // Check if the total number of comments for the video exceeds 100
    var count int
    err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM "Comment" WHERE "videoId" = $1`, videoId).Scan(&count)
    if err != nil {
        return err
    }

    if count >= 100 {
        // Delete the oldest comments to keep the count at 100
        _, err = pool.Exec(ctx, `
            WITH old_comments AS (
                SELECT id
                FROM "Comment"
                WHERE "videoId" = $1
                ORDER BY id ASC
                LIMIT $2
            )
            DELETE FROM "Comment"
            WHERE id IN (SELECT id FROM old_comments)
        `, videoId, count-99)
        if err != nil {
            return err
        }
    }

    // Insert the new comment
    _, err = pool.Exec(ctx, `
        INSERT INTO "Comment" ("videoId", "userId", "text")
        VALUES ($1, $2, $3)
    `, videoId, userId, message)
    if err != nil {
        return err
    }

    return nil
}


type UserDetails struct {
	Username     string  `json:"username"`
	ProfileImage *string `json:"profileImage"`
	UserId       string  `json:"userId"`
}

func GetUserDetailsFromDatabase(userId string) (UserDetails, error) {

	var userDetails UserDetails

	ctx := context.Background()
	err := pool.QueryRow(ctx,
		`SELECT username, "profileImage", id
		 FROM "User"
		 WHERE id = $1`,
		userId).Scan(&userDetails.Username, &userDetails.ProfileImage, &userDetails.UserId)
	if err != nil {
		fmt.Println(err)
		return userDetails, err
	}
	return userDetails, nil

}

func GetCommmentsForCreator(w http.ResponseWriter, r *http.Request) {

	var creatorId = r.Context().Value("userId").(string)
	ctx := context.Background()

	rows, err := pool.Query(ctx,
		`SELECT c."userId",c.text,c."createdAt",c."videoId",u.username,u."profileImage",v.title,v.thumbnail
		 FROM "Comment" c
		 JOIN "User" u ON c."userId" = u.id
		 JOIN "Video" v ON c."videoId" = v.id
		 Where c."videoId" IN (SELECT v.id FROM "Video" v WHERE v."userId" = $1)
		 ORDER BY c."createdAt" DESC
		`, creatorId)

	if err != nil {
		fmt.Println(err, "errror getting comments")
		sendError(w, err.Error())
		return
	}
	defer rows.Close()

	type Chat struct {
		User struct {
			Username     string  `json:"username"`
			ProfileImage *string `json:"profileImage"`
			UserId       string  `json:"userId"`
		} `json:"user"`

		Message   string    `json:"message"`
		CreatedAt time.Time `json:"createdAt"`
		Video     struct {
			VideoId   string `json:"videoId"`
			Thumbnail string `json:"thumbnail"`
			Title     string `json:"title"`
		} `json:"video"`
	}
	chats := make([]Chat, 0)
	for rows.Next() {
		var chat Chat
		err = rows.Scan(&chat.User.UserId, &chat.Message, &chat.CreatedAt, &chat.Video.VideoId, &chat.User.Username, &chat.User.ProfileImage, &chat.Video.Title, &chat.Video.Thumbnail)
		if err != nil {
			fmt.Println(err, "error scanning")
			sendError(w, err.Error())
			return
		}
		chats = append(chats, chat)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err, "error getting rows")

		sendError(w, err.Error())
		return
	}

	result, _ := json.MarshalIndent(chats, "", "  ")
	w.Write(result)

}

func GetUserDetailsByUsername(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	var username string
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		fmt.Println(err.Error())
		sendError(w, "error decoding userId"+err.Error())
		return
	}
	var userDetails struct {
		UserDetails
		Subscribers int `json:"subscribers"`
	}
	ctx := context.Background()
	err = pool.QueryRow(ctx,
		`SELECT "profileImage", id,COALESCE(
			(SELECT COUNT(*) AS subscribers
			 FROM "Subscription"
			 WHERE "creatorId" = (SELECT "id" FROM "User" WHERE username = $1)
			),
			0
		) 
		
		FROM "User"
		WHERE username = $1`,
		username).Scan(&userDetails.ProfileImage, &userDetails.UserId, &userDetails.Subscribers)

	if err != nil {
		fmt.Println(err)
		sendError(w, "No user found with that username")
		return
	}
	userDetails.Username = username
	result, _ := json.MarshalIndent(userDetails, "", "  ")
	w.Write(result)
}

func GetChannelSummary(w http.ResponseWriter, r *http.Request) {

	var response struct {
		Subscribers          int `json:"subscribers"`
		SubscribersLast7Days int `json:"subscribersLast7Days"`
		TotalVideos          int `json:"totalVideos"`
	}
	ctx := context.Background()
	userId := r.Context().Value("userId").(string)
	err := pool.QueryRow(ctx,
		`SELECT 	
		COALESCE(
			(SELECT COUNT(*) AS subscribers
			 FROM "Subscription"
			 WHERE "creatorId" = $1
			),
			0
		) AS subscribers,
		COALESCE(
			(SELECT COUNT(*) AS "subscribersLast7Days"
			 FROM "Subscription"
			 WHERE "creatorId" = $1 AND "createdAt" > NOW() - INTERVAL '7 days'
			),
			0
		) AS "subscribersLast7Days",
		COALESCE(
			(SELECT COUNT(*) AS "totalVideos"
			 FROM "Video"
			 WHERE "userId" = $1
			),
			0
		) AS "totalVideos"
		`, userId).Scan(&response.Subscribers, &response.SubscribersLast7Days, &response.TotalVideos)
	if err != nil {
		sendError(w, err.Error())
		return
	}
	result, _ := json.MarshalIndent(response, "", "  ")
	w.Write(result)

}

func UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	var userDetails struct {
		Username     *string `json:"username"`
		ProfileImage *string `json:"profileImage"`
	}
	err := json.NewDecoder(r.Body).Decode(&userDetails)
	if err != nil {
		fmt.Println(err.Error())
		sendError(w, "error decoding user details")
		return
	}
	userId := r.Context().Value("userId").(string)
	ctx := context.Background()
	_, err = pool.Exec(ctx,
		`UPDATE "User"
		 SET "username" = COALESCE($1, "username"),
		 "profileImage" = COALESCE($2, "profileImage")
		 WHERE id = $3`,
		userDetails.Username, userDetails.ProfileImage, userId)
	if err != nil {
		sendError(w, err.Error())
		return
	}
	var response struct {
		Message string `json:"message"`
	}
	response.Message = "User details updated"

	result, _ := json.MarshalIndent(response, "", "  ")
	w.Write(result)

}

func CheckUsernamePassword(username string, password string) (bool, string, error) {
	ctx := context.Background()
	var hashedPassword, id string

	err := pool.QueryRow(ctx,
		`SELECT p.password, u.id
		FROM "User" u
		JOIN "Password" p ON u.username = p.username
		WHERE u.username = $1;
		  `,
		username).Scan(&hashedPassword, &id)
	if err != nil {
		return false, "", err
	}
	correct := CheckPasswordHash(password, hashedPassword)
	if !correct {
		return false, "", fmt.Errorf("Wrong password")
	}

	return true, id, nil
}

func CreateUserWithPassword(username string, password string) (string, error) {
	ctx := context.Background()
	var id string

	hash, err := HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("error hashing password")
	}
	fmt.Println(hash)

	err = pool.QueryRow(ctx,
		`SELECT "id" FROM "User" WHERE username = $1`,
		username).Scan(&id)

	if err == pgx.ErrNoRows { // User doesn't exist, perform insert
		fmt.Println("User doesn't exist")
		err2 := pool.QueryRow(ctx,
			`INSERT INTO "User" (username)
             VALUES ($1)
             RETURNING "id"`,
			username).Scan(&id)
		if err2 != nil {
			return "", err2
		}

		_, err2 = pool.Exec(ctx,
			`INSERT INTO "Password"(username,password)
			VALUES ($1,$2)
			`, username, hash)

		if err2 != nil {
			return "", err2
		}

		return id, nil

	}
	return "", fmt.Errorf("Already signed Up , try to sign in")

}

func SaveVod(w http.ResponseWriter, r *http.Request) {

	var request struct {
		VideoId     string `json:"videoId"`
		Thumbnail   string `json:"thumbnail"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Visibility  int    `json:"visibility"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		sendError(w, "error decoding videoId")
		return
	}
	ctx := context.Background()
	_, err = pool.Exec(ctx,
		`INSERT INTO "Video" (id, title, description, category, thumbnail, "isStreaming", "userId","isVOD","isProcessed","visibility")
		 VALUES ($1, $2, $3, $4, $5, false, $6,true,false,$7)
		`,
		request.VideoId, request.Title, request.Description, request.Category, request.Thumbnail, r.Context().Value("userId"), request.Visibility)

	if err != nil {
		fmt.Println(err, "error saving video")
		sendError(w, "error saving video")
		return
	}
	err = rmq.PublishMessage(request.VideoId, "vods")
	if err != nil {
		fmt.Println(err, "error publishing message")
		sendError(w, "error publishing message")
		return
	}

	var response struct {
		Message string `json:"message"`
	}
	response.Message = "Saved"
	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		sendError(w, "error marshalling response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func DeleteStreamsByTitle(title string) error {
    ctx := context.Background()

    // Assuming you have a pool variable initialized elsewhere
    _, err := pool.Exec(ctx,
        `DELETE FROM "Video" WHERE title = $1`,
        title,
    )
    if err != nil {
        return err
    }

    return nil
}
