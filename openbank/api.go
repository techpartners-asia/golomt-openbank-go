package openbank

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	AuthApi = API{
		Url:         "/v1/auth/login",
		Method:      http.MethodPost,
		ServiceName: "LGIN",
	}

	StatementApi = API{
		Url:         "/v1/account/operative/statement/?client_id={{clientId}}&state={{state}}&scope={{scope}}",
		Method:      http.MethodPost,
		ServiceName: "OPERACCTSTA",
	}
	AccountListApi = API{
		Url:         "/v1/account/list?client_id={{clientId}}&state={{state}}&scope={{scope}}",
		Method:      http.MethodPost,
		ServiceName: "ACCTLST",
	}
	AccountTypeInq = API{
		Url:         "/v1/account/type/inq",
		Method:      http.MethodPost,
		ServiceName: "ACCTTYPEINQ",
	}

	AccountBalcApi = API{
		Url:         "/v1/account/balance/inq?client_id={{clientId}}&state={{state}}&scope={{scope}}",
		Method:      http.MethodPost,
		ServiceName: "ACCTBALINQ",
	}

	UtilityRateAPI = API{
		Url:         "/v1/utility/rate/inq",
		Method:      http.MethodPost,
		ServiceName: "RATEINQ",
	}
)

// AuthQPayV2 [Login to qpay]
func (g *openbank) auth() (authRes AuthResp, err error) {
	if g.authObject != nil {
		if g.expireTime.Before(time.Now()) {
			authRes = *g.authObject
			return
		} else {
			url := g.url + "/v1/auth/refresh"
			req, reqErr := http.NewRequest(http.MethodGet, url, nil)
			if reqErr != nil {
				fmt.Println(err.Error())
				err = reqErr
				return
			}
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("X-Golomt-Service", AuthApi.ServiceName)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return authRes, err
			}

			body, _ := io.ReadAll(res.Body)
			if res.StatusCode != http.StatusOK {
				var errResp ErrorResp
				json.Unmarshal(body, &errResp)
				fmt.Println(errResp)
				return authRes, fmt.Errorf("%s-Golomt CG auth response: %s", time.Now().Format("20060102150405"), errResp.Message)
			}
			json.Unmarshal(body, &authRes)
			fmt.Println("-----------------------------Login response--------------------------------------------------------")
			fmt.Println(authRes)
			defer res.Body.Close()
			return authRes, nil
		}
	}
	url := g.url + AuthApi.Url
	reqBody := AuthReq{
		Name: g.username,
		Password: func() string {
			pass, err := g.EncryptAESCBC(g.password)
			if err != nil {
				return ""
			}
			return pass
		}(),
	}

	if reqBody.Password == "" {
		panic("AES error")
	}
	requestByte, _ := json.Marshal(reqBody)
	requestBody := bytes.NewReader(requestByte)
	req, reqErr := http.NewRequest(AuthApi.Method, url, requestBody)
	if reqErr != nil {
		fmt.Println(err.Error())
		err = reqErr
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Golomt-Service", AuthApi.ServiceName)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		var errResp ErrorResp
		json.Unmarshal(body, &errResp)
		fmt.Println(errResp)
		return authRes, fmt.Errorf("%s-Golomt CG auth response: %s", time.Now().Format("20060102150405"), errResp.Message)
	}
	json.Unmarshal(body, &authRes)
	fmt.Println("-----------------------------Login response--------------------------------------------------------")
	fmt.Println(authRes)
	defer res.Body.Close()
	return authRes, nil
}

func (g *openbank) HttpRequest(body interface{}, api API, urlExt string) (response string, err error) {
	authObj, authErr := g.auth()
	if authErr != nil {
		err = authErr
		return
	}

	g.authObject = &authObj
	postBody, _ := json.Marshal(body)
	fmt.Println("----------------------body-----------------------")
	fmt.Println(body)
	hash := sha256.Sum256(postBody)
	hex := hex.EncodeToString(hash[:])
	// checkSum, _ := g.encryptAESCBC(string(hash))
	checkSum, err := g.EncryptAESCBC(hex)
	if err != nil {
		return "", errors.New(err.Error())
	}
	var requestByte []byte
	var requestBody *bytes.Reader
	if body == nil {
		requestBody = bytes.NewReader(nil)
	} else {
		requestByte, _ = json.Marshal(body)
		requestBody = bytes.NewReader(requestByte)
	}
	if api.ServiceName == "OPERACCTSTA" || api.ServiceName == "ACCTLST" || api.ServiceName == "ACCTBALINQ" {
		api.Url = strings.Replace(api.Url, "{{clientId}}", g.clientID, 1)
		api.Url = strings.Replace(api.Url, "{{state}}", g.state, 1)
		api.Url = strings.Replace(api.Url, "{{scope}}", g.scope, 1)
	}
	req, _ := http.NewRequest(api.Method, g.url+api.Url+urlExt, requestBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Golomt-Checksum", checkSum)
	req.Header.Add("X-Golomt-Service", api.ServiceName)
	req.Header.Add("Authorization", "Bearer "+g.authObject.Token)
	fmt.Println("-----------------------------Request--------------------------------------------------------")
	fmt.Println(req)
	fmt.Println("-----------------------------Request body--------------------------------------------------------")
	fmt.Println(req.Body)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.New(err.Error())
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.New(err.Error())
	}
	if res.Status != "200 OK" {
		response, err = g.DecryptAESCBC(string(resBody))
		return "", errors.New(string(response))
	}

	response, err = g.DecryptAESCBC(string(resBody))
	if err != nil {
		return "", errors.New(err.Error())
	}

	return response, nil
}

