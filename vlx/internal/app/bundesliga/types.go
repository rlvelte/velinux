package bundesliga

import "encoding/json"

// TODO: Maybe use more info later...

type Config struct {
	Team          TeamInfo `json:"team"`
	Notifications bool     `json:"notifications"`
}

type Group struct {
	GroupName    string `json:"groupName"`
	GroupOrderID int    `json:"groupOrderID"`
	GroupID      int    `json:"groupID"`
}

type TeamInfo struct {
	TeamID        int    `json:"teamId"`
	TeamName      string `json:"teamName"`
	ShortName     string `json:"shortName"`
	TeamIconURL   string `json:"teamIconUrl"`
	TeamGroupName any    `json:"teamGroupName"`
}

type MatchResult struct {
	ResultID          int    `json:"resultID"`
	ResultName        string `json:"resultName"`
	PointsTeam1       int    `json:"pointsTeam1"`
	PointsTeam2       int    `json:"pointsTeam2"`
	ResultOrderID     int    `json:"resultOrderID"`
	ResultTypeID      int    `json:"resultTypeID"`
	ResultDescription string `json:"resultDescription"`
}

type Goal struct {
	GoalID         int    `json:"goalID"`
	ScoreTeam1     int    `json:"scoreTeam1"`
	ScoreTeam2     int    `json:"scoreTeam2"`
	MatchMinute    int    `json:"matchMinute"`
	GoalGetterID   int    `json:"goalGetterID"`
	GoalGetterName string `json:"goalGetterName"`
	IsPenalty      bool   `json:"isPenalty"`
	IsOwnGoal      bool   `json:"isOwnGoal"`
	IsOvertime     bool   `json:"isOvertime"`
}

type Match struct {
	MatchID            int           `json:"matchID"`
	MatchDateTime      string        `json:"matchDateTime"`
	TimeZoneID         string        `json:"timeZoneID"`
	LeagueID           int           `json:"leagueId"`
	LeagueName         string        `json:"leagueName"`
	LeagueSeason       int           `json:"leagueSeason"`
	LeagueShortcut     string        `json:"leagueShortcut"`
	MatchDateTimeUTC   string        `json:"matchDateTimeUTC"`
	Group              Group         `json:"group"`
	Team1              TeamInfo      `json:"team1"`
	Team2              TeamInfo      `json:"team2"`
	LastUpdateDateTime string        `json:"lastUpdateDateTime"`
	MatchIsFinished    bool          `json:"matchIsFinished"`
	MatchResults       []MatchResult `json:"matchResults"`
	Goals              []Goal        `json:"goals"`
}

type TableEntry struct {
	TeamInfoID    int    `json:"teamInfoId"`
	TeamName      string `json:"teamName"`
	ShortName     string `json:"shortName"`
	TeamIconURL   string `json:"teamIconURL"`
	Points        int    `json:"points"`
	OpponentGoals int    `json:"opponentGoals"`
	Goals         int    `json:"goals"`
	Matches       int    `json:"matches"`
	Won           int    `json:"won"`
	Lost          int    `json:"lost"`
	Draw          int    `json:"draw"`
	GoalDiff      int    `json:"goalDiff"`
}

func decodeConfig(name, _ string, data []byte) (Config, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
