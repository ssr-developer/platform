// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"encoding/json"
	"fmt"
	"io"
	"unicode/utf8"
)

const (
	OAUTH_ACTION_SIGNUP       = "signup"
	OAUTH_ACTION_LOGIN        = "login"
	OAUTH_ACTION_EMAIL_TO_SSO = "email_to_sso"
	OAUTH_ACTION_SSO_TO_EMAIL = "sso_to_email"
)

type OAuthApp struct {
	Id           string      `json:"id"`
	CreatorId    string      `json:"creator_id"`
	CreateAt     int64       `json:"create_at"`
	UpdateAt     int64       `json:"update_at"`
	ClientSecret string      `json:"client_secret"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	IconURL      string      `json:"icon_url"`
	CallbackUrls StringArray `json:"callback_urls"`
	Homepage     string      `json:"homepage"`
	IsTrusted    bool        `json:"is_trusted"`
}

// IsValid validates the app and returns an error if it isn't configured
// correctly.
func (a *OAuthApp) IsValid() *AppError {

	if len(a.Id) != 26 {
		return NewLocAppError("OAuthApp.IsValid", "model.oauth.is_valid.app_id.app_error", nil, "")
	}

	if a.CreateAt == 0 {
		return NewLocAppError("OAuthApp.IsValid", "model.oauth.is_valid.create_at.app_error", nil, "app_id="+a.Id)
	}

	if a.UpdateAt == 0 {
		return NewLocAppError("OAuthApp.IsValid", "model.oauth.is_valid.update_at.app_error", nil, "app_id="+a.Id)
	}

	if len(a.CreatorId) != 26 {
		return NewLocAppError("OAuthApp.IsValid", "model.oauth.is_valid.creator_id.app_error", nil, "app_id="+a.Id)
	}

	if len(a.ClientSecret) == 0 || len(a.ClientSecret) > 128 {
		return NewLocAppError("OAuthApp.IsValid", "model.oauth.is_valid.client_secret.app_error", nil, "app_id="+a.Id)
	}

	if len(a.Name) == 0 || len(a.Name) > 64 {
		return NewLocAppError("OAuthApp.IsValid", "model.oauth.is_valid.name.app_error", nil, "app_id="+a.Id)
	}

	if len(a.CallbackUrls) == 0 || len(fmt.Sprintf("%s", a.CallbackUrls)) > 1024 {
		return NewLocAppError("OAuthApp.IsValid", "model.oauth.is_valid.callback.app_error", nil, "app_id="+a.Id)
	}

	for _, callback := range a.CallbackUrls {
		if !IsValidHttpUrl(callback) {
			return NewLocAppError("OAuthApp.IsValid", "model.oauth.is_valid.callback.app_error", nil, "")
		}
	}

	if len(a.Homepage) == 0 || len(a.Homepage) > 256 || !IsValidHttpUrl(a.Homepage) {
		return NewLocAppError("OAuthApp.IsValid", "model.oauth.is_valid.homepage.app_error", nil, "app_id="+a.Id)
	}

	if utf8.RuneCountInString(a.Description) > 512 {
		return NewLocAppError("OAuthApp.IsValid", "model.oauth.is_valid.description.app_error", nil, "app_id="+a.Id)
	}

	if len(a.IconURL) > 0 {
		if len(a.IconURL) > 512 || !IsValidHttpUrl(a.IconURL) {
			return NewLocAppError("OAuthApp.IsValid", "model.oauth.is_valid.icon_url.app_error", nil, "app_id="+a.Id)
		}
	}

	return nil
}

// PreSave will set the Id and ClientSecret if missing.  It will also fill
// in the CreateAt, UpdateAt times. It should be run before saving the app to the db.
func (a *OAuthApp) PreSave() {
	if a.Id == "" {
		a.Id = NewId()
	}

	if a.ClientSecret == "" {
		a.ClientSecret = NewId()
	}

	a.CreateAt = GetMillis()
	a.UpdateAt = a.CreateAt
}

// PreUpdate should be run before updating the app in the db.
func (a *OAuthApp) PreUpdate() {
	a.UpdateAt = GetMillis()
}

// ToJson convert a User to a json string
func (a *OAuthApp) ToJson() string {
	b, err := json.Marshal(a)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

// Generate a valid strong etag so the browser can cache the results
func (a *OAuthApp) Etag() string {
	return Etag(a.Id, a.UpdateAt)
}

// Remove any private data from the app object
func (a *OAuthApp) Sanitize() {
	a.ClientSecret = ""
}

func (a *OAuthApp) IsValidRedirectURL(url string) bool {
	for _, u := range a.CallbackUrls {
		if u == url {
			return true
		}
	}

	return false
}

// OAuthAppFromJson will decode the input and return a User
func OAuthAppFromJson(data io.Reader) *OAuthApp {
	decoder := json.NewDecoder(data)
	var app OAuthApp
	err := decoder.Decode(&app)
	if err == nil {
		return &app
	} else {
		return nil
	}
}

func OAuthAppMapToJson(a map[string]*OAuthApp) string {
	b, err := json.Marshal(a)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func OAuthAppMapFromJson(data io.Reader) map[string]*OAuthApp {
	decoder := json.NewDecoder(data)
	var apps map[string]*OAuthApp
	err := decoder.Decode(&apps)
	if err == nil {
		return apps
	} else {
		return nil
	}
}

func OAuthAppListToJson(l []*OAuthApp) string {
	b, err := json.Marshal(l)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func OAuthAppListFromJson(data io.Reader) []*OAuthApp {
	decoder := json.NewDecoder(data)
	var o []*OAuthApp
	err := decoder.Decode(&o)
	if err == nil {
		return o
	} else {
		return nil
	}
}
