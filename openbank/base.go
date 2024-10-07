package openbank

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type openbank struct {
	organizationName string
	username         string
	password         string
	ivKey            string
	sessionKey       string
	url              string
	registerNo       string
	expireTime       time.Time
	authObject       *AuthResp
	clientID         string
	state            string
	scope            string
}

type Openbank interface {
	Statement(body StatementReq) (StatementResp, error)
	AccountList(body AccountListReq) (AccountListResp, error)
	AccountTypeInq(body AccountTypeInqReq) (AccountTypeInqResp, error)
	AccountBalcInq(body AccountBalcInqReq) (AccountBalcInqResp, error)
}

func New(username, password, client_id, orgname, ivKey, sessoinKey, url, registerNo string) Openbank {
	return &openbank{
		organizationName: orgname,
		username:         username,
		password:         password,
		ivKey:            ivKey,
		sessionKey:       sessoinKey,
		url:              url,
		authObject:       nil,
		registerNo:       registerNo,
		clientID:         client_id,
	}
}

// func (golomt *golomtbankCG) Login() (AuthResp, error) {
// 	resp, err := golomt.auth()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return AuthResp{}, err
// 	}
// 	return resp, err
// }

// func (golomt *golomtbankCG) ServiceList(code string, services []string) (StatementResp, error) {
// 	body := ServiceListReq{
// 		RegisterNo: golomt.RegisterNo,
// 		Services:   services,
// 		Code:       code,
// 	}

// 	resp, err := golomt.HttpRequest(body, ServiceListApi, "")

// 	if err != nil {
// 		fmt.Println(err.Error())
// 		var errorBody ErrorResp
// 		err = json.Unmarshal([]byte(resp), &errorBody)
// 		return errorBody, err
// 	}
// 	var response StatementResp
// 	err = json.Unmarshal([]byte(resp), &response)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return StatementResp{}, err
// 	}
// 	return response, err
// }

func (golomt *openbank) Statement(body StatementReq) (StatementResp, error) {

	// body := StatementReq{
	// 	AccountID:  account,
	// 	RegisterNo: golomt.RegisterNo,
	// 	StartDate:  startdate.Format("2006-01-02"),
	// 	EndDate:    endate.Format("2006-01-02"),
	// }
	// resp, err := golomt.HttpOAuthRequest(body, StatementApi, "")
	// if err != nil {
	// 	return nil, err
	// }
	// golomt.ClientID = resp.ClientID
	// golomt.State = resp.State
	// golomt.Scope = resp.Scope

	res, err := golomt.HttpRequest(body, StatementApi, "")
	if err != nil {
		var errorBody ErrorResp
		err = json.Unmarshal([]byte(err.Error()), &errorBody)

		return StatementResp{}, errors.New(errorBody.Message + " :" + errorBody.DebugMessage)
	}
	var response StatementResp
	err = json.Unmarshal([]byte(res), &response)
	if err != nil {
		fmt.Println(err.Error())
		return StatementResp{}, err
	}
	return response, err
	// var response StatementResp
	// err = json.Unmarshal([]byte(resp), &response)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return StatementResp{}, err
	// }
	// return response, err
}

func (golomt *openbank) AccountList(body AccountListReq) (AccountListResp, error) {

	// body := AccountListReq{
	// 	RegisterNo: registerNo,
	// }
	// resp, err := golomt.HttpOAuthRequest(body, AccountListApi, "")
	// if err != nil {
	// 	return nil, err
	// }
	// if golomt.ClientID == "" && golomt.State == "" && golomt.Scope == "" {
	// var responseOath Response
	// err = json.Unmarshal([]byte(resp), &responseOath)
	// golomt.ClientID = resp.ClientID
	// golomt.State = resp.State
	// golomt.Scope = resp.Scope

	res, err := golomt.HttpRequest(body, AccountListApi, "")
	if err != nil {
		var errorBody ErrorResp
		err = json.Unmarshal([]byte(err.Error()), &errorBody)
		return AccountListResp{}, errors.New(errorBody.Message + " :" + errorBody.DebugMessage)
	}
	var response AccountListResp
	err = json.Unmarshal([]byte(res), &response)
	if err != nil {
		fmt.Println(err.Error())
		return AccountListResp{}, err
	}
	// }
	// var response AccountListResp
	// err = json.Unmarshal([]byte(res), &response)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return AccountListResp{}, err
	// }
	return response, err
}

func (golomt *openbank) AccountTypeInq(body AccountTypeInqReq) (AccountTypeInqResp, error) {

	// body := AccountTypeInqReq{
	// 	AccountID: accountId,
	// }
	resp, err := golomt.HttpRequest(body, AccountTypeInq, "")

	if err != nil {
		var errorBody ErrorResp
		err = json.Unmarshal([]byte(resp), &errorBody)
		return AccountTypeInqResp{}, errors.New(errorBody.Message + " :" + errorBody.DebugMessage)
	}
	var response AccountTypeInqResp
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		fmt.Println(err.Error())
		return AccountTypeInqResp{}, err
	}
	return response, err
}

func (golomt *openbank) AccountBalcInq(body AccountBalcInqReq) (AccountBalcInqResp, error) {

	resp, err := golomt.HttpRequest(body, AccountBalcApi, "")

	if err != nil {
		var errorBody ErrorResp
		err = json.Unmarshal([]byte(resp), &errorBody)
		return AccountBalcInqResp{}, errors.New(errorBody.Message + " :" + errorBody.DebugMessage)
	}
	var response AccountBalcInqResp
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		fmt.Println(err.Error())
		return AccountBalcInqResp{}, err
	}
	return response, err
}
