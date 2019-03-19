package route

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
	"welcome_robot/dao"
	. "welcome_robot/models"
)

type DetailSession struct {
	SessionId   bson.ObjectId `bson:"session_id"`
	Supporter   User          `bson:"supporter"`
	User        User          `bson:"user"`
	CheckInTime time.Time     `bson:"check_in_time"`
}

func InsertSessions(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var session Session
	var user User
	var userInfor UserInfo
	if err := json.NewDecoder(r.Body).Decode(&userInfor); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	session.SessionID = bson.NewObjectId()
	if userInfor.Permission != 2 {
		session.CheckInTime = time.Now()
		//respondWithJson(w, http.StatusCreated, user2)
	}
	user, err := dao.InsertUser(userInfor)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//respondWithJson(w, http.StatusCreated, user)
	session.SupporterID = getSupporterId()
	session.UserID = user.UserID
	if err := dao.InsertSession(session); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, session)
}
func getSupporterId() bson.ObjectId {
	return bson.ObjectIdHex("5c8f626431ce9701e81c10a1")
}
func GetAllSession(w http.ResponseWriter, r *http.Request) {
	users, err := dao.GetAllUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, users)
}

func GetDetailSession(w http.ResponseWriter, r *http.Request) {
	var supporter User
	var user User
	var session Session
	var detailSession DetailSession
	params := mux.Vars(r)
	session, err := dao.GetSessionByUserID(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	supporter, err = dao.FindUserById(session.SupporterID.Hex())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	user, err = dao.FindUserById(session.UserID.Hex())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	detailSession.SessionId = session.SessionID
	detailSession.Supporter = supporter
	detailSession.User = user
	detailSession.CheckInTime = session.CheckInTime
	log.Print(detailSession)
	respondWithJson(w, http.StatusOK, detailSession)
}
func RemoveSessions(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	var session Session
	session, err := dao.GetSessionById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	_ = dao.RemoveUser(session.UserID.Hex())
	_ = dao.DeleteVideoTimeBySessionId(session.SessionID.Hex())
	if err := dao.RemoveSession(params["id"]); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