func (g *openbank) HttpRequestSimple(body interface{}, api API, urlExt string) (response string, err error) {
	authObj, authErr := g.auth()
	if authErr != nil {
		err = authErr
		return
	}

	g.authObject = &authObj
	postBody, _ := json.Marshal(body)
	fmt.Println("----------------------body-----------------------")
	fmt.Println(body)
	hash := sha256.Sum256(postBody)
	hex := hex.EncodeToString(hash[:])
	// checkSum, _ := g.encryptAESCBC(string(hash))
	checkSum, err := g.EncryptAESCBC(hex)
	if err != nil {
		return "", errors.New(err.Error())
	}
	var requestByte []byte
	var requestBody *bytes.Reader
	if body == nil {
		requestBody = bytes.NewReader(nil)
	} else {
		requestByte, _ = json.Marshal(body)
		requestBody = bytes.NewReader(requestByte)
	}
	if api.ServiceName == "OPERACCTSTA" || api.ServiceName == "ACCTLST" || api.ServiceName == "ACCTBALINQ" {
		api.Url = strings.Replace(api.Url, "{{clientId}}", g.clientID, 1)
		api.Url = strings.Replace(api.Url, "{{state}}", g.state, 1)
		api.Url = strings.Replace(api.Url, "{{scope}}", g.scope, 1)
	}
	req, _ := http.NewRequest(api.Method, g.url+api.Url+urlExt, requestBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Golomt-Checksum", checkSum)
	req.Header.Add("X-Golomt-Service", api.ServiceName)
	req.Header.Add("Authorization", "Bearer "+g.authObject.Token)
	fmt.Println("-----------------------------Request--------------------------------------------------------")
	fmt.Println(req)
	fmt.Println("-----------------------------Request body--------------------------------------------------------")
	fmt.Println(req.Body)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.New(err.Error())
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.New(err.Error())
	}
	if res.Status != "200 OK" {
		response, err = g.DecryptAESCBC(string(resBody))
		return "", errors.New(string(response))
	}

	// response, err = g.DecryptAESCBC(string(resBody))
	// if err != nil {
	// 	return "", errors.New(err.Error())
	// }

	return string(resBody), nil
}

func (g *openbank) HttpOAuthRequest(body interface{}, api API, urlExt string) (OAuthReponse, error) {
	authObj, authErr := g.auth()
	if authErr != nil {
		err := authErr
		return OAuthReponse{}, err
	}

	g.authObject = &authObj
	postBody, _ := json.Marshal(body)
	fmt.Println("----------------------body-----------------------")
	fmt.Println(body)
	hash := sha256.Sum256(postBody)
	hex := hex.EncodeToString(hash[:])
	// checkSum, _ := g.encryptAESCBC(string(hash))
	checkSum, err := g.EncryptAESCBC(hex)
	if err != nil {
		return OAuthReponse{}, errors.New(err.Error())
	}
	var requestByte []byte
	var requestBody *bytes.Reader
	if body == nil {
		requestBody = bytes.NewReader(nil)
	} else {
		requestByte, _ = json.Marshal(body)
		requestBody = bytes.NewReader(requestByte)
	}
	if api.ServiceName == "OPERACCTSTA" || api.ServiceName == "ACCTLST" {
		// api.Url = strings.Replace(api.Url, "{{clientId}}", g.ClientID, 1)
		api.Url = strings.Replace(api.Url, "{{clientId}}", g.clientID, 1)
		api.Url = strings.Replace(api.Url, "{{state}}", g.state, 1)
		api.Url = strings.Replace(api.Url, "{{scope}}", g.scope, 1)
	}
	req, _ := http.NewRequest(api.Method, g.url+api.Url+urlExt, requestBody)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Golomt-Checksum", checkSum)
	req.Header.Add("X-Golomt-Service", api.ServiceName)
	req.Header.Add("Authorization", "Bearer "+g.authObject.Token)
	fmt.Println("-----------------------------Request--------------------------------------------------------")
	fmt.Println(req)
	fmt.Println("-----------------------------Request body--------------------------------------------------------")
	fmt.Println(req.Body)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return OAuthReponse{}, errors.New(err.Error())
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return OAuthReponse{}, errors.New(err.Error())
	}
	if res.Status != "200 OK" {
		response, err := g.DecryptAESCBC(string(resBody))
		if err != nil {
			return OAuthReponse{}, errors.New(err.Error())
		}
		return OAuthReponse{}, errors.New(string(response))
	}

	response, err := g.DecryptAESCBC(string(resBody))
	if err != nil {
		return OAuthReponse{}, errors.New(err.Error())
	}
	var resp OAuthReponse
	json.Unmarshal([]byte(response), &resp)
	return resp, nil
}
