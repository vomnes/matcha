package user

import (
	"log"

	"../../../../lib"
	"github.com/jmoiron/sqlx"
)

type link struct {
	Count       uint8
	LikedByUser bool
}

func getNumberLikesConnections(db *sqlx.DB, userID string) (float64, float64, int, string) {
	var likes []lib.Like
	err := db.Select(&likes, `Select userid, liked_userid From Likes Where userid = $1 OR liked_userid = $1`, userID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect likes in database " + err.Error()))
		return 0, 0, 500, "Failed to collect likes in the database"
	}
	links := make(map[string]link)
	for _, like := range likes {
		if like.UserID == userID {
			links[like.LikedUserID] = link{
				Count:       links[like.LikedUserID].Count + 1,
				LikedByUser: true,
			}
		} else if like.LikedUserID == userID {
			links[like.UserID] = link{
				Count:       links[like.UserID].Count + 1,
				LikedByUser: links[like.LikedUserID].LikedByUser,
			}
		}
	}
	var nbLikes, nbConnections float64
	for _, link := range links {
		if link.LikedByUser == true {
			if link.Count >= 2 {
				nbConnections++
				nbLikes++
			}
		} else {
			nbLikes++
		}
	}
	return nbLikes, nbConnections, 0, ""
}

func getNumberVisits(db *sqlx.DB, userID string) (float64, int, string) {
	var visits []lib.Visit
	err := db.Select(&visits, `Select Distinct on (userid) userid, visited_userid From Visits Where visited_userid = $1`, userID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect visits in database " + err.Error()))
		return 0, 500, "Failed to collect visits in the database"
	}
	return float64(len(visits)), 0, ""
}

func getNumberFakeReports(db *sqlx.DB, userID string) (float64, int, string) {
	var reports []lib.FakeReport
	err := db.Select(&reports, `Select id From Fake_reports Where target_userid = $1`, userID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - SELECT] Failed to collect fake reports in database " + err.Error()))
		return 0, 500, "Failed to collect fake reports in the database"
	}
	return float64(len(reports)), 0, ""
}

func updateRating(db *sqlx.DB, userID string) (int, string) {
	var nbLikes, nbConnection, nbVisits, nbFakeReports, rating float64
	nbLikes, nbConnection, errCode, errContent := getNumberLikesConnections(db, userID)
	if errCode != 0 || errContent != "" {
		return errCode, errContent
	}
	nbVisits, errCode, errContent = getNumberVisits(db, userID)
	if errCode != 0 || errContent != "" {
		return errCode, errContent
	}
	nbFakeReports, errCode, errContent = getNumberFakeReports(db, userID)
	if errCode != 0 || errContent != "" {
		return errCode, errContent
	}
	if nbVisits == 0 {
		nbVisits = 1
	}
	rating = 2.0*(nbLikes*(0.9+(nbConnection/100.0))/nbVisits) + 3.0*(1.0-(nbFakeReports*5)/100)
	if rating > 5 {
		rating = 5
	} else if rating < 0 {
		rating = 0
	}
	updateRequest := `UPDATE users SET
    rating = $1
    WHERE  users.id = $2`
	_, err := db.Queryx(updateRequest, rating, userID)
	if err != nil {
		log.Println(lib.PrettyError("[DB REQUEST - Update] Failed to update User[" + userID + "] Rating Data " + err.Error()))
		return 500, "Failed to update rating data in database"
	}
	return 0, ""
}
