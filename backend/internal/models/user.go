package models

type User struct {
	Id           int    `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	Email        string `json:"email"  db:"email"`
	Password     string `json:"-" db:"password_hash"`
	Wallet       int    `json:"wallet" db:"wallet"`
	Role         string `json:"role" db:"role"`
	Subscription bool   `json:"subscription" db:"subscription"`
	EndOfSub     string `json:"end_of_sub" db:"end_of_sub"`
}

type PlannedManga struct {
	Id      int    `json:"id" db:"id"`
	UserId  int    `json:"user_id" db:"user_id"`
	Site    string `json:"site" db:"site"`
	MangaId string `json:"manga_id" db:"manga_id"`
}

type ArchivedFile struct {
	Id     int    `json:"id" db:"id"`
	UserId int    `json:"user_id" db:"user_id"`
	Name   string `json:"name" db:"name"`
	Url    string `json:"url" db:"url"`
}

type Session struct {
	Id           int    `json:"id" db:"id"`
	UserId       string `json:"user_id" db:"user_id"`
	RefreshToken string `josn:"refresh_token" db:"refresh_token"`
	Fingerprint  string `json:"fingerprint" db:"fingerprint"`
	Ip           string `json:"ip" db:"ip"`
}

type Ads struct {
	Id        int    `json:"id" db:"id"`
	AdType    int    `json:"ad_type" db:"ad_type"`
	BannerUrl string `json:"banner_url" db:"banner_url"`
	Position  int    `json:"position" db:"position"`
}
