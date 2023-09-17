package database

import (
	"errors"
	"twitter_task/model"

	"gorm.io/gorm"
)

func SaveAccessToken(address string, twitterUserId string, twitterUsername string, oauthToken string, oauthSecret string) (*model.User, error) {

	// if user have twitter binded
	var user model.User
	userResult := db.Model(&user).Where("address = ?", address).First(&user)
	if userResult.Error != nil && !errors.Is(userResult.Error, gorm.ErrRecordNotFound) {
		return nil, userResult.Error
	}
	if userResult.RowsAffected > 0 && user.TwitterID > 0 {
		return nil, UserAddressAlreadyBindWithTwitter
	}

	// if twitter already exists and binded with another user
	var twitter model.Twitter
	twitterResult := db.Model(&twitter).Where("tw_user_id = ?", twitterUserId).First(&twitter)
	if twitterResult.Error != nil && !errors.Is(twitterResult.Error, gorm.ErrRecordNotFound) {
		return nil, twitterResult.Error
	}

	// twitter exists
	if twitterResult.RowsAffected > 0 {
		twitter.OAuthToken = oauthToken
		twitter.OAuthSecret = oauthSecret
		twitter.TWUsername = twitterUsername
	} else {
		// twitter not exists
		twitter.OAuthToken = oauthToken
		twitter.OAuthSecret = oauthSecret
		twitter.TWUsername = twitterUsername
		twitter.TWUserID = twitterUserId
	}

	if err := db.Save(&twitter).Error; err != nil {
		return nil, err
	}

	if rowsAffected := db.Model(&model.User{}).Where("twitter_id = ?", twitter.ID).RowsAffected; rowsAffected > 0 {
		// twitter already binded an address
		return nil, UserAddressAlreadyBindWithTwitter
	}
	// user exists
	if userResult.RowsAffected > 0 {
		user.TwitterID = twitter.ID
		user.Twitter = twitter
		if user.Username == "" {
			user.Username = twitter.TWUsername
		}
		if err := db.Save(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}
	// user not exists, create user
	user.Address = address
	user.TwitterID = twitter.ID
	user.Twitter = twitter
	user.Username = twitter.TWUsername
	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserWithAddress(address string) (*model.User, error) {
	var user model.User
	if err := db.Model(&user).Preload("Twitter").Where("address = ?", address).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
